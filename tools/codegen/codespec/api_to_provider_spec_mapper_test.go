package codespec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type convertToSpecTestCase struct {
	expectedResult       *codespec.CodeSpecification
	inputOpenAPISpecPath string
	inputConfigPath      string
	inputResourceName    string
}

func TestConvertToProviderSpec(t *testing.T) {
	testCases := map[string]convertToSpecTestCase{
		"Valid input": {
			inputOpenAPISpecPath: "testdata/api-spec.yml",
			inputConfigPath:      "testdata/config.yml",
			inputResourceName:    "test_resource",
			// TODO: replace with test case object after ToProviderSpecModel() implemented
			expectedResult: codespec.TestExampleCodeSpecification(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := codespec.ToProviderSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
			assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
		})
	}
}
