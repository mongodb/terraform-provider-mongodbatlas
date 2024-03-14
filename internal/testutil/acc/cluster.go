package acc

import (
	"fmt"
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
)

type ClusterRequest struct {
	OrgID        string
	ProviderName string
	ExtraConfig  string
	CloudBackup  bool
	Geosharded   bool
}

type ClusterInfo struct {
	ProjectIDStr        string
	ClusterName         string
	ClusterNameStr      string
	ClusterTerraformStr string
}

// GetClusterInfo is used to obtain a project and cluster configuration resource.
// When `MONGODB_ATLAS_CLUSTER_NAME` and `MONGODB_ATLAS_PROJECT_ID` are defined, creation of resources is avoided. This is useful for local execution but not intended for CI executions.
func GetClusterInfo(tb testing.TB, req *ClusterRequest) ClusterInfo {
	tb.Helper()
	if req == nil {
		req = new(ClusterRequest)
	}
	if req.OrgID == "" {
		req.OrgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
	}
	if req.ProviderName == "" {
		req.ProviderName = constant.AWS
	}
	clusterName := os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	if clusterName != "" && projectID != "" {
		return ClusterInfo{
			ProjectIDStr:        fmt.Sprintf("%q", projectID),
			ClusterName:         clusterName,
			ClusterNameStr:      fmt.Sprintf("%q", clusterName),
			ClusterTerraformStr: "",
		}
	}
	projectName := RandomProjectName()
	clusterName = RandomClusterName()
	clusterTypeStr := "REPLICASET"
	if req.Geosharded {
		clusterTypeStr = "GEOSHARDED"
	}
	clusterTerraformStr := fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}
	
		resource "mongodbatlas_cluster" "test_cluster" {
			project_id   									= mongodbatlas_project.test.id
			name         									= %[3]q
			cloud_backup         					= %[4]t
			auto_scaling_disk_gb_enabled	= false
			provider_name               	= %[5]q
			provider_instance_size_name 	= "M10"
		
			cluster_type = %[6]q
			replication_specs {
				num_shards = 1
				zone_name  = "Zone 1"
				regions_config {
					region_name     = "US_WEST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
			}
			%[7]s
		}
	`, req.OrgID, projectName, clusterName, req.CloudBackup, req.ProviderName, clusterTypeStr, req.ExtraConfig)
	return ClusterInfo{
		ProjectIDStr:        "mongodbatlas_project.test.id",
		ClusterName:         clusterName,
		ClusterNameStr:      "mongodbatlas_cluster.test_cluster.name",
		ClusterTerraformStr: clusterTerraformStr,
	}
}

func ExistingClusterUsed() bool {
	clusterName := os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	return clusterName != "" && projectID != ""
}
