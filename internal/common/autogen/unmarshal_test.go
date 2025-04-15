package autogen_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				"attrNull": null
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
}

func TestUnmarshalNestedAllTypes(t *testing.T) {
	type modelst struct {
		AttrObj               types.Object `tfsdk:"attr_obj"`
		AttrObjNullNotSent    types.Object `tfsdk:"attr_obj_null_not_sent"`
		AttrObjNullSent       types.Object `tfsdk:"attr_obj_null_sent"`
		AttrObjUnknownNotSent types.Object `tfsdk:"attr_obj_unknown_not_sent"`
		AttrObjUnknownSent    types.Object `tfsdk:"attr_obj_unknown_sent"`
		AttrObjParent         types.Object `tfsdk:"attr_obj_parent"`
		AttrListString        types.List   `tfsdk:"attr_list_string"`
		AttrListObj           types.List   `tfsdk:"attr_list_obj"`
		AttrSetString         types.Set    `tfsdk:"attr_set_string"`
		AttrSetObj            types.Set    `tfsdk:"attr_set_obj"`
		AttrListListString    types.List   `tfsdk:"attr_list_list_string"`
		AttrSetListObj        types.Set    `tfsdk:"attr_set_list_obj"`
		AttrListObjKnown      types.List   `tfsdk:"attr_list_obj_known"`
	}
	model := modelst{
		AttrObj: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
			// these attribute values are irrelevant, they will be overwritten with JSON values
			"attr_string": types.StringValue("different_string"),
			"attr_int":    types.Int64Value(123456),
			"attr_float":  types.Float64Unknown(), // can even be null
			"attr_bool":   types.BoolUnknown(),    // can even be unknown
		}),
		AttrObjNullNotSent:    types.ObjectNull(objTypeTest.AttrTypes),
		AttrObjNullSent:       types.ObjectNull(objTypeTest.AttrTypes),
		AttrObjUnknownNotSent: types.ObjectUnknown(objTypeTest.AttrTypes), // unknown values are changed to null
		AttrObjUnknownSent:    types.ObjectUnknown(objTypeTest.AttrTypes),
		AttrObjParent:         types.ObjectNull(objTypeParentTest.AttrTypes),
		AttrListString:        types.ListUnknown(types.StringType),
		AttrListObj:           types.ListUnknown(objTypeTest),
		AttrSetString:         types.SetUnknown(types.StringType),
		AttrSetObj:            types.SetUnknown(objTypeTest),
		AttrListListString:    types.ListUnknown(types.ListType{ElemType: types.StringType}),
		AttrSetListObj:        types.SetUnknown(types.ListType{ElemType: objTypeTest}),
		AttrListObjKnown: types.ListValueMust(objTypeTest, []attr.Value{
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("val"),
				"attr_int":    types.Int64Value(1),
				"attr_float":  types.Float64Value(1.1),
				"attr_bool":   types.BoolValue(true),
			}),
		}),
	}
	// attrUnexisting is ignored because it is in JSON but not in the model, no error is returned
	const (
		jsonResp = `
			{
				"attrObj": {
					"attrString": "value_string",
					"attrInt": 123,
					"attrFloat": 1.1,
					"attrBool": true,
					"attrUnexisting": "val"
				}, 
				"attrObjNullSent": {
					"attrString": "null_obj",
					"attrInt": 1,
					"attrFloat": null
				},
				"attrObjUnknownSent": {
					"attrString": "unknown_obj"
				},
				"attrObjParent": {
					"attrParentString": "parent string",
					"attrParentObj": {
						"attrString": "inside parent string"
					}
				},
				"attrListString": [
					"list1",
					"list2"
				],
				"attrListObj": [
					{
						"attrString": "list1",
						"attrInt": 1,
						"attrFloat": 1.1,
						"attrBool": true
					},
					{
						"attrString": "list2",
						"attrInt": 2,
						"attrFloat": 2.2,
						"attrBool": false
					}
				],
				"attrSetString": [
					"set1",
					"set2"
				],
				"attrSetObj": [
					{
						"attrString": "set1",
						"attrInt": 11,
						"attrFloat": 11.1,
						"attrBool": false
					},
					{			
						"attrString": "set2",
						"attrInt": 22,
						"attrFloat": 22.2,		
						"attrBool": true		
					}
				],
				"attrListListString": [
					["list1a", "list1b"],
					["list2a", "list2b", "list2c"]
				],
				"attrSetListObj": [
					[{
						"attrString": "setList1",
						"attrInt": 1,
						"attrFloat": 1.1,
						"attrBool": true
					},
					{
						"attrString": "setList2",	
						"attrInt": 2,
						"attrFloat": 2.2,
						"attrBool": false
					}],
					[{
						"attrString": "setList3",	
						"attrInt": 3,
						"attrFloat": 3.3,
						"attrBool": true
					},
					{
						"attrString": "setList4",
						"attrInt": 4,					
						"attrFloat": 4.4,
						"attrBool": false
					},
					{
						"attrString": "setList5",
						"attrInt": 5,
						"attrFloat": 5.5,
						"attrBool": true
					}]
				],
				"attrListObjKnown": [
					{
						"attrString": "val2",
						"attrInt": 2
					}
				]
			}
		`
	)
	modelExpected := modelst{
		AttrObj: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
			"attr_string": types.StringValue("value_string"),
			"attr_int":    types.Int64Value(123),
			"attr_float":  types.Float64Value(1.1),
			"attr_bool":   types.BoolValue(true),
		}),
		AttrObjNullNotSent: types.ObjectNull(objTypeTest.AttrTypes),
		AttrObjNullSent: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
			"attr_string": types.StringValue("null_obj"),
			"attr_int":    types.Int64Value(1),
			"attr_float":  types.Float64Null(),
			"attr_bool":   types.BoolNull(),
		}),
		AttrObjUnknownNotSent: types.ObjectUnknown(objTypeTest.AttrTypes),
		AttrObjUnknownSent: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
			"attr_string": types.StringValue("unknown_obj"),
			"attr_int":    types.Int64Null(),
			"attr_float":  types.Float64Null(),
			"attr_bool":   types.BoolNull(),
		}),
		AttrObjParent: types.ObjectValueMust(objTypeParentTest.AttrTypes, map[string]attr.Value{
			"attr_parent_string": types.StringValue("parent string"),
			"attr_parent_int":    types.Int64Null(),
			"attr_parent_obj": types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("inside parent string"),
				"attr_int":    types.Int64Null(),
				"attr_float":  types.Float64Null(),
				"attr_bool":   types.BoolNull(),
			}),
		}),
		AttrListString: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("list1"),
			types.StringValue("list2"),
		}),
		AttrListObj: types.ListValueMust(objTypeTest, []attr.Value{
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("list1"),
				"attr_int":    types.Int64Value(1),
				"attr_float":  types.Float64Value(1.1),
				"attr_bool":   types.BoolValue(true),
			}),
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("list2"),
				"attr_int":    types.Int64Value(2),
				"attr_float":  types.Float64Value(2.2),
				"attr_bool":   types.BoolValue(false),
			}),
		}),
		AttrSetString: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("set1"),
			types.StringValue("set2"),
		}),
		AttrSetObj: types.SetValueMust(objTypeTest, []attr.Value{
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("set1"),
				"attr_int":    types.Int64Value(11),
				"attr_float":  types.Float64Value(11.1),
				"attr_bool":   types.BoolValue(false),
			}),
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("set2"),
				"attr_int":    types.Int64Value(22),
				"attr_float":  types.Float64Value(22.2),
				"attr_bool":   types.BoolValue(true),
			}),
		}),
		AttrListListString: types.ListValueMust(types.ListType{ElemType: types.StringType}, []attr.Value{
			types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("list1a"),
				types.StringValue("list1b"),
			}),
			types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("list2a"),
				types.StringValue("list2b"),
				types.StringValue("list2c"),
			}),
		}),
		AttrSetListObj: types.SetValueMust(types.ListType{ElemType: objTypeTest}, []attr.Value{
			types.ListValueMust(objTypeTest, []attr.Value{
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList1"),
					"attr_int":    types.Int64Value(1),
					"attr_float":  types.Float64Value(1.1),
					"attr_bool":   types.BoolValue(true),
				}),
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList2"),
					"attr_int":    types.Int64Value(2),
					"attr_float":  types.Float64Value(2.2),
					"attr_bool":   types.BoolValue(false),
				}),
			}),
			types.ListValueMust(objTypeTest, []attr.Value{
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList3"),
					"attr_int":    types.Int64Value(3),
					"attr_float":  types.Float64Value(3.3),
					"attr_bool":   types.BoolValue(true),
				}),
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList4"),
					"attr_int":    types.Int64Value(4),
					"attr_float":  types.Float64Value(4.4),
					"attr_bool":   types.BoolValue(false),
				}),
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList5"),
					"attr_int":    types.Int64Value(5),
					"attr_float":  types.Float64Value(5.5),
					"attr_bool":   types.BoolValue(true),
				}),
			}),
		}),
		AttrListObjKnown: types.ListValueMust(objTypeTest, []attr.Value{
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("val2"),
				"attr_int":    types.Int64Value(2),
				"attr_float":  types.Float64Value(1.1),
				"attr_bool":   types.BoolValue(true),
			}),
		}),
	}
	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, modelExpected, model)
}

