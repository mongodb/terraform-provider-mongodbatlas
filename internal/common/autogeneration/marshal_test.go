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
		AttributeInt          types.Int64   `tfsdk:"attribute_int"`
		AttributeIntWithFloat types.Int64   `tfsdk:"attribute_int_with_float"`
		AttributeTrue         types.Bool    `tfsdk:"attribute_true"`
		AttributeFalse        types.Bool    `tfsdk:"attribute_false"`
	}
	const (
		epsilon = 10e-15 // float tolerance
		// attribute_not_in_model is ignored but not an error
		tfResponseJSON = `
			{
				"attribute_string": "value_string",
				"attribute_true": true,
				"attribute_false": false,
				"attribute_int": 123,
				"attribute_int_with_float": 10.6,
				"attribute_float": 456.1,
				"attribute_float_with_int": 13,
				"attribute_not_in_model": "val"
			}
		`
	)
	require.NoError(t, autogeneration.Unmarshal([]byte(tfResponseJSON), &model))
	assert.Equal(t, "value_string", model.AttributeString.ValueString())
	assert.True(t, model.AttributeTrue.ValueBool())
	assert.False(t, model.AttributeFalse.ValueBool())
	assert.Equal(t, int64(123), model.AttributeInt.ValueInt64())
	assert.Equal(t, int64(10), model.AttributeIntWithFloat.ValueInt64())
	assert.InEpsilon(t, float64(456.1), model.AttributeFloat.ValueFloat64(), epsilon)
	assert.InEpsilon(t, float64(13), model.AttributeFloatWithInt.ValueFloat64(), epsilon)
}
