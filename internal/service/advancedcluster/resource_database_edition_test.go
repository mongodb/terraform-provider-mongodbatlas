package advancedcluster_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

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
			// Clear storage while compute stays configured, on a freshly created cluster, to verify
			// the provider correctly handles storage removal with the explicit-null approach.
			{
				Config: configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, new(1024), false),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(1024)), computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
					"compute_enabled":            knownvalue.Bool(true),
					"compute_scale_down_enabled": knownvalue.Bool(false),
					"compute_max_instance_size":  knownvalue.StringExact("M20"),
				})...),
			},
			{
				Config: configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, nil, false),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil), computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
					"compute_enabled":            knownvalue.Bool(true),
					"compute_scale_down_enabled": knownvalue.Bool(false),
					"compute_max_instance_size":  knownvalue.StringExact("M20"),
				})...),
			},
			{
				PreConfig:         acc.PreConfigWaitForShardSizeLimitMetrics(t),
				Config:            configDatabaseEdition(projectID, clusterName, new("INFINITE"), 2, new(1024)),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, new(1024)),
			},
			// Raising the limit skips the current-size check, so no wait here.
			{
				Config:            configDatabaseEdition(projectID, clusterName, new("INFINITE"), 2, new(2048)),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, new(2048)),
			},
			{
				PreConfig:         acc.PreConfigWaitForShardSizeLimitMetrics(t),
				Config:            configDatabaseEdition(projectID, clusterName, new("INFINITE"), 2, new(1024)),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, new(1024)),
			},
			{
				Config:            configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, storageConfig, true),
				Check:             resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, new(1024)),
			},
			{
				Config:            configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, "", true),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, nil),
			},
			// Adding a limit where none is set triggers the Atlas current-size check, so wait before it.
			{
				PreConfig: acc.PreConfigWaitForShardSizeLimitMetrics(t),
				Config:    configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 2, new(1024), true),
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
			acc.TestStepImportCluster(resourceName),
		},
	})
}

// TestAccClusterAdvancedCluster_infiniteShardSizeLimitWithZeroNodeAnalytics verifies that clearing
// storage_config works when analytics_specs is explicitly configured with node_count = 0.
// This tests the edge case where analyticsSpecs would be included in the PATCH without instanceSize.
func TestAccClusterAdvancedCluster_infiniteShardSizeLimitWithZeroNodeAnalytics(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseEditionWithZeroNodeAnalytics(projectID, clusterName, new(1024)),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(1024)),
					computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
						"compute_enabled": knownvalue.Bool(false),
					})...),
			},
			{
				Config: configDatabaseEditionWithZeroNodeAnalytics(projectID, clusterName, nil),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil),
					computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
						"compute_enabled": knownvalue.Bool(false),
					})...),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

// TestAccClusterAdvancedCluster_infiniteAnalyticsNodeRemoval verifies that removing analytics nodes
// (node_count = 0) produces a valid PATCH request with the required instance_size.
func TestAccClusterAdvancedCluster_infiniteAnalyticsNodeRemoval(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseEditionWithAnalyticsNodes(projectID, clusterName, 2),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil),
					analyticsNodeCountChecks(clusterName, 2)...),
			},
			{
				Config: configDatabaseEditionWithAnalyticsNodes(projectID, clusterName, 0),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil),
					analyticsNodeCountChecks(clusterName, 0)...),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

// TestAccClusterAdvancedCluster_infiniteReadOnlyNodeRemoval verifies that removing read-only nodes
// (node_count = 0) produces a valid PATCH request.
func TestAccClusterAdvancedCluster_infiniteReadOnlyNodeRemoval(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseEditionWithReadOnlyNodes(projectID, clusterName, 2),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil),
					readOnlyNodeCountChecks(clusterName, 2)...),
			},
			{
				Config: configDatabaseEditionWithReadOnlyNodes(projectID, clusterName, 0),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil),
					readOnlyNodeCountChecks(clusterName, 0)...),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