func TestUnmarshalZeroLenCollections(t *testing.T) {
	type modelst struct {
		ListNullAbsent  types.List `tfsdk:"list_null_absent"`
		ListNullEmpty   types.List `tfsdk:"list_null_present"`
		ListNullNull    types.List `tfsdk:"list_null_present_null"`
		ListEmptyAbsent types.List `tfsdk:"list_empty_absent"`
		ListEmptyEmpty  types.List `tfsdk:"list_empty_present"`
		ListEmptyNull   types.List `tfsdk:"list_empty_present_null"`
	}
	model := modelst{
		ListNullAbsent:  types.ListNull(types.StringType),
		ListNullEmpty:   types.ListNull(types.StringType),
		ListNullNull:    types.ListNull(types.StringType),
		ListEmptyAbsent: types.ListValueMust(types.StringType, []attr.Value{}),
		ListEmptyEmpty:  types.ListValueMust(types.StringType, []attr.Value{}),
		ListEmptyNull:   types.ListValueMust(types.StringType, []attr.Value{}),
	}
	const (
		jsonResp = `
			{
				"list_null_empty": [],
				"list_null_null": null,
				"list_empty_empty": [],
				"list_empty_null": null
			}
		`
	)
	modelExpected := modelst{
		ListNullAbsent:  types.ListNull(types.StringType),
		ListNullEmpty:   types.ListNull(types.StringType),
		ListNullNull:    types.ListNull(types.StringType),
		ListEmptyAbsent: types.ListValueMust(types.StringType, []attr.Value{}),
		ListEmptyEmpty:  types.ListValueMust(types.StringType, []attr.Value{}),
		ListEmptyNull:   types.ListValueMust(types.StringType, []attr.Value{}),
	}
	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, modelExpected, model)
}

