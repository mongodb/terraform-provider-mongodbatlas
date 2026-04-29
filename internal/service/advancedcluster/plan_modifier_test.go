package advancedcluster_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

// Attr type maps mirrors the unexported vars in schema.go to avoid uppercasing them.
// They are redeclared here because this package is advancedcluster_test and cannot access package-private symbols.
// Done in order to be able to add unit tests for WarnIgnoredSpecChange function, as they allow to test the function without depending on acceptance tests
var (
	autoScalingAttrTypes = map[string]attr.Type{
		"compute_enabled":            types.BoolType,
		"compute_max_instance_size":  types.StringType,
		"compute_min_instance_size":  types.StringType,
		"compute_scale_down_enabled": types.BoolType,
		"disk_gb_enabled":            types.BoolType,
	}
	specsAttrTypes = map[string]attr.Type{
		"disk_iops":       types.Int64Type,
		"disk_size_gb":    types.Float64Type,
		"ebs_volume_type": types.StringType,
		"instance_size":   types.StringType,
		"node_count":      types.Int64Type,
	}
	regionConfigAttrTypes = map[string]attr.Type{
		"analytics_auto_scaling": types.ObjectType{AttrTypes: autoScalingAttrTypes},
		"analytics_specs":        types.ObjectType{AttrTypes: specsAttrTypes},
		"auto_scaling":           types.ObjectType{AttrTypes: autoScalingAttrTypes},
		"backing_provider_name":  types.StringType,
		"electable_specs":        types.ObjectType{AttrTypes: specsAttrTypes},
		"priority":               types.Int64Type,
		"provider_name":          types.StringType,
		"read_only_specs":        types.ObjectType{AttrTypes: specsAttrTypes},
		"region_name":            types.StringType,
	}
	replicationSpecAttrTypes = map[string]attr.Type{
		"container_id":   types.MapType{ElemType: types.StringType},
		"external_id":    types.StringType,
		"region_configs": types.ListType{ElemType: types.ObjectType{AttrTypes: regionConfigAttrTypes}},
		"zone_id":        types.StringType,
		"zone_name":      types.StringType,
	}
)

var (
	repSpec0      = tfjsonpath.New("replication_specs").AtSliceIndex(0)
	repSpec1      = tfjsonpath.New("replication_specs").AtSliceIndex(1)
	regionConfig0 = repSpec0.AtMapKey("region_configs").AtSliceIndex(0)
	regionConfig1 = repSpec1.AtMapKey("region_configs").AtSliceIndex(0)
)

func autoScalingKnownValue(computeEnabled, diskEnabled, scaleDown bool, minInstanceSize, maxInstanceSize string) knownvalue.Check {
	return knownvalue.ObjectExact(map[string]knownvalue.Check{
		"compute_enabled":            knownvalue.Bool(computeEnabled),
		"disk_gb_enabled":            knownvalue.Bool(diskEnabled),
		"compute_scale_down_enabled": knownvalue.Bool(scaleDown),
		"compute_min_instance_size":  knownvalue.StringExact(minInstanceSize),
		"compute_max_instance_size":  knownvalue.StringExact(maxInstanceSize),
	})
}

func specInstanceSizeNodeCount(instanceSize string, nodeCount int) knownvalue.Check {
	return knownvalue.ObjectPartial(map[string]knownvalue.Check{
		"instance_size": knownvalue.StringExact(instanceSize),
		"node_count":    knownvalue.Int64Exact(int64(nodeCount)),
	})
}

