package advancedcluster_test

import (
	"fmt"
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationAdvancedCluster_tenantUpgrade(t *testing.T) {
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
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config:            testAccAdvancedClusterConfigTenantBlocks(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigTenant(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigTenant(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForTenantConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigTenantUpgrade(orgID, projectName, rName),
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
		},
	})
}

func TestAccMigrationAdvancedCluster_singleAWSProviderToMultiCloud(t *testing.T) {
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
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config:            testAccAdvancedClusterConfigSingleProviderBlocks(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigSingleProvider(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigSingleProvider(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForSingleProviderConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForMultiCloudConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName)...,
				),
			},
		},
	})
}

func TestAccMigrationAdvancedCluster_multiCloud(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		rName        = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config:            testAccAdvancedClusterConfigMultiCloudBlocks(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.read_preference", "secondary"),

					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.container_id.%", "2"),
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

					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M10"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.priority", "6"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.provider_name", "GCP"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.region_name", "NORTH_AMERICA_NORTHEAST_1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationAdvancedCluster_multiCloud_stateUpgrader(t *testing.T) {
	var (
		cluster      matlas.AdvancedCluster
		resourceName = "mongodbatlas_advanced_cluster.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		rName        = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.7.0"),
				Config:            testAccAdvancedClusterConfigMultiCloudBlocksSchemaV0(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector.0.read_preference", "secondary"),

					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M10"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.node_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.instance_size", "M10"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.node_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.priority", "7"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.region_name", "US_EAST_1"),

					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M10"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.priority", "6"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.provider_name", "GCP"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.region_name", "NORTH_AMERICA_NORTHEAST_1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigMultiCloudSchemaV1(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigMultiCloud(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.read_preference", "secondary"),

					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.container_id.%", "2"),
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

					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.instance_size", "M10"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.priority", "6"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.provider_name", "GCP"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.1.region_name", "NORTH_AMERICA_NORTHEAST_1"),
				),
			},
		},
	})
}

func TestAccMigrationAdvancedCluster_multicloudSharded(t *testing.T) {
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
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config:            testAccAdvancedClusterConfigMultiCloudShardedBlocks(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "SHARDED"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigMultiCloudSharded(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigMultiCloudSharded(orgID, projectName, rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsMulticloudSharded(&cluster, resourceName, dataSourceName, dataSourceClustersName, rNameUpdated, false)...,
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigMultiCloudShardedUpdated(orgID, projectName, rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsMulticloudSharded(&cluster, resourceName, dataSourceName, dataSourceClustersName, rNameUpdated, true)...,
				),
			},
		},
	})
}

func TestAccMigrationAdvancedCluster_geoSharded(t *testing.T) {
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
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config:            testAccAdvancedClusterConfigGeoShardedBlocks(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigGeoSharded(orgID, projectName, rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigGeoSharded(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForGeoshardedConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName, false)...,
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigGeoShardedUpdated(orgID, projectName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testFuncsForGeoshardedConfig(&cluster, resourceName, dataSourceName, dataSourceClustersName, rName, true)...,
				),
			},
		},
	})
}

func TestAccMigrationAdvancedCluster_pausedToUnpaused(t *testing.T) {
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
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config:            testAccAdvancedClusterConfigSingleProviderPausedBlocks(orgID, projectName, rName, true, instanceSize),
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
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, instanceSize),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, true, instanceSize),
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
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigSingleProviderPaused(orgID, projectName, rName, false, instanceSize),
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
		},
	})
}

