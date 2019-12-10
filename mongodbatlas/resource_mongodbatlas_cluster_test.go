package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
)

func TestAccResourceMongoDBAtlasCluster_basicAWS(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	tiers := []string{"M10", "M20", "M30", "M40"}
	region := "EU_CENTRAL_1"

	for _, tier := range tiers {
		name := fmt.Sprintf("testAcc-%s-%s-%s", "AWS", tier, acctest.RandString(1))

		t.Run(name, func(t *testing.T) {

			resource.ParallelTest(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccMongoDBAtlasClusterConfigAWS(projectID, name, "true", tier, region),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
							testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttr(resourceName, "name", name),
							resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
						),
					},
					{
						Config: testAccMongoDBAtlasClusterConfigAWS(projectID, name, "false", tier, region),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
							testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttr(resourceName, "name", name),
							resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
						),
					},
				},
			})
		})
	}
}

func TestAccResourceMongoDBAtlasCluster_basicAzure(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var region string = "US_EAST_2"

	for tier, size := range matlas.DefaultDiskSizeGB["AZURE"] {
		name := fmt.Sprintf("testAcc-%s-%s-%ggb-%s", "AZURE", tier, size, acctest.RandString(1))

		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccMongoDBAtlasClusterConfigAzure(projectID, name, "true", tier, region),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
							testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttr(resourceName, "name", name),
							resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
						),
					},
					{
						Config: testAccMongoDBAtlasClusterConfigAzure(projectID, name, "false", tier, region),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
							testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttr(resourceName, "name", name),
							resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
						),
					},
				},
			})
		})
	}

}
func TestAccResourceMongoDBAtlasCluster_basicGCP(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var region string

	for tier, size := range matlas.DefaultDiskSizeGB["GCP"] {
		name := fmt.Sprintf("testAcc-%s-%s-%ggb-%s", "GCP", tier, size, acctest.RandString(1))

		switch tier {
		case "M200":
			region = "EUROPE_WEST_3"
		case "M300":
			region = "CENTRAL_US"
		default:
			region = "US_EAST_4"
		}

		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccMongoDBAtlasClusterConfigGCP(projectID, name, "true", tier, region),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
							testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttr(resourceName, "name", name),
							resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
						),
					},
					{
						Config: testAccMongoDBAtlasClusterConfigGCP(projectID, name, "false", tier, region),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
							testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttr(resourceName, "name", name),
							resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
							resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
						),
					},
				},
			})
		})
	}
}

func TestAccResourceMongoDBAtlasCluster_tenant(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.tenant"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var region string

	for tier, size := range matlas.DefaultDiskSizeGB["TENANT"] {
		name := fmt.Sprintf("testAcc-%s-%s-%ggb-%s", "TEN", tier, size, acctest.RandString(1))

		region = []string{"EU_WEST_1", "EU_CENTRAL_1"}[acctest.RandIntRange(0, 2)]

		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccMongoDBAtlasClusterConfigTenant(projectID, name, "false", cast.ToString(size), tier, region),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
							testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttr(resourceName, "name", name),
							resource.TestCheckResourceAttr(resourceName, "disk_size_gb", cast.ToString(size)),
							resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
						),
					},
				},
			})
		})
	}
}

func TestAccResourceMongoDBAtlasCluster_tenantWithoutDiskSizeGb(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.tenant"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var region string

	for tier, size := range matlas.DefaultDiskSizeGB["TENANT"] {
		name := fmt.Sprintf("testAcc-%s-%s-%ggb-%s", "TEN", tier, size, acctest.RandString(1))

		region = []string{"EU_WEST_1", "EU_CENTRAL_1"}[acctest.RandIntRange(0, 2)]

		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccMongoDBAtlasClusterConfigTenantWithoutDiskSizeGb(projectID, name, tier, region),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
							testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttr(resourceName, "name", name),
							resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
						),
					},
				},
			})
		})
	}
}
func TestAccResourceMongoDBAtlasCluster_basicAdvancedConf(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigAdvancedConf(projectID, name, "false", &matlas.ProcessArgs{
					FailIndexKeyTooLong:              pointy.Bool(true),
					JavascriptEnabled:                pointy.Bool(true),
					MinimumEnabledTLSProtocol:        "TLS1_2",
					NoTableScan:                      pointy.Bool(true),
					OplogSizeMB:                      pointy.Int64(1000),
					SampleRefreshIntervalBIConnector: pointy.Int64(310),
					SampleSizeBIConnector:            pointy.Int64(110),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.fail_index_key_too_long", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.minimum_enabled_tls_protocol", "TLS1_2"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.no_table_scan", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttrSet(resourceName, "advanced_configuration.%"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterConfigAdvancedConf(projectID, name, "true", &matlas.ProcessArgs{
					FailIndexKeyTooLong:              pointy.Bool(false),
					JavascriptEnabled:                pointy.Bool(false),
					MinimumEnabledTLSProtocol:        "TLS1_1",
					NoTableScan:                      pointy.Bool(false),
					OplogSizeMB:                      pointy.Int64(990),
					SampleRefreshIntervalBIConnector: pointy.Int64(0),
					SampleSizeBIConnector:            pointy.Int64(0),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.javascript_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.oplog_size_mb", "990"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.sample_size_bi_connector", "0"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.sample_refresh_interval_bi_connector", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "advanced_configuration.%"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCluster_MultiRegion(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-multi-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigMultiRegion(projectID, name, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "3"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterConfigMultiRegion(projectID, name, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "3"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasCluster_Global(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-global-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigGlobal(projectID, name, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.1.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "80"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "GEOSHARDED"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.1.regions_config.#", "1"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasCluster_importBasic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	importStateID := fmt.Sprintf("%s-%s", projectID, clusterName)

	resourceName := "mongodbatlas_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterConfigAWS(projectID, clusterName, "true", "M10", "CENTRAL_US"),
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           importStateID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccCheckMongoDBAtlasClusterExists(resourceName string, cluster *matlas.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])

		if clusterResp, _, err := conn.Clusters.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"]); err == nil {
			*cluster = *clusterResp
			return nil
		}

		return fmt.Errorf("cluster(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasClusterAttributes(cluster *matlas.Cluster, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cluster.Name != name {
			return fmt.Errorf("bad name: %s", cluster.Name)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster" {
			continue
		}

		// Try to find the cluster
		_, _, err := conn.Clusters.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["name"])

		if err == nil {
			return fmt.Errorf("cluster (%s:%s) still exists", rs.Primary.Attributes["name"], rs.Primary.ID)
		}
	}

	return nil
}

