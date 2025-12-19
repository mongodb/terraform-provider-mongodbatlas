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

var (
	simpleTestResourceAttributes = codespec.Attributes{
		{
			TFSchemaName:             "group_id",
			TFModelName:              "GroupId",
			APIName:                  "groupId",
			ComputedOptionalRequired: codespec.Required,
			String:                   &codespec.StringAttribute{},
			Description:              conversion.StringPtr(testPathParamDesc),
			ReqBodyUsage:             codespec.OmitAlways,
			CreateOnly:               true,
		},
		{
			TFSchemaName:             "string_attr",
			TFModelName:              "StringAttr",
			APIName:                  "stringAttr",
			ComputedOptionalRequired: codespec.Required,
			String:                   &codespec.StringAttribute{},
			Description:              conversion.StringPtr(testFieldDesc),
			ReqBodyUsage:             codespec.AllRequestBodies,
			PresentInAnyResponse:     true,
		},
	}

	simpleTestResourceOperations = codespec.APIOperations{
		Create: &codespec.APIOperation{
			Path:       "/api/atlas/v2/groups/{groupId}/simpleTestResource",
			HTTPMethod: "POST",
		},
		Read: &codespec.APIOperation{
			Path:       "/api/atlas/v2/groups/{groupId}/simpleTestResource",
			HTTPMethod: "GET",
		},
		Update: &codespec.APIOperation{
			Path:       "/api/atlas/v2/groups/{groupId}/simpleTestResource",
			HTTPMethod: "PATCH",
		},
		Delete: &codespec.APIOperation{
			Path:       "/api/atlas/v2/groups/{groupId}/simpleTestResource",
			HTTPMethod: "DELETE",
		},
		VersionHeader: "application/vnd.atlas.2023-01-01+json",
	}
)

type convertToSpecTestCase struct {
	expectedResult       *codespec.Model
	inputOpenAPISpecPath string
	inputConfigPath      string
	inputResourceName    string
}

// getTestResourceComputedAttributes returns the common TestResource attributes used in data sources.
// All attributes are Computed as appropriate for data source response fields.
func getTestResourceComputedAttributes() codespec.Attributes {
	return codespec.Attributes{
		{
			TFSchemaName:             "bool_default_attr",
			TFModelName:              "BoolDefaultAttr",
			APIName:                  "boolDefaultAttr",
			ComputedOptionalRequired: codespec.Computed,
			Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
			ReqBodyUsage:             codespec.OmitAlways,
		},
		{
			TFSchemaName:             "count",
			TFModelName:              "Count",
			APIName:                  "count",
			ComputedOptionalRequired: codespec.Computed,
			Int64:                    &codespec.Int64Attribute{},
			Description:              conversion.StringPtr(testFieldDesc),
			ReqBodyUsage:             codespec.OmitAlways,
		},
		{
			TFSchemaName:             "create_date",
			TFModelName:              "CreateDate",
			APIName:                  "createDate",
			String:                   &codespec.StringAttribute{},
			ComputedOptionalRequired: codespec.Computed,
			Description:              conversion.StringPtr(testFieldDesc),
			ReqBodyUsage:             codespec.OmitAlways,
		},
		{
			TFSchemaName:             "num_double_default_attr",
			TFModelName:              "NumDoubleDefaultAttr",
			APIName:                  "numDoubleDefaultAttr",
			Float64:                  &codespec.Float64Attribute{Default: conversion.Pointer(2.0)},
			ComputedOptionalRequired: codespec.Computed,
			ReqBodyUsage:             codespec.OmitAlways,
		},
		{
			TFSchemaName:             "str_computed_attr",
			TFModelName:              "StrComputedAttr",
			APIName:                  "strComputedAttr",
			ComputedOptionalRequired: codespec.Computed,
			String:                   &codespec.StringAttribute{},
			Description:              conversion.StringPtr(testFieldDesc),
			ReqBodyUsage:             codespec.OmitAlways,
		},
		{
			TFSchemaName:             "str_req_attr1",
			TFModelName:              "StrReqAttr1",
			APIName:                  "strReqAttr1",
			ComputedOptionalRequired: codespec.Computed,
			String:                   &codespec.StringAttribute{},
			Description:              conversion.StringPtr(testFieldDesc),
			ReqBodyUsage:             codespec.OmitAlways,
		},
		{
			TFSchemaName:             "str_req_attr2",
			TFModelName:              "StrReqAttr2",
			APIName:                  "strReqAttr2",
			ComputedOptionalRequired: codespec.Computed,
			String:                   &codespec.StringAttribute{},
			Description:              conversion.StringPtr(testFieldDesc),
			ReqBodyUsage:             codespec.OmitAlways,
		},
		{
			TFSchemaName:             "str_req_attr3",
			TFModelName:              "StrReqAttr3",
			APIName:                  "strReqAttr3",
			String:                   &codespec.StringAttribute{},
			ComputedOptionalRequired: codespec.Computed,
			Description:              conversion.StringPtr(testFieldDesc),
			ReqBodyUsage:             codespec.OmitAlways,
		},
	}
}

// getProjectIDPathParam returns the path parameter (aliased from groupId to projectId) for data sources.
func getProjectIDPathParam() codespec.Attribute {
	return codespec.Attribute{
		TFSchemaName:             "project_id",
		TFModelName:              "ProjectId",
		APIName:                  "groupId",
		ComputedOptionalRequired: codespec.Required,
		String:                   &codespec.StringAttribute{},
		Description:              conversion.StringPtr(testPathParamDesc),
		ReqBodyUsage:             codespec.OmitAlways,
	}
}

