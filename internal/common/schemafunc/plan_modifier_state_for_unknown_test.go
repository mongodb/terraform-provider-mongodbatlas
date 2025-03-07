package schemafunc_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/stretchr/testify/assert"
)

type TFAutoScalingModel struct {
	ComputeMinInstanceSize types.String `tfsdk:"compute_min_instance_size"`
	ComputeMaxInstanceSize types.String `tfsdk:"compute_max_instance_size"`
	ComputeEnabled         types.Bool   `tfsdk:"compute_enabled"`
	DiskGBEnabled          types.Bool   `tfsdk:"disk_gb_enabled"`
}

var AutoScalingObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"compute_enabled":           types.BoolType,
	"disk_gb_enabled":           types.BoolType,
	"compute_min_instance_size": types.StringType,
	"compute_max_instance_size": types.StringType,
}}

type TFSpec struct {
	InstanceSize types.String `tfsdk:"instance_size"`
	NodeCount    types.Int64  `tfsdk:"node_count"`
}

var SpecObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"instance_size": types.StringType,
	"node_count":    types.Int64Type,
}}

type TFRegionConfig struct {
	AutoScaling  types.Object `tfsdk:"auto_scaling"`
	ProviderName types.String `tfsdk:"provider_name"`
	RegionName   types.String `tfsdk:"region_name"`
	Spec         types.Object `tfsdk:"spec"`
}

var RegionConfigsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"auto_scaling":  AutoScalingObjType,
	"provider_name": types.StringType,
	"region_name":   types.StringType,
	"spec":          SpecObjType,
}}

type TFReplicationSpec struct {
	RegionConfigs types.List   `tfsdk:"region_configs"`
	ZoneName      types.String `tfsdk:"zone_name"`
}

var ReplicationSpecsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"zone_name":      types.StringType,
	"region_configs": types.ListType{ElemType: RegionConfigsObjType},
}}

type TFAdvancedConfig struct {
	JavascriptEnabled types.Bool `tfsdk:"javascript_enabled"`
}

var AdvancedConfigObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"javascript_enabled": types.BoolType,
}}

type TFSimpleModel struct {
	ReplicationSpecs types.List   `tfsdk:"replication_specs"`
	ProjectID        types.String `tfsdk:"project_id"`
	Name             types.String `tfsdk:"name"`
	AdvancedConfig   types.Object `tfsdk:"advanced_config"`
	BackupEnabled    types.Bool   `tfsdk:"backup_enabled"`
}

