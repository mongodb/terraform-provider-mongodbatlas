package acc

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

// TestMainExecution must be called from TestMain in the test package if ProjectIDExecution is going to be used.
func TestMainExecution(m *testing.M) {
	if !InUnitTest() {
		atlasInfo.init = true
		atlasInfo.resourceName = resourceName()
	}

	exitCode := m.Run()

	if !InUnitTest() && atlasInfo.needsDeletion {
		fmt.Printf("Deleting execution project: %s, resource: %s\n", atlasInfo.projectName, atlasInfo.resourceName)
		deleteProject(atlasInfo.projectID)
	}

	os.Exit(exitCode)
}

// ProjectIDExecution returns a project id created for the execution of the resource tests.
func ProjectIDExecution(tb testing.TB) string {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, atlasInfo.init, "TestMainExecution must called to be able to use ProjectIDExecution")

	atlasInfo.mu.Lock()
	defer atlasInfo.mu.Unlock()

	// lazy creation so it's only done if really needed
	if atlasInfo.projectName == "" {
		var globalName, globalID string
		if atlasInfo.resourceName != "" {
			globalName = (prefixProjectKeep + "-" + atlasInfo.resourceName)[:projectNameMaxLen]
			globalID = projectID(globalName)
		}

		if globalID == "" {
			atlasInfo.projectName = RandomProjectName()
			tb.Logf("Creating execution project: %s, resource: %s, global project (not found): %s\n", atlasInfo.projectName, atlasInfo.resourceName, globalName)
			atlasInfo.projectID = createProject(tb, atlasInfo.projectName)
			atlasInfo.needsDeletion = true
		} else {
			atlasInfo.projectName = globalName
			tb.Logf("Reusing global project: %s, resource: %s\n", atlasInfo.projectName, atlasInfo.resourceName)
			atlasInfo.projectID = globalID
		}
	}

	return atlasInfo.projectID
}

var atlasInfo = struct {
	projectID     string
	projectName   string
	resourceName  string
	mu            sync.Mutex
	init          bool
	needsDeletion bool
}{}

const (
	projectNameMaxLen = 64
)

func resourceName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	pattern := `([^/]+)_test\.TestMain$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(runtime.FuncForPC(pc).Name())
	if len(matches) <= 1 {
		return ""
	}
	return matches[1]
}

func createProject(tb testing.TB, name string) string {
	tb.Helper()
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	require.NotNil(tb, "Project creation failed: %s, org not set", name)
	params := &admin.Group{Name: name, OrgId: orgID}
	resp, _, err := ConnV2().ProjectsApi.CreateProject(context.Background(), params).Execute()
	require.NoError(tb, err, "Project creation failed: %s, err: %s", name, err)
	id := resp.GetId()
	require.NotEmpty(tb, id, "Project creation failed: %s", name)
	return id
}

func deleteProject(id string) {
	_, _, err := ConnV2().ProjectsApi.DeleteProject(context.Background(), id).Execute()
	if err != nil {
		fmt.Printf("Project deletion failed: %s, error: %s", id, err)
	}
}

func projectID(name string) string {
	resp, _, _ := ConnV2().ProjectsApi.GetProjectByName(context.Background(), name).Execute()
	return resp.GetId()
}
