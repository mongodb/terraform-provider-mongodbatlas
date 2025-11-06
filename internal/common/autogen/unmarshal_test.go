package autogen_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const epsilon = 10e-15 // float tolerance in test equality

func TestUnmarshalBasic(t *testing.T) {
	var model struct {
		AttrFloat        types.Float64 `tfsdk:"attr_float"`
		AttrFloatWithInt types.Float64 `tfsdk:"attr_float_with_int"`
		AttrString       types.String  `tfsdk:"attr_string"`
		AttrNotInJSON    types.String  `tfsdk:"attr_not_in_json"`
		AttrInt          types.Int64   `tfsdk:"attr_int"`
		AttrIntWithFloat types.Int64   `tfsdk:"attr_int_with_float"`
		AttrTrue         types.Bool    `tfsdk:"attr_true"`
		AttrFalse        types.Bool    `tfsdk:"attr_false"`
		AttrMANYUpper    types.Int64   `tfsdk:"attr_many_upper"`
	}
	const (
		// attribute_not_in_model is ignored because it is not in the model, no error is thrown.
		// attribute_null is ignored because it is null, no error is thrown even if it is not in the model.
		jsonResp = `
			{
				"attrString": "value_string",
				"attrTrue": true,
				"attrFalse": false,
				"attrInt": 123,
				"attrIntWithFloat": 10.6,
				"attrFloat": 456.1,
				"attrFloatWithInt": 13,
				"attrNotInModel": "val",
				"attrNull": null,
				"attrMANYUpper": 1
			}
		`
	)
	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, "value_string", model.AttrString.ValueString())
	assert.True(t, model.AttrTrue.ValueBool())
	assert.False(t, model.AttrFalse.ValueBool())
	assert.Equal(t, int64(123), model.AttrInt.ValueInt64())
	assert.Equal(t, int64(10), model.AttrIntWithFloat.ValueInt64()) // response floats stored in model ints have their decimals stripped.
	assert.InEpsilon(t, float64(456.1), model.AttrFloat.ValueFloat64(), epsilon)
	assert.InEpsilon(t, float64(13), model.AttrFloatWithInt.ValueFloat64(), epsilon)
	assert.True(t, model.AttrNotInJSON.IsNull()) // attributes not in JSON response are not changed, so null is kept.
	assert.Equal(t, int64(1), model.AttrMANYUpper.ValueInt64())
}

func TestUnmarshalDynamicJSONAttr(t *testing.T) {
	var model struct {
		AttrDynamicJSONObject  jsontypes.Normalized `tfsdk:"attr_dynamic_json_object"`
		AttrDynamicJSONBoolean jsontypes.Normalized `tfsdk:"attr_dynamic_json_boolean"`
		AttrDynamicJSONNumber  jsontypes.Normalized `tfsdk:"attr_dynamic_json_number"`
		AttrDynamicJSONString  jsontypes.Normalized `tfsdk:"attr_dynamic_json_string"`
		AttrDynamicJSONArray   jsontypes.Normalized `tfsdk:"attr_dynamic_json_array"`
	}
	const jsonResp = `
		{
			"attrDynamicJSONObject": {"hello":"there"},
			"attrDynamicJSONBoolean": true,
			"attrDynamicJSONNumber": 1.234,
			"attrDynamicJSONString": "hello",
			"attrDynamicJSONArray": [1, 2, 3]
		}
	`
	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.JSONEq(t, "{\"hello\":\"there\"}", model.AttrDynamicJSONObject.ValueString())
	assert.JSONEq(t, "true", model.AttrDynamicJSONBoolean.ValueString())
	assert.JSONEq(t, "1.234", model.AttrDynamicJSONNumber.ValueString())
	assert.JSONEq(t, "\"hello\"", model.AttrDynamicJSONString.ValueString())
	assert.JSONEq(t, "[1, 2, 3]", model.AttrDynamicJSONArray.ValueString())
}

