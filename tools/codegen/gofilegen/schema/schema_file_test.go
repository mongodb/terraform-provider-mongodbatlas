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

var doubleNestedListAttr = codespec.Attribute{
	TFSchemaName:             "double_nested_list_attr",
	TFModelName:              "DoubleNestedListAttr",
	Description:              admin.PtrString("double nested list attribute"),
	ComputedOptionalRequired: codespec.Optional,
	ListNested: &codespec.ListNestedAttribute{
		NestedObject: codespec.NestedAttributeObject{
			Attributes: []codespec.Attribute{
				stringAttr,
			},
		},
	},
}

type schemaGenerationTestCase struct {
	inputModel     codespec.Resource
	goldenFileName string
	withObjType    bool
}

func TestSchemaGenerationFromCodeSpec(t *testing.T) {
	testCases := map[string]schemaGenerationTestCase{
		"Primitive attributes": {
			inputModel: codespec.Resource{
				Name: "test_name",
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
							TFSchemaName: "simple_list_attr",
							TFModelName:  "SimpleListAttr",
							List: &codespec.ListAttribute{
								ElementType: codespec.String,
							},
							Description:              admin.PtrString("simple arr description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							TFSchemaName: "simple_set_attr",
							TFModelName:  "SimpleSetAttr",
							Set: &codespec.SetAttribute{
								ElementType: codespec.Float64,
							},
							Description:              admin.PtrString("simple set description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							TFSchemaName: "simple_map_attr",
							TFModelName:  "SimpleMapAttr",
							Map: &codespec.MapAttribute{
								ElementType: codespec.Bool,
							},
							Description:              admin.PtrString("simple map description"),
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
					},
				},
			},
			withObjType:    true,
			goldenFileName: "primitive-attributes",
		},
		"Nested attributes": {
			inputModel: codespec.Resource{
				Name: "test_name",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{
						{
							TFSchemaName:             "nested_single_attr",
							TFModelName:              "NestedSingleAttr",
							Description:              admin.PtrString("nested single attribute"),
							ComputedOptionalRequired: codespec.Required,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{
										stringAttr,
										intAttr,
										{
											TFSchemaName:             "attr_not_included_in_req_bodies",
											TFModelName:              "AttrNotIncludedInReqBodies",
											String:                   &codespec.StringAttribute{},
											Description:              admin.PtrString("string description"),
											ComputedOptionalRequired: codespec.Computed,
											ReqBodyUsage:             codespec.OmitAlways,
										},
									},
								},
							},
						},
						{
							TFSchemaName:             "nested_list_attr",
							TFModelName:              "NestedListAttr",
							Description:              admin.PtrString("nested list attribute"),
							ComputedOptionalRequired: codespec.Optional,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{stringAttr, intAttr, doubleNestedListAttr},
								},
							},
						},
						{
							TFSchemaName:             "set_nested_attribute",
							TFModelName:              "SetNestedAttribute",
							Description:              admin.PtrString("set nested attribute"),
							ComputedOptionalRequired: codespec.ComputedOptional,
							SetNested: &codespec.SetNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{stringAttr, intAttr},
								},
							},
						},
						{
							TFSchemaName:             "map_nested_attribute",
							TFModelName:              "MapNestedAttribute",
							Description:              admin.PtrString("map nested attribute"),
							ComputedOptionalRequired: codespec.ComputedOptional,
							MapNested: &codespec.MapNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{stringAttr, intAttr},
								},
							},
						},
					},
				},
			},
			withObjType:    true,
			goldenFileName: "nested-attributes",
		},
		"Timeout attribute": {
			inputModel: codespec.Resource{
				Name: "test_name",
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
		"Avoid generation of ObjType definitions": {
			inputModel: codespec.Resource{
				Name: "test_name",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{
						{
							TFSchemaName:             "nested_list_attr",
							TFModelName:              "NestedListAttr",
							Description:              admin.PtrString("nested list attribute"),
							ComputedOptionalRequired: codespec.Optional,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{doubleNestedListAttr},
								},
							},
						},
					},
				},
			},
			withObjType:    false,
			goldenFileName: "no-obj-type-models",
		},
		"Multiple nested models with same parent attribute name": {
			inputModel: codespec.Resource{
				Name: "test_name",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{
						{
							TFSchemaName:             "first_nested_attr",
							TFModelName:              "FirstNestedAttr",
							Description:              admin.PtrString("first nested attribute"),
							ComputedOptionalRequired: codespec.Optional,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{doubleNestedListAttr},
								},
							},
						},
						{
							TFSchemaName:             "second_nested_attr",
							TFModelName:              "SecondNestedAttr",
							Description:              admin.PtrString("second nested attribute"),
							ComputedOptionalRequired: codespec.Optional,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{doubleNestedListAttr},
								},
							},
						},
					},
				},
			},
			withObjType:    true,
			goldenFileName: "multiple-nested-models-same-parent-attr-name",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			result := schema.GenerateGoCode(&tc.inputModel, tc.withObjType)
			g := goldie.New(t, goldie.WithNameSuffix(".golden.go"))
			g.Assert(t, tc.goldenFileName, []byte(result))
		})
	}
}
