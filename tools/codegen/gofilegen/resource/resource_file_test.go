package resource_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/stringcase"
	"github.com/sebdah/goldie/v2"
)

type resourceGenerationTestCase struct {
	inputModel     codespec.Resource
	goldenFileName string
}

func TestResourceGenerationFromCodeSpec(t *testing.T) {
	testCases := map[string]resourceGenerationTestCase{
		"Defining different operation URLs with different path params": {
			inputModel: codespec.Resource{
				Name: stringcase.SnakeCaseString("test_name"),

				Operations: codespec.APIOperations{

					Create: codespec.APIOperation{
						HTTPMethod: "POST",
						Path:       "/api/v1/testname/{projectId}",
					},
					Update: codespec.APIOperation{
						HTTPMethod: "PATCH",
						Path:       "/api/v1/testname/{projectId}/{roleName}",
					},
					Read: codespec.APIOperation{
						HTTPMethod: "GET",
						Path:       "/api/v1/testname/{projectId}/{roleName}",
					},
					Delete: codespec.APIOperation{
						HTTPMethod: "DELETE",
						Path:       "/api/v1/testname/{projectId}/{roleName}",
					},
					VersionHeader: "application/vnd.atlas.2024-05-30+json",
				},
			},
			goldenFileName: "different-urls-with-path-params",
		},
		"Update operation using PUT": {
			inputModel: codespec.Resource{
				Name: stringcase.SnakeCaseString("test_name"),

				Operations: codespec.APIOperations{

					Create: codespec.APIOperation{
						HTTPMethod: "POST",
						Path:       "/api/v1/testname/{projectId}",
					},
					Update: codespec.APIOperation{
						HTTPMethod: "PUT",
						Path:       "/api/v1/testname/{projectId}",
					},
					Read: codespec.APIOperation{
						HTTPMethod: "GET",
						Path:       "/api/v1/testname/{projectId}",
					},
					Delete: codespec.APIOperation{
						HTTPMethod: "DELETE",
						Path:       "/api/v1/testname/{projectId}",
					},
					VersionHeader: "application/vnd.atlas.2024-05-30+json",
				},
			},
			goldenFileName: "update-with-put",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			result := resource.GenerateGoCode(&tc.inputModel)
			g := goldie.New(t, goldie.WithNameSuffix(".golden.go"))
			g.Assert(t, tc.goldenFileName, []byte(result))
		})
	}
}
