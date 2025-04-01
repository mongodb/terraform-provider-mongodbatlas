package autogeneration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogeneration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalBasic(t *testing.T) {
	var model struct {
		AttributeFloat        types.Float64 `tfsdk:"attribute_float"`
		AttributeFloatWithInt types.Float64 `tfsdk:"attribute_float_with_int"`
		AttributeString       types.String  `tfsdk:"attribute_string"`
		AttributeNotInJSON    types.String  `tfsdk:"attribute_not_in_json"`
		AttributeInt          types.Int64   `tfsdk:"attribute_int"`
		AttributeIntWithFloat types.Int64   `tfsdk:"attribute_int_with_float"`
		AttributeTrue         types.Bool    `tfsdk:"attribute_true"`
		AttributeFalse        types.Bool    `tfsdk:"attribute_false"`
	}
	const (
		epsilon = 10e-15 // float tolerance
		// attribute_not_in_model is ignored because it is not in the model, no error is thrown.
		// attribute_null is ignored because it is null, no error is thrown even if it is not in the model.
		tfResponseJSON = `
			{
				"attribute_string": "value_string",
				"attribute_true": true,
				"attribute_false": false,
				"attribute_int": 123,
				"attribute_int_with_float": 10.6,
				"attribute_float": 456.1,
				"attribute_float_with_int": 13,
				"attribute_not_in_model": "val",
				"attribute_null": null
			}
		`
	)
	require.NoError(t, autogeneration.Unmarshal([]byte(tfResponseJSON), &model))
	assert.Equal(t, "value_string", model.AttributeString.ValueString())
	assert.True(t, model.AttributeTrue.ValueBool())
	assert.False(t, model.AttributeFalse.ValueBool())
	assert.Equal(t, int64(123), model.AttributeInt.ValueInt64())
	assert.Equal(t, int64(10), model.AttributeIntWithFloat.ValueInt64()) // response floats stored in model ints have their decimals stripped.
	assert.InEpsilon(t, float64(456.1), model.AttributeFloat.ValueFloat64(), epsilon)
	assert.InEpsilon(t, float64(13), model.AttributeFloatWithInt.ValueFloat64(), epsilon)
	assert.True(t, model.AttributeNotInJSON.IsNull()) // attributes not in JSON response are not changed, so null is kept.
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
			errorStr:     "not supported yet type map[string]interface {} for field att",
		},
		"JSON arrays not supported yet": {
			responseJSON: `{"attr": [{"key": "value"}]}`,
			errorStr:     "not supported yet type []interface {} for field att",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.ErrorContains(t, autogeneration.Unmarshal([]byte(tc.responseJSON), &model), tc.errorStr)
		})
	}
}
