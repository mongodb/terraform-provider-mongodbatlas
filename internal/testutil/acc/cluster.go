package acc

import (
	"fmt"
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

type ClusterRequest struct {
	ResourceDependencyName string
	ClusterNameExplicit    string
	ReplicationSpecs       []admin.ReplicationSpec
	DiskSizeGb             int
	CloudBackup            bool
	Geosharded             bool
}

type ClusterInfo struct {
	ProjectIDStr        string
	ProjectID           string
	ClusterName         string
	ClusterResourceName string
	ClusterNameStr      string
	ClusterTerraformStr string
}

// GetClusterInfo is used to obtain a project and cluster configuration resource.
// When `MONGODB_ATLAS_CLUSTER_NAME` and `MONGODB_ATLAS_PROJECT_ID` are defined, creation of resources is avoided. This is useful for local execution but not intended for CI executions.
// Clusters will be created in project ProjectIDExecution.
func GetClusterInfo(tb testing.TB, req *ClusterRequest) ClusterInfo {
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
	clusterTerraformStr, clusterName, err := ClusterResourceHcl(projectID, req)
	if err != nil {
		tb.Error(err)
	}
	clusterResourceName := "mongodbatlas_advanced_cluster.cluster_info"
	return ClusterInfo{
		ProjectIDStr:        fmt.Sprintf("%q", projectID),
		ProjectID:           projectID,
		ClusterName:         clusterName,
		ClusterNameStr:      fmt.Sprintf("%s.name", clusterResourceName),
		ClusterResourceName: clusterResourceName,
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
	ProviderName string
	NodeCount    int
}

func (r *ReplicationSpecRequest) AddDefaults() {
	if r.NodeCount == 0 {
		r.NodeCount = 3
	}
	if r.ZoneName == "" {
		r.ZoneName = "Zone 1"
	}
	if r.Region == "" {
		r.Region = "US_WEST_2"
	}
	if r.InstanceSize == "" {
		r.InstanceSize = "M10"
	}
	if r.ProviderName == "" {
		r.ProviderName = constant.AWS
	}
}

func ReplicationSpec(req *ReplicationSpecRequest) admin.ReplicationSpec {
	if req == nil {
		req = new(ReplicationSpecRequest)
	}
	req.AddDefaults()
	defaultNumShards := 1
	regionConfig := CloudRegionConfig(req.Region, req.ProviderName, req.InstanceSize, req.NodeCount)
	configs := []admin.CloudRegionConfig{regionConfig}
	return admin.ReplicationSpec{
		NumShards:     &defaultNumShards,
		ZoneName:      &req.ZoneName,
		RegionConfigs: &configs,
	}
}

func CloudRegionConfig(regionName, provider, instanceSize string, nodeCount int) admin.CloudRegionConfig {
	return admin.CloudRegionConfig{
		RegionName:   &regionName,
		ProviderName: &provider,
		ElectableSpecs: &admin.HardwareSpec{
			InstanceSize: &instanceSize,
			NodeCount:    &nodeCount,
		},
	}
}