// getPluralDSAttributes returns the attributes for plural data sources (with results list).
func getPluralDSAttributes() codespec.Attributes {
	projectIDPathParam := getProjectIDPathParam()
	testResourceComputedAttributes := getTestResourceComputedAttributes()

	return codespec.Attributes{
		projectIDPathParam,
		{
			TFSchemaName:             "results",
			TFModelName:              "Results",
			APIName:                  "results",
			ComputedOptionalRequired: codespec.Computed,
			CustomType:               codespec.NewCustomNestedListType("Results"),
			ListNested: &codespec.ListNestedAttribute{
				NestedObject: codespec.NestedAttributeObject{
					Attributes: testResourceComputedAttributes,
				},
			},
			ReqBodyUsage: codespec.OmitAlways,
		},
	}
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
							APIName:                  "boolDefaultAttr",
							ComputedOptionalRequired: codespec.ComputedOptional,
							Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "count",
							TFModelName:              "Count",
							APIName:                  "count",
							ComputedOptionalRequired: codespec.Optional,
							Int64:                    &codespec.Int64Attribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "create_date",
							TFModelName:              "CreateDate",
							APIName:                  "createDate",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Computed,
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "num_double_default_attr",
							TFModelName:              "NumDoubleDefaultAttr",
							APIName:                  "numDoubleDefaultAttr",
							Float64:                  &codespec.Float64Attribute{Default: conversion.Pointer(2.0)},
							ComputedOptionalRequired: codespec.ComputedOptional,
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "str_computed_attr",
							TFModelName:              "StrComputedAttr",
							APIName:                  "strComputedAttr",
							ComputedOptionalRequired: codespec.Computed,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "str_req_attr1",
							TFModelName:              "StrReqAttr1",
							APIName:                  "strReqAttr1",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "str_req_attr2",
							TFModelName:              "StrReqAttr2",
							APIName:                  "strReqAttr2",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "str_req_attr3",
							TFModelName:              "StrReqAttr3",
							APIName:                  "strReqAttr3",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Required,
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
					},
				},
				Name:        "test_resource_no_schema_opts",
				PackageName: "testresourcenoschemaopts",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
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

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
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
							APIName:                  "attrAlwaysInUpdates",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Always in updates"),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "cluster_name",
							TFModelName:              "ClusterName",
							APIName:                  "clusterName",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "list_primitive_string_attr",
							TFModelName:              "ListPrimitiveStringAttr",
							APIName:                  "listPrimitiveStringAttr",
							ComputedOptionalRequired: codespec.Computed,
							CustomType:               codespec.NewCustomListType(codespec.String),
							List: &codespec.ListAttribute{
								ElementType: codespec.String,
							},
							Description:          conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:         codespec.OmitAlways,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "nested_list_array_attr",
							TFModelName:              "NestedListArrayAttr",
							APIName:                  "nestedListArrayAttr",
							ComputedOptionalRequired: codespec.Required,
							CustomType:               codespec.NewCustomNestedListType("NestedListArrayAttr"),
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_num_attr",
											TFModelName:              "InnerNumAttr",
											APIName:                  "innerNumAttr",
											ComputedOptionalRequired: codespec.Required,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.AllRequestBodies,
											PresentInAnyResponse:     true,
										},
										{
											TFSchemaName:             "list_primitive_string_attr",
											TFModelName:              "ListPrimitiveStringAttr",
											APIName:                  "listPrimitiveStringAttr",
											ComputedOptionalRequired: codespec.Optional,
											CustomType:               codespec.NewCustomListType(codespec.String),
											List: &codespec.ListAttribute{
												ElementType: codespec.String,
											},
											Description:          conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:         codespec.AllRequestBodies,
											PresentInAnyResponse: true,
										},
										{
											TFSchemaName:             "list_primitive_string_computed_attr",
											TFModelName:              "ListPrimitiveStringComputedAttr",
											APIName:                  "listPrimitiveStringComputedAttr",
											ComputedOptionalRequired: codespec.Computed,
											CustomType:               codespec.NewCustomListType(codespec.String),
											List: &codespec.ListAttribute{
												ElementType: codespec.String,
											},
											Description:          conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:         codespec.OmitAlways,
											PresentInAnyResponse: true,
										},
									},
								},
							},
							Description:          conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:         codespec.AllRequestBodies,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "nested_map_object_attr",
							TFModelName:              "NestedMapObjectAttr",
							APIName:                  "nestedMapObjectAttr",
							ComputedOptionalRequired: codespec.Computed,
							CustomType:               codespec.NewCustomNestedMapType("NestedMapObjectAttr"),
							MapNested: &codespec.MapNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "attr",
											TFModelName:              "Attr",
											APIName:                  "attr",
											ComputedOptionalRequired: codespec.Computed,
											String:                   &codespec.StringAttribute{},
											ReqBodyUsage:             codespec.OmitAlways,
											PresentInAnyResponse:     true,
										},
									},
								},
							},
							ReqBodyUsage:         codespec.OmitAlways,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "nested_set_array_attr",
							TFModelName:              "NestedSetArrayAttr",
							APIName:                  "nestedSetArrayAttr",
							ComputedOptionalRequired: codespec.Computed,
							CustomType:               codespec.NewCustomNestedSetType("NestedSetArrayAttr"),
							SetNested: &codespec.SetNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_num_attr",
											TFModelName:              "InnerNumAttr",
											APIName:                  "innerNumAttr",
											ComputedOptionalRequired: codespec.Computed,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
											PresentInAnyResponse:     true,
										},
										{
											TFSchemaName:             "list_primitive_string_attr",
											TFModelName:              "ListPrimitiveStringAttr",
											APIName:                  "listPrimitiveStringAttr",
											ComputedOptionalRequired: codespec.Computed,
											CustomType:               codespec.NewCustomListType(codespec.String),
											List: &codespec.ListAttribute{
												ElementType: codespec.String,
											},
											Description:          conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:         codespec.OmitAlways,
											PresentInAnyResponse: true,
										},
									},
								},
							},
							ReqBodyUsage:         codespec.OmitAlways,
							Description:          conversion.StringPtr(testFieldDesc),
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "optional_string_attr",
							TFModelName:              "OptionalStringAttr",
							APIName:                  "optionalStringAttr",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Optional string"),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "set_primitive_string_attr",
							TFModelName:              "SetPrimitiveStringAttr",
							APIName:                  "setPrimitiveStringAttr",
							ComputedOptionalRequired: codespec.Computed,
							CustomType:               codespec.NewCustomSetType(codespec.String),
							Set: &codespec.SetAttribute{
								ElementType: codespec.String,
							},
							Description:          conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:         codespec.OmitAlways,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "single_nested_attr",
							TFModelName:              "SingleNestedAttr",
							APIName:                  "singleNestedAttr",
							ComputedOptionalRequired: codespec.Computed,
							CustomType:               codespec.NewCustomObjectType("SingleNestedAttr"),
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_int_attr",
											TFModelName:              "InnerIntAttr",
											APIName:                  "innerIntAttr",
											ComputedOptionalRequired: codespec.Computed,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
											PresentInAnyResponse:     true,
										},
										{
											TFSchemaName:             "inner_str_attr",
											TFModelName:              "InnerStrAttr",
											APIName:                  "innerStrAttr",
											ComputedOptionalRequired: codespec.Computed,
											String:                   &codespec.StringAttribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
											PresentInAnyResponse:     true,
										},
									},
								},
							},
							Description:          conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:         codespec.OmitAlways,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "single_nested_attr_with_nested_maps",
							TFModelName:              "SingleNestedAttrWithNestedMaps",
							APIName:                  "singleNestedAttrWithNestedMaps",
							ComputedOptionalRequired: codespec.Computed,
							CustomType:               codespec.NewCustomObjectType("SingleNestedAttrWithNestedMaps"),
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "map_attr1",
											TFModelName:              "MapAttr1",
											APIName:                  "mapAttr1",
											ComputedOptionalRequired: codespec.Computed,
											CustomType:               codespec.NewCustomMapType(codespec.String),
											Map: &codespec.MapAttribute{
												ElementType: codespec.String,
											},
											ReqBodyUsage:         codespec.OmitAlways,
											PresentInAnyResponse: true,
										},
										{
											TFSchemaName:             "map_attr2",
											TFModelName:              "MapAttr2",
											APIName:                  "mapAttr2",
											ComputedOptionalRequired: codespec.Computed,
											CustomType:               codespec.NewCustomMapType(codespec.String),
											Map: &codespec.MapAttribute{
												ElementType: codespec.String,
											},
											ReqBodyUsage:         codespec.OmitAlways,
											PresentInAnyResponse: true,
										},
									},
								},
							},
							Description:          conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:         codespec.OmitAlways,
							PresentInAnyResponse: true,
						},
					},
				},
				Name:        "test_resource_with_nested_attr",
				PackageName: "testresourcewithnestedattr",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
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

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
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
							APIName:                  "attrAlwaysInUpdates",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Always in updates"),
							ReqBodyUsage:             codespec.IncludeNullOnUpdate,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "project_id",
							TFModelName:              "ProjectId", // TFModelName changed by alias
							APIName:                  "groupId",   // Original API name preserved for apiname tag
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "nested_list_array_attr",
							TFModelName:              "NestedListArrayAttr",
							APIName:                  "nestedListArrayAttr",
							ComputedOptionalRequired: codespec.Required,
							CustomType:               codespec.NewCustomNestedListType("NestedListArrayAttr"),
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_num_attr_alias",
											TFModelName:              "InnerNumAttrAlias", // TFModelName changed by alias
											APIName:                  "innerNumAttr",      // Original API name preserved for apiname tag
											ComputedOptionalRequired: codespec.Required,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr("Overridden inner_num_attr_alias description"),
											ReqBodyUsage:             codespec.AllRequestBodies,
											PresentInAnyResponse:     true,
										},
										{
											TFSchemaName:             "list_primitive_string_computed_attr",
											TFModelName:              "ListPrimitiveStringComputedAttr",
											APIName:                  "listPrimitiveStringComputedAttr",
											ComputedOptionalRequired: codespec.Computed,
											CustomType:               codespec.NewCustomListType(codespec.String),
											List: &codespec.ListAttribute{
												ElementType: codespec.String,
											},
											Description:          conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:         codespec.OmitAlways,
											PresentInAnyResponse: true,
										},
									},
								},
							},
							Description:          conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:         codespec.AllRequestBodies,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "optional_string_attr",
							TFModelName:              "OptionalStringAttr",
							APIName:                  "optionalStringAttr",
							ComputedOptionalRequired: codespec.ComputedOptional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Optional string that has config override to optional/computed"),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "outer_object",
							TFModelName:              "OuterObject",
							APIName:                  "outerObject",
							ComputedOptionalRequired: codespec.Computed,
							ReqBodyUsage:             codespec.OmitAlways,
							PresentInAnyResponse:     true,
							CustomType:               codespec.NewCustomObjectType("OuterObject"),
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "nested_level1",
											TFModelName:              "NestedLevel1",
											APIName:                  "nestedLevel1",
											ComputedOptionalRequired: codespec.Computed,
											ReqBodyUsage:             codespec.OmitAlways,
											PresentInAnyResponse:     true,
											CustomType:               codespec.NewCustomObjectType("OuterObjectNestedLevel1"),
											SingleNested: &codespec.SingleNestedAttribute{
												NestedObject: codespec.NestedAttributeObject{
													Attributes: codespec.Attributes{
														{
															TFSchemaName:             "level_field1_alias",
															TFModelName:              "LevelField1Alias", // TFModelName changed by alias
															APIName:                  "levelField1",      // Original API name preserved for apiname tag
															ComputedOptionalRequired: codespec.Computed,
															ReqBodyUsage:             codespec.OmitAlways,
															PresentInAnyResponse:     true,
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
					},
				},
				Name:        "test_resource_with_nested_attr_overrides",
				PackageName: "testresourcewithnestedattroverrides",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
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

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}

// TestConvertToProviderSpec_pathBasedAlias verifies that aliases can target specific nested attributes
// using path-based keys (e.g., "nestedListArrayAttr.innerNumAttr") instead of just model names.
// This ensures that when multiple nested objects have attributes with the same name, only the targeted one is aliased.
func TestConvertToProviderSpec_pathBasedAlias(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      "testdata/config-path-based-alias.yml",
		inputResourceName:    "test_resource_with_nested_attr_path_alias",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr(testResourceDesc),
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "project_id",
							TFModelName:              "ProjectId", // TFModelName changed by alias
							APIName:                  "groupId",   // Original API name preserved for apiname tag
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "nested_list_array_attr",
							TFModelName:              "NestedListArrayAttr",
							APIName:                  "nestedListArrayAttr",
							ComputedOptionalRequired: codespec.Required,
							CustomType:               codespec.NewCustomNestedListType("NestedListArrayAttr"),
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											// This attribute should be ALIASED via path-based alias
											TFSchemaName:             "renamed_inner_num_attr",
											TFModelName:              "RenamedInnerNumAttr", // TFModelName changed by alias
											APIName:                  "innerNumAttr",        // Original API name preserved for apiname tag
											ComputedOptionalRequired: codespec.Required,
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.AllRequestBodies,
											PresentInAnyResponse:     true,
										},
									},
								},
							},
							Description:          conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:         codespec.AllRequestBodies,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "nested_set_array_attr",
							TFModelName:              "NestedSetArrayAttr",
							APIName:                  "nestedSetArrayAttr",
							ComputedOptionalRequired: codespec.Computed,
							CustomType:               codespec.NewCustomNestedSetType("NestedSetArrayAttr"),
							SetNested: &codespec.SetNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											// This attribute should NOT be aliased (different path)
											TFSchemaName:             "inner_num_attr",
											TFModelName:              "InnerNumAttr",
											APIName:                  "innerNumAttr",
											ComputedOptionalRequired: codespec.Computed, // computed because parent is readOnly
											Int64:                    &codespec.Int64Attribute{},
											Description:              conversion.StringPtr(testFieldDesc),
											ReqBodyUsage:             codespec.OmitAlways,
											PresentInAnyResponse:     true,
										},
									},
								},
							},
							Description:          conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:         codespec.OmitAlways,
							PresentInAnyResponse: true,
						},
					},
				},
				Name:        "test_resource_with_nested_attr_path_alias",
				PackageName: "testresourcewithnestedattrpathalias",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/nestedTestResource",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2024-05-30+json",
				},
			},
			},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
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
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:                "special_param",
							TFModelName:                 "SpecialParam",
							APIName:                     "specialParam",
							ComputedOptionalRequired:    codespec.Optional,
							String:                      &codespec.StringAttribute{},
							ReqBodyUsage:                codespec.OmitInUpdateBody,
							Description:                 conversion.StringPtr(testPathParamDesc),
							CreateOnly:                  true,
							RequestOnlyRequiredOnCreate: true,
						},
						{
							TFSchemaName:             "str_req_attr1",
							TFModelName:              "StrReqAttr1",
							APIName:                  "strReqAttr1",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
					},
				},
				Name:        "test_resource_path_param_in_post_req",
				PackageName: "testresourcepathparaminpostreq",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/pathparaminpostreq",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/pathparaminpostreq/{specialParam}",
						HTTPMethod: "GET",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/pathparaminpostreq/{specialParam}",
						HTTPMethod: "DELETE",
					},
					Update: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/pathparaminpostreq/{specialParam}",
						HTTPMethod: "PATCH",
					},
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
			},
			},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
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
							APIName:                  "flag",
							ComputedOptionalRequired: codespec.Optional,
							Bool:                     &codespec.BoolAttribute{},
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
					},
				},
				Name:        "test_singleton_resource_no_delete_op",
				PackageName: "testsingletonresourcenodeleteop",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testSingletonResource",
						HTTPMethod: "PATCH",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testSingletonResource",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testSingletonResource",
						HTTPMethod: "PATCH",
					},
					Delete:        nil,
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}

