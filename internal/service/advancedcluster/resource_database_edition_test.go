package advancedcluster_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312024/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccClusterAdvancedCluster_infiniteBasic(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseEdition(projectID, clusterName, new("INFINITE"), 2, nil),
				Check:  checkDatabaseEdition(new("INFINITE"), "INFINITE"),
				ConfigStateChecks: append(
					pluralDatabaseEditionChecks(clusterName, new("INFINITE"), "INFINITE"),
					shardSizeLimitChecks(clusterName, nil)...,
				),
			},
			{
				Config:      configDatabaseEditionWithDiskGBEnabled(projectID, clusterName, new("INFINITE"), 2),
				ExpectError: regexp.MustCompile(`autoScaling\.diskGB is not configurable for an Atlas Infinite\s+cluster`),
			},
			{
				Config:      configDatabaseEdition(projectID, clusterName, new("CORE"), 2, nil),
				ExpectError: regexp.MustCompile("databaseEdition cannot be changed"),
			},
			{
				Config:      configDatabaseEdition(projectID, clusterName, nil, 2, nil),
				ExpectError: regexp.MustCompile("databaseEdition cannot be changed"),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_infiniteShardSizeLimit(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 3)
	storageConfig := databaseEditionStorageConfig(new(1024))
	storageOnlyRecovery := acc.TestStepCheckEmptyPlan(configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, storageConfig, true))
	storageOnlyRecovery.ConfigStateChecks = shardSizeLimitChecks(clusterName, new(1024))
	computeChecks := computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
		"compute_enabled":            knownvalue.Bool(true),
		"compute_scale_down_enabled": knownvalue.Bool(false),
		"compute_max_instance_size":  knownvalue.StringExact("M20"),
	})

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseEdition(projectID, clusterName, new("INFINITE"), 2, new(1024)),
				Check:  checkDatabaseEdition(new("INFINITE"), "INFINITE"),
				ConfigStateChecks: append(
					pluralDatabaseEditionChecks(clusterName, new("INFINITE"), "INFINITE"),
					shardSizeLimitChecks(clusterName, new(1024))...,
				),
			},
			{
				Config:            configDatabaseEdition(projectID, clusterName, new("INFINITE"), 2, new(2048)),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, new(2048)),
			},
			{
				Config:            configDatabaseEdition(projectID, clusterName, new("INFINITE"), 2, new(1024)),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, new(1024)),
			},
			{
				Config:            configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, storageConfig, true),
				Check:             resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, new(1024)),
			},
			{
				Config:      configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, "disk_gb_enabled = true\n"+storageConfig, true),
				ExpectError: regexp.MustCompile(`autoScaling\.diskGB is not configurable for an Atlas Infinite\s+cluster`),
			},
			storageOnlyRecovery,
			{
				Config:      configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, "disk_gb_enabled = false\n"+storageConfig, true),
				ExpectError: regexp.MustCompile(`(?s)INVALID_ATTRIBUTE.*autoScaling\.diskGB`),
			},
			storageOnlyRecovery,
			{
				Config:      configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, databaseEditionStorageConfig(new(0)), true),
				ExpectError: regexp.MustCompile("SHARD_SIZE_LIMIT_OUT_OF_RANGE"),
			},
			storageOnlyRecovery,
			{
				Config:      configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, databaseEditionStorageConfig(new(-1)), true),
				ExpectError: regexp.MustCompile("SHARD_SIZE_LIMIT_OUT_OF_RANGE"),
			},
			storageOnlyRecovery,
			{
				Config:            configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, "", true),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, nil),
			},
			{
				Config: configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, new(1024), true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkDatabaseEdition(new("INFINITE"), "INFINITE"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(1024)), computeChecks...),
			},
			// An unrelated replication_specs change must preserve the configured shard limit.
			{
				Config: configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, "compute_enabled = false\n"+storageConfig, true),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(1024)), computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
					"compute_enabled": knownvalue.Bool(false),
				})...),
			},
			{
				Config:            configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, new(2048), true),
				Check:             checkDatabaseEdition(new("INFINITE"), "INFINITE"),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(2048)), computeChecks...),
			},
			{
				// Four electable nodes are invalid in both supported Infinite topology modes.
				Config:      configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 4, new(2048), true),
				ExpectError: regexp.MustCompile(`(?s)INVALID_ATTRIBUTE.*number of electable nodes`),
			},
			{
				Config: configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, new(2048), true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.node_count", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "replication_specs.0.region_configs.0.electable_specs.node_count", "2"),
				),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(2048)), computeChecks...),
			},
			acc.TestStepImportCluster(resourceName),
			{
				Config:            configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, nil, true),
				Check:             checkDatabaseEdition(new("INFINITE"), "INFINITE"),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil), computeChecks...),
			},
			{
				Config:            configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, new(1024), true),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, new(1024)),
			},
			{
				Config: configDatabaseEdition(projectID, clusterName, new("INFINITE"), 2, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "tags.env"),
					resource.TestCheckNoResourceAttr(dataSourceName, "tags.env"),
				),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil), computeChecks...),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_infiniteComputeAutoScaling(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)
	const computeConfig = `
		compute_enabled            = true
		compute_max_instance_size  = "M20"
	`
	const scaleDownConfig = computeConfig + `
		compute_scale_down_enabled = true
		compute_min_instance_size  = "M10"
	`
	computeOnlyConfig := configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, computeConfig, false)
	configWithStorage := configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, scaleDownConfig+databaseEditionStorageConfig(new(1024)), false)
	configWithoutStorage := configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, scaleDownConfig, false)
	computeOnlyChecks := func(maxInstanceSize string) []statecheck.StateCheck {
		return append(shardSizeLimitChecks(clusterName, nil), computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
			"compute_enabled":           knownvalue.Bool(true),
			"compute_max_instance_size": knownvalue.StringExact(maxInstanceSize),
		})...)
	}
	computeChecks := computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
		"compute_enabled":            knownvalue.Bool(true),
		"compute_scale_down_enabled": knownvalue.Bool(true),
		"compute_min_instance_size":  knownvalue.StringExact("M10"),
		"compute_max_instance_size":  knownvalue.StringExact("M20"),
	})

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			// Create without storage so the request cannot rely on storage-triggered cleanup.
			{
				Config:            computeOnlyConfig,
				Check:             checkDatabaseEdition(new("INFINITE"), "INFINITE"),
				ConfigStateChecks: computeOnlyChecks("M20"),
			},
			acc.TestStepImportCluster(resourceName),
			{
				Config:            configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, strings.ReplaceAll(computeConfig, "M20", "M30_GEN_2"), false),
				ConfigStateChecks: computeOnlyChecks("M30_GEN_2"),
			},
			{
				Config:            computeOnlyConfig,
				ConfigStateChecks: computeOnlyChecks("M20"),
			},
			acc.TestStepImportCluster(resourceName),
			// Recreate to cover compute auto-scaling both with and without storage in the create request.
			{Config: computeOnlyConfig, Destroy: true, Check: acc.CheckDestroyCluster},
			{
				Config:            configWithStorage,
				Check:             checkDatabaseEdition(new("INFINITE"), "INFINITE"),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(1024)), computeChecks...),
			},
			acc.TestStepImportCluster(resourceName),
			{
				Config: configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, new(1024), false),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(1024)), computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
					"compute_enabled":            knownvalue.Bool(true),
					"compute_scale_down_enabled": knownvalue.Bool(false),
					"compute_max_instance_size":  knownvalue.StringExact("M20"),
				})...),
			},
			{
				Config:            configWithStorage,
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(1024)), computeChecks...),
			},
			{
				Config:            configWithoutStorage,
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil), computeChecks...),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_infiniteEmptyAutoScaling(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				// Explicit null children exercise the same Terraform object as auto_scaling = {}.
				Config:            configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, "compute_enabled = null\ndisk_gb_enabled = null", false),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, nil),
			},
			{
				Config: configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, "compute_enabled = false", false),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil), computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
					"compute_enabled": knownvalue.Bool(false),
				})...),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_infiniteAnalyticsAutoScaling(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 4)
	clusterConfig := func(limit *int, maxInstanceSize, diskConfig string, omitBlocks bool) string {
		var regionConfig string
		if !omitBlocks {
			regionConfig = fmt.Sprintf(`
				electable_specs = { instance_size = "M10", node_count = 2 }
				read_only_specs = { instance_size = "M10", node_count = 1 }
				analytics_specs = { instance_size = "M10", node_count = 1 }
				analytics_auto_scaling = {
					compute_enabled = true
					compute_max_instance_size = %q
					%s
				}`, maxInstanceSize, diskConfig)
		}
		if limit != nil {
			regionConfig += "\nauto_scaling = {" + databaseEditionStorageConfig(limit) + "\n}"
		}
		return fmt.Sprintf(`
			resource "mongodbatlas_advanced_cluster" "test" {
				project_id = %q
				name = %q
				cluster_type = "REPLICASET"
				database_edition = "INFINITE"
				replication_specs = [{ zone_name = "Existing zone", region_configs = [{
					provider_name = "AWS"
					region_name = "US_EAST_1"
					priority = 7
					%s
				}] }]
			}`, projectID, clusterName, regionConfig) + dataSourcesConfig
	}
	checks := func(limit *int, maxInstanceSize string) []statecheck.StateCheck {
		region := knownvalue.ObjectPartial(map[string]knownvalue.Check{
			"electable_specs": knownvalue.ObjectPartial(map[string]knownvalue.Check{"node_count": knownvalue.Int64Exact(2)}),
			"read_only_specs": knownvalue.ObjectPartial(map[string]knownvalue.Check{"node_count": knownvalue.Int64Exact(1)}),
			"analytics_specs": knownvalue.ObjectPartial(map[string]knownvalue.Check{"node_count": knownvalue.Int64Exact(1)}),
		})
		replicationSpec := knownvalue.ObjectPartial(map[string]knownvalue.Check{
			"zone_name":      knownvalue.StringExact("Existing zone"),
			"region_configs": knownvalue.ListExact([]knownvalue.Check{region}),
		})
		path := tfjsonpath.New("replication_specs").AtSliceIndex(0)
		checks := append(shardSizeLimitChecks(clusterName, limit),
			statecheck.ExpectKnownValue(resourceName, path, replicationSpec),
			statecheck.ExpectKnownValue(dataSourceName, path, replicationSpec),
			acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), map[string]knownvalue.Check{"replication_specs.0": replicationSpec}),
		)
		return append(checks, infiniteAutoScalingChecks(clusterName, "analytics_auto_scaling", map[string]knownvalue.Check{
			"compute_enabled": knownvalue.Bool(true), "compute_max_instance_size": knownvalue.StringExact(maxInstanceSize),
		})...)
	}
	baseConfig := clusterConfig(new(1024), "M20", "", false)
	partialHardwareConfig := strings.NewReplacer(
		`, node_count = 2`, "", `, node_count = 1`, "", `zone_name = "Existing zone",`, "",
	).Replace(clusterConfig(nil, "M20", "", false))
	recovery := acc.TestStepCheckEmptyPlan(baseConfig)
	recovery.ConfigStateChecks = checks(new(1024), "M20")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{Config: baseConfig, ConfigStateChecks: checks(new(1024), "M20")},
			{Config: clusterConfig(new(1024), "M30_GEN_2", "", false), ConfigStateChecks: checks(new(1024), "M30_GEN_2")},
			// Clearing storage while omitting computed blocks must preserve active nodes and analytics scaling.
			{Config: clusterConfig(nil, "", "", true), ConfigStateChecks: checks(nil, "M30_GEN_2")},
			acc.TestStepImportCluster(resourceName),
			{Config: baseConfig, ConfigStateChecks: checks(new(1024), "M20")},
			// Clearing storage must retain planned node counts and zone metadata when omitted from configuration.
			{Config: partialHardwareConfig, ConfigStateChecks: checks(nil, "M20")},
			{Config: baseConfig, ConfigStateChecks: checks(new(1024), "M20")},
			{
				Config:      clusterConfig(new(1024), "M20", "disk_gb_enabled = true", false),
				ExpectError: regexp.MustCompile(`(?s)INVALID_ATTRIBUTE.*autoScaling\.diskGB`),
			},
			recovery,
			{
				Config:      clusterConfig(new(1024), "M20", "disk_gb_enabled = false", false),
				ExpectError: regexp.MustCompile(`(?s)INVALID_ATTRIBUTE.*autoScaling\.diskGB`),
			},
			recovery,
		},
	})
}

