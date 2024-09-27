package schema_test

import (
	"testing"

	genconfigmapper "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/schema"
	"github.com/sebdah/goldie/v2"
	"go.mongodb.org/atlas-sdk/v20240530005/admin"
)

type schemaGenerationTestCase struct {
	inputModel     genconfigmapper.Resource
	goldenFileName string
}

func TestSchemaGenerationFromCodeSpec(t *testing.T) {
	testCases := map[string]schemaGenerationTestCase{
		"Primitive attributes": {
			inputModel: genconfigmapper.Resource{
				Name: "test_name",
				Schema: &genconfigmapper.Schema{
					Attributes: []genconfigmapper.Attribute{
						{
							Name:        "string_attr",
							String:      &genconfigmapper.StringAttribute{},
							Description: admin.PtrString("string description"),
						},
						{
							Name:        "bool_attr",
							Bool:        &genconfigmapper.BoolAttribute{},
							Description: admin.PtrString("bool description"),
						},
						{
							Name:        "int_attr",
							Int64:       &genconfigmapper.Int64Attribute{},
							Description: admin.PtrString("int description"),
						},
						{
							Name:        "float_attr",
							Float64:     &genconfigmapper.Float64Attribute{},
							Description: admin.PtrString("float description"),
						},
						{
							Name:        "number_attr",
							Number:      &genconfigmapper.NumberAttribute{},
							Description: admin.PtrString("number description"),
						},
						{
							Name: "simple_list_attr",
							List: &genconfigmapper.ListAttribute{
								ElementType: genconfigmapper.Float64,
							},
							Description: admin.PtrString("simple arr description"),
						},
						{
							Name: "simple_set_attr",
							Set: &genconfigmapper.SetAttribute{
								ElementType: genconfigmapper.Float64,
							},
							Description: admin.PtrString("simple set description"),
						},
						{
							Name: "simple_map_attr",
							Map: &genconfigmapper.MapAttribute{
								ElementType: genconfigmapper.Float64,
							},
							Description: admin.PtrString("simple map description"),
						},
					},
				},
			},
			goldenFileName: "primitive-attributes",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			result := schema.GenerateGoCode(tc.inputModel)
			g := goldie.New(t)
			g.Assert(t, tc.goldenFileName, []byte(result))
		})
	}
}