var (
	ctx             = context.Background()
	regionConfigSrc = TFRegionConfig{
		AutoScaling:  autoScalingFalseAndNull,
		ProviderName: types.StringValue("aws"),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         asObjectValue(ctx, TFSpec{InstanceSize: types.StringValue("M10"), NodeCount: types.Int64Value(3)}, SpecObjType.AttrTypes),
	}
	regionConfigNodeCount0 = TFRegionConfig{
		AutoScaling:  autoScalingFalseAndNull,
		ProviderName: types.StringValue("aws"),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         asObjectValue(ctx, TFSpec{InstanceSize: types.StringValue("M10"), NodeCount: types.Int64Value(0)}, SpecObjType.AttrTypes),
	}
	regionConfigDest = TFRegionConfig{
		AutoScaling:  autoScalingFalseAndNull,
		ProviderName: types.StringUnknown(),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         types.ObjectUnknown(SpecObjType.AttrTypes),
	}
	regionConfigNodeCountUnknown = TFRegionConfig{
		AutoScaling:  autoScalingFalseAndNull,
		ProviderName: types.StringValue("aws"),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         asObjectValue(ctx, TFSpec{InstanceSize: types.StringValue("M10"), NodeCount: types.Int64Unknown()}, SpecObjType.AttrTypes),
	}
	regionConfigSpecUnknown = TFRegionConfig{
		AutoScaling:  autoScalingFalseAndNull,
		ProviderName: types.StringValue("aws"),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         types.ObjectUnknown(SpecObjType.AttrTypes),
	}
	advancedConfigTrue       = asObjectValue(ctx, TFAdvancedConfig{JavascriptEnabled: types.BoolValue(true)}, AdvancedConfigObjType.AttrTypes)
	autoScalingFalseAndEmpty = asObjectValue(ctx, TFAutoScalingModel{
		ComputeEnabled:         types.BoolValue(false),
		DiskGBEnabled:          types.BoolValue(false),
		ComputeMinInstanceSize: types.StringValue(""),
		ComputeMaxInstanceSize: types.StringValue(""),
	}, AutoScalingObjType.AttrTypes)

	autoScalingFalseAndNull = asObjectValue(ctx, TFAutoScalingModel{
		ComputeEnabled:         types.BoolValue(false),
		DiskGBEnabled:          types.BoolValue(false),
		ComputeMinInstanceSize: types.StringNull(),
		ComputeMaxInstanceSize: types.StringNull(),
	}, AutoScalingObjType.AttrTypes)
	autoScalingTrueAndNonEmpty = asObjectValue(ctx, TFAutoScalingModel{
		ComputeEnabled:         types.BoolValue(true),
		DiskGBEnabled:          types.BoolValue(true),
		ComputeMinInstanceSize: types.StringValue("M10"),
		ComputeMaxInstanceSize: types.StringValue("M20"),
	}, AutoScalingObjType.AttrTypes)
	autoScalingUnknown      = types.ObjectUnknown(AutoScalingObjType.AttrTypes)
	autoScalingLeafsUnknown = asObjectValue(ctx, TFAutoScalingModel{
		ComputeEnabled:         types.BoolUnknown(),
		DiskGBEnabled:          types.BoolUnknown(),
		ComputeMinInstanceSize: types.StringUnknown(),
		ComputeMaxInstanceSize: types.StringUnknown(),
	}, AutoScalingObjType.AttrTypes)
	keepProjectIDUnknown = func(name string, value attr.Value) bool {
		return name == "project_id"
	}
	useStateOnlyWhenNodeCount0 = func(name string, value attr.Value) bool {
		return name == "node_count" && !value.Equal(types.Int64Value(0))
	}

	autoScalingBoolsKeepUnknown = func(name string, value attr.Value) bool {
		return slices.Contains([]string{"compute_enabled", "disk_gb_enabled"}, name) && value.Equal(types.BoolValue(true))
	}
	autoScalingStringsKeepUnknown = func(name string, value attr.Value) bool {
		return slices.Contains([]string{"compute_min_instance_size", "compute_max_instance_size"}, name) && value.(types.String).ValueString() != ""
	}
)

