package autogen_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
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
		AttrFloat           types.Float64        `tfsdk:"attr_float"`
		AttrString          types.String         `tfsdk:"attr_string"`
		AttrOmit            types.String         `tfsdk:"attr_omit" autogen:"omitjson"`
		AttrUnkown          types.String         `tfsdk:"attr_unknown"`
		AttrNull            types.String         `tfsdk:"attr_null"`
		AttrJSON            jsontypes.Normalized `tfsdk:"attr_json"`
		AttrOmitNoTerraform string               `autogen:"omitjson"`
		AttrInt             types.Int64          `tfsdk:"attr_int"`
		AttrBoolTrue        types.Bool           `tfsdk:"attr_bool_true"`
		AttrBoolFalse       types.Bool           `tfsdk:"attr_bool_false"`
		AttrBoolNull        types.Bool           `tfsdk:"attr_bool_null"`
		AttrMANYUpper       types.Int64          `tfsdk:"attr_many_upper"`
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
		AttrJSON:            jsontypes.NewNormalizedValue("{\"hello\": \"there\"}"),
		AttrMANYUpper:       types.Int64Value(2),
	}
	const expectedJSON = `
		{ 
			"attrString": "hello", 
			"attrInt": 1, 
			"attrFloat": 1.234, 
			"attrBoolTrue": true, 
			"attrBoolFalse": false, 
			"attrJSON": {"hello": "there"}, 
			"attrMANYUpper": 2
		}
	`
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

func TestMarshalUpdateNull(t *testing.T) {
	model := struct {
		AttrList          types.List   `tfsdk:"attr_list"`
		AttrSet           types.Set    `tfsdk:"attr_set"`
		AttrString        types.String `tfsdk:"attr_string"`
		AttrObj           types.Object `tfsdk:"attr_obj"`
		AttrIncludeString types.String `tfsdk:"attr_include_update" autogen:"includenullonupdate"`
		AttrIncludeObj    types.Object `tfsdk:"attr_include_obj" autogen:"includenullonupdate"`
	}{
		AttrList:          types.ListNull(types.StringType),
		AttrSet:           types.SetNull(types.StringType),
		AttrString:        types.StringNull(),
		AttrObj:           types.ObjectNull(objTypeTest.AttrTypes),
		AttrIncludeString: types.StringNull(),
		AttrIncludeObj:    types.ObjectNull(objTypeTest.AttrTypes),
	}
	// null list and set root elements are sent as empty arrays in update.
	// fields with includenullonupdate tag are included even when null during updates.
	const expectedJSON = `
		{
			"attrList": [],
			"attrSet": [],
			"attrIncludeString": null,
			"attrIncludeObj": null
		}
	`
	raw, err := autogen.Marshal(&model, true)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))

	// Test that includenullonupdate fields are NOT included when isUpdate is false
	rawCreate, errCreate := autogen.Marshal(&model, false)
	require.NoError(t, errCreate)
	const expectedJSONCreate = `
		{
		}
	`
	assert.JSONEq(t, expectedJSONCreate, string(rawCreate))
}

func TestMarshalCustomTypeObject(t *testing.T) {
	ctx := context.Background()

	type modelEmptyTest struct{}

	type modelCustomTypeTest struct {
		AttrPrimitiveOmit    types.String                            `tfsdk:"attr_primitive_omit" autogen:"omitjson"`
		AttrObjectOmit       customtypes.ObjectValue[modelEmptyTest] `tfsdk:"attr_object_omit" autogen:"omitjson"`
		AttrObjectOmitUpdate customtypes.ObjectValue[modelEmptyTest] `tfsdk:"attr_object_omit_update" autogen:"omitjsonupdate"`
		AttrNull             customtypes.ObjectValue[modelEmptyTest] `tfsdk:"attr_null" autogen:"includenullonupdate"`
		AttrInt              types.Int64                             `tfsdk:"attr_int"`
		AttrMANYUpper        types.Int64                             `tfsdk:"attr_many_upper"`
	}

	type modelCustomTypeParentTest struct {
		AttrString types.String                                 `tfsdk:"attr_string"`
		AttrObject customtypes.ObjectValue[modelCustomTypeTest] `tfsdk:"attr_object"`
	}

	nullObject := customtypes.NewObjectValueNull[modelEmptyTest](ctx)
	emptyObject := customtypes.NewObjectValue[modelEmptyTest](ctx, modelEmptyTest{})

	model := struct {
		AttrObjectBasic  customtypes.ObjectValue[modelCustomTypeTest]       `tfsdk:"attr_object_basic"`
		AttrObjectNull   customtypes.ObjectValue[modelCustomTypeTest]       `tfsdk:"attr_object_null"`
		AttrObjectNested customtypes.ObjectValue[modelCustomTypeParentTest] `tfsdk:"attr_object_nested"`
	}{
		AttrObjectBasic: customtypes.NewObjectValue[modelCustomTypeTest](ctx, modelCustomTypeTest{
			AttrInt:              types.Int64Value(1),
			AttrPrimitiveOmit:    types.StringValue("omitted"),
			AttrObjectOmit:       emptyObject,
			AttrObjectOmitUpdate: emptyObject,
			AttrNull:             nullObject,
			AttrMANYUpper:        types.Int64Value(2),
		}),
		AttrObjectNull: customtypes.NewObjectValueNull[modelCustomTypeTest](ctx),
		AttrObjectNested: customtypes.NewObjectValue[modelCustomTypeParentTest](ctx, modelCustomTypeParentTest{
			AttrString: types.StringValue("parent"),
			AttrObject: customtypes.NewObjectValue[modelCustomTypeTest](ctx, modelCustomTypeTest{
				AttrInt:              types.Int64Value(2),
				AttrPrimitiveOmit:    types.StringValue("omitted"),
				AttrObjectOmit:       emptyObject,
				AttrObjectOmitUpdate: emptyObject,
				AttrNull:             nullObject,
				AttrMANYUpper:        types.Int64Value(3),
			}),
		}),
	}

	const expectedCreateJSON = `
		{
			"attrObjectBasic": {
				"attrInt": 1,
				"attrObjectOmitUpdate": {},
				"attrMANYUpper": 2
			},
			"attrObjectNested": {
				"attrObject": {
					"attrInt": 2,
					"attrObjectOmitUpdate": {},
					"attrMANYUpper": 3
				},
				"attrString": "parent"
			}
		}
	`
	rawCreate, err := autogen.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedCreateJSON, string(rawCreate))

	const expectedUpdateJSON = `
		{
			"attrObjectBasic": {
				"attrInt": 1,
				"attrNull": null,
				"attrMANYUpper": 2
			},
			"attrObjectNested": {
				"attrObject": {
					"attrInt": 2,
					"attrNull": null,
					"attrMANYUpper": 3
				},
				"attrString": "parent"
			}
		}
	`
	rawUpdate, err := autogen.Marshal(&model, true)
	require.NoError(t, err)
	assert.JSONEq(t, expectedUpdateJSON, string(rawUpdate))
}

