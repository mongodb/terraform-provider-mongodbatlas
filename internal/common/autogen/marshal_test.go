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

func TestMarshalBasic(t *testing.T) {
	model := struct {
		AttrFloat           types.Float64 `tfsdk:"attr_float"`
		AttrString          types.String  `tfsdk:"attr_string"`
		AttrOmit            types.String  `tfsdk:"attr_omit" autogen:"omitjson"`
		AttrUnknown         types.String  `tfsdk:"attr_unknown"`
		AttrNull            types.String  `tfsdk:"attr_null"`
		AttrOmitNoTerraform string        `autogen:"omitjson"`
		AttrInt             types.Int64   `tfsdk:"attr_int"`
		AttrBoolTrue        types.Bool    `tfsdk:"attr_bool_true"`
		AttrBoolFalse       types.Bool    `tfsdk:"attr_bool_false"`
		AttrBoolNull        types.Bool    `tfsdk:"attr_bool_null"`
		AttrMANYUpper       types.Int64   `tfsdk:"attr_many_upper"`
	}{
		AttrFloat:           types.Float64Value(1.234),
		AttrString:          types.StringValue("hello"),
		AttrOmit:            types.StringValue("omit"),
		AttrOmitNoTerraform: "omit",
		AttrUnknown:         types.StringUnknown(), // unknown values are not marshaled
		AttrNull:            types.StringNull(),    // null values are not marshaled
		AttrInt:             types.Int64Value(1),
		AttrBoolTrue:        types.BoolValue(true),
		AttrBoolFalse:       types.BoolValue(false),
		AttrBoolNull:        types.BoolNull(), // null values are not marshaled
		AttrMANYUpper:       types.Int64Value(2),
	}
	const expectedJSON = `
		{ 
			"attrString": "hello", 
			"attrInt": 1, 
			"attrFloat": 1.234, 
			"attrBoolTrue": true, 
			"attrBoolFalse": false,
			"attrMANYUpper": 2
		}
	`
	raw, err := autogen.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))
}

func TestMarshalDynamicJSONAttr(t *testing.T) {
	model := struct {
		AttrDynamicJSONObject  jsontypes.Normalized `tfsdk:"attr_dynamic_json_object"`
		AttrDynamicJSONBoolean jsontypes.Normalized `tfsdk:"attr_dynamic_json_boolean"`
		AttrDynamicJSONString  jsontypes.Normalized `tfsdk:"attr_dynamic_json_string"`
		AttrDynamicJSONNumber  jsontypes.Normalized `tfsdk:"attr_dynamic_json_number"`
		AttrDynamicJSONArray   jsontypes.Normalized `tfsdk:"attr_dynamic_json_array"`
	}{
		AttrDynamicJSONObject:  jsontypes.NewNormalizedValue("{\"hello\": \"there\"}"),
		AttrDynamicJSONBoolean: jsontypes.NewNormalizedValue("true"),
		AttrDynamicJSONString:  jsontypes.NewNormalizedValue("\"hello\""),
		AttrDynamicJSONNumber:  jsontypes.NewNormalizedValue("1.234"),
		AttrDynamicJSONArray:   jsontypes.NewNormalizedValue("[1, 2, 3]"),
	}
	const expectedJSON = `
		{ 
			"attrDynamicJSONObject": {"hello": "there"}, 
			"attrDynamicJSONBoolean": true, 
			"attrDynamicJSONString": "hello", 
			"attrDynamicJSONNumber": 1.234,
			"attrDynamicJSONArray": [1, 2, 3]
		}
	`
	raw, err := autogen.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))
}

