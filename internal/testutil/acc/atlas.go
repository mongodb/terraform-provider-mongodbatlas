package acc

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

func ProjectIDExecution(tb testing.TB) string {
	tb.Helper()
	SkipInUnitTest(tb)
	atlasInfo.mu.Lock()
	defer atlasInfo.mu.Unlock()

	if atlasInfo.counter == 0 {
		atlasInfo.projectName = RandomProjectName()
		atlasInfo.projectID = createProject(tb, atlasInfo.projectName)
		tb.Logf("CREATING PROJECT EXECUTION: %s", atlasInfo.projectName)
	}

	atlasInfo.counter++
	tb.Cleanup(func() {
		cleanupExecution(tb)
	})

	return atlasInfo.projectID
}

var atlasInfo = struct {
	projectID   string
	projectName string
	counter     int
	mu          sync.Mutex
}{}

func cleanupExecution(tb testing.TB) {
	tb.Helper()
	atlasInfo.mu.Lock()
	defer atlasInfo.mu.Unlock()

	atlasInfo.counter--
	if atlasInfo.counter == 0 {
		deleteProject(tb, atlasInfo.projectID)
		tb.Logf("DELETING PROJECT EXECUTION: %s", atlasInfo.projectName)
		atlasInfo.projectID = ""
		atlasInfo.projectName = ""
	}
}

func createProject(tb testing.TB, name string) string {
	tb.Helper()
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	require.NotEmpty(tb, orgID)
	params := &admin.Group{Name: name, OrgId: orgID}
	resp, _, err := ConnV2().ProjectsApi.CreateProject(context.Background(), params).Execute()
	require.NoError(tb, err, "Project creation failed: %s, error: %s", name, err)
	id := resp.GetId()
	require.NotEmpty(tb, id, "Project creation failed: %s", name)
	return id
}

func deleteProject(tb testing.TB, id string) {
	tb.Helper()
	_, _, err := ConnV2().ProjectsApi.DeleteProject(context.Background(), id).Execute()
	assert.NoError(tb, err)
}
