package acc

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

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

func ProjectID(tb testing.TB, name string) string {
	tb.Helper()
	SkipInUnitTest(tb)
	resp, _, _ := ConnV2().ProjectsApi.GetProjectByName(context.Background(), name).Execute()
	id := resp.GetId()
	require.NotEmpty(tb, id, "Project name not found: %s", name)
	return id
}
