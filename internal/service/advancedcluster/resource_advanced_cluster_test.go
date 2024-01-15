package advancedcluster_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccClusterAdvancedCluster_basicTenant(t *testing.T) {
	var (
		cluster                matlas.AdvancedCluster
		resourceName           = "mongodbatlas_advanced_cluster.test"
		dataSourceName         = "data.mongodbatlas_advanced_cluster.test"
		dataSourceClustersName = "data.mongodbatlas_advanced_clusters.test"
		orgID                  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName            = acctest.RandomWithPrefix("test-acc")
		rName                  = acctest.RandomWithPrefix("test-acc")
		rNameUpdated           = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigTenant(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "termination_protection_enabled"),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "termination_protection_enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.termination_protection_enabled"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigTenant(orgID, projectName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(dataSourceName, "termination_protection_enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.termination_protection_enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: acc.ImportStateClusterIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccClusterAdvancedCluster_singleProvider(t *testing.T) {
	var (
		cluster        matlas.AdvancedCluster
		resourceName   = "mongodbatlas_advanced_cluster.test"
		dataSourceName = fmt.Sprintf("data.%s", resourceName)
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
		rName          = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProvider(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttrWith(dataSourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs", "retain_backups_enabled"},
			},
		},
	})
}

func TestAccClusterAdvancedCluster_multicloud(t *testing.T) {
	var (
		cluster                matlas.AdvancedCluster
		resourceName           = "mongodbatlas_advanced_cluster.test"
		dataSourceName         = "data.mongodbatlas_advanced_cluster.test"
		dataSourceClustersName = "data.mongodbatlas_advanced_clusters.test"
		orgID                  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName            = acctest.RandomWithPrefix("test-acc")
		rName                  = acctest.RandomWithPrefix("test-acc")
		rNameUpdated           = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.name"),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(orgID, projectName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.name"),
					resource.TestCheckResourceAttr(dataSourceName, "name", rNameUpdated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs", "retain_backups_enabled"},
			},
		},
	})
}

func TestAccClusterAdvancedCluster_multicloudSharded(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		rName        = acctest.RandomWithPrefix("test-acc")
		rNameUpdated = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloudSharded(orgID, projectName, rName),
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
				Config: testAccMongoDBAtlasAdvancedClusterConfigMultiCloudSharded(orgID, projectName, rNameUpdated),
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
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
		},
	})
}

func TestAccClusterAdvancedCluster_UnpausedToPaused(t *testing.T) {
	acc.SkipTest(t)
	var (
		cluster             matlas.AdvancedCluster
		resourceName        = "mongodbatlas_advanced_cluster.test"
		orgID               = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName         = acctest.RandomWithPrefix("test-acc")
		rName               = acctest.RandomWithPrefix("test-acc")
		instanceSize        = "M10"
		anotherInstanceSize = "M20"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, false, instanceSize),
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
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, instanceSize),
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
				Config:      testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, anotherInstanceSize),
				ExpectError: regexp.MustCompile("CANNOT_UPDATE_PAUSED_CLUSTER"),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
		},
	})
}

func TestAccClusterAdvancedCluster_PausedToUnpaused(t *testing.T) {
	acc.SkipTest(t)
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		rName        = acctest.RandomWithPrefix("test-acc")
		instanceSize = "M10"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, instanceSize),
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
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, false, instanceSize),
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
				Config:      testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, instanceSize),
				ExpectError: regexp.MustCompile("CANNOT_PAUSE_RECENTLY_RESUMED_CLUSTER"),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, false, instanceSize),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
		},
	})
}