func TestCopyUnknowns(t *testing.T) {
	tests := map[string]struct {
		src              *TFSimpleModel
		dest             *TFSimpleModel
		expectedDest     *TFSimpleModel
		keepUnknownCalls schemafunc.KeepUnknownFunc
		panicMessage     string
		keepUnknown      []string
	}{
		"copy unknown basic fields": {
			src: &TFSimpleModel{
				ProjectID:     types.StringValue("src-project"),
				Name:          types.StringValue("src-name"),
				BackupEnabled: types.BoolValue(true),
			},
			dest: &TFSimpleModel{
				ProjectID:     types.StringUnknown(),
				Name:          types.StringValue("dest-name"),
				BackupEnabled: types.BoolUnknown(),
			},
			expectedDest: &TFSimpleModel{
				ProjectID:     types.StringValue("src-project"),
				Name:          types.StringValue("dest-name"),
				BackupEnabled: types.BoolValue(true),
			},
		},
		"respect keepUnknown": {
			src: &TFSimpleModel{
				ProjectID:        types.StringValue("src-project"),
				Name:             types.StringValue("src-name"),
				BackupEnabled:    types.BoolValue(true),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigSrc, regionConfigSrc}),
				AdvancedConfig:   advancedConfigTrue,
			},
			dest: &TFSimpleModel{
				ProjectID:        types.StringUnknown(),
				Name:             types.StringUnknown(),
				BackupEnabled:    types.BoolUnknown(),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringUnknown(), []TFRegionConfig{regionConfigNodeCountUnknown, regionConfigSpecUnknown}),
				AdvancedConfig:   types.ObjectUnknown(AdvancedConfigObjType.AttrTypes),
			},
			expectedDest: &TFSimpleModel{
				ProjectID:        types.StringValue("src-project"),
				Name:             types.StringUnknown(),
				BackupEnabled:    types.BoolValue(true),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringUnknown(), []TFRegionConfig{regionConfigNodeCountUnknown, regionConfigNodeCountUnknown}),
				AdvancedConfig:   types.ObjectUnknown(AdvancedConfigObjType.AttrTypes),
			},
			keepUnknown: []string{"name", "advanced_config", "zone_name", "node_count"},
		},
		"respect keepUnknown on object": {
			src: &TFSimpleModel{
				ProjectID:        types.StringValue("src-project"),
				Name:             types.StringValue("src-name"),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ProjectID:        types.StringUnknown(),
				Name:             types.StringUnknown(),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringUnknown(), []TFRegionConfig{regionConfigSpecUnknown}),
			},
			expectedDest: &TFSimpleModel{
				ProjectID:        types.StringValue("src-project"),
				Name:             types.StringValue("src-name"),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigSpecUnknown}),
			},
			keepUnknown: []string{"spec"},
		},
		"respect keepUnknownCall root": {
			src: &TFSimpleModel{
				ProjectID: types.StringValue("src-project"),
				Name:      types.StringValue("src-name"),
			},
			dest: &TFSimpleModel{
				ProjectID: types.StringUnknown(),
				Name:      types.StringUnknown(),
			},
			expectedDest: &TFSimpleModel{
				ProjectID: types.StringUnknown(),
				Name:      types.StringValue("src-name"),
			},
			keepUnknownCalls: keepProjectIDUnknown,
		},
		"respect keepUnknownCall nested": {
			src: &TFSimpleModel{
				ProjectID:        types.StringValue("src-project"),
				Name:             types.StringValue("src-name"),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ProjectID:        types.StringUnknown(),
				Name:             types.StringUnknown(),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringUnknown(), []TFRegionConfig{regionConfigNodeCountUnknown}),
			},
			expectedDest: &TFSimpleModel{
				ProjectID:        types.StringValue("src-project"),
				Name:             types.StringValue("src-name"),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigNodeCountUnknown}),
			},
			keepUnknownCalls: useStateOnlyWhenNodeCount0,
		},
		"respect multiple keepUnknownCall": {
			src: &TFSimpleModel{
				ProjectID:        types.StringValue("src-project"),
				Name:             types.StringValue("src-name"),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ProjectID:        types.StringUnknown(),
				Name:             types.StringUnknown(),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringUnknown(), []TFRegionConfig{regionConfigNodeCountUnknown}),
			},
			expectedDest: &TFSimpleModel{
				ProjectID:        types.StringUnknown(),
				Name:             types.StringValue("src-name"),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigNodeCountUnknown}),
			},
			keepUnknownCalls: schemafunc.KeepUnknownFuncOr(keepProjectIDUnknown, useStateOnlyWhenNodeCount0),
		},
		"copy node_count 0": {
			src: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigNodeCount0}),
			},
			dest: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringUnknown(), []TFRegionConfig{regionConfigNodeCountUnknown}),
			},
			expectedDest: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigNodeCount0}),
			},
			keepUnknownCalls: useStateOnlyWhenNodeCount0,
		},
		"keepUnknownCall on string": {
			src: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{
					regionConfigWithAutoScaling(autoScalingFalseAndEmpty),
					regionConfigWithAutoScaling(autoScalingFalseAndNull),
					regionConfigWithAutoScaling(autoScalingTrueAndNonEmpty),
				}),
			},
			dest: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringUnknown(), []TFRegionConfig{
					regionConfigWithAutoScaling(autoScalingUnknown),
					regionConfigWithAutoScaling(autoScalingUnknown),
					regionConfigWithAutoScaling(autoScalingUnknown),
				}),
			},
			expectedDest: &TFSimpleModel{
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{
					regionConfigWithAutoScaling(autoScalingFalseAndEmpty),
					regionConfigWithAutoScaling(autoScalingFalseAndNull),
					regionConfigWithAutoScaling(autoScalingLeafsUnknown),
				}),
			},
			keepUnknownCalls: schemafunc.KeepUnknownFuncOr(autoScalingStringsKeepUnknown, autoScalingBoolsKeepUnknown),
		},
		"non-pointer input": {
			src:          &TFSimpleModel{},
			dest:         nil,
			panicMessage: "params must be pointers to structs: *schemafunc_test.TFSimpleModel, *schemafunc_test.TFSimpleModel and not nil: (&{<null> <null> <null> <null> <null>}, <nil>)\n",
		},
		"unknown nested field at root": {
			src: &TFSimpleModel{
				AdvancedConfig: advancedConfigTrue,
			},
			dest: &TFSimpleModel{
				ProjectID:      types.StringValue("project"),
				AdvancedConfig: types.ObjectUnknown(AdvancedConfigObjType.AttrTypes),
			},
			expectedDest: &TFSimpleModel{
				ProjectID:      types.StringValue("project"),
				AdvancedConfig: advancedConfigTrue,
			},
		},
		"nested unknown fields object": {
			src: &TFSimpleModel{
				ProjectID:      types.StringValue("src-project"),
				AdvancedConfig: advancedConfigTrue,
			},
			dest: &TFSimpleModel{
				ProjectID: types.StringUnknown(),
				AdvancedConfig: asObjectValue(ctx, TFAdvancedConfig{
					JavascriptEnabled: types.BoolUnknown(),
				}, AdvancedConfigObjType.AttrTypes),
			},
			expectedDest: &TFSimpleModel{
				ProjectID:      types.StringValue("src-project"),
				AdvancedConfig: advancedConfigTrue,
			},
		},
		"nested unknown fields list": {
			src: &TFSimpleModel{
				ProjectID: types.StringValue("src-project"),
				ReplicationSpecs: newReplicationSpecs(ctx,
					types.StringValue("zone1"),
					[]TFRegionConfig{regionConfigSrc},
				),
			},
			dest: &TFSimpleModel{
				ProjectID: types.StringUnknown(),
				ReplicationSpecs: newReplicationSpecs(ctx,
					types.StringUnknown(),
					[]TFRegionConfig{regionConfigDest},
				),
			},
			expectedDest: &TFSimpleModel{
				ProjectID: types.StringValue("src-project"),
				ReplicationSpecs: newReplicationSpecs(ctx,
					types.StringValue("zone1"),
					[]TFRegionConfig{regionConfigSrc}),
			},
		},
		"nested unknown field in spec (list.list.object)": {
			src: &TFSimpleModel{
				ProjectID:        types.StringValue("src-project"),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone1"), []TFRegionConfig{regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ProjectID:        types.StringValue("dest-project"),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone2"), []TFRegionConfig{regionConfigNodeCountUnknown}),
			},
			expectedDest: &TFSimpleModel{
				ProjectID:        types.StringValue("dest-project"),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("zone2"), []TFRegionConfig{regionConfigSrc}),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panicMessage != "" {
				assert.PanicsWithValue(t, tc.panicMessage, func() {
					schemafunc.CopyUnknowns(ctx, tc.src, tc.dest, tc.keepUnknown, nil)
				})
				return
			}
			schemafunc.CopyUnknowns(ctx, tc.src, tc.dest, tc.keepUnknown, tc.keepUnknownCalls)
			assert.Equal(t, *tc.expectedDest, *tc.dest)
		})
	}
}