// TestAccClusterAdvancedCluster_infiniteCombinedRemoval verifies that removing analytics nodes,
// read-only nodes, and storage_config in the same apply produces a valid PATCH request.
func TestAccClusterAdvancedCluster_infiniteCombinedRemoval(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseEditionWithAllNodeTypes(projectID, clusterName, 2, 2, new(1024)),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(1024)),
					append(analyticsNodeCountChecks(clusterName, 2),
						readOnlyNodeCountChecks(clusterName, 2)...)...),
			},
			{
				Config: configDatabaseEditionWithAllNodeTypes(projectID, clusterName, 0, 0, nil),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil),
					append(analyticsNodeCountChecks(clusterName, 0),
						readOnlyNodeCountChecks(clusterName, 0)...)...),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

// TestAccClusterAdvancedCluster_infiniteRemovalWithAutoScaling verifies that removing analytics nodes
// works when compute auto-scaling is enabled.
func TestAccClusterAdvancedCluster_infiniteRemovalWithAutoScaling(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configDatabaseEditionWithAnalyticsAndAutoScaling(projectID, clusterName, 2, true),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil),
					analyticsNodeCountChecks(clusterName, 2)...),
			},
			{
				Config: configDatabaseEditionWithAnalyticsAndAutoScaling(projectID, clusterName, 0, true),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil),
					analyticsNodeCountChecks(clusterName, 0)...),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_infiniteShardSizeLimitErrors(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 3)
	storageConfig := "compute_enabled = false\n" + databaseEditionStorageConfig(new(1024))
	storageOnlyRecovery := acc.TestStepCheckEmptyPlan(configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, storageConfig, true))
	storageOnlyRecovery.ConfigStateChecks = shardSizeLimitChecks(clusterName, new(1024))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:            configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, storageConfig, true),
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
				// Four electable nodes are invalid in both supported Infinite topology modes.
				Config:      configDatabaseEditionWithComputeAutoScaling(projectID, clusterName, new("INFINITE"), 4, new(2048), true),
				ExpectError: regexp.MustCompile(`(?s)INVALID_ATTRIBUTE.*number of electable nodes`),
			},
			storageOnlyRecovery,
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
	// Start without storage to mirror the CLOUDP-449163 repro.
	computeOnlyConfig := configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, computeConfig, false)
	configWithStorage := configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, scaleDownConfig+databaseEditionStorageConfig(new(1024)), false)
	effectiveBase := newInfiniteEffectiveReq(projectID, clusterName).withInstanceSize("M10").withEffectiveInstanceSize("M10").withFlag()
	configWithStorageEffective := effectiveBase.withAutoScaling(
		scaleDownConfig+databaseEditionStorageConfig(new(1024)),
		map[string]knownvalue.Check{
			"compute_enabled":           knownvalue.Bool(true),
			"compute_max_instance_size": knownvalue.StringExact("M20"),
		}).withStorageLimit(1024)
	configWithStorageEffectiveComputeUpdated := effectiveBase.withAutoScaling(
		strings.ReplaceAll(scaleDownConfig+databaseEditionStorageConfig(new(1024)), "M20", "M30_GEN_2"),
		map[string]knownvalue.Check{
			"compute_enabled":           knownvalue.Bool(true),
			"compute_max_instance_size": knownvalue.StringExact("M30_GEN_2"),
		}).withStorageLimit(1024)
	configWithStorageEffectiveUpdated := effectiveBase.withAutoScaling(
		scaleDownConfig+databaseEditionStorageConfig(new(2048)),
		map[string]knownvalue.Check{
			"compute_enabled":            knownvalue.Bool(true),
			"compute_scale_down_enabled": knownvalue.Bool(true),
			"compute_min_instance_size":  knownvalue.StringExact("M10"),
			"compute_max_instance_size":  knownvalue.StringExact("M20"),
		}).withStorageLimit(2048)
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
			// Adding a limit where none is set triggers the Atlas current-size check, so wait before it.
			{
				PreConfig:         acc.PreConfigWaitForShardSizeLimitMetrics(t),
				Config:            configWithStorage,
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, new(1024)), computeChecks...),
			},
			// Enabling use_effective_fields on an existing cluster must preserve the configured hardware in state.
			{
				Config:            configWithStorageEffective.config(),
				ConfigStateChecks: configWithStorageEffective.check(),
			},
			// A compute-only PATCH with the flag on must keep returning hardware specs in the response.
			{
				Config:            configWithStorageEffectiveComputeUpdated.config(),
				ConfigStateChecks: configWithStorageEffectiveComputeUpdated.check(),
			},
			// Raising the limit skips the current-size check, so no wait here.
			{
				Config:            configWithStorageEffectiveUpdated.config(),
				ConfigStateChecks: configWithStorageEffectiveUpdated.check(),
			},
			// Import doesn't send the flag so non-effective specs not in the config are set in the imported state.
			acc.TestStepImportCluster(resourceName, "use_effective_fields", "replication_specs"),
		},
	})
}