func TestAccClusterAdvancedCluster_advancedConf(t *testing.T) {
	var (
		cluster                matlas.AdvancedCluster
		resourceName           = "mongodbatlas_advanced_cluster.test"
		dataSourceName         = "data.mongodbatlas_advanced_cluster.test"
		dataSourceNameClusters = "data.mongodbatlas_advanced_clusters.test"
		orgID                  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName            = acctest.RandomWithPrefix("test-acc")
		rName                  = acctest.RandomWithPrefix("test-acc")
		rNameUpdated           = acctest.RandomWithPrefix("test-acc")
		processArgs            = &matlas.ProcessArgs{
			DefaultReadConcern:               "available",
			DefaultWriteConcern:              "1",
			FailIndexKeyTooLong:              pointy.Bool(false),
			JavascriptEnabled:                pointy.Bool(true),
			MinimumEnabledTLSProtocol:        "TLS1_1",
			NoTableScan:                      pointy.Bool(false),
			OplogSizeMB:                      pointy.Int64(1000),
			SampleRefreshIntervalBIConnector: pointy.Int64(310),
			SampleSizeBIConnector:            pointy.Int64(110),
			TransactionLifetimeLimitSeconds:  pointy.Int64(300),
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
			TransactionLifetimeLimitSeconds:  pointy.Int64(300),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigAdvancedConf(orgID, projectName, rName, processArgs),
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
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.transaction_lifetime_limit_seconds", "300"),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceNameClusters, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameClusters, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameClusters, "results.0.name"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigAdvancedConf(orgID, projectName, rNameUpdated, processArgsUpdated),
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
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.transaction_lifetime_limit_seconds", "300"),
					resource.TestCheckResourceAttr(dataSourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttrSet(dataSourceNameClusters, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameClusters, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameClusters, "results.0.name"),
				),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_DefaultWrite(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
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
			TransactionLifetimeLimitSeconds:  pointy.Int64(300),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigAdvancedConfDefaultWrite(orgID, projectName, rName, processArgs),
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
				Config: testAccMongoDBAtlasAdvancedClusterConfigAdvancedConfDefaultWrite(orgID, projectName, rNameUpdated, processArgsUpdated),
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

func TestAccClusterAdvancedClusterConfig_ReplicationSpecsAutoScaling(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
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
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAutoScaling(orgID, projectName, rName, autoScaling),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					testAccCheckMongoDBAtlasAdvancedClusterScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAutoScaling(orgID, projectName, rNameUpdated, autoScalingUpdated),
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

func TestAccClusterAdvancedClusterConfig_ReplicationSpecsAnalyticsAutoScaling(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
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
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(orgID, projectName, rName, autoScaling),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					testAccCheckMongoDBAtlasAdvancedClusterAnalyticsScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(orgID, projectName, rNameUpdated, autoScalingUpdated),
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

func TestAccClusterAdvancedCluster_WithTags(t *testing.T) {
	var (
		cluster                matlas.AdvancedCluster
		resourceName           = "mongodbatlas_advanced_cluster.test"
		dataSourceName         = "data.mongodbatlas_advanced_cluster.test"
		dataSourceClustersName = "data.mongodbatlas_advanced_clusters.test"
		orgID                  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName            = acctest.RandomWithPrefix("test-acc")
		rName                  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigWithTags(orgID, projectName, rName, []matlas.Tag{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.tags.#", "0"),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigWithTags(orgID, projectName, rName, []matlas.Tag{
					{
						Key:   "key 1",
						Value: "value 1",
					},
					{
						Key:   "key 2",
						Value: "value 2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceClustersName, "results.0.tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceClustersName, "results.0.tags.*", acc.ClusterTagsMap2),
				),
			},
			{
				Config: testAccMongoDBAtlasAdvancedClusterConfigWithTags(orgID, projectName, rName, []matlas.Tag{
					{
						Key:   "key 3",
						Value: "value 3",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceClustersName, "results.0.tags.*", acc.ClusterTagsMap3),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName string, cluster *matlas.AdvancedCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestMongoDBClient.(*config.MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

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

func testAccMongoDBAtlasAdvancedClusterConfigTenant(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = mongodbatlas_project.cluster_project.id
  name         = %[3]q
  cluster_type = "REPLICASET"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M5"
      }
      provider_name         = "TENANT"
      backing_provider_name = "AWS"
      region_name           = "EU_WEST_1"
      priority              = 7
    }
  }
}

data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, orgID, projectName, name)
}

func testAccMongoDBAtlasAdvancedClusterConfigWithTags(orgID, projectName, name string, tags []matlas.Tag) string {
	var tagsConf string
	for _, label := range tags {
		tagsConf += fmt.Sprintf(`
			tags {
				key   = "%s"
				value = "%s"
			}
		`, label.Key, label.Value)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
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
			}

			%[4]s
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
		}
	`, orgID, projectName, name, tagsConf)
}

func testAccMongoDBAtlasAdvancedClusterConfigSingleProvider(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = mongodbatlas_project.cluster_project.id
  name         = %[3]q
  cluster_type = "REPLICASET"
  retain_backups_enabled = "true"

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
      region_name   = "EU_WEST_1"
    }
  }
}
data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

	`, orgID, projectName, name)
}

func testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = mongodbatlas_project.cluster_project.id
  name         = %[3]q
  cluster_type = "REPLICASET"
  retain_backups_enabled = false

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
      region_name   = "EU_WEST_1"
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

data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, orgID, projectName, name)
}

func testAccMongoDBAtlasAdvancedClusterConfigMultiCloudSharded(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = mongodbatlas_project.cluster_project.id
  name         = %[3]q
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
      region_name   = "EU_WEST_1"
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
	`, orgID, projectName, name)
}

func testAccMongoDBAtlasAdvancedClusterConfigSingleProviderPaused(orgID, projectName, name string, paused bool, instanceSize string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = mongodbatlas_project.cluster_project.id
  name         = %[3]q
  cluster_type = "REPLICASET"
  paused       = %[4]t

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = %[5]q
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }
}
	`, orgID, projectName, name, paused, instanceSize)
}

func testAccMongoDBAtlasAdvancedClusterConfigAdvancedConf(orgID, projectName, name string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = mongodbatlas_project.cluster_project.id
  name                   = %[3]q
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
      region_name   = "EU_WEST_1"
    }
  }

  advanced_configuration  {
    fail_index_key_too_long              = %[4]t
    javascript_enabled                   = %[5]t
    minimum_enabled_tls_protocol         = %[6]q
    no_table_scan                        = %[7]t
    oplog_size_mb                        = %[8]d
    sample_size_bi_connector			 = %[9]d
    sample_refresh_interval_bi_connector = %[10]d
	transaction_lifetime_limit_seconds   = %[11]d
  }
}
data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}

	`, orgID, projectName, name,
		*p.FailIndexKeyTooLong, *p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector, *p.TransactionLifetimeLimitSeconds)
}

func testAccMongoDBAtlasAdvancedClusterConfigAdvancedConfDefaultWrite(orgID, projectName, name string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = mongodbatlas_project.cluster_project.id
  name                   = %[3]q
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
      region_name   = "EU_WEST_1"
    }
  }

  advanced_configuration  {
    javascript_enabled                   = %[4]t
    minimum_enabled_tls_protocol         = %[5]q
    no_table_scan                        = %[6]t
    oplog_size_mb                        = %[7]d
    sample_size_bi_connector			 = %[8]d
    sample_refresh_interval_bi_connector = %[9]d
    default_read_concern                 = %[10]q
    default_write_concern                = %[11]q
  }
}

	`, orgID, projectName, name, *p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector, p.DefaultReadConcern, p.DefaultWriteConcern)
}

func testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAutoScaling(orgID, projectName, name string, p *matlas.AutoScaling) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = mongodbatlas_project.cluster_project.id
  name                   = %[3]q
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
        compute_enabled = %[4]t
        disk_gb_enabled = %[5]t
		compute_max_instance_size = %[6]q
	  }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }


}

	`, orgID, projectName, name, *p.Compute.Enabled, *p.DiskGBEnabled, p.Compute.MaxInstanceSize)
}

func testAccMongoDBAtlasAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(orgID, projectName, name string, p *matlas.AutoScaling) string {
	return fmt.Sprintf(`

resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}

resource "mongodbatlas_advanced_cluster" "test" {
  project_id             = mongodbatlas_project.cluster_project.id
  name                   = %[3]q
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
        compute_enabled = %[4]t
        disk_gb_enabled = %[5]t
		compute_max_instance_size = %[6]q
	  }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }


}

	`, orgID, projectName, name, *p.Compute.Enabled, *p.DiskGBEnabled, p.Compute.MaxInstanceSize)
}
