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
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

func createProject(tb testing.TB, name string) string {
	tb.Helper()
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	require.NotNil(tb, "Project creation failed: %s, org not set", name)
	params := &admin.Group{Name: name, OrgId: orgID}
	resp, _, err := ConnV2().ProjectsApi.CreateGroup(tb.Context(), params).Execute()
	require.NoError(tb, err, "Project creation failed: %s, err: %s", name, err)
	id := resp.GetId()
	require.NotEmpty(tb, id, "Project creation failed: %s", name)
	return id
}

func deleteProject(id string) {
	_, err := ConnV2().ProjectsApi.DeleteGroup(context.Background(), id).Execute()
	if admin.IsErrorCode(err, "CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_CLUSTERS") {
		fmt.Printf("Project deletion failed will retry in 30s: %s, error: %s", id, err)
		time.Sleep(30 * time.Second)
		_, err = ConnV2().ProjectsApi.DeleteGroup(context.Background(), id).Execute()
	}
	if err != nil {
		fmt.Printf("Project deletion failed: %s, error: %s", id, err)
	}
}

func createCluster(tb testing.TB, projectID, name string) string {
	tb.Helper()
	req := clusterReq(name, projectID)
	_, _, err := ConnV2().ClustersApi.CreateCluster(tb.Context(), projectID, &req).Execute()
	require.NoError(tb, err, "Cluster creation failed: %s, err: %s", name, err)
	// TODO: TEMPORARY CHANGE, DON'T MERGE
	// TODO: TEMPORARY CHANGE, DON'T MERGE
	stateConf := advancedcluster.CreateStateChangeConfig(tb.Context(), ConnV2(), projectID, name, 1*time.Hour)
	_, err = stateConf.WaitForStateContext(tb.Context())
	require.NoError(tb, err, "Cluster creation failed: %s, err: %s", name, err)

	return name
}

func deleteCluster(projectID, name string) {
	_, err := ConnV2().ClustersApi.DeleteCluster(context.Background(), projectID, name).Execute()
	if err != nil {
		fmt.Printf("Cluster deletion failed: %s %s, error: %s", projectID, name, err)
	}
	// TODO: TEMPORARY CHANGE, DON'T MERGE
	// TODO: TEMPORARY CHANGE, DON'T MERGE
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

func createStreamInstance(tb testing.TB, projectID, name string) {
	tb.Helper()
	req := admin.StreamsTenant{
		Name: admin.PtrString(name),
		DataProcessRegion: &admin.StreamsDataProcessRegion{
			Region:        "VIRGINIA_USA",
			CloudProvider: constant.AWS,
		},
		StreamConfig: &admin.StreamConfig{
			Tier: admin.PtrString("SP10"),
		},
		SampleConnections: &admin.StreamsSampleConnections{
			Solar: admin.PtrBool(true),
		},
	}
	_, _, err := ConnV2().StreamsApi.CreateStreamWorkspace(tb.Context(), projectID, &req).Execute()
	require.NoError(tb, err, "Stream instance creation failed: %s, err: %s", name, err)
}

func projectIDLocal() string {
	return os.Getenv("MONGODB_ATLAS_PROJECT_ID")
}
