package acc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/clean"
	"github.com/stretchr/testify/require"
)

const (
	MaxClusterNodesPerProject = 30 // Choose to be conservative, 40 clusters per project is the limit before `CROSS_REGION_NETWORK_PERMISSIONS_LIMIT_EXCEEDED` error, see https://www.mongodb.com/docs/atlas/reference/atlas-limits/
	MaxFreeTierClusterCount   = 1  // Project can have at most 1 free tier cluster
)

type projectInfo struct {
	id                   string
	name                 string
	nodeCount            int
	freeTierClusterCount int
}

var sharedInfo = struct {
	clusterName        string
	streamInstanceName string
	projects           []projectInfo
	mu                 sync.Mutex
	muSleep            sync.Mutex
	init               bool
}{
	projects: []projectInfo{},
}

// Run handles the common logic for running acceptance tests: Init shared struct, run tests and clean up shared resources.
// It returns an exit code to pass to os.Exit. The exit code is zero when all tests pass and clean up succeeds, and non-zero for any kind of failure.
func Run(m *testing.M) (code int) {
	sharedInfo.init = true
	exitCode := m.Run()
	if err := cleanupSharedResources(); err != nil {
		log.Printf("[ERROR] Cleanup failed: %v", err)
		exitCode = 1
	}
	return exitCode
}

func cleanupSharedResources() error {
	var hasError bool
	firstProjectID := projectIDLocal()
	if firstProjectID == "" && len(sharedInfo.projects) > 0 {
		firstProjectID = sharedInfo.projects[0].id
	}
	if sharedInfo.clusterName != "" {
		fmt.Printf("Deleting execution cluster: %s, project id: %s\n", sharedInfo.clusterName, firstProjectID)
		if err := deleteCluster(firstProjectID, sharedInfo.clusterName); err != nil {
			fmt.Printf("[ERROR] Cluster deletion failed: %v\n", err)
			hasError = true
		}
	}
	if sharedInfo.streamInstanceName != "" {
		fmt.Printf("Deleting execution stream instance: %s, project id: %s\n", sharedInfo.streamInstanceName, firstProjectID)
		if _, err := clean.RemoveStreamInstances(context.TODO(), false, ConnV2(), firstProjectID); err != nil {
			fmt.Printf("[ERROR] Stream instance deletion failed: %v\n", err)
			hasError = true
		}
	}
	for i, project := range sharedInfo.projects {
		fmt.Printf("Deleting execution project (%d): %s, id: %s\n", i+1, project.name, project.id)
		if err := deleteProject(project.id); err != nil {
			fmt.Printf("[ERROR] Project deletion failed: %v\n", err)
			hasError = true
		}
	}
	config.CloseTokenSource() // Revoke SA token when acceptance tests finish.
	if hasError {
		return errors.New("failed to delete shared resources")
	}
	return nil
}

// ProjectIDExecution returns a project id created for the execution of the tests in the resource package.
// Even if a GH test group is run, every resource/package will create its own project, not a shared project for all the test group.
// When `MONGODB_ATLAS_PROJECT_ID` is defined, it is used instead of creating a project. This is useful for local execution but not intended for CI executions.
func ProjectIDExecution(tb testing.TB) string {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "sharedInfo not initialized, use acc.Run() to run tests that require shared resources")

	if id := projectIDLocal(); id != "" {
		return id
	}

	createSharedProjects(tb, 1)
	return sharedInfo.projects[0].id
}

// MultipleProjectIDsExecution returns multiple project ids created for test execution in the resource package.
// Even if a GH test group is run, every resource/package will create its own projects, not shared projects for all the test group.
// Panics when `MONGODB_ATLAS_PROJECT_ID` is defined and more than 1 project is requested.
func MultipleProjectIDsExecution(tb testing.TB, count int) []string {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "sharedInfo not initialized, use acc.Run() to run tests that require shared resources")
	require.Positive(tb, count, "count must be greater than 0")

	if id := projectIDLocal(); id != "" {
		projectIDs := []string{id}
		for i := range count - 1 {
			if id = projectIDLocalN(i + 1); id == "" {
				panic(fmt.Sprintf("MONGODB_ATLAS_PROJECT_ID_%d expected to be set (test requires %d projects)", i+1, count))
			}
			projectIDs = append(projectIDs, id)
		}
		return projectIDs
	}

	createSharedProjects(tb, count)
	projectIDs := make([]string, count)
	for i, project := range sharedInfo.projects {
		projectIDs[i] = project.id
	}
	return projectIDs
}

func createSharedProjects(tb testing.TB, count int) {
	tb.Helper()
	if len(sharedInfo.projects) >= count {
		return
	}

	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()

	for len(sharedInfo.projects) < count {
		projectName := RandomProjectName()
		tb.Logf("Creating execution project (%d): %s\n", len(sharedInfo.projects)+1, projectName)
		projectID := createProject(tb, projectName)
		sharedInfo.projects = append(sharedInfo.projects, projectInfo{
			id:   projectID,
			name: projectName,
		})
	}
}

