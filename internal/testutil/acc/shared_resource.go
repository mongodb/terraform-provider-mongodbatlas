package acc

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
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

// ClusterNameExecution returns the name of a created cluster for the execution of the tests in the resource package.
// This function relies on using an execution project and returns its id.
// When `MONGODB_ATLAS_CLUSTER_NAME` and `MONGODB_ATLAS_PROJECT_ID` are defined it will be used instead of creating resources. This is useful for local execution but not intended for CI executions.
func ClusterNameExecution(tb testing.TB) (projectID, clusterName string) {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")

	localProjectID := projectIDLocal(tb)
	localClusterName := clusterNameLocal(tb)
	if localProjectID != "" && localClusterName != "" {
		return localProjectID, localClusterName
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

var sharedInfo = struct {
	projectID   string
	projectName string
	clusterName string
	mu          sync.Mutex
	init        bool
}{}