type unmarshalModelEmpty struct{}

type unmarshalModelCustomType struct {
	AttrFloat     types.Float64                                `tfsdk:"attr_float"`
	AttrString    types.String                                 `tfsdk:"attr_string"`
	AttrNested    customtypes.ObjectValue[unmarshalModelEmpty] `tfsdk:"attr_nested"`
	AttrInt       types.Int64                                  `tfsdk:"attr_int"`
	AttrBool      types.Bool                                   `tfsdk:"attr_bool"`
	AttrMANYUpper types.Int64                                  `tfsdk:"attr_many_upper"`
}

func TestUnmarshalCustomObject(t *testing.T) {
	ctx := context.Background()

	type modelst struct {
		AttrCustomObj               customtypes.ObjectValue[unmarshalModelCustomType] `tfsdk:"attr_custom_obj"`
		AttrCustomObjNullNotSent    customtypes.ObjectValue[unmarshalModelCustomType] `tfsdk:"attr_custom_obj_null_not_sent"`
		AttrCustomObjNullSent       customtypes.ObjectValue[unmarshalModelCustomType] `tfsdk:"attr_custom_obj_null_sent"`
		AttrCustomObjUnknownNotSent customtypes.ObjectValue[unmarshalModelCustomType] `tfsdk:"attr_custom_obj_unknown_not_sent"`
		AttrCustomObjUnknownSent    customtypes.ObjectValue[unmarshalModelCustomType] `tfsdk:"attr_custom_obj_unknown_sent"`
		AttrCustomObjParent         customtypes.ObjectValue[unmarshalModelCustomType] `tfsdk:"attr_custom_obj_parent"`
		AttrCustomObjZeroInit       customtypes.ObjectValue[unmarshalModelCustomType] `tfsdk:"attr_custom_obj_zero"`
	}

	model := modelst{
		AttrCustomObj: customtypes.NewObjectValue[unmarshalModelCustomType](ctx, unmarshalModelCustomType{
			AttrString:    types.StringValue("different_string"),
			AttrInt:       types.Int64Value(999),
			AttrFloat:     types.Float64Unknown(),
			AttrBool:      types.BoolUnknown(),
			AttrNested:    customtypes.NewObjectValueUnknown[unmarshalModelEmpty](ctx),
			AttrMANYUpper: types.Int64Value(999),
		}),
		AttrCustomObjNullNotSent:    customtypes.NewObjectValueNull[unmarshalModelCustomType](ctx),
		AttrCustomObjNullSent:       customtypes.NewObjectValueNull[unmarshalModelCustomType](ctx),
		AttrCustomObjUnknownNotSent: customtypes.NewObjectValueUnknown[unmarshalModelCustomType](ctx),
		AttrCustomObjUnknownSent:    customtypes.NewObjectValueUnknown[unmarshalModelCustomType](ctx),
		AttrCustomObjParent:         customtypes.NewObjectValueNull[unmarshalModelCustomType](ctx),
	}

	const (
		jsonResp = `
			{
				"attrCustomObj": {
					"attrString": "value_string",
					"attrInt": 123,
					"attrFloat": 1.1,
					"attrBool": true,
					"attrNested": {},
					"attrMANYUpper": 456
				},
				"attrCustomObjNullSent": {
					"attrString": "null_obj",
					"attrInt": 1,
					"attrFloat": null
				},
				"attrCustomObjUnknownSent": {
					"attrString": "unknown_obj"
				},
				"attrCustomObjParent": {
					"attrString": "parent string",
					"attrNested": {}
				},
				"attrCustomObjZeroInit": {
					"attrString": "zero init string",
					"attrNested": {}
				}
			}
		`
	)

	modelExpected := modelst{
		AttrCustomObj: customtypes.NewObjectValue[unmarshalModelCustomType](ctx, unmarshalModelCustomType{
			AttrString:    types.StringValue("value_string"),
			AttrInt:       types.Int64Value(123),
			AttrFloat:     types.Float64Value(1.1),
			AttrBool:      types.BoolValue(true),
			AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
			AttrMANYUpper: types.Int64Value(456),
		}),
		AttrCustomObjNullNotSent: customtypes.NewObjectValueNull[unmarshalModelCustomType](ctx),
		AttrCustomObjNullSent: customtypes.NewObjectValue[unmarshalModelCustomType](ctx, unmarshalModelCustomType{
			AttrString:    types.StringValue("null_obj"),
			AttrInt:       types.Int64Value(1),
			AttrFloat:     types.Float64Null(),
			AttrBool:      types.BoolNull(),
			AttrNested:    customtypes.NewObjectValueNull[unmarshalModelEmpty](ctx),
			AttrMANYUpper: types.Int64Null(),
		}),
		AttrCustomObjUnknownNotSent: customtypes.NewObjectValueUnknown[unmarshalModelCustomType](ctx),
		AttrCustomObjUnknownSent: customtypes.NewObjectValue[unmarshalModelCustomType](ctx, unmarshalModelCustomType{
			AttrString:    types.StringValue("unknown_obj"),
			AttrInt:       types.Int64Null(),
			AttrFloat:     types.Float64Null(),
			AttrBool:      types.BoolNull(),
			AttrNested:    customtypes.NewObjectValueNull[unmarshalModelEmpty](ctx),
			AttrMANYUpper: types.Int64Null(),
		}),
		AttrCustomObjParent: customtypes.NewObjectValue[unmarshalModelCustomType](ctx, unmarshalModelCustomType{
			AttrString:    types.StringValue("parent string"),
			AttrInt:       types.Int64Null(),
			AttrFloat:     types.Float64Null(),
			AttrBool:      types.BoolNull(),
			AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
			AttrMANYUpper: types.Int64Null(),
		}),
		AttrCustomObjZeroInit: customtypes.NewObjectValue[unmarshalModelCustomType](ctx, unmarshalModelCustomType{
			AttrString:    types.StringValue("zero init string"),
			AttrInt:       types.Int64Null(),
			AttrFloat:     types.Float64Null(),
			AttrBool:      types.BoolNull(),
			AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
			AttrMANYUpper: types.Int64Null(),
		}),
	}

	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, modelExpected, model)
}