func TestMarshalNestedAllTypes(t *testing.T) {
	model := struct {
		AttrCustomList customtypes.ListValue[types.String] `tfsdk:"attr_custom_list"`
		AttrCustomSet  customtypes.SetValue[types.String]  `tfsdk:"attr_custom_set"`
		AttrCustomMap  customtypes.MapValue[types.String]  `tfsdk:"attr_custom_map"`
	}{
		AttrCustomList: customtypes.NewListValue[types.String](t.Context(), []attr.Value{types.StringValue("val1"), types.StringValue("val2")}),
		AttrCustomSet:  customtypes.NewSetValue[types.String](t.Context(), []attr.Value{types.StringValue("val11"), types.StringValue("val22")}),
		AttrCustomMap: customtypes.NewMapValue[types.String](t.Context(), map[string]attr.Value{
			"keyOne": types.StringValue("val1"),
			"KeyTwo": types.StringValue("val2"),
		}),
	}
	const expectedJSON = `
		{
			"attrCustomList": ["val1", "val2"],
			"attrCustomSet": ["val11", "val22"],
			"attrCustomMap": {
				"keyOne": "val1",
				"KeyTwo": "val2"
			}
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
		AttrCustomList    customtypes.ListValue[types.String] `tfsdk:"attr_custom_list"`
		AttrCustomSet     customtypes.SetValue[types.String]  `tfsdk:"attr_custom_set"`
		AttrCustomMap     customtypes.MapValue[types.String]  `tfsdk:"attr_custom_map"`
		AttrString        types.String                        `tfsdk:"attr_string"`
		AttrIncludeString types.String                        `tfsdk:"attr_include_update" autogen:"includenullonupdate"`
	}{
		AttrCustomList:    customtypes.NewListValueNull[types.String](t.Context()),
		AttrCustomSet:     customtypes.NewSetValueNull[types.String](t.Context()),
		AttrCustomMap:     customtypes.NewMapValueNull[types.String](t.Context()),
		AttrString:        types.StringNull(),
		AttrIncludeString: types.StringNull(),
	}
	// null list and set root elements are sent as empty arrays in update.
	// fields with includenullonupdate tag are included even when null during updates.
	const expectedJSON = `
		{
			"attrCustomList": [],
			"attrCustomSet": [],
			"attrIncludeString": null
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

func TestMarshalCustomTypeNestedSet(t *testing.T) {
	ctx := context.Background()

	type modelEmptyTest struct{}

	type modelNestedObject struct {
		AttrNestedInt types.Int64 `tfsdk:"attr_nested_int"`
	}

	type modelNestedSetItem struct {
		AttrOmit       customtypes.NestedSetValue[modelEmptyTest] `tfsdk:"attr_omit" autogen:"omitjson"`
		AttrOmitUpdate customtypes.NestedSetValue[modelEmptyTest] `tfsdk:"attr_omit_update" autogen:"omitjsonupdate"`
		AttrPrimitive  types.String                               `tfsdk:"attr_primitive"`
		AttrObject     customtypes.ObjectValue[modelNestedObject] `tfsdk:"attr_object"`
		AttrMANYUpper  types.Int64                                `tfsdk:"attr_many_upper"`
	}

	model := struct {
		AttrNestedSet      customtypes.NestedSetValue[modelNestedSetItem] `tfsdk:"attr_nested_set"`
		AttrNestedSetNull  customtypes.NestedSetValue[modelNestedSetItem] `tfsdk:"attr_nested_set_null"`
		AttrNestedSetEmpty customtypes.NestedSetValue[modelNestedSetItem] `tfsdk:"attr_nested_set_empty"`
	}{
		AttrNestedSet: customtypes.NewNestedSetValue[modelNestedSetItem](ctx, []modelNestedSetItem{
			{
				AttrPrimitive: types.StringValue("string1"),
				AttrMANYUpper: types.Int64Value(1),
				AttrObject: customtypes.NewObjectValue[modelNestedObject](ctx, modelNestedObject{
					AttrNestedInt: types.Int64Value(2),
				}),
				AttrOmit:       customtypes.NewNestedSetValue[modelEmptyTest](ctx, []modelEmptyTest{}),
				AttrOmitUpdate: customtypes.NewNestedSetValue[modelEmptyTest](ctx, []modelEmptyTest{}),
			},
			{
				AttrPrimitive: types.StringValue("string2"),
				AttrMANYUpper: types.Int64Value(3),
				AttrObject: customtypes.NewObjectValue[modelNestedObject](ctx, modelNestedObject{
					AttrNestedInt: types.Int64Value(4),
				}),
				AttrOmit:       customtypes.NewNestedSetValue[modelEmptyTest](ctx, []modelEmptyTest{}),
				AttrOmitUpdate: customtypes.NewNestedSetValue[modelEmptyTest](ctx, []modelEmptyTest{}),
			},
		}),
		AttrNestedSetNull:  customtypes.NewNestedSetValueNull[modelNestedSetItem](ctx),
		AttrNestedSetEmpty: customtypes.NewNestedSetValue[modelNestedSetItem](ctx, []modelNestedSetItem{}),
	}

	const expectedCreateJSON = `
		{
			"attrNestedSet": [
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
			"attrNestedSetEmpty": []
		}
	`
	rawCreate, err := autogen.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedCreateJSON, string(rawCreate))

	const expectedUpdateJSON = `
		{
			"attrNestedSet": [
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
			"attrNestedSetNull": [],
			"attrNestedSetEmpty": []
		}
	`
	rawUpdate, err := autogen.Marshal(&model, true)
	require.NoError(t, err)
	assert.JSONEq(t, expectedUpdateJSON, string(rawUpdate))
}

func TestMarshalCustomTypeNestedMap(t *testing.T) {
	ctx := context.Background()

	type modelEmptyTest struct{}

	type modelNestedObject struct {
		AttrNestedInt types.Int64 `tfsdk:"attr_nested_int"`
	}

	type modelNestedMapItem struct {
		AttrOmit       customtypes.NestedMapValue[modelEmptyTest] `tfsdk:"attr_omit" autogen:"omitjson"`
		AttrOmitUpdate customtypes.NestedMapValue[modelEmptyTest] `tfsdk:"attr_omit_update" autogen:"omitjsonupdate"`
		AttrPrimitive  types.String                               `tfsdk:"attr_primitive"`
		AttrObject     customtypes.ObjectValue[modelNestedObject] `tfsdk:"attr_object"`
		AttrMANYUpper  types.Int64                                `tfsdk:"attr_many_upper"`
	}

	model := struct {
		AttrNestedMap      customtypes.NestedMapValue[modelNestedMapItem] `tfsdk:"attr_nested_map"`
		AttrNestedMapNull  customtypes.NestedMapValue[modelNestedMapItem] `tfsdk:"attr_nested_map_null"`
		AttrNestedMapEmpty customtypes.NestedMapValue[modelNestedMapItem] `tfsdk:"attr_nested_map_empty"`
	}{
		AttrNestedMap: customtypes.NewNestedMapValue[modelNestedMapItem](ctx, map[string]modelNestedMapItem{
			"keyOne": {
				AttrPrimitive: types.StringValue("string1"),
				AttrMANYUpper: types.Int64Value(1),
				AttrObject: customtypes.NewObjectValue[modelNestedObject](ctx, modelNestedObject{
					AttrNestedInt: types.Int64Value(2),
				}),
				AttrOmit:       customtypes.NewNestedMapValue[modelEmptyTest](ctx, map[string]modelEmptyTest{}),
				AttrOmitUpdate: customtypes.NewNestedMapValue[modelEmptyTest](ctx, map[string]modelEmptyTest{}),
			},
			"KeyTwo": {
				AttrPrimitive: types.StringValue("string2"),
				AttrMANYUpper: types.Int64Value(3),
				AttrObject: customtypes.NewObjectValue[modelNestedObject](ctx, modelNestedObject{
					AttrNestedInt: types.Int64Value(4),
				}),
				AttrOmit:       customtypes.NewNestedMapValue[modelEmptyTest](ctx, map[string]modelEmptyTest{}),
				AttrOmitUpdate: customtypes.NewNestedMapValue[modelEmptyTest](ctx, map[string]modelEmptyTest{}),
			},
		}),
		AttrNestedMapNull:  customtypes.NewNestedMapValueNull[modelNestedMapItem](ctx),
		AttrNestedMapEmpty: customtypes.NewNestedMapValue[modelNestedMapItem](ctx, map[string]modelNestedMapItem{}),
	}

	const expectedCreateJSON = `
		{
			"attrNestedMap": {
				"keyOne": {
					"attrPrimitive": "string1",
					"attrMANYUpper": 1,
					"attrObject": {
						"attrNestedInt": 2
					},
					"attrOmitUpdate": {}
				},
				"KeyTwo": {
					"attrPrimitive": "string2",
					"attrMANYUpper": 3,
					"attrObject": {
						"attrNestedInt": 4
					},
					"attrOmitUpdate": {}
				}
			},
			"attrNestedMapEmpty": {}
		}
	`
	rawCreate, err := autogen.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedCreateJSON, string(rawCreate))

	const expectedUpdateJSON = `
		{
			"attrNestedMap": {
				"keyOne": {
					"attrPrimitive": "string1",
					"attrMANYUpper": 1,
					"attrObject": {
						"attrNestedInt": 2
					}
				},
				"KeyTwo": {
					"attrPrimitive": "string2",
					"attrMANYUpper": 3,
					"attrObject": {
						"attrNestedInt": 4
					}
				}
			},
			"attrNestedMapEmpty": {}
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
