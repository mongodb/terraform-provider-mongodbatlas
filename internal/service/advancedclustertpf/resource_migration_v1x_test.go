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

// TODO: this may fail because 2nd step might be using num_shards
// geo-sharded ISS
// this requires dataSources = dataSourcesTFNewSchema ONLY as all steps use new schema
func TestMigV1xClusterAdvancedClusterConfig_geoShardedTransitionFromOldToNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	versionBeforeTPFGARelease := os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")
	isSDKv2 := true

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); mig.PreCheckLast1XVersion(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configGeoShardedTransitionOldToNewSchema(t, !isSDKv2, projectID, clusterName, true),
				Check:             checkGeoShardedTransitionOldToNewSchema(!isSDKv2, true),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true),
				Check:                    checkGeoShardedTransitionOldToNewSchema(true, true),
			},
			mig.TestStepCheckEmptyPlan(resourceName),
			// acc.TestStepImportCluster(resourceName),
		},
	})
}

func configGeoShardedTransitionOldToNewSchema(t *testing.T, isTPF bool, projectID, name string, useNewSchema bool) string {
	t.Helper()
	var numShardsStr string
	if !useNewSchema {
		numShardsStr = `num_shards = 2`
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

			%[3]s
		}
	`, projectID, name, replicationSpecs)) + dataSources
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

// sharded -pre-ISS
func TestMigV1xAdvancedCluster_oldToNewSchemaWithAutoscalingEnabled(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	versionBeforeTPFGARelease := os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")
	isSDKv2 := true

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicSleep(t, nil, projectID, clusterName); mig.PreCheckLast1XVersion(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configShardedTransitionOldToNewSchema(t, !isSDKv2, projectID, clusterName, false, true, false),
				Check:             acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER"),
			},
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
			mig.TestStepCheckEmptyPlan(resourceName),
			// acc.TestStepImportCluster(resourceName),
		},
	})
}

// TODO: this test may be redundant with TestMigAdvancedCluster_shardedMigrationFromOldToNewSchema
// func TestAccAdvancedCluster_oldToNewSchemaWithAutoscalingDisabledToEnabled(t *testing.T) {
// 	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
// 	versionBeforeTPFGARelease := os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { acc.PreCheckBasicSleep(t, nil, projectID, clusterName); mig.PreCheckLast1XVersion(t) },
// 		CheckDestroy: acc.CheckDestroyCluster,
// 		Steps: []resource.TestStep{
// 			{
// 				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
// 				Config:            configShardedTransitionOldToNewSchema(t, false, projectID, clusterName, false, false),
// 				Check:             acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER"),
// 			},
// 			{
// 				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
// 				Config:                   configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, false),
// 				Check:                    acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER"),
// 			},
// 			{
// 				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
// 				Config:                   configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, true),
// 				Check:                    acc.CheckIndependentShardScalingMode(resourceName, clusterName, "SHARD"),
// 			},
// 			// acc.TestStepImportCluster(resourceName),
// 		},
// 	})
// }

// TODO: this test may be redundant with TestMigAdvancedCluster_shardedMigrationFromOldToNewSchema & the test above
// func TestAccClusterAdvancedClusterConfig_shardedTransitionFromOldToNewSchema(t *testing.T) {
// 	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:                 func() { acc.PreCheckBasic(t) },
// 		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
// 		CheckDestroy:             acc.CheckDestroyCluster,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, false, false),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					checkShardedTransitionOldToNewSchema(true, false),
// 					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
// 			},
// 			{
// 				Config: configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, false),
// 				Check:  checkShardedTransitionOldToNewSchema(true, true),
// 			},
// 			acc.TestStepImportCluster(resourceName),
// 		},
// 	})
// }

// sharded -ISS
func TestMigV1xAdvancedCluster_shardedMigrationFromOldToNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	versionBeforeTPFGARelease := os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")
	isSDKv2 := true

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); mig.PreCheckLast1XVersion(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configShardedTransitionOldToNewSchema(t, !isSDKv2, projectID, clusterName, true, false, false),
				Check:             checkShardedTransitionOldToNewSchema(!isSDKv2, true),
			},
			// {
			// 	ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
			// 	Config:            configShardedTransitionOldToNewSchema(t, false, projectID, clusterName, true, false),
			// 	Check:             checkShardedTransitionOldToNewSchema(false, true),
			// },
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, false, false),
				Check:                    checkShardedTransitionOldToNewSchema(true, true),
			},
			mig.TestStepCheckEmptyPlan(resourceName),
		},
	})
}

func configShardedTransitionOldToNewSchema(t *testing.T, isTPF bool, projectID, name string, useNewSchema, autoscaling, isUpdate bool) string {
	t.Helper()
	var numShardsStr string
	if !useNewSchema {
		numShardsStr = `num_shards = 2`
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
				compute_min_instance_size = "M10"
				compute_enabled = true
				disk_gb_enabled = true
				compute_max_instance_size = "M20"
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

			%[3]s
		}

	`, projectID, name, replicationSpecs)) + dataSources
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
			checkAggrMig(isTPF, false, []string{"replication_specs.1.id", "replication_specs.0.external_id", "replication_specs.1.external_id"},
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

// geo-sharded pre-ISS
// func TestMigAdvancedCluster_symmetricGeoShardedOldSchema(t *testing.T) {
// 	// migTest(t, symmetricGeoShardedOldSchemaTestCase)
// 	resource.ParallelTest(t, symmetricGeoShardedOldSchemaTestCase(t, true))
// }

// func symmetricGeoShardedOldSchemaTestCase(t *testing.T, useSDKv2 ...bool) resource.TestCase {
// 	t.Helper()

// 	var (
// 		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 18)
// 		isSDKv2                = isOptionalTrue(useSDKv2...)
// 		// isTPF                  = !isSDKv2
// 		versionBeforeTPFGARelease = os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")
// 	)

// 	return resource.TestCase{
// 		PreCheck:                 func() { acc.PreCheckBasic(t); mig.PreCheckLast1XVersion(t) },
// 		CheckDestroy:             acc.CheckDestroyCluster,
// 		Steps: []resource.TestStep{
// 			{
// 				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
// 				Config: configGeoShardedOldSchema(t, projectID, clusterName, 2, 2, false, isSDKv2),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					checkGeoShardedOldSchema(!isSDKv2, clusterName, 2, 2, true, false),
// 					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
// 			},
// 			{
// 				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
// 				Config: configGeoShardedOldSchema(t, projectID, clusterName, 2, 2, false, false),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					checkGeoShardedOldSchema(true, clusterName, 2, 2, true, false),
// 					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
// 			},
// 			{
// 				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
// 				Config: configGeoShardedOldSchema(t, projectID, clusterName, 3, 3, false, false),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					checkGeoShardedOldSchema(true, clusterName, 3, 3, true, true),
// 					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
// 			},
// 			// acc.TestStepImportCluster(resourceName, "replication_specs"), // Import with old schema will NOT use `num_shards`
// 			// acc.TestStepImportCluster(resourceName), // Import with old schema will NOT use `num_shards`

// 		},
// 	}
// }

// geo- pre-ISS
func TestMigV1xAdvancedCluster_geoShardedMigrationFromOldToNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	versionBeforeTPFGARelease := os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")
	isSDKv2 := true

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
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkGeoShardedTransitionOldToNewSchema(true, true),
					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
			},
			mig.TestStepCheckEmptyPlan(resourceName),
		},
	})
}

