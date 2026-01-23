package datasource_test

import (
	"testing"

	"github.com/sebdah/goldie/v2"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/datasource"
)

type dsGenerationTestCase struct {
	goldenFileName string
	inputModel     codespec.Resource
}

func TestDataSourceGenerationFromCodeSpec(t *testing.T) {
	testCases := map[string]dsGenerationTestCase{
		"Basic read operation with path params": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Operations: codespec.APIOperations{
						Read: &codespec.APIOperation{
							HTTPMethod: "GET",
							Path:       "/api/v1/groups/{projectId}/testname/{name}",
						},
						VersionHeader: "application/vnd.atlas.2024-05-30+json",
					},
				},
			},
			goldenFileName: "ds-basic-read-with-path-params",
		},
		"Single path parameter": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Operations: codespec.APIOperations{
						Read: &codespec.APIOperation{
							HTTPMethod: "GET",
							Path:       "/api/v1/groups/{projectId}",
						},
						VersionHeader: "application/vnd.atlas.2024-05-30+json",
					},
				},
			},
			goldenFileName: "ds-single-path-param",
		},
		"No path parameters": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Operations: codespec.APIOperations{
						Read: &codespec.APIOperation{
							HTTPMethod: "GET",
							Path:       "/api/v1/testname",
						},
						VersionHeader: "application/vnd.atlas.2024-05-30+json",
					},
				},
			},
			goldenFileName: "ds-no-path-params",
		},
		"Three path parameters": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Operations: codespec.APIOperations{
						Read: &codespec.APIOperation{
							HTTPMethod: "GET",
							Path:       "/api/v1/orgs/{orgId}/groups/{projectId}/testname/{name}",
						},
						VersionHeader: "application/vnd.atlas.2024-05-30+json",
					},
				},
			},
			goldenFileName: "ds-three-path-params",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			result, err := datasource.GenerateGoCode(&tc.inputModel)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			g := goldie.New(t, goldie.WithNameSuffix(".golden.go"))
			g.Assert(t, tc.goldenFileName, result)
		})
	}
}

func TestDataSourceGenerationErrors(t *testing.T) {
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
			expectedErrMsg: "data source read operation is required for test_name",
		},
		"Missing Read operation": {
			inputModel: codespec.Resource{
				Name:        "test_name",
				PackageName: "testname",
				DataSources: &codespec.DataSources{
					Operations: codespec.APIOperations{
						Read:          nil,
						VersionHeader: "application/vnd.atlas.2024-05-30+json",
					},
				},
			},
			expectedErrMsg: "data source read operation is required for test_name",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			_, err := datasource.GenerateGoCode(&tc.inputModel)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
			if err.Error() != tc.expectedErrMsg {
				t.Fatalf("expected error message %q but got %q", tc.expectedErrMsg, err.Error())
			}
		})
	}
}

