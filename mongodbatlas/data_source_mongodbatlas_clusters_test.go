package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasClusters_basic(t *testing.T) {
	var (
		cluster         matlas.Cluster
		resourceName    = "mongodbatlas_cluster.basic_ds"
		dataSourceName  = "data.mongodbatlas_clusters.basic_ds"
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name            = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		minSizeInstance = "M20"
		maxSizeInstance = "M80"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasClustersConfig(projectID, name, "true", "true", "true", minSizeInstance, maxSizeInstance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.name"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.labels.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.auto_scaling_compute_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.auto_scaling_compute_scale_down_enabled", "false"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasClustersConfig(projectID, name, backupEnabled, autoScalingEnabled, scaleDownEnabled, minSizeName, maxSizeName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "basic_ds" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 10

            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "US_EAST_2"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			provider_backup_enabled      = %s
			auto_scaling_disk_gb_enabled = true

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M40"

			labels {
				key   = "key 1"
				value = "value 1"
			}
			labels {
				key   = "key 2"
				value = "value 2"
			}

		auto_scaling_compute_enabled            = %s
		auto_scaling_compute_scale_down_enabled = %s
		provider_auto_scaling_compute_min_instance_size = "%s"
		provider_auto_scaling_compute_max_instance_size = "%s"
		}

		data "mongodbatlas_clusters" "basic_ds" {
			project_id = mongodbatlas_cluster.basic_ds.project_id
		}
	`, projectID, name, backupEnabled, autoScalingEnabled, scaleDownEnabled, minSizeName, maxSizeName)
}
