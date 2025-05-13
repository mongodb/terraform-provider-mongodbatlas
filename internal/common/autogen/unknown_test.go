package autogen_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveUnknowns(t *testing.T) {
	type modelst struct {
		AttrStringUnknown types.String `tfsdk:"attr_string_unknown"`
		AttrObjectUnknown types.Object `tfsdk:"attr_object_unknown"`
		AttrListUnknown   types.List   `tfsdk:"attr_list_unknown"`
		AttrObject        types.Object `tfsdk:"attr_object"`
		AttrListString    types.List   `tfsdk:"attr_list_string"`
		AttrSetString     types.Set    `tfsdk:"attr_set_string"`
		AttrListObjObj    types.List   `tfsdk:"attr_list_obj_obj"`
		AttrMapUnknown    types.Map    `tfsdk:"attr_map_unknown"`
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
		AttrMapUnknown: types.MapUnknown(types.StringType),
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
		AttrMapUnknown: types.MapNull(types.StringType),
	}
	require.NoError(t, autogen.ResolveUnknowns(&model))
	assert.Equal(t, modelExpected, model)
}
