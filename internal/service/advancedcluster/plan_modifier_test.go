package advancedcluster_test

import (
	"context"
	"fmt"
	"maps"
	"testing"

	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

var (
	repSpec0      = tfjsonpath.New("replication_specs").AtSliceIndex(0)
	repSpec1      = tfjsonpath.New("replication_specs").AtSliceIndex(1)
	regionConfig0 = repSpec0.AtMapKey("region_configs").AtSliceIndex(0)
	regionConfig1 = repSpec1.AtMapKey("region_configs").AtSliceIndex(0)
)

func autoScalingKnownValue(computeEnabled, diskEnabled, scaleDown bool, minInstanceSize, maxInstanceSize string, includeStorageConfig bool) knownvalue.Check {
	attributes := map[string]knownvalue.Check{
		"compute_enabled":            knownvalue.Bool(computeEnabled),
		"disk_gb_enabled":            knownvalue.Bool(diskEnabled),
		"compute_scale_down_enabled": knownvalue.Bool(scaleDown),
		"compute_min_instance_size":  knownvalue.StringExact(minInstanceSize),
		"compute_max_instance_size":  knownvalue.StringExact(maxInstanceSize),
	}
	if includeStorageConfig {
		attributes["storage_config"] = knownvalue.Null()
	}
	return knownvalue.ObjectExact(attributes)
}

func specInstanceSizeNodeCount(instanceSize string, nodeCount int) knownvalue.Check {
	return knownvalue.ObjectPartial(map[string]knownvalue.Check{
		"instance_size": knownvalue.StringExact(instanceSize),
		"node_count":    knownvalue.Int64Exact(int64(nodeCount)),
	})
}

func TestPlanChecksClusterTwoRepSpecsWithAutoScalingAndSpecs(t *testing.T) {
	var (
		baseConfig                  = unit.NewMockPlanChecksConfig(t, &mockConfig, unit.ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs)
		resourceName                = baseConfig.ResourceName
		autoScalingEnabled          = autoScalingKnownValue(true, true, true, "M10", "M30", true)
		analyticsAutoScalingEnabled = autoScalingKnownValue(true, true, true, "M10", "M30", false)
		testCases                   = []unit.PlanCheckTest{
			{
				ConfigFilename: "main_removed_blocks_from_config_no_plan_changes.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
				},
			},
			{
				// Every Optional-only attribute is unresolved, it must stay unknown instead of using the state value.
				ConfigFilename: "main_unresolved_optional_attributes.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("tags")),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("labels")),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("retain_backups_enabled")),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("adaptive_capacity")),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("use_effective_fields")),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("pinned_fcv")),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("timeouts")),
					plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("accept_data_risks_and_force_replica_set_reconfig")),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("backing_provider_name")),
					// Only the unresolved attributes stay unknown: the computed ones still keep their state value.
					plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("cluster_id"), knownvalue.StringExact("67d01a2d01d3561b07caf76e")),
					plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("encryption_at_rest_provider"), knownvalue.StringExact("NONE")),
					plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("root_cert_type"), knownvalue.StringExact("ISRGROOTX1")),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("electable_specs").AtMapKey("instance_size"), knownvalue.StringExact("M10")),
				},
			},
			{
				// auto_scaling is Optional+Computed, so the schema rule in CopyUnknowns does not keep it unknown.
				// Planning its state value, and the read_only_specs derived from it, is rejected during apply.
				ConfigFilename: "main_unresolved_optional_computed.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("auto_scaling")),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("read_only_specs")),
					plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("cluster_id"), knownvalue.StringExact("67d01a2d01d3561b07caf76e")),
				},
			},
			{
				ConfigFilename: "main_node_count_unknown.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("read_only_specs").AtMapKey("node_count"), knownvalue.Int64Exact(2)),
				},
			},
			{
				ConfigFilename: "main_removed_blocks_from_config_and_instance_change.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					// checks regionConfig0
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M10", 2)),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("electable_specs"), specInstanceSizeNodeCount("M10", 5)),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("analytics_auto_scaling"), analyticsAutoScalingEnabled),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("analytics_specs")), // analytics specs was defined in region_configs.0 but not in region_configs.1

					// checks regionConfig1
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M20", 1)),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("electable_specs"), specInstanceSizeNodeCount("M20", 3)),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("analytics_auto_scaling"), analyticsAutoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("analytics_specs"), knownvalue.NotNull()),
				},
			},
		}
	)
	for _, testCase := range testCases {
		t.Run(testCase.ConfigFilename, func(t *testing.T) {
			unit.MockPlanChecksAndRun(t, baseConfig.WithPlanCheckTest(testCase))
		})
	}
}