// TestAccClusterAdvancedCluster_infiniteEffectiveFields exercises use_effective_fields on an INFINITE cluster
// with the flag set at creation: configured hardware must be echoed in state while effective specs report
// actual running values, and updates with the flag on must keep returning hardware specs. The cluster starts
// with electable compute auto-scaling only (the typical case) and adds analytics nodes mid-test.
func TestAccClusterAdvancedCluster_infiniteEffectiveFields(t *testing.T) {
	base := baseInfiniteEffectiveReq(t)
	var (
		initial        = base.withUnsetSpecsNull()
		updated        = base.withInstanceSize("M20")
		computeUpdated = updated.withComputeAutoScaling("M30_GEN_2")
		analyticsAdded = computeUpdated.withAnalytics()
		backToRunning  = analyticsAdded.withInstanceSize("M10").withComputeAutoScaling("M20")
		flagOff        = backToRunning.withoutFlag()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, initial.projectID, initial.clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				// Create with the flag on and no analytics nodes: unsent read_only_specs and
				// analytics_specs must not be echoed into state.
				Config:            initial.config(),
				Check:             checkDatabaseEdition(new("INFINITE"), "INFINITE"),
				ConfigStateChecks: initial.check(),
			},
			{
				// With compute auto-scaling enabled the instance size update is echoed in state but not
				// applied, so configured (M20) and effective (M10) values deliberately differ.
				Config:            updated.config(),
				ConfigStateChecks: updated.check(),
			},
			{
				// Compute-only PATCH with the flag on must preserve hardware specs in the response.
				Config:            computeUpdated.config(),
				ConfigStateChecks: computeUpdated.check(),
			},
			{
				// Adding analytics nodes with the flag on must populate analytics_specs and effective_analytics_specs.
				Config:            analyticsAdded.config(),
				ConfigStateChecks: analyticsAdded.check(),
			},
			{
				// Back to the running instance size so turning the flag off converges on current values.
				Config:            backToRunning.config(),
				ConfigStateChecks: backToRunning.check(),
			},
			{
				// Toggle the flag off: state now shows the current hardware values.
				Config:            flagOff.config(),
				ConfigStateChecks: flagOff.check(),
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
				// INFINITE requires an explicit compute choice on create; false exercises an inert auto_scaling block.
				Config:            configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, "compute_enabled = false", false),
				ConfigStateChecks: shardSizeLimitChecks(clusterName, nil),
			},
			{
				Config: configDatabaseEditionWithAutoScaling(projectID, clusterName, new("INFINITE"), 2, "compute_enabled = true\ncompute_max_instance_size = \"M20\"", false),
				ConfigStateChecks: append(shardSizeLimitChecks(clusterName, nil), computeAutoScalingChecks(clusterName, map[string]knownvalue.Check{
					"compute_enabled":           knownvalue.Bool(true),
					"compute_max_instance_size": knownvalue.StringExact("M20"),
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
			// INFINITE requires an explicit auto_scaling.compute choice on create; keep it on updates so
			// the block stays configured, but omit it when the test exercises omitted computed blocks.
			autoScaling := "compute_enabled = false"
			if limit != nil {
				autoScaling += "\n" + databaseEditionStorageConfig(limit)
			}
			regionConfig += "\nauto_scaling = {" + autoScaling + "\n}"
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
			{
				PreConfig:         acc.PreConfigWaitForShardSizeLimitMetrics(t),
				Config:            baseConfig,
				ConfigStateChecks: checks(new(1024), "M20"),
			},
			// Clearing storage must retain planned node counts and zone metadata when omitted from configuration.
			{Config: partialHardwareConfig, ConfigStateChecks: checks(nil, "M20")},
			{
				PreConfig:         acc.PreConfigWaitForShardSizeLimitMetrics(t),
				Config:            baseConfig,
				ConfigStateChecks: checks(new(1024), "M20"),
			},
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

func configDatabaseEdition(projectID, clusterName string, databaseEdition *string, nodeCount int, shardSizeLimitGB *int) string {
	autoScaling := "compute_enabled = false"
	if shardSizeLimitGB != nil {
		autoScaling += "\n" + databaseEditionStorageConfig(shardSizeLimitGB)
	}
	return configDatabaseEditionWithAutoScaling(projectID, clusterName, databaseEdition, nodeCount, autoScaling, false)
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

// configDatabaseEditionWithZeroNodeAnalytics creates a config with analytics_specs explicitly set to node_count = 0.
// This tests the edge case where analyticsSpecs would be included in the PATCH without instanceSize.
func configDatabaseEditionWithZeroNodeAnalytics(projectID, clusterName string, shardSizeLimitGB *int) string {
	storage := ""
	if shardSizeLimitGB != nil {
		storage = fmt.Sprintf(`
			storage_config = {
				shard_size_limit_gb = %d
			}`, *shardSizeLimitGB)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id     = %[1]q
			name           = %[2]q
			cluster_type   = "REPLICASET"
			database_edition = "INFINITE"
			backup_enabled = true
			pit_enabled    = true

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = "M10"
						node_count    = 2
					}
					analytics_specs = {
						node_count    = 0
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
					auto_scaling = {
						compute_enabled = false
						%[3]s
					}
				}]
			}]
		}
	`, projectID, clusterName, storage) + dataSourcesConfig
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

// dataSourcesConfigEffective mirrors dataSourcesConfig with use_effective_fields set on both data sources.
const dataSourcesConfigEffective = `
data "mongodbatlas_advanced_cluster" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	name 	     = mongodbatlas_advanced_cluster.test.name
	use_effective_fields = true
	depends_on = [mongodbatlas_advanced_cluster.test]
}

data "mongodbatlas_advanced_clusters" "test" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	use_effective_fields = true
	depends_on = [mongodbatlas_advanced_cluster.test]
}`

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

// infiniteEffectiveSpecsChecks asserts the use_effective_fields contract on an INFINITE cluster: configured
// hardware is echoed in state while effective specs report actual running values, with configured and
// effective instance sizes deliberately differing. Effective spec attributes exist only in the data sources.
// Disk fields are not asserted: Atlas omits them from INFINITE effective specs.
func infiniteEffectiveSpecsChecks(clusterName, configuredInstanceSize, effectiveInstanceSize string, flagEnabled bool) []statecheck.StateCheck {
	regionPath := tfjsonpath.New("replication_specs").AtSliceIndex(0).AtMapKey("region_configs").AtSliceIndex(0)
	electable := knownvalue.ObjectPartial(map[string]knownvalue.Check{
		"instance_size": knownvalue.StringExact(configuredInstanceSize),
		"node_count":    knownvalue.Int64Exact(2),
	})
	effective := knownvalue.ObjectPartial(map[string]knownvalue.Check{
		"instance_size": knownvalue.StringExact(effectiveInstanceSize),
		"node_count":    knownvalue.Int64Exact(2),
	})
	pluralChecks := map[string]knownvalue.Check{
		"replication_specs.0.region_configs.0.electable_specs":           electable,
		"replication_specs.0.region_configs.0.effective_electable_specs": effective,
	}
	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resourceName, regionPath.AtMapKey("electable_specs"), electable),
		statecheck.ExpectKnownValue(dataSourceName, regionPath.AtMapKey("electable_specs"), electable),
		statecheck.ExpectKnownValue(dataSourceName, regionPath.AtMapKey("effective_electable_specs"), effective),
	}
	if flagEnabled {
		checks = append(checks, statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("use_effective_fields"), knownvalue.Bool(true)))
		pluralChecks["use_effective_fields"] = knownvalue.Bool(true)
	}
	return append(checks, acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), pluralChecks))
}

// infiniteEffectiveAnalyticsChecks asserts the same echo/effective contract for INFINITE analytics nodes.
func infiniteEffectiveAnalyticsChecks(clusterName, configuredInstanceSize, effectiveInstanceSize string) []statecheck.StateCheck {
	regionPath := tfjsonpath.New("replication_specs").AtSliceIndex(0).AtMapKey("region_configs").AtSliceIndex(0)
	analytics := knownvalue.ObjectPartial(map[string]knownvalue.Check{
		"instance_size": knownvalue.StringExact(configuredInstanceSize),
		"node_count":    knownvalue.Int64Exact(1),
	})
	effective := knownvalue.ObjectPartial(map[string]knownvalue.Check{
		"instance_size": knownvalue.StringExact(effectiveInstanceSize),
		"node_count":    knownvalue.Int64Exact(1),
	})
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resourceName, regionPath.AtMapKey("analytics_specs"), analytics),
		statecheck.ExpectKnownValue(dataSourceName, regionPath.AtMapKey("analytics_specs"), analytics),
		statecheck.ExpectKnownValue(dataSourceName, regionPath.AtMapKey("effective_analytics_specs"), effective),
		acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), map[string]knownvalue.Check{
			"replication_specs.0.region_configs.0.analytics_specs":           analytics,
			"replication_specs.0.region_configs.0.effective_analytics_specs": effective,
		}),
	}
}

// infiniteEffectiveReq builds INFINITE configs and state checks for use_effective_fields lifecycle steps,
// mirroring the builder style of effectiveReq in effective_fields_test.go.
type infiniteEffectiveReq struct {
	storageLimitGB        *int
	autoScalingChecks     map[string]knownvalue.Check
	projectID             string
	clusterName           string
	instanceSize          string // configured electable size, echoed in state while the flag is on
	autoScalingAttributes string
	effectiveInstanceSize string // actual running size, lags the configured one while auto-scaling is enabled
	analytics             bool
	useEffectiveFields    bool
	expectUnsetSpecsNull  bool
}

func newInfiniteEffectiveReq(projectID, clusterName string) infiniteEffectiveReq {
	return infiniteEffectiveReq{projectID: projectID, clusterName: clusterName}
}

func baseInfiniteEffectiveReq(t *testing.T) infiniteEffectiveReq {
	t.Helper()
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 3) // 2 electable + 1 analytics nodes
	return newInfiniteEffectiveReq(projectID, clusterName).
		withInstanceSize("M10").
		withEffectiveInstanceSize("M10").
		withComputeAutoScaling("M20").
		withFlag()
}

func (req infiniteEffectiveReq) withInstanceSize(instanceSize string) infiniteEffectiveReq {
	req.instanceSize = instanceSize
	return req
}

func (req infiniteEffectiveReq) withEffectiveInstanceSize(instanceSize string) infiniteEffectiveReq {
	req.effectiveInstanceSize = instanceSize
	return req
}

func (req infiniteEffectiveReq) withAutoScaling(attributes string, checks map[string]knownvalue.Check) infiniteEffectiveReq {
	req.autoScalingAttributes = attributes
	req.autoScalingChecks = checks
	return req
}

// withComputeAutoScaling sets the electable compute auto-scaling block and its expected state for a max size.
func (req infiniteEffectiveReq) withComputeAutoScaling(maxInstanceSize string) infiniteEffectiveReq {
	attributes := fmt.Sprintf(`
						compute_enabled            = true
						compute_scale_down_enabled = false
						compute_max_instance_size  = %q`, maxInstanceSize)
	return req.withAutoScaling(attributes, map[string]knownvalue.Check{
		"compute_enabled":            knownvalue.Bool(true),
		"compute_scale_down_enabled": knownvalue.Bool(false),
		"compute_max_instance_size":  knownvalue.StringExact(maxInstanceSize),
	})
}

func (req infiniteEffectiveReq) withStorageLimit(limitGB int) infiniteEffectiveReq {
	req.storageLimitGB = &limitGB
	return req
}

func (req infiniteEffectiveReq) withAnalytics() infiniteEffectiveReq {
	req.analytics = true
	return req
}

func (req infiniteEffectiveReq) withFlag() infiniteEffectiveReq {
	req.useEffectiveFields = true
	return req
}

func (req infiniteEffectiveReq) withoutFlag() infiniteEffectiveReq {
	req.useEffectiveFields = false
	req.expectUnsetSpecsNull = false
	return req
}

// withUnsetSpecsNull asserts that unconfigured read_only_specs and analytics_specs stay null. It only holds
// for a cluster created with the flag on; enabling the flag on an existing cluster can echo empty maps.
func (req infiniteEffectiveReq) withUnsetSpecsNull() infiniteEffectiveReq {
	req.expectUnsetSpecsNull = true
	return req
}

func (req infiniteEffectiveReq) config() string {
	autoScalingConfig := ""
	if req.autoScalingAttributes != "" {
		autoScalingConfig = fmt.Sprintf(`
					auto_scaling = {
						%s
					}`, req.autoScalingAttributes)
	}
	analyticsConfig := ""
	if req.analytics {
		analyticsConfig = `
					analytics_specs = {
						instance_size = "M10"
						node_count    = 1
					}
					analytics_auto_scaling = {
						compute_enabled           = true
						compute_max_instance_size = "M20"
					}`
	}
	effectiveFieldsConfig := ""
	dataSources := dataSourcesConfig
	if req.useEffectiveFields {
		effectiveFieldsConfig = "use_effective_fields = true"
		dataSources = dataSourcesConfigEffective
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id       = %[1]q
			name             = %[2]q
			cluster_type     = "REPLICASET"
			backup_enabled   = true
			pit_enabled      = true
			database_edition = "INFINITE"
			%[3]s

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = %[4]q
						node_count    = 2
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
					%[5]s
					%[6]s
				}]
			}]
		}
	`, req.projectID, req.clusterName, effectiveFieldsConfig, req.instanceSize, autoScalingConfig, analyticsConfig) + dataSources
}

