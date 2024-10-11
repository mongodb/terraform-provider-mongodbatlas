package codespec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

const (
	testFieldDesc     = "Test field description"
	testResourceDesc  = "Configures the project level settings for the Test Resource feature."
	testPathParamDesc = "Path param test description"
)

type convertToSpecTestCase struct {
	expectedResult       *codespec.Model
	inputOpenAPISpecPath string
	inputConfigPath      string
	inputResourceName    string
}

func TestConvertToProviderSpec(t *testing.T) {
	testCases := map[string]convertToSpecTestCase{
		"Valid input": {
			inputOpenAPISpecPath: "testdata/api-spec.yml",
			inputConfigPath:      "testdata/config-no-schema-opts.yml",
			inputResourceName:    "test_resource",

			expectedResult: &codespec.Model{
				Resources: []codespec.Resource{{
					Schema: &codespec.Schema{
						Description: conversion.StringPtr(testResourceDesc),
						Attributes: codespec.Attributes{
							codespec.Attribute{
								Name:                     "bool_default_attr",
								ComputedOptionalRequired: codespec.ComputedOptional,
								Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
							},
							codespec.Attribute{
								Name:                     "count",
								ComputedOptionalRequired: codespec.Optional,
								Int64:                    &codespec.Int64Attribute{},
								Description:              conversion.StringPtr(testFieldDesc),
							},
							codespec.Attribute{
								Name:                     "create_date",
								String:                   &codespec.StringAttribute{},
								ComputedOptionalRequired: codespec.Computed,
								Description:              conversion.StringPtr(testFieldDesc),
							},
							codespec.Attribute{
								Name:                     "group_id",
								ComputedOptionalRequired: codespec.Required,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testPathParamDesc),
							},
							codespec.Attribute{
								Name:                     "num_double_default_attr",
								Float64:                  &codespec.Float64Attribute{Default: conversion.Pointer(2.0)},
								ComputedOptionalRequired: codespec.ComputedOptional,
							},
							codespec.Attribute{
								Name:                     "str_computed_attr",
								ComputedOptionalRequired: codespec.Computed,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testFieldDesc),
							},
							codespec.Attribute{
								Name:                     "str_req_attr1",
								ComputedOptionalRequired: codespec.Required,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testFieldDesc),
							},
							codespec.Attribute{
								Name:                     "str_req_attr2",
								ComputedOptionalRequired: codespec.Required,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testFieldDesc),
							},
							codespec.Attribute{
								Name:                     "str_req_attr3",
								String:                   &codespec.StringAttribute{},
								ComputedOptionalRequired: codespec.Required,
								Description:              conversion.StringPtr(testFieldDesc),
							},
						},
					},
					Name: "test_resource",
				},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, _ := codespec.ToProviderSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, tc.inputResourceName)
			assert.Equal(t, tc.expectedResult.Resources[0], result.Resources[0], "Expected result to match the specified structure")
		})
	}
}
