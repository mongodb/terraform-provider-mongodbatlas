package autogen_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalBasic(t *testing.T) {
	var model struct {
		AttrFloat        types.Float64 `tfsdk:"attr_float"`
		AttrFloatWithInt types.Float64 `tfsdk:"attr_float_with_int"`
		AttrString       types.String  `tfsdk:"attr_string"`
		AttrNotInJSON    types.String  `tfsdk:"attr_not_in_json"`
		AttrInt          types.Int64   `tfsdk:"attr_int"`
		AttrIntWithFloat types.Int64   `tfsdk:"attr_int_with_float"`
		AttrTrue         types.Bool    `tfsdk:"attr_true"`
		AttrFalse        types.Bool    `tfsdk:"attr_false"`
	}
	const (
		// attribute_not_in_model is ignored because it is not in the model, no error is thrown.
		// attribute_null is ignored because it is null, no error is thrown even if it is not in the model.
		jsonResp = `
			{
				"attrString": "value_string",
				"attrTrue": true,
				"attrFalse": false,
				"attrInt": 123,
				"attrIntWithFloat": 10.6,
				"attrFloat": 456.1,
				"attrFloatWithInt": 13,
				"attrNotInModel": "val",
				"attrNull": null
			}
		`
	)
	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, "value_string", model.AttrString.ValueString())
	assert.True(t, model.AttrTrue.ValueBool())
	assert.False(t, model.AttrFalse.ValueBool())
	assert.Equal(t, int64(123), model.AttrInt.ValueInt64())
	assert.Equal(t, int64(10), model.AttrIntWithFloat.ValueInt64()) // response floats stored in model ints have their decimals stripped.
	assert.InEpsilon(t, float64(456.1), model.AttrFloat.ValueFloat64(), epsilon)
	assert.InEpsilon(t, float64(13), model.AttrFloatWithInt.ValueFloat64(), epsilon)
	assert.True(t, model.AttrNotInJSON.IsNull()) // attributes not in JSON response are not changed, so null is kept.
}