func TestMarshalCustomTypeNestedList(t *testing.T) {
	ctx := context.Background()

	type modelEmptyTest struct{}

	type modelNestedObject struct {
		AttrNestedInt types.Int64 `tfsdk:"attr_nested_int"`
	}

	type modelNestedListItem struct {
		AttrOmit       customtypes.NestedListValue[modelEmptyTest] `tfsdk:"attr_omit" autogen:"omitjson"`
		AttrOmitUpdate customtypes.NestedListValue[modelEmptyTest] `tfsdk:"attr_omit_update" autogen:"omitjsonupdate"`
		AttrPrimitive  types.String                                `tfsdk:"attr_primitive"`
		AttrObject     customtypes.ObjectValue[modelNestedObject]  `tfsdk:"attr_object"`
		AttrMANYUpper  types.Int64                                 `tfsdk:"attr_many_upper"`
	}

	model := struct {
		AttrNestedList      customtypes.NestedListValue[modelNestedListItem] `tfsdk:"attr_nested_list"`
		AttrNestedListNull  customtypes.NestedListValue[modelNestedListItem] `tfsdk:"attr_nested_list_null"`
		AttrNestedListEmpty customtypes.NestedListValue[modelNestedListItem] `tfsdk:"attr_nested_list_empty"`
	}{
		AttrNestedList: customtypes.NewNestedListValue[modelNestedListItem](ctx, []modelNestedListItem{
			{
				AttrPrimitive: types.StringValue("string1"),
				AttrMANYUpper: types.Int64Value(1),
				AttrObject: customtypes.NewObjectValue[modelNestedObject](ctx, modelNestedObject{
					AttrNestedInt: types.Int64Value(2),
				}),
				AttrOmit:       customtypes.NewNestedListValue[modelEmptyTest](ctx, []modelEmptyTest{}),
				AttrOmitUpdate: customtypes.NewNestedListValue[modelEmptyTest](ctx, []modelEmptyTest{}),
			},
			{
				AttrPrimitive: types.StringValue("string2"),
				AttrMANYUpper: types.Int64Value(3),
				AttrObject: customtypes.NewObjectValue[modelNestedObject](ctx, modelNestedObject{
					AttrNestedInt: types.Int64Value(4),
				}),
				AttrOmit:       customtypes.NewNestedListValue[modelEmptyTest](ctx, []modelEmptyTest{}),
				AttrOmitUpdate: customtypes.NewNestedListValue[modelEmptyTest](ctx, []modelEmptyTest{}),
			},
		}),
		AttrNestedListNull:  customtypes.NewNestedListValueNull[modelNestedListItem](ctx),
		AttrNestedListEmpty: customtypes.NewNestedListValue[modelNestedListItem](ctx, []modelNestedListItem{}),
	}

	const expectedCreateJSON = `
		{
			"attrNestedList": [
				{
					"attrPrimitive": "string1",
					"attrMANYUpper": 1,
					"attrObject": {
						"attrNestedInt": 2
					},
					"attrOmitUpdate": []
				},
				{
					"attrPrimitive": "string2",
					"attrMANYUpper": 3,
					"attrObject": {
						"attrNestedInt": 4
					},
					"attrOmitUpdate": []
				}
			],
			"attrNestedListEmpty": []
		}
	`
	rawCreate, err := autogen.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedCreateJSON, string(rawCreate))

	const expectedUpdateJSON = `
		{
			"attrNestedList": [
				{
					"attrPrimitive": "string1",
					"attrMANYUpper": 1,
					"attrObject": {
						"attrNestedInt": 2
					}
				},
				{
					"attrPrimitive": "string2",
					"attrMANYUpper": 3,
					"attrObject": {
						"attrNestedInt": 4
					}
				}
			],
			"attrNestedListNull": [],
			"attrNestedListEmpty": []
		}
	`
	rawUpdate, err := autogen.Marshal(&model, true)
	require.NoError(t, err)
	assert.JSONEq(t, expectedUpdateJSON, string(rawUpdate))
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