func (req infiniteEffectiveReq) check() []statecheck.StateCheck {
	var checks []statecheck.StateCheck
	if req.storageLimitGB != nil {
		checks = append(checks, shardSizeLimitChecks(req.clusterName, req.storageLimitGB)...)
	}
	if req.autoScalingChecks != nil {
		checks = append(checks, computeAutoScalingChecks(req.clusterName, req.autoScalingChecks)...)
	}
	checks = append(checks, infiniteEffectiveSpecsChecks(req.clusterName, req.instanceSize, req.effectiveInstanceSize, req.useEffectiveFields)...)
	regionConfigPath := tfjsonpath.New("replication_specs").AtSliceIndex(0).AtMapKey("region_configs").AtSliceIndex(0)
	if req.analytics {
		checks = append(checks, infiniteEffectiveAnalyticsChecks(req.clusterName, "M10", "M10")...)
		checks = append(checks, infiniteAutoScalingChecks(req.clusterName, "analytics_auto_scaling", map[string]knownvalue.Check{
			"compute_enabled":           knownvalue.Bool(true),
			"compute_max_instance_size": knownvalue.StringExact("M20"),
		})...)
	} else if req.expectUnsetSpecsNull {
		checks = append(checks, statecheck.ExpectKnownValue(resourceName, regionConfigPath.AtMapKey("analytics_specs"), knownvalue.Null()))
	}
	if req.expectUnsetSpecsNull {
		checks = append(checks, statecheck.ExpectKnownValue(resourceName, regionConfigPath.AtMapKey("read_only_specs"), knownvalue.Null()))
	}
	return checks
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

// configDatabaseEditionWithAnalyticsNodes creates a config with analytics_specs set to the specified node_count.
func configDatabaseEditionWithAnalyticsNodes(projectID, clusterName string, analyticsNodeCount int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id     = %[1]q
			name           = %[2]q
			cluster_type   = "REPLICASET"
			database_edition = "INFINITE"
			backup_enabled = true
			pit_enabled    = true

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = "M10"
						node_count    = 2
					}
					analytics_specs = {
						instance_size = "M10"
						node_count    = %[3]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
					auto_scaling = {
						compute_enabled = false
					}
				}]
			}]
		}
	`, projectID, clusterName, analyticsNodeCount) + dataSourcesConfig
}

// configDatabaseEditionWithReadOnlyNodes creates a config with read_only_specs set to the specified node_count.
func configDatabaseEditionWithReadOnlyNodes(projectID, clusterName string, readOnlyNodeCount int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id     = %[1]q
			name           = %[2]q
			cluster_type   = "REPLICASET"
			database_edition = "INFINITE"
			backup_enabled = true
			pit_enabled    = true

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = "M10"
						node_count    = 2
					}
					read_only_specs = {
						instance_size = "M10"
						node_count    = %[3]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
					auto_scaling = {
						compute_enabled = false
					}
				}]
			}]
		}
	`, projectID, clusterName, readOnlyNodeCount) + dataSourcesConfig
}