func TestAccClusterAdvancedCluster_infiniteShardSizeLimitDrift(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)
	clusterConfig := configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, new(1024), false)
	storagePath := tfjsonpath.New("replication_specs").AtSliceIndex(0).AtMapKey("region_configs").AtSliceIndex(0).
		AtMapKey("auto_scaling").AtMapKey("storage_config").AtMapKey("shard_size_limit_gb")
	checks := append(shardSizeLimitChecks(clusterName, new(1024)), computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
		"compute_enabled":            knownvalue.Bool(true),
		"compute_scale_down_enabled": knownvalue.Bool(false),
		"compute_max_instance_size":  knownvalue.StringExact("M20"),
	})...)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:            clusterConfig,
				ConfigStateChecks: checks,
			},
			{
				PreConfig: func() {
					ctx := t.Context()
					api := acc.ConnV2().ClustersAPI
					clusterResp, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
					require.NoError(t, err)
					specs := clusterResp.GetReplicationSpecs()
					require.Len(t, specs, 1)
					// Send only writable fields for this topology; Infinite rejects explicitly configured disk settings.
					patch := &admin.ClusterDescription20240805{ReplicationSpecs: &[]admin.ReplicationSpec20240805{{
						Id: specs[0].Id,
						RegionConfigs: &[]admin.CloudRegionConfig20240805{{
							ProviderName: new("AWS"), RegionName: new("US_EAST_1"), Priority: new(7),
							ElectableSpecs: &admin.HardwareSpec20240805{InstanceSize: new("M10"), NodeCount: new(2)},
							AutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled: new(true), ScaleDownEnabled: new(false), MaxInstanceSize: new("M20"),
								},
								StorageConfig: &admin.StorageConfig{ShardSizeLimitGB: new(2048)},
							},
						}},
					}}}
					_, _, err = api.UpdateCluster(ctx, projectID, clusterName, patch).Execute()
					require.NoError(t, err)
					diags := &diag.Diagnostics{}
					updated := advancedcluster.AwaitChanges(ctx, acc.MongoDBClient, &advancedcluster.ClusterWaitParams{
						ProjectID: projectID, ClusterName: clusterName, Timeout: 30 * time.Minute,
					}, "external shard limit update", diags)
					require.False(t, diags.HasError(), "%v", diags)
					require.Equal(t, 2048, updated.GetReplicationSpecs()[0].GetRegionConfigs()[0].AutoScaling.StorageConfig.GetShardSizeLimitGB())
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				// Refresh checks contain managed resources only; apply steps also check both data sources.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.storage_config.shard_size_limit_gb", "2048"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.compute_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.compute_scale_down_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.compute_max_instance_size", "M20"),
				),
				RefreshPlanChecks: resource.RefreshPlanChecks{PostRefresh: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectKnownValue(resourceName, storagePath, knownvalue.Int64Exact(1024)),
				}},
			},
			{
				Config:            clusterConfig,
				ConfigStateChecks: checks,
			},
		},
	})
}

