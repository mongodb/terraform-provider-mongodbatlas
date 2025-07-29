package advancedcluster_test

import (
	"fmt"
)

// func TestAccAdvancedCluster_moveNotSupportedLegacySchema(t *testing.T) {
// 	acc.SkipIfAdvancedClusterV2Schema(t) // This test is specific to the legacy schema
// 	var (
// 		projectID   = acc.ProjectIDExecution(t)
// 		clusterName = acc.RandomClusterName()
// 	)
// 	resource.ParallelTest(t, resource.TestCase{
// 		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
// 			tfversion.SkipBelow(tfversion.Version1_8_0),
// 		},
// 		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
// 		CheckDestroy:             acc.CheckDestroyCluster,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: configMoveFirst(projectID, clusterName),
// 			},
// 			{
// 				Config:      configMoveSecond(projectID, clusterName),
// 				ExpectError: regexp.MustCompile("Move Resource State Not Supported"),
// 			},
// 			{
// 				Config: configMoveFirst(projectID, clusterName),
// 			},
// 		},
// 	})
// }

func configMoveFirst(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "old" {
			project_id = %[1]q
			name = %[2]q
			cluster_type = "REPLICASET"
			provider_name = "AWS"
			provider_instance_size_name = "M10"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_WEST_2"
					electable_nodes = 3
					priority        = 7
				}
			}
		}
	`, projectID, clusterName)
}

func configMoveSecond(projectID, clusterName string) string {
	return fmt.Sprintf(`
		moved {
			from = mongodbatlas_cluster.old
			to   = mongodbatlas_advanced_cluster.test
		}

		resource "mongodbatlas_advanced_cluster" "test" {
    project_id   = %[1]q
    name         = %[2]q
    cluster_type = "REPLICASET"
    replication_specs {
      region_configs {
        electable_specs {
          instance_size = "M10"
          node_count    = 3
        }
        provider_name = "AWS"
        priority      = 7
        region_name   = "US_WEST_2"
      }
    }
  }
	`, projectID, clusterName)
}
