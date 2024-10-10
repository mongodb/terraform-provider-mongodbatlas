package advancedcluster_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

// last version that did not support new sharding schema or attributes
const versionBeforeISSRelease = "1.17.6"

func TestMigAdvancedCluster_replicaSetAWSProvider(t *testing.T) {
	testCase := replicaSetAWSProviderTestCase(t)
	mig.CreateAndRunTest(t, &testCase)
}

func TestMigAdvancedCluster_replicaSetMultiCloud(t *testing.T) {
	testCase := replicaSetMultiCloudTestCase(t)
	mig.CreateAndRunTest(t, &testCase)
}

func TestMigAdvancedCluster_singleShardedMultiCloud(t *testing.T) {
	testCase := singleShardedMultiCloudTestCase(t)
	mig.CreateAndRunTest(t, &testCase)
}

func TestMigAdvancedCluster_symmetricGeoShardedOldSchema(t *testing.T) {
	testCase := symmetricGeoShardedOldSchemaTestCase(t)
	mig.CreateAndRunTest(t, &testCase)
}

func TestMigAdvancedCluster_asymmetricShardedNewSchema(t *testing.T) {
	testCase := asymmetricShardedNewSchemaTestCase(t)
	mig.CreateAndRunTest(t, &testCase)
}

func TestMigAdvancedCluster_replicaSetAWSProviderUpdate(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     mig.PreCheckBasicSleep(t),
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeISSRelease),
				Config:            configReplicaSetAWSProvider(projectID, clusterName, 60, 3),
				Check:             checkReplicaSetAWSProvider(projectID, clusterName, 60, 3, false, false),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configReplicaSetAWSProvider(projectID, clusterName, 60, 5),
				Check:                    checkReplicaSetAWSProvider(projectID, clusterName, 60, 5, true, true),
			},
		},
	})
}

func TestMigAdvancedCluster_geoShardedOldSchemaUpdate(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeISSRelease),
				Config:            configGeoShardedOldSchema(orgID, projectName, clusterName, 2, 2, false),
				Check:             checkGeoShardedOldSchema(clusterName, 2, 2, false, false),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configGeoShardedOldSchema(orgID, projectName, clusterName, 2, 1, false),
				Check:                    checkGeoShardedOldSchema(clusterName, 2, 1, true, false),
			},
		},
	})
}

func TestMigAdvancedCluster_shardedMigrationFromOldToNewSchema(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeISSRelease),
				Config:            configShardedTransitionOldToNewSchema(orgID, projectName, clusterName, false),
				Check:             checkShardedTransitionOldToNewSchema(false),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configShardedTransitionOldToNewSchema(orgID, projectName, clusterName, true),
				ExpectError:              regexp.MustCompile("SERVICE_UNAVAILABLE"),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configShardedTransitionOldToNewSchema(orgID, projectName, clusterName, true),
				Check:                    checkShardedTransitionOldToNewSchema(true),
			},
		},
	})
}

func TestMigAdvancedCluster_geoShardedMigrationFromOldToNewSchema(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeISSRelease),
				Config:            configGeoShardedTransitionOldToNewSchema(orgID, projectName, clusterName, false),
				Check:             checkGeoShardedTransitionOldToNewSchema(false),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configShardedTransitionOldToNewSchema(orgID, projectName, clusterName, true),
				ExpectError:              regexp.MustCompile("SERVICE_UNAVAILABLE"),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configGeoShardedTransitionOldToNewSchema(orgID, projectName, clusterName, true),
				Check:                    checkGeoShardedTransitionOldToNewSchema(true),
			},
		},
	})
}

func TestMigAdvancedCluster_partialAdvancedConf(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.19.0") // version where change_stream_options_pre_and_post_images_expire_after_seconds was introduced
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
		// necessary to test oplog_min_retention_hours
		autoScalingConfigured = `
			auto_scaling {
				disk_gb_enabled = true
			}`
		extraArgs = `
			advanced_configuration  {
				fail_index_key_too_long              = false
				javascript_enabled                   = true
				minimum_enabled_tls_protocol         = "TLS1_1"
				no_table_scan                        = false
				oplog_min_retention_hours 		     = 4
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
		config        = configPartialAdvancedConfig(projectID, clusterName, extraArgs, autoScalingConfigured)
		configUpdated = configPartialAdvancedConfig(projectID, clusterName, extraArgsUpdated, "")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     mig.PreCheckBasicSleep(t),
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
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_min_retention_hours", "4"),
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
			mig.TestStepCheckEmptyPlan(configUpdated),
		},
	})
}

func configPartialAdvancedConfig(projectID, clusterName, extraArgs, autoScaling string) string {
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
					%[4]s
				}
			}
			%[3]s
		}
	`, projectID, clusterName, extraArgs, autoScaling)
}
