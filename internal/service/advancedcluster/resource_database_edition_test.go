package advancedcluster_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312024/admin"
	"go.mongodb.org/atlas-sdk/v20250312024/mockadmin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
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

func TestInfiniteClusterTypeValidation(t *testing.T) {
	for _, name := range []string{"ASSUME_ROLE_ARN", "TF_VAR_ASSUME_ROLE_ARN", "MONGODB_ATLAS_CLIENT_ID", "MONGODB_ATLAS_CLIENT_SECRET"} {
		t.Setenv(name, "")
	}
	for _, tc := range []struct {
		clusterType, edition string
		unsupported          bool
	}{
		{"SHARDED", "INFINITE", true},
		{"GEOSHARDED", "INFINITE", true},
		{"REPLICASET", "INFINITE", false},
		{"SHARDED", "CORE", false},
		{"GEOSHARDED", "CORE", false},
		{"SHARDED", "", false},
	} {
		t.Run(tc.clusterType+"/"+tc.edition, func(t *testing.T) {
			var expected *regexp.Regexp
			if tc.unsupported {
				expected = regexp.MustCompile(`(?s)Unsupported INFINITE cluster type.*cluster_type = "` + tc.clusterType + `".*not supported.*newer provider version`)
			}
			transport := &noTopologyRequests{}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, transport),
				Steps: []resource.TestStep{{
					Config: topologyValidationConfig(tc.clusterType, tc.edition), PlanOnly: true,
					ExpectError: expected, ExpectNonEmptyPlan: !tc.unsupported,
				}},
			})
			require.Zero(t, transport.calls.Load(), "validation and planning must not call Atlas")
		})
	}
}

func TestInfiniteClusterTypeUnknownPlanAndApply(t *testing.T) {
	ctx := t.Context()
	for _, clusterType := range []string{"SHARDED", "GEOSHARDED"} {
		for _, unknownAttribute := range []string{"cluster_type", "database_edition"} {
			for _, operation := range []string{"create", "update"} {
				t.Run(clusterType+"/"+unknownAttribute+"/"+operation, func(t *testing.T) {
					r := advancedcluster.Resource()
					var schemaResponse frameworkresource.SchemaResponse
					r.Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResponse)
					typ := schemaResponse.Schema.Type().TerraformType(ctx)
					resolved := topologyValidationAttributes(clusterType, "INFINITE")
					unknown := topologyValidationAttributes(clusterType, "INFINITE")
					unknown[unknownAttribute] = tftypes.UnknownValue
					prior := planTestValue(typ, topologyValidationAttributes("REPLICASET", "CORE"))
					if operation == "create" {
						prior = tftypes.NewValue(typ, nil)
					}
					dynamic := func(value tftypes.Value) *tfprotov6.DynamicValue {
						result, err := tfprotov6.NewDynamicValue(typ, value)
						require.NoError(t, err)
						return &result
					}
					server, err := acc.TestAccProviderV6Factories["mongodbatlas"]()
					require.NoError(t, err)
					plan := func(attributes map[string]any) *tfprotov6.PlanResourceChangeResponse {
						value := dynamic(planTestValue(typ, attributes))
						response, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
							TypeName: "mongodbatlas_advanced_cluster", PriorState: dynamic(prior), ProposedNewState: value, Config: value,
						})
						require.NoError(t, err)
						return response
					}
					initial := plan(unknown)
					require.Empty(t, initial.Diagnostics)
					initialValue, err := initial.PlannedState.Unmarshal(typ)
					require.NoError(t, err)
					attribute, _, err := tftypes.WalkAttributePath(initialValue, tftypes.NewAttributePath().WithAttributeName(unknownAttribute))
					require.NoError(t, err)
					require.False(t, attribute.(tftypes.Value).IsKnown())
					final := plan(resolved)
					require.Len(t, final.Diagnostics, 1)
					require.Equal(t, "Unsupported INFINITE cluster type", final.Diagnostics[0].Summary)
					require.Contains(t, final.Diagnostics[0].Detail, clusterType)

					api := mockadmin.NewClustersAPI(t)
					r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
					resolvedPlan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: planTestValue(typ, resolved)}
					if operation == "create" {
						var response frameworkresource.CreateResponse
						r.Create(ctx, frameworkresource.CreateRequest{Plan: resolvedPlan}, &response)
						require.Len(t, response.Diagnostics.Errors(), 1)
						require.Equal(t, "Unsupported INFINITE cluster type", response.Diagnostics.Errors()[0].Summary())
					} else {
						var response frameworkresource.UpdateResponse
						r.Update(ctx, frameworkresource.UpdateRequest{
							Plan: resolvedPlan, State: tfsdk.State{Schema: schemaResponse.Schema, Raw: prior},
						}, &response)
						require.Len(t, response.Diagnostics.Errors(), 1)
						require.Equal(t, "Unsupported INFINITE cluster type", response.Diagnostics.Errors()[0].Summary())
					}
					require.Empty(t, api.Calls, "unsupported topologies must fail before any Atlas request")
				})
			}
		}
	}
}

