package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	ClusterTagsMap1 = map[string]string{
		"key":   "key 1",
		"value": "value 1",
	}

	ClusterTagsMap2 = map[string]string{
		"key":   "key 2",
		"value": "value 2",
	}

	ClusterTagsMap3 = map[string]string{
		"key":   "key 3",
		"value": "value 3",
	}
)

func CheckDestroyCluster(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster" && rs.Type != "mongodbatlas_advanced_cluster" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		_, _, err := ConnV2().ClustersApi.GetCluster(context.Background(), projectID, clusterName).Execute()
		if err == nil {
			return fmt.Errorf("cluster (%s:%s) still exists", clusterName, rs.Primary.ID)
		}
	}
	return nil
}

func ConfigClusterGlobal(resourceName, projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" %[1]q {
			project_id              = %[2]q
			name                    = %[3]q
			disk_size_gb            = 80
			num_shards              = 1
			cloud_backup            = %[4]s
			cluster_type            = "GEOSHARDED"

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M30"

			replication_specs {
				zone_name  = "Zone 1"
				num_shards = 2
				regions_config {
				region_name     = "US_EAST_1"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
				}
			}

			replication_specs {
				zone_name  = "Zone 2"
				num_shards = 2
				regions_config {
				region_name     = "US_EAST_2"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
				}
			}
		}
	`, resourceName, projectID, name, backupEnabled)
}

func ImportStateClusterIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"]), nil
	}
}
