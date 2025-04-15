package autogen_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveUnknowns(t *testing.T) {
	type modelst struct {
		AttrString types.String `tfsdk:"attr_string"`
		AttrObject types.Object `tfsdk:"attr_object"`
		AttrList   types.List   `tfsdk:"attr_list"`
	}
	model := modelst{
		AttrString: types.StringUnknown(),
		AttrObject: types.ObjectUnknown(objTypeTest.AttributeTypes()),
		AttrList:   types.ListUnknown(objTypeTest),
	}
	modelExpected := modelst{
		AttrString: types.StringNull(),
		AttrObject: types.ObjectNull(objTypeTest.AttributeTypes()),
		AttrList:   types.ListNull(objTypeTest),
	}
	require.NoError(t, autogen.ResolveUnknowns(&model))
	assert.Equal(t, modelExpected, model)
}
