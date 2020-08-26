package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMongoDBAtlasGlobalCluster_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_global_cluster_config.config"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name           = fmt.Sprintf("test-acc-global-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasGlobalClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasGlobalClusterConfig(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_name"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasGlobalClusterConfig(projectID, name string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cluster" "test" {
		project_id              = "%s"
		name                    = "%s"
		disk_size_gb            = 80
		provider_backup_enabled = false
		cluster_type            = "GEOSHARDED"

		// Provider Settings "block"
		provider_name               = "AWS"
		provider_disk_iops          = 240
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

	resource "mongodbatlas_global_cluster_config" "config" {
		project_id = mongodbatlas_cluster.test.project_id
		cluster_name = mongodbatlas_cluster.test.name

		managed_namespaces {
			db 				 = "mydata"
			collection 		 = "publishers"
			custom_shard_key = "city"
		}

		custom_zone_mappings {
			location ="CA"
			zone =  "Zone 1"
		}
	}

	data "mongodbatlas_global_cluster_config" "config" {
		project_id = mongodbatlas_global_cluster_config.config.project_id
		cluster_name = mongodbatlas_global_cluster_config.config.cluster_name
	}
	`, projectID, name)
}