func TestUnmarshalNestedAllTypes(t *testing.T) {
	type modelst struct {
		AttrObj               types.Object `tfsdk:"attr_obj"`
		AttrObjNullNotSent    types.Object `tfsdk:"attr_obj_null_not_sent"`
		AttrObjNullSent       types.Object `tfsdk:"attr_obj_null_sent"`
		AttrObjUnknownNotSent types.Object `tfsdk:"attr_obj_unknown_not_sent"`
		AttrObjUnknownSent    types.Object `tfsdk:"attr_obj_unknown_sent"`
		AttrObjParent         types.Object `tfsdk:"attr_obj_parent"`
		AttrListString        types.List   `tfsdk:"attr_list_string"`
		AttrListObj           types.List   `tfsdk:"attr_list_obj"`
		AttrSetString         types.Set    `tfsdk:"attr_set_string"`
		AttrSetObj            types.Set    `tfsdk:"attr_set_obj"`
		AttrListListString    types.List   `tfsdk:"attr_list_list_string"`
		AttrSetListObj        types.Set    `tfsdk:"attr_set_list_obj"`
		AttrListObjKnown      types.List   `tfsdk:"attr_list_obj_known"`
		AttrMapSimple         types.Map    `tfsdk:"attr_map_simple"`
		AttrMapSimpleExisting types.Map    `tfsdk:"attr_map_simple_existing"`
		AttrMapObj            types.Map    `tfsdk:"attr_map_obj"`
	}
	model := modelst{
		AttrObj: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
			// these attribute values are irrelevant, they will be overwritten with JSON values
			"attr_string": types.StringValue("different_string"),
			"attr_int":    types.Int64Value(123456),
			"attr_float":  types.Float64Unknown(), // can even be null
			"attr_bool":   types.BoolUnknown(),    // can even be unknown
		}),
		AttrObjNullNotSent:    types.ObjectNull(objTypeTest.AttrTypes),
		AttrObjNullSent:       types.ObjectNull(objTypeTest.AttrTypes),
		AttrObjUnknownNotSent: types.ObjectUnknown(objTypeTest.AttrTypes), // unknown values are changed to null
		AttrObjUnknownSent:    types.ObjectUnknown(objTypeTest.AttrTypes),
		AttrObjParent:         types.ObjectNull(objTypeParentTest.AttrTypes),
		AttrListString:        types.ListUnknown(types.StringType),
		AttrListObj:           types.ListUnknown(objTypeTest),
		AttrSetString:         types.SetUnknown(types.StringType),
		AttrSetObj:            types.SetUnknown(objTypeTest),
		AttrListListString:    types.ListUnknown(types.ListType{ElemType: types.StringType}),
		AttrSetListObj:        types.SetUnknown(types.ListType{ElemType: objTypeTest}),
		AttrListObjKnown: types.ListValueMust(objTypeTest, []attr.Value{
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("val"),
				"attr_int":    types.Int64Value(1),
				"attr_float":  types.Float64Value(1.1),
				"attr_bool":   types.BoolValue(true),
			}),
		}),
		AttrMapSimple: types.MapNull(types.StringType),
		AttrMapSimpleExisting: types.MapValueMust(types.StringType, map[string]attr.Value{
			"existing":       types.StringValue("valexisting"),
			"existingCHANGE": types.StringValue("before"),
		}),
		AttrMapObj: types.MapUnknown(objTypeTest),
	}
	// attrUnexisting is ignored because it is in JSON but not in the model, no error is returned
	const (
		jsonResp = `
			{
				"attrObj": {
					"attrString": "value_string",
					"attrInt": 123,
					"attrFloat": 1.1,
					"attrBool": true,
					"attrUnexisting": "val"
				}, 
				"attrObjNullSent": {
					"attrString": "null_obj",
					"attrInt": 1,
					"attrFloat": null
				},
				"attrObjUnknownSent": {
					"attrString": "unknown_obj"
				},
				"attrObjParent": {
					"attrParentString": "parent string",
					"attrParentObj": {
						"attrString": "inside parent string"
					}
				},
				"attrListString": [
					"list1",
					"list2"
				],
				"attrListObj": [
					{
						"attrString": "list1",
						"attrInt": 1,
						"attrFloat": 1.1,
						"attrBool": true
					},
					{
						"attrString": "list2",
						"attrInt": 2,
						"attrFloat": 2.2,
						"attrBool": false
					}
				],
				"attrSetString": [
					"set1",
					"set2"
				],
				"attrSetObj": [
					{
						"attrString": "set1",
						"attrInt": 11,
						"attrFloat": 11.1,
						"attrBool": false
					},
					{			
						"attrString": "set2",
						"attrInt": 22,
						"attrFloat": 22.2,		
						"attrBool": true		
					}
				],
				"attrListListString": [
					["list1a", "list1b"],
					["list2a", "list2b", "list2c"]
				],
				"attrSetListObj": [
					[{
						"attrString": "setList1",
						"attrInt": 1,
						"attrFloat": 1.1,
						"attrBool": true
					},
					{
						"attrString": "setList2",	
						"attrInt": 2,
						"attrFloat": 2.2,
						"attrBool": false
					}],
					[{
						"attrString": "setList3",	
						"attrInt": 3,
						"attrFloat": 3.3,
						"attrBool": true
					},
					{
						"attrString": "setList4",
						"attrInt": 4,					
						"attrFloat": 4.4,
						"attrBool": false
					},
					{
						"attrString": "setList5",
						"attrInt": 5,
						"attrFloat": 5.5,
						"attrBool": true
					}]
				],
				"attrListObjKnown": [
					{
						"attrString": "val2",
						"attrInt": 2
					}
				],
				"attrMapSimple": {
					"keyOne": "val1",
					"KeyTwo": "val2"
				},
				"attrMapSimpleExisting": {
					"key": "val",
					"existingCHANGE": "after"
				},
				"attrMapObj": {
					"obj1": {
						"attrString": "str1",
						"attrInt": 11,
						"attrFloat": 11.1,
						"attrBool": false
					},
					"obj2": {			
						"attrString": "str2",
						"attrInt": 22,
						"attrFloat": 22.2,		
						"attrBool": true		
					}
				}
			}
		`
	)
	modelExpected := modelst{
		AttrObj: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
			"attr_string": types.StringValue("value_string"),
			"attr_int":    types.Int64Value(123),
			"attr_float":  types.Float64Value(1.1),
			"attr_bool":   types.BoolValue(true),
		}),
		AttrObjNullNotSent: types.ObjectNull(objTypeTest.AttrTypes),
		AttrObjNullSent: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
			"attr_string": types.StringValue("null_obj"),
			"attr_int":    types.Int64Value(1),
			"attr_float":  types.Float64Null(),
			"attr_bool":   types.BoolNull(),
		}),
		AttrObjUnknownNotSent: types.ObjectUnknown(objTypeTest.AttrTypes),
		AttrObjUnknownSent: types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
			"attr_string": types.StringValue("unknown_obj"),
			"attr_int":    types.Int64Null(),
			"attr_float":  types.Float64Null(),
			"attr_bool":   types.BoolNull(),
		}),
		AttrObjParent: types.ObjectValueMust(objTypeParentTest.AttrTypes, map[string]attr.Value{
			"attr_parent_string": types.StringValue("parent string"),
			"attr_parent_int":    types.Int64Null(),
			"attr_parent_obj": types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("inside parent string"),
				"attr_int":    types.Int64Null(),
				"attr_float":  types.Float64Null(),
				"attr_bool":   types.BoolNull(),
			}),
		}),
		AttrListString: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("list1"),
			types.StringValue("list2"),
		}),
		AttrListObj: types.ListValueMust(objTypeTest, []attr.Value{
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("list1"),
				"attr_int":    types.Int64Value(1),
				"attr_float":  types.Float64Value(1.1),
				"attr_bool":   types.BoolValue(true),
			}),
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("list2"),
				"attr_int":    types.Int64Value(2),
				"attr_float":  types.Float64Value(2.2),
				"attr_bool":   types.BoolValue(false),
			}),
		}),
		AttrSetString: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("set1"),
			types.StringValue("set2"),
		}),
		AttrSetObj: types.SetValueMust(objTypeTest, []attr.Value{
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("set1"),
				"attr_int":    types.Int64Value(11),
				"attr_float":  types.Float64Value(11.1),
				"attr_bool":   types.BoolValue(false),
			}),
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("set2"),
				"attr_int":    types.Int64Value(22),
				"attr_float":  types.Float64Value(22.2),
				"attr_bool":   types.BoolValue(true),
			}),
		}),
		AttrListListString: types.ListValueMust(types.ListType{ElemType: types.StringType}, []attr.Value{
			types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("list1a"),
				types.StringValue("list1b"),
			}),
			types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("list2a"),
				types.StringValue("list2b"),
				types.StringValue("list2c"),
			}),
		}),
		AttrSetListObj: types.SetValueMust(types.ListType{ElemType: objTypeTest}, []attr.Value{
			types.ListValueMust(objTypeTest, []attr.Value{
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList1"),
					"attr_int":    types.Int64Value(1),
					"attr_float":  types.Float64Value(1.1),
					"attr_bool":   types.BoolValue(true),
				}),
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList2"),
					"attr_int":    types.Int64Value(2),
					"attr_float":  types.Float64Value(2.2),
					"attr_bool":   types.BoolValue(false),
				}),
			}),
			types.ListValueMust(objTypeTest, []attr.Value{
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList3"),
					"attr_int":    types.Int64Value(3),
					"attr_float":  types.Float64Value(3.3),
					"attr_bool":   types.BoolValue(true),
				}),
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList4"),
					"attr_int":    types.Int64Value(4),
					"attr_float":  types.Float64Value(4.4),
					"attr_bool":   types.BoolValue(false),
				}),
				types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
					"attr_string": types.StringValue("setList5"),
					"attr_int":    types.Int64Value(5),
					"attr_float":  types.Float64Value(5.5),
					"attr_bool":   types.BoolValue(true),
				}),
			}),
		}),
		AttrListObjKnown: types.ListValueMust(objTypeTest, []attr.Value{
			types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("val2"),
				"attr_int":    types.Int64Value(2),
				"attr_float":  types.Float64Value(1.1),
				"attr_bool":   types.BoolValue(true),
			}),
		}),
		AttrMapSimple: types.MapValueMust(types.StringType, map[string]attr.Value{
			"keyOne": types.StringValue("val1"),
			"KeyTwo": types.StringValue("val2"), // don't change the key case when it's a map
		}),
		AttrMapSimpleExisting: types.MapValueMust(types.StringType, map[string]attr.Value{
			"key":            types.StringValue("val"),
			"existing":       types.StringValue("valexisting"), // existing map values are kept
			"existingCHANGE": types.StringValue("after"),       // existing map values are changed if in JSON
		}),
		AttrMapObj: types.MapValueMust(objTypeTest, map[string]attr.Value{
			"obj1": types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("str1"),
				"attr_int":    types.Int64Value(11),
				"attr_float":  types.Float64Value(11.1),
				"attr_bool":   types.BoolValue(false),
			}),
			"obj2": types.ObjectValueMust(objTypeTest.AttrTypes, map[string]attr.Value{
				"attr_string": types.StringValue("str2"),
				"attr_int":    types.Int64Value(22),
				"attr_float":  types.Float64Value(22.2),
				"attr_bool":   types.BoolValue(true),
			}),
		}),
	}
	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, modelExpected, model)
}

