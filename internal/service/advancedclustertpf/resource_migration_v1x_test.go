package advancedclustertpf_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

var versionBeforeTPFGARelease = os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")

func TestV1xMigClusterAdvancedClusterConfig_geoShardedNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	isSDKv2 := acc.IsTestSDKv2ToTPF()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); mig.PreCheckLast1XVersion(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configGeoShardedTransitionOldToNewSchema(t, !isSDKv2, projectID, clusterName, true),
				Check:             checkGeoShardedTransitionOldToNewSchema(!isSDKv2, true),
			},
			mig.TestStepCheckEmptyPlan(configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true)),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true),
				Check:                    checkGeoShardedTransitionOldToNewSchema(true, true),
			},
			mig.TestStepCheckEmptyPlan(configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true)),
		},
	})
}

func configGeoShardedTransitionOldToNewSchema(t *testing.T, isTPF bool, projectID, name string, useNewSchema bool) string {
	t.Helper()
	var numShardsStr string
	var diskSizeGB string
	if !useNewSchema {
		numShardsStr = `num_shards = 2`
		diskSizeGB = `disk_size_gb = 15`
	}
	replicationSpec := `
		replication_specs {
			%[1]s
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
				region_name   = %[2]q
			}
			zone_name = %[3]q
		}
	`

	var replicationSpecs string
	if !useNewSchema {
		replicationSpecs = fmt.Sprintf(`
			%[1]s
			%[2]s
		`, fmt.Sprintf(replicationSpec, numShardsStr, "US_EAST_1", "zone 1"), fmt.Sprintf(replicationSpec, numShardsStr, "EU_WEST_1", "zone 2"))
	} else {
		replicationSpecs = fmt.Sprintf(`
			%[1]s
			%[2]s
			%[3]s
			%[4]s
		`, fmt.Sprintf(replicationSpec, numShardsStr, "US_EAST_1", "zone 1"), fmt.Sprintf(replicationSpec, numShardsStr, "US_EAST_1", "zone 1"),
			fmt.Sprintf(replicationSpec, numShardsStr, "EU_WEST_1", "zone 2"), fmt.Sprintf(replicationSpec, numShardsStr, "EU_WEST_1", "zone 2"))
	}

	var dataSources = dataSourcesTFOldSchema
	if useNewSchema {
		dataSources = dataSourcesTFNewSchema
	}

	return acc.ConvertAdvancedClusterToTPF(t, isTPF, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = %[1]q
			name = %[2]q
			backup_enabled = false
			cluster_type   = "GEOSHARDED"

			%[4]s

			%[3]s
		}
	`, projectID, name, replicationSpecs, diskSizeGB)) + dataSources
}

func checkGeoShardedTransitionOldToNewSchema(isTPF, useNewSchema bool) resource.TestCheckFunc {
	if useNewSchema {
		return checkAggrMig(isTPF, false,
			[]string{
				"replication_specs.0.external_id", "replication_specs.1.external_id", "replication_specs.2.external_id", "replication_specs.3.external_id",
			},
			map[string]string{
				"replication_specs.#":           "4",
				"replication_specs.0.zone_name": "zone 1",
				"replication_specs.1.zone_name": "zone 1",
				"replication_specs.2.zone_name": "zone 2",
				"replication_specs.3.zone_name": "zone 2",
			},
		)
	}
	return checkAggrMig(isTPF, false,
		[]string{},
		map[string]string{
			"replication_specs.#":           "2",
			"replication_specs.0.zone_name": "zone 1",
			"replication_specs.1.zone_name": "zone 2",
		},
	)
}

// SDKv2/TPF pre OLD to new - sharded
func TestV1xMigAdvancedCluster_oldToNewSchemaWithAutoscalingEnabled(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	isSDKv2 := acc.IsTestSDKv2ToTPF()
	// if not SDKv2, ensure preview flag is set - add to precheck
	// isTPFPreview := os.Getenv("MONGODB_ATLAS_TEST_PREVIEW_TO_TPF")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicSleep(t, nil, projectID, clusterName); mig.PreCheckLast1XVersion(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configShardedTransitionOldToNewSchema(t, !isSDKv2, projectID, clusterName, false, true, false),
				Check:             acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER"),
			},
			mig.TestStepCheckEmptyPlan(configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, true, false)),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, true, false),
				Check:                    acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER"),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, true, true),
				Check:                    acc.CheckIndependentShardScalingMode(resourceName, clusterName, "SHARD"),
			},
			mig.TestStepCheckEmptyPlan(configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, true, true)),
		},
	})
}

func TestV1xMigAdvancedCluster_shardedNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	versionBeforeTPFGARelease := os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")
	isSDKv2 := acc.IsTestSDKv2ToTPF()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); mig.PreCheckLast1XVersion(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configShardedTransitionOldToNewSchema(t, !isSDKv2, projectID, clusterName, true, false, false),
				Check:             checkShardedTransitionOldToNewSchema(!isSDKv2, true),
			},
			mig.TestStepCheckEmptyPlan(configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, false, false)),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, false, false),
				Check:                    checkShardedTransitionOldToNewSchema(true, true),
			},
			mig.TestStepCheckEmptyPlan(configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, false, false)),
		},
	})
}

func configShardedTransitionOldToNewSchema(t *testing.T, isTPF bool, projectID, name string, useNewSchema, autoscaling, isUpdate bool) string {
	t.Helper()
	var numShardsStr string
	var diskSizeGBStr string
	if !useNewSchema {
		numShardsStr = `num_shards = 2`
		diskSizeGBStr = `disk_size_gb = 15`
	}
	var autoscalingStr string
	if autoscaling {
		autoscalingStr = `auto_scaling {
			compute_enabled = true
			disk_gb_enabled = true
			compute_max_instance_size = "M20"
		}`

		if isUpdate {
			autoscalingStr = `auto_scaling {
				compute_enabled = true
				disk_gb_enabled = true
				compute_max_instance_size = "M30"
			}`
		}
	}

	replicationSpec := fmt.Sprintf(`
		replication_specs {
			%[1]s

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
				%[2]s
			}
		}
	`, numShardsStr, autoscalingStr)

	var replicationSpecs string
	if useNewSchema {
		replicationSpecs = fmt.Sprintf(`
			%[1]s
			%[1]s
		`, replicationSpec)
	} else {
		replicationSpecs = replicationSpec
	}

	var dataSources = dataSourcesTFOldSchema
	if useNewSchema {
		dataSources = dataSourcesTFNewSchema
	}

	return acc.ConvertAdvancedClusterToTPF(t, isTPF, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = %[1]q
			name = %[2]q
			backup_enabled = false
			cluster_type   = "SHARDED"

			%[4]s

			%[3]s
		}

	`, projectID, name, replicationSpecs, diskSizeGBStr)) + dataSources
}