func TestPlanChecksClusterTwoRepSpecsWithAutoScalingAndSpecs(t *testing.T) {
	var (
		baseConfig         = unit.NewMockPlanChecksConfig(t, &mockConfig, unit.ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs)
		resourceName       = baseConfig.ResourceName
		autoScalingEnabled = autoScalingKnownValue(true, true, true, "M10", "M30")
		testCases          = []unit.PlanCheckTest{
			{
				ConfigFilename: "main_removed_blocks_from_config_no_plan_changes.tf",
				Checks: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
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
					plancheck.ExpectKnownValue(resourceName, regionConfig0.AtMapKey("analytics_auto_scaling"), autoScalingEnabled),
					plancheck.ExpectUnknownValue(resourceName, regionConfig0.AtMapKey("analytics_specs")), // analytics specs was defined in region_configs.0 but not in region_configs.1

					// checks regionConfig1
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("read_only_specs"), specInstanceSizeNodeCount("M20", 1)),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("electable_specs"), specInstanceSizeNodeCount("M20", 3)),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("auto_scaling"), autoScalingEnabled),
					plancheck.ExpectKnownValue(resourceName, regionConfig1.AtMapKey("analytics_auto_scaling"), autoScalingEnabled),
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

func buildPlanWithAutoScaling(t *testing.T, useEffectiveFields, computeEnabled, diskGBEnabled bool) *advancedcluster.TFModel {
	t.Helper()
	ctx := context.Background()

	autoScaling, diags := types.ObjectValueFrom(ctx, autoScalingAttrTypes, advancedcluster.TFAutoScalingModel{
		ComputeEnabled:          types.BoolValue(computeEnabled),
		DiskGBEnabled:           types.BoolValue(diskGBEnabled),
		ComputeScaleDownEnabled: types.BoolValue(false),
		ComputeMinInstanceSize:  types.StringValue("M10"),
		ComputeMaxInstanceSize:  types.StringValue("M40"),
	})
	require.Empty(t, diags)

	regionConfig := advancedcluster.TFRegionConfigsModel{
		AutoScaling:          autoScaling,
		AnalyticsAutoScaling: types.ObjectNull(autoScalingAttrTypes),
		AnalyticsSpecs:       types.ObjectNull(specsAttrTypes),
		ElectableSpecs:       types.ObjectNull(specsAttrTypes),
		ReadOnlySpecs:        types.ObjectNull(specsAttrTypes),
		ProviderName:         types.StringValue("AWS"),
		RegionName:           types.StringValue("US_EAST_1"),
		Priority:             types.Int64Value(7),
		BackingProviderName:  types.StringNull(),
	}
	regionConfigs, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: regionConfigAttrTypes}, []advancedcluster.TFRegionConfigsModel{regionConfig})
	require.Empty(t, diags)

	repSpec := advancedcluster.TFReplicationSpecsModel{
		RegionConfigs: regionConfigs,
		ContainerId:   types.MapNull(types.StringType),
		ExternalId:    types.StringNull(),
		ZoneId:        types.StringNull(),
		ZoneName:      types.StringNull(),
	}
	repSpecs, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: replicationSpecAttrTypes}, []advancedcluster.TFReplicationSpecsModel{repSpec})
	require.Empty(t, diags)

	return &advancedcluster.TFModel{
		UseEffectiveFields: types.BoolValue(useEffectiveFields),
		ReplicationSpecs:   repSpecs,
		Labels:             types.MapNull(types.StringType),
		Tags:               types.MapNull(types.StringType),
	}
}

func TestAdvancedCluster_WarnIgnoredSpecChange(t *testing.T) {
	testCases := []struct {
		name               string
		changedFields      schemafunc.AttributeChanges
		useEffectiveFields bool
		computeEnabled     bool
		diskGBEnabled      bool
		expectWarning      bool
	}{
		{
			name:               "warns when compute auto-scaling on and instance_size changed",
			useEffectiveFields: true,
			computeEnabled:     true,
			changedFields:      schemafunc.AttributeChanges{"instance_size"},
			expectWarning:      true,
		},
		{
			name:               "warns when disk auto-scaling on and disk_size_gb changed",
			useEffectiveFields: true,
			diskGBEnabled:      true,
			changedFields:      schemafunc.AttributeChanges{"disk_size_gb"},
			expectWarning:      true,
		},
		{
			name:               "warns when compute auto-scaling on and disk_iops changed",
			useEffectiveFields: true,
			computeEnabled:     true,
			changedFields:      schemafunc.AttributeChanges{"disk_iops"},
			expectWarning:      true,
		},
		{
			name:               "no warning when use_effective_fields is false",
			useEffectiveFields: false,
			computeEnabled:     true,
			changedFields:      schemafunc.AttributeChanges{"instance_size"},
			expectWarning:      false,
		},
		{
			name:               "no warning when auto-scaling is disabled",
			useEffectiveFields: true,
			computeEnabled:     false,
			diskGBEnabled:      false,
			changedFields:      schemafunc.AttributeChanges{"instance_size"},
			expectWarning:      false,
		},
		{
			name:               "no warning when auto-scaling is on but no managed spec fields changed",
			useEffectiveFields: true,
			computeEnabled:     true,
			changedFields:      schemafunc.AttributeChanges{"node_count"},
			expectWarning:      false,
		},
		{
			name:               "no warning when replication_specs list length changes (avoid false positive on new spec)",
			useEffectiveFields: true,
			computeEnabled:     true,
			changedFields:      schemafunc.AttributeChanges{"replication_specs[+1]", "replication_specs[1].region_configs[0].instance_size", "instance_size"},
			expectWarning:      false,
		},
		{
			name:               "no warning when region_configs list length changes (avoid false positive on new region config)",
			useEffectiveFields: true,
			computeEnabled:     true,
			changedFields:      schemafunc.AttributeChanges{"replication_specs[0].region_configs[+1]", "replication_specs[0].region_configs[1].instance_size", "instance_size"},
			expectWarning:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			plan := buildPlanWithAutoScaling(t, tc.useEffectiveFields, tc.computeEnabled, tc.diskGBEnabled)
			var diags diag.Diagnostics

			advancedcluster.WarnIgnoredSpecChange(ctx, &diags, plan, tc.changedFields)

			assert.False(t, diags.HasError())
			if tc.expectWarning {
				assert.Equal(t, 1, diags.WarningsCount())
			} else {
				assert.Equal(t, 0, diags.WarningsCount())
			}
		})
	}
}