// ProjectIDExecutionWithFreeCluster is identical to ProjectIDExecutionWithCluster but also contemplates the restriction of `MaxFreeTierClusterCount`
func ProjectIDExecutionWithFreeCluster(tb testing.TB, totalNodeCount, freeTierClusterCount int) (projectID, clusterName string) {
	tb.Helper()
	if ExistingClusterUsed() {
		return existingProjectIDClusterName()
	}
	// Only skip after ExistingClusterUsed() to allow MacT (Mocked-Acceptance Tests) to return early instead of being skipped.
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "sharedInfo not initialized, use acc.Run() to run tests that require shared resources")
	return NextProjectIDClusterName(totalNodeCount, freeTierClusterCount, func(projectName string) string {
		return createProject(tb, projectName)
	})
}

// ProjectIDExecutionWithCluster creates a project and reuses it with other tests respecting `MaxClusterNodesPerProject` restrictions. The clusterName is always unique.
// TotalNodeCount = sum(specs.node_count) * num_shards (1 if new schema)
// This avoids `CROSS_REGION_NETWORK_PERMISSIONS_LIMIT_EXCEEDED` and `project has reached the limit for the number of free clusters` errors when creating too many clusters within the same project.
// When `MONGODB_ATLAS_PROJECT_ID` and `MONGODB_ATLAS_CLUSTER_NAME` are defined, they are used instead of creating a project and clusterName.
func ProjectIDExecutionWithCluster(tb testing.TB, totalNodeCount int) (projectID, clusterName string) {
	tb.Helper()
	return ProjectIDExecutionWithFreeCluster(tb, totalNodeCount, 0)
}

// ClusterNameExecution returns the name of a created cluster for the execution of the tests in the resource package.
// This function relies on using an execution project and returns its id.
// When `MONGODB_ATLAS_CLUSTER_NAME` and `MONGODB_ATLAS_PROJECT_ID` are defined it will be used instead of creating resources. This is useful for local execution but not intended for CI executions.
func ClusterNameExecution(tb testing.TB, populateSampleData bool) (projectID, clusterName string) {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "sharedInfo not initialized, use acc.Run() to run tests that require shared resources")

	if ExistingClusterUsed() {
		return existingProjectIDClusterName()
	}

	projectID = ProjectIDExecution(tb) // ensure the execution project is created before cluster creation

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

// ProjectIDExecutionWithStreamInstance returns the project ID and stream instance name for test execution.
// Uses the same ProjectID as the ProjectIDExecution.
// The stream instance will include the `sample_stream_solar` connection. It is included to avoid ALREADY_EXIST errors as many different tests depends on this stream connection. Use a data source whenever you need it.
// You can use `MONGODB_ATLAS_STREAM_INSTANCE_NAME` to use an "externally" managed stream instance.
// We reuse a SPI to reduce the resource allocation (specially relevant for cloud dev).
func ProjectIDExecutionWithStreamInstance(tb testing.TB) (projectID, streamInstanceName string) {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "sharedInfo not initialized, use acc.Run() to run tests that require shared resources")
	projectID = ProjectIDExecution(tb)

	if existingStreamInstanceUsed() {
		return projectID, existingStreamInstanceName()
	}
	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()
	if sharedInfo.streamInstanceName == "" {
		name := RandomStreamInstanceName()
		tb.Logf("Creating execution stream instance: %s\n", name)
		sharedInfo.streamInstanceName = name
		createStreamInstance(tb, projectID, name)
	}

	return projectID, sharedInfo.streamInstanceName
}

// SerialSleep waits a few seconds so clusters in a project are not created concurrently, see HELP-65223.
// This must be called once the test is marked as parallel, e.g. in PreCheck inside Terraform tests.
func SerialSleep(tb testing.TB) {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "sharedInfo not initialized, use acc.Run() to run tests that require shared resources")

	sharedInfo.muSleep.Lock()
	defer sharedInfo.muSleep.Unlock()

	time.Sleep(5 * time.Second)
}

// NextProjectIDClusterName is an internal method used when we want to reuse a projectID respecting `MaxClustersNodesPerProject` and `MaxFreeTierClusterCount`
func NextProjectIDClusterName(totalNodeCount, freeTierClusterCount int, projectCreator func(string) string) (projectID, clusterName string) {
	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()
	var project projectInfo
	if len(sharedInfo.projects) == 0 ||
		sharedInfo.projects[len(sharedInfo.projects)-1].nodeCount+totalNodeCount > MaxClusterNodesPerProject ||
		sharedInfo.projects[len(sharedInfo.projects)-1].freeTierClusterCount+freeTierClusterCount > MaxFreeTierClusterCount {
		project = projectInfo{
			name:      RandomProjectName(),
			nodeCount: totalNodeCount,
		}
		project.id = projectCreator(project.name)
		sharedInfo.projects = append(sharedInfo.projects, project)
	} else {
		project = sharedInfo.projects[len(sharedInfo.projects)-1]
		sharedInfo.projects[len(sharedInfo.projects)-1].nodeCount += totalNodeCount
		sharedInfo.projects[len(sharedInfo.projects)-1].freeTierClusterCount += freeTierClusterCount
	}
	return project.id, RandomClusterName()
}