func TestConvertToProviderSpec_NoUpdateOperation(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_no_update_op",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr("POST API description"),
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "string_attr",
							TFModelName:              "StringAttr",
							APIName:                  "stringAttr",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitInUpdateBody,
							CreateOnly:               true,
							PresentInAnyResponse:     true,
						},
					},
				},
				Name:        "test_resource_no_update_op",
				PackageName: "testresourcenoupdateop",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceNoUpdate",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceNoUpdate",
						HTTPMethod: "GET",
					},
					Update: nil,
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceNoUpdate",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}

func TestConvertToProviderSpec_typeOverride(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_with_overridden_collection_types",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr(testResourceDesc),
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "flag",
							TFModelName:              "Flag",
							APIName:                  "flag",
							ComputedOptionalRequired: codespec.Required,
							Bool:                     &codespec.BoolAttribute{},
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "list_string",
							TFModelName:              "ListString",
							APIName:                  "listString",
							ComputedOptionalRequired: codespec.Required,
							// List overridden to set
							CustomType:           codespec.NewCustomSetType(codespec.String),
							Set:                  &codespec.SetAttribute{ElementType: codespec.String},
							ReqBodyUsage:         codespec.AllRequestBodies,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "set_string",
							TFModelName:              "SetString",
							APIName:                  "setString",
							ComputedOptionalRequired: codespec.Required,
							// Set overridden to list
							CustomType:           codespec.NewCustomListType(codespec.String),
							List:                 &codespec.ListAttribute{ElementType: codespec.String},
							ReqBodyUsage:         codespec.AllRequestBodies,
							PresentInAnyResponse: true,
						},
					},
				},
				Name:        "test_resource_with_overridden_collection_types",
				PackageName: "testresourcewithoverriddencollectiontypes",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceWithCollections",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceWithCollections",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceWithCollections",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceWithCollections",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}