// configDatabaseEditionWithAllNodeTypes creates a config with analytics_specs, read_only_specs, and storage_config.
func configDatabaseEditionWithAllNodeTypes(projectID, clusterName string, analyticsNodeCount, readOnlyNodeCount int, shardSizeLimitGB *int) string {
	storage := ""
	if shardSizeLimitGB != nil {
		storage = fmt.Sprintf(`
			storage_config = {
				shard_size_limit_gb = %d
			}`, *shardSizeLimitGB)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id     = %[1]q
			name           = %[2]q
			cluster_type   = "REPLICASET"
			database_edition = "INFINITE"
			backup_enabled = true
			pit_enabled    = true

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = "M10"
						node_count    = 2
					}
					analytics_specs = {
						instance_size = "M10"
						node_count    = %[3]d
					}
					read_only_specs = {
						instance_size = "M10"
						node_count    = %[4]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
					auto_scaling = {
						compute_enabled = false
						%[5]s
					}
				}]
			}]
		}
	`, projectID, clusterName, analyticsNodeCount, readOnlyNodeCount, storage) + dataSourcesConfig
}

// configDatabaseEditionWithAnalyticsAndAutoScaling creates a config with analytics_specs and compute auto-scaling.
func configDatabaseEditionWithAnalyticsAndAutoScaling(projectID, clusterName string, analyticsNodeCount int, withAutoScaling bool) string {
	autoScaling := "compute_enabled = false"
	if withAutoScaling {
		autoScaling = `
			compute_enabled            = true
			compute_scale_down_enabled = false
			compute_max_instance_size  = "M20"`
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id     = %[1]q
			name           = %[2]q
			cluster_type   = "REPLICASET"
			database_edition = "INFINITE"
			backup_enabled = true
			pit_enabled    = true

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = "M10"
						node_count    = 2
					}
					analytics_specs = {
						instance_size = "M10"
						node_count    = %[3]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
					auto_scaling = {
						%[4]s
					}
				}]
			}]
		}
	`, projectID, clusterName, analyticsNodeCount, autoScaling) + dataSourcesConfig
}

