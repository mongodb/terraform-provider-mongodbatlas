package schemafunc_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/stretchr/testify/assert"
)

type missingMethodType struct{}

type wrongReturnType struct{}

func (w wrongReturnType) IsUnknown() string {
	return "I'm a string!"
}

type multipleReturns struct{}

func (m multipleReturns) IsUnknown() (bool, error) {
	return false, nil
}

type hasUnknownsOk struct {
	Field types.Bool
}

type hasUnknownsPanicWrongType struct {
	Field wrongReturnType
}

type hasUnknownsPanicMultipleReturn struct {
	Field multipleReturns
}

type hasUnknownsPanicMissingMethod struct {
	Field missingMethodType
}

func TestHasUnknown(t *testing.T) {
	tests := map[string]struct {
		input        any
		panicMessage string
		inputBool    types.Bool
		expected     bool
	}{
		"valid unknown true": {
			inputBool: types.BoolUnknown(),
			expected:  true,
		},
		"valid unknown false": {
			inputBool: types.BoolValue(true),
			expected:  false,
		},
		"missing IsUnknown method": {
			input:        &hasUnknownsPanicMissingMethod{missingMethodType{}},
			panicMessage: "IsUnknown method not found for {}",
		},
		"wrong return type": {
			input:        &hasUnknownsPanicWrongType{wrongReturnType{}},
			panicMessage: "IsUnknown method must return a bool, got I'm a string!",
		},
		"multiple return values": {
			input:        &hasUnknownsPanicMultipleReturn{multipleReturns{}},
			panicMessage: "IsUnknown method must return a single value, got [<bool Value> <error Value>]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panicMessage != "" {
				assert.PanicsWithValue(t, tc.panicMessage, func() {
					schemafunc.HasUnknowns(tc.input)
				})
				return
			}

			wrapper := &hasUnknownsOk{Field: tc.inputBool}
			result := schemafunc.HasUnknowns(wrapper)
			assert.Equal(t, tc.expected, result)
		})
	}
}

type TFSpec struct {
	InstanceSize types.String `tfsdk:"instance_size"`
	NodeCount    types.Int64  `tfsdk:"node_count"`
}

var SpecObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"instance_size": types.StringType,
	"node_count":    types.Int64Type,
}}

type TFRegionConfig struct {
	ProviderName types.String `tfsdk:"provider_name"`
	RegionName   types.String `tfsdk:"region_name"`
	Spec         types.Object `tfsdk:"spec"`
}

var RegionConfigsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
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
		ProviderName: types.StringValue("aws"),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         asObjectValue(ctx, TFSpec{InstanceSize: types.StringValue("M10"), NodeCount: types.Int64Value(3)}, SpecObjType.AttrTypes),
	}
	regionConfigDest = TFRegionConfig{
		ProviderName: types.StringUnknown(),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         types.ObjectUnknown(SpecObjType.AttrTypes),
	}
	regionConfigNodeCountUnknown = TFRegionConfig{
		ProviderName: types.StringValue("aws"),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         asObjectValue(ctx, TFSpec{InstanceSize: types.StringValue("M10"), NodeCount: types.Int64Unknown()}, SpecObjType.AttrTypes),
	}
	regionConfigSpecUnknown = TFRegionConfig{
		ProviderName: types.StringValue("aws"),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         types.ObjectUnknown(SpecObjType.AttrTypes),
	}
	advancedConfigTrue = asObjectValue(ctx, TFAdvancedConfig{JavascriptEnabled: types.BoolValue(true)}, AdvancedConfigObjType.AttrTypes)
)

func TestCopyUnknowns(t *testing.T) {
	tests := map[string]struct {
		src          *TFSimpleModel
		dest         *TFSimpleModel
		expectedDest *TFSimpleModel
		panicMessage string
		keepUnknown  []string
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
					schemafunc.CopyUnknowns(ctx, tc.src, tc.dest, tc.keepUnknown)
				})
				return
			}
			schemafunc.CopyUnknowns(ctx, tc.src, tc.dest, tc.keepUnknown)
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
			panic("failed to create region config object")
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
