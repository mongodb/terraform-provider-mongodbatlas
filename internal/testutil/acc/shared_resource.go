package acc

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/require"
)

const (
	MaxClusterNodesPerProject = 30 // Choose to be conservative, 40 clusters per project is the limit before `CROSS_REGION_NETWORK_PERMISSIONS_LIMIT_EXCEEDED` error, see https://www.mongodb.com/docs/atlas/reference/atlas-limits/
)

// SetupSharedResources must be called from TestMain test package in order to use ProjectIDExecution.
// It returns the cleanup function that must be called at the end of TestMain.
func SetupSharedResources() func() {
	sharedInfo.init = true
	setupTestsSDKv2ToTPF()
	return cleanupSharedResources
}

func cleanupSharedResources() {
	if sharedInfo.clusterName != "" {
		projectID := sharedInfo.projectID
		if projectID == "" {
			projectID = projectIDLocal()
		}
		fmt.Printf("Deleting execution cluster: %s, project id: %s\n", sharedInfo.clusterName, projectID)
		deleteCluster(projectID, sharedInfo.clusterName)
	}

	if sharedInfo.projectID != "" {
		fmt.Printf("Deleting execution project: %s, id: %s\n", sharedInfo.projectName, sharedInfo.projectID)
		deleteProject(sharedInfo.projectID)
	}
	for i, project := range sharedInfo.projects {
		fmt.Printf("Deleting execution project (%d): %s, id: %s\n", i+1, project.name, project.id)
		deleteProject(project.id)
	}
}

// ProjectIDExecution returns a project id created for the execution of the tests in the resource package.
// Even if a GH test group is run, every resource/package will create its own project, not a shared project for all the test group.
// When `MONGODB_ATLAS_PROJECT_ID` is defined, it is used instead of creating a project. This is useful for local execution but not intended for CI executions.
func ProjectIDExecution(tb testing.TB) string {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")

	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()

	if id := projectIDLocal(); id != "" {
		return id
	}

	// lazy creation so it's only done if really needed
	if sharedInfo.projectID == "" {
		sharedInfo.projectName = RandomProjectName()
		tb.Logf("Creating execution project: %s\n", sharedInfo.projectName)
		sharedInfo.projectID = createProject(tb, sharedInfo.projectName)
	}

	return sharedInfo.projectID
}

// ProjectIDExecutionWithCluster creates a project and reuses it for  `MaxClusterNodesPerProject ` nodes. The clusterName is always unique.
// TotalNodeCount = sum(specs.node_count) * num_shards (1 if new schema)
// This avoids the `CROSS_REGION_NETWORK_PERMISSIONS_LIMIT_EXCEEDED` error when creating too many clusters within the same project.
// When `MONGODB_ATLAS_PROJECT_ID` and `MONGODB_ATLAS_CLUSTER_NAME` are defined, they are used instead of creating a project and clusterName.
func ProjectIDExecutionWithCluster(tb testing.TB, totalNodeCount int) (projectID, clusterName string) {
	tb.Helper()
	if ExistingClusterUsed() {
		return existingProjectIDClusterName()
	}
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")
	return NextProjectIDClusterName(totalNodeCount, func(projectName string) string {
		return createProject(tb, projectName)
	})
}

// ClusterNameExecution returns the name of a created cluster for the execution of the tests in the resource package.
// This function relies on using an execution project and returns its id.
// When `MONGODB_ATLAS_CLUSTER_NAME` and `MONGODB_ATLAS_PROJECT_ID` are defined it will be used instead of creating resources. This is useful for local execution but not intended for CI executions.
func ClusterNameExecution(tb testing.TB, populateSampleData bool) (projectID, clusterName string) {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")

	if ExistingClusterUsed() {
		return existingProjectIDClusterName()
	}

	projectID = sharedInfo.projectID
	if projectID == "" {
		projectID = ProjectIDExecution(tb) // ensure the execution project is created before cluster creation
	}

	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()

	// lazy creation so it's only done if really needed
	if sharedInfo.clusterName == "" {
		name := RandomClusterName()
		tb.Logf("Creating execution cluster: %s\n", name)
		sharedInfo.clusterName = createCluster(tb, projectID, name)

		if populateSampleData {
			err := PopulateWithSampleData(projectID, sharedInfo.clusterName)
			require.NoError(tb, err)
		}
	}

	return projectID, sharedInfo.clusterName
}

// SerialSleep waits a few seconds so clusters in a project are not created concurrently, see HELP-65223.
// This must be called once the test is marked as parallel, e.g. in PreCheck inside Terraform tests.
func SerialSleep(tb testing.TB) {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")

	sharedInfo.muSleep.Lock()
	defer sharedInfo.muSleep.Unlock()

	time.Sleep(5 * time.Second)
}

type projectInfo struct {
	id        string
	name      string
	nodeCount int
}

var sharedInfo = struct {
	projectID   string
	projectName string
	clusterName string
	projects    []projectInfo
	mu          sync.Mutex
	muSleep     sync.Mutex
	init        bool
}{
	projects: []projectInfo{},
}

// NextProjectIDClusterName is an internal method used when we want to reuse a projectID `MaxClustersPerProject` times
func NextProjectIDClusterName(totalNodeCount int, projectCreator func(string) string) (projectID, clusterName string) {
	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()
	var project projectInfo
	if len(sharedInfo.projects) == 0 || sharedInfo.projects[len(sharedInfo.projects)-1].nodeCount+totalNodeCount > MaxClusterNodesPerProject {
		project = projectInfo{
			name:      RandomProjectName(),
			nodeCount: totalNodeCount,
		}
		project.id = projectCreator(project.name)
		sharedInfo.projects = append(sharedInfo.projects, project)
	} else {
		project = sharedInfo.projects[len(sharedInfo.projects)-1]
		sharedInfo.projects[len(sharedInfo.projects)-1].nodeCount += totalNodeCount
	}
	return project.id, RandomClusterName()
}

// setupTestsSDKv2ToTPF sets the Preview environment variable to false so the previous version in migration tests uses SDKv2.
// However the current version will use TPF as the variable is only read once during import when it was true.
func setupTestsSDKv2ToTPF() {
	if IsTestSDKv2ToTPF() && config.PreviewProviderV2AdvancedCluster() {
		os.Setenv(config.PreviewProviderV2AdvancedClusterEnvVar, "false")
	}
}
