package acc

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// GetClusterInfo is used to obtain a project and cluster configuration resource.
// If env variables `MONGODB_ATLAS_CLUSTER_NAME` and `MONGODB_ATLAS_PROJECT_ID` are defined, creation of resources is avoided (useful for local execution)
func GetClusterInfo(orgID string) (projectIDStr, clusterName, clusterNameStr, clusterTerraformStr string) {
	// Allows faster test execution in local, don't use in CI
	clusterName = os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	if clusterName != "" && projectID != "" {
		clusterNameStr = fmt.Sprintf("%q", clusterName)
		projectIDStr = fmt.Sprintf("%q", projectID)
	} else {
		clusterName = acctest.RandomWithPrefix("test-acc")
		projectName := acctest.RandomWithPrefix("test-acc")
		projectIDStr = "mongodbatlas_project.test.id"
		clusterNameStr = "mongodbatlas_cluster.test_cluster.name"
		clusterTerraformStr = fmt.Sprintf(`
			resource "mongodbatlas_project" "test" {
				org_id = %[1]q
				name   = %[2]q
			}
		
			resource "mongodbatlas_cluster" "test_cluster" {
				project_id   									= mongodbatlas_project.test.id
				name         									= %[3]q
				disk_size_gb 									= 10
				backup_enabled               	= false
				auto_scaling_disk_gb_enabled	= false
				provider_name               	= "AWS"
				provider_instance_size_name 	= "M10"
			
				cluster_type = "REPLICASET"
				replication_specs {
					num_shards = 1
					regions_config {
						region_name     = "US_WEST_2"
						electable_nodes = 3
						priority        = 7
						read_only_nodes = 0
					}
				}
			}
		`, orgID, projectName, clusterName)
	}
	return projectIDStr, clusterName, clusterNameStr, clusterTerraformStr
}
