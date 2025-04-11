package autogeneration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogeneration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type modelTest struct {
	AttrFloat  types.Float64 `tfsdk:"attr_float"`
	AttrString types.String  `tfsdk:"attr_string"`
	AttrInt    types.Int64   `tfsdk:"attr_int"`
	AttrBool   types.Bool    `tfsdk:"attr_bool"`
}

type modelParentTest struct {
	AttrParentObj    types.Object `tfsdk:"attr_parent_obj"`
	AttrParentString types.String `tfsdk:"attr_parent_string"`
	AttrParentInt    types.Int64  `tfsdk:"attr_parent_int"`
}

var (
	objTypeTest = types.ObjectType{AttrTypes: map[string]attr.Type{
		"attr_float":  types.Float64Type,
		"attr_string": types.StringType,
		"attr_int":    types.Int64Type,
		"attr_bool":   types.BoolType,
	}}

	objTypeParentTest = types.ObjectType{AttrTypes: map[string]attr.Type{
		"attr_parent_obj":    objTypeTest,
		"attr_parent_string": types.StringType,
		"attr_parent_int":    types.Int64Type,
	}}
)

const epsilon = 10e-15 // float tolerance in test equality

func TestMarshalBasic(t *testing.T) {
	model := struct {
		AttrFloat  types.Float64 `tfsdk:"attr_float"`
		AttrString types.String  `tfsdk:"attr_string"`
		// values with tag `omitjson` are not marshaled, and they don't need to be Terraform types
		AttrOmit            types.String `tfsdk:"attr_omit" autogeneration:"omitjson"`
		AttrOmitNoTerraform string       `autogeneration:"omitjson"`
		AttrUnkown          types.String `tfsdk:"attr_unknown"`
		AttrNull            types.String `tfsdk:"attr_null"`
		AttrInt             types.Int64  `tfsdk:"attr_int"`
		AttrBoolTrue        types.Bool   `tfsdk:"attr_bool_true"`
		AttrBoolFalse       types.Bool   `tfsdk:"attr_bool_false"`
		AttrBoolNull        types.Bool   `tfsdk:"attr_bool_null"`
	}{
		AttrFloat:           types.Float64Value(1.234),
		AttrString:          types.StringValue("hello"),
		AttrOmit:            types.StringValue("omit"),
		AttrOmitNoTerraform: "omit",
		AttrUnkown:          types.StringUnknown(), // unknown values are not marshaled
		AttrNull:            types.StringNull(),    // null values are not marshaled
		AttrInt:             types.Int64Value(1),
		AttrBoolTrue:        types.BoolValue(true),
		AttrBoolFalse:       types.BoolValue(false),
		AttrBoolNull:        types.BoolNull(), // null values are not marshaled
	}
	const expectedJSON = `{ "attrString": "hello", "attrInt": 1, "attrFloat": 1.234, "attrBoolTrue": true, "attrBoolFalse": false }`
	raw, err := autogeneration.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))
}