func TestUnmarshalCustomList(t *testing.T) {
	ctx := context.Background()

	type modelst struct {
		AttrCustomListString               customtypes.ListValue[types.String]                   `tfsdk:"attr_custom_list_string"`
		AttrCustomNestedList               customtypes.NestedListValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_list"`
		AttrCustomNestedListNullNotSent    customtypes.NestedListValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_list_null_not_sent"`
		AttrCustomNestedListNullSent       customtypes.NestedListValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_list_null_sent"`
		AttrCustomNestedListUnknownNotSent customtypes.NestedListValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_list_unknown_not_sent"`
		AttrCustomNestedListUnknownSent    customtypes.NestedListValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_list_unknown_sent"`
		AttrCustomNestedListZeroInit       customtypes.NestedListValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_list_zero"`
	}

	model := modelst{
		AttrCustomListString: customtypes.NewListValueUnknown[types.String](ctx),
		AttrCustomNestedList: customtypes.NewNestedListValue[unmarshalModelCustomType](ctx, []unmarshalModelCustomType{
			{
				AttrString:    types.StringValue("different_string"),
				AttrInt:       types.Int64Value(999),
				AttrFloat:     types.Float64Unknown(),
				AttrBool:      types.BoolUnknown(),
				AttrNested:    customtypes.NewObjectValueUnknown[unmarshalModelEmpty](ctx),
				AttrMANYUpper: types.Int64Value(999),
			},
			{
				AttrString:    types.StringValue("existing not overwritten"),
				AttrInt:       types.Int64Unknown(),
				AttrFloat:     types.Float64Unknown(),
				AttrBool:      types.BoolUnknown(),
				AttrNested:    customtypes.NewObjectValueUnknown[unmarshalModelEmpty](ctx),
				AttrMANYUpper: types.Int64Value(999),
			},
		}),
		AttrCustomNestedListNullNotSent:    customtypes.NewNestedListValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedListNullSent:       customtypes.NewNestedListValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedListUnknownNotSent: customtypes.NewNestedListValueUnknown[unmarshalModelCustomType](ctx),
		AttrCustomNestedListUnknownSent:    customtypes.NewNestedListValueUnknown[unmarshalModelCustomType](ctx),
	}

	const (
		jsonResp = `
			{
				"attrCustomListString": [
					"list1",
					"list2"
				],
				"attrCustomNestedList": [
					{
						"attrString": "nestedList1",
						"attrInt": 1,
						"attrFloat": 1.1,
						"attrBool": true,
						"attrNested": {},
						"attrMANYUpper": 123
					},
					{
						"attrFloat": 2.2,
						"attrBool": false,
						"attrNested": {},
						"attrMANYUpper": 456
					}
				],
				"attrCustomNestedListNullSent": null,
				"attrCustomNestedListUnknownSent": [
					{
						"attrString": "unknownSent"
					}
				],
				"attrCustomNestedListZeroInit": [
					{
						"attrString": "zero init string",
						"attrNested": {}
					}
				]
			}
		`
	)

	modelExpected := modelst{
		AttrCustomListString: customtypes.NewListValue[types.String](ctx, []attr.Value{
			types.StringValue("list1"),
			types.StringValue("list2"),
		}),
		AttrCustomNestedList: customtypes.NewNestedListValue[unmarshalModelCustomType](ctx, []unmarshalModelCustomType{
			{
				AttrString:    types.StringValue("nestedList1"),
				AttrInt:       types.Int64Value(1),
				AttrFloat:     types.Float64Value(1.1),
				AttrBool:      types.BoolValue(true),
				AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
				AttrMANYUpper: types.Int64Value(123),
			},
			{
				AttrString:    types.StringValue("existing not overwritten"),
				AttrInt:       types.Int64Unknown(),
				AttrFloat:     types.Float64Value(2.2),
				AttrBool:      types.BoolValue(false),
				AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
				AttrMANYUpper: types.Int64Value(456),
			},
		}),
		AttrCustomNestedListNullNotSent:    customtypes.NewNestedListValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedListNullSent:       customtypes.NewNestedListValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedListUnknownNotSent: customtypes.NewNestedListValueUnknown[unmarshalModelCustomType](ctx),
		AttrCustomNestedListUnknownSent: customtypes.NewNestedListValue[unmarshalModelCustomType](ctx, []unmarshalModelCustomType{
			{
				AttrString:    types.StringValue("unknownSent"),
				AttrInt:       types.Int64Null(),
				AttrFloat:     types.Float64Null(),
				AttrBool:      types.BoolNull(),
				AttrNested:    customtypes.NewObjectValueNull[unmarshalModelEmpty](ctx),
				AttrMANYUpper: types.Int64Null(),
			},
		}),
		AttrCustomNestedListZeroInit: customtypes.NewNestedListValue[unmarshalModelCustomType](ctx, []unmarshalModelCustomType{
			{
				AttrString:    types.StringValue("zero init string"),
				AttrInt:       types.Int64Null(),
				AttrFloat:     types.Float64Null(),
				AttrBool:      types.BoolNull(),
				AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
				AttrMANYUpper: types.Int64Null(),
			},
		}),
	}

	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, modelExpected, model)
}

