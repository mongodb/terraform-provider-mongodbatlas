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

func TestMarshalWithApiNameTag(t *testing.T) {
	// Test that apiname tag is used for JSON field name instead of struct field name
	model := struct {
		ProjectID types.String `tfsdk:"project_id" apiname:"groupId" autogen:"omitjson"`
		Name      types.String `tfsdk:"name" apiname:"clusterName"`
		RegularID types.String `tfsdk:"regular_id"` // No apiname tag, uses field name
	}{
		ProjectID: types.StringValue("proj123"),
		Name:      types.StringValue("my-cluster"),
		RegularID: types.StringValue("reg456"),
	}
	const expectedJSON = `
		{
			"clusterName": "my-cluster",
			"regularID": "reg456"
		}
	`
	// Note: ProjectID is omitted due to omitjson tag
	raw, err := autogen.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))
}

func TestMarshalDynamicJSONAttr(t *testing.T) {
	model := struct {
		AttrDynamicJSONObject         jsontypes.Normalized                        `tfsdk:"attr_dynamic_json_object"`
		AttrDynamicJSONBoolean        jsontypes.Normalized                        `tfsdk:"attr_dynamic_json_boolean"`
		AttrDynamicJSONString         jsontypes.Normalized                        `tfsdk:"attr_dynamic_json_string"`
		AttrDynamicJSONNumber         jsontypes.Normalized                        `tfsdk:"attr_dynamic_json_number"`
		AttrDynamicJSONArray          jsontypes.Normalized                        `tfsdk:"attr_dynamic_json_array"`
		AttrListOfDynamicJSONObjects  customtypes.ListValue[jsontypes.Normalized] `tfsdk:"attr_list_of_dynamic_json_objects"`
		AttrSetOfDynamicJSONObjects   customtypes.SetValue[jsontypes.Normalized]  `tfsdk:"attr_set_of_dynamic_json_objects"`
		AttrMapOfDynamicJSONObjects   customtypes.MapValue[jsontypes.Normalized]  `tfsdk:"attr_map_of_dynamic_json_objects"`
		AttrListOfDynamicJSONBooleans customtypes.ListValue[jsontypes.Normalized] `tfsdk:"attr_list_of_dynamic_json_booleans"`
		AttrSetOfDynamicJSONBooleans  customtypes.SetValue[jsontypes.Normalized]  `tfsdk:"attr_set_of_dynamic_json_booleans"`
		AttrMapOfDynamicJSONBooleans  customtypes.MapValue[jsontypes.Normalized]  `tfsdk:"attr_map_of_dynamic_json_booleans"`
	}{
		AttrDynamicJSONObject:        jsontypes.NewNormalizedValue("{\"hello\": \"there\"}"),
		AttrDynamicJSONBoolean:       jsontypes.NewNormalizedValue("true"),
		AttrDynamicJSONString:        jsontypes.NewNormalizedValue("\"hello\""),
		AttrDynamicJSONNumber:        jsontypes.NewNormalizedValue("1.234"),
		AttrDynamicJSONArray:         jsontypes.NewNormalizedValue("[1, 2, 3]"),
		AttrListOfDynamicJSONObjects: customtypes.NewListValue[jsontypes.Normalized](t.Context(), []attr.Value{jsontypes.NewNormalizedValue("{\"hello\": \"there\"}")}),
		AttrSetOfDynamicJSONObjects:  customtypes.NewSetValue[jsontypes.Normalized](t.Context(), []attr.Value{jsontypes.NewNormalizedValue("{\"hello\": \"there\"}")}),
		AttrMapOfDynamicJSONObjects: customtypes.NewMapValue[jsontypes.Normalized](t.Context(), map[string]attr.Value{
			"key1": jsontypes.NewNormalizedValue("{\"hello\": \"there\"}"),
			"key2": jsontypes.NewNormalizedValue("{\"hello\": \"there\"}"),
		}),
		AttrListOfDynamicJSONBooleans: customtypes.NewListValue[jsontypes.Normalized](t.Context(), []attr.Value{jsontypes.NewNormalizedValue("true")}),
		AttrSetOfDynamicJSONBooleans:  customtypes.NewSetValue[jsontypes.Normalized](t.Context(), []attr.Value{jsontypes.NewNormalizedValue("true")}),
		AttrMapOfDynamicJSONBooleans: customtypes.NewMapValue[jsontypes.Normalized](t.Context(), map[string]attr.Value{
			"key1": jsontypes.NewNormalizedValue("true"),
			"key2": jsontypes.NewNormalizedValue("false"),
		}),
	}
	const expectedJSON = `
		{ 
			"attrDynamicJSONObject": {"hello": "there"}, 
			"attrDynamicJSONBoolean": true, 
			"attrDynamicJSONString": "hello", 
			"attrDynamicJSONNumber": 1.234,
			"attrDynamicJSONArray": [1, 2, 3],
			"attrListOfDynamicJSONObjects": [{"hello": "there"}],
			"attrSetOfDynamicJSONObjects": [{"hello": "there"}],
			"attrMapOfDynamicJSONObjects": {"key1": {"hello": "there"}, "key2": {"hello": "there"}},
			"attrListOfDynamicJSONBooleans": [true],
			"attrSetOfDynamicJSONBooleans": [true],
			"attrMapOfDynamicJSONBooleans": {"key1": true, "key2": false}
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

func TestMarshalUpdateAbsentAttrs(t *testing.T) {
	type modelEmptyTest struct{}

	model := struct {
		AttrList                  customtypes.ListValue[types.String]         `tfsdk:"attr_list"`
		AttrListSendNull          customtypes.ListValue[types.String]         `tfsdk:"attr_list_send_null" autogen:"sendnullasnullonupdate"`
		AttrListSendEmpty         customtypes.ListValue[types.String]         `tfsdk:"attr_list_send_empty" autogen:"sendnullasemptyonupdate"`
		AttrSet                   customtypes.SetValue[types.String]          `tfsdk:"attr_set"`
		AttrSetSendNull           customtypes.SetValue[types.String]          `tfsdk:"attr_set_send_null" autogen:"sendnullasnullonupdate"`
		AttrSetSendEmpty          customtypes.SetValue[types.String]          `tfsdk:"attr_set_send_empty" autogen:"sendnullasemptyonupdate"`
		AttrNestedList            customtypes.NestedListValue[modelEmptyTest] `tfsdk:"attr_nested_list"`
		AttrNestedListSendNull    customtypes.NestedListValue[modelEmptyTest] `tfsdk:"attr_nested_list_send_null" autogen:"sendnullasnullonupdate"`
		AttrNestedListSendEmpty   customtypes.NestedListValue[modelEmptyTest] `tfsdk:"attr_nested_list_send_empty" autogen:"sendnullasemptyonupdate"`
		AttrNestedSet             customtypes.NestedSetValue[modelEmptyTest]  `tfsdk:"attr_nested_set"`
		AttrNestedSetSendNull     customtypes.NestedSetValue[modelEmptyTest]  `tfsdk:"attr_nested_set_send_null" autogen:"sendnullasnullonupdate"`
		AttrNestedSetSendEmpty    customtypes.NestedSetValue[modelEmptyTest]  `tfsdk:"attr_nested_set_send_empty" autogen:"sendnullasemptyonupdate"`
		AttrMap                   customtypes.MapValue[types.String]          `tfsdk:"attr_map"`
		AttrMapSendNull           customtypes.MapValue[types.String]          `tfsdk:"attr_map_send_null" autogen:"sendnullasnullonupdate"`
		AttrMapSendEmpty          customtypes.MapValue[types.String]          `tfsdk:"attr_map_send_empty" autogen:"sendnullasemptyonupdate"`
		AttrNestedMap             customtypes.NestedMapValue[modelEmptyTest]  `tfsdk:"attr_nested_map"`
		AttrNestedMapSendNull     customtypes.NestedMapValue[modelEmptyTest]  `tfsdk:"attr_nested_map_send_null" autogen:"sendnullasnullonupdate"`
		AttrNestedMapSendEmpty    customtypes.NestedMapValue[modelEmptyTest]  `tfsdk:"attr_nested_map_send_empty" autogen:"sendnullasemptyonupdate"`
		AttrNestedObject          customtypes.ObjectValue[modelEmptyTest]     `tfsdk:"attr_nested_object"`
		AttrNestedObjectSendNull  customtypes.ObjectValue[modelEmptyTest]     `tfsdk:"attr_nested_object_send_null" autogen:"sendnullasnullonupdate"`
		AttrNestedObjectSendEmpty customtypes.ObjectValue[modelEmptyTest]     `tfsdk:"attr_nested_object_send_empty" autogen:"sendnullasemptyonupdate"`
		AttrString                types.String                                `tfsdk:"attr_string"`
		AttrStringSendNull        types.String                                `tfsdk:"attr_send_update" autogen:"sendnullasnullonupdate"`
		AttrStringSendEmpty       types.String                                `tfsdk:"attr_string_send_empty" autogen:"sendnullasemptyonupdate"`
		AttrInt                   types.Int64                                 `tfsdk:"attr_int"`
		AttrIntSendNull           types.Int64                                 `tfsdk:"attr_int_send_null" autogen:"sendnullasnullonupdate"`
		AttrIntSendEmpty          types.Int64                                 `tfsdk:"attr_int_send_empty" autogen:"sendnullasemptyonupdate"`
		AttrBool                  types.Bool                                  `tfsdk:"attr_bool"`
		AttrBoolSendNull          types.Bool                                  `tfsdk:"attr_bool_send_null" autogen:"sendnullasnullonupdate"`
		AttrBoolSendEmpty         types.Bool                                  `tfsdk:"attr_bool_send_empty" autogen:"sendnullasemptyonupdate"`
	}{
		AttrList:                  customtypes.NewListValueNull[types.String](t.Context()),
		AttrListSendNull:          customtypes.NewListValueNull[types.String](t.Context()),
		AttrListSendEmpty:         customtypes.NewListValueNull[types.String](t.Context()),
		AttrSet:                   customtypes.NewSetValueNull[types.String](t.Context()),
		AttrSetSendNull:           customtypes.NewSetValueNull[types.String](t.Context()),
		AttrSetSendEmpty:          customtypes.NewSetValueNull[types.String](t.Context()),
		AttrMap:                   customtypes.NewMapValueNull[types.String](t.Context()),
		AttrMapSendNull:           customtypes.NewMapValueNull[types.String](t.Context()),
		AttrMapSendEmpty:          customtypes.NewMapValueNull[types.String](t.Context()),
		AttrNestedObject:          customtypes.NewObjectValueNull[modelEmptyTest](t.Context()),
		AttrNestedList:            customtypes.NewNestedListValueNull[modelEmptyTest](t.Context()),
		AttrNestedListSendNull:    customtypes.NewNestedListValueNull[modelEmptyTest](t.Context()),
		AttrNestedListSendEmpty:   customtypes.NewNestedListValueNull[modelEmptyTest](t.Context()),
		AttrNestedSet:             customtypes.NewNestedSetValueNull[modelEmptyTest](t.Context()),
		AttrNestedSetSendNull:     customtypes.NewNestedSetValueNull[modelEmptyTest](t.Context()),
		AttrNestedSetSendEmpty:    customtypes.NewNestedSetValueNull[modelEmptyTest](t.Context()),
		AttrNestedMap:             customtypes.NewNestedMapValueNull[modelEmptyTest](t.Context()),
		AttrNestedMapSendNull:     customtypes.NewNestedMapValueNull[modelEmptyTest](t.Context()),
		AttrNestedMapSendEmpty:    customtypes.NewNestedMapValueNull[modelEmptyTest](t.Context()),
		AttrNestedObjectSendNull:  customtypes.NewObjectValueNull[modelEmptyTest](t.Context()),
		AttrNestedObjectSendEmpty: customtypes.NewObjectValueNull[modelEmptyTest](t.Context()),
		AttrString:                types.StringNull(),
		AttrBool:                  types.BoolNull(),
		AttrInt:                   types.Int64Null(),
		AttrStringSendNull:        types.StringNull(),
		AttrStringSendEmpty:       types.StringNull(),
		AttrIntSendNull:           types.Int64Null(),
		AttrIntSendEmpty:          types.Int64Null(),
		AttrBoolSendNull:          types.BoolNull(),
		AttrBoolSendEmpty:         types.BoolNull(),
	}
	// Default behavior: null values are not included in the update payload.
	// Fields with sendnullasnullonupdate tag are included as null during updates.
	// Fields with sendnullasemptyonupdate tag send empty values ([] for list/set, {} for map) during updates.
	const expectedJSON = `
		{
			"attrListSendNull": null,
			"attrListSendEmpty": [],
			"attrSetSendNull": null,
			"attrSetSendEmpty": [],
			"attrMapSendNull": null,
			"attrMapSendEmpty": {},
			"attrNestedMapSendNull": null,
			"attrNestedMapSendEmpty": {},
			"attrNestedListSendNull": null,
			"attrNestedListSendEmpty": [],
			"attrNestedSetSendNull": null,
			"attrNestedSetSendEmpty": [],
			"attrNestedObjectSendNull": null,
			"attrStringSendNull": null,
			"attrIntSendNull": null,
			"attrBoolSendNull": null
		}
	`
	raw, err := autogen.Marshal(&model, true)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))

	// Test that sendnullasnullonupdate and sendnullasemptyonupdate fields are NOT included when isUpdate is false
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
		AttrNull             customtypes.ObjectValue[modelEmptyTest] `tfsdk:"attr_null" autogen:"sendnullasnullonupdate"`
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

func TestMarshalListAsMap(t *testing.T) {
	model := struct {
		AttrListAsMapWithValues customtypes.MapValue[types.String] `tfsdk:"attr_list_as_map_with_values" autogen:"listasmap"`
		AttrListAsMapEmpty      customtypes.MapValue[types.String] `tfsdk:"attr_list_as_map_empty" autogen:"listasmap"`
		AttrListAsMapNull       customtypes.MapValue[types.String] `tfsdk:"attr_list_as_map_null" autogen:"listasmap"`
	}{
		AttrListAsMapWithValues: customtypes.NewMapValue[types.String](t.Context(), map[string]attr.Value{
			"key1": types.StringValue("val1"),
			"key2": types.StringValue("val2"),
		}),
		AttrListAsMapEmpty: customtypes.NewMapValue[types.String](t.Context(), map[string]attr.Value{}),
		AttrListAsMapNull:  customtypes.NewMapValueNull[types.String](t.Context()),
	}
	const expectedCreateJSON = `
		{ 
			"attrListAsMapWithValues": [
				{
					"key": "key1",
					"value": "val1"
				},
				{
					"key": "key2",
					"value": "val2"
				}
			],
			"attrListAsMapEmpty": []
		}
	`
	const expectedUpdateJSON = `
		{ 
			"attrListAsMapWithValues": [
				{
					"key": "key1",
					"value": "val1"
				},
				{
					"key": "key2",
					"value": "val2"
				}
			],
			"attrListAsMapEmpty": []
		}
	`
	rawCreate, err := autogen.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedCreateJSON, string(rawCreate))
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

func TestMarshalEmbeddedExpandedModel(t *testing.T) {
	type modelExpandedFields struct {
		ID types.String `tfsdk:"id" apiname:"id" autogen:"omitjson"`
	}
	type modelExpanded struct {
		modelExpandedFields
		ConnectionName types.String `tfsdk:"connection_name" apiname:"name"`
		Type           types.String `tfsdk:"type"`
	}
	model := modelExpanded{
		modelExpandedFields: modelExpandedFields{
			ID: types.StringValue("ws-123-conn"),
		},
		ConnectionName: types.StringValue("conn"),
		Type:           types.StringValue("Sample"),
	}
	const expectedJSON = `
	{
		"name": "conn",
		"type": "Sample"
	}
	`

	raw, err := autogen.Marshal(&model, false)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(raw))
}