func TestMarshalNestedAllTypes(t *testing.T) {
	attrListObj, diags := types.ListValueFrom(t.Context(), objTypeTest, []modelTest{
		{
			AttrString: types.StringValue("str1"),
			AttrInt:    types.Int64Value(1),
		},
		{
			AttrString: types.StringValue("str2"),
			AttrInt:    types.Int64Value(2),
		},
	})
	assert.False(t, diags.HasError())
	attrSetObj, diags := types.SetValueFrom(t.Context(), objTypeTest, []modelTest{
		{
			AttrString: types.StringValue("str11"),
			AttrInt:    types.Int64Value(11),
		},
		{
			AttrString: types.StringValue("str22"),
			AttrInt:    types.Int64Value(22),
		},
	})
	assert.False(t, diags.HasError())
	attrMapObj, diags := types.MapValueFrom(t.Context(), objTypeTest, map[string]modelTest{
		"keyOne": {
			AttrString: types.StringValue("str1"),
			AttrInt:    types.Int64Value(1),
		},
		"KeyTwo": { // don't change the key case when it's a map
			AttrString: types.StringValue("str2"),
			AttrInt:    types.Int64Value(2),
		},
	})
	assert.False(t, diags.HasError())
	model := struct {
		AttrString     types.String `tfsdk:"attr_string"`
		AttrListSimple types.List   `tfsdk:"attr_list_simple"`
		AttrListObj    types.List   `tfsdk:"attr_list_obj"`
		AttrSetSimple  types.Set    `tfsdk:"attr_set_simple"`
		AttrSetObj     types.Set    `tfsdk:"attr_set_obj"`
		AttrMapSimple  types.Map    `tfsdk:"attr_map_simple"`
		AttrMapObj     types.Map    `tfsdk:"attr_map_obj"`
	}{
		AttrString:     types.StringValue("val"),
		AttrListSimple: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("val1"), types.StringValue("val2")}),
		AttrListObj:    attrListObj,
		AttrSetSimple:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("val11"), types.StringValue("val22")}),
		AttrSetObj:     attrSetObj,
		AttrMapSimple: types.MapValueMust(types.StringType, map[string]attr.Value{
			"keyOne": types.StringValue("val1"),
			"KeyTwo": types.StringValue("val2"), // don't change the key case when it's a map
		}),
		AttrMapObj: attrMapObj,
	}
	const expectedJSON = `
		{
			"attrString": "val", 
			"attrListSimple": ["val1", "val2"],
			"attrListObj": [
				{ "attrString": "str1", "attrInt": 1 },
				{ "attrString": "str2", "attrInt": 2 }
			],
			"attrSetSimple": ["val11", "val22"],
			"attrSetObj": [
				{ "attrString": "str11", "attrInt": 11 },
				{ "attrString": "str22", "attrInt": 22 }
			],
			"attrMapSimple": {
				"keyOne": "val1",
				"KeyTwo": "val2"
			},
			"attrMapObj": {
				"keyOne": { "attrString": "str1", "attrInt": 1 },
				"KeyTwo": { "attrString": "str2", "attrInt": 2 }
			}
		}
	`
	raw, err := autogeneration.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))
}

func TestMarshalNestedMultiLevel(t *testing.T) {
	attrListObj, diags := types.ListValueFrom(t.Context(), objTypeParentTest, []modelParentTest{
		{
			AttrParentObj: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("str11"),
				"attr_int":    types.Int64Value(11),
				"attr_float":  types.Float64Value(11.1),
				"attr_bool":   types.BoolValue(true),
			}),
			AttrParentString: types.StringValue("str1"),
			AttrParentInt:    types.Int64Value(1),
		},
		{
			AttrParentObj: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("str22"),
				"attr_int":    types.Int64Value(22),
				"attr_float":  types.Float64Value(22.2),
				"attr_bool":   types.BoolValue(false),
			}),
			AttrParentString: types.StringValue("str2"),
			AttrParentInt:    types.Int64Value(2),
		},
	})
	assert.False(t, diags.HasError())

	model := struct {
		AttrString      types.String `tfsdk:"attr_string"`
		AttrListParents types.List   `tfsdk:"attr_list_parents"`
	}{
		AttrString:      types.StringValue("val"),
		AttrListParents: attrListObj,
	}
	const expectedJSON = `
		{
			"attrString": "val", 
			"attrListParents": [
				{
					"attrParentString": "str1",
					"attrParentInt": 1,
					"attrParentObj": {
						"attrString": "str11",			
						"attrInt": 11,
						"attrFloat": 11.1,
						"attrBool": true
					}				
				},
				{
					"attrParentString": "str2",
					"attrParentInt": 2,
					"attrParentObj": {		
						"attrString": "str22",	
						"attrInt": 22,
						"attrFloat": 22.2,
						"attrBool": false
					}
				}
			]
		}
	`
	raw, err := autogeneration.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))
}