func TestUnmarshalCustomSet(t *testing.T) {
	ctx := context.Background()

	type modelst struct {
		AttrCustomSetString               customtypes.SetValue[types.String]                   `tfsdk:"attr_custom_set_string"`
		AttrCustomNestedSet               customtypes.NestedSetValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_set"`
		AttrCustomNestedSetNullNotSent    customtypes.NestedSetValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_set_null_not_sent"`
		AttrCustomNestedSetNullSent       customtypes.NestedSetValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_set_null_sent"`
		AttrCustomNestedSetUnknownNotSent customtypes.NestedSetValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_set_unknown_not_sent"`
		AttrCustomNestedSetUnknownSent    customtypes.NestedSetValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_set_unknown_sent"`
		AttrCustomNestedSetZeroInit       customtypes.NestedSetValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_set_zero"`
	}

	model := modelst{
		AttrCustomSetString: customtypes.NewSetValueUnknown[types.String](ctx),
		AttrCustomNestedSet: customtypes.NewNestedSetValue[unmarshalModelCustomType](ctx, []unmarshalModelCustomType{{
			AttrString:    types.StringValue("different_string"),
			AttrInt:       types.Int64Value(999),
			AttrFloat:     types.Float64Unknown(),
			AttrBool:      types.BoolUnknown(),
			AttrNested:    customtypes.NewObjectValueUnknown[unmarshalModelEmpty](ctx),
			AttrMANYUpper: types.Int64Value(999),
		}}),
		AttrCustomNestedSetNullNotSent:    customtypes.NewNestedSetValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedSetNullSent:       customtypes.NewNestedSetValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedSetUnknownNotSent: customtypes.NewNestedSetValueUnknown[unmarshalModelCustomType](ctx),
		AttrCustomNestedSetUnknownSent:    customtypes.NewNestedSetValueUnknown[unmarshalModelCustomType](ctx),
	}

	const (
		jsonResp = `
			{
				"attrCustomSetString": [
					"set1",
					"set2"
				],
				"attrCustomNestedSet": [
					{
						"attrString": "nestedSet1",
						"attrInt": 1,
						"attrFloat": 1.1,
						"attrBool": true,
						"attrNested": {},
						"attrMANYUpper": 123
					},
					{
						"attrString": "nestedSet2",
						"attrInt": 2,
						"attrFloat": 2.2,
						"attrBool": false,
						"attrNested": {},
						"attrMANYUpper": 456
					}
				],
				"attrCustomNestedSetNullSent": null,
				"attrCustomNestedSetUnknownSent": [
					{
						"attrString": "unknownSetSent"
					}
				],
				"attrCustomNestedSetZeroInit": [
					{
						"attrString": "zero init set string",
						"attrNested": {}
					}
				]
			}
		`
	)

	modelExpected := modelst{
		AttrCustomSetString: customtypes.NewSetValue[types.String](ctx, []attr.Value{
			types.StringValue("set1"),
			types.StringValue("set2"),
		}),
		AttrCustomNestedSet: customtypes.NewNestedSetValue[unmarshalModelCustomType](ctx, []unmarshalModelCustomType{
			{
				AttrString:    types.StringValue("nestedSet1"),
				AttrInt:       types.Int64Value(1),
				AttrFloat:     types.Float64Value(1.1),
				AttrBool:      types.BoolValue(true),
				AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
				AttrMANYUpper: types.Int64Value(123),
			},
			{
				AttrString:    types.StringValue("nestedSet2"),
				AttrInt:       types.Int64Value(2),
				AttrFloat:     types.Float64Value(2.2),
				AttrBool:      types.BoolValue(false),
				AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
				AttrMANYUpper: types.Int64Value(456),
			},
		}),
		AttrCustomNestedSetNullNotSent:    customtypes.NewNestedSetValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedSetNullSent:       customtypes.NewNestedSetValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedSetUnknownNotSent: customtypes.NewNestedSetValueUnknown[unmarshalModelCustomType](ctx),
		AttrCustomNestedSetUnknownSent: customtypes.NewNestedSetValue[unmarshalModelCustomType](ctx, []unmarshalModelCustomType{
			{
				AttrString:    types.StringValue("unknownSetSent"),
				AttrInt:       types.Int64Null(),
				AttrFloat:     types.Float64Null(),
				AttrBool:      types.BoolNull(),
				AttrNested:    customtypes.NewObjectValueNull[unmarshalModelEmpty](ctx),
				AttrMANYUpper: types.Int64Null(),
			},
		}),
		AttrCustomNestedSetZeroInit: customtypes.NewNestedSetValue[unmarshalModelCustomType](ctx, []unmarshalModelCustomType{
			{
				AttrString:    types.StringValue("zero init set string"),
				AttrInt:       types.Int64Null(),
				AttrFloat:     types.Float64Null(),
				AttrBool:      types.BoolNull(),
				AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
				AttrMANYUpper: types.Int64Null(),
			},
		}),
	}

	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, modelExpected, model)
}

