package codespec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type convertToSpecTestCase struct {
	expectedResult       *codespec.CodeSpecification
	inputOpenAPISpecPath string
	inputConfigPath      string
	inputResourceName    string
}

func TestConvertToProviderSpec(t *testing.T) {
	testFieldDesc := "Test field description"
	testCases := map[string]convertToSpecTestCase{
		"Valid input": {
			inputOpenAPISpecPath: "testdata/api-spec.yml",
			inputConfigPath:      "testdata/config.yml",
			inputResourceName:    "test_resource",
			expectedResult: &codespec.CodeSpecification{
				Resources: codespec.Resource{
					Schema: &codespec.Schema{
						Attributes: codespec.Attributes{
							codespec.Attribute{
								Name:        "project_id",
								IsRequired:  conversion.Pointer(true),
								String:      &codespec.StringAttribute{},
								Description: conversion.StringPtr("Overridden project_id description"),
							},
							codespec.Attribute{
								Name:        "bucket_name",
								IsRequired:  conversion.Pointer(true),
								String:      &codespec.StringAttribute{},
								Description: &testFieldDesc,
							},
							codespec.Attribute{
								Name:        "iam_role_id",
								IsRequired:  conversion.Pointer(true),
								String:      &codespec.StringAttribute{},
								Description: &testFieldDesc,
							},
							codespec.Attribute{
								Name:        "state",
								IsComputed:  conversion.Pointer(true),
								String:      &codespec.StringAttribute{},
								Description: &testFieldDesc,
							},
							codespec.Attribute{
								Name:        "prefix_path",
								String:      &codespec.StringAttribute{},
								IsComputed:  conversion.Pointer(true),
								IsOptional:  conversion.Pointer(true),
								Description: &testFieldDesc,
							},
						},
					},
					Name: "test_resource",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := codespec.ToProviderSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
			assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
		})
	}
}