// analyticsNodeCountChecks returns checks for the analytics_specs node_count.
func analyticsNodeCountChecks(clusterName string, nodeCount int) []statecheck.StateCheck {
	path := tfjsonpath.New("replication_specs").AtSliceIndex(0).AtMapKey("region_configs").AtSliceIndex(0).AtMapKey("analytics_specs")
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resourceName, path.AtMapKey("node_count"), knownvalue.Int64Exact(int64(nodeCount))),
		statecheck.ExpectKnownValue(dataSourceName, path.AtMapKey("node_count"), knownvalue.Int64Exact(int64(nodeCount))),
		acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), map[string]knownvalue.Check{
			"replication_specs.0.region_configs.0.analytics_specs.node_count": knownvalue.Int64Exact(int64(nodeCount)),
		}),
	}
}

// readOnlyNodeCountChecks returns checks for the read_only_specs node_count.
func readOnlyNodeCountChecks(clusterName string, nodeCount int) []statecheck.StateCheck {
	path := tfjsonpath.New("replication_specs").AtSliceIndex(0).AtMapKey("region_configs").AtSliceIndex(0).AtMapKey("read_only_specs")
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resourceName, path.AtMapKey("node_count"), knownvalue.Int64Exact(int64(nodeCount))),
		statecheck.ExpectKnownValue(dataSourceName, path.AtMapKey("node_count"), knownvalue.Int64Exact(int64(nodeCount))),
		acc.PluralResultCheck(dataSourcePluralName, "name", knownvalue.StringExact(clusterName), map[string]knownvalue.Check{
			"replication_specs.0.region_configs.0.read_only_specs.node_count": knownvalue.Int64Exact(int64(nodeCount)),
		}),
	}
}