func TestUnmarshalZeroLenCollections(t *testing.T) {
	type modelst struct {
		ListNullAbsent  types.List `tfsdk:"list_null_absent"`
		ListNullEmpty   types.List `tfsdk:"list_null_present"`
		ListNullNull    types.List `tfsdk:"list_null_present_null"`
		ListEmptyAbsent types.List `tfsdk:"list_empty_absent"`
		ListEmptyEmpty  types.List `tfsdk:"list_empty_present"`
		ListEmptyNull   types.List `tfsdk:"list_empty_present_null"`
	}
	model := modelst{
		ListNullAbsent:  types.ListNull(types.StringType),
		ListNullEmpty:   types.ListNull(types.StringType),
		ListNullNull:    types.ListNull(types.StringType),
		ListEmptyAbsent: types.ListValueMust(types.StringType, []attr.Value{}),
		ListEmptyEmpty:  types.ListValueMust(types.StringType, []attr.Value{}),
		ListEmptyNull:   types.ListValueMust(types.StringType, []attr.Value{}),
	}
	const (
		jsonResp = `
			{
				"list_null_empty": [],
				"list_null_null": null,
				"list_empty_empty": [],
				"list_empty_null": null
			}
		`
	)
	modelExpected := modelst{
		ListNullAbsent:  types.ListNull(types.StringType),
		ListNullEmpty:   types.ListNull(types.StringType),
		ListNullNull:    types.ListNull(types.StringType),
		ListEmptyAbsent: types.ListValueMust(types.StringType, []attr.Value{}),
		ListEmptyEmpty:  types.ListValueMust(types.StringType, []attr.Value{}),
		ListEmptyNull:   types.ListValueMust(types.StringType, []attr.Value{}),
	}
	require.NoError(t, autogen.Unmarshal([]byte(jsonResp), &model))
	assert.Equal(t, modelExpected, model)
}

