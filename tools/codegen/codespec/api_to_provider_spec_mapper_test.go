package codespec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

const (
	testFieldDesc       = "Test field description"
	testResourceDesc    = "POST API description"
	testPathParamDesc   = "Path param test description"
	testDataAPISpecPath = "testdata/api-spec.yml"
	testDataConfigPath  = "testdata/config.yml"
)

type convertToSpecTestCase struct {
	expectedResult       *codespec.Model
	inputOpenAPISpecPath string
	inputConfigPath      string
	inputResourceName    string
}

func TestConvertToProviderSpec(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_no_schema_opts",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr(testResourceDesc),
					Attributes: codespec.Attributes{
						{
							Name:                     "bool_default_attr",
							PascalCaseName:           "BoolDefaultAttr",
							ComputedOptionalRequired: codespec.ComputedOptional,
							Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
						},
						{
							Name:                     "count",
							PascalCaseName:           "Count",
							ComputedOptionalRequired: codespec.Optional,
							Int64:                    &codespec.Int64Attribute{},
							Description:              conversion.StringPtr(testFieldDesc),
						},
						{
							Name:                     "create_date",
							PascalCaseName:           "CreateDate",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Computed,
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "group_id",
							PascalCaseName:           "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "num_double_default_attr",
							PascalCaseName:           "NumDoubleDefaultAttr",
							Float64:                  &codespec.Float64Attribute{Default: conversion.Pointer(2.0)},
							ComputedOptionalRequired: codespec.ComputedOptional,
						},
						{
							Name:                     "str_computed_attr",
							PascalCaseName:           "StrComputedAttr",
							ComputedOptionalRequired: codespec.Computed,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "str_req_attr1",
							PascalCaseName:           "StrReqAttr1",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
						},
						{
							Name:                     "str_req_attr2",
							PascalCaseName:           "StrReqAttr2",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
						},
						{
							Name:                     "str_req_attr3",
							PascalCaseName:           "StrReqAttr3",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Required,
							Description:              conversion.StringPtr(testFieldDesc),
						},
					},
				},
				Name: "test_resource_no_schema_opts",
				Operations: codespec.APIOperations{
					Create: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "POST",
					},
					Read: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "GET",
					},
					Update: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
			}},
		},
	}
	runTestCase(t, tc)
}

