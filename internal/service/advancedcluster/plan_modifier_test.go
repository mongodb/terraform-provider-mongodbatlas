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

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

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

type regionConfigTestParams struct {
	electableInstanceSize   string // empty = null analytics_specs
	analyticsInstanceSize   string
	diskSizeGb              float64
	diskIops                int64
	computeEnabled          bool
	analyticsComputeEnabled bool
	diskGBEnabled           bool
}

func buildModelForWarnTest(t *testing.T, useEffectiveFields bool, rcParams regionConfigTestParams) *advancedcluster.TFModel {
	t.Helper()
	ctx := context.Background()

	autoScaling, diags := types.ObjectValueFrom(ctx, autoScalingAttrTypes, advancedcluster.TFAutoScalingModel{
		ComputeEnabled:          types.BoolValue(rcParams.computeEnabled),
		DiskGBEnabled:           types.BoolValue(rcParams.diskGBEnabled),
		ComputeScaleDownEnabled: types.BoolValue(false),
		ComputeMinInstanceSize:  types.StringValue("M10"),
		ComputeMaxInstanceSize:  types.StringValue("M40"),
	})
	require.Empty(t, diags)

	analyticsAutoScaling, diags := types.ObjectValueFrom(ctx, autoScalingAttrTypes, advancedcluster.TFAutoScalingModel{
		ComputeEnabled:          types.BoolValue(rcParams.analyticsComputeEnabled),
		DiskGBEnabled:           types.BoolValue(false),
		ComputeScaleDownEnabled: types.BoolValue(false),
		ComputeMinInstanceSize:  types.StringValue("M10"),
		ComputeMaxInstanceSize:  types.StringValue("M40"),
	})
	require.Empty(t, diags)

	electableSpecs, diags := types.ObjectValueFrom(ctx, specsAttrTypes, advancedcluster.TFSpecsModel{
		InstanceSize:  types.StringValue(rcParams.electableInstanceSize),
		NodeCount:     types.Int64Value(3),
		DiskSizeGb:    types.Float64Value(rcParams.diskSizeGb),
		DiskIops:      types.Int64Value(rcParams.diskIops),
		EbsVolumeType: types.StringNull(),
	})
	require.Empty(t, diags)

	var analyticsSpecs types.Object
	if rcParams.analyticsInstanceSize != "" {
		analyticsSpecs, diags = types.ObjectValueFrom(ctx, specsAttrTypes, advancedcluster.TFSpecsModel{
			InstanceSize:  types.StringValue(rcParams.analyticsInstanceSize),
			NodeCount:     types.Int64Value(1),
			DiskSizeGb:    types.Float64Value(rcParams.diskSizeGb),
			DiskIops:      types.Int64Value(rcParams.diskIops),
			EbsVolumeType: types.StringNull(),
		})
		require.Empty(t, diags)
	} else {
		analyticsSpecs = types.ObjectNull(specsAttrTypes)
	}

	regionConfig := advancedcluster.TFRegionConfigsModel{
		AutoScaling:          autoScaling,
		AnalyticsAutoScaling: analyticsAutoScaling,
		AnalyticsSpecs:       analyticsSpecs,
		ElectableSpecs:       electableSpecs,
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
		stateRC            regionConfigTestParams
		planRC             regionConfigTestParams
		useEffectiveFields bool
		expectWarning      bool
	}{
		{
			name:               "warns when compute auto-scaling on and electable instance_size changed",
			useEffectiveFields: true,
			stateRC:            regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10"},
			planRC:             regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M20"},
			expectWarning:      true,
		},
		{
			name:               "warns when disk auto-scaling on and disk_size_gb changed",
			useEffectiveFields: true,
			stateRC:            regionConfigTestParams{diskGBEnabled: true, electableInstanceSize: "M10", diskSizeGb: 10},
			planRC:             regionConfigTestParams{diskGBEnabled: true, electableInstanceSize: "M10", diskSizeGb: 20},
			expectWarning:      true,
		},
		{
			name:               "warns when disk auto-scaling on and disk_iops changed",
			useEffectiveFields: true,
			stateRC:            regionConfigTestParams{diskGBEnabled: true, electableInstanceSize: "M10", diskIops: 3000},
			planRC:             regionConfigTestParams{diskGBEnabled: true, electableInstanceSize: "M10", diskIops: 4000},
			expectWarning:      true,
		},
		{
			name:               "warns when analytics compute auto-scaling on and analytics instance_size changed",
			useEffectiveFields: true,
			stateRC:            regionConfigTestParams{analyticsComputeEnabled: true, electableInstanceSize: "M10", analyticsInstanceSize: "M10"},
			planRC:             regionConfigTestParams{analyticsComputeEnabled: true, electableInstanceSize: "M10", analyticsInstanceSize: "M20"},
			expectWarning:      true,
		},
		{
			name:               "no warning when use_effective_fields is false",
			useEffectiveFields: false,
			stateRC:            regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10"},
			planRC:             regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M20"},
			expectWarning:      false,
		},
		{
			name:               "no warning when auto-scaling is disabled",
			useEffectiveFields: true,
			stateRC:            regionConfigTestParams{electableInstanceSize: "M10"},
			planRC:             regionConfigTestParams{electableInstanceSize: "M20"},
			expectWarning:      false,
		},
		{
			name:               "no warning when auto-scaling is on but no managed spec fields changed",
			useEffectiveFields: true,
			stateRC:            regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10"},
			planRC:             regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10"},
			expectWarning:      false,
		},
		{
			name:               "no warning when only analytics compute auto-scaling on but electable instance_size changed",
			useEffectiveFields: true,
			stateRC:            regionConfigTestParams{analyticsComputeEnabled: true, electableInstanceSize: "M10"},
			planRC:             regionConfigTestParams{analyticsComputeEnabled: true, electableInstanceSize: "M20"},
			expectWarning:      false,
		},
		{
			name:               "no warning when only compute auto-scaling on but analytics instance_size changed",
			useEffectiveFields: true,
			stateRC:            regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10", analyticsInstanceSize: "M10"},
			planRC:             regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10", analyticsInstanceSize: "M20"},
			expectWarning:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			state := buildModelForWarnTest(t, tc.useEffectiveFields, tc.stateRC)
			plan := buildModelForWarnTest(t, tc.useEffectiveFields, tc.planRC)
			var diags diag.Diagnostics

			advancedcluster.WarnIgnoredSpecChange(ctx, &diags, state, plan)

			assert.False(t, diags.HasError())
			if tc.expectWarning {
				assert.Equal(t, 1, diags.WarningsCount())
			} else {
				assert.Equal(t, 0, diags.WarningsCount())
			}
		})
	}

	// List length changes: new specs/region_configs added in plan have no state counterpart,
	// so minLen iteration naturally skips them — no false positive warning.
	t.Run("no warning when replication_specs list length changes", func(t *testing.T) {
		ctx := context.Background()
		state := buildModelForWarnTest(t, true, regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10"})
		// plan has 2 rep specs; first is unchanged, second is new (no state counterpart)
		plan := addExtraRepSpec(t, buildModelForWarnTest(t, true, regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10"}),
			regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M20"})
		var diags diag.Diagnostics
		advancedcluster.WarnIgnoredSpecChange(ctx, &diags, state, plan)
		assert.False(t, diags.HasError())
		assert.Equal(t, 0, diags.WarningsCount())
	})

	t.Run("no warning when region_configs list length changes", func(t *testing.T) {
		ctx := context.Background()
		state := buildModelForWarnTest(t, true, regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10"})
		// plan has 2 region configs; first is unchanged, second is new (no state counterpart)
		plan := addExtraRegionConfig(t, buildModelForWarnTest(t, true, regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M10"}),
			regionConfigTestParams{computeEnabled: true, electableInstanceSize: "M20"})
		var diags diag.Diagnostics
		advancedcluster.WarnIgnoredSpecChange(ctx, &diags, state, plan)
		assert.False(t, diags.HasError())
		assert.Equal(t, 0, diags.WarningsCount())
	})
}

// addExtraRepSpec appends a new replication spec to the model's ReplicationSpecs list.
func addExtraRepSpec(t *testing.T, model *advancedcluster.TFModel, rcParams regionConfigTestParams) *advancedcluster.TFModel {
	t.Helper()
	extra := buildModelForWarnTest(t, model.UseEffectiveFields.ValueBool(), rcParams)
	ctx := context.Background()

	existing := make([]advancedcluster.TFReplicationSpecsModel, 0)
	diags := model.ReplicationSpecs.ElementsAs(ctx, &existing, false)
	require.Empty(t, diags)

	extra1 := make([]advancedcluster.TFReplicationSpecsModel, 0)
	diags = extra.ReplicationSpecs.ElementsAs(ctx, &extra1, false)
	require.Empty(t, diags)

	combined, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: replicationSpecAttrTypes}, append(existing, extra1...))
	require.Empty(t, diags)

	model.ReplicationSpecs = combined
	return model
}

// addExtraRegionConfig appends a new region config to the first replication spec of the model.
func addExtraRegionConfig(t *testing.T, model *advancedcluster.TFModel, rcParams regionConfigTestParams) *advancedcluster.TFModel {
	t.Helper()
	extra := buildModelForWarnTest(t, model.UseEffectiveFields.ValueBool(), rcParams)
	ctx := context.Background()

	repSpecs := make([]advancedcluster.TFReplicationSpecsModel, 0)
	diags := model.ReplicationSpecs.ElementsAs(ctx, &repSpecs, false)
	require.Empty(t, diags)

	existingRCs := make([]advancedcluster.TFRegionConfigsModel, 0)
	diags = repSpecs[0].RegionConfigs.ElementsAs(ctx, &existingRCs, false)
	require.Empty(t, diags)

	extraRepSpecs := make([]advancedcluster.TFReplicationSpecsModel, 0)
	diags = extra.ReplicationSpecs.ElementsAs(ctx, &extraRepSpecs, false)
	require.Empty(t, diags)

	extraRCs := make([]advancedcluster.TFRegionConfigsModel, 0)
	diags = extraRepSpecs[0].RegionConfigs.ElementsAs(ctx, &extraRCs, false)
	require.Empty(t, diags)

	combined, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: regionConfigAttrTypes}, append(existingRCs, extraRCs...))
	require.Empty(t, diags)

	repSpecs[0].RegionConfigs = combined
	repSpecsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: replicationSpecAttrTypes}, repSpecs)
	require.Empty(t, diags)

	model.ReplicationSpecs = repSpecsList
	return model
}
