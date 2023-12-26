package advancedcluster_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
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
				Config: testAccAdvancedClusterConfigTenant(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForTenantConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				),
			},
			{
				Config: testAccAdvancedClusterConfigTenant(orgID, projectName, rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForTenantConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rNameUpdated)...,
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

func TestAccClusterAdvancedCluster_tenantUpgrade(t *testing.T) {
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
				Config: testAccAdvancedClusterConfigTenant(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForTenantConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				),
			},
			{
				Config: testAccAdvancedClusterConfigTenantUpgrade(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),

					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "termination_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_provider", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "state_name", "IDLE"),
					resource.TestCheckResourceAttr(resourceName, "version_release_system", "LTS"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "disk_size_gb"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.read_preference", "secondary"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_strings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_strings.0.standard_srv"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_strings.0.standard"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.priority", "7"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.region_name", "US_EAST_1"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.0.electable_specs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M10")),
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

func TestAccClusterAdvancedCluster_singleProviderToMultiCloud(t *testing.T) {
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
				Config: testAccAdvancedClusterConfigSingleProvider(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForSingleProviderConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				),
			},
			{
				Config: testAccAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForMultiCloudConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
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
				Config: testAccAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForMultiCloudConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				),
			},
			{
				Config: testAccAdvancedClusterConfigMultiCloud(orgID, projectName, rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForMultiCloudConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rNameUpdated)...,
				),
			},
			{
				Config: testAccAdvancedClusterConfigMultiCloud(orgID, projectName, rNameUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
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
				Config: testAccAdvancedClusterConfigMultiCloudSharded(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsMulticloudSharded(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName, false)...,
				),
			},
			{
				Config: testAccAdvancedClusterConfigMultiCloudSharded(orgID, projectName, rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsMulticloudSharded(&cluster, resourceName, dataSourceName, dataSourceClustersName, rNameUpdated, false)...,
				),
			},
			{
				Config: testAccAdvancedClusterConfigMultiCloudShardedUpdated(orgID, projectName, rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsMulticloudSharded(&cluster, resourceName, dataSourceName, dataSourceClustersName, rNameUpdated, true)...,
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

func TestAccClusterAdvancedCluster_geoSharded(t *testing.T) {
	var (
		cluster                matlas.AdvancedCluster
		resourceName           = "mongodbatlas_advanced_cluster.test"
		dataSourceName         = "data.mongodbatlas_advanced_cluster.test"
		dataSourceClustersName = "data.mongodbatlas_advanced_clusters.test"
		orgID                  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName            = acctest.RandomWithPrefix("test-acc")
		rName                  = acctest.RandomWithPrefix("test-acc")
		// rNameUpdated           = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccAdvancedClusterConfigGeoSharded(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// testAccCheckAdvancedClusterExists(resourceName, &cluster),
					// testAccCheckAdvancedClusterAttributes(&cluster, rName),

					testFuncsForGeoshardedConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName, false)...,
				),
			},
			{
				Config: testAccAdvancedClusterConfigGeoShardedUpdated(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// testAccCheckAdvancedClusterExists(resourceName, &cluster),
					// testAccCheckAdvancedClusterAttributes(&cluster, rName),
					testFuncsForGeoshardedConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName, true)...,
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

func TestAccClusterAdvancedCluster_unpausedToPaused(t *testing.T) {
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
				Config: testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, false, instanceSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
				),
			},
			{
				Config: testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, instanceSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
				),
			},
			{
				Config:      testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, anotherInstanceSize),
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

func TestAccClusterAdvancedCluster_pausedToUnpaused(t *testing.T) {
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
				Config: testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, instanceSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
				),
			},
			{
				Config: testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, false, instanceSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
				),
			},
			{
				Config:      testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, instanceSize),
				ExpectError: regexp.MustCompile("CANNOT_PAUSE_RECENTLY_RESUMED_CLUSTER"),
			},
			{
				Config: testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, false, instanceSize),
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
			// DefaultReadConcern:  "available",
			// DefaultWriteConcern: "1",
			// FailIndexKeyTooLong:              pointy.Bool(false),
			JavascriptEnabled:                pointy.Bool(true),
			MinimumEnabledTLSProtocol:        "TLS1_1",
			NoTableScan:                      pointy.Bool(false),
			OplogSizeMB:                      pointy.Int64(1000),
			OplogMinRetentionHours:           pointy.Float64(2.0),
			SampleRefreshIntervalBIConnector: pointy.Int64(310),
			SampleSizeBIConnector:            pointy.Int64(110),
			TransactionLifetimeLimitSeconds:  pointy.Int64(300),
		}
		processArgsUpdated = &matlas.ProcessArgs{
			// DefaultReadConcern:  "available",
			// DefaultWriteConcern: "0",
			// FailIndexKeyTooLong:              pointy.Bool(false),
			JavascriptEnabled:         pointy.Bool(true),
			MinimumEnabledTLSProtocol: "TLS1_2",
			NoTableScan:               pointy.Bool(false),
			OplogSizeMB:               pointy.Int64(1000),
			// OplogMinRetentionHours:           pointy.Float64(0.0),
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
				Config: testAccAdvancedClusterConfigAdvancedConf(orgID, projectName, rName, processArgs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.default_read_concern"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.default_write_concern"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_min_retention_hours", "2"),
					// resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					// resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_min_retention_hours", "2"),
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
				Config: testAccAdvancedClusterConfigAdvancedConfNoOplogHrs(orgID, projectName, rNameUpdated, processArgsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rNameUpdated),

					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.oplog_min_retention_hours"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.default_read_concern"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.default_write_concern"),
					// resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					// resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "1"),
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

func TestAccClusterAdvancedCluster_defaultWrite(t *testing.T) {
	var (
		cluster                matlas.AdvancedCluster
		resourceName           = "mongodbatlas_advanced_cluster.test"
		dataSourceName         = "data.mongodbatlas_advanced_cluster.test"
		dataSourceClustersName = "data.mongodbatlas_advanced_clusters.test"
		orgID                  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName            = acctest.RandomWithPrefix("test-acc")
		rName                  = acctest.RandomWithPrefix("test-acc")
		rNameUpdated           = acctest.RandomWithPrefix("test-acc")
		processArgs            = &matlas.ProcessArgs{
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
				Config: testAccAdvancedClusterConfigAdvancedConfDefaultWrite(orgID, projectName, rName, processArgs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),

					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.default_write_concern", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),

					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.advanced_configuration.0.default_write_concern", "1"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.advanced_configuration.0.sample_size_bi_connector", "110"),
				),
			},
			{
				Config: testAccAdvancedClusterConfigAdvancedConfDefaultWrite(orgID, projectName, rNameUpdated, processArgsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "majority"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long"),
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

func TestAccClusterAdvancedClusterConfig_replicationSpecsAutoScaling(t *testing.T) {
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
				Config: testAccAdvancedClusterConfigReplicationSpecsAutoScaling(orgID, projectName, rName, autoScaling),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.disk_gb_enabled", "true"),

					testAccCheckAdvancedClusterScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				Config: testAccAdvancedClusterConfigReplicationSpecsAutoScaling(orgID, projectName, rNameUpdated, autoScalingUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.disk_gb_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_max_instance_size", "M20"),
					testAccCheckAdvancedClusterScaling(&cluster, *autoScalingUpdated.Compute.Enabled),
				),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_replicationSpecsAnalyticsAutoScaling(t *testing.T) {
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
				Config: testAccAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(orgID, projectName, rName, autoScaling),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.disk_gb_enabled", "true"),
					// resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_max_instance_size"),
					testAccCheckAdvancedClusterAnalyticsScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				Config: testAccAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(orgID, projectName, rNameUpdated, autoScalingUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.disk_gb_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_max_instance_size", "M20"),
					testAccCheckAdvancedClusterAnalyticsScaling(&cluster, *autoScalingUpdated.Compute.Enabled),
				),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_withTags(t *testing.T) {
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
				Config: testAccAdvancedClusterConfigWithTags(orgID, projectName, rName, []matlas.Tag{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.tags.#", "0"),
				),
			},
			{
				Config: testAccAdvancedClusterConfigWithTags(orgID, projectName, rName, []matlas.Tag{
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
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
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
				Config: testAccAdvancedClusterConfigWithTags(orgID, projectName, rName, []matlas.Tag{
					{
						Key:   "key 3",
						Value: "value 3",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
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

func TestAccClusterAdvancedCluster_withLabels(t *testing.T) {
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
				Config: testAccAdvancedClusterConfigWithLabels(orgID, projectName, rName, []matlas.Label{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "labels.#", "0"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.labels.#", "0"),
				),
			},
			{
				Config: testAccAdvancedClusterConfigWithLabels(orgID, projectName, rName, []matlas.Label{
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
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "labels.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "labels.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceName, "labels.#", "2"), // check if data source returnes all labels including default
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "labels.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "labels.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.labels.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceClustersName, "results.0.labels.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceClustersName, "results.0.labels.*", acc.ClusterTagsMap2),
				),
			},
			{
				Config: testAccAdvancedClusterConfigWithLabels(orgID, projectName, rName, []matlas.Label{
					{
						Key:   "key 3",
						Value: "value 3",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "labels.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceName, "labels.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "labels.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.labels.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceClustersName, "results.0.labels.*", acc.ClusterTagsMap3),
				),
			},
		},
	})
}

func testFuncsForMultiCloudConfig(cluster *matlas.AdvancedCluster, resourceName, dataSourceName, dataSourceClustersName, rName string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		testAccCheckAdvancedClusterExists(resourceName, cluster),
		testAccCheckAdvancedClusterAttributes(cluster, rName),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "name", rName),
		resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
		resource.TestCheckResourceAttr(resourceName, "bi_connector_config.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.read_preference", "secondary"),
		resource.TestCheckResourceAttr(resourceName, "termination_protection_enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_provider", "NONE"),
		resource.TestCheckResourceAttr(resourceName, "state_name", "IDLE"),
		resource.TestCheckResourceAttr(resourceName, "version_release_system", "LTS"),
		resource.TestCheckResourceAttr(resourceName, "pit_enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "paused", "false"),
		resource.TestCheckResourceAttrSet(resourceName, "disk_size_gb"),
		resource.TestCheckResourceAttrSet(resourceName, "connection_strings.#"),
		resource.TestCheckResourceAttrSet(resourceName, "connection_strings.0.standard_srv"),
		resource.TestCheckResourceAttrSet(resourceName, "connection_strings.0.standard"),
		resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
		resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.container_id.%", "2"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.zone_name", advancedcluster.DefaultZoneName),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.priority", "7"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.provider_name", "AWS"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.region_name", "US_EAST_1"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.node_count", "1"),
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.node_count", "3"),
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.priority", "6"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.provider_name", "GCP"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.region_name", "NORTH_AMERICA_NORTHEAST_1"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(
			resourceName,
			"replication_specs.*",
			map[string]string{
				"zone_name":  advancedcluster.DefaultZoneName,
				"num_shards": "1",
			},
		),

		resource.TestCheckResourceAttr(dataSourceName, "name", rName),
		resource.TestCheckResourceAttr(dataSourceName, "termination_protection_enabled", "false"),
		resource.TestCheckResourceAttr(dataSourceName, "cluster_type", "REPLICASET"),
		resource.TestCheckResourceAttr(dataSourceName, "encryption_at_rest_provider", "NONE"),
		resource.TestCheckResourceAttr(dataSourceName, "state_name", "IDLE"),
		resource.TestCheckResourceAttr(dataSourceName, "version_release_system", "LTS"),
		resource.TestCheckResourceAttr(dataSourceName, "pit_enabled", "false"),
		resource.TestCheckResourceAttr(dataSourceName, "paused", "false"),
		resource.TestCheckResourceAttrSet(dataSourceName, "disk_size_gb"),
		resource.TestCheckResourceAttr(dataSourceName, "bi_connector_config.0.enabled", "true"),
		resource.TestCheckResourceAttr(dataSourceName, "bi_connector_config.0.read_preference", "secondary"),
		resource.TestCheckResourceAttrSet(dataSourceName, "connection_strings.#"),
		resource.TestCheckResourceAttrSet(dataSourceName, "connection_strings.0.standard_srv"),
		resource.TestCheckResourceAttrSet(dataSourceName, "connection_strings.0.standard"),
		// resource.TestCheckResourceAttrSet(dataSourceName, "labels.#"),
		resource.TestCheckResourceAttrSet(dataSourceName, "replication_specs.#"),

		resource.TestCheckTypeSetElemNestedAttrs(
			dataSourceName,
			"replication_specs.*",
			map[string]string{
				"zone_name":  advancedcluster.DefaultZoneName,
				"num_shards": "1",
			},
		),
		resource.TestCheckResourceAttrSet(dataSourceName, "replication_specs.0.region_configs.#"),
		resource.TestCheckTypeSetElemNestedAttrs(
			dataSourceName,
			"replication_specs.0.region_configs.*",
			map[string]string{
				"priority":                        "6",
				"provider_name":                   "GCP",
				"region_name":                     "NORTH_AMERICA_NORTHEAST_1",
				"electable_specs.0.instance_size": "M10",
				"electable_specs.0.node_count":    "2",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			dataSourceName,
			"replication_specs.0.region_configs.*",
			map[string]string{
				"priority":                        "7",
				"provider_name":                   "AWS",
				"region_name":                     "US_EAST_1",
				"analytics_specs.0.instance_size": "M10",
				"analytics_specs.0.node_count":    "1",
				"electable_specs.0.instance_size": "M10",
				"electable_specs.0.node_count":    "3",
			},
		),

		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.#"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.#"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.name"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.termination_protection_enabled"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.cluster_type", "REPLICASET"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.encryption_at_rest_provider", "NONE"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.state_name", "IDLE"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.version_release_system", "LTS"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.pit_enabled", "false"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.paused", "false"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.disk_size_gb"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.bi_connector_config.0.enabled", "true"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.bi_connector_config.0.read_preference", "secondary"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.connection_strings.#"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.connection_strings.0.standard_srv"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.connection_strings.0.standard"),
		// resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.labels.#", "1"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.#"),
		resource.TestCheckTypeSetElemNestedAttrs(
			dataSourceClustersName,
			"results.0.replication_specs.*",
			map[string]string{
				"zone_name":  advancedcluster.DefaultZoneName,
				"num_shards": "1",
			},
		),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.0.region_configs.#"),
		resource.TestCheckTypeSetElemNestedAttrs(
			dataSourceClustersName,
			"results.0.replication_specs.0.region_configs.*",
			map[string]string{
				"priority":      "6",
				"provider_name": "GCP",
				"region_name":   "NORTH_AMERICA_NORTHEAST_1",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			dataSourceClustersName,
			"results.0.replication_specs.0.region_configs.*",
			map[string]string{
				"priority":      "7",
				"provider_name": "AWS",
				"region_name":   "US_EAST_1",
			},
		),
	}
}

func testAccCheckAdvancedClusterExists(resourceName string, cluster *matlas.AdvancedCluster) resource.TestCheckFunc {
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

func testAccCheckAdvancedClusterAttributes(cluster *matlas.AdvancedCluster, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cluster.Name != name {
			return fmt.Errorf("bad name: %s", cluster.Name)
		}

		return nil
	}
}

func testAccCheckAdvancedClusterScaling(cluster *matlas.AdvancedCluster, computeEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *cluster.ReplicationSpecs[0].RegionConfigs[0].AutoScaling.Compute.Enabled != computeEnabled {
			return fmt.Errorf("compute_enabled: %d", cluster.ReplicationSpecs[0].RegionConfigs[0].AutoScaling.Compute.Enabled)
		}

		return nil
	}
}

func testAccCheckAdvancedClusterAnalyticsScaling(cluster *matlas.AdvancedCluster, computeEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *cluster.ReplicationSpecs[0].RegionConfigs[0].AnalyticsAutoScaling.Compute.Enabled != computeEnabled {
			return fmt.Errorf("compute_enabled: %d", cluster.ReplicationSpecs[0].RegionConfigs[0].AnalyticsAutoScaling.Compute.Enabled)
		}

		return nil
	}
}

func testAccAdvancedClusterConfigTenant(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = mongodbatlas_project.cluster_project.id
  name         = %[3]q
  cluster_type = "REPLICASET"

  replication_specs = [{
    region_configs = [{
      electable_specs = [{
        instance_size = "M5"
      }]
      provider_name         = "TENANT"
      backing_provider_name = "AWS"
      region_name           = "US_EAST_1"
      priority              = 7
    }]
  }]
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

func testFuncsForTenantConfig(cluster *matlas.AdvancedCluster, resourceName, dataSourceName, dataSourceClustersName, rName string) []resource.TestCheckFunc {
	res := []resource.TestCheckFunc{
		testAccCheckAdvancedClusterExists(resourceName, cluster),
		testAccCheckAdvancedClusterAttributes(cluster, rName),

		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.#"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.#"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.name"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.termination_protection_enabled"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.cluster_type", "REPLICASET"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.backup_enabled", "true"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.encryption_at_rest_provider", "NONE"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.state_name", "IDLE"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.version_release_system", "LTS"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.pit_enabled", "false"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.paused", "false"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.disk_size_gb"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.bi_connector_config.0.enabled", "false"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.bi_connector_config.0.read_preference", "secondary"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.connection_strings.#"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.connection_strings.0.standard_srv"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.connection_strings.0.standard"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.labels.#", "0"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.#"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.0.region_configs.#"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.0.region_configs.0.priority"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.replication_specs.0.region_configs.0.provider_name", "TENANT"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.replication_specs.0.region_configs.0.region_name", "US_EAST_1"),
		resource.TestCheckResourceAttrSet(dataSourceClustersName, "results.0.replication_specs.0.region_configs.0.electable_specs.#"),
		resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M5"),
	}

	res = append(res, testFuncsForTenantConfigGeneric(resourceName, rName)...)
	res = append(res, testFuncsForTenantConfigGeneric(dataSourceName, rName)...)
	return res
}

func testFuncsForTenantConfigGeneric(sourceName, rName string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(sourceName, "project_id"),
		resource.TestCheckResourceAttr(sourceName, "name", rName),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.#"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.0.region_configs.#"),
		resource.TestCheckResourceAttr(sourceName, "termination_protection_enabled", "false"),

		resource.TestCheckResourceAttr(sourceName, "cluster_type", "REPLICASET"),
		resource.TestCheckResourceAttr(sourceName, "backup_enabled", "true"),
		resource.TestCheckResourceAttr(sourceName, "encryption_at_rest_provider", "NONE"),
		resource.TestCheckResourceAttr(sourceName, "state_name", "IDLE"),
		resource.TestCheckResourceAttr(sourceName, "version_release_system", "LTS"),
		resource.TestCheckResourceAttr(sourceName, "pit_enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "paused", "false"),
		resource.TestCheckResourceAttrSet(sourceName, "disk_size_gb"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.0.enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.0.read_preference", "secondary"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.#"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.0.standard_srv"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.0.standard"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.#"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.zone_name", advancedcluster.DefaultZoneName),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.0.region_configs.#"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.priority", "7"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.provider_name", "TENANT"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.backing_provider_name", "AWS"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.region_name", "US_EAST_1"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.0.region_configs.0.electable_specs.#"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M5"),
	}
}

func testAccAdvancedClusterConfigTenantUpgrade(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = mongodbatlas_project.cluster_project.id
  name         = %[3]q
  cluster_type = "REPLICASET"

  replication_specs = [{
    region_configs = [{
      electable_specs = [{
        instance_size = "M10"
      }]
      provider_name         = "AWS"
      region_name           = "US_EAST_1"
      priority              = 7
    }]
  }]
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

func testAccAdvancedClusterConfigWithLabels(orgID, projectName, name string, tags []matlas.Label) string {
	tagsConf := "labels = [%s]"
	var tagsArr string
	for _, label := range tags {
		tagsArr += fmt.Sprintf(`
		{
			key   = "%s"
			value = "%s"
		},
	`, label.Key, label.Value)
	}

	if len(tags) > 0 {
		tagsArr = tagsArr[:len(tagsArr)-1]
		tagsConf = fmt.Sprintf(tagsConf, tagsArr)
	} else {
		tagsConf = ""
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

			replication_specs = [{
				region_configs = [{
					electable_specs = [{
						instance_size = "M10"
						node_count    = 3
					}]
					analytics_specs = [{
						instance_size = "M10"
						node_count    = 1
					}]
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
				}]
			}]

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

func testAccAdvancedClusterConfigWithTags(orgID, projectName, name string, tags []matlas.Tag) string {
	tagsConf := "tags = [%s]"
	var tagsArr string
	for _, label := range tags {
		tagsArr += fmt.Sprintf(`
		{
			key   = "%s"
			value = "%s"
		},
	`, label.Key, label.Value)
	}

	if len(tags) > 0 {
		tagsArr = tagsArr[:len(tagsArr)-1]
		tagsConf = fmt.Sprintf(tagsConf, tagsArr)
	} else {
		tagsConf = ""
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

			replication_specs = [{
				region_configs = [{
					electable_specs = [{
						instance_size = "M10"
						node_count    = 3
					}]
					analytics_specs = [{
						instance_size = "M10"
						node_count    = 1
					}]
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
				}]
			}]

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

func testAccAdvancedClusterConfigSingleProvider(orgID, projectName, name string) string {
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

  replication_specs = [{
    region_configs = [{
      electable_specs = [{
        instance_size = "M10"
        node_count    = 3
      }]
      analytics_specs = [{
        instance_size = "M10"
        node_count    = 1
      }]
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }]
  }]
}
data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

	`, orgID, projectName, name)
}

func testFuncsForSingleProviderConfig(cluster *matlas.AdvancedCluster, resourceName, dataSourceName, dataSourceClustersName, rName string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		testAccCheckAdvancedClusterExists(resourceName, cluster),
		testAccCheckAdvancedClusterAttributes(cluster, rName),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "name", rName),
		resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
		resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "bi_connector_config.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.read_preference", "secondary"),

		resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
		resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.container_id.%", "1"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.zone_name", advancedcluster.DefaultZoneName),

		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.node_count", "3"),
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.node_count", "1"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.priority", "7"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.provider_name", "AWS"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.region_name", "US_EAST_1"),

		resource.TestCheckResourceAttrWith(dataSourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
	}
}

func testAccAdvancedClusterConfigMultiCloud(orgID, projectName, name string) string {
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

	bi_connector_config = [{
		enabled = true
	}]
  
	replication_specs = [{
	  region_configs = [{
		electable_specs = [{
		  instance_size = "M10"
		  node_count    = 3
		}]
		analytics_specs = [{
		  instance_size = "M10"
		  node_count    = 1
		}]
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  },
	  {
		electable_specs = [{
		  instance_size = "M10"
		  node_count    = 2
		}]
		provider_name = "GCP"
		priority      = 6
		region_name   = "NORTH_AMERICA_NORTHEAST_1"
	  }]
	}]
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

func testAccAdvancedClusterConfigMultiCloudSharded(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
	project_id   = mongodbatlas_project.cluster_project.id
	name         = %[3]q
	cluster_type = "SHARDED"
  
	replication_specs = [{
	  num_shards = 1
	  region_configs = [{
		electable_specs = [{
		  instance_size = "M30"
		  node_count    = 3
		}]
		analytics_specs = [{
		  instance_size = "M30"
		  node_count    = 1
		}]
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  },
	  {
		electable_specs = [{
		  instance_size = "M30"
		  node_count    = 2
		}]
		provider_name = "AZURE"
		priority      = 6
		region_name   = "US_EAST_2"
	  }]
	}]
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

func testAccAdvancedClusterConfigGeoSharded(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
	project_id   = mongodbatlas_project.cluster_project.id
	name         = %[3]q
  
	cluster_type   = "GEOSHARDED"
	backup_enabled = true
  
	replication_specs =[{ # zone n1
	  zone_name  = "zone n1"
	  num_shards = 3 # 3-shard Multi-Cloud Cluster
  
	  region_configs =[{ # shard n1 
		electable_specs =[{
		  instance_size = "M10"
		  node_count    = 3
		}]
		analytics_specs =[{
		  instance_size = "M10"
		   node_count    = 1
		}]
		analytics_auto_scaling =[{
		  compute_enabled = true
		  compute_scale_down_enabled = false
		  compute_max_instance_size = "M20"
		}]
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  },
   { # shard n3
		electable_specs =[{
		  instance_size = "M10"
		  node_count    = 2
		}]
		analytics_specs =[{
		  instance_size = "M10"
		  node_count    = 1
		}]
		provider_name = "GCP"
		priority      = 6
		region_name   = "US_EAST_4"
	  }
	 ]
	}
	, { # zone n2
	  zone_name  = "zone n2"
	  num_shards = 2 # 2-shard Multi-Cloud Cluster
  
	  region_configs =[{ # shard n1 
		electable_specs =[{
		  instance_size = "M10"
		  node_count    = 3
		}]
		analytics_specs =[{
		  instance_size = "M10"
		  node_count    = 1
		}]
		provider_name = "AWS"
		priority      = 7
		region_name   = "EU_WEST_1"
	  }]
	}]
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

func testAccAdvancedClusterConfigGeoShardedUpdated(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
	project_id   = mongodbatlas_project.cluster_project.id
	name         = %[3]q

	cluster_type   = "GEOSHARDED"
	backup_enabled = true
  
	replication_specs =[{ # zone n1
		zone_name  = "zone n1"
		num_shards = 3 # 3-shard Multi-Cloud Cluster
	
		region_configs =[{ 
		  electable_specs =[{
			instance_size = "M10"
			node_count    = 3
		  }]
		  analytics_specs =[{
			instance_size = "M10"
			 node_count    = 1
		  }]
		  analytics_auto_scaling =[{
			compute_enabled = true
			compute_scale_down_enabled = false
			compute_max_instance_size = "M20"
		  }]
		  provider_name = "AWS"
		  priority      = 7
		  region_name   = "US_EAST_1"
		},
		{ 
		  electable_specs =[{
			instance_size = "M10"
			node_count    = 2
		  }]
		  analytics_specs =[{
			instance_size = "M10"
			node_count    = 1
		  }]
		  provider_name = "GCP"
		  priority      = 6
		  region_name   = "US_EAST_4"
		}]
	  }, { 
		# zone_name  = "zone n2"
		num_shards = 2 # 2-shard Multi-Cloud Cluster
	
		region_configs =[{ # shard n1 
		  electable_specs =[{
			instance_size = "M10"
			node_count    = 3
		  }]
		  analytics_specs =[{
			instance_size = "M10"
			node_count    = 1
		  }]
		  provider_name = "AWS"
		  priority      = 7
		  region_name   = "EU_WEST_1"
		},{ 
			 electable_specs =[{
			   instance_size = "M10"
			   node_count    = 2
			 }]
			 analytics_specs =[{
			   instance_size = "M10"
			   node_count    = 1
			 }]
			 provider_name = "AZURE"
			 priority      = 6
			 region_name   = "EUROPE_NORTH"
		   }]
	  }]
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

func testFuncsForGeoshardedConfig(cluster *matlas.AdvancedCluster, resourceName, dataSourceName, dataSourceClustersName, rName string, isSpecUpdateTest bool) []resource.TestCheckFunc {
	res := []resource.TestCheckFunc{
		testAccCheckAdvancedClusterExists(resourceName, cluster),
		testAccCheckAdvancedClusterAttributes(cluster, rName),

		// resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.priority", "6"),
		// resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.provider_name", "AZURE"),
		// resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.region_name", "US_EAST_2"),
		// resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M30"),
		// resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
		// resource.TestCheckNoResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops"),
		// resource.TestCheckNoResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.ebs_volume_type"),
	}

	res = append(res, testFuncsGeoshardedGeneric(resourceName, rName)...)
	res = append(res, testFuncsDSGeoshardedGeneric(dataSourceName, rName)...)

	if isSpecUpdateTest {
		updateTests := []resource.TestCheckFunc{
			resource.TestCheckResourceAttrSet(resourceName, "replication_specs.1.region_configs.#"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.0.container_id.%", "2"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.container_id.%", "2"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.zone_name", advancedcluster.DefaultZoneName),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.num_shards", "2"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.region_configs.1.priority", "6"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.region_configs.1.provider_name", "AZURE"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.region_configs.1.region_name", "EUROPE_NORTH"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.region_configs.1.electable_specs.0.instance_size", "M10"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.region_configs.1.electable_specs.0.node_count", "2"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.region_configs.1.analytics_specs.0.instance_size", "M10"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.1.region_configs.1.analytics_specs.0.node_count", "1"),
		}
		res = append(res, updateTests...)
	}

	return res
}

func testFuncsDSGeoshardedGeneric(sourceName, rName string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(sourceName, "project_id"),
		resource.TestCheckResourceAttr(sourceName, "name", rName),
		resource.TestCheckResourceAttr(sourceName, "cluster_type", "GEOSHARDED"),
		resource.TestCheckResourceAttr(sourceName, "backup_enabled", "true"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.#", "1"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.0.enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.0.read_preference", "secondary"),
		resource.TestCheckResourceAttr(sourceName, "termination_protection_enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "encryption_at_rest_provider", "NONE"),
		resource.TestCheckResourceAttr(sourceName, "state_name", "IDLE"),
		resource.TestCheckResourceAttr(sourceName, "version_release_system", "LTS"),
		resource.TestCheckResourceAttr(sourceName, "pit_enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "paused", "false"),
		resource.TestCheckResourceAttrSet(sourceName, "disk_size_gb"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.#"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.0.standard_srv"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.0.standard"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.#"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.0.region_configs.#"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.1.region_configs.#"),

		resource.TestCheckTypeSetElemNestedAttrs(
			sourceName,
			"replication_specs.*",
			map[string]string{
				// "container_id.%": "2",
				"zone_name": "zone n1",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			sourceName,
			"replication_specs.*",
			map[string]string{
				// "container_id.%": "1",
				// "zone_name":  advancedcluster.DefaultZoneName,
				"num_shards": "2",
			},
		),

		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.priority", "7"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.provider_name", "AWS"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.region_name", "US_EAST_1"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.analytics_specs.0.instance_size", "M10"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.analytics_specs.0.node_count", "1"),
		// resource.TestCheckResourceAttrWith(sourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M10"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.node_count", "3"),
		// resource.TestCheckResourceAttrWith(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.priority", "6"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.provider_name", "GCP"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.region_name", "US_EAST_4"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M10"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
		// resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops"),
		// resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.ebs_volume_type"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.analytics_specs.0.instance_size", "M10"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.analytics_specs.0.node_count", "1"),
		// resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.analytics_specs.0.disk_iops"),
		// resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.analytics_specs.0.ebs_volume_type"),

		// resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.priority", "7"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.provider_name", "AWS"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.region_name", "EU_WEST_1"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.electable_specs.0.instance_size", "M10"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.electable_specs.0.node_count", "3"),
		// resource.TestCheckResourceAttrWith(sourceName, "replication_specs.1.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		// resource.TestCheckResourceAttrSet(sourceName, "replication_specs.1.region_configs.0.electable_specs.0.ebs_volume_type"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.analytics_specs.0.instance_size", "M10"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.analytics_specs.0.node_count", "1"),
		// resource.TestCheckResourceAttrWith(sourceName, "replication_specs.1.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		// resource.TestCheckResourceAttrSet(sourceName, "replication_specs.1.region_configs.0.analytics_specs.0.ebs_volume_type"),
	}
}

func testFuncsGeoshardedGeneric(sourceName, rName string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(sourceName, "project_id"),
		resource.TestCheckResourceAttr(sourceName, "name", rName),
		resource.TestCheckResourceAttr(sourceName, "cluster_type", "GEOSHARDED"),
		resource.TestCheckResourceAttr(sourceName, "backup_enabled", "true"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.#", "1"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.0.enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.0.read_preference", "secondary"),
		resource.TestCheckResourceAttr(sourceName, "termination_protection_enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "encryption_at_rest_provider", "NONE"),
		resource.TestCheckResourceAttr(sourceName, "state_name", "IDLE"),
		resource.TestCheckResourceAttr(sourceName, "version_release_system", "LTS"),
		resource.TestCheckResourceAttr(sourceName, "pit_enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "paused", "false"),
		resource.TestCheckResourceAttrSet(sourceName, "disk_size_gb"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.#"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.0.standard_srv"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.0.standard"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.#"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.0.region_configs.#"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.0.container_id.%", "2"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.zone_name", "zone n1"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.priority", "7"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.provider_name", "AWS"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.region_name", "US_EAST_1"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.analytics_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.analytics_specs.0.node_count", "1"),
		resource.TestCheckResourceAttrWith(sourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.node_count", "3"),
		resource.TestCheckResourceAttrWith(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.priority", "6"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.provider_name", "GCP"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.region_name", "US_EAST_4"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
		resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops"),
		resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.ebs_volume_type"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.analytics_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.analytics_specs.0.node_count", "1"),
		resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.analytics_specs.0.disk_iops"),
		resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.analytics_specs.0.ebs_volume_type"),

		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.1.region_configs.#"),
		// resource.TestCheckResourceAttr(sourceName, "replication_specs.1.zone_name", advancedcluster.DefaultZoneName),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.1.num_shards", "2"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.priority", "7"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.provider_name", "AWS"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.region_name", "EU_WEST_1"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.electable_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.electable_specs.0.node_count", "3"),
		resource.TestCheckResourceAttrWith(sourceName, "replication_specs.1.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.1.region_configs.0.electable_specs.0.ebs_volume_type"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.analytics_specs.0.instance_size", "M10"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.1.region_configs.0.analytics_specs.0.node_count", "1"),
		resource.TestCheckResourceAttrWith(sourceName, "replication_specs.1.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.1.region_configs.0.analytics_specs.0.ebs_volume_type"),
	}
}

func testFuncsMulticloudSharded(cluster *matlas.AdvancedCluster, resourceName, dataSourceName, dataSourceClustersName, rName string, isSpecUpdateTest bool) []resource.TestCheckFunc {
	res := []resource.TestCheckFunc{
		testAccCheckAdvancedClusterExists(resourceName, cluster),
		testAccCheckAdvancedClusterAttributes(cluster, rName),

		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.priority", "6"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.provider_name", "AZURE"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.region_name", "US_EAST_2"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M30"),
		resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
		resource.TestCheckNoResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops"),
		resource.TestCheckNoResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.ebs_volume_type"),
	}

	res = append(res, testFuncsMulticloudShardedGeneric(resourceName, rName)...)
	res = append(res, testFuncsMulticloudShardedGeneric(dataSourceName, rName)...)

	if isSpecUpdateTest {
		updateTests := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "replication_specs.0.container_id.%", "3"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.2.priority", "5"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.2.provider_name", "GCP"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.2.region_name", "US_EAST_4"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.2.electable_specs.0.instance_size", "M30"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.2.electable_specs.0.node_count", "2"),
			resource.TestCheckNoResourceAttr(resourceName, "replication_specs.0.region_configs.2.electable_specs.0.disk_iops"),
			resource.TestCheckNoResourceAttr(resourceName, "replication_specs.0.region_configs.2.electable_specs.0.ebs_volume_type"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.2.analytics_specs.0.instance_size", "M30"),
			resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.2.analytics_specs.0.node_count", "1"),
			resource.TestCheckNoResourceAttr(resourceName, "replication_specs.0.region_configs.2.analytics_specs.0.disk_iops"),
			resource.TestCheckNoResourceAttr(resourceName, "replication_specs.0.region_configs.2.analytics_specs.0.ebs_volume_type"),
		}
		res = append(res, updateTests...)
	}

	return res
}

func testFuncsMulticloudShardedGeneric(sourceName, rName string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(sourceName, "project_id"),
		resource.TestCheckResourceAttr(sourceName, "name", rName),
		resource.TestCheckResourceAttr(sourceName, "cluster_type", "SHARDED"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.#", "1"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.0.enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "bi_connector_config.0.read_preference", "secondary"),
		resource.TestCheckResourceAttr(sourceName, "termination_protection_enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "encryption_at_rest_provider", "NONE"),
		resource.TestCheckResourceAttr(sourceName, "state_name", "IDLE"),
		resource.TestCheckResourceAttr(sourceName, "version_release_system", "LTS"),
		resource.TestCheckResourceAttr(sourceName, "pit_enabled", "false"),
		resource.TestCheckResourceAttr(sourceName, "paused", "false"),
		resource.TestCheckResourceAttrSet(sourceName, "disk_size_gb"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.#"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.0.standard_srv"),
		resource.TestCheckResourceAttrSet(sourceName, "connection_strings.0.standard"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.#"),
		resource.TestCheckResourceAttrSet(sourceName, "replication_specs.0.region_configs.#"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.zone_name", advancedcluster.DefaultZoneName),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.priority", "7"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.provider_name", "AWS"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.region_name", "US_EAST_1"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.analytics_specs.0.instance_size", "M30"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.analytics_specs.0.node_count", "1"),
		resource.TestCheckResourceAttrWith(sourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M30"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.node_count", "3"),
		resource.TestCheckResourceAttrWith(sourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.priority", "6"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.provider_name", "AZURE"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.region_name", "US_EAST_2"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M30"),
		resource.TestCheckResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
		resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops"),
		resource.TestCheckNoResourceAttr(sourceName, "replication_specs.0.region_configs.1.electable_specs.0.ebs_volume_type"),
	}
}

func testAccAdvancedClusterConfigMultiCloudShardedUpdated(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
	project_id   = mongodbatlas_project.cluster_project.id
	name         = %[3]q
	cluster_type = "SHARDED"
  
	replication_specs = [{
	  num_shards = 1
	  region_configs = [{
		electable_specs = [{
		  instance_size = "M30"
		  node_count    = 3
		}]
		analytics_specs = [{
		  instance_size = "M30"
		  node_count    = 1
		}]
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  },
	  {
		electable_specs = [{
		  instance_size = "M30"
		  node_count    = 2
		}]
		provider_name = "AZURE"
		priority      = 6
		region_name   = "US_EAST_2"
	  },
	  {
		electable_specs =[{
		  instance_size = "M30"
		  node_count    = 2
		}]
		analytics_specs =[{
		  instance_size = "M30"
		  node_count    = 1
		}]
		provider_name = "GCP"
		priority      = 5
		region_name   = "US_EAST_4"
	  }]
	}]
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

func testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, name string, paused bool, instanceSize string) string {
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

  replication_specs = [{
    region_configs = [{
      electable_specs = [{
        instance_size = %[5]q
        node_count    = 3
      }]
      analytics_specs = [{
        instance_size = "M10"
        node_count    = 1
      }]
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }]
  }]
}
	`, orgID, projectName, name, paused, instanceSize)
}

func testAccAdvancedClusterConfigAdvancedConf(orgID, projectName, name string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
	project_id             = mongodbatlas_project.cluster_project.id
	name                   = %[3]q
	cluster_type           = "REPLICASET"
  
	 replication_specs = [{
	  region_configs  = [{
		electable_specs = [{
		  instance_size = "M10"
		  node_count    = 3
		}]
		analytics_specs = [{
		  instance_size = "M10"
		  node_count    = 1
		}]
		auto_scaling = [{
			compute_enabled = true
			disk_gb_enabled = true
			compute_max_instance_size = "M20"
		   }]
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  }]
	}]
  
	advanced_configuration  = [{
	  javascript_enabled                   = %[4]t
	  minimum_enabled_tls_protocol         = %[5]q
	  no_table_scan                        = %[6]t
	  oplog_size_mb                        = %[7]d
	  sample_size_bi_connector			 = %[8]d
	  sample_refresh_interval_bi_connector = %[9]d
	  transaction_lifetime_limit_seconds   = %[10]d
	  oplog_min_retention_hours            = %[11]d
	}]
  }
data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}

	`, orgID, projectName, name,
		*p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector, *p.TransactionLifetimeLimitSeconds,
		cast.ToInt(*p.OplogMinRetentionHours))
}

func testAccAdvancedClusterConfigAdvancedConfNoOplogHrs(orgID, projectName, name string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
	project_id             = mongodbatlas_project.cluster_project.id
	name                   = %[3]q
	cluster_type           = "REPLICASET"
  
	 replication_specs = [{
	  region_configs  = [{
		electable_specs = [{
		  instance_size = "M10"
		  node_count    = 3
		}]
		analytics_specs = [{
		  instance_size = "M10"
		  node_count    = 1
		}]
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  }]
	}]
  
	advanced_configuration  = [{
	  javascript_enabled                   = %[4]t
	  minimum_enabled_tls_protocol         = %[5]q
	  no_table_scan                        = %[6]t
	  oplog_size_mb                        = %[7]d
	  sample_size_bi_connector			 = %[8]d
	  sample_refresh_interval_bi_connector = %[9]d
	  transaction_lifetime_limit_seconds   = %[10]d
	}]
  }
data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}

	`, orgID, projectName, name,
		*p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector, *p.TransactionLifetimeLimitSeconds,
	)
}

func testAccAdvancedClusterConfigAdvancedConfDefaultWrite(orgID, projectName, name string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
	project_id             = mongodbatlas_project.cluster_project.id
	name                   = %[3]q
	cluster_type           = "REPLICASET"
  
	 replication_specs = [{
	  region_configs = [{
		electable_specs = [{
		  instance_size = "M10"
		  node_count    = 3
		}]
		analytics_specs = [{
		  instance_size = "M10"
		  node_count    = 1
		}]
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  }]
	}]
  
	advanced_configuration  = [{
	  javascript_enabled                   = %[4]t
	  minimum_enabled_tls_protocol         = %[5]q
	  no_table_scan                        = %[6]t
	  oplog_size_mb                        = %[7]d
	  sample_size_bi_connector			 = %[8]d
	  sample_refresh_interval_bi_connector = %[9]d
	  default_read_concern                 = %[10]q
	  default_write_concern                = %[11]q
	}]
  }

data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}

	`, orgID, projectName, name, *p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector, p.DefaultReadConcern, p.DefaultWriteConcern)
}

func testAccAdvancedClusterConfigReplicationSpecsAutoScaling(orgID, projectName, name string, p *matlas.AutoScaling) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}	
resource "mongodbatlas_advanced_cluster" "test" {
	project_id             = mongodbatlas_project.cluster_project.id
	name                   = %[3]q
	cluster_type           = "REPLICASET"
  
	 replication_specs = [{
	  region_configs = [{
		electable_specs = [{
		  instance_size = "M10"
		  node_count    = 3
		}]
		analytics_specs = [{
		  instance_size = "M10"
		  node_count    = 1
		}]
		  auto_scaling = [{
		   compute_enabled = %[4]t
		   disk_gb_enabled = %[5]t
		   compute_max_instance_size = %[6]q
		  }]
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  }]
	}]
  }

	`, orgID, projectName, name, *p.Compute.Enabled, *p.DiskGBEnabled, p.Compute.MaxInstanceSize)
}

func testAccAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(orgID, projectName, name string, p *matlas.AutoScaling) string {
	return fmt.Sprintf(`

resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}

resource "mongodbatlas_advanced_cluster" "test" {
	project_id             = mongodbatlas_project.cluster_project.id
	name                   = %[3]q
	cluster_type           = "REPLICASET"
  
	 replication_specs = [{
	  region_configs = [{
		electable_specs = [{
		  instance_size = "M10"
		  node_count    = 3
		}]
		analytics_specs = [{
		  instance_size = "M10"
		  node_count    = 1
		}]
		analytics_auto_scaling = [{
		  compute_enabled = %[4]t
		  disk_gb_enabled = %[5]t
		  compute_max_instance_size = %[6]q
		}]
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  }]
	}]
  }

	`, orgID, projectName, name, *p.Compute.Enabled, *p.DiskGBEnabled, p.Compute.MaxInstanceSize)
}
