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
		"Defining wait configuration in create update and delete": {
			inputModel: codespec.Resource{
				Name: stringcase.SnakeCaseString("test_name"),
				Operations: codespec.APIOperations{
					Create: codespec.APIOperation{
						HTTPMethod: "POST",
						Path:       "/api/v1/testname/{projectId}",
						Wait: &codespec.Wait{
							StateProperty:     "state",
							PendingStates:     []string{"INITIATING"},
							TargetStates:      []string{"IDLE"},
							TimeoutSeconds:    300,
							MinTimeoutSeconds: 60,
							DelaySeconds:      10,
						},
					},
					Update: codespec.APIOperation{
						HTTPMethod: "PUT",
						Path:       "/api/v1/testname/{projectId}",
						Wait: &codespec.Wait{
							StateProperty:     "state",
							PendingStates:     []string{"UPDATING"},
							TargetStates:      []string{"IDLE"},
							TimeoutSeconds:    300,
							MinTimeoutSeconds: 60,
							DelaySeconds:      10,
						},
					},
					Read: codespec.APIOperation{
						HTTPMethod: "GET",
						Path:       "/api/v1/testname/{projectId}",
					},
					Delete: codespec.APIOperation{
						HTTPMethod: "DELETE",
						Path:       "/api/v1/testname/{projectId}",
						Wait: &codespec.Wait{
							StateProperty:     "state",
							PendingStates:     []string{"PENDING"},
							TargetStates:      []string{"UNCONFIGURED", "DELETED"},
							TimeoutSeconds:    300,
							MinTimeoutSeconds: 60,
							DelaySeconds:      10,
						},
					},
					VersionHeader: "application/vnd.atlas.2024-05-30+json",
				},
			},
			goldenFileName: "wait-configuration",
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
