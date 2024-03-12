package acc

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

// TestMainExecution must be called from TestMain in the test package if ProjectIDExecution is going to be used.
func TestMainExecution(m *testing.M) {
	if !InUnitTest() {
		atlasInfo.packageName = packageName()
		atlasInfo.projectName = RandomProjectName()
		fmt.Printf("CREATING EXECUTION PROJECT: %s, resource: %s\n", atlasInfo.projectName, atlasInfo.packageName)
		atlasInfo.projectID = createProject(atlasInfo.projectName)
	}

	exitCode := m.Run()

	if !InUnitTest() {
		fmt.Printf("DELETING EXECUTION PROJECT: %s, resource: %s\n", atlasInfo.projectName, atlasInfo.packageName)
		deleteProject(atlasInfo.projectID)
		atlasInfo.projectID = ""
		atlasInfo.projectName = ""
	}

	os.Exit(exitCode)
}

// ProjectIDExecution returns a project id created for the execution of the test group.
func ProjectIDExecution(tb testing.TB) string {
	tb.Helper()
	SkipInUnitTest(tb)
	require.NotEmpty(tb, atlasInfo.projectID)
	return atlasInfo.projectID
}

var atlasInfo = struct {
	projectID   string
	projectName string
	packageName string
}{}

func packageName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	pattern := `([^/]+)_test\.TestMain$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(f.Name())
	if len(matches) <= 1 {
		return ""
	}
	return matches[1]
}

func createProject(name string) string {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	if orgID == "" {
		fmt.Printf("Project creation failed: %s, org not set", name)
		return ""
	}
	params := &admin.Group{Name: name, OrgId: orgID}
	resp, _, err := ConnV2().ProjectsApi.CreateProject(context.Background(), params).Execute()
	id := resp.GetId()
	if err != nil || id == "" {
		fmt.Printf("Project creation failed: %s, error: %s", name, err)
		return ""
	}
	return id
}

func deleteProject(id string) {
	_, _, err := ConnV2().ProjectsApi.DeleteProject(context.Background(), id).Execute()
	if err != nil {
		fmt.Printf("Project deletion failed: %s, error: %s", id, err)
	}
}