func TestMigV1xAdvancedCluster_replicaSetAWSProvider(t *testing.T) {
	var (
		projectID, clusterName    = acc.ProjectIDExecutionWithCluster(t, 6)
		versionBeforeTPFGARelease = os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")
		isSDKv2                   = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
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
				}, isSDKv2),
				Check: checkReplicaSetAWSProvider(!isSDKv2, false, projectID, clusterName, 60, 3, true, true),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: configAWSProvider(t, ReplicaSetAWSConfig{
					ProjectID:          projectID,
					ClusterName:        clusterName,
					ClusterType:        "REPLICASET",
					DiskSizeGB:         60,
					NodeCountElectable: 3,
					WithAnalyticsSpecs: true,
				}, false),
				Check: checkReplicaSetAWSProvider(true, false, projectID, clusterName, 60, 3, true, true),
			},
			mig.TestStepCheckEmptyPlan(resourceName),
		},
	})

}

func TestMigV1xAdvancedCluster_replicaSetMultiCloud(t *testing.T) {
	var (
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()

		isSDKv2 = true
		// isTPF              = !isSDKv2
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acc.PreCheckBasic(t) },

		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configReplicaSetMultiCloud(t, orgID, projectName, clusterName, isSDKv2),
				Check:             checkReplicaSetMultiCloud(!isSDKv2, false, clusterName, 3),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configReplicaSetMultiCloud(t, orgID, projectName, clusterNameUpdated, false),
				Check:                    checkReplicaSetMultiCloud(true, false, clusterNameUpdated, 3),
			},
			mig.TestStepCheckEmptyPlan(resourceName),
		},
	})
}