func TestUpdateVerifiesImportedDatabaseEditionBeforeWrites(t *testing.T) {
	for _, clusterType := range []string{"SHARDED", "GEOSHARDED"} {
		for _, tc := range []struct {
			requestedEdition       any
			lookupError            error
			name, effectiveEdition string
		}{
			{name: "omitted edition", effectiveEdition: "INFINITE"},
			{name: "explicit CORE cannot bypass effective edition", requestedEdition: "CORE", effectiveEdition: "INFINITE"},
			{name: "CORE permits update", effectiveEdition: "CORE"},
			{name: "failed lookup prevents update", lookupError: errors.New("edition lookup failed")},
		} {
			t.Run(clusterType+"/"+tc.name, func(t *testing.T) {
				ctx := t.Context()
				r := advancedcluster.Resource()
				var schemaResponse frameworkresource.SchemaResponse
				r.Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResponse)
				typ := schemaResponse.Schema.Type().TerraformType(ctx)
				prior := planTestValue(typ, topologyValidationAttributes("REPLICASET", nil))
				attributes := topologyValidationAttributes(clusterType, tc.requestedEdition)
				if tc.effectiveEdition == "INFINITE" {
					attributes["pinned_fcv"] = map[string]any{"expiration_date": "2099-01-01T00:00:00Z", "version": "8.0"}
				}
				plan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: planTestValue(typ, attributes)}
				api := mockadmin.NewClustersAPI(t)
				r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
				api.EXPECT().GetCluster(mock.Anything, dummyProjectID, "example").Return(admin.GetClusterApiRequest{ApiService: api}).Once()
				api.EXPECT().GetClusterExecute(mock.Anything).Return(&admin.ClusterDescription20240805{
					ClusterType: new("REPLICASET"), EffectiveDatabaseEdition: new(tc.effectiveEdition),
				}, nil, tc.lookupError).Once()
				apiError := errors.New("update request inspected")
				if tc.effectiveEdition == "CORE" {
					api.EXPECT().UpdateCluster(mock.Anything, dummyProjectID, "example", mock.Anything).Return(admin.UpdateClusterApiRequest{ApiService: api}).Once()
					api.EXPECT().UpdateClusterExecute(mock.Anything).Return(nil, nil, apiError).Once()
				}
				var response frameworkresource.UpdateResponse
				r.Update(ctx, frameworkresource.UpdateRequest{
					Plan: plan, State: tfsdk.State{Schema: schemaResponse.Schema, Raw: prior},
				}, &response)
				require.Len(t, response.Diagnostics.Errors(), 1)
				switch {
				case tc.lookupError != nil:
					require.Contains(t, response.Diagnostics.Errors()[0].Detail(), tc.lookupError.Error())
				case tc.effectiveEdition == "CORE":
					require.Contains(t, response.Diagnostics.Errors()[0].Detail(), apiError.Error())
				default:
					require.Equal(t, "Unsupported INFINITE cluster type", response.Diagnostics.Errors()[0].Summary())
					require.Contains(t, response.Diagnostics.Errors()[0].Detail(), clusterType)
				}
				if tc.effectiveEdition != "CORE" {
					api.AssertNotCalled(t, "UpdateCluster", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				}
				api.AssertNotCalled(t, "PinFeatureCompatibilityVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			})
		}
	}
}

func TestCreateSkipsInfiniteValidationWhenEditionOmitted(t *testing.T) {
	for _, clusterType := range []string{"SHARDED", "GEOSHARDED"} {
		t.Run(clusterType, func(t *testing.T) {
			ctx := t.Context()
			r := advancedcluster.Resource()
			var schemaResponse frameworkresource.SchemaResponse
			r.Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResponse)
			typ := schemaResponse.Schema.Type().TerraformType(ctx)
			plan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: planTestValue(typ, topologyValidationAttributes(clusterType, nil))}
			api := mockadmin.NewClustersAPI(t)
			r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
			apiError := errors.New("create request inspected")
			api.EXPECT().CreateCluster(mock.Anything, dummyProjectID, mock.Anything).Return(admin.CreateClusterApiRequest{ApiService: api}).Once()
			api.EXPECT().CreateClusterExecute(mock.Anything).Return(nil, nil, apiError).Once()
			var response frameworkresource.CreateResponse
			r.Create(ctx, frameworkresource.CreateRequest{Plan: plan}, &response)
			require.Len(t, response.Diagnostics.Errors(), 1)
			require.Contains(t, response.Diagnostics.Errors()[0].Detail(), apiError.Error(),
				"an omitted database_edition defaults to CORE server-side and must not be treated as INFINITE")
		})
	}
}

