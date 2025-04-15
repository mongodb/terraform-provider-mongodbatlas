package autogen_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareResponseModel(t *testing.T) {
	type modelst struct {
		AttrStringUnknown types.String `tfsdk:"attr_string_unknown"`
		AttrObjectUnknown types.Object `tfsdk:"attr_object_unknown"`
		AttrListUnknown   types.List   `tfsdk:"attr_list_unknown"`
		AttrListEmpty     types.List   `tfsdk:"attr_list_empty"`
		AttrObject        types.Object `tfsdk:"attr_object"`
	}
	model := modelst{
		AttrStringUnknown: types.StringUnknown(),
		AttrObjectUnknown: types.ObjectUnknown(objTypeTest.AttributeTypes()),
		AttrListUnknown:   types.ListUnknown(objTypeTest),
		AttrListEmpty:     types.ListValueMust(objTypeTest, []attr.Value{}),
		AttrObject: types.ObjectValueMust(objTypeTest.AttributeTypes(), map[string]attr.Value{
			"attr_string": types.StringUnknown(),
			"attr_float":  types.Float64Unknown(),
			"attr_int":    types.Int64Unknown(),
			"attr_bool":   types.BoolUnknown(),
		}),
	}
	modelExpected := modelst{
		AttrStringUnknown: types.StringNull(),
		AttrObjectUnknown: types.ObjectNull(objTypeTest.AttributeTypes()),
		AttrListUnknown:   types.ListNull(objTypeTest),
		AttrListEmpty:     types.ListNull(objTypeTest),
		AttrObject: types.ObjectValueMust(objTypeTest.AttributeTypes(), map[string]attr.Value{
			"attr_string": types.StringNull(),
			"attr_float":  types.Float64Null(),
			"attr_int":    types.Int64Null(),
			"attr_bool":   types.BoolNull(),
		}),
	}
	require.NoError(t, autogen.PrepareResponseModel(&model))
	assert.Equal(t, modelExpected, model)
}