func TestConvertToProviderSpec_nested(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_with_nested_attr",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr(testResourceDesc),
					Attributes: codespec.Attributes{
						{
							Name:                     "attr_always_in_updates",
							PascalCaseName:           "AttrAlwaysInUpdates",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Always in updates"),
						},
						{
							Name:                     "cluster_name",
							PascalCaseName:           "ClusterName",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "group_id",
							PascalCaseName:           "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "list_primitive_string_attr",
							PascalCaseName:           "ListPrimitiveStringAttr",
							ComputedOptionalRequired: codespec.Computed,
							List: &codespec.ListAttribute{
								ElementType: codespec.String,
							},
							Description:  conversion.StringPtr(testFieldDesc),
							ReqBodyUsage: codespec.OmitAlways,
						},
						{
							Name:                     "nested_list_array_attr",
							PascalCaseName:           "NestedListArrayAttr",
							ComputedOptionalRequired: codespec.Required,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "inner_num_attr",
											PascalCaseName:           "InnerNumAttr",
											ComputedOptionalRequired: codespec.Required,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
										},
										{
											Name:                     "list_primitive_string_attr",
											PascalCaseName:           "ListPrimitiveStringAttr",
											ComputedOptionalRequired: codespec.Optional,
											List: &codespec.ListAttribute{
												ElementType: codespec.String,
											},
											Description: conversion.StringPtr(testFieldDesc),
										},
										{
											Name:                     "list_primitive_string_computed_attr",
											PascalCaseName:           "ListPrimitiveStringComputedAttr",
											ComputedOptionalRequired: codespec.Computed,
											List: &codespec.ListAttribute{
												ElementType: codespec.String,
											},
											Description:  conversion.StringPtr(testFieldDesc),
											ReqBodyUsage: codespec.OmitAlways,
										},
									},
								},
							},
							Description: conversion.StringPtr(testFieldDesc),
						},
						{
							Name:                     "nested_map_object_attr",
							PascalCaseName:           "NestedMapObjectAttr",
							ComputedOptionalRequired: codespec.Computed,
							MapNested: &codespec.MapNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "attr",
											PascalCaseName:           "Attr",
											ComputedOptionalRequired: codespec.Computed,
											String:                   &codespec.StringAttribute{},
											ReqBodyUsage:             codespec.OmitAlways,
										},
									},
								},
							},
							ReqBodyUsage: codespec.OmitAlways,
						},
						{
							Name:                     "nested_set_array_attr",
							PascalCaseName:           "NestedSetArrayAttr",
							ComputedOptionalRequired: codespec.Computed,
							SetNested: &codespec.SetNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "inner_num_attr",
											PascalCaseName:           "InnerNumAttr",
											ComputedOptionalRequired: codespec.Computed,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
										},
										{
											Name:                     "list_primitive_string_attr",
											PascalCaseName:           "ListPrimitiveStringAttr",
											ComputedOptionalRequired: codespec.Computed,
											List: &codespec.ListAttribute{
												ElementType: codespec.String,
											},
											Description:  conversion.StringPtr(testFieldDesc),
											ReqBodyUsage: codespec.OmitAlways,
										},
									},
								},
							},
							ReqBodyUsage: codespec.OmitAlways,
							Description:  conversion.StringPtr(testFieldDesc),
						},
						{
							Name:                     "optional_string_attr",
							PascalCaseName:           "OptionalStringAttr",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Optional string"),
						},
						{
							Name:                     "set_primitive_string_attr",
							PascalCaseName:           "SetPrimitiveStringAttr",
							ComputedOptionalRequired: codespec.Computed,
							Set: &codespec.SetAttribute{
								ElementType: codespec.String,
							},
							ReqBodyUsage: codespec.OmitAlways,
							Description:  conversion.StringPtr(testFieldDesc),
						},
						{
							Name:                     "single_nested_attr",
							PascalCaseName:           "SingleNestedAttr",
							ComputedOptionalRequired: codespec.Computed,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "inner_int_attr",
											PascalCaseName:           "InnerIntAttr",
											ComputedOptionalRequired: codespec.Computed,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
										},
										{
											Name:                     "inner_str_attr",
											PascalCaseName:           "InnerStrAttr",
											ComputedOptionalRequired: codespec.Computed,
											String:                   &codespec.StringAttribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
										},
									},
								},
							},
							ReqBodyUsage: codespec.OmitAlways,
							Description:  conversion.StringPtr(testFieldDesc),
						},
						{
							Name:                     "single_nested_attr_with_nested_maps",
							PascalCaseName:           "SingleNestedAttrWithNestedMaps",
							ComputedOptionalRequired: codespec.Computed,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "map_attr1",
											PascalCaseName:           "MapAttr1",
											ComputedOptionalRequired: codespec.Computed,
											Map: &codespec.MapAttribute{
												ElementType: codespec.String,
											},
											ReqBodyUsage: codespec.OmitAlways,
										},
										{
											Name:                     "map_attr2",
											PascalCaseName:           "MapAttr2",
											ComputedOptionalRequired: codespec.Computed,
											Map: &codespec.MapAttribute{
												ElementType: codespec.String,
											},
											ReqBodyUsage: codespec.OmitAlways,
										},
									},
								},
							},
							ReqBodyUsage: codespec.OmitAlways,
							Description:  conversion.StringPtr(testFieldDesc),
						},
					},
				},
				Name: "test_resource_with_nested_attr",
				Operations: codespec.APIOperations{
					Create: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "POST",
					},
					Read: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "GET",
					},
					Update: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2024-05-30+json",
				},
			}},
		},
	}
	runTestCase(t, tc)
}

