package codespec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

const (
	testFieldDesc     = "Test field description"
	testResourceDesc  = "POST API description"
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
							{
								Name:                     "bool_default_attr",
								ComputedOptionalRequired: codespec.ComputedOptional,
								Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
							},
							{
								Name:                     "count",
								ComputedOptionalRequired: codespec.Optional,
								Int64:                    &codespec.Int64Attribute{},
								Description:              conversion.StringPtr(testFieldDesc),
							},
							{
								Name:                     "create_date",
								String:                   &codespec.StringAttribute{},
								ComputedOptionalRequired: codespec.Computed,
								Description:              conversion.StringPtr(testFieldDesc),
							},
							{
								Name:                     "group_id",
								ComputedOptionalRequired: codespec.Required,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testPathParamDesc),
							},
							{
								Name:                     "num_double_default_attr",
								Float64:                  &codespec.Float64Attribute{Default: conversion.Pointer(2.0)},
								ComputedOptionalRequired: codespec.ComputedOptional,
							},
							{
								Name:                     "str_computed_attr",
								ComputedOptionalRequired: codespec.Computed,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testFieldDesc),
							},
							{
								Name:                     "str_req_attr1",
								ComputedOptionalRequired: codespec.Required,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testFieldDesc),
							},
							{
								Name:                     "str_req_attr2",
								ComputedOptionalRequired: codespec.Required,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testFieldDesc),
							},
							{
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
			result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, tc.inputResourceName)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
		})
	}
}

func TestConvertToProviderSpec_nested(t *testing.T) {
	testCases := map[string]convertToSpecTestCase{
		"Valid input": {
			inputOpenAPISpecPath: "testdata/api-spec.yml",
			inputConfigPath:      "testdata/config-nested-schema.yml",
			inputResourceName:    "test_resource_with_nested_attr",

			expectedResult: &codespec.Model{
				Resources: []codespec.Resource{{
					Schema: &codespec.Schema{
						Description: conversion.StringPtr(testResourceDesc),
						Attributes: codespec.Attributes{
							{
								Name:                     "cluster_name",
								ComputedOptionalRequired: codespec.Required,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testPathParamDesc),
							},
							{
								Name:                     "group_id",
								ComputedOptionalRequired: codespec.Required,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testPathParamDesc),
							},
							{
								Name:                     "list_primitive_string_attr",
								ComputedOptionalRequired: codespec.Computed,
								List: &codespec.ListAttribute{
									ElementType: codespec.String,
								},
								Description: conversion.StringPtr(testFieldDesc),
							},
							{
								Name:                     "nested_object_array_list_attr",
								ComputedOptionalRequired: codespec.Required,
								ListNested: &codespec.ListNestedAttribute{
									NestedObject: codespec.NestedAttributeObject{
										Attributes: codespec.Attributes{
											{
												Name:                     "inner_num_attr",
												ComputedOptionalRequired: codespec.Required,
												Int64:                    &codespec.Int64Attribute{},
												Description:              conversion.StringPtr(testFieldDesc),
											},
											{
												Name:                     "list_primitive_string_attr",
												ComputedOptionalRequired: codespec.Optional,
												List: &codespec.ListAttribute{
													ElementType: codespec.String,
												},
												Description: conversion.StringPtr(testFieldDesc),
											},
										},
									},
								},
								Description: conversion.StringPtr(testFieldDesc),
							},
							{
								Name:                     "nested_object_array_set_attr",
								ComputedOptionalRequired: codespec.Computed,
								SetNested: &codespec.SetNestedAttribute{
									NestedObject: codespec.NestedAttributeObject{
										Attributes: codespec.Attributes{
											{
												Name:                     "inner_num_attr",
												ComputedOptionalRequired: codespec.Required,
												Int64:                    &codespec.Int64Attribute{},
												Description:              conversion.StringPtr(testFieldDesc),
											},
											{
												Name:                     "list_primitive_string_attr",
												ComputedOptionalRequired: codespec.Optional,
												List: &codespec.ListAttribute{
													ElementType: codespec.String,
												},
												Description: conversion.StringPtr(testFieldDesc),
											},
										},
									},
								},
								Description: conversion.StringPtr(testFieldDesc),
							},
							{
								Name:                     "set_primitive_string_attr",
								ComputedOptionalRequired: codespec.Computed,
								Set: &codespec.SetAttribute{
									ElementType: codespec.String,
								},
								Description: conversion.StringPtr(testFieldDesc),
							},
							{
								Name:                     "single_nested_attr",
								ComputedOptionalRequired: codespec.Computed,
								SingleNested: &codespec.SingleNestedAttribute{
									NestedObject: codespec.NestedAttributeObject{
										Attributes: codespec.Attributes{
											{
												Name:                     "inner_int_attr",
												ComputedOptionalRequired: codespec.Required,
												Int64:                    &codespec.Int64Attribute{},
												Description:              conversion.StringPtr(testFieldDesc),
											},
											{
												Name:                     "inner_str_attr",
												ComputedOptionalRequired: codespec.Required,
												String:                   &codespec.StringAttribute{},
												Description:              conversion.StringPtr(testFieldDesc),
											},
										},
									},
								},
								Description: conversion.StringPtr(testFieldDesc),
							},
							{
								Name:                     "single_nested_attr_with_nested_maps",
								ComputedOptionalRequired: codespec.Computed,
								SingleNested: &codespec.SingleNestedAttribute{
									NestedObject: codespec.NestedAttributeObject{
										Attributes: codespec.Attributes{
											{
												Name:                     "aws",
												ComputedOptionalRequired: codespec.Optional,
												Map: &codespec.MapAttribute{
													ElementType: codespec.String,
												},
											},
											{
												Name:                     "azure",
												ComputedOptionalRequired: codespec.Optional,
												Map: &codespec.MapAttribute{
													ElementType: codespec.String,
												},
											},
										},
									},
								},
								Description: conversion.StringPtr(testFieldDesc),
							},
						},
					},
					Name: "test_resource_with_nested_attr",
				},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, tc.inputResourceName)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
		})
	}
}