func topologyValidationAttributes(clusterType, edition any) map[string]any {
	return map[string]any{
		"name": "example", "project_id": dummyProjectID, "cluster_type": clusterType, "database_edition": edition,
		"replication_specs": []any{map[string]any{"region_configs": []any{map[string]any{
			"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7),
			"electable_specs": map[string]any{"instance_size": "M10", "node_count": int64(3)},
		}}}},
	}
}

func topologyValidationConfig(clusterType, edition string) string {
	editionAttribute := ""
	if edition != "" {
		editionAttribute = fmt.Sprintf("database_edition = %q", edition)
	}
	return fmt.Sprintf(`
provider "mongodbatlas" {
  public_key = "test-public-key"
  private_key = "test-private-key"
  base_url = "https://atlas.invalid/"
}
resource "mongodbatlas_advanced_cluster" "test" {
  name = "example"
  project_id = %q
  cluster_type = %q
  %s
  replication_specs = [{
    region_configs = [{
      provider_name = "AWS"
      region_name = "US_EAST_1"
      priority = 7
      electable_specs = { instance_size = "M10", node_count = 3 }
    }]
  }]
}
`, dummyProjectID, clusterType, editionAttribute)
}

type noTopologyRequests struct {
	calls atomic.Int64
}

func (m *noTopologyRequests) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *noTopologyRequests) ResetHTTPClient(*http.Client) {}

func (m *noTopologyRequests) RoundTrip(req *http.Request) (*http.Response, error) {
	m.calls.Add(1)
	return nil, fmt.Errorf("unexpected Atlas request: %s %s", req.Method, req.URL.Path)
}

