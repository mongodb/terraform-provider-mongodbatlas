package schema_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
	"github.com/sebdah/goldie/v2"
	"go.mongodb.org/atlas-sdk/v20240530005/admin"
)

var stringAttr = codespec.Attribute{
	TFSchemaName:             "string_attr",
	TFModelName:              "StringAttr",
	String:                   &codespec.StringAttribute{},
	Description:              admin.PtrString("string attribute"),
	ComputedOptionalRequired: codespec.Optional,
}

var intAttr = codespec.Attribute{
	TFSchemaName:             "int_attr",
	TFModelName:              "IntAttr",
	Int64:                    &codespec.Int64Attribute{},
	Description:              admin.PtrString("int attribute"),
	ComputedOptionalRequired: codespec.Required,
}

func doubleCustomNestedListAttr(ancestorName string) codespec.Attribute {
	return codespec.Attribute{
		TFSchemaName:             "double_nested_list_attr",
		TFModelName:              "DoubleNestedListAttr",
		Description:              admin.PtrString("double nested list attribute"),
		ComputedOptionalRequired: codespec.Optional,
		CustomType:               codespec.NewCustomNestedListType(ancestorName + "DoubleNestedListAttr"),
		ListNested: &codespec.ListNestedAttribute{
			NestedObject: codespec.NestedAttributeObject{
				Attributes: []codespec.Attribute{
					stringAttr,
				},
			},
		},
	}
}

type schemaGenerationTestCase struct {
	inputModel     codespec.Resource
	goldenFileName string
}