func TestPluralDataSourceGenerationFromCodeSpec(t *testing.T) {
	stringAttr := &codespec.StringAttribute{}
	int64Attr := &codespec.Int64Attribute{}
	boolAttr := &codespec.BoolAttribute{}
	listAttr := &codespec.ListAttribute{}

	testCases := map[string]dsGenerationTestCase{
		"Basic list operation without query params": {
			inputModel: codespec.Resource{
				Name:        "test_api",
				PackageName: "testapi",
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						PluralDSAttributes: &codespec.Attributes{},
					},
					Operations: codespec.APIOperations{
						List: &codespec.APIOperation{
							HTTPMethod: "GET",
							Path:       "/api/v1/groups/{projectId}/tests",
						},
						VersionHeader: "application/vnd.atlas.2024-05-30+json",
					},
				},
			},
			goldenFileName: "plural-ds-no-query-params",
		},
		"List operation with path and query params": {
			inputModel: codespec.Resource{
				Name:        "test_api",
				PackageName: "testapi",
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						PluralDSAttributes: &codespec.Attributes{
							{
								TFSchemaName:             "project_id",
								TFModelName:              "ProjectId",
								APIName:                  "groupId",
								ComputedOptionalRequired: codespec.Required,
								String:                   stringAttr,
							},
							{
								TFSchemaName:             "status",
								TFModelName:              "Status",
								APIName:                  "status",
								ComputedOptionalRequired: codespec.Optional,
								String:                   stringAttr,
							},
							{
								TFSchemaName:             "page_size",
								TFModelName:              "PageSize",
								APIName:                  "pageSize",
								ComputedOptionalRequired: codespec.Optional,
								Int64:                    int64Attr,
							},
						},
					},
					Operations: codespec.APIOperations{
						List: &codespec.APIOperation{
							HTTPMethod: "GET",
							Path:       "/api/v1/groups/{projectId}/tests",
						},
						VersionHeader: "application/vnd.atlas.2024-08-05+json",
					},
				},
			},
			goldenFileName: "plural-ds-with-query-params",
		},
		"List operation with multiple query param types": {
			inputModel: codespec.Resource{
				Name:        "service_account",
				PackageName: "serviceaccount",
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						PluralDSAttributes: &codespec.Attributes{
							{
								TFSchemaName:             "org_id",
								TFModelName:              "OrgId",
								APIName:                  "orgId",
								ComputedOptionalRequired: codespec.Required,
								String:                   stringAttr,
							},
							{
								TFSchemaName:             "name",
								TFModelName:              "Name",
								APIName:                  "name",
								ComputedOptionalRequired: codespec.Optional,
								String:                   stringAttr,
							},
							{
								TFSchemaName:             "limit",
								TFModelName:              "Limit",
								APIName:                  "limit",
								ComputedOptionalRequired: codespec.Optional,
								Int64:                    int64Attr,
							},
							{
								TFSchemaName:             "include_deleted",
								TFModelName:              "IncludeDeleted",
								APIName:                  "includeDeleted",
								ComputedOptionalRequired: codespec.Optional,
								Bool:                     boolAttr,
							},
							{
								TFSchemaName:             "types",
								TFModelName:              "Types",
								APIName:                  "types",
								ComputedOptionalRequired: codespec.Optional,
								List:                     listAttr,
							},
						},
					},
					Operations: codespec.APIOperations{
						List: &codespec.APIOperation{
							HTTPMethod: "GET",
							Path:       "/api/atlas/v2/orgs/{orgId}/serviceAccounts",
						},
						VersionHeader: "application/vnd.atlas.2024-08-05+json",
					},
				},
			},
			goldenFileName: "plural-ds-with-multiple-query-param-types",
		},
		"List operation with no path params": {
			inputModel: codespec.Resource{
				Name:        "control_plane_ip_addresses_api",
				PackageName: "controlplaneipaddressesapi",
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						PluralDSAttributes: &codespec.Attributes{},
					},
					Operations: codespec.APIOperations{
						List: &codespec.APIOperation{
							HTTPMethod: "GET",
							Path:       "/api/atlas/v2/unauth/controlPlaneIPAddresses",
						},
						VersionHeader: "application/vnd.atlas.2024-08-05+json",
					},
				},
			},
			goldenFileName: "plural-ds-no-path-params",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			result, err := datasource.GeneratePluralGoCode(&tc.inputModel)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			g := goldie.New(t, goldie.WithNameSuffix(".golden.go"))
			g.Assert(t, tc.goldenFileName, result)
		})
	}
}

func TestPluralDataSourceGenerationErrors(t *testing.T) {
	testCases := map[string]struct {
		expectedErrMsg string
		inputModel     codespec.Resource
	}{
		"Missing DataSources": {
			inputModel: codespec.Resource{
				Name:        "test_api",
				PackageName: "testapi",
				DataSources: nil,
			},
			expectedErrMsg: "data source list operation is required for plural data source test_api",
		},
		"Missing List operation": {
			inputModel: codespec.Resource{
				Name:        "test_api",
				PackageName: "testapi",
				DataSources: &codespec.DataSources{
					Operations: codespec.APIOperations{
						List:          nil,
						VersionHeader: "application/vnd.atlas.2024-05-30+json",
					},
				},
			},
			expectedErrMsg: "data source list operation is required for plural data source test_api",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			_, err := datasource.GeneratePluralGoCode(&tc.inputModel)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
			if err.Error() != tc.expectedErrMsg {
				t.Fatalf("expected error message %q but got %q", tc.expectedErrMsg, err.Error())
			}
		})
	}
}
