package autogen_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveUnknowns(t *testing.T) {
	ctx := context.Background()

	type modelEmptyTest struct{}

	type modelCustomTypeTest struct {
		AttrKnownString   types.String                            `tfsdk:"attr_known_string"`
		AttrUnknownObject customtypes.ObjectValue[modelEmptyTest] `tfsdk:"attr_unknown_object"`
		AttrMANYUpper     types.Int64                             `tfsdk:"attr_many_upper"`
	}

	type modelst struct {
		AttrStringUnknown           types.String                                     `tfsdk:"attr_string_unknown"`
		AttrCustomObjectUnknown     customtypes.ObjectValue[modelEmptyTest]          `tfsdk:"attr_custom_object_unknown"`
		AttrCustomObject            customtypes.ObjectValue[modelCustomTypeTest]     `tfsdk:"attr_custom_object"`
		AttrCustomListUnknown       customtypes.ListValue[types.String]              `tfsdk:"attr_custom_list_string"`
		AttrCustomSetUnknown        customtypes.SetValue[types.String]               `tfsdk:"attr_custom_set_string"`
		AttrCustomMapUnknown        customtypes.MapValue[types.String]               `tfsdk:"attr_custom_map_string"`
		AttrCustomNestedListUnknown customtypes.NestedListValue[modelEmptyTest]      `tfsdk:"attr_custom_nested_list_unknown"`
		AttrCustomNestedSetUnknown  customtypes.NestedSetValue[modelEmptyTest]       `tfsdk:"attr_custom_nested_set_unknown"`
		AttrCustomNestedMapUnknown  customtypes.NestedMapValue[modelEmptyTest]       `tfsdk:"attr_custom_nested_map_unknown"`
		AttrCustomNestedList        customtypes.NestedListValue[modelCustomTypeTest] `tfsdk:"attr_custom_nested_list"`
		AttrCustomNestedMap         customtypes.NestedMapValue[modelCustomTypeTest]  `tfsdk:"attr_custom_nested_map"`
	}

	model := modelst{
		AttrStringUnknown:       types.StringUnknown(),
		AttrCustomObjectUnknown: customtypes.NewObjectValueUnknown[modelEmptyTest](ctx),
		AttrCustomObject: customtypes.NewObjectValue[modelCustomTypeTest](ctx, modelCustomTypeTest{
			AttrKnownString:   types.StringValue("val1"),
			AttrUnknownObject: customtypes.NewObjectValueUnknown[modelEmptyTest](ctx),
			AttrMANYUpper:     types.Int64Unknown(),
		}),
		AttrCustomNestedListUnknown: customtypes.NewNestedListValueUnknown[modelEmptyTest](ctx),
		AttrCustomNestedSetUnknown:  customtypes.NewNestedSetValueUnknown[modelEmptyTest](ctx),
		AttrCustomNestedMapUnknown:  customtypes.NewNestedMapValueUnknown[modelEmptyTest](ctx),
		AttrCustomListUnknown:       customtypes.NewListValueUnknown[types.String](ctx),
		AttrCustomSetUnknown:        customtypes.NewSetValueUnknown[types.String](ctx),
		AttrCustomMapUnknown:        customtypes.NewMapValueUnknown[types.String](ctx),
		AttrCustomNestedList: customtypes.NewNestedListValue[modelCustomTypeTest](ctx, []modelCustomTypeTest{
			{
				AttrKnownString:   types.StringValue("val1"),
				AttrUnknownObject: customtypes.NewObjectValueUnknown[modelEmptyTest](ctx),
				AttrMANYUpper:     types.Int64Unknown(),
			},
			{
				AttrKnownString:   types.StringUnknown(),
				AttrUnknownObject: customtypes.NewObjectValue[modelEmptyTest](ctx, modelEmptyTest{}),
				AttrMANYUpper:     types.Int64Value(2),
			},
		}),
		AttrCustomNestedMap: customtypes.NewNestedMapValue[modelCustomTypeTest](ctx, map[string]modelCustomTypeTest{
			"keyOne": {
				AttrKnownString:   types.StringValue("val1"),
				AttrUnknownObject: customtypes.NewObjectValueUnknown[modelEmptyTest](ctx),
				AttrMANYUpper:     types.Int64Unknown(),
			},
			"keyTwo": {
				AttrKnownString:   types.StringUnknown(),
				AttrUnknownObject: customtypes.NewObjectValue[modelEmptyTest](ctx, modelEmptyTest{}),
				AttrMANYUpper:     types.Int64Value(2),
			},
		}),
	}
	modelExpected := modelst{
		AttrStringUnknown:       types.StringNull(),
		AttrCustomObjectUnknown: customtypes.NewObjectValueNull[modelEmptyTest](ctx),
		AttrCustomObject: customtypes.NewObjectValue[modelCustomTypeTest](ctx, modelCustomTypeTest{
			AttrKnownString:   types.StringValue("val1"),
			AttrUnknownObject: customtypes.NewObjectValueNull[modelEmptyTest](ctx),
			AttrMANYUpper:     types.Int64Null(),
		}),
		AttrCustomListUnknown:       customtypes.NewListValueNull[types.String](ctx),
		AttrCustomSetUnknown:        customtypes.NewSetValueNull[types.String](ctx),
		AttrCustomMapUnknown:        customtypes.NewMapValueNull[types.String](ctx),
		AttrCustomNestedListUnknown: customtypes.NewNestedListValueNull[modelEmptyTest](ctx),
		AttrCustomNestedSetUnknown:  customtypes.NewNestedSetValueNull[modelEmptyTest](ctx),
		AttrCustomNestedMapUnknown:  customtypes.NewNestedMapValueNull[modelEmptyTest](ctx),
		AttrCustomNestedList: customtypes.NewNestedListValue[modelCustomTypeTest](ctx, []modelCustomTypeTest{
			{
				AttrKnownString:   types.StringValue("val1"),
				AttrUnknownObject: customtypes.NewObjectValueNull[modelEmptyTest](ctx),
				AttrMANYUpper:     types.Int64Null(),
			},
			{
				AttrKnownString:   types.StringNull(),
				AttrUnknownObject: customtypes.NewObjectValue[modelEmptyTest](ctx, modelEmptyTest{}),
				AttrMANYUpper:     types.Int64Value(2),
			},
		}),
		AttrCustomNestedMap: customtypes.NewNestedMapValue[modelCustomTypeTest](ctx, map[string]modelCustomTypeTest{
			"keyOne": {
				AttrKnownString:   types.StringValue("val1"),
				AttrUnknownObject: customtypes.NewObjectValueNull[modelEmptyTest](ctx),
				AttrMANYUpper:     types.Int64Null(),
			},
			"keyTwo": {
				AttrKnownString:   types.StringNull(),
				AttrUnknownObject: customtypes.NewObjectValue[modelEmptyTest](ctx, modelEmptyTest{}),
				AttrMANYUpper:     types.Int64Value(2),
			},
		}),
	}
	require.NoError(t, autogen.ResolveUnknowns(&model))
	assert.Equal(t, modelExpected, model)
}