func TestConvertToProviderSpec_dynamicJSONProperties(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_dynamic_json_properties",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr(testResourceDesc),
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "array_of_dynamic_values",
							TFModelName:              "ArrayOfDynamicValues",
							APIName:                  "arrayOfDynamicValues",
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomListType(codespec.CustomTypeJSON),
							List: &codespec.ListAttribute{
								ElementType: codespec.CustomTypeJSON,
							},
							Description:          conversion.StringPtr("Array of dynamic values."),
							ReqBodyUsage:         codespec.AllRequestBodies,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "dynamic_value",
							TFModelName:              "DynamicValue",
							APIName:                  "dynamicValue",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							CustomType:               &codespec.CustomTypeJSONVar,
							Description:              conversion.StringPtr("Dynamic value."),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "object_of_dynamic_values",
							TFModelName:              "ObjectOfDynamicValues",
							APIName:                  "objectOfDynamicValues",
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomMapType(codespec.CustomTypeJSON),
							Map: &codespec.MapAttribute{
								ElementType: codespec.CustomTypeJSON,
							},
							Description:          conversion.StringPtr("Object of dynamic values."),
							ReqBodyUsage:         codespec.AllRequestBodies,
							PresentInAnyResponse: true,
						},
					},
				},
				Name:        "test_dynamic_json_properties",
				PackageName: "testdynamicjsonproperties",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/dynamicJsonProperties",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/dynamicJsonProperties",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
						Path:       "/api/atlas/v2/dynamicJsonProperties",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/dynamicJsonProperties",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2024-05-30+json",
				},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}

