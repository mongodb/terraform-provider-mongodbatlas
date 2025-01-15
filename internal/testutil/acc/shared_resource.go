package acc

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	MaxClustersPerProject = 6
)

// SetupSharedResources must be called from TestMain test package in order to use ProjectIDExecution.
// It returns the cleanup function that must be called at the end of TestMain.
func SetupSharedResources() func() {
	sharedInfo.init = true
	return cleanupSharedResources
}

func cleanupSharedResources() {
	if sharedInfo.projectID != "" && sharedInfo.clusterName != "" {
		fmt.Printf("Deleting execution cluster: %s, project id: %s\n", sharedInfo.clusterName, sharedInfo.projectID)
		deleteCluster(sharedInfo.projectID, sharedInfo.clusterName)
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

	if id := projectIDLocal(tb); id != "" {
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

// ProjectIDExecutionWithCluster creates a project and reuses it `MaxClustersPerProject` times. The clusterName is always unique.
// This avoids the `CROSS_REGION_NETWORK_PERMISSIONS_LIMIT_EXCEEDED` error when creating too many clusters within the same project.
// When `MONGODB_ATLAS_PROJECT_ID` and `MONGODB_ATLAS_CLUSTER_NAME` are defined, they are used instead of creating a project and clusterName.
func ProjectIDExecutionWithCluster(tb testing.TB) (projectID, clusterName string) {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")

	if ExistingClusterUsed() {
		return existingProjectIDClusterName()
	}
	return NextProjectIDClusterName(func(projectName string) string {
		return createProject(tb, projectName)
	})
}

// ClusterNameExecution returns the name of a created cluster for the execution of the tests in the resource package.
// This function relies on using an execution project and returns its id.
// When `MONGODB_ATLAS_CLUSTER_NAME` and `MONGODB_ATLAS_PROJECT_ID` are defined it will be used instead of creating resources. This is useful for local execution but not intended for CI executions.
func ClusterNameExecution(tb testing.TB) (projectID, clusterName string) {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")

	if ExistingClusterUsed() {
		return existingProjectIDClusterName()
	}

	// before locking for cluster creation we need to ensure we have an execution project created
	if sharedInfo.projectID == "" {
		_ = ProjectIDExecution(tb)
	}

	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()

	// lazy creation so it's only done if really needed
	if sharedInfo.clusterName == "" {
		name := RandomClusterName()
		tb.Logf("Creating execution cluster: %s\n", name)
		sharedInfo.clusterName = createCluster(tb, sharedInfo.projectID, name)
	}

	return sharedInfo.projectID, sharedInfo.clusterName
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
	id           string
	name         string
	clusterCount int
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
func NextProjectIDClusterName(projectCreator func(string) string) (projectID, clusterName string) {
	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()
	var project projectInfo
	if len(sharedInfo.projects) == 0 || sharedInfo.projects[len(sharedInfo.projects)-1].clusterCount == MaxClustersPerProject {
		project = projectInfo{
			name:         RandomProjectName(),
			clusterCount: 1,
		}
		project.id = projectCreator(project.name)
		sharedInfo.projects = append(sharedInfo.projects, project)
	} else {
		project = sharedInfo.projects[len(sharedInfo.projects)-1]
		sharedInfo.projects[len(sharedInfo.projects)-1].clusterCount++
	}
	return project.id, RandomClusterName()
}
