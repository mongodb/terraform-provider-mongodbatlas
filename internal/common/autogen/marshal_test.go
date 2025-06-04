package autogen_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
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
		AttrOmit            types.String `tfsdk:"attr_omit" autogen:"omitjson"`
		AttrOmitNoTerraform string       `autogen:"omitjson"`
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
	raw, err := autogen.Marshal(&model, false)
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
	raw, err := autogen.Marshal(&model, false)
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
	raw, err := autogen.Marshal(&model, false)
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
		AttrOmitUpdate types.String `tfsdk:"attr_omit_update" autogen:"omitjsonupdate"`
		AttrOmit       types.String `tfsdk:"attr_omit" autogen:"omitjson"`
	}{
		Attr:           types.StringValue("val1"),
		AttrOmitUpdate: types.StringValue("val2"),
		AttrOmit:       types.StringValue("omit"),
	}
	create, errCreate := autogen.Marshal(&model, false)
	require.NoError(t, errCreate)
	assert.JSONEq(t, expectedCreate, string(create))

	update, errUpdate := autogen.Marshal(&model, true)
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
			raw, err := autogen.Marshal(model, false)
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
				_, _ = autogen.Marshal(model, false)
			})
		})
	}
}


func TestIsDiscriminatorTag(t *testing.T) {
	testCases := map[string]struct {
		tagValue string
		expected *autogen.DiscriminatorTag
	}{
		"empty tag": {
			tagValue: "",
			expected: nil,
		},
		"valid tag with name": {
			tagValue: "discriminator:type=Cluster",
			expected: &autogen.DiscriminatorTag{
				DiscriminatorPropName: "type",
				DiscriminatorPropValue:    "Cluster",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tag := autogen.IsDiscriminatorTag(tc.tagValue)
			if tc.expected == nil {
				assert.Nil(t, tag)
			} else {
				require.NotNil(t, tag)
				assert.Equal(t, tc.expected.DiscriminatorPropName, tag.DiscriminatorPropName)
				assert.Equal(t, tc.expected.DiscriminatorPropValue, tag.DiscriminatorPropValue)
			}
		})
	}
}

func TestMarshalDiscriminator(t *testing.T) {
	testCases := map[string]struct {
		model streamConn
		expectedJSON string
	}{
		"cluster model": {
			model: streamConnModelCluster,
			expectedJSON: jsonRespCluster,
		},
		"https model": {
			model: streamConnModelHttps,
			expectedJSON: jsonRespHttps,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			raw, err := autogen.Marshal(&tc.model, false)
			require.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(raw))
		})
	}
}