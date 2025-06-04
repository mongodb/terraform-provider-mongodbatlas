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
	String:                   &codespec.StringAttribute{},
	Description:              admin.PtrString("string attribute"),
	ComputedOptionalRequired: codespec.Optional,
}

var intAttr = codespec.Attribute{
	Name:                     "int_attr",
	Int64:                    &codespec.Int64Attribute{},
	Description:              admin.PtrString("int attribute"),
	ComputedOptionalRequired: codespec.Required,
}

var doubleNestedListAttr = codespec.Attribute{
	Name:                     "double_nested_list_attr",
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
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
						},
						{
							Name:                     "bool_attr",
							Bool:                     &codespec.BoolAttribute{},
							Description:              admin.PtrString("bool description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:                     "int_attr",
							Int64:                    &codespec.Int64Attribute{},
							Description:              admin.PtrString("int description"),
							ComputedOptionalRequired: codespec.ComputedOptional,
						},
						{
							Name:                     "float_attr",
							Float64:                  &codespec.Float64Attribute{},
							Description:              admin.PtrString("float description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:                     "number_attr",
							Number:                   &codespec.NumberAttribute{},
							Description:              admin.PtrString("number description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name: "simple_list_attr",
							List: &codespec.ListAttribute{
								ElementType: codespec.String,
							},
							Description:              admin.PtrString("simple arr description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name: "simple_set_attr",
							Set: &codespec.SetAttribute{
								ElementType: codespec.Float64,
							},
							Description:              admin.PtrString("simple set description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name: "simple_map_attr",
							Map: &codespec.MapAttribute{
								ElementType: codespec.Bool,
							},
							Description:              admin.PtrString("simple map description"),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:                     "attr_not_included_in_req_bodies",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "attr_only_in_post_req_bodies",
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
							ReqBodyUsage:             codespec.OmitInUpdateBody,
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
							Description:              admin.PtrString("nested single attribute"),
							ComputedOptionalRequired: codespec.Required,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{
										stringAttr,
										intAttr,
										{
											Name:                     "attr_not_included_in_req_bodies",
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
							String:                   &codespec.StringAttribute{},
							Description:              admin.PtrString("string description"),
							ComputedOptionalRequired: codespec.Required,
						},
						{
							Name: "timeouts",
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
		"Discriminator mapping": {
			inputModel: codespec.Resource{
				Name: "test_name",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{
						{
							Name:                     "type",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Computed,
							Description:              admin.PtrString("Type of the stream connection"),
						},
						{
							Name:                     "type_cluster",
							ComputedOptionalRequired: codespec.Optional,
							Description:              admin.PtrString("Use this when you want a cluster stream connection"),
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{
										{
											Name:                     "cluster_name",
											String:                   &codespec.StringAttribute{},
											ComputedOptionalRequired: codespec.Required,
											Description:              admin.PtrString("Name of the cluster to connect to"),
										},
									},
								},
							},
							Discriminator: &codespec.DiscriminatorMapping{
								DiscriminatorProperty: "type",
								DiscriminatorValue:    "Cluster",
							},
						},
						{
							Name:                     "type_https",
							ComputedOptionalRequired: codespec.Optional,
							Description:              admin.PtrString("Use this when you want a https stream connection"),
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{
										{
											Name:                     "url",
											String:                   &codespec.StringAttribute{},
											ComputedOptionalRequired: codespec.Required,
											Description:              admin.PtrString("Url of the https stream connection"),
										},
									},
								},
							},
							Discriminator: &codespec.DiscriminatorMapping{
								DiscriminatorProperty: "type",
								DiscriminatorValue:    "Https",
							},
						},
					},
				},
			},
			withObjType:    false,
			goldenFileName: "discriminator-mapping-stream-connection",
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
