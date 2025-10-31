package autogen_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveUnknowns(t *testing.T) {
	ctx := context.Background()

	type modelEmptyTest struct{}

	type modelCustomTypeTest struct {
		AttrKnownString   types.String                            `tfsdk:"attr_known_string"`
		AttrUnknownObject customtypes.ObjectValue[modelEmptyTest] `tfsdk:"attr_unknown_object"`
		AttrMANYUpper     types.Int64                             `tfsdk:"attr_many_upper"`
	}

	type modelst struct {
		AttrStringUnknown           types.String                                     `tfsdk:"attr_string_unknown"`
		AttrObjectUnknown           types.Object                                     `tfsdk:"attr_object_unknown"`
		AttrListUnknown             types.List                                       `tfsdk:"attr_list_unknown"`
		AttrObject                  types.Object                                     `tfsdk:"attr_object"`
		AttrListString              types.List                                       `tfsdk:"attr_list_string"`
		AttrSetString               types.Set                                        `tfsdk:"attr_set_string"`
		AttrListObjObj              types.List                                       `tfsdk:"attr_list_obj_obj"`
		AttrMapUnknown              types.Map                                        `tfsdk:"attr_map_unknown"`
		AttrCustomObjectUnknown     customtypes.ObjectValue[modelEmptyTest]          `tfsdk:"attr_custom_object_unknown"`
		AttrCustomObject            customtypes.ObjectValue[modelCustomTypeTest]     `tfsdk:"attr_custom_object"`
		AttrCustomListUnknown       customtypes.ListValue[types.String]              `tfsdk:"attr_custom_list_string"`
		AttrCustomSetUnknown        customtypes.SetValue[types.String]               `tfsdk:"attr_custom_set_string"`
		AttrCustomNestedListUnknown customtypes.NestedListValue[modelEmptyTest]      `tfsdk:"attr_custom_nested_list_unknown"`
		AttrCustomNestedSetUnknown  customtypes.NestedSetValue[modelEmptyTest]       `tfsdk:"attr_custom_nested_set_unknown"`
		AttrCustomNestedMapUnknown  customtypes.NestedMapValue[modelEmptyTest]       `tfsdk:"attr_custom_nested_map_unknown"`
		AttrCustomNestedList        customtypes.NestedListValue[modelCustomTypeTest] `tfsdk:"attr_custom_nested_list"`
		AttrCustomNestedSet         customtypes.NestedSetValue[modelCustomTypeTest]  `tfsdk:"attr_custom_nested_set"`
		AttrCustomNestedMap         customtypes.NestedMapValue[modelCustomTypeTest]  `tfsdk:"attr_custom_nested_map"`
	}

	model := modelst{
		AttrStringUnknown: types.StringUnknown(),
		AttrObjectUnknown: types.ObjectUnknown(objTypeTest.AttributeTypes()),
		AttrListUnknown:   types.ListUnknown(objTypeTest),
		AttrObject: types.ObjectValueMust(objTypeTest.AttributeTypes(), map[string]attr.Value{
			"attr_string": types.StringUnknown(),
			"attr_float":  types.Float64Unknown(),
			"attr_int":    types.Int64Unknown(),
			"attr_bool":   types.BoolUnknown(),
		}),
		AttrListString: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("val1"),
			types.StringUnknown(),
			types.StringValue("val2"),
			types.StringNull(),
		}),
		AttrSetString: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("se1"),
			types.StringUnknown(),
		}),
		AttrListObjObj: types.ListValueMust(objTypeParentTest, []attr.Value{
			types.ObjectValueMust(objTypeParentTest.AttributeTypes(), map[string]attr.Value{
				"attr_parent_obj": types.ObjectValueMust(objTypeTest.AttributeTypes(), map[string]attr.Value{
					"attr_string": types.StringUnknown(),
					"attr_float":  types.Float64Value(1.234),
					"attr_int":    types.Int64Value(1),
					"attr_bool":   types.BoolUnknown(),
				}),
				"attr_parent_string": types.StringUnknown(),
				"attr_parent_int":    types.Int64Value(1),
			}),
			types.ObjectValueMust(objTypeParentTest.AttributeTypes(), map[string]attr.Value{
				"attr_parent_obj": types.ObjectValueMust(objTypeTest.AttributeTypes(), map[string]attr.Value{
					"attr_string": types.StringValue("val1"),
					"attr_float":  types.Float64Value(1.234),
					"attr_int":    types.Int64Value(1),
					"attr_bool":   types.BoolValue(true),
				}),
				"attr_parent_string": types.StringUnknown(),
				"attr_parent_int":    types.Int64Unknown(),
			}),
			types.ObjectUnknown(objTypeParentTest.AttributeTypes()),
		}),
		AttrMapUnknown:          types.MapUnknown(types.StringType),
		AttrCustomObjectUnknown: customtypes.NewObjectValueUnknown[modelEmptyTest](ctx),
		AttrCustomObject: customtypes.NewObjectValue[modelCustomTypeTest](ctx, modelCustomTypeTest{
			AttrKnownString:   types.StringValue("val1"),
			AttrUnknownObject: customtypes.NewObjectValueUnknown[modelEmptyTest](ctx),
			AttrMANYUpper:     types.Int64Unknown(),
		}),
		AttrCustomNestedListUnknown: customtypes.NewNestedListValueUnknown[modelEmptyTest](ctx),
		AttrCustomNestedSetUnknown:  customtypes.NewNestedSetValueUnknown[modelEmptyTest](ctx),
		AttrCustomNestedMapUnknown:  customtypes.NewNestedMapValueUnknown[modelEmptyTest](ctx),
		AttrCustomListUnknown:       customtypes.NewListValueUnknown[types.String](ctx),
		AttrCustomSetUnknown:        customtypes.NewSetValueUnknown[types.String](ctx),
		AttrCustomNestedList: customtypes.NewNestedListValue[modelCustomTypeTest](ctx, []modelCustomTypeTest{
			{
				AttrKnownString:   types.StringValue("val1"),
				AttrUnknownObject: customtypes.NewObjectValueUnknown[modelEmptyTest](ctx),
				AttrMANYUpper:     types.Int64Unknown(),
			},
			{
				AttrKnownString:   types.StringUnknown(),
				AttrUnknownObject: customtypes.NewObjectValue[modelEmptyTest](ctx, modelEmptyTest{}),
				AttrMANYUpper:     types.Int64Value(2),
			},
		}),
		AttrCustomNestedMap: customtypes.NewNestedMapValue[modelCustomTypeTest](ctx, map[string]modelCustomTypeTest{
			"keyOne": {
				AttrKnownString:   types.StringValue("val1"),
				AttrUnknownObject: customtypes.NewObjectValueUnknown[modelEmptyTest](ctx),
				AttrMANYUpper:     types.Int64Unknown(),
			},
			"keyTwo": {
				AttrKnownString:   types.StringUnknown(),
				AttrUnknownObject: customtypes.NewObjectValue[modelEmptyTest](ctx, modelEmptyTest{}),
				AttrMANYUpper:     types.Int64Value(2),
			},
		}),
	}
	modelExpected := modelst{
		AttrStringUnknown: types.StringNull(),
		AttrObjectUnknown: types.ObjectNull(objTypeTest.AttributeTypes()),
		AttrListUnknown:   types.ListNull(objTypeTest),
		AttrObject: types.ObjectValueMust(objTypeTest.AttributeTypes(), map[string]attr.Value{
			"attr_string": types.StringNull(),
			"attr_float":  types.Float64Null(),
			"attr_int":    types.Int64Null(),
			"attr_bool":   types.BoolNull(),
		}),
		AttrListString: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("val1"),
			types.StringNull(),
			types.StringValue("val2"),
			types.StringNull(),
		}),
		AttrSetString: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("se1"),
			types.StringNull(),
		}),
		AttrListObjObj: types.ListValueMust(objTypeParentTest, []attr.Value{
			types.ObjectValueMust(objTypeParentTest.AttributeTypes(), map[string]attr.Value{
				"attr_parent_obj": types.ObjectValueMust(objTypeTest.AttributeTypes(), map[string]attr.Value{
					"attr_string": types.StringNull(),
					"attr_float":  types.Float64Value(1.234),
					"attr_int":    types.Int64Value(1),
					"attr_bool":   types.BoolNull(),
				}),
				"attr_parent_string": types.StringNull(),
				"attr_parent_int":    types.Int64Value(1),
			}),
			types.ObjectValueMust(objTypeParentTest.AttributeTypes(), map[string]attr.Value{
				"attr_parent_obj": types.ObjectValueMust(objTypeTest.AttributeTypes(), map[string]attr.Value{
					"attr_string": types.StringValue("val1"),
					"attr_float":  types.Float64Value(1.234),
					"attr_int":    types.Int64Value(1),
					"attr_bool":   types.BoolValue(true),
				}),
				"attr_parent_string": types.StringNull(),
				"attr_parent_int":    types.Int64Null(),
			}),
			types.ObjectNull(objTypeParentTest.AttributeTypes()),
		}),
		AttrMapUnknown:          types.MapNull(types.StringType),
		AttrCustomObjectUnknown: customtypes.NewObjectValueNull[modelEmptyTest](ctx),
		AttrCustomObject: customtypes.NewObjectValue[modelCustomTypeTest](ctx, modelCustomTypeTest{
			AttrKnownString:   types.StringValue("val1"),
			AttrUnknownObject: customtypes.NewObjectValueNull[modelEmptyTest](ctx),
			AttrMANYUpper:     types.Int64Null(),
		}),
		AttrCustomListUnknown:       customtypes.NewListValueNull[types.String](ctx),
		AttrCustomSetUnknown:        customtypes.NewSetValueNull[types.String](ctx),
		AttrCustomNestedListUnknown: customtypes.NewNestedListValueNull[modelEmptyTest](ctx),
		AttrCustomNestedSetUnknown:  customtypes.NewNestedSetValueNull[modelEmptyTest](ctx),
		AttrCustomNestedMapUnknown:  customtypes.NewNestedMapValueNull[modelEmptyTest](ctx),
		AttrCustomNestedList: customtypes.NewNestedListValue[modelCustomTypeTest](ctx, []modelCustomTypeTest{
			{
				AttrKnownString:   types.StringValue("val1"),
				AttrUnknownObject: customtypes.NewObjectValueNull[modelEmptyTest](ctx),
				AttrMANYUpper:     types.Int64Null(),
			},
			{
				AttrKnownString:   types.StringNull(),
				AttrUnknownObject: customtypes.NewObjectValue[modelEmptyTest](ctx, modelEmptyTest{}),
				AttrMANYUpper:     types.Int64Value(2),
			},
		}),
		AttrCustomNestedMap: customtypes.NewNestedMapValue[modelCustomTypeTest](ctx, map[string]modelCustomTypeTest{
			"keyOne": {
				AttrKnownString:   types.StringValue("val1"),
				AttrUnknownObject: customtypes.NewObjectValueNull[modelEmptyTest](ctx),
				AttrMANYUpper:     types.Int64Null(),
			},
			"keyTwo": {
				AttrKnownString:   types.StringNull(),
				AttrUnknownObject: customtypes.NewObjectValue[modelEmptyTest](ctx, modelEmptyTest{}),
				AttrMANYUpper:     types.Int64Value(2),
			},
		}),
	}
	require.NoError(t, autogen.ResolveUnknowns(&model))
	assert.Equal(t, modelExpected, model)
}
