package schemafunc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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

// TestHasUnknown is kept even if HasUnknowns is not used so we can test isUnknown that is also used in CopyFromUnknown
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
	regionConfigProviderNameUnknown = TFRegionConfig{
		AutoScaling:  autoScalingFalseAndNull,
		ProviderName: types.StringUnknown(),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         asObjectValue(ctx, TFSpec{InstanceSize: types.StringValue("M10"), NodeCount: types.Int64Value(3)}, SpecObjType.AttrTypes),
	}
	regionConfigSpecUnknown = TFRegionConfig{
		AutoScaling:  autoScalingFalseAndNull,
		ProviderName: types.StringValue("aws"),
		RegionName:   types.StringValue("US_EAST_1"),
		Spec:         types.ObjectUnknown(SpecObjType.AttrTypes),
	}
	advancedConfigTrue      = asObjectValue(ctx, TFAdvancedConfig{JavascriptEnabled: types.BoolValue(true)}, AdvancedConfigObjType.AttrTypes)
	autoScalingFalseAndNull = asObjectValue(ctx, TFAutoScalingModel{
		ComputeEnabled:         types.BoolValue(false),
		DiskGBEnabled:          types.BoolValue(false),
		ComputeMinInstanceSize: types.StringNull(),
		ComputeMaxInstanceSize: types.StringNull(),
	}, AutoScalingObjType.AttrTypes)

	// project_id and provider_name are Optional without Computed at their level, the rest can be copied.
	simpleModelAttributes = map[string]schema.Attribute{
		"project_id":      schema.StringAttribute{Optional: true},
		"name":            schema.StringAttribute{Optional: true, Computed: true},
		"backup_enabled":  schema.BoolAttribute{Computed: true},
		"advanced_config": schema.SingleNestedAttribute{Computed: true, Attributes: map[string]schema.Attribute{"javascript_enabled": schema.BoolAttribute{Optional: true, Computed: true}}},
		"replication_specs": schema.ListNestedAttribute{Optional: true, Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"zone_name": schema.StringAttribute{Optional: true, Computed: true},
			"region_configs": schema.ListNestedAttribute{Optional: true, Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"provider_name": schema.StringAttribute{Optional: true},
				"region_name":   schema.StringAttribute{Optional: true, Computed: true},
				"auto_scaling":  schema.SingleNestedAttribute{Optional: true, Computed: true, Attributes: map[string]schema.Attribute{}},
				"spec": schema.SingleNestedAttribute{Optional: true, Computed: true, Attributes: map[string]schema.Attribute{
					"instance_size": schema.StringAttribute{Optional: true, Computed: true},
					"node_count":    schema.Int64Attribute{Optional: true, Computed: true},
				}},
			}}},
		}}},
	}
)

func TestCopyUnknowns(t *testing.T) {
	tests := map[string]struct {
		src          *TFSimpleModel
		dest         *TFSimpleModel
		expectedDest *TFSimpleModel
		attributes   map[string]schema.Attribute
		panicMessage string
		keepUnknown  []string
	}{
		"schema keeps Optional-only attributes unknown": {
			src: &TFSimpleModel{
				ProjectID:        types.StringValue("src-project"),
				Name:             types.StringValue("src-name"),
				BackupEnabled:    types.BoolValue(true),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigSrc}),
			},
			dest: &TFSimpleModel{
				ProjectID:        types.StringUnknown(),
				Name:             types.StringUnknown(),
				BackupEnabled:    types.BoolUnknown(),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringUnknown(), []TFRegionConfig{regionConfigDest}),
			},
			expectedDest: &TFSimpleModel{
				// project_id and the nested provider_name are Optional-only, the others are copied from src.
				ProjectID:        types.StringUnknown(),
				Name:             types.StringValue("src-name"),
				BackupEnabled:    types.BoolValue(true),
				ReplicationSpecs: newReplicationSpecs(ctx, types.StringValue("Zone 1"), []TFRegionConfig{regionConfigProviderNameUnknown}),
			},
			attributes: simpleModelAttributes,
		},
		"schema and keepUnknown are combined": {
			src: &TFSimpleModel{
				ProjectID:      types.StringValue("src-project"),
				Name:           types.StringValue("src-name"),
				BackupEnabled:  types.BoolValue(true),
				AdvancedConfig: advancedConfigTrue,
			},
			dest: &TFSimpleModel{
				ProjectID:      types.StringUnknown(),
				Name:           types.StringUnknown(),
				BackupEnabled:  types.BoolUnknown(),
				AdvancedConfig: types.ObjectUnknown(AdvancedConfigObjType.AttrTypes),
			},
			expectedDest: &TFSimpleModel{
				ProjectID:      types.StringUnknown(),
				Name:           types.StringUnknown(),
				BackupEnabled:  types.BoolValue(true),
				AdvancedConfig: advancedConfigTrue,
			},
			attributes:  simpleModelAttributes,
			keepUnknown: []string{"name"},
		},
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
					schemafunc.CopyUnknowns(ctx, tc.src, tc.dest, tc.attributes, tc.keepUnknown)
				})
				return
			}
			schemafunc.CopyUnknowns(ctx, tc.src, tc.dest, tc.attributes, tc.keepUnknown)
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

func TestUnknownInConfig(t *testing.T) {
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name": tftypes.String, "tags": tftypes.Map{ElementType: tftypes.String},
		"specs": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"size": tftypes.String}}},
	}}
	specs := func(size tftypes.Value) tftypes.Value {
		element := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"size": tftypes.String}}
		return tftypes.NewValue(tftypes.List{ElementType: element},
			[]tftypes.Value{tftypes.NewValue(element, map[string]tftypes.Value{"size": size})})
	}
	object := func(name, tags, specs tftypes.Value) tftypes.Value {
		return tftypes.NewValue(objType, map[string]tftypes.Value{"name": name, "tags": tags, "specs": specs})
	}
	knownName := tftypes.NewValue(tftypes.String, "example")
	nullTags := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil)
	knownSize := tftypes.NewValue(tftypes.String, "M10")
	unknownSize := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	testCases := map[string]struct {
		config   tftypes.Value
		expected []string
	}{
		"all known":                {object(knownName, nullTags, specs(knownSize)), nil},
		"unknown attribute":        {object(knownName, tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, tftypes.UnknownValue), specs(knownSize)), []string{"tags"}},
		"unknown nested in a list": {object(knownName, nullTags, specs(unknownSize)), []string{"specs"}},
		"sorted names":             {object(tftypes.NewValue(tftypes.String, tftypes.UnknownValue), nullTags, specs(unknownSize)), []string{"name", "specs"}},
		// Degenerate values must not panic, ModifyPlan is called for destroy plans too.
		"zero value":     {tftypes.Value{}, nil},
		"null object":    {tftypes.NewValue(objType, nil), nil},
		"unknown object": {tftypes.NewValue(objType, tftypes.UnknownValue), nil},
		"not an object":  {tftypes.NewValue(tftypes.String, "example"), nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, schemafunc.UnknownInConfig(tc.config))
		})
	}
}
