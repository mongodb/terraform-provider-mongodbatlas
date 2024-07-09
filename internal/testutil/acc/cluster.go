package acc

import (
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

type ClusterRequest struct {
	ProviderName           string
	ResourceDependencyName string
	ClusterNameExplicit    string
	CloudBackup            bool
	Geosharded             bool
}

type ClusterInfo struct {
	ProjectIDStr        string
	ProjectID           string
	ClusterName         string
	ClusterNameStr      string
	ClusterTerraformStr string
}

// GetClusterInfo is used to obtain a project and cluster configuration resource.
// When `MONGODB_ATLAS_CLUSTER_NAME` and `MONGODB_ATLAS_PROJECT_ID` are defined, creation of resources is avoided. This is useful for local execution but not intended for CI executions.
// Clusters will be created in project ProjectIDExecution.
func GetClusterInfo(tb testing.TB, req *ClusterRequest, specs ...admin.ReplicationSpec) ClusterInfo {
	tb.Helper()
	if req == nil {
		req = new(ClusterRequest)
	}
	clusterName := os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	if clusterName != "" && projectID != "" {
		return ClusterInfo{
			ProjectIDStr:        fmt.Sprintf("%q", projectID),
			ProjectID:           projectID,
			ClusterName:         clusterName,
			ClusterNameStr:      fmt.Sprintf("%q", clusterName),
			ClusterTerraformStr: "",
		}
	}
	projectID = ProjectIDExecution(tb)
	clusterTerraformStr, clusterName, err := ClusterResourceHcl(projectID, req, specs)
	if err != nil {
		tb.Error(err)
	}
	return ClusterInfo{
		ProjectIDStr:        fmt.Sprintf("%q", projectID),
		ProjectID:           projectID,
		ClusterName:         clusterName,
		ClusterNameStr:      "mongodbatlas_advanced_cluster.cluster_info.name",
		ClusterTerraformStr: clusterTerraformStr,
	}
}

func ExistingClusterUsed() bool {
	clusterName := os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	return clusterName != "" && projectID != ""
}

type ReplicationSpecRequest struct {
	ZoneName     string
	Region       string
	InstanceSize string
	NumShards    int
	NodeCount    int
}

func ReplicationSpec(req *ReplicationSpecRequest) admin.ReplicationSpec {
	if req == nil {
		req = new(ReplicationSpecRequest)
	}
	if req.NumShards == 0 {
		req.NumShards = 1
	}
	if req.NodeCount == 0 {
		req.NodeCount = 3
	}
	if req.ZoneName == "" {
		req.ZoneName = "zone1"
	}
	if req.Region == "" {
		req.Region = "US_WEST_1"
	}
	if req.InstanceSize == "" {
		req.InstanceSize = "M10"
	}
	spec := admin.ReplicationSpec{
		NumShards: &req.NumShards,
		ZoneName:  &req.ZoneName,
		RegionConfigs: &[]admin.CloudRegionConfig{
			{
				RegionName: &req.Region,
				ElectableSpecs: &admin.HardwareSpec{
					InstanceSize: &req.InstanceSize,
					NodeCount:    &req.NodeCount,
				},
			},
		},
	}
	return spec
}