func TestCreateRejectsUnexpectedInfiniteTopologyPreservesState(t *testing.T) {
	t.Setenv("ASSUME_ROLE_ARN", "")
	t.Setenv("TF_VAR_ASSUME_ROLE_ARN", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_ID", "")
	t.Setenv("MONGODB_ATLAS_CLIENT_SECRET", "")
	require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach())
	for _, clusterType := range []string{"SHARDED", "GEOSHARDED"} {
		t.Run(clusterType, func(t *testing.T) {
			ctx := t.Context()
			httpMock := &unexpectedInfiniteCreateHTTPMock{clusterType: clusterType}
			server, err := unit.TestAccProviderV6FactoriesWithMock(t, httpMock)["mongodbatlas"]()
			require.NoError(t, err)
			schemas, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
			require.NoError(t, err)
			require.Empty(t, schemas.Diagnostics)
			dynamic := func(typ tftypes.Type, input any) *tfprotov6.DynamicValue {
				value, err := tfprotov6.NewDynamicValue(typ, planTestValue(typ, input))
				require.NoError(t, err)
				return &value
			}
			providerConfig := dynamic(schemas.Provider.ValueType(), map[string]any{
				"public_key": "test-public-key", "private_key": "test-private-key", "base_url": "https://atlas.invalid/",
			})
			configured, err := server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{Config: providerConfig})
			require.NoError(t, err)
			require.Empty(t, configured.Diagnostics)
			typ := schemas.ResourceSchemas["mongodbatlas_advanced_cluster"].ValueType()
			clusterConfig := dynamic(typ, map[string]any{
				"name": "example", "project_id": "111111111111111111111111", "cluster_type": clusterType,
				"paused": true,
				"replication_specs": []any{map[string]any{"region_configs": []any{map[string]any{
					"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7),
					"electable_specs": map[string]any{"instance_size": "M10", "node_count": int64(3)},
				}}}},
				"advanced_configuration": map[string]any{"minimum_enabled_tls_protocol": "TLS1_2"},
			})
			empty := dynamic(typ, nil)
			plan, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
				TypeName: "mongodbatlas_advanced_cluster", PriorState: empty, ProposedNewState: clusterConfig, Config: clusterConfig,
			})
			require.NoError(t, err)
			require.Empty(t, plan.Diagnostics)
			created, err := server.ApplyResourceChange(ctx, &tfprotov6.ApplyResourceChangeRequest{
				TypeName: "mongodbatlas_advanced_cluster", PriorState: empty, PlannedState: plan.PlannedState, Config: clusterConfig,
			})
			require.NoError(t, err)
			require.Len(t, created.Diagnostics, 1, "the unsupported-topology diagnostic must not be accompanied by invalid state errors")
			require.Equal(t, "Unsupported INFINITE cluster type", created.Diagnostics[0].Summary)
			require.Contains(t, created.Diagnostics[0].Detail, clusterType)
			require.Equal(t, []string{"POST", "GET"}, httpMock.methods, "no pause or advanced configuration writes may follow the unexpected edition")
			state, err := created.NewState.Unmarshal(typ)
			require.NoError(t, err)
			require.False(t, state.IsNull(), "a successfully created cluster must remain tracked even when its topology is rejected")
			require.True(t, state.IsFullyKnown(), "retained state must not include unknown plan values")
			for name, expected := range map[string]string{
				"name": "example", "project_id": "111111111111111111111111", "cluster_id": "333333333333333333333333",
			} {
				value, _, err := tftypes.WalkAttributePath(state, tftypes.NewAttributePath().WithAttributeName(name))
				require.NoError(t, err)
				require.Equal(t, tftypes.NewValue(tftypes.String, expected), value)
			}
			// Destroy without refreshing first remains available for a cluster rejected during creation.
			destroyed, err := server.ApplyResourceChange(ctx, &tfprotov6.ApplyResourceChangeRequest{
				TypeName: "mongodbatlas_advanced_cluster", PriorState: created.NewState, PlannedState: empty, Config: empty,
			})
			require.NoError(t, err)
			require.Empty(t, destroyed.Diagnostics)
			destroyedState, err := destroyed.NewState.Unmarshal(typ)
			require.NoError(t, err)
			require.True(t, destroyedState.IsNull())
			require.Equal(t, []string{"POST", "GET", "DELETE", "GET"}, httpMock.methods)
		})
	}
}

type unexpectedInfiniteCreateHTTPMock struct {
	clusterType string
	methods     []string
	deleted     bool
}

func (m *unexpectedInfiniteCreateHTTPMock) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *unexpectedInfiniteCreateHTTPMock) ResetHTTPClient(*http.Client) {}

func (m *unexpectedInfiniteCreateHTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	const clustersPath = "/api/atlas/v2/groups/111111111111111111111111/clusters"
	const clusterPath = clustersPath + "/example"
	m.methods = append(m.methods, req.Method)
	status := http.StatusOK
	var cluster map[string]any
	if err := json.Unmarshal(fmt.Appendf(nil, storageConfigClusterResponse, 1024), &cluster); err != nil {
		return nil, err
	}
	delete(cluster, "databaseEdition")
	cluster["clusterType"], cluster["effectiveDatabaseEdition"] = m.clusterType, "INFINITE"
	body, err := json.Marshal(cluster)
	if err != nil {
		return nil, err
	}
	switch {
	case req.Method == http.MethodPost && req.URL.Path == clustersPath:
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, err
		}
		if payload["databaseEdition"] != nil || payload["paused"] != nil {
			return nil, fmt.Errorf("create must omit databaseEdition and paused")
		}
	case req.Method == http.MethodGet && req.URL.Path == clusterPath:
		if m.deleted {
			status, body = http.StatusNotFound, []byte(`{"errorCode":"CLUSTER_NOT_FOUND","error":404}`)
		}
	case req.Method == http.MethodDelete && req.URL.Path == clusterPath:
		m.deleted = true
		body = []byte(`{}`)
	default:
		return nil, fmt.Errorf("unexpected mocked Atlas request: %s %s", req.Method, req.URL.Path)
	}
	return &http.Response{
		StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}},
		Body: io.NopCloser(strings.NewReader(string(body))), Request: req,
	}, nil
}