func TestMarshalOmitJSONUpdate(t *testing.T) {
	const (
		expectedCreate = `{ "attr": "val1", "attrOmitUpdate": "val2" }`
		expectedUpdate = `{ "attr": "val1" }`
	)
	model := struct {
		Attr           types.String `tfsdk:"attr"`
		AttrOmitUpdate types.String `tfsdk:"attr_omit_update" autogeneration:"omitjsonupdate"`
		AttrOmit       types.String `tfsdk:"attr_omit" autogeneration:"omitjson"`
	}{
		Attr:           types.StringValue("val1"),
		AttrOmitUpdate: types.StringValue("val2"),
		AttrOmit:       types.StringValue("omit"),
	}
	create, errCreate := autogeneration.Marshal(&model, false)
	require.NoError(t, errCreate)
	assert.JSONEq(t, expectedCreate, string(create))

	update, errUpdate := autogeneration.Marshal(&model, true)
	require.NoError(t, errUpdate)
	assert.JSONEq(t, expectedUpdate, string(update))
}

func TestMarshalUnsupported(t *testing.T) {
	testCases := map[string]any{
		"Int32 not supported yet as it's not being used in any model": &struct {
			Attr types.Int32
		}{
			Attr: types.Int32Value(1),
		},
		"Float32 not supported yet as it's not being used in any model": &struct {
			Attr types.Float32
		}{
			Attr: types.Float32Value(1.0),
		},
	}
	for name, model := range testCases {
		t.Run(name, func(t *testing.T) {
			raw, err := autogeneration.Marshal(model, false)
			require.Error(t, err)
			assert.Nil(t, raw)
		})
	}
}

func TestMarshalPanic(t *testing.T) {
	str := "string"
	testCases := map[string]any{
		"no Terraform types": &struct {
			Attr string
		}{
			Attr: "a",
		},
		"no pointer": struct {
			Attr types.String
		}{
			Attr: types.StringValue("a"),
		},
		"no struct": &str,
	}
	for name, model := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Panics(t, func() {
				_, _ = autogeneration.Marshal(model, false)
			})
		})
	}
}

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
	require.NoError(t, autogeneration.Unmarshal([]byte(jsonResp), &model))
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
		AttrSetString         types.Set    `tfsdk:"attr_set_string"`
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
		AttrSetString:         types.SetUnknown(types.StringType),
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
				"attrSetString": [
					"set1",
					"set2"
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
		AttrObjUnknownNotSent: types.ObjectNull(objTypeTest.AttrTypes),
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
		AttrSetString: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("set1"),
			types.StringValue("set2"),
		}),
	}
	require.NoError(t, autogeneration.Unmarshal([]byte(jsonResp), &model))
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
			errorStr:     "unmarshal gets incorrect number for field attr_string, value: 1",
			responseJSON: `{ "attrObj": { "attrString": 1 } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
		"model attr types in objects must match JSON types - bool": {
			errorStr:     "unmarshal gets incorrect string for field attr_bool, value: not a bool",
			responseJSON: `{ "attrObj": { "attrBool": "not a bool" } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
		"model attr types in objects must match JSON types - int": {
			errorStr:     "unmarshal gets incorrect string for field attr_int, value: not an int",
			responseJSON: `{ "attrObj": { "attrInt": "not an int" } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
		"model attr types in objects must match JSON types - float": {
			errorStr:     "unmarshal gets incorrect string for field attr_float, value: not an int",
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
			assert.ErrorContains(t, autogeneration.Unmarshal([]byte(tc.responseJSON), tc.model), tc.errorStr)
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
			assert.Error(t, autogeneration.Unmarshal([]byte(tc.responseJSON), tc.model))
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
		"JSON arrays not supported yet": {
			model: &struct {
				AttrList types.List `tfsdk:"attr_list"`
			}{},
			responseJSON: `{"attrList": [{"key": "value"}]}`,
			errorStr:     "unmarshal not supported yet for type map[string]interface {} in list",
		},
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
			assert.ErrorContains(t, autogeneration.Unmarshal([]byte(tc.responseJSON), tc.model), tc.errorStr)
		})
	}
}
