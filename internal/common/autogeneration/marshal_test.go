package autogeneration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogeneration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	var model tfModelTest
	require.NoError(t, autogeneration.Unmarshal([]byte(tfResponseJSON), &model))
	assert.Equal(t, "value_string", model.AttributeString.ValueString())
	assert.True(t, model.AttributeTrue.ValueBool())
	assert.False(t, model.AttributeFalse.ValueBool())
}

type tfModelTest struct {
	AttributeString types.String `tfsdk:"attribute_string"`
	AttributeTrue   types.Bool   `tfsdk:"attribute_true"`
	AttributeFalse  types.Bool   `tfsdk:"attribute_false"`
}

const tfResponseJSON = `
{
	"attribute_string": "value_string",
	"attribute_true": true,
	"attribute_false": false,
	"attribute_not_in_model": "val"
}
`
