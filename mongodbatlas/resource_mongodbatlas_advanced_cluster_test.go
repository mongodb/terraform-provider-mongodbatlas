package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccClusterRSAdvancedCluster_basicTenant(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		rName        = acctest.RandomWithPrefix("test-acc")
		rNameUpdated = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigTenant(projectID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "termination_protection_enabled"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigTenant(projectID, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasClusterImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccClusterRSAdvancedCluster_singleProvider(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		rName        = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProvider(projectID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(projectID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasClusterImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
		},
	})
}

func TestAccClusterRSAdvancedCluster_multicloud(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		rName        = acctest.RandomWithPrefix("test-acc")
		rNameUpdated = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(projectID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(projectID, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasClusterImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
		},
	})
}

func TestAccClusterRSAdvancedCluster_multicloudSharded(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		rName        = acctest.RandomWithPrefix("test-acc")
		rNameUpdated = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloudSharded(projectID, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloudSharded(projectID, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasClusterImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
		},
	})
}

func TestAccClusterRSAdvancedCluster_Paused(t *testing.T) {
	SkipTest(t)
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		rName        = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(projectID, rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(projectID, rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasClusterImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
		},
	})
}

func TestAccClusterRSAdvancedCluster_advancedConf(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		rName        = acctest.RandomWithPrefix("test-acc")
		rNameUpdated = acctest.RandomWithPrefix("test-acc")
		processArgs  = &matlas.ProcessArgs{
			DefaultReadConcern:               "available",
			DefaultWriteConcern:              "1",
			FailIndexKeyTooLong:              pointy.Bool(false),
			JavascriptEnabled:                pointy.Bool(true),
			MinimumEnabledTLSProtocol:        "TLS1_1",
			NoTableScan:                      pointy.Bool(false),
			OplogSizeMB:                      pointy.Int64(1000),
			SampleRefreshIntervalBIConnector: pointy.Int64(310),
			SampleSizeBIConnector:            pointy.Int64(110),
		}
		processArgsUpdated = &matlas.ProcessArgs{
			DefaultReadConcern:               "available",
			DefaultWriteConcern:              "0",
			FailIndexKeyTooLong:              pointy.Bool(false),
			JavascriptEnabled:                pointy.Bool(true),
			MinimumEnabledTLSProtocol:        "TLS1_2",
			NoTableScan:                      pointy.Bool(false),
			OplogSizeMB:                      pointy.Int64(1000),
			SampleRefreshIntervalBIConnector: pointy.Int64(310),
			SampleSizeBIConnector:            pointy.Int64(110),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigAdvancedConf(projectID, rName, processArgs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigAdvancedConf(projectID, rNameUpdated, processArgsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_2"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
				),
			},
		},
	})
}

func TestAccClusterRSAdvancedCluster_DefaultWrite(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		rName        = acctest.RandomWithPrefix("test-acc")
		rNameUpdated = acctest.RandomWithPrefix("test-acc")
		processArgs  = &matlas.ProcessArgs{
			DefaultReadConcern:               "available",
			DefaultWriteConcern:              "1",
			JavascriptEnabled:                pointy.Bool(true),
			MinimumEnabledTLSProtocol:        "TLS1_1",
			NoTableScan:                      pointy.Bool(false),
			OplogSizeMB:                      pointy.Int64(1000),
			SampleRefreshIntervalBIConnector: pointy.Int64(310),
			SampleSizeBIConnector:            pointy.Int64(110),
		}
		processArgsUpdated = &matlas.ProcessArgs{
			DefaultReadConcern:               "available",
			DefaultWriteConcern:              "majority",
			JavascriptEnabled:                pointy.Bool(true),
			MinimumEnabledTLSProtocol:        "TLS1_2",
			NoTableScan:                      pointy.Bool(false),
			OplogSizeMB:                      pointy.Int64(1000),
			SampleRefreshIntervalBIConnector: pointy.Int64(310),
			SampleSizeBIConnector:            pointy.Int64(110),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigAdvancedConfDefaultWrite(projectID, rName, processArgs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigAdvancedConfDefaultWrite(projectID, rNameUpdated, processArgsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "majority"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_2"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
				),
			},
		},
	})
}

