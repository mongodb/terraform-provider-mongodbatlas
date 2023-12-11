package acc

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func GetClusterInfo(projectID string) (clusterName, clusterNameStr, clusterTerraformStr string) {
	// Allows faster test execution in local, don't use in CI
	clusterName = os.Getenv("MONGODB_ATLAS_CLUSTER_NAME")
	if clusterName != "" {
		clusterNameStr = fmt.Sprintf("%q", clusterName)
	} else {
		clusterName = acctest.RandomWithPrefix("test-acc")
		clusterNameStr = "mongodbatlas_cluster.test_cluster.name"
		clusterTerraformStr = fmt.Sprintf(`
			resource "mongodbatlas_cluster" "test_cluster" {
				project_id   									= %[1]q
				name         									= %[2]q
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
		`, projectID, clusterName)
	}
	return clusterName, clusterNameStr, clusterTerraformStr
}