func TestConvertToProviderSpec_moveState(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_move_state",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr("POST API description"),
					Attributes:  simpleTestResourceAttributes,
				},
				Name:        "test_resource_move_state",
				PackageName: "testresourcemovestate",
				Operations:  simpleTestResourceOperations,
				MoveState:   &codespec.MoveState{SourceResources: []string{"test_resource"}},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}

func TestConvertToProviderSpec_deprecatedResource(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_deprecated",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description:        conversion.StringPtr("POST API description"),
					DeprecationMessage: conversion.StringPtr("This resource is deprecated. Please use test_resource_new resource instead."),
					Attributes:         simpleTestResourceAttributes,
				},
				Name:        "test_resource_deprecated",
				PackageName: "testresourcedeprecated",
				Operations:  simpleTestResourceOperations,
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}

func TestConvertToProviderSpec_multipleConsecutiveCaps(t *testing.T) {
	// This test verifies that aliasing works correctly with attributes that have multiple
	// consecutive capital letters (e.g., MongoDBMajorVersion). The fix ensures that apiPath
	// is built from APIName values which preserve the original casing, avoiding the lossy
	// snake to camel case conversion that would incorrectly convert "MongoDBMajorVersion" to "MongoDbMajorVersion".
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_with_multiple_caps",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr("POST API description"),
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "mongo_db_version",
							TFModelName:              "MongoDbVersion",      // Aliased from MongoDBMajorVersion
							APIName:                  "MongoDBMajorVersion", // Original API name preserved
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("MongoDB major version with multiple consecutive capital letters"),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "nested_object",
							TFModelName:              "NestedObject",
							APIName:                  "nestedObject",
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomObjectType("NestedObject"),
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_attribute",
											TFModelName:              "InnerAttribute", // Aliased from innerAttr
											APIName:                  "innerAttr",      // Original API name preserved
											ComputedOptionalRequired: codespec.Required,
											String:                   &codespec.StringAttribute{},
											Description:              conversion.StringPtr("Inner attribute"),
											ReqBodyUsage:             codespec.AllRequestBodies,
											PresentInAnyResponse:     true,
										},
									},
								},
							},
							Description:          conversion.StringPtr(""),
							ReqBodyUsage:         codespec.AllRequestBodies,
							PresentInAnyResponse: true,
						},
						{
							TFSchemaName:             "nested_object_db",
							TFModelName:              "NestedObjectDB",
							APIName:                  "nestedObjectDB", // Preserves multiple consecutive caps in nested object name
							ComputedOptionalRequired: codespec.Optional,
							CustomType:               codespec.NewCustomObjectType("NestedObjectDB"),
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "api_key_alias",
											TFModelName:              "ApiKeyAlias", // Aliased from apiKey
											APIName:                  "apiKey",      // Original API name preserved
											ComputedOptionalRequired: codespec.Required,
											String:                   &codespec.StringAttribute{},
											Description:              conversion.StringPtr("API key for the nested object DB"),
											ReqBodyUsage:             codespec.AllRequestBodies,
											PresentInAnyResponse:     true,
										},
									},
								},
							},
							Description:          conversion.StringPtr(""),
							ReqBodyUsage:         codespec.AllRequestBodies,
							PresentInAnyResponse: true,
						},
					},
				},
				Name:        "test_resource_with_multiple_caps",
				PackageName: "testresourcewithmultiplecaps",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceWithMultipleCaps",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceWithMultipleCaps",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceWithMultipleCaps",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResourceWithMultipleCaps",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")

	// Verify that the alias lookup worked correctly by checking that:
	// 1. MongoDBMajorVersion was aliased to mongoDbVersion (schema name)
	// 2. The APIName is preserved as "MongoDBMajorVersion" (not "MongoDbMajorVersion")
	// 3. The nested alias also worked: nestedObject.innerAttr -> nestedObject.innerAttribute
	// 4. The nestedObjectDB (with multiple consecutive caps) alias works: nestedObjectDB.apiKey -> nestedObjectDB.privateApiKey
	mongoDbVersionAttr := result.Resources[0].Schema.Attributes[1] // Index 1 (after groupId)
	assert.Equal(t, "mongo_db_version", mongoDbVersionAttr.TFSchemaName, "Schema name should be aliased")
	assert.Equal(t, "MongoDbVersion", mongoDbVersionAttr.TFModelName, "Model name should be aliased")
	assert.Equal(t, "MongoDBMajorVersion", mongoDbVersionAttr.APIName, "APIName should preserve original casing with multiple consecutive caps")

	nestedObjectAttr := result.Resources[0].Schema.Attributes[2] // Index 2 (after groupId and mongoDbVersion)
	assert.Equal(t, "nested_object", nestedObjectAttr.TFSchemaName)
	innerAttr := nestedObjectAttr.SingleNested.NestedObject.Attributes[0]
	assert.Equal(t, "inner_attribute", innerAttr.TFSchemaName, "Nested schema name should be aliased")
	assert.Equal(t, "InnerAttribute", innerAttr.TFModelName, "Nested model name should be aliased")
	assert.Equal(t, "innerAttr", innerAttr.APIName, "Nested APIName should be preserved")

	nestedObjectDBAttr := result.Resources[0].Schema.Attributes[3] // Index 3 (after groupId, mongoDbVersion, and nestedObject)
	assert.Equal(t, "nested_object_db", nestedObjectDBAttr.TFSchemaName, "Schema name should be snake_case")
	assert.Equal(t, "NestedObjectDB", nestedObjectDBAttr.TFModelName, "Model name should preserve multiple consecutive caps")
	assert.Equal(t, "nestedObjectDB", nestedObjectDBAttr.APIName, "APIName should preserve original casing with multiple consecutive caps")
	apiKeyAttr := nestedObjectDBAttr.SingleNested.NestedObject.Attributes[0]
	assert.Equal(t, "api_key_alias", apiKeyAttr.TFSchemaName, "Nested schema name should be aliased")
	assert.Equal(t, "ApiKeyAlias", apiKeyAttr.TFModelName, "Nested model name should be aliased")
	assert.Equal(t, "apiKey", apiKeyAttr.APIName, "Nested APIName should be preserved")
}