func TestConvertToProviderSpec_nested_schemaOverrides(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      "testdata/config-nested-schema-overrides.yml",
		inputResourceName:    "test_resource_with_nested_attr_overrides",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr(testResourceDesc),
					Attributes: codespec.Attributes{
						{
							Name:                     "attr_always_in_updates",
							PascalCaseName:           "AttrAlwaysInUpdates",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Always in updates"),
							ReqBodyUsage:             codespec.IncludeNullOnUpdate,
						},
						{
							Name:                     "project_id",
							PascalCaseName:           "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "nested_list_array_attr",
							PascalCaseName:           "NestedListArrayAttr",
							ComputedOptionalRequired: codespec.Required,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "inner_num_attr_alias",
											PascalCaseName:           "InnerNumAttr",
											ComputedOptionalRequired: codespec.Required,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr("Overridden inner_num_attr_alias description"),
										},
										{
											Name:                     "list_primitive_string_computed_attr",
											PascalCaseName:           "ListPrimitiveStringComputedAttr",
											ComputedOptionalRequired: codespec.Computed,
											List: &codespec.ListAttribute{
												ElementType: codespec.String,
											},
											Description:  conversion.StringPtr(testFieldDesc),
											ReqBodyUsage: codespec.OmitAlways,
										},
									},
								},
							},
							Description: conversion.StringPtr(testFieldDesc),
						},
						{
							Name:                     "optional_string_attr",
							PascalCaseName:           "OptionalStringAttr",
							ComputedOptionalRequired: codespec.ComputedOptional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Optional string that has config override to optional/computed"),
						},
						{
							Name:                     "outer_object",
							PascalCaseName:           "OuterObject",
							ComputedOptionalRequired: codespec.Computed,
							ReqBodyUsage:             codespec.OmitAlways,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "nested_level1",
											PascalCaseName:           "NestedLevel1",
											ComputedOptionalRequired: codespec.Computed,
											ReqBodyUsage:             codespec.OmitAlways,
											SingleNested: &codespec.SingleNestedAttribute{
												NestedObject: codespec.NestedAttributeObject{
													Attributes: codespec.Attributes{
														{
															Name:                     "level_field1_alias",
															PascalCaseName:           "LevelField1",
															ComputedOptionalRequired: codespec.Computed,
															ReqBodyUsage:             codespec.OmitAlways,
															String:                   &codespec.StringAttribute{},
															Description:              conversion.StringPtr("Overridden level_field1_alias description"),
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Name: "timeouts",
							Timeouts: &codespec.TimeoutsAttribute{
								ConfigurableTimeouts: []codespec.Operation{codespec.Create, codespec.Read, codespec.Update, codespec.Delete},
							},
							ReqBodyUsage: codespec.OmitAlways,
						},
					},
				},
				Name: "test_resource_with_nested_attr_overrides",
				Operations: codespec.APIOperations{
					Create: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "POST",
					},
					Read: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "GET",
					},
					Update: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2035-01-01+json", // version header defined in config
				},
			},
			},
		},
	}
	runTestCase(t, tc)
}

func TestConvertToProviderSpec_pathParamPresentInPostRequest(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_path_param_in_post_req",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr(testResourceDesc),
					Attributes: codespec.Attributes{
						{
							Name:                     "group_id",
							PascalCaseName:           "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "special_param",
							PascalCaseName:           "SpecialParam",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitInUpdateBody,
							Description:              conversion.StringPtr(testPathParamDesc),
						},
						{
							Name:                     "str_req_attr1",
							PascalCaseName:           "StrReqAttr1",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
						},
					},
				},
				Name: "test_resource_path_param_in_post_req",
				Operations: codespec.APIOperations{
					Create: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/pathparaminpostreq",
						HTTPMethod: "POST",
					},
					Read: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/pathparaminpostreq/{specialParam}",
						HTTPMethod: "GET",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/pathparaminpostreq/{specialParam}",
						HTTPMethod: "DELETE",
					},
					Update: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/pathparaminpostreq/{specialParam}",
						HTTPMethod: "PATCH",
					},
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
			},
			},
		},
	}
	runTestCase(t, tc)
}

func TestConvertToProviderSpec_singletonResourceNoDeleteOperation(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_singleton_resource_no_delete_op",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr("PATCH API description"),
					Attributes: codespec.Attributes{
						{
							Name:                     "flag",
							PascalCaseName:           "Flag",
							ComputedOptionalRequired: codespec.Optional,
							Bool:                     &codespec.BoolAttribute{},
						},
						{
							Name:                     "group_id",
							PascalCaseName:           "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
					},
				},
				Name: "test_singleton_resource_no_delete_op",
				Operations: codespec.APIOperations{
					Create: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testSingletonResource",
						HTTPMethod: "PATCH",
					},
					Read: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testSingletonResource",
						HTTPMethod: "GET",
					},
					Update: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testSingletonResource",
						HTTPMethod: "PATCH",
					},
					Delete:        nil,
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
			}},
		},
	}
	runTestCase(t, tc)
}

func runTestCase(t *testing.T, tc convertToSpecTestCase) {
	t.Helper()
	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}
