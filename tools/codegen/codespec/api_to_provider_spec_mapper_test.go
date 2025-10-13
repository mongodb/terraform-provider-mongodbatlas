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
							TFSchemaName:             "bool_default_attr",
							TFModelName:              "BoolDefaultAttr",
							ComputedOptionalRequired: codespec.ComputedOptional,
							Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "count",
							TFModelName:              "Count",
							ComputedOptionalRequired: codespec.Optional,
							Int64:                    &codespec.Int64Attribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "create_date",
							TFModelName:              "CreateDate",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Computed,
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "num_double_default_attr",
							TFModelName:              "NumDoubleDefaultAttr",
							Float64:                  &codespec.Float64Attribute{Default: conversion.Pointer(2.0)},
							ComputedOptionalRequired: codespec.ComputedOptional,
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "str_computed_attr",
							TFModelName:              "StrComputedAttr",
							ComputedOptionalRequired: codespec.Computed,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							TFSchemaName:             "str_req_attr1",
							TFModelName:              "StrReqAttr1",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "str_req_attr2",
							TFModelName:              "StrReqAttr2",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "str_req_attr3",
							TFModelName:              "StrReqAttr3",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Required,
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
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
							TFSchemaName:             "attr_always_in_updates",
							TFModelName:              "AttrAlwaysInUpdates",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Always in updates"),
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "cluster_name",
							TFModelName:              "ClusterName",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "list_primitive_string_attr",
							TFModelName:              "ListPrimitiveStringAttr",
							ComputedOptionalRequired: codespec.Computed,
							List: &codespec.ListAttribute{
								ElementType: codespec.String,
							},
							Description:  conversion.StringPtr(testFieldDesc),
							ReqBodyUsage: codespec.OmitAlways,
						},
						{
							TFSchemaName:             "nested_list_array_attr",
							TFModelName:              "NestedListArrayAttr",
							ComputedOptionalRequired: codespec.Required,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_num_attr",
											TFModelName:              "InnerNumAttr",
											ComputedOptionalRequired: codespec.Required,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.AllRequestBodies,
										},
										{
											TFSchemaName:             "list_primitive_string_attr",
											TFModelName:              "ListPrimitiveStringAttr",
											ComputedOptionalRequired: codespec.Optional,
											List: &codespec.ListAttribute{
												ElementType: codespec.String,
											},
											Description:  conversion.StringPtr(testFieldDesc),
											ReqBodyUsage: codespec.AllRequestBodies,
										},
										{
											TFSchemaName:             "list_primitive_string_computed_attr",
											TFModelName:              "ListPrimitiveStringComputedAttr",
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
							Description:  conversion.StringPtr(testFieldDesc),
							ReqBodyUsage: codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "nested_map_object_attr",
							TFModelName:              "NestedMapObjectAttr",
							ComputedOptionalRequired: codespec.Computed,
							MapNested: &codespec.MapNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "attr",
											TFModelName:              "Attr",
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
							TFSchemaName:             "nested_set_array_attr",
							TFModelName:              "NestedSetArrayAttr",
							ComputedOptionalRequired: codespec.Computed,
							SetNested: &codespec.SetNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_num_attr",
											TFModelName:              "InnerNumAttr",
											ComputedOptionalRequired: codespec.Computed,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
										},
										{
											TFSchemaName:             "list_primitive_string_attr",
											TFModelName:              "ListPrimitiveStringAttr",
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
							TFSchemaName:             "optional_string_attr",
							TFModelName:              "OptionalStringAttr",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Optional string"),
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "set_primitive_string_attr",
							TFModelName:              "SetPrimitiveStringAttr",
							ComputedOptionalRequired: codespec.Computed,
							Set: &codespec.SetAttribute{
								ElementType: codespec.String,
							},
							Description:  conversion.StringPtr(testFieldDesc),
							ReqBodyUsage: codespec.OmitAlways,
						},
						{
							TFSchemaName:             "single_nested_attr",
							TFModelName:              "SingleNestedAttr",
							ComputedOptionalRequired: codespec.Computed,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_int_attr",
											TFModelName:              "InnerIntAttr",
											ComputedOptionalRequired: codespec.Computed,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
										},
										{
											TFSchemaName:             "inner_str_attr",
											TFModelName:              "InnerStrAttr",
											ComputedOptionalRequired: codespec.Computed,
											String:                   &codespec.StringAttribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
										},
									},
								},
							},
							Description:  conversion.StringPtr(testFieldDesc),
							ReqBodyUsage: codespec.OmitAlways,
						},
						{
							TFSchemaName:             "single_nested_attr_with_nested_maps",
							TFModelName:              "SingleNestedAttrWithNestedMaps",
							ComputedOptionalRequired: codespec.Computed,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "map_attr1",
											TFModelName:              "MapAttr1",
											ComputedOptionalRequired: codespec.Computed,
											Map: &codespec.MapAttribute{
												ElementType: codespec.String,
											},
											ReqBodyUsage: codespec.OmitAlways,
										},
										{
											TFSchemaName:             "map_attr2",
											TFModelName:              "MapAttr2",
											ComputedOptionalRequired: codespec.Computed,
											Map: &codespec.MapAttribute{
												ElementType: codespec.String,
											},
											ReqBodyUsage: codespec.OmitAlways,
										},
									},
								},
							},
							Description:  conversion.StringPtr(testFieldDesc),
							ReqBodyUsage: codespec.OmitAlways,
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
							TFSchemaName:             "attr_always_in_updates",
							TFModelName:              "AttrAlwaysInUpdates",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Always in updates"),
							ReqBodyUsage:             codespec.IncludeNullOnUpdate,
						},
						{
							TFSchemaName:             "project_id",
							TFModelName:              "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "nested_list_array_attr",
							TFModelName:              "NestedListArrayAttr",
							ComputedOptionalRequired: codespec.Required,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_num_attr_alias",
											TFModelName:              "InnerNumAttr",
											ComputedOptionalRequired: codespec.Required,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr("Overridden inner_num_attr_alias description"),
											ReqBodyUsage:             codespec.AllRequestBodies,
										},
										{
											TFSchemaName:             "list_primitive_string_computed_attr",
											TFModelName:              "ListPrimitiveStringComputedAttr",
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
							Description:  conversion.StringPtr(testFieldDesc),
							ReqBodyUsage: codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "optional_string_attr",
							TFModelName:              "OptionalStringAttr",
							ComputedOptionalRequired: codespec.ComputedOptional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Optional string that has config override to optional/computed"),
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "outer_object",
							TFModelName:              "OuterObject",
							ComputedOptionalRequired: codespec.Computed,
							ReqBodyUsage:             codespec.OmitAlways,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "nested_level1",
											TFModelName:              "NestedLevel1",
											ComputedOptionalRequired: codespec.Computed,
											ReqBodyUsage:             codespec.OmitAlways,
											SingleNested: &codespec.SingleNestedAttribute{
												NestedObject: codespec.NestedAttributeObject{
													Attributes: codespec.Attributes{
														{
															TFSchemaName:             "level_field1_alias",
															TFModelName:              "LevelField1",
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
							TFSchemaName: "timeouts",
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
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "special_param",
							TFModelName:              "SpecialParam",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitInUpdateBody,
							Description:              conversion.StringPtr(testPathParamDesc),
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "str_req_attr1",
							TFModelName:              "StrReqAttr1",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
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
							TFSchemaName:             "flag",
							TFModelName:              "Flag",
							ComputedOptionalRequired: codespec.Optional,
							Bool:                     &codespec.BoolAttribute{},
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
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
