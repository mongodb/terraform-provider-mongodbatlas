package schema_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
	"github.com/sebdah/goldie/v2"
	"go.mongodb.org/atlas-sdk/v20240530005/admin"
)

var stringAttr = codespec.Attribute{
	Name:                     "string_attr",
	PascalCaseName:           "StringAttr",
	String:                   &codespec.StringAttribute{},
	Description:              admin.PtrString("string attribute"),
	ComputedOptionalRequired: codespec.Optional,
}

var intAttr = codespec.Attribute{
	Name:                     "int_attr",
	PascalCaseName:           "IntAttr",
	Int64:                    &codespec.Int64Attribute{},
	Description:              admin.PtrString("int attribute"),
	ComputedOptionalRequired: codespec.Required,
}

var doubleNestedListAttr = codespec.Attribute{
	Name:                     "double_nested_list_attr",
	PascalCaseName:           "DoubleNestedListAttr",
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
							Name:                     "string_attr",
							PascalCaseName:           "StringAttr",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
						},
						{
							Name:                     "bool_attr",
							PascalCaseName:           "BoolAttr",
							Bool:                     &codespec.BoolAttribute{},
							Description:              admin.PtrString("bool description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:                     "int_attr",
							PascalCaseName:           "IntAttr",
							Int64:                    &codespec.Int64Attribute{},
							Description:              admin.PtrString("int description"),
							ComputedOptionalRequired: codespec.ComputedOptional,
						},
						{
							Name:                     "float_attr",
							PascalCaseName:           "FloatAttr",
							Float64:                  &codespec.Float64Attribute{},
							Description:              admin.PtrString("float description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:                     "number_attr",
							PascalCaseName:           "NumberAttr",
							Number:                   &codespec.NumberAttribute{},
							Description:              admin.PtrString("number description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:           "simple_list_attr",
							PascalCaseName: "SimpleListAttr",
							List: &codespec.ListAttribute{
								ElementType: codespec.String,
							},
							Description:              admin.PtrString("simple arr description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:           "simple_set_attr",
							PascalCaseName: "SimpleSetAttr",
							Set: &codespec.SetAttribute{
								ElementType: codespec.Float64,
							},
							Description:              admin.PtrString("simple set description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:           "simple_map_attr",
							PascalCaseName: "SimpleMapAttr",
							Map: &codespec.MapAttribute{
								ElementType: codespec.Bool,
							},
							Description:              admin.PtrString("simple map description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:                     "attr_not_included_in_req_bodies",
							PascalCaseName:           "AttrNotIncludedInReqBodies",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "attr_only_in_post_req_bodies",
							PascalCaseName:           "AttrOnlyInPostReqBodies",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
							ReqBodyUsage:             codespec.OmitInUpdateBody,
						},
						{
							Name:                     "json_attr",
							PascalCaseName:           "JsonAttr",
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
							Name:                     "nested_single_attr",
							PascalCaseName:           "NestedSingleAttr",
							Description:              admin.PtrString("nested single attribute"),
							ComputedOptionalRequired: codespec.Required,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{
										stringAttr,
										intAttr,
										{
											Name:                     "attr_not_included_in_req_bodies",
											PascalCaseName:           "AttrNotIncludedInReqBodies",
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
							Name:                     "nested_list_attr",
							PascalCaseName:           "NestedListAttr",
							Description:              admin.PtrString("nested list attribute"),
							ComputedOptionalRequired: codespec.Optional,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{stringAttr, intAttr, doubleNestedListAttr},
								},
							},
						},
						{
							Name:                     "set_nested_attribute",
							PascalCaseName:           "SetNestedAttribute",
							Description:              admin.PtrString("set nested attribute"),
							ComputedOptionalRequired: codespec.ComputedOptional,
							SetNested: &codespec.SetNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{stringAttr, intAttr},
								},
							},
						},
						{
							Name:                     "map_nested_attribute",
							PascalCaseName:           "MapNestedAttribute",
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
							Name:                     "string_attr",
							PascalCaseName:           "StringAttr",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
						},
						{
							Name:           "timeouts",
							PascalCaseName: "Timeouts",
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
							Name:                     "nested_list_attr",
							PascalCaseName:           "NestedListAttr",
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
							Name:                     "first_nested_attr",
							PascalCaseName:           "FirstNestedAttr",
							Description:              admin.PtrString("first nested attribute"),
							ComputedOptionalRequired: codespec.Optional,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{doubleNestedListAttr},
								},
							},
						},
						{
							Name:                     "second_nested_attr",
							PascalCaseName:           "SecondNestedAttr",
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