//nolint:funlen // Long test data
func TestSchemaGenerationFromCodeSpec(t *testing.T) {
	schemaGenFromCodeSpecTestCases := map[string]schemaGenerationTestCase{
		"Primitive attributes": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{
						{
							TFSchemaName:             "string_attr",
							TFModelName:              "StringAttr",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
						},
						{
							TFSchemaName:             "bool_attr",
							TFModelName:              "BoolAttr",
							Bool:                     &codespec.BoolAttribute{},
							Description:              admin.PtrString("bool description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							TFSchemaName:             "int_attr",
							TFModelName:              "IntAttr",
							Int64:                    &codespec.Int64Attribute{},
							Description:              admin.PtrString("int description"),
							ComputedOptionalRequired: codespec.ComputedOptional,
						},
						{
							TFSchemaName:             "float_attr",
							TFModelName:              "FloatAttr",
							Float64:                  &codespec.Float64Attribute{},
							Description:              admin.PtrString("float description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							TFSchemaName:             "number_attr",
							TFModelName:              "NumberAttr",
							Number:                   &codespec.NumberAttribute{},
							Description:              admin.PtrString("number description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							TFSchemaName:             "attr_not_included_in_req_bodies",
							TFModelName:              "AttrNotIncludedInReqBodies",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							TFSchemaName:             "attr_only_in_post_req_bodies",
							TFModelName:              "AttrOnlyInPostReqBodies",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
							ReqBodyUsage:             codespec.OmitInUpdateBody,
						},
						{
							TFSchemaName:             "json_attr",
							TFModelName:              "JsonAttr",
							String:                   &codespec.StringAttribute{},
							CustomType:               &codespec.CustomTypeJSONVar,
							Description:              admin.PtrString("json description"),
							ComputedOptionalRequired: codespec.Required,
						},
						{
							TFSchemaName:             "sensitive_string_attr",
							TFModelName:              "SensitiveStringAttr",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("sensitive string description"),
							ComputedOptionalRequired: codespec.Required,
							ReqBodyUsage:             codespec.OmitInUpdateBody,
							Sensitive:                true,
						},
					},
				},
			},
			goldenFileName: "primitive-attributes",
		},
		"Custom type attributes": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{
						{
							TFSchemaName:             "nested_object_attr",
							TFModelName:              "NestedObjectAttr",
							Description:              admin.PtrString("nested object attribute"),
							ComputedOptionalRequired: codespec.Required,
							CustomType:               codespec.NewCustomObjectType("NestedObjectAttr"),
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{
										stringAttr,
										{
											TFSchemaName: "sub_nested_object_attr",
											TFModelName:  "SubNestedObjectAttr",
											CustomType:   codespec.NewCustomObjectType("NestedObjectAttrSubNestedObjectAttr"),
											SingleNested: &codespec.SingleNestedAttribute{
												NestedObject: codespec.NestedAttributeObject{
													Attributes: []codespec.Attribute{
														intAttr,
													},
												},
											},
											ComputedOptionalRequired: codespec.Required,
										},
									},
								},
							},
						},
						{
							TFSchemaName:             "string_list_attr",
							TFModelName:              "StringListAttr",
							Description:              admin.PtrString("string list attribute"),
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomListType(codespec.String),
							List:                     &codespec.ListAttribute{ElementType: codespec.String},
						},
						{
							TFSchemaName:             "nested_list_attr",
							TFModelName:              "NestedListAttr",
							Description:              admin.PtrString("nested list attribute"),
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomNestedListType("NestedListAttr"),
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{stringAttr, intAttr, doubleCustomNestedListAttr("NestedListAttr")},
								},
							},
						},
						{
							TFSchemaName:             "string_set_attr",
							TFModelName:              "StringSetAttr",
							Description:              admin.PtrString("string set attribute"),
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomSetType(codespec.String),
							Set:                      &codespec.SetAttribute{ElementType: codespec.String},
						},
						{
							TFSchemaName:             "nested_set_attr",
							TFModelName:              "NestedSetAttr",
							Description:              admin.PtrString("nested set attribute"),
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomNestedSetType("NestedSetAttr"),
							SetNested: &codespec.SetNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{intAttr},
								},
							},
						},
						{
							TFSchemaName:             "string_map_attr",
							TFModelName:              "StringMapAttr",
							Description:              admin.PtrString("string map attribute"),
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomMapType(codespec.String),
							Map:                      &codespec.MapAttribute{ElementType: codespec.String},
						},
						{
							TFSchemaName:             "map_nested_attribute",
							TFModelName:              "MapNestedAttribute",
							Description:              admin.PtrString("nested map attribute"),
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomNestedMapType("MapNestedAttribute"),
							MapNested: &codespec.MapNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{stringAttr},
								},
							},
						},
					},
				},
			},
			goldenFileName: "custom-types-attributes",
		},
		"Timeout attribute": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{
						{
							TFSchemaName:             "string_attr",
							TFModelName:              "StringAttr",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
						},
						{
							TFSchemaName: "timeouts",
							TFModelName:  "Timeouts",
							Timeouts: &codespec.TimeoutsAttribute{
								ConfigurableTimeouts: []codespec.Operation{codespec.Create, codespec.Update, codespec.Delete},
							},
							ReqBodyUsage: codespec.OmitAlways,
						},
					},
				},
			},
			goldenFileName: "timeouts",
		},
		"Multiple nested models with same parent attribute name": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{
						{
							TFSchemaName:             "first_nested_attr",
							TFModelName:              "FirstNestedAttr",
							Description:              admin.PtrString("first nested attribute"),
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomNestedListType("FirstNestedAttr"),
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{doubleCustomNestedListAttr("FirstNestedAttr")},
								},
							},
						},
						{
							TFSchemaName:             "second_nested_attr",
							TFModelName:              "SecondNestedAttr",
							Description:              admin.PtrString("second nested attribute"),
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomNestedListType("SecondNestedAttr"),
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{doubleCustomNestedListAttr("SecondNestedAttr")},
								},
							},
						},
					},
				},
			},
			goldenFileName: "multiple-nested-models-same-parent-attr-name",
		},
		"Plan modifiers using create only": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{
						{
							TFSchemaName:             "string_attr",
							TFModelName:              "StringAttr",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "bool_attr",
							TFModelName:              "BoolAttr",
							Bool:                     &codespec.BoolAttribute{},
							Description:              admin.PtrString("bool description"),
							ComputedOptionalRequired: codespec.Optional,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "int_attr",
							TFModelName:              "IntAttr",
							Int64:                    &codespec.Int64Attribute{},
							Description:              admin.PtrString("int description"),
							ComputedOptionalRequired: codespec.ComputedOptional,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "float_attr",
							TFModelName:              "FloatAttr",
							Float64:                  &codespec.Float64Attribute{},
							Description:              admin.PtrString("float description"),
							ComputedOptionalRequired: codespec.Optional,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "number_attr",
							TFModelName:              "NumberAttr",
							Number:                   &codespec.NumberAttribute{},
							Description:              admin.PtrString("number description"),
							ComputedOptionalRequired: codespec.Optional,
							CreateOnly:               true,
						},
						{
							TFSchemaName: "simple_list_attr",
							TFModelName:  "SimpleListAttr",
							CustomType:   codespec.NewCustomListType(codespec.String),
							List: &codespec.ListAttribute{
								ElementType: codespec.String,
							},
							Description:              admin.PtrString("simple arr description"),
							ComputedOptionalRequired: codespec.Optional,
							CreateOnly:               true,
						},
						{
							TFSchemaName: "simple_set_attr",
							TFModelName:  "SimpleSetAttr",
							CustomType:   codespec.NewCustomSetType(codespec.Float64),
							Set: &codespec.SetAttribute{
								ElementType: codespec.Float64,
							},
							Description:              admin.PtrString("simple set description"),
							ComputedOptionalRequired: codespec.Optional,
							CreateOnly:               true,
						},
						{
							TFSchemaName: "simple_map_attr",
							TFModelName:  "SimpleMapAttr",
							CustomType:   codespec.NewCustomMapType(codespec.Bool),
							Map: &codespec.MapAttribute{
								ElementType: codespec.Bool,
							},
							Description:              admin.PtrString("simple map description"),
							ComputedOptionalRequired: codespec.Optional,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "nested_single_attr",
							TFModelName:              "NestedSingleAttr",
							Description:              admin.PtrString("nested single attribute"),
							ComputedOptionalRequired: codespec.Required,
							CustomType:               codespec.NewCustomObjectType("NestedSingleAttr"),
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{intAttr},
								},
							},
							CreateOnly: true,
						},
						{
							TFSchemaName:             "nested_list_attr",
							TFModelName:              "NestedListAttr",
							Description:              admin.PtrString("nested list attribute"),
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomNestedListType("NestedListAttr"),
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{intAttr},
								},
							},
							CreateOnly: true,
						},
						{
							TFSchemaName:             "set_nested_attribute",
							TFModelName:              "SetNestedAttribute",
							Description:              admin.PtrString("set nested attribute"),
							ComputedOptionalRequired: codespec.ComputedOptional,
							CustomType:               codespec.NewCustomNestedSetType("SetNestedAttribute"),
							SetNested: &codespec.SetNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{intAttr},
								},
							},
							CreateOnly: true,
						},
						{
							TFSchemaName:             "map_nested_attribute",
							TFModelName:              "MapNestedAttribute",
							Description:              admin.PtrString("map nested attribute"),
							ComputedOptionalRequired: codespec.ComputedOptional,
							CustomType:               codespec.NewCustomNestedMapType("MapNestedAttribute"),
							MapNested: &codespec.MapNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{intAttr},
								},
							},
							CreateOnly: true,
						},
					},
				},
			},
			goldenFileName: "plan-modifiers-create-only",
		},
	}

	for testName, tc := range schemaGenFromCodeSpecTestCases {
		t.Run(testName, func(t *testing.T) {
			result := schema.GenerateGoCode(&tc.inputModel)
			g := goldie.New(t, goldie.WithNameSuffix(".golden.go"))
			g.Assert(t, tc.goldenFileName, result)
		})
	}
}