type streamConn struct {
	Type        types.String `tfsdk:"type"`
	TypeCluster types.Object `tfsdk:"type_cluster" autogen:"discriminator:type=Cluster"`
	TypeHTTPS   types.Object `tfsdk:"type_https" autogen:"discriminator:type=Https"`
}

func (s *streamConn) DiscriminatorAttr(objJSON map[string]any) string {
	// Probably can return a single attribute
	switch objJSON["type"] {
	case "Cluster":
		return "TypeCluster"
	case "Https":
		return "TypeHTTPS"
	}
	return ""
}

var (
	dbRoleObjType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"role_name": types.StringType,
			"type":      types.StringType,
		},
	}
	clusterObjType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cluster_name":       types.StringType,
			"db_role_to_execute": dbRoleObjType,
		},
	}
	httpsObjType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"url": types.StringType,
		},
	}
	streamConnModelCluster = streamConn{
		Type: types.StringValue("Cluster"),
		TypeCluster: types.ObjectValueMust(clusterObjType.AttrTypes, map[string]attr.Value{
			"cluster_name": types.StringValue("myCluster"),
			"db_role_to_execute": types.ObjectValueMust(dbRoleObjType.AttrTypes, map[string]attr.Value{
				"role_name": types.StringValue("myRole"),
				"type":      types.StringValue("myType"),
			}),
		}),
		TypeHTTPS: types.ObjectNull(httpsObjType.AttrTypes),
	}
	streamConnModelHTTPS = streamConn{
		Type:        types.StringValue("Https"),
		TypeCluster: types.ObjectNull(clusterObjType.AttrTypes),
		TypeHTTPS: types.ObjectValueMust(httpsObjType.AttrTypes, map[string]attr.Value{
			"url": types.StringValue("https://example.com"),
		}),
	}
)