func checkShardedTransitionOldToNewSchema(isTPF, useNewSchema bool) resource.TestCheckFunc {
	var amtOfReplicationSpecs int
	if useNewSchema {
		amtOfReplicationSpecs = 2
	} else {
		amtOfReplicationSpecs = 1
	}
	var checksForNewSchema []resource.TestCheckFunc
	if useNewSchema {
		checksForNewSchema = []resource.TestCheckFunc{
			checkAggrMig(isTPF, false, []string{"replication_specs.0.external_id", "replication_specs.1.external_id"},
				map[string]string{
					"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
					"replication_specs.1.region_configs.0.electable_specs.0.instance_size": "M10",
					"replication_specs.1.region_configs.0.analytics_specs.0.instance_size": "M10",
				}),
		}
	}

	return checkAggrMig(isTPF, false,
		[]string{},
		map[string]string{
			"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
			"replication_specs.0.region_configs.0.electable_specs.0.instance_size": "M10",
			"replication_specs.0.region_configs.0.analytics_specs.0.instance_size": "M10",
		},
		checksForNewSchema...,
	)
}

func TestV1xMigAdvancedCluster_geoShardedMigrationFromOldToNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	versionBeforeTPFGARelease := os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")
	isSDKv2 := acc.IsTestSDKv2ToTPF()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); mig.PreCheckLast1XVersion(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configGeoShardedTransitionOldToNewSchema(t, !isSDKv2, projectID, clusterName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkGeoShardedTransitionOldToNewSchema(false, false),
					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
			},
			mig.TestStepCheckEmptyPlan(configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true)),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkGeoShardedTransitionOldToNewSchema(true, true),
					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
			},
			mig.TestStepCheckEmptyPlan(configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true)),
		},
	})
}

func TestV1xMigAdvancedCluster_replicaSetAWSProvider(t *testing.T) {
	var (
		projectID, clusterName    = acc.ProjectIDExecutionWithCluster(t, 6)
		versionBeforeTPFGARelease = os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")
		isSDKv2                   = acc.IsTestSDKv2ToTPF()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config: configAWSProvider(t, ReplicaSetAWSConfig{
					ProjectID:          projectID,
					ClusterName:        clusterName,
					ClusterType:        "REPLICASET",
					DiskSizeGB:         60,
					NodeCountElectable: 3,
					WithAnalyticsSpecs: true,
				}, !isSDKv2),
				Check: checkReplicaSetAWSProvider(!isSDKv2, false, projectID, clusterName, 60, 3, true, true),
			},
			mig.TestStepCheckEmptyPlan(configAWSProvider(t, ReplicaSetAWSConfig{
				ProjectID:          projectID,
				ClusterName:        clusterName,
				ClusterType:        "REPLICASET",
				DiskSizeGB:         60,
				NodeCountElectable: 3,
				WithAnalyticsSpecs: true,
			}, true)),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: configAWSProvider(t, ReplicaSetAWSConfig{
					ProjectID:          projectID,
					ClusterName:        clusterName,
					ClusterType:        "REPLICASET",
					DiskSizeGB:         60,
					NodeCountElectable: 3,
					WithAnalyticsSpecs: true,
				}, true),
				Check: checkReplicaSetAWSProvider(true, false, projectID, clusterName, 60, 3, true, true),
			},
		},
	})
}

func TestV1xMigAdvancedCluster_replicaSetMultiCloud(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
		isSDKv2     = acc.IsTestSDKv2ToTPF()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configReplicaSetMultiCloud(t, orgID, projectName, clusterName, !isSDKv2),
				Check:             checkReplicaSetMultiCloud(!isSDKv2, false, clusterName, 3),
			},
			mig.TestStepCheckEmptyPlan(configReplicaSetMultiCloud(t, orgID, projectName, clusterName, true)),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configReplicaSetMultiCloud(t, orgID, projectName, clusterName, true),
				Check:                    checkReplicaSetMultiCloud(true, false, clusterName, 3),
			},
		},
	})
}