func asObjectValue[T any](ctx context.Context, t T, attrs map[string]attr.Type) types.Object {
	objType, diagsLocal := types.ObjectValueFrom(ctx, attrs, t)
	if diagsLocal.HasError() {
		panic("failed to convert object to model")
	}
	return objType
}

func newReplicationSpecs(ctx context.Context, zoneName types.String, regionConfigs []TFRegionConfig) types.List {
	regionConfigsObjects := make([]attr.Value, len(regionConfigs))
	for i, config := range regionConfigs {
		configObject, diags := types.ObjectValueFrom(ctx, RegionConfigsObjType.AttrTypes, config)
		if diags.HasError() {
			panic(fmt.Sprintf("failed to create region config object %v", diags))
		}
		regionConfigsObjects[i] = configObject
	}

	replicationSpec, diags := types.ObjectValueFrom(ctx, ReplicationSpecsObjType.AttrTypes, TFReplicationSpec{
		ZoneName:      zoneName,
		RegionConfigs: types.ListValueMust(RegionConfigsObjType, regionConfigsObjects),
	})
	if diags.HasError() {
		panic("failed to create replication spec object")
	}
	return types.ListValueMust(ReplicationSpecsObjType, []attr.Value{replicationSpec})
}

func combineReplicationSpecs(specs ...types.List) types.List {
	combined := []attr.Value{}
	for _, spec := range specs {
		combined = append(combined, spec.Elements()...)
	}
	return types.ListValueMust(ReplicationSpecsObjType, combined)
}

func regionConfigWithAutoScaling(autoScaling types.Object) TFRegionConfig {
	return TFRegionConfig{
		AutoScaling:  autoScaling,
		ProviderName: types.StringValue("aws"),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         asObjectValue(ctx, TFSpec{InstanceSize: types.StringValue("M10"), NodeCount: types.Int64Value(3)}, SpecObjType.AttrTypes),
	}
}
