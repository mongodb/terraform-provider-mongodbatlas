package acc

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/replay"
	"github.com/stretchr/testify/require"
)

// SetupSharedResources must be called from TestMain test package in order to use ProjectIDExecution.
// It returns the cleanup function that must be called at the end of TestMain.
func SetupSharedResources() func() {
	sharedInfo.init = true
	return cleanupSharedResources
}

func cleanupSharedResources() {
	if sharedInfo.projectID != "" {
		fmt.Printf("Deleting execution project: %s, id: %s\n", sharedInfo.projectName, sharedInfo.projectID)
		deleteProject(sharedInfo.projectID)
		if replay.IsInCaptureMode() {
			if err := replay.CaptureExecutionVariables(sharedInfo.projectID); err != nil {
				log.Fatal("Failed to generate file for capturing execution variables")
			}
		}
	}
}

// ProjectIDExecution returns a project id created for the execution of the tests in the resource package.
// Even if a GH test group is run, every resource/package will create its own project, not a shared project for all the test group.
func ProjectIDExecution(tb testing.TB) string {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")

	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()

	// lazy creation so it's only done if really needed
	if sharedInfo.projectID == "" {
		if replay.IsInSimulateMode() {
			vars, _ := replay.ObtainExecutionVariables()
			tb.Logf("Using project id from simulation execution variables %s\n", vars.ProjectID)
			sharedInfo.projectID = vars.ProjectID
		} else {
			sharedInfo.projectName = RandomProjectName()
			tb.Logf("Creating execution project: %s\n", sharedInfo.projectName)
			sharedInfo.projectID = createProject(tb, sharedInfo.projectName)
		}
	}

	return sharedInfo.projectID
}

var sharedInfo = struct {
	projectID   string
	projectName string
	mu          sync.Mutex
	init        bool
}{}