func TestUnmarshalCustomMap(t *testing.T) {
	ctx := context.Background()

	type modelst struct {
		AttrCustomMapString               customtypes.MapValue[types.String]                   `tfsdk:"attr_custom_map_string"`
		AttrCustomNestedMap               customtypes.NestedMapValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_map"`
		AttrCustomNestedMapNullNotSent    customtypes.NestedMapValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_map_null_not_sent"`
		AttrCustomNestedMapNullSent       customtypes.NestedMapValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_map_null_sent"`
		AttrCustomNestedMapUnknownNotSent customtypes.NestedMapValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_map_unknown_not_sent"`
		AttrCustomNestedMapUnknownSent    customtypes.NestedMapValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_map_unknown_sent"`
		AttrCustomNestedMapZeroInit       customtypes.NestedMapValue[unmarshalModelCustomType] `tfsdk:"attr_custom_nested_map_zero"`
	}

	model := modelst{
		AttrCustomMapString: customtypes.NewMapValueUnknown[types.String](ctx),
		AttrCustomNestedMap: customtypes.NewNestedMapValue[unmarshalModelCustomType](ctx, map[string]unmarshalModelCustomType{
			"keyOne": {
				AttrString:    types.StringValue("different_string"),
				AttrInt:       types.Int64Value(999),
				AttrFloat:     types.Float64Unknown(),
				AttrBool:      types.BoolUnknown(),
				AttrNested:    customtypes.NewObjectValueUnknown[unmarshalModelEmpty](ctx),
				AttrMANYUpper: types.Int64Value(999),
			},
			"keyTwo": {
				AttrString:    types.StringValue("existing not overwritten"),
				AttrInt:       types.Int64Unknown(),
				AttrFloat:     types.Float64Unknown(),
				AttrBool:      types.BoolUnknown(),
				AttrNested:    customtypes.NewObjectValueUnknown[unmarshalModelEmpty](ctx),
				AttrMANYUpper: types.Int64Value(999),
			},
		}),
		AttrCustomNestedMapNullNotSent:    customtypes.NewNestedMapValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedMapNullSent:       customtypes.NewNestedMapValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedMapUnknownNotSent: customtypes.NewNestedMapValueUnknown[unmarshalModelCustomType](ctx),
		AttrCustomNestedMapUnknownSent:    customtypes.NewNestedMapValueUnknown[unmarshalModelCustomType](ctx),
	}

	const (
		jsonResp = `
			{
				"attrCustomMapString": {
					"keyOne": "map1",
					"KeyTwo": "map2"
				},
				"attrCustomNestedMap": {
					"keyOne": {
						"attrString": "nestedMap1",
						"attrInt": 1,
						"attrFloat": 1.1,
						"attrBool": true,
						"attrNested": {},
						"attrMANYUpper": 123
					},
					"keyTwo": {
						"attrFloat": 2.2,
						"attrBool": false,
						"attrNested": {},
						"attrMANYUpper": 456
					}
				},
				"attrCustomNestedMapNullSent": null,
				"attrCustomNestedMapUnknownSent": {
					"keyOne": {
						"attrString": "unknownMapSent"
					}
				},
				"attrCustomNestedMapZeroInit": {
					"keyOne": {
						"attrString": "zero init map string",
						"attrNested": {}
					}
				}
			}
		`
	)

	modelExpected := modelst{
		AttrCustomMapString: customtypes.NewMapValue[types.String](ctx, map[string]attr.Value{
			"keyOne": types.StringValue("map1"),
			"KeyTwo": types.StringValue("map2"),
		}),
		AttrCustomNestedMap: customtypes.NewNestedMapValue[unmarshalModelCustomType](ctx, map[string]unmarshalModelCustomType{
			"keyOne": {
				AttrString:    types.StringValue("nestedMap1"),
				AttrInt:       types.Int64Value(1),
				AttrFloat:     types.Float64Value(1.1),
				AttrBool:      types.BoolValue(true),
				AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
				AttrMANYUpper: types.Int64Value(123),
			},
			"keyTwo": {
				AttrString:    types.StringValue("existing not overwritten"),
				AttrInt:       types.Int64Unknown(),
				AttrFloat:     types.Float64Value(2.2),
				AttrBool:      types.BoolValue(false),
				AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
				AttrMANYUpper: types.Int64Value(456),
			},
		}),
		AttrCustomNestedMapNullNotSent:    customtypes.NewNestedMapValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedMapNullSent:       customtypes.NewNestedMapValueNull[unmarshalModelCustomType](ctx),
		AttrCustomNestedMapUnknownNotSent: customtypes.NewNestedMapValueUnknown[unmarshalModelCustomType](ctx),
		AttrCustomNestedMapUnknownSent: customtypes.NewNestedMapValue[unmarshalModelCustomType](ctx, map[string]unmarshalModelCustomType{
			"keyOne": {
				AttrString:    types.StringValue("unknownMapSent"),
				AttrInt:       types.Int64Null(),
				AttrFloat:     types.Float64Null(),
				AttrBool:      types.BoolNull(),
				AttrNested:    customtypes.NewObjectValueNull[unmarshalModelEmpty](ctx),
				AttrMANYUpper: types.Int64Null(),
			},
		}),
		AttrCustomNestedMapZeroInit: customtypes.NewNestedMapValue[unmarshalModelCustomType](ctx, map[string]unmarshalModelCustomType{
			"keyOne": {
				AttrString:    types.StringValue("zero init map string"),
				AttrInt:       types.Int64Null(),
				AttrFloat:     types.Float64Null(),
				AttrBool:      types.BoolNull(),
				AttrNested:    customtypes.NewObjectValue[unmarshalModelEmpty](ctx, unmarshalModelEmpty{}),
				AttrMANYUpper: types.Int64Null(),
			},
		}),
	}

	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, modelExpected, model)
}