// TestConvertToProviderSpec_pathParamWithAlias verifies that aliasing works correctly when an attribute
// is both a path parameter and appears in request/response bodies. The api_name tag should preserve
// the original API name while the schema uses the aliased name, and path parameters should use the aliased name.
func TestConvertToProviderSpec_pathParamWithAlias(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_path_param_with_alias",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr(testResourceDesc),
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "email",
							TFModelName:              "Email",
							APIName:                  "email",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "db_user",
							TFModelName:              "DbUser",   // Aliased from username
							APIName:                  "username", // Original API name preserved for apiname tag
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies, // Merged from path param and request body
							CreateOnly:               false,                     // Not create-only because it's in request bodies
							PresentInAnyResponse:     true,
						},
					},
				},
				Name:        "test_resource_path_param_with_alias",
				PackageName: "testresourcepathparamwithalias",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/users",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/users/{dbUser}",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/users/{dbUser}",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/users/{dbUser}",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")

	// Verify that the alias worked correctly:
	// 1. The merged attribute (from both path param and request body) has aliased schema name but preserved API name
	var dbUserAttr *codespec.Attribute
	for i := range result.Resources[0].Schema.Attributes {
		if result.Resources[0].Schema.Attributes[i].TFSchemaName == "db_user" {
			dbUserAttr = &result.Resources[0].Schema.Attributes[i]
			break
		}
	}
	require.NotNil(t, dbUserAttr, "db_user attribute should exist")
	assert.Equal(t, "db_user", dbUserAttr.TFSchemaName, "Schema name should be aliased")
	assert.Equal(t, "DbUser", dbUserAttr.TFModelName, "Model name should be aliased")
	assert.Equal(t, "username", dbUserAttr.APIName, "APIName should preserve original API name for apiname tag")
	assert.Equal(t, codespec.AllRequestBodies, dbUserAttr.ReqBodyUsage, "ReqBodyUsage should be AllRequestBodies when merged from path param and request body")
	assert.False(t, dbUserAttr.CreateOnly, "CreateOnly should be false when attribute is in request bodies")

	// 2. Path parameters in operations use the aliased name
	assert.Contains(t, result.Resources[0].Operations.Read.Path, "{dbUser}", "Read path should use aliased path param")
	assert.Contains(t, result.Resources[0].Operations.Update.Path, "{dbUser}", "Update path should use aliased path param")
	assert.Contains(t, result.Resources[0].Operations.Delete.Path, "{dbUser}", "Delete path should use aliased path param")
}

func TestConvertToProviderSpec_ignoreSchemaAndIdAttributes(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_schema_ignore_and_id_attributes",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr("POST API description"),
					// due to schema_ignore string_attr is not included in the attributes
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
					},
				},
				Name:         "test_schema_ignore_and_id_attributes",
				PackageName:  "testschemaignoreandidattributes",
				Operations:   simpleTestResourceOperations,
				IDAttributes: []string{"group_id", "id"},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")
}

