package codespec_test

import (
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
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
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "group_id",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
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
							ReqBodyUsage:             codespec.OmitAlways,
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
					Delete: codespec.APIOperation{
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
							Name:                     "cluster_name",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "group_id",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "list_primitive_string_attr",
							ComputedOptionalRequired: codespec.Computed,
							List: &codespec.ListAttribute{
								ElementType: codespec.String,
							},
							Description:  conversion.StringPtr(testFieldDesc),
							ReqBodyUsage: codespec.OmitAlways,
						},
						{
							Name:                     "nested_list_array_attr",
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
										{
											Name:                     "list_primitive_string_computed_attr",
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
							ComputedOptionalRequired: codespec.Computed,
							MapNested: &codespec.MapNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "attr",
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
							ComputedOptionalRequired: codespec.Computed,
							SetNested: &codespec.SetNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "inner_num_attr",
											ComputedOptionalRequired: codespec.Computed,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
										},
										{
											Name:                     "list_primitive_string_attr",
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
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Optional string"),
						},
						{
							Name:                     "set_primitive_string_attr",
							ComputedOptionalRequired: codespec.Computed,
							Set: &codespec.SetAttribute{
								ElementType: codespec.String,
							},
							ReqBodyUsage: codespec.OmitAlways,
							Description:  conversion.StringPtr(testFieldDesc),
						},
						{
							Name:                     "single_nested_attr",
							ComputedOptionalRequired: codespec.Computed,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "inner_int_attr",
											ComputedOptionalRequired: codespec.Computed,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
										},
										{
											Name:                     "inner_str_attr",
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
							ComputedOptionalRequired: codespec.Computed,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "map_attr1",
											ComputedOptionalRequired: codespec.Computed,
											Map: &codespec.MapAttribute{
												ElementType: codespec.String,
											},
											ReqBodyUsage: codespec.OmitAlways,
										},
										{
											Name:                     "map_attr2",
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
					Delete: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2024-05-30+json",
				},
			},
			},
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
							Name:                     "project_id",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "nested_list_array_attr",
							ComputedOptionalRequired: codespec.Required,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "inner_num_attr_alias",
											ComputedOptionalRequired: codespec.Required,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr("Overridden inner_num_attr_alias description"),
										},
										{
											Name:                     "list_primitive_string_computed_attr",
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
							ComputedOptionalRequired: codespec.ComputedOptional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Optional string that has config override to optional/computed"),
						},
						{
							Name:                     "outer_object",
							ComputedOptionalRequired: codespec.Computed,
							ReqBodyUsage:             codespec.OmitAlways,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											Name:                     "nested_level1",
											ComputedOptionalRequired: codespec.Computed,
											ReqBodyUsage:             codespec.OmitAlways,
											SingleNested: &codespec.SingleNestedAttribute{
												NestedObject: codespec.NestedAttributeObject{
													Attributes: codespec.Attributes{
														{
															Name:                     "level_field1_alias",
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
					Delete: codespec.APIOperation{
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
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							Name:                     "special_param",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitInUpdateBody,
							Description:              conversion.StringPtr(testPathParamDesc),
						},
						{
							Name:                     "str_req_attr1",
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
					Delete: codespec.APIOperation{
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

var streamConnectionExpectedSchema = codespec.Schema{
	Description: conversion.StringPtr("create sc"),
	Attributes: []codespec.Attribute{
		{
			Name:                     "group_id",
			ComputedOptionalRequired: codespec.Required,
			String:                   &codespec.StringAttribute{},
			Description:              conversion.StringPtr("Unique 24-hexadecimal digit"),
			ReqBodyUsage:             codespec.OmitAlways,
		},
		{
			Name:                     "name",
			ComputedOptionalRequired: codespec.Optional, // inferred wrongly
			String:                   &codespec.StringAttribute{},
			Description:              conversion.StringPtr("Human-readable label that identifies the stream connection. In the case of the Sample type, this is the name of the sample source."),
		},
		{
			Name:                     "tenant_name",
			ComputedOptionalRequired: codespec.Required,
			String:                   &codespec.StringAttribute{},
			Description:              conversion.StringPtr("Human-readable label that identifies the stream instance."),
			ReqBodyUsage:             codespec.OmitAlways,
		},
		{
			Name:                     "type",
			String:                   &codespec.StringAttribute{},
			ComputedOptionalRequired: codespec.Computed,
			Description:              conversion.Pointer("Type of the connection."),
		},
		{
			Name:                     "type_cluster",
			ComputedOptionalRequired: codespec.Optional,
			SingleNested: &codespec.SingleNestedAttribute{
				NestedObject: codespec.NestedAttributeObject{
					Attributes: []codespec.Attribute{
						{
							Name:                     "cluster_name",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Optional, // this is inferred as optional, but it should be required
							Description:              conversion.Pointer("Name of the cluster configured for this connection."),
						},
					},
				},
			},
			Discriminator: &codespec.DiscriminatorMapping{
				DiscriminatorProperty: "type",
				DiscriminatorValue:    "Cluster",
			},
		},
		{
			Name:                     "type_https",
			ComputedOptionalRequired: codespec.Optional,
			SingleNested: &codespec.SingleNestedAttribute{
				NestedObject: codespec.NestedAttributeObject{
					Attributes: []codespec.Attribute{
						{
							Name: "headers",
							Map: &codespec.MapAttribute{
								ElementType: codespec.String,
							},
							Description:              conversion.StringPtr("A map of key-value pairs that will be passed as headers for the request."),
							ComputedOptionalRequired: codespec.Optional,
						},
						{
							Name:                     "url",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Optional, // this is inferred as optional, but it should be required
							Description:              conversion.Pointer("The url to be used for the request."),
						},
					},
				},
			},
			Discriminator: &codespec.DiscriminatorMapping{
				DiscriminatorProperty: "type",
				DiscriminatorValue:    "Https",
			},
		},
	},
}

func TestConvertToProviderSpec_discriminatorStreamConnection(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: "testdata/api-spec-stream_connection.yaml",
		inputConfigPath:      "testdata/config-stream_connection.yml",
		inputResourceName:    "stream_connection",
		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &streamConnectionExpectedSchema,
				Name:   "stream_connection",
				Operations: codespec.APIOperations{
					Create: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/connections",
						HTTPMethod: "POST",
					},
					Read: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/connections/{connectionName}",
						HTTPMethod: "GET",
					},
					Update: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/connections/{connectionName}",
						HTTPMethod: "PATCH",
					},
					Delete: codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/streams/{tenantName}/connections/{connectionName}",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2023-02-01+json",
				},
			},
			},
		},
	}
	runTestCase(t, tc)
}

func runTestCase(t *testing.T, tc convertToSpecTestCase) {
	t.Helper()
	apiSpecParsed, err := openapi.ParseAtlasAdminAPI(tc.inputOpenAPISpecPath)
	require.NoError(t, err, "Failed to parse OpenAPI spec")
	configModel, err := config.ParseGenConfigYAML(tc.inputConfigPath)
	require.NoError(t, err, "Failed to parse config file")
	configModel.APISpecs = []config.APISpec{
		{
			Name:      "test_api_spec",
			URL:       "Ignored",
			IsDefault: conversion.Pointer(true),
		},
	}
	parsedAPISpecs := map[string]*libopenapi.DocumentModel[v3.Document]{
		"test_api_spec": apiSpecParsed,
	}
	result, err := codespec.ToCodeSpecModel(parsedAPISpecs, configModel, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}