func TestUnmarshalErrors(t *testing.T) {
	testCases := map[string]struct {
		model        any
		responseJSON string
		errorStr     string
	}{
		"response ints are not converted to model strings": {
			errorStr:     "unmarshal can't assign value to model field Attr",
			responseJSON: `{"attr": 123}`,
			model: &struct {
				Attr types.String
			}{},
		},
		"response strings are not converted to model ints": {
			errorStr:     "unmarshal can't assign value to model field Attr",
			responseJSON: `{"attr": "hello"}`,
			model: &struct {
				Attr types.Int64
			}{},
		},
		"response strings are not converted to model bools": {
			errorStr:     "unmarshal can't assign value to model field Attr",
			responseJSON: `{"attr": "true"}`,
			model: &struct {
				Attr types.Bool
			}{},
		},
		"response bools are not converted to model string": {
			errorStr:     "unmarshal can't assign value to model field Attr",
			responseJSON: `{"attr": true}`,
			model: &struct {
				Attr types.String
			}{},
		},
		"model attributes have to be of Terraform types": {
			errorStr:     "unmarshal can't assign value to model field Attr",
			responseJSON: `{"attr": "hello"}`,
			model: &struct {
				Attr string
			}{},
		},
		"model attr types in objects must match JSON types - string": {
			errorStr:     "unmarshal gets incorrect number for value: 1",
			responseJSON: `{ "attrObj": { "attrString": 1 } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
		"model attr types in objects must match JSON types - bool": {
			errorStr:     "unmarshal gets incorrect string for value: not a bool",
			responseJSON: `{ "attrObj": { "attrBool": "not a bool" } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
		"model attr types in objects must match JSON types - int": {
			errorStr:     "unmarshal gets incorrect string for value: not an int",
			responseJSON: `{ "attrObj": { "attrInt": "not an int" } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
		"model attr types in objects must match JSON types - float": {
			errorStr:     "unmarshal gets incorrect string for value: not an int",
			responseJSON: `{ "attrObj": { "attrFloat": "not an int" } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
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

// TestUnmarshalUnsupportedResponse has JSON response types not supported yet.
// It will be updated when we add support for them.
func TestUnmarshalUnsupportedResponse(t *testing.T) {
	testCases := map[string]struct {
		model        any
		responseJSON string
		errorStr     string
	}{
		"JSON maps not supported yet": {
			model: &struct {
				AttrMap types.Map `tfsdk:"attr_map"`
			}{},
			responseJSON: `{"attrMap": {"key": "value"}}`,
			errorStr:     "unmarshal expects object for field attrMap",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.ErrorContains(t, autogen.Unmarshal([]byte(tc.responseJSON), tc.model), tc.errorStr)
		})
	}
}