const (
	jsonRespCluster = `
			{
				"type": "Cluster",
				"clusterName": "myCluster",
				"dbRoleToExecute": {
					"roleName": "myRole",
					"type": "myType"
				}
		}`
	jsonRespHTTPS = `{
					"type": "Https",
					"url": "https://example.com"
				}`
)

func TestUnmarshalModelWithDiscriminator(t *testing.T) {
	testCases := map[string]struct {
		modelExpected streamConn
		jsonResp      string
	}{
		"cluster": {
			modelExpected: streamConnModelCluster,
			jsonResp:      jsonRespCluster,
		},
		"https": {
			modelExpected: streamConnModelHTTPS,
			jsonResp:      jsonRespHTTPS,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			model := streamConn{
				Type:        types.StringUnknown(),
				TypeCluster: types.ObjectNull(clusterObjType.AttrTypes),
				TypeHTTPS:   types.ObjectNull(httpsObjType.AttrTypes),
			}
			require.NoError(t, autogen.Unmarshal([]byte(tc.jsonResp), &model))
			assert.Equal(t, tc.modelExpected, model)
		})
	}
}

func TestUnmarshalErrors(t *testing.T) {
	testCases := map[string]struct {
		model        any
		responseJSON string
		errorStr     string
	}{
		"response ints are not converted to model strings": {
			errorStr:     "unmarshal of attribute attr expects type StringType but got Number with value: 1",
			responseJSON: `{"attr": 123}`,
			model: &struct {
				Attr types.String
			}{},
		},
		"response strings are not converted to model ints": {
			errorStr:     "unmarshal of attribute attr expects type Int64Type but got String with value: hello",
			responseJSON: `{"attr": "hello"}`,
			model: &struct {
				Attr types.Int64
			}{},
		},
		"response strings are not converted to model bools": {
			errorStr:     "unmarshal of attribute attr expects type BoolType but got String with value: true",
			responseJSON: `{"attr": "true"}`,
			model: &struct {
				Attr types.Bool
			}{},
		},
		"response bools are not converted to model string": {
			errorStr:     "unmarshal of attribute attr expects type StringType but got Bool with value: true",
			responseJSON: `{"attr": true}`,
			model: &struct {
				Attr types.String
			}{},
		},
		"model attributes have to be of Terraform types": {
			errorStr:     "unmarshal trying to set non-Terraform attribute Attr",
			responseJSON: `{"attr": "hello"}`,
			model: &struct {
				Attr string
			}{},
		},
		"model attr types in objects must match JSON types - string": {
			errorStr:     "unmarshal of attribute attr_string expects type StringType but got Number with value: 1",
			responseJSON: `{ "attrObj": { "attrString": 1 } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
		"model attr types in objects must match JSON types - bool": {
			errorStr:     "unmarshal of attribute attr_bool expects type BoolType but got String with value: not a bool",
			responseJSON: `{ "attrObj": { "attrBool": "not a bool" } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
		"model attr types in objects must match JSON types - int": {
			errorStr:     "unmarshal of attribute attr_int expects type Int64Type but got String with value: not an int",
			responseJSON: `{ "attrObj": { "attrInt": "not an int" } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
		"model attr types in objects must match JSON types - float": {
			errorStr:     "unmarshal of attribute attr_float expects type Float64Type but got String with value: not an int",
			responseJSON: `{ "attrObj": { "attrFloat": "not an int" } }`,
			model: &struct {
				AttrObj types.Object `tfsdk:"attr_obj"`
			}{
				AttrObj: types.ObjectNull(objTypeTest.AttrTypes),
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.ErrorContains(t, autogen.Unmarshal([]byte(tc.responseJSON), tc.model), tc.errorStr)
		})
	}
}

// TestUnmarshalUnsupportedModel has Terraform types not supported yet.
// It will be updated when we add support for them.
func TestUnmarshalUnsupportedModel(t *testing.T) {
	testCases := map[string]struct {
		model        any
		responseJSON string
	}{
		"Int32 not supported yet as it's not being used in any model": {
			responseJSON: `{"attr": 1}`,
			model: &struct {
				Attr types.Int32
			}{},
		},
		"Float32 not supported yet as it's not being used in any model": {
			responseJSON: `{"attr": 1}`,
			model: &struct {
				Attr types.Float32
			}{},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Error(t, autogen.Unmarshal([]byte(tc.responseJSON), tc.model))
		})
	}
}