func TestAccClusterRSAdvancedClusterConfig_ReplicationSpecsAutoScaling(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		rName        = acctest.RandomWithPrefix("test-acc")
		rNameUpdated = acctest.RandomWithPrefix("test-acc")
		autoScaling  = &matlas.AutoScaling{
			Compute:       &matlas.Compute{Enabled: pointy.Bool(false), MaxInstanceSize: ""},
			DiskGBEnabled: pointy.Bool(true),
		}
		autoScalingUpdated = &matlas.AutoScaling{
			Compute:       &matlas.Compute{Enabled: pointy.Bool(true), MaxInstanceSize: "M20"},
			DiskGBEnabled: pointy.Bool(true),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAutoScaling(projectID, rName, autoScaling),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					testAccCheckMongoDBAtlasAdvancedClusterScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAutoScaling(projectID, rNameUpdated, autoScalingUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					testAccCheckMongoDBAtlasAdvancedClusterScaling(&cluster, *autoScalingUpdated.Compute.Enabled),
				),
			},
		},
	})
}

func TestAccClusterRSAdvancedClusterConfig_ReplicationSpecsAnalyticsAutoScaling(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		rName        = acctest.RandomWithPrefix("test-acc")
		rNameUpdated = acctest.RandomWithPrefix("test-acc")
		autoScaling  = &matlas.AutoScaling{
			Compute:       &matlas.Compute{Enabled: pointy.Bool(false), MaxInstanceSize: ""},
			DiskGBEnabled: pointy.Bool(true),
		}
		autoScalingUpdated = &matlas.AutoScaling{
			Compute:       &matlas.Compute{Enabled: pointy.Bool(true), MaxInstanceSize: "M20"},
			DiskGBEnabled: pointy.Bool(true),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(projectID, rName, autoScaling),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					testAccCheckMongoDBAtlasAdvancedClusterAnalyticsScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(projectID, rNameUpdated, autoScalingUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					testAccCheckMongoDBAtlasAdvancedClusterAnalyticsScaling(&cluster, *autoScalingUpdated.Compute.Enabled),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName string, cluster *matlas.AdvancedCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		log.Printf("[DEBUG] projectID: %s, name %s", ids["project_id"], ids["cluster_name"])

		if clusterResp, _, err := conn.AdvancedClusters.Get(context.Background(), ids["project_id"], ids["cluster_name"]); err == nil {
			*cluster = *clusterResp
			return nil
		}

		return fmt.Errorf("cluster(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasAdvancedClusterAttributes(cluster *matlas.AdvancedCluster, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cluster.Name != name {
			return fmt.Errorf("bad name: %s", cluster.Name)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasAdvancedClusterScaling(cluster *matlas.AdvancedCluster, computeEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *cluster.ReplicationSpecs[0].RegionConfigs[0].AutoScaling.Compute.Enabled != computeEnabled {
			return fmt.Errorf("compute_enabled: %d", cluster.ReplicationSpecs[0].RegionConfigs[0].AutoScaling.Compute.Enabled)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasAdvancedClusterAnalyticsScaling(cluster *matlas.AdvancedCluster, computeEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *cluster.ReplicationSpecs[0].RegionConfigs[0].AnalyticsAutoScaling.Compute.Enabled != computeEnabled {
			return fmt.Errorf("compute_enabled: %d", cluster.ReplicationSpecs[0].RegionConfigs[0].AnalyticsAutoScaling.Compute.Enabled)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasAdvancedClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_advanced_cluster" {
			continue
		}

		// Try to find the cluster
		_, _, err := conn.AdvancedClusters.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"])

		if err == nil {
			return fmt.Errorf("cluster (%s:%s) still exists", rs.Primary.Attributes["cluster_name"], rs.Primary.ID)
		}
	}

	return nil
}

func testAccMongoDBAtlasAdvancedClusterConfigTenant(projectID, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = %[1]q
  name         = %[2]q
  cluster_type = "REPLICASET"
  retain_backups_enabled = "false"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M5"
      }
      provider_name         = "TENANT"
      backing_provider_name = "AWS"
      region_name           = "US_EAST_1"
      priority              = 7
    }
  }
}
	`, projectID, name)
}

func testAccMongoDBAtlasAdvancedClusterConfigSingleProvider(projectID, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = %[1]q
  name         = %[2]q
  cluster_type = "REPLICASET"
  retain_backups_enabled = "false"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
  }
}
	`, projectID, name)
}

func testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(projectID, name string) string {
	return fmt.Sprintf(`
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
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 2
      }
      provider_name = "GCP"
      priority      = 6
      region_name   = "NORTH_AMERICA_NORTHEAST_1"
    }
  }
}
	`, projectID, name)
}

func testAccMongoDBAtlasAdvancedClusterConfigMultiCloudSharded(projectID, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = %[1]q
  name         = %[2]q
  cluster_type = "SHARDED"

  replication_specs {
    num_shards = 1
    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M30"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 2
      }
      provider_name = "AZURE"
      priority      = 6
      region_name   = "US_EAST_2"
    }
  }
}
	`, projectID, name)
}

func testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(projectID, name string, paused bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = %[1]q
  name         = %[2]q
  cluster_type = "REPLICASET"
  paused       = %[3]t

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
  }
}
	`, projectID, name, paused)
}

func testAccMongoDBAtlasAdvancedClusterConfigAdvancedConf(projectID, name string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = %[1]q
  name                   = %[2]q
  cluster_type           = "REPLICASET"

   replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
  }

  advanced_configuration  {
    fail_index_key_too_long              = %[3]t
    javascript_enabled                   = %[4]t
    minimum_enabled_tls_protocol         = %[5]q
    no_table_scan                        = %[6]t
    oplog_size_mb                        = %[7]d
    sample_size_bi_connector			 = %[8]d
    sample_refresh_interval_bi_connector = %[9]d
  }
}

	`, projectID, name,
		*p.FailIndexKeyTooLong, *p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector)
}

func testAccMongoDBAtlasAdvancedClusterConfigAdvancedConfDefaultWrite(projectID, name string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = %[1]q
  name                   = %[2]q
  cluster_type           = "REPLICASET"

   replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
  }

  advanced_configuration  {
    javascript_enabled                   = %[3]t
    minimum_enabled_tls_protocol         = %[4]q
    no_table_scan                        = %[5]t
    oplog_size_mb                        = %[6]d
    sample_size_bi_connector			 = %[7]d
    sample_refresh_interval_bi_connector = %[8]d
    default_read_concern                 = %[9]q
    default_write_concern                = %[10]q
  }
}

	`, projectID, name, *p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector, p.DefaultReadConcern, p.DefaultWriteConcern)
}

func testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAutoScaling(projectID, name string, p *matlas.AutoScaling) string {
	return fmt.Sprintf(`
resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = %[1]q
  name                   = %[2]q
  cluster_type           = "REPLICASET"

   replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
	  auto_scaling {
        compute_enabled = %[3]t
        disk_gb_enabled = %[4]t
		compute_max_instance_size = %[5]q
	  }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
  }


}

	`, projectID, name, *p.Compute.Enabled, *p.DiskGBEnabled, p.Compute.MaxInstanceSize)
}

func testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(projectID, name string, p *matlas.AutoScaling) string {
	return fmt.Sprintf(`
resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = %[1]q
  name                   = %[2]q
  cluster_type           = "REPLICASET"

   replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
	  analytics_auto_scaling {
        compute_enabled = %[3]t
        disk_gb_enabled = %[4]t
		compute_max_instance_size = %[5]q
	  }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
  }


}

	`, projectID, name, *p.Compute.Enabled, *p.DiskGBEnabled, p.Compute.MaxInstanceSize)
}