// TestConvertToProviderSpec_withDataSources verifies that data sources are correctly generated
// when a datasources block is defined in the config. This test verifies:
// 1. Singular data source schema is generated with all response attributes as Computed
// 2. Plural data source schema is generated (reusing getPluralDSAttributes helper)
// 3. Path parameters are Required in the data source schema
// 4. Aliasing in data source config (groupId -> projectId) works correctly and handles duplicates
// 5. Data source operations are correctly extracted from config
func TestConvertToProviderSpec_withDataSources(t *testing.T) {
	pluralDSAttributes := getPluralDSAttributes()

	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      "testdata/config-datasources.yml",
		inputResourceName:    "test_resource_with_datasource",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Schema: &codespec.Schema{
					Description: conversion.StringPtr(testResourceDesc),
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "bool_default_attr",
							TFModelName:              "BoolDefaultAttr",
							APIName:                  "boolDefaultAttr",
							ComputedOptionalRequired: codespec.ComputedOptional,
							Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "count",
							TFModelName:              "Count",
							APIName:                  "count",
							ComputedOptionalRequired: codespec.Optional,
							Int64:                    &codespec.Int64Attribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "create_date",
							TFModelName:              "CreateDate",
							APIName:                  "createDate",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Computed,
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testPathParamDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							CreateOnly:               true,
						},
						{
							TFSchemaName:             "num_double_default_attr",
							TFModelName:              "NumDoubleDefaultAttr",
							APIName:                  "numDoubleDefaultAttr",
							Float64:                  &codespec.Float64Attribute{Default: conversion.Pointer(2.0)},
							ComputedOptionalRequired: codespec.ComputedOptional,
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "str_computed_attr",
							TFModelName:              "StrComputedAttr",
							APIName:                  "strComputedAttr",
							ComputedOptionalRequired: codespec.Computed,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.OmitAlways,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "str_req_attr1",
							TFModelName:              "StrReqAttr1",
							APIName:                  "strReqAttr1",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "str_req_attr2",
							TFModelName:              "StrReqAttr2",
							APIName:                  "strReqAttr2",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
						{
							TFSchemaName:             "str_req_attr3",
							TFModelName:              "StrReqAttr3",
							APIName:                  "strReqAttr3",
							String:                   &codespec.StringAttribute{},
							ComputedOptionalRequired: codespec.Required,
							Description:              conversion.StringPtr(testFieldDesc),
							ReqBodyUsage:             codespec.AllRequestBodies,
							PresentInAnyResponse:     true,
						},
					},
				},
				Name:        "test_resource_with_datasource",
				PackageName: "testresourcewithdatasource",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
				// Data sources model with independent schema
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						SingularDSDescription: conversion.StringPtr("GET API description"),
						PluralDSDescription:   conversion.StringPtr("LIST API description"),
						PluralDSAttributes:    &pluralDSAttributes,
						SingularDSAttributes: &codespec.Attributes{
							// All response attributes are Computed
							{
								TFSchemaName:             "bool_default_attr",
								TFModelName:              "BoolDefaultAttr",
								APIName:                  "boolDefaultAttr",
								ComputedOptionalRequired: codespec.Computed,
								Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
								ReqBodyUsage:             codespec.OmitAlways,
							},
							{
								TFSchemaName:             "count",
								TFModelName:              "Count",
								APIName:                  "count",
								ComputedOptionalRequired: codespec.Computed,
								Int64:                    &codespec.Int64Attribute{},
								Description:              conversion.StringPtr(testFieldDesc),
								ReqBodyUsage:             codespec.OmitAlways,
							},
							{
								TFSchemaName:             "create_date",
								TFModelName:              "CreateDate",
								APIName:                  "createDate",
								String:                   &codespec.StringAttribute{},
								ComputedOptionalRequired: codespec.Computed,
								Description:              conversion.StringPtr(testFieldDesc),
								ReqBodyUsage:             codespec.OmitAlways,
							},
							{
								TFSchemaName:             "num_double_default_attr",
								TFModelName:              "NumDoubleDefaultAttr",
								APIName:                  "numDoubleDefaultAttr",
								Float64:                  &codespec.Float64Attribute{Default: conversion.Pointer(2.0)},
								ComputedOptionalRequired: codespec.Computed,
								ReqBodyUsage:             codespec.OmitAlways,
							},
							// Path param groupId aliased to projectId - Required (not computed)
							{
								TFSchemaName:             "project_id",
								TFModelName:              "ProjectId",
								APIName:                  "groupId",
								ComputedOptionalRequired: codespec.Required,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testPathParamDesc),
								ReqBodyUsage:             codespec.OmitAlways,
							},
							{
								TFSchemaName:             "str_computed_attr",
								TFModelName:              "StrComputedAttr",
								APIName:                  "strComputedAttr",
								ComputedOptionalRequired: codespec.Computed,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testFieldDesc),
								ReqBodyUsage:             codespec.OmitAlways,
							},
							{
								TFSchemaName:             "str_req_attr1",
								TFModelName:              "StrReqAttr1",
								APIName:                  "strReqAttr1",
								ComputedOptionalRequired: codespec.Computed,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testFieldDesc),
								ReqBodyUsage:             codespec.OmitAlways,
							},
							{
								TFSchemaName:             "str_req_attr2",
								TFModelName:              "StrReqAttr2",
								APIName:                  "strReqAttr2",
								ComputedOptionalRequired: codespec.Computed,
								String:                   &codespec.StringAttribute{},
								Description:              conversion.StringPtr(testFieldDesc),
								ReqBodyUsage:             codespec.OmitAlways,
							},
							{
								TFSchemaName:             "str_req_attr3",
								TFModelName:              "StrReqAttr3",
								APIName:                  "strReqAttr3",
								String:                   &codespec.StringAttribute{},
								ComputedOptionalRequired: codespec.Computed,
								Description:              conversion.StringPtr(testFieldDesc),
								ReqBodyUsage:             codespec.OmitAlways,
							},
						},
					},
					Operations: codespec.APIOperations{
						Read: &codespec.APIOperation{
							Path:       "/api/atlas/v2/groups/{projectId}/testResource",
							HTTPMethod: "GET",
						},
						List: &codespec.APIOperation{
							Path:       "/api/atlas/v2/groups/{projectId}/testResources",
							HTTPMethod: "GET",
						},
						VersionHeader: "application/vnd.atlas.2023-01-01+json",
					},
				},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)
	assert.Equal(t, tc.expectedResult, result, "Expected result to match the specified structure")

	// Additional assertions to verify key singular data source behaviors
	ds := result.Resources[0].DataSources
	require.NotNil(t, ds, "DataSources should be populated")
	require.NotNil(t, ds.Schema, "DataSources.Schema should be populated")

	// Verify path param is Required (not Computed) even in singular data source
	var projectIDAttr *codespec.Attribute
	for i := range *ds.Schema.SingularDSAttributes {
		if (*ds.Schema.SingularDSAttributes)[i].TFSchemaName == "project_id" {
			projectIDAttr = &(*ds.Schema.SingularDSAttributes)[i]
			break
		}
	}
	require.NotNil(t, projectIDAttr, "project_id attribute should exist in data source")
	assert.Equal(t, codespec.Required, projectIDAttr.ComputedOptionalRequired, "Path param should be Required in data source")
	assert.Equal(t, "groupId", projectIDAttr.APIName, "APIName should preserve original name for aliased path param")

	// Verify response attributes are Computed
	for _, attr := range *ds.Schema.SingularDSAttributes {
		if attr.TFSchemaName != "project_id" { // Skip path param
			assert.Equal(t, codespec.Computed, attr.ComputedOptionalRequired,
				"Response attribute %s should be Computed in data source", attr.TFSchemaName)
		}
	}

	// Verify operation path uses aliased placeholder
	assert.Contains(t, ds.Operations.Read.Path, "{projectId}",
		"Data source Read path should use aliased path param placeholder")
}