func TestPlanRemoveAutoScalingStorageConfig(t *testing.T) {
	ctx := context.Background()
	schemaResponse := frameworkresource.SchemaResponse{}
	advancedcluster.Resource().Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResponse)
	typ := schemaResponse.Schema.Type().TerraformType(ctx)
	model := func(autoScaling any) tftypes.Value {
		return planTestValue(typ, map[string]any{
			"name": "example", "project_id": "111111111111111111111111", "cluster_type": "REPLICASET", "database_edition": "INFINITE",
			"replication_specs": []any{map[string]any{"region_configs": []any{map[string]any{
				"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7),
				"electable_specs": map[string]any{"instance_size": "M10", "node_count": int64(2)},
				"auto_scaling":    autoScaling,
			}}}},
		})
	}
	dynamic := func(value tftypes.Value) *tfprotov6.DynamicValue {
		result, err := tfprotov6.NewDynamicValue(typ, value)
		require.NoError(t, err)
		return &result
	}
	storagePath := tftypes.NewAttributePath().WithAttributeName("replication_specs").WithElementKeyInt(0).
		WithAttributeName("region_configs").WithElementKeyInt(0).WithAttributeName("auto_scaling").WithAttributeName("storage_config")
	for _, computeEnabled := range []bool{false, true} {
		attributes := map[string]any{
			"storage_config": map[string]any{"shard_size_limit_gb": int64(1024)},
		}
		if computeEnabled {
			attributes["compute_enabled"] = true
			attributes["compute_scale_down_enabled"] = false
			attributes["compute_max_instance_size"] = "M20"
		}
		withoutStorage := maps.Clone(attributes)
		delete(withoutStorage, "storage_config")
		for name, tc := range map[string]struct{ config, plan any }{
			"remove parent with unknown plan": {nil, tftypes.UnknownValue},
			"remove parent with known plan":   {nil, attributes},
			"remove storage only":             {withoutStorage, withoutStorage},
		} {
			t.Run(fmt.Sprintf("compute=%t/%s", computeEnabled, name), func(t *testing.T) {
				server, err := acc.TestAccProviderV6Factories["mongodbatlas"]()
				require.NoError(t, err)
				result, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
					TypeName:   "mongodbatlas_advanced_cluster",
					PriorState: dynamic(model(attributes)), ProposedNewState: dynamic(model(tc.plan)), Config: dynamic(model(tc.config)),
				})
				require.NoError(t, err)
				require.Empty(t, result.Diagnostics)
				plan, err := result.PlannedState.Unmarshal(typ)
				require.NoError(t, err)
				storage, _, err := tftypes.WalkAttributePath(plan, storagePath)
				require.NoError(t, err)
				require.True(t, storage.(tftypes.Value).IsNull(), "removing storage_config or auto_scaling must clear the configured shard limit")
				if computeEnabled {
					computePath := tftypes.NewAttributePath().WithAttributeName("replication_specs").WithElementKeyInt(0).
						WithAttributeName("region_configs").WithElementKeyInt(0).WithAttributeName("auto_scaling").WithAttributeName("compute_enabled")
					compute, _, err := tftypes.WalkAttributePath(plan, computePath)
					require.NoError(t, err)
					require.Equal(t, tftypes.NewValue(tftypes.Bool, true), compute)
				}
			})
		}
	}
}

// planTestValue fills omitted attributes with typed nulls so the fixture follows the current resource schema.
func planTestValue(typ tftypes.Type, value any) tftypes.Value {
	if value == nil || value == tftypes.UnknownValue {
		return tftypes.NewValue(typ, value)
	}
	switch shape := typ.(type) {
	case tftypes.Object:
		input := value.(map[string]any)
		out := map[string]tftypes.Value{}
		for name, childType := range shape.AttributeTypes {
			out[name] = planTestValue(childType, input[name])
		}
		return tftypes.NewValue(typ, out)
	case tftypes.List:
		out := []tftypes.Value{}
		for _, item := range value.([]any) {
			out = append(out, planTestValue(shape.ElementType, item))
		}
		return tftypes.NewValue(typ, out)
	default:
		return tftypes.NewValue(typ, value)
	}
}

// TestPlanUnknownConfigKeepsStateCopy pins the scope of unknownInConfig: an unresolved expression must keep
// its own attribute unknown without discarding the state values that ModifyPlan copies into everything else.
func TestPlanUnknownConfigKeepsStateCopy(t *testing.T) {
	ctx := t.Context()
	schemaResponse := frameworkresource.SchemaResponse{}
	advancedcluster.Resource().Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResponse)
	typ := schemaResponse.Schema.Type().TerraformType(ctx)
	computedAttributes := []string{"cluster_id", "backup_enabled", "encryption_at_rest_provider", "root_cert_type"}
	model := func(instanceSize string, majorVersion any) map[string]any {
		return map[string]any{
			"name": "example", "project_id": "111111111111111111111111", "cluster_type": "REPLICASET",
			"mongo_db_major_version": majorVersion,
			"replication_specs": []any{map[string]any{"region_configs": []any{map[string]any{
				"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7),
				"electable_specs": map[string]any{"instance_size": instanceSize, "node_count": int64(3)},
			}}}},
		}
	}
	dynamic := func(value map[string]any) *tfprotov6.DynamicValue {
		result, err := tfprotov6.NewDynamicValue(typ, planTestValue(typ, value))
		require.NoError(t, err)
		return &result
	}
	prior := model("M10", "8.0")
	maps.Copy(prior, map[string]any{
		"cluster_id": "333333333333333333333333", "backup_enabled": true,
		"encryption_at_rest_provider": "NONE", "root_cert_type": "ISRGROOTX1",
	})
	for name, majorVersion := range map[string]any{"known": "8.0", "unknown": tftypes.UnknownValue} {
		t.Run("mongo_db_major_version="+name, func(t *testing.T) {
			proposed := model("M20", majorVersion)
			for _, attribute := range computedAttributes {
				proposed[attribute] = tftypes.UnknownValue
			}
			server, err := acc.TestAccProviderV6Factories["mongodbatlas"]()
			require.NoError(t, err)
			result, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
				TypeName: "mongodbatlas_advanced_cluster", PriorState: dynamic(prior),
				ProposedNewState: dynamic(proposed), Config: dynamic(model("M20", majorVersion)),
			})
			require.NoError(t, err)
			require.Empty(t, result.Diagnostics)
			plan, err := result.PlannedState.Unmarshal(typ)
			require.NoError(t, err)
			for _, attribute := range computedAttributes {
				value, _, err := tftypes.WalkAttributePath(plan, tftypes.NewAttributePath().WithAttributeName(attribute))
				require.NoError(t, err)
				require.True(t, value.(tftypes.Value).IsKnown(),
					"%s must keep its state value: only the unresolved attribute needs to stay unknown", attribute)
			}
		})
	}
}