func testAccMongoDBAtlasClusterConfigAdvancedConf(projectID, name, backupEnabled string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 10

			replication_factor           = 3
			backup_enabled               = %s
			auto_scaling_disk_gb_enabled = true
			mongo_db_major_version       = "4.0"
		
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_disk_iops          = 100
			provider_encrypt_ebs_volume = false
			provider_instance_size_name = "M10"
			provider_region_name        = "EU_CENTRAL_1"
		
			advanced_configuration = {
				fail_index_key_too_long              = %t
				javascript_enabled                   = %t
				minimum_enabled_tls_protocol         = "%s"
				no_table_scan                        = %t
				oplog_size_mb                        = %d
				sample_size_bi_connector			       = %d
				sample_refresh_interval_bi_connector = %d
			}
		}
	`, projectID, name, backupEnabled,
		*p.FailIndexKeyTooLong, *p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector)
}

func testAccMongoDBAtlasClusterConfigAWS(projectID, name, backupEnabled, tier, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%s"
			name         = "%s"
			num_shards   = 1
		
			replication_factor           = 3
			backup_enabled               = %s
			auto_scaling_disk_gb_enabled = false
			mongo_db_major_version       = "4.0"
		
			//Provider Settings "block"
			provider_name               = "AWS"
			provider_encrypt_ebs_volume = false
			provider_instance_size_name = "%s"
			provider_region_name        = "%s"
		}
	`, projectID, name, backupEnabled, tier, region)
}

func testAccMongoDBAtlasClusterConfigAzure(projectID, name, backupEnabled, tier, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%s"
			name         = "%s"
			num_shards   = 1
			
			replication_factor           = 3
			backup_enabled               = %s
			auto_scaling_disk_gb_enabled = false
			mongo_db_major_version       = "4.0"
			
			//Provider Settings "block"
			provider_name               = "AZURE"
			provider_disk_type_name     = "P6"
			provider_instance_size_name = "%s"
			provider_region_name        = "%s"
		}
	`, projectID, name, backupEnabled, tier, region)
}

func testAccMongoDBAtlasClusterConfigGCP(projectID, name, backupEnabled, tier, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%s"
			name         = "%s"
			num_shards   = 1
			
			replication_factor           = 3
			backup_enabled               = %s
			auto_scaling_disk_gb_enabled = true
			mongo_db_major_version       = "4.0"
			
			//Provider Settings "block"
			provider_name               = "GCP"
			provider_instance_size_name = "%s"
			provider_region_name        = "%s"
		}
	`, projectID, name, backupEnabled, tier, region)
}

func testAccMongoDBAtlasClusterConfigTenant(projectID, name, autoScaling, size, tier, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "tenant" {
			project_id             = "%s"
			name                   = "%s"
			mongo_db_major_version = "4.0"
		
			backing_provider_name        = "AWS"
			auto_scaling_disk_gb_enabled = "%s"
			disk_size_gb                 = "%s"
		
			provider_name               = "TENANT"
			provider_instance_size_name = "%s"
			provider_region_name        = "%s"
		}
	`, projectID, name, autoScaling, size, tier, region)
}

func testAccMongoDBAtlasClusterConfigTenantWithoutDiskSizeGb(projectID, name, tier, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "tenant" {
			project_id             = "%s"
			name                   = "%s"
			mongo_db_major_version = "4.0"
		
			backing_provider_name        = "AWS"
		
			provider_name               = "TENANT"
			provider_instance_size_name = "%s"
			provider_region_name        = "%s"
		}
	`, projectID, name, tier, region)
}

func testAccMongoDBAtlasClusterConfigMultiRegion(projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id     = "%s"
			name           = "%s"
			disk_size_gb   = 100
			num_shards     = 1
			backup_enabled = %s
			cluster_type   = "REPLICASET"

			//Provider Settings "block"
			provider_name               = "AWS"
			provider_disk_iops          = 300
			provider_instance_size_name = "M10"

			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_WEST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
				regions_config {
					region_name     = "EU_CENTRAL_1"
					electable_nodes = 2
					priority        = 6
					read_only_nodes = 0
				}
				regions_config {
					region_name     = "US_WEST_1"
					electable_nodes = 2
					priority        = 5
					read_only_nodes = 2
				}
			}
		}
	`, projectID, name, backupEnabled)
}

func testAccMongoDBAtlasClusterConfigGlobal(projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id              = "%s"
			name                    = "%s"
			disk_size_gb            = 80
			num_shards              = 1
			backup_enabled          = %s
			provider_backup_enabled = true
			cluster_type            = "GEOSHARDED"
			
			//Provider Settings "block"
			provider_name               = "AWS"
			provider_disk_iops          = 240
			provider_instance_size_name = "M30"
			
			replication_specs {
				zone_name  = "Zone 1"
				num_shards = 2
				regions_config {
				region_name     = "EU_CENTRAL_1"
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
	`, projectID, name, backupEnabled)
}
