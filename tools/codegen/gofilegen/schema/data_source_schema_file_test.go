package schema_test

import (
	"testing"

	"github.com/sebdah/goldie/v2"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
)

type dsSchemaGenerationTestCase struct {
	goldenFileName string
	inputModel     codespec.Resource
}

func TestDataSourceSchemaGenerationFromCodeSpec(t *testing.T) {
	testCases := map[string]dsSchemaGenerationTestCase{
		"Primitive attributes": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						SingularDSAttributes: &codespec.Attributes{
							{
								TFSchemaName:             "string_attr",
								TFModelName:              "StringAttr",
								String:                   &codespec.StringAttribute{},
								Description:              new("string description"),
								ComputedOptionalRequired: codespec.Computed,
							},
							{
								TFSchemaName:             "bool_attr",
								TFModelName:              "BoolAttr",
								Bool:                     &codespec.BoolAttribute{},
								Description:              new("bool description"),
								ComputedOptionalRequired: codespec.Computed,
							},
							{
								TFSchemaName:             "int_attr",
								TFModelName:              "IntAttr",
								Int64:                    &codespec.Int64Attribute{},
								Description:              new("int description"),
								ComputedOptionalRequired: codespec.Computed,
							},
							{
								TFSchemaName:             "float_attr",
								TFModelName:              "FloatAttr",
								Float64:                  &codespec.Float64Attribute{},
								Description:              new("float description"),
								ComputedOptionalRequired: codespec.Computed,
							},
							{
								TFSchemaName:             "number_attr",
								TFModelName:              "NumberAttr",
								Number:                   &codespec.NumberAttribute{},
								Description:              new("number description"),
								ComputedOptionalRequired: codespec.Computed,
							},
						},
					},
				},
			},
			goldenFileName: "ds-primitive-attributes",
		},
		"Custom type attributes": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						SingularDSAttributes: &codespec.Attributes{
							{
								TFSchemaName:             "nested_object_attr",
								TFModelName:              "NestedObjectAttr",
								Description:              new("nested object attribute"),
								ComputedOptionalRequired: codespec.Computed,
								CustomType:               codespec.NewCustomObjectType("NestedObjectAttr"),
								SingleNested: &codespec.SingleNestedAttribute{
									NestedObject: codespec.NestedAttributeObject{
										Attributes: []codespec.Attribute{
											{
												TFSchemaName:             "string_attr",
												TFModelName:              "StringAttr",
												String:                   &codespec.StringAttribute{},
												Description:              new("string attribute"),
												ComputedOptionalRequired: codespec.Computed,
											},
										},
									},
								},
							},
							{
								TFSchemaName:             "string_list_attr",
								TFModelName:              "StringListAttr",
								Description:              new("string list attribute"),
								ComputedOptionalRequired: codespec.Computed,
								CustomType:               codespec.NewCustomListType(codespec.String),
								List:                     &codespec.ListAttribute{ElementType: codespec.String},
							},
							{
								TFSchemaName:             "nested_list_attr",
								TFModelName:              "NestedListAttr",
								Description:              new("nested list attribute"),
								ComputedOptionalRequired: codespec.Computed,
								CustomType:               codespec.NewCustomNestedListType("NestedListAttr"),
								ListNested: &codespec.ListNestedAttribute{
									NestedObject: codespec.NestedAttributeObject{
										Attributes: []codespec.Attribute{
											{
												TFSchemaName:             "int_attr",
												TFModelName:              "IntAttr",
												Int64:                    &codespec.Int64Attribute{},
												Description:              new("int attribute"),
												ComputedOptionalRequired: codespec.Computed,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			goldenFileName: "ds-custom-types-attributes",
		},
		"Deprecation message": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						DeprecationMessage: new("This data source is deprecated. Please use the test_name_new data source instead."),
						SingularDSAttributes: &codespec.Attributes{
							{
								TFSchemaName:             "string_attr",
								TFModelName:              "StringAttr",
								String:                   &codespec.StringAttribute{},
								Description:              new("string description"),
								ComputedOptionalRequired: codespec.Computed,
							},
						},
					},
				},
			},
			goldenFileName: "ds-deprecation-message",
		},
		"Required path parameter attribute": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						SingularDSAttributes: &codespec.Attributes{
							{
								TFSchemaName:             "project_id",
								TFModelName:              "ProjectId",
								String:                   &codespec.StringAttribute{},
								Description:              new("project identifier"),
								ComputedOptionalRequired: codespec.Required,
							},
							{
								TFSchemaName:             "name",
								TFModelName:              "Name",
								String:                   &codespec.StringAttribute{},
								Description:              new("resource name"),
								ComputedOptionalRequired: codespec.Computed,
							},
						},
					},
				},
			},
			goldenFileName: "ds-required-path-param",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			result, err := schema.GenerateDataSourceSchemaGoCode(&tc.inputModel)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			g := goldie.New(t, goldie.WithNameSuffix(".golden.go"))
			g.Assert(t, tc.goldenFileName, result)
		})
	}
}

func TestDataSourceSchemaGenerationErrors(t *testing.T) {
	testCases := map[string]struct {
		expectedErrMsg string
		inputModel     codespec.Resource
	}{
		"Missing DataSources": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: nil,
			},
			expectedErrMsg: "data source schema is required for test_name",
		},
		"Missing DataSources Schema": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Schema: nil,
				},
			},
			expectedErrMsg: "data source schema is required for test_name",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			_, err := schema.GenerateDataSourceSchemaGoCode(&tc.inputModel)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
			if err.Error() != tc.expectedErrMsg {
				t.Fatalf("expected error message %q but got %q", tc.expectedErrMsg, err.Error())
			}
		})
	}
}
