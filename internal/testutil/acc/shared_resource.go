package acc

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// SetupSharedResources must be called from TestMain test package in order to use ProjectIDExecution.
func SetupSharedResources() {
	sharedInfo.init = true
}

// CleanupSharedResources must be called from TestMain test package in order to use ProjectIDExecution.
func CleanupSharedResources() {
	if sharedInfo.projectID != "" {
		fmt.Printf("Deleting execution project: %s, id: %s\n", sharedInfo.projectName, sharedInfo.projectID)
		deleteProject(sharedInfo.projectID)
	}
}

// ProjectIDExecution returns a project id created for the execution of the tests in the resource package.
func ProjectIDExecution(tb testing.TB) string {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")

	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()

	// lazy creation so it's only done if really needed
	if sharedInfo.projectID == "" {
		sharedInfo.projectName = RandomProjectName()
		tb.Logf("Creating execution project: %s\n", sharedInfo.projectName)
		sharedInfo.projectID = createProject(tb, sharedInfo.projectName)
	}

	return sharedInfo.projectID
}

var sharedInfo = struct {
	projectID   string
	projectName string
	mu          sync.Mutex
	init        bool
}{}
