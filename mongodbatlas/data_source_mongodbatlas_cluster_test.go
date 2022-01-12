package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasCluster_basic(t *testing.T) {
	var (
		cluster         matlas.Cluster
		resourceName    = "mongodbatlas_cluster.basic_ds"
		dataSourceName  = "data.mongodbatlas_cluster.basic_ds"
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
				Config: testAccDataSourceMongoDBAtlasClusterConfig(projectID, name, "true", "true", "true", minSizeInstance, maxSizeInstance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(dataSourceName, "pit_enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(dataSourceName, "labels.#", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "auto_scaling_disk_gb_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "auto_scaling_compute_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "auto_scaling_compute_scale_down_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "provider_auto_scaling_compute_min_instance_size", minSizeInstance),
					resource.TestCheckResourceAttr(dataSourceName, "provider_auto_scaling_compute_max_instance_size", maxSizeInstance),
					resource.TestCheckResourceAttr(dataSourceName, "version_release_system", "LTS"),
				),
			},
		},
	})
}

func TestAccDataSourceMongoDBAtlasCluster_advancedConf(t *testing.T) {
	var (
		cluster        matlas.Cluster
		resourceName   = "mongodbatlas_cluster.test"
		dataSourceName = "data.mongodbatlas_cluster.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name           = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasClusterConfigAdvancedConf(projectID, name, false, &matlas.ProcessArgs{
					FailIndexKeyTooLong:              pointy.Bool(true),
					JavascriptEnabled:                pointy.Bool(true),
					MinimumEnabledTLSProtocol:        "TLS1_1",
					NoTableScan:                      pointy.Bool(false),
					OplogSizeMB:                      pointy.Int64(1000),
					SampleRefreshIntervalBIConnector: pointy.Int64(310),
					SampleSizeBIConnector:            pointy.Int64(110),
				}),
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
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttr(dataSourceName, "pit_enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replication_specs.#"),
					resource.TestCheckResourceAttr(dataSourceName, "version_release_system", "LTS"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasClusterConfig(projectID, name, backupEnabled, autoScalingEnabled, scaleDownEnabled, minSizeName, maxSizeName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "basic_ds" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 100

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

		data "mongodbatlas_cluster" "basic_ds" {
			project_id = mongodbatlas_cluster.basic_ds.project_id
			name 	     = mongodbatlas_cluster.basic_ds.name
		}
	`, projectID, name, backupEnabled, autoScalingEnabled, scaleDownEnabled, minSizeName, maxSizeName)
}

func testAccDataSourceMongoDBAtlasClusterConfigAdvancedConf(projectID, name string, autoscalingEnabled bool, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "test" {
  project_id   = %[1]q
  name         = %[2]q
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

  provider_name               = "AWS"
  provider_instance_size_name = "M10"

  backup_enabled = false
  auto_scaling_disk_gb_enabled = %[3]t
  mongo_db_major_version       = "4.0"

  advanced_configuration  {
    fail_index_key_too_long              = %[4]t
    javascript_enabled                   = %[5]t
    minimum_enabled_tls_protocol         = %[6]q
    no_table_scan                        = %[7]t
    oplog_size_mb                        = %[8]d
    sample_size_bi_connector             = %[9]d
    sample_refresh_interval_bi_connector = %[10]d
  }
}

data "mongodbatlas_cluster" "test" {
  project_id = mongodbatlas_cluster.test.project_id
  name 	     = mongodbatlas_cluster.test.name
}
	`, projectID, name, autoscalingEnabled,
		*p.FailIndexKeyTooLong, *p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector)
}
