package advancedcluster_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigAdvancedCluster_replicaSetAWSProvider(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
		config      = configReplicaSetAWSProvider(projectID, clusterName, 3)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             checkReplicaSetAWSProvider(projectID, clusterName, 3, false, false),
			},
			mig.TestStepCheckEmptyPlan(config),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				Check:                    checkReplicaSetAWSProvider(projectID, clusterName, 3, true, true), // external_id will be present in latest version
			},
		},
	})
}

func TestMigAdvancedCluster_replicaSetMultiCloud(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
		config      = configReplicaSetMultiCloud(orgID, projectName, clusterName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             checkReplicaSetMultiCloud(clusterName, 3, false),
			},
			mig.TestStepCheckEmptyPlan(config),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				Check:                    checkReplicaSetMultiCloud(clusterName, 3, true), // external_id will be present in latest version
			},
		},
	})
}

func TestMigAdvancedCluster_singleShardedMultiCloud(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
		config      = configSingleShardedMultiCloud(orgID, projectName, clusterName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             checkSingleShardedMultiCloud(clusterName, false),
			},
			mig.TestStepCheckEmptyPlan(config),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				Check:                    checkSingleShardedMultiCloud(clusterName, true), // external_id will be present in latest version
			},
		},
	})
}

func TestMigAdvancedCluster_singleShardPerZoneGeoSharded(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
		config      = configGeoShardedOldSchema(orgID, projectName, clusterName, 1, 1, false)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             checkGeoShardedOldSchema(clusterName, 1, 1, false, false),
			},
			mig.TestStepCheckEmptyPlan(config),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				Check:                    checkGeoShardedOldSchema(clusterName, 1, 1, true, true), // external_id will be present in latest version
			},
		},
	})
}

func TestMigAdvancedCluster_symmetricGeoShardedOldSchema(t *testing.T) {
	acc.SkipTestForCI(t) // TODO: CLOUDP-260154 for ensuring old schema is supported
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
		config      = configGeoShardedOldSchema(orgID, projectName, clusterName, 2, 2, false)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             checkGeoShardedOldSchema(clusterName, 2, 2, false, false),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigAdvancedCluster_partialAdvancedConf(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
		extraArgs   = `
			advanced_configuration  {
				fail_index_key_too_long              = false
				javascript_enabled                   = true
				minimum_enabled_tls_protocol         = "TLS1_1"
				no_table_scan                        = false
			}

			bi_connector_config {
				enabled = true
			}`

		extraArgsUpdated = `
			advanced_configuration  {
				fail_index_key_too_long              = false
				javascript_enabled                   = true
				minimum_enabled_tls_protocol         = "TLS1_1"
				no_table_scan                        = false
				default_read_concern                 = "available"
				sample_size_bi_connector			 = 110
					sample_refresh_interval_bi_connector = 310
				}
				
				bi_connector_config {
					enabled = false
					read_preference = "secondary"
			}`
		config        = configPartialAdvancedConfig(projectID, clusterName, extraArgs)
		configUpdated = configPartialAdvancedConfig(projectID, clusterName, extraArgsUpdated)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "true"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.read_preference", "secondary"),
				),
			},
		},
	})
}

func configPartialAdvancedConfig(projectID, clusterName, extraArgs string) string {
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
					region_name   = "US_WEST_2"
				}
			}
			%[3]s
		}
	`, projectID, clusterName, extraArgs)
}
