package datasource_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/datasource"
	"github.com/sebdah/goldie/v2"
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
