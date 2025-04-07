package autogeneration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogeneration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TFTestModel struct {
	AttrString types.String `tfsdk:"attr_string"`
	AttrInt    types.Int64  `tfsdk:"attr_int"`
}

var TestObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"attr_string": types.StringType,
	"attr_int":    types.Int64Type,
}}

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
	}{
		AttrFloat:           types.Float64Value(1.234),
		AttrString:          types.StringValue("hello"),
		AttrOmit:            types.StringValue("omit"),
		AttrOmitNoTerraform: "omit",
		AttrUnkown:          types.StringUnknown(), // unknown values are not marshaled
		AttrNull:            types.StringNull(),    // null values are not marshaled
		AttrInt:             types.Int64Value(1),
	}
	const expectedJSON = `{ "attrString": "hello", "attrInt": 1, "attrFloat": 1.234 }`
	raw, err := autogeneration.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))
}

func TestMarshalNested(t *testing.T) {
	attrListObj, diags := types.ListValueFrom(t.Context(), TestObjType, []TFTestModel{
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
	attrSetObj, diags := types.SetValueFrom(t.Context(), TestObjType, []TFTestModel{
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
	model := struct {
		AttrString     types.String `tfsdk:"attr_string"`
		AttrListSimple types.List   `tfsdk:"attr_list_simple"`
		AttrListObj    types.List   `tfsdk:"attr_list_obj"`
		AttrSetObj     types.Set    `tfsdk:"attr_set_obj"`
	}{
		AttrString:     types.StringValue("val"),
		AttrListSimple: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("val1"), types.StringValue("val2")}),
		AttrListObj:    attrListObj,
		AttrSetObj:     attrSetObj,
	}
	const expectedJSON = `
		{
			"attrString": "val", 
			"attrListSimple": ["val1", "val2"],
			"attrListObj": [
				{ "attr_string": "str1", "attr_int": 1 },
				{ "attr_string": "str2", "attr_int": 2 }
			],
			"attrSetObj": [
				{ "attr_string": "str11", "attr_int": 11 },
				{ "attr_string": "str22", "attr_int": 22 }
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
		"Map not supported yet, only no-nested types": &struct {
			Attr types.Map
		}{
			Attr: types.MapValueMust(types.StringType, map[string]attr.Value{
				"key": types.StringValue("value"),
			}),
		},
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
		epsilon = 10e-15 // float tolerance
		// attribute_not_in_model is ignored because it is not in the model, no error is thrown.
		// attribute_null is ignored because it is null, no error is thrown even if it is not in the model.
		tfResponseJSON = `
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
	require.NoError(t, autogeneration.Unmarshal([]byte(tfResponseJSON), &model))
	assert.Equal(t, "value_string", model.AttrString.ValueString())
	assert.True(t, model.AttrTrue.ValueBool())
	assert.False(t, model.AttrFalse.ValueBool())
	assert.Equal(t, int64(123), model.AttrInt.ValueInt64())
	assert.Equal(t, int64(10), model.AttrIntWithFloat.ValueInt64()) // response floats stored in model ints have their decimals stripped.
	assert.InEpsilon(t, float64(456.1), model.AttrFloat.ValueFloat64(), epsilon)
	assert.InEpsilon(t, float64(13), model.AttrFloatWithInt.ValueFloat64(), epsilon)
	assert.True(t, model.AttrNotInJSON.IsNull()) // attributes not in JSON response are not changed, so null is kept.
}

func TestUnmarshalErrors(t *testing.T) {
	const errorStr = "can't assign value to model field Attr"
	testCases := map[string]struct {
		model        any
		responseJSON string
	}{
		"response ints are not converted to model strings": {
			responseJSON: `{"attr": 123}`, //
			model: &struct {
				Attr types.String
			}{},
		},
		"response strings are not converted to model ints": {
			responseJSON: `{"attr": "hello"}`,
			model: &struct {
				Attr types.Int64
			}{},
		},
		"response strings are not converted to model bools": {
			responseJSON: `{"attr": "true"}`,
			model: &struct {
				Attr types.Bool
			}{},
		},
		"response bools are not converted to model string": {
			responseJSON: `{"attr": true}`,
			model: &struct {
				Attr types.String
			}{},
		},
		"model attributes have to be of Terraform types": {
			responseJSON: `{"attr": "hello"}`,
			model: &struct {
				Attr string
			}{},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.ErrorContains(t, autogeneration.Unmarshal([]byte(tc.responseJSON), tc.model), errorStr)
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
	var model struct {
		Attr types.String
	}
	testCases := map[string]struct {
		responseJSON string
		errorStr     string
	}{
		"JSON objects not support yet": {
			responseJSON: `{"attr": {"key": "value"}}`,
			errorStr:     "unmarshal not supported yet for type map[string]interface {} for field attr",
		},
		"JSON arrays not supported yet": {
			responseJSON: `{"attr": [{"key": "value"}]}`,
			errorStr:     "unmarshal not supported yet for type []interface {} for field attr",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.ErrorContains(t, autogeneration.Unmarshal([]byte(tc.responseJSON), &model), tc.errorStr)
		})
	}
}
