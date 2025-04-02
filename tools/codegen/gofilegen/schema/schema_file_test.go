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

type schemaGenerationTestCase struct {
	inputModel     codespec.Resource
	goldenFileName string
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
					},
				},
			},
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
									Attributes: []codespec.Attribute{stringAttr, intAttr},
								},
							},
						},
						{
							Name:                     "nested_list_attr",
							Description:              admin.PtrString("nested list attribute"),
							ComputedOptionalRequired: codespec.Optional,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: []codespec.Attribute{stringAttr, intAttr},
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
			goldenFileName: "nested-attributes",
		},
		"timeout attribute": {
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
						},
					},
				},
			},
			goldenFileName: "timeouts",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			result := schema.GenerateGoCode(&tc.inputModel)
			g := goldie.New(t, goldie.WithNameSuffix(".golden.go"))
			g.Assert(t, tc.goldenFileName, []byte(result))
		})
	}
}