func TestAccClusterAdvancedCluster_core(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:            configDatabaseEdition(projectID, clusterName, new("CORE"), 3, nil),
				Check:             checkDatabaseEdition(new("CORE"), "CORE"),
				ConfigStateChecks: pluralDatabaseEditionChecks(clusterName, new("CORE"), "CORE"),
			},
			{
				Config:      configDatabaseEdition(projectID, clusterName, new("INFINITE"), 3, nil),
				ExpectError: regexp.MustCompile("databaseEdition cannot be changed"),
			},
			{
				Config:      configDatabaseEdition(projectID, clusterName, nil, 3, nil),
				ExpectError: regexp.MustCompile("databaseEdition cannot be changed"),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func configDatabaseEdition(projectID, clusterName string, databaseEdition *string, nodeCount int, shardSizeLimitGB *int) string {
	return configDatabaseEditionWithAutoScaling(projectID, clusterName, databaseEdition, nodeCount, databaseEditionStorageConfig(shardSizeLimitGB), false)
}

func configDatabaseEditionWithComputeAutoScaling(projectID, clusterName string, databaseEdition *string, nodeCount int, shardSizeLimitGB *int, withTags bool) string {
	return configDatabaseEditionWithAutoScaling(projectID, clusterName, databaseEdition, nodeCount, fmt.Sprintf(`
			compute_enabled            = true
			compute_scale_down_enabled = false
			compute_max_instance_size  = "M20"
			%s`, databaseEditionStorageConfig(shardSizeLimitGB)), withTags)
}

func databaseEditionStorageConfig(shardSizeLimitGB *int) string {
	var storageConfig string
	if shardSizeLimitGB != nil {
		storageConfig = fmt.Sprintf(`
			storage_config = {
				shard_size_limit_gb = %d
			}`, *shardSizeLimitGB)
	}
	return storageConfig
}

func configDatabaseEditionWithDiskGBEnabled(projectID, clusterName string, databaseEdition *string, nodeCount int) string {
	return configDatabaseEditionWithAutoScaling(projectID, clusterName, databaseEdition, nodeCount, "disk_gb_enabled = true", false)
}

func configDatabaseEditionWithAutoScaling(projectID, clusterName string, databaseEdition *string, nodeCount int, autoScalingAttributes string, withTags bool) string {
	var databaseEditionConfig, autoScalingConfig, tagsConfig string
	if databaseEdition != nil {
		databaseEditionConfig = fmt.Sprintf("database_edition = %q", *databaseEdition)
	}
	if autoScalingAttributes != "" {
		autoScalingConfig = fmt.Sprintf(`
					auto_scaling = {
						%s
					}`, autoScalingAttributes)
	}
	if withTags {
		tagsConfig = `tags = { "env" = "test" }`
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id     = %[1]q
			name           = %[2]q
			cluster_type   = "REPLICASET"
			backup_enabled = true
			pit_enabled    = true
			%[3]s
			%[6]s

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = "M10"
						node_count    = %[4]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
					%[5]s
				}]
			}]
		}
	`, projectID, clusterName, databaseEditionConfig, nodeCount, autoScalingConfig, tagsConfig) + dataSourcesConfig
}

func checkDatabaseEdition(databaseEdition *string, effectiveDatabaseEdition string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		acc.CheckExistsCluster(resourceName),
		resource.TestCheckResourceAttr(dataSourceName, "effective_database_edition", effectiveDatabaseEdition),
	}
	if databaseEdition == nil {
		checks = append(checks,
			resource.TestCheckNoResourceAttr(resourceName, "database_edition"),
			resource.TestCheckNoResourceAttr(dataSourceName, "database_edition"),
		)
	} else {
		checks = append(checks,
			resource.TestCheckResourceAttr(resourceName, "database_edition", *databaseEdition),
			resource.TestCheckResourceAttr(dataSourceName, "database_edition", *databaseEdition),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func pluralDatabaseEditionChecks(clusterName string, databaseEdition *string, effectiveDatabaseEdition string) []statecheck.StateCheck {
	var databaseEditionCheck knownvalue.Check = knownvalue.Null()
	if databaseEdition != nil {
		databaseEditionCheck = knownvalue.StringExact(*databaseEdition)
	}
	return []statecheck.StateCheck{
		acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), map[string]knownvalue.Check{
			"database_edition":           databaseEditionCheck,
			"effective_database_edition": knownvalue.StringExact(effectiveDatabaseEdition),
		}),
	}
}

func computeAutoScalingChecks(clusterName string, attributes map[string]knownvalue.Check) []statecheck.StateCheck {
	return infiniteAutoScalingChecks(clusterName, "auto_scaling", attributes)
}

func infiniteAutoScalingChecks(clusterName, attribute string, attributes map[string]knownvalue.Check) []statecheck.StateCheck {
	path := tfjsonpath.New("replication_specs").AtSliceIndex(0).AtMapKey("region_configs").AtSliceIndex(0).AtMapKey(attribute)
	value := knownvalue.ObjectPartial(attributes)
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resourceName, path, value),
		statecheck.ExpectKnownValue(dataSourceName, path, value),
		acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), map[string]knownvalue.Check{
			"replication_specs.0.region_configs.0." + attribute: value,
		}),
		// Atlas hides Infinite disk settings in individual GETs; list responses can include them.
		statecheck.ExpectKnownValue(resourceName, path.AtMapKey("disk_gb_enabled"), knownvalue.Null()),
		statecheck.ExpectKnownValue(dataSourceName, path.AtMapKey("disk_gb_enabled"), knownvalue.Null()),
	}
}

func shardSizeLimitChecks(clusterName string, shardSizeLimitGB *int) []statecheck.StateCheck {
	attrName := "replication_specs.0.region_configs.0.auto_scaling"
	path := tfjsonpath.New("replication_specs").AtSliceIndex(0).AtMapKey("region_configs").AtSliceIndex(0).AtMapKey("auto_scaling")
	var shardSizeLimitCheck knownvalue.Check = storageConfigAbsentCheck{}
	if shardSizeLimitGB != nil {
		attrName += ".storage_config.shard_size_limit_gb"
		path = path.AtMapKey("storage_config").AtMapKey("shard_size_limit_gb")
		shardSizeLimitCheck = knownvalue.Int64Exact(int64(*shardSizeLimitGB))
	}
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resourceName, path, shardSizeLimitCheck),
		statecheck.ExpectKnownValue(dataSourceName, path, shardSizeLimitCheck),
		acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), map[string]knownvalue.Check{
			attrName: shardSizeLimitCheck,
		}),
	}
}

type storageConfigAbsentCheck struct{}

func (storageConfigAbsentCheck) CheckValue(value any) error {
	if value == nil {
		return nil
	}
	autoScaling, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("expected auto_scaling to be null or an object, got %T", value)
	}
	if autoScaling["storage_config"] != nil {
		return fmt.Errorf("expected storage_config to be null, got %v", autoScaling["storage_config"])
	}
	return nil
}

func (storageConfigAbsentCheck) String() string {
	return "storage_config absent"
}
