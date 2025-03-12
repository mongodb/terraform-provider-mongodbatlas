package acc

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
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

func createCluster(tb testing.TB, projectID, name string) string {
	tb.Helper()
	req := clusterReq(name, projectID)
	_, _, err := ConnV2().ClustersApi.CreateCluster(context.Background(), projectID, &req).Execute()
	require.NoError(tb, err, "Cluster creation failed: %s, err: %s", name, err)

	stateConf := advancedcluster.CreateStateChangeConfig(context.Background(), ConnV2(), projectID, name, 1*time.Hour)
	_, err = stateConf.WaitForStateContext(context.Background())
	require.NoError(tb, err, "Cluster creation failed: %s, err: %s", name, err)

	return name
}

func deleteCluster(projectID, name string) {
	_, err := ConnV2().ClustersApi.DeleteCluster(context.Background(), projectID, name).Execute()
	if err != nil {
		fmt.Printf("Cluster deletion failed: %s %s, error: %s", projectID, name, err)
	}
	stateConf := advancedcluster.DeleteStateChangeConfig(context.Background(), ConnV2(), projectID, name, 1*time.Hour)
	_, err = stateConf.WaitForStateContext(context.Background())
	if err != nil {
		fmt.Printf("Cluster deletion failed: %s %s, error: %s", projectID, name, err)
	}
}

func clusterReq(name, projectID string) admin.ClusterDescription20240805 {
	return admin.ClusterDescription20240805{
		Name:        admin.PtrString(name),
		GroupId:     admin.PtrString(projectID),
		ClusterType: admin.PtrString("REPLICASET"),
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{
			{
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						ProviderName: admin.PtrString(constant.AWS),
						RegionName:   admin.PtrString(constant.UsWest2),
						Priority:     admin.PtrInt(7),
						ElectableSpecs: &admin.HardwareSpec20240805{
							InstanceSize: admin.PtrString(constant.M10),
							NodeCount:    admin.PtrInt(3),
						},
					},
				},
			},
		},
	}
}

// ProjectID returns the id for a project name.
// When `MONGODB_ATLAS_PROJECT_ID` is defined, it is used instead of creating a project. This is useful for local execution but not intended for CI executions.
func ProjectID(tb testing.TB, name string) string {
	tb.Helper()
	SkipInUnitTest(tb)

	if id := projectIDLocal(); id != "" {
		return id
	}

	resp, _, _ := ConnV2().ProjectsApi.GetProjectByName(context.Background(), name).Execute()
	id := resp.GetId()
	require.NotEmpty(tb, id, "Project name not found: %s", name)
	return id
}

func projectIDLocal() string {
	return os.Getenv("MONGODB_ATLAS_PROJECT_ID")
}