func TestUnmarshalEmptyJSON(t *testing.T) {
	model := struct {
		Attr types.String `tfsdk:"attr"`
	}{
		Attr: types.StringValue("hello"),
	}
	require.NoError(t, autogen.Unmarshal([]byte(""), &model))
	require.NoError(t, autogen.Unmarshal(nil, &model))
	assert.Equal(t, types.StringValue("hello"), model.Attr)
}

func TestUnmarshalErrors(t *testing.T) {
	ctx := context.Background()

	type testNestedObject struct {
		AttrFloat  types.Float64 `tfsdk:"attr_float"`
		AttrString types.String  `tfsdk:"attr_string"`
		AttrInt    types.Int64   `tfsdk:"attr_int"`
		AttrBool   types.Bool    `tfsdk:"attr_bool"`
	}

	testCases := map[string]struct {
		model        any
		responseJSON string
		errorStr     string
	}{
		"response ints are not converted to model strings": {
			errorStr:     "unmarshal of attribute attr expects type StringType but got Number with value: 1",
			responseJSON: `{"attr": 123}`,
			model: &struct {
				Attr types.String
			}{},
		},
		"response strings are not converted to model ints": {
			errorStr:     "unmarshal of attribute attr expects type Int64Type but got String with value: hello",
			responseJSON: `{"attr": "hello"}`,
			model: &struct {
				Attr types.Int64
			}{},
		},
		"response strings are not converted to model bools": {
			errorStr:     "unmarshal of attribute attr expects type BoolType but got String with value: true",
			responseJSON: `{"attr": "true"}`,
			model: &struct {
				Attr types.Bool
			}{},
		},
		"response bools are not converted to model string": {
			errorStr:     "unmarshal of attribute attr expects type StringType but got Bool with value: true",
			responseJSON: `{"attr": true}`,
			model: &struct {
				Attr types.String
			}{},
		},
		"model attributes have to be of Terraform types": {
			errorStr:     "unmarshal trying to set non-Terraform attribute Attr",
			responseJSON: `{"attr": "hello"}`,
			model: &struct {
				Attr string
			}{},
		},
		"model attr types in objects must match JSON types - string": {
			errorStr:     "unmarshal of attribute attr_string expects type StringType but got Number with value: 1",
			responseJSON: `{ "attrObj": { "attrString": 1 } }`,
			model: &struct {
				AttrObj customtypes.ObjectValue[testNestedObject] `tfsdk:"attr_obj"`
			}{
				AttrObj: customtypes.NewObjectValueNull[testNestedObject](ctx),
			},
		},
		"model attr types in objects must match JSON types - bool": {
			errorStr:     "unmarshal of attribute attr_bool expects type BoolType but got String with value: not a bool",
			responseJSON: `{ "attrObj": { "attrBool": "not a bool" } }`,
			model: &struct {
				AttrObj customtypes.ObjectValue[testNestedObject] `tfsdk:"attr_obj"`
			}{
				AttrObj: customtypes.NewObjectValueNull[testNestedObject](ctx),
			},
		},
		"model attr types in objects must match JSON types - int": {
			errorStr:     "unmarshal of attribute attr_int expects type Int64Type but got String with value: not an int",
			responseJSON: `{ "attrObj": { "attrInt": "not an int" } }`,
			model: &struct {
				AttrObj customtypes.ObjectValue[testNestedObject] `tfsdk:"attr_obj"`
			}{
				AttrObj: customtypes.NewObjectValueNull[testNestedObject](ctx),
			},
		},
		"model attr types in objects must match JSON types - float": {
			errorStr:     "unmarshal of attribute attr_float expects type Float64Type but got String with value: not an int",
			responseJSON: `{ "attrObj": { "attrFloat": "not an int" } }`,
			model: &struct {
				AttrObj customtypes.ObjectValue[testNestedObject] `tfsdk:"attr_obj"`
			}{
				AttrObj: customtypes.NewObjectValueNull[testNestedObject](ctx),
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.ErrorContains(t, autogen.Unmarshal([]byte(tc.responseJSON), tc.model), tc.errorStr)
		})
	}
}

// TestUnmarshalUnsupportedModel has Terraform types not supported yet.
// It will be updated when we add support for them.
func TestUnmarshalUnsupportedModel(t *testing.T) {
	testCases := map[string]struct {
		model        any
		responseJSON string
	}{
		"Int32 not supported yet as it's not being used in any model": {
			responseJSON: `{"attr": 1}`,
			model: &struct {
				Attr types.Int32
			}{},
		},
		"Float32 not supported yet as it's not being used in any model": {
			responseJSON: `{"attr": 1}`,
			model: &struct {
				Attr types.Float32
			}{},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Error(t, autogen.Unmarshal([]byte(tc.responseJSON), tc.model))
		})
	}
}
