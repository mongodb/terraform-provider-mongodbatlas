package acc

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cluster"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
)

func createProject(tb testing.TB, name string) string {
	tb.Helper()
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	require.NotNil(tb, "Project creation failed: %s, org not set", name)
	params := &admin.Group{Name: name, OrgId: orgID}
	resp, _, err := ConnV2().ProjectsAPI.CreateGroup(tb.Context(), params).Execute()
	require.NoError(tb, err, "Project creation failed: %s, err: %s", name, err)
	id := resp.GetId()
	require.NotEmpty(tb, id, "Project creation failed: %s", name)
	return id
}

func deleteProject(id string) error {
	_, err := ConnV2().ProjectsAPI.DeleteGroup(context.Background(), id).Execute()
	if admin.IsErrorCode(err, "CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_CLUSTERS") {
		fmt.Printf("Project deletion failed will retry in 30s: %s, error: %s", id, err)
		time.Sleep(30 * time.Second)
		_, err = ConnV2().ProjectsAPI.DeleteGroup(context.Background(), id).Execute()
	}
	if err != nil {
		return fmt.Errorf("project deletion failed: %s, error: %w", id, err)
	}
	return nil
}

func createCluster(tb testing.TB, projectID, name string, backupEnabled, pitEnabled bool) string {
	tb.Helper()
	regions := DefaultRegions()
	for i, region := range regions {
		lastRegion := i == len(regions)-1
		req := clusterReq(name, projectID, region, backupEnabled, pitEnabled)
		_, _, err := ConnV2().ClustersAPI.CreateCluster(tb.Context(), projectID, &req).Execute()
		if isOutOfCapacityError(err) && !lastRegion {
			tb.Logf("Cluster creation in %s failed with OUT_OF_CAPACITY, trying next region: %s, err: %s", region, name, err)
			continue
		}
		require.NoError(tb, err, "Cluster creation failed: %s, err: %s", name, err)
		stateConf := cluster.CreateStateChangeConfig(tb.Context(), ConnV2(), projectID, name, 1*time.Hour)
		_, err = stateConf.WaitForStateContext(tb.Context())
		if isOutOfCapacityError(err) && !lastRegion {
			tb.Logf("Cluster creation in %s failed with OUT_OF_CAPACITY while provisioning, trying next region: %s, err: %s", region, name, err)
			if err := deleteCluster(projectID, name); err != nil {
				tb.Logf("Failed to delete cluster %s after OUT_OF_CAPACITY: %s", name, err)
			}
			continue
		}
		require.NoError(tb, err, "Cluster creation failed: %s, err: %s", name, err)
		return name
	}
	require.Fail(tb, "Cluster creation failed with OUT_OF_CAPACITY in all candidate regions", "name: %s, regions: %s", name, strings.Join(regions, ","))
	return name
}

func isOutOfCapacityError(err error) bool {
	return err != nil && (admin.IsErrorCode(err, "OUT_OF_CAPACITY") || strings.Contains(err.Error(), "OUT_OF_CAPACITY"))
}

func deleteCluster(projectID, name string) error {
	_, err := ConnV2().ClustersAPI.DeleteCluster(context.Background(), projectID, name).Execute()
	if err != nil {
		return fmt.Errorf("cluster deletion failed: %s %s, error: %w", projectID, name, err)
	}
	stateConf := cluster.DeleteStateChangeConfig(context.Background(), ConnV2(), projectID, name, 1*time.Hour)
	_, err = stateConf.WaitForStateContext(context.Background())
	if err != nil {
		return fmt.Errorf("cluster deletion failed: %s %s, error: %w", projectID, name, err)
	}
	return nil
}

func clusterReq(name, projectID, region string, backupEnabled, pitEnabled bool) admin.ClusterDescription20240805 {
	return admin.ClusterDescription20240805{
		Name:          new(name),
		GroupId:       new(projectID),
		ClusterType:   new("REPLICASET"),
		BackupEnabled: new(backupEnabled),
		PitEnabled:    new(pitEnabled),
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{
			{
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						ProviderName: admin.PtrString(constant.AWS),
						RegionName:   new(region),
						Priority:     new(7),
						ElectableSpecs: &admin.HardwareSpec20240805{
							InstanceSize: admin.PtrString(constant.M10),
							NodeCount:    new(3),
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
		Name: new(name),
		DataProcessRegion: &admin.StreamsDataProcessRegion{
			Region:        "VIRGINIA_USA",
			CloudProvider: constant.AWS,
		},
		StreamConfig: &admin.StreamConfig{
			Tier: new("SP10"),
		},
		SampleConnections: &admin.StreamsSampleConnections{
			Solar: new(true),
		},
	}
	_, _, err := ConnV2().StreamsAPI.CreateStreamWorkspace(tb.Context(), projectID, &req).Execute()
	require.NoError(tb, err, "Stream instance creation failed: %s, err: %s", name, err)
}

func projectIDLocal() string {
	return os.Getenv("MONGODB_ATLAS_PROJECT_ID")
}

func projectIDLocalN(n int) string {
	return os.Getenv("MONGODB_ATLAS_PROJECT_ID_" + strconv.Itoa(n))
}