// TestConvertToProviderSpec_withPluralDataSource verifies that plural data sources are correctly generated
// when a datasources block with a list operation is defined in the config. This test verifies:
// 1. Plural data source schema is generated with results list containing resource attributes
// 2. Path parameters are Required in the plural data source schema
// 3. All nested attributes within results are Computed
// 4. Aliasing in plural data source config (groupId -> projectId) works correctly
// 5. List operation is correctly extracted and path uses aliased placeholders
func TestConvertToProviderSpec_withPluralDataSource(t *testing.T) {
	pluralDSAttributes := getPluralDSAttributes()

	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      "testdata/config-datasources.yml",
		inputResourceName:    "test_resource_with_datasource",

		expectedResult: &codespec.Model{
			Resources: []codespec.Resource{{
				Name:        "test_resource_with_datasource",
				PackageName: "testresourcewithdatasource",
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "POST",
					},
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "GET",
					},
					Update: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "PATCH",
					},
					Delete: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/testResource",
						HTTPMethod: "DELETE",
					},
					VersionHeader: "application/vnd.atlas.2023-01-01+json",
				},
				DataSources: &codespec.DataSources{
					Schema: &codespec.DataSourceSchema{
						PluralDSDescription: conversion.StringPtr("LIST API description"),
						PluralDSAttributes:  &pluralDSAttributes,
					},
					Operations: codespec.APIOperations{
						List: &codespec.APIOperation{
							Path:       "/api/atlas/v2/groups/{projectId}/testResources",
							HTTPMethod: "GET",
						},
						VersionHeader: "application/vnd.atlas.2023-01-01+json",
					},
				},
			}},
		},
	}

	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)

	// We only verify the plural data source parts
	ds := result.Resources[0].DataSources
	require.NotNil(t, ds, "DataSources should be populated")
	require.NotNil(t, ds.Schema, "DataSources.Schema should be populated")
	require.NotNil(t, ds.Schema.PluralDSAttributes, "PluralDSAttributes should be populated")

	// Verify path param is Required (not Computed) in plural data source
	var projectIDAttr *codespec.Attribute
	for i := range *ds.Schema.PluralDSAttributes {
		if (*ds.Schema.PluralDSAttributes)[i].TFSchemaName == "project_id" {
			projectIDAttr = &(*ds.Schema.PluralDSAttributes)[i]
			break
		}
	}
	require.NotNil(t, projectIDAttr, "project_id attribute should exist in plural data source")
	assert.Equal(t, codespec.Required, projectIDAttr.ComputedOptionalRequired, "Path param should be Required in plural data source")
	assert.Equal(t, "groupId", projectIDAttr.APIName, "APIName should preserve original name for aliased path param")

	// Verify results array exists and is Computed
	var resultsAttr *codespec.Attribute
	for i := range *ds.Schema.PluralDSAttributes {
		if (*ds.Schema.PluralDSAttributes)[i].TFSchemaName == "results" {
			resultsAttr = &(*ds.Schema.PluralDSAttributes)[i]
			break
		}
	}
	require.NotNil(t, resultsAttr, "results attribute should exist in plural data source")
	assert.Equal(t, codespec.Computed, resultsAttr.ComputedOptionalRequired, "results attribute should be Computed")
	require.NotNil(t, resultsAttr.ListNested, "results should be a ListNested attribute")

	// Verify nested attributes within results are all Computed
	for i := range resultsAttr.ListNested.NestedObject.Attributes {
		attr := &resultsAttr.ListNested.NestedObject.Attributes[i]
		assert.Equal(t, codespec.Computed, attr.ComputedOptionalRequired,
			"Nested attribute %s within results should be Computed in plural data source", attr.TFSchemaName)
	}

	// Verify the expected structure matches (comparing against the helper function output)
	assert.Equal(t, tc.expectedResult.Resources[0].DataSources.Schema.PluralDSAttributes, ds.Schema.PluralDSAttributes,
		"PluralDSAttributes should match expected structure")

	// Verify List operation path uses aliased placeholder
	require.NotNil(t, ds.Operations.List, "List operation should be defined for plural data source")
	assert.Contains(t, ds.Operations.List.Path, "{projectId}",
		"Plural data source List path should use aliased path param placeholder")
}

func TestConvertToProviderSpec_withTagsAndLabelsAsMapType(t *testing.T) {
	tc := convertToSpecTestCase{
		inputOpenAPISpecPath: testDataAPISpecPath,
		inputConfigPath:      testDataConfigPath,
		inputResourceName:    "test_resource_with_tags_and_labels_as_map_type",
	}
	result, err := codespec.ToCodeSpecModel(tc.inputOpenAPISpecPath, tc.inputConfigPath, &tc.inputResourceName)
	require.NoError(t, err)

	resourceAttrs := []codespec.Attribute(result.Resources[0].Schema.Attributes)
	dataSourceAttrs := []codespec.Attribute(*result.Resources[0].DataSources.Schema.SingularDSAttributes)

	assert.Len(t, resourceAttrs, 2, "Expected resource to have 2 attributes")
	assert.Len(t, dataSourceAttrs, 2, "Expected data source to have 2 attributes")

	// Verify resource and data source attributes: two map attributes named "labels" and "tags"
	for _, attr := range append(resourceAttrs, dataSourceAttrs...) {
		assert.NotNil(t, attr.Map, "Attribute %s should be a map type", attr.TFSchemaName)
		assert.Nil(t, attr.ListNested, "Attribute %s should be a map type", attr.TFSchemaName)
		assert.True(t, attr.ListTypeAsMap, "Attribute %s should have correct flag", attr.TFSchemaName)
		assert.Contains(t, []string{"labels", "tags"}, attr.TFSchemaName, "Resource attribute should be either 'labels' or 'tags'")
	}
}
