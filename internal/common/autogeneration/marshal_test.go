package autogeneration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogeneration"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	var model tfModelTest
	require.NoError(t, autogeneration.Unmarshal([]byte(tfResponseJSON), &model))
	require.Equal(t, "value_string", model.AttributeString.ValueString())
}

type tfModelTest struct {
	AttributeString types.String `tfsdk:"attribute_string"`
}

const tfResponseJSON = `
{
	"attribute_string": "value_string",
	"attribute_not_in_model": "val"
}
`