func TestAccMigrationAdvancedCluster_advancedConf(t *testing.T) {
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
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config:            testAccAdvancedClusterConfigAdvancedConfBlocks(orgID, projectName, rName, processArgs),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
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
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigAdvancedConf(orgID, projectName, rName, processArgs),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigAdvancedConf(orgID, projectName, rName, processArgs),
				Check: resource.ComposeAggregateTestCheckFunc(
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
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigAdvancedConfNoOplogHrs(orgID, projectName, rNameUpdated, processArgsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rNameUpdated),

					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.oplog_min_retention_hours"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.default_read_concern"),
					resource.TestCheckNoResourceAttr(resourceName, "advanced_configuration.0.default_write_concern"),
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

func TestAccMigrationAdvancedCluster_replicationSpecsAutoScaling(t *testing.T) {
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
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config:            testAccAdvancedClusterConfigReplicationSpecsAutoScalingBlocks(orgID, projectName, rName, autoScaling),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.disk_gb_enabled", "true"),

					testAccCheckAdvancedClusterScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigReplicationSpecsAutoScaling(orgID, projectName, rName, autoScaling),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigReplicationSpecsAutoScaling(orgID, projectName, rName, autoScaling),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.disk_gb_enabled", "true"),

					testAccCheckAdvancedClusterScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigReplicationSpecsAutoScaling(orgID, projectName, rNameUpdated, autoScalingUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
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

func TestAccMigrationAdvancedCluster_replicationSpecsAnalyticsAutoScaling(t *testing.T) {
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
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config:            testAccAdvancedClusterConfigReplicationSpecsAnalyticsAutoScalingBlocks(orgID, projectName, rName, autoScaling),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.disk_gb_enabled", "true"),
					testAccCheckAdvancedClusterAnalyticsScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(orgID, projectName, rName, autoScaling),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(orgID, projectName, rName, autoScaling),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAdvancedClusterExists(resourceName, &cluster),
					testAccCheckAdvancedClusterAttributes(&cluster, rName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.disk_gb_enabled", "true"),
					testAccCheckAdvancedClusterAnalyticsScaling(&cluster, *autoScaling.Compute.Enabled),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccAdvancedClusterConfigReplicationSpecsAnalyticsAutoScaling(orgID, projectName, rNameUpdated, autoScalingUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
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

func TestAccMigrationAdvancedCluster_withTags(t *testing.T) {
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
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config: testAccAdvancedClusterConfigWithTagsBlocks(orgID, projectName, rName, []matlas.Tag{
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
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: testAccAdvancedClusterConfigWithTags(orgID, projectName, rName, []matlas.Tag{
					{
						Key:   "key 1",
						Value: "value 1",
					},
					{
						Key:   "key 2",
						Value: "value 2",
					},
				}), ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
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

func TestAccMigrationAdvancedCluster_withLabels(t *testing.T) {
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
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyTeamAdvancedCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.14.0"),
				Config: testAccAdvancedClusterConfigWithLabelsBlocks(orgID, projectName, rName, []matlas.Label{
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
					resource.TestCheckResourceAttr(dataSourceName, "labels.#", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
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
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
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

func testAccAdvancedClusterConfigSingleProviderBlocks(orgID, projectName, name string) string {
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
      region_name   = "US_EAST_1"
    }
  }
}
data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

	`, orgID, projectName, name)
}

func testAccAdvancedClusterConfigTenantBlocks(orgID, projectName, name string) string {
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
      region_name           = "US_EAST_1"
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

func testAccAdvancedClusterConfigMultiCloudBlocks(orgID, projectName, name string) string {
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

  bi_connector_config {
	enabled = true
  }

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

data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, orgID, projectName, name)
}

func testAccAdvancedClusterConfigMultiCloudBlocksSchemaV0(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
	project_id   = mongodbatlas_project.cluster_project.id
	name         = %[3]q
	cluster_type = "REPLICASET"

	bi_connector {
		read_preference = "secondary"
	}
  
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

data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, orgID, projectName, name)
}

func testAccAdvancedClusterConfigMultiCloudSchemaV1(orgID, projectName, name string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "cluster_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_advanced_cluster" "test" {
	project_id   = mongodbatlas_project.cluster_project.id
	name         = %[3]q
	cluster_type = "REPLICASET"

	bi_connector_config = [{
		read_preference = "secondary"
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

func testAccAdvancedClusterConfigMultiCloudShardedBlocks(orgID, projectName, name string) string {
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

  data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, orgID, projectName, name)
}

func testAccAdvancedClusterConfigSingleProviderPausedBlocks(orgID, projectName, name string, paused bool, instanceSize string) string {
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
      region_name   = "US_EAST_1"
    }
  }
}
	`, orgID, projectName, name, paused, instanceSize)
}

func testAccAdvancedClusterConfigAdvancedConfBlocks(orgID, projectName, name string, p *matlas.ProcessArgs) string {
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
	  region_configs  {
		electable_specs {
		  instance_size = "M10"
		  node_count    = 3
		}
		analytics_specs {
		  instance_size = "M10"
		  node_count    = 1
		}
		auto_scaling {
			compute_enabled = true
			disk_gb_enabled = true
			compute_max_instance_size = "M20"
		   }
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  }
	}
  
	advanced_configuration  {
	  javascript_enabled                   = %[4]t
	  minimum_enabled_tls_protocol         = %[5]q
	  no_table_scan                        = %[6]t
	  oplog_size_mb                        = %[7]d
	  sample_size_bi_connector			 = %[8]d
	  sample_refresh_interval_bi_connector = %[9]d
	  transaction_lifetime_limit_seconds   = %[10]d
	  oplog_min_retention_hours            = %[11]d
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
		*p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector, *p.TransactionLifetimeLimitSeconds,
		cast.ToInt(*p.OplogMinRetentionHours))
}

func testAccAdvancedClusterConfigGeoShardedBlocks(orgID, projectName, name string) string {
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
  
	replication_specs { # zone n1
	  zone_name  = "zone n1"
	  num_shards = 3 # 3-shard Multi-Cloud Cluster
  
	  region_configs { # shard n1 
		electable_specs {
		  instance_size = "M10"
		  node_count    = 3
		}
		analytics_specs {
		  instance_size = "M10"
		   node_count    = 1
		}
		analytics_auto_scaling {
		  compute_enabled = true
		  compute_scale_down_enabled = false
		  compute_max_instance_size = "M20"
		}
		provider_name = "AWS"
		priority      = 7
		region_name   = "US_EAST_1"
	  }

	  region_configs { # shard n3
		electable_specs {
		  instance_size = "M10"
		  node_count    = 2
		}
		analytics_specs {
		  instance_size = "M10"
		  node_count    = 1
		}
		provider_name = "GCP"
		priority      = 6
		region_name   = "US_EAST_4"
	  }
	}


	replication_specs { # zone n2
	  zone_name  = "zone n2"
	  num_shards = 2 # 2-shard Multi-Cloud Cluster
  
	  region_configs { # shard n1 
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

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
}
	`, orgID, projectName, name)
}

func testAccAdvancedClusterConfigReplicationSpecsAutoScalingBlocks(orgID, projectName, name string, p *matlas.AutoScaling) string {
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
		region_name   = "US_EAST_1"
	  }
	}
  }

	`, orgID, projectName, name, *p.Compute.Enabled, *p.DiskGBEnabled, p.Compute.MaxInstanceSize)
}

func testAccAdvancedClusterConfigReplicationSpecsAnalyticsAutoScalingBlocks(orgID, projectName, name string, p *matlas.AutoScaling) string {
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
		region_name   = "US_EAST_1"
	  }
	}
  }

	`, orgID, projectName, name, *p.Compute.Enabled, *p.DiskGBEnabled, p.Compute.MaxInstanceSize)
}

func testAccAdvancedClusterConfigWithTagsBlocks(orgID, projectName, name string, tags []matlas.Tag) string {
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

func testAccAdvancedClusterConfigWithLabelsBlocks(orgID, projectName, name string, tags []matlas.Label) string {
	var labelsConf string
	for _, label := range tags {
		labelsConf += fmt.Sprintf(`
			labels {
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
	`, orgID, projectName, name, labelsConf)
}
