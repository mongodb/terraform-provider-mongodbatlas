package codespec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
)

func TestApplyTimeoutTransformation(t *testing.T) {
	tests := map[string]struct {
		inputOperations  codespec.APIOperations
		expectedTimeouts []codespec.Operation
	}{
		"No wait blocks - no timeout attribute added": {
			inputOperations: codespec.APIOperations{
				Create: &codespec.APIOperation{},
				Read:   &codespec.APIOperation{},
				Update: &codespec.APIOperation{},
			},
			expectedTimeouts: nil,
		},
		"Create wait only": {
			inputOperations: codespec.APIOperations{
				Create: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Read:   &codespec.APIOperation{},
				Update: &codespec.APIOperation{},
			},
			expectedTimeouts: []codespec.Operation{codespec.Create},
		},
		"Create, Update, Delete waits": {
			inputOperations: codespec.APIOperations{
				Create: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Read: &codespec.APIOperation{},
				Update: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Delete: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
			},
			expectedTimeouts: []codespec.Operation{codespec.Create, codespec.Update, codespec.Delete},
		},
		"All operations with waits": {
			inputOperations: codespec.APIOperations{
				Create: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Read: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Update: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Delete: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
			},
			expectedTimeouts: []codespec.Operation{codespec.Create, codespec.Update, codespec.Read, codespec.Delete},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resource := &codespec.Resource{
				Name: "test_resource",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{},
				},
				Operations: tc.inputOperations,
			}
			codespec.ApplyTimeoutTransformation(resource)
			if tc.expectedTimeouts == nil {
				assert.Empty(t, resource.Schema.Attributes)
			} else {
				assert.Len(t, resource.Schema.Attributes, 1)
				expectedAttr := codespec.Attribute{
					TFSchemaName: "timeouts",
					TFModelName:  "Timeouts",
					ReqBodyUsage: codespec.OmitAlways,
					Timeouts:     &codespec.TimeoutsAttribute{ConfigurableTimeouts: tc.expectedTimeouts},
				}
				assert.Equal(t, expectedAttr, resource.Schema.Attributes[0])
			}
		})
	}
}

func TestApplyDeleteOnCreateTimeoutTransformation(t *testing.T) {
	tests := map[string]struct {
		inputOperations                codespec.APIOperations
		shouldAddDeleteOnCreateTimeout bool
	}{
		"Create with wait and Delete operation - attribute added": {
			inputOperations: codespec.APIOperations{
				Create: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Read:   &codespec.APIOperation{},
				Update: &codespec.APIOperation{},
				Delete: &codespec.APIOperation{},
			},
			shouldAddDeleteOnCreateTimeout: true,
		},
		"Create with wait but no Delete operation - attribute not added": {
			inputOperations: codespec.APIOperations{
				Create: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Read:   &codespec.APIOperation{},
				Update: &codespec.APIOperation{},
			},
			shouldAddDeleteOnCreateTimeout: false,
		},
		"Create without wait but with Delete operation - attribute not added": {
			inputOperations: codespec.APIOperations{
				Create: &codespec.APIOperation{},
				Read:   &codespec.APIOperation{},
				Update: &codespec.APIOperation{},
				Delete: &codespec.APIOperation{},
			},
			shouldAddDeleteOnCreateTimeout: false,
		},
		"No Create wait and no Delete operation - attribute not added": {
			inputOperations: codespec.APIOperations{
				Create: &codespec.APIOperation{},
				Read:   &codespec.APIOperation{},
				Update: &codespec.APIOperation{},
			},
			shouldAddDeleteOnCreateTimeout: false,
		},
		"Create with wait, Update with wait, and Delete operation - attribute added": {
			inputOperations: codespec.APIOperations{
				Create: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Read: &codespec.APIOperation{},
				Update: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Delete: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
			},
			shouldAddDeleteOnCreateTimeout: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resource := &codespec.Resource{
				Name: "test_resource",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{},
				},
				Operations: tc.inputOperations,
			}
			codespec.ApplyDeleteOnCreateTimeoutTransformation(resource)
			if !tc.shouldAddDeleteOnCreateTimeout {
				assert.Empty(t, resource.Schema.Attributes)
			} else {
				assert.Len(t, resource.Schema.Attributes, 1)
				description := codespec.DeleteOnCreateTimeoutDescription
				expectedAttr := codespec.Attribute{
					TFSchemaName:             "delete_on_create_timeout",
					TFModelName:              "DeleteOnCreateTimeout",
					Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(true)},
					Description:              &description,
					ReqBodyUsage:             codespec.OmitAlways,
					CreateOnly:               true,
					ComputedOptionalRequired: codespec.ComputedOptional,
				}
				assert.Equal(t, expectedAttr, resource.Schema.Attributes[0])
			}
		})
	}
}

func TestApplyTransformationsToDataSources_AliasTransformation(t *testing.T) {
	tests := map[string]struct {
		inputDataSources   *codespec.DataSources
		inputConfig        *config.DataSources
		expectedReadPath   string
		expectedListPath   string
		expectedAttributes codespec.Attributes
	}{
		"Alias applied to attribute and path param": {
			inputDataSources: &codespec.DataSources{
				Schema: &codespec.DataSourceSchema{
					SingularDSAttributes: &codespec.Attributes{
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							TFSchemaName:             "name",
							TFModelName:              "Name",
							APIName:                  "name",
							ComputedOptionalRequired: codespec.Computed,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitAlways,
						},
					},
				},
				Operations: codespec.APIOperations{
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/resource",
						HTTPMethod: "GET",
					},
				},
			},
			inputConfig: &config.DataSources{
				SchemaOptions: config.SchemaOptions{
					Aliases: map[string]string{
						"groupId": "projectId",
					},
				},
			},
			// Note: attributes are NOT sorted by transformation - they keep original order with alias applied
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "project_id",
					TFModelName:              "ProjectId",
					APIName:                  "groupId", // APIName preserved
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
				{
					TFSchemaName:             "name",
					TFModelName:              "Name",
					APIName:                  "name",
					ComputedOptionalRequired: codespec.Computed,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
			},
			expectedReadPath: "/api/atlas/v2/groups/{projectId}/resource",
		},
		"Alias applied to List operation path": {
			inputDataSources: &codespec.DataSources{
				Schema: &codespec.DataSourceSchema{
					SingularDSAttributes: &codespec.Attributes{
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitAlways,
						},
					},
				},
				Operations: codespec.APIOperations{
					List: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/resources",
						HTTPMethod: "GET",
					},
				},
			},
			inputConfig: &config.DataSources{
				SchemaOptions: config.SchemaOptions{
					Aliases: map[string]string{
						"groupId": "projectId",
					},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "project_id",
					TFModelName:              "ProjectId",
					APIName:                  "groupId",
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
			},
			expectedListPath: "/api/atlas/v2/groups/{projectId}/resources",
		},
		"No aliases - attributes unchanged": {
			inputDataSources: &codespec.DataSources{
				Schema: &codespec.DataSourceSchema{
					SingularDSAttributes: &codespec.Attributes{
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitAlways,
						},
					},
				},
				Operations: codespec.APIOperations{
					Read: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/resource",
						HTTPMethod: "GET",
					},
				},
			},
			inputConfig: &config.DataSources{
				SchemaOptions: config.SchemaOptions{},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "group_id",
					TFModelName:              "GroupId",
					APIName:                  "groupId",
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
			},
			expectedReadPath: "/api/atlas/v2/groups/{groupId}/resource",
		},
		"Alias applied to plural data source attributes and List path": {
			inputDataSources: &codespec.DataSources{
				Schema: &codespec.DataSourceSchema{
					PluralDSAttributes: &codespec.Attributes{
						{
							TFSchemaName:             "group_id",
							TFModelName:              "GroupId",
							APIName:                  "groupId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							TFSchemaName:             "results",
							TFModelName:              "Results",
							APIName:                  "results",
							ComputedOptionalRequired: codespec.Computed,
							CustomType:               codespec.NewCustomNestedListType("Results"),
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "name",
											TFModelName:              "Name",
											APIName:                  "name",
											ComputedOptionalRequired: codespec.Computed,
											String:                   &codespec.StringAttribute{},
											ReqBodyUsage:             codespec.OmitAlways,
										},
									},
								},
							},
							ReqBodyUsage: codespec.OmitAlways,
						},
					},
				},
				Operations: codespec.APIOperations{
					List: &codespec.APIOperation{
						Path:       "/api/atlas/v2/groups/{groupId}/resources",
						HTTPMethod: "GET",
					},
				},
			},
			inputConfig: &config.DataSources{
				SchemaOptions: config.SchemaOptions{
					Aliases: map[string]string{
						"groupId": "projectId",
					},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "project_id",
					TFModelName:              "ProjectId",
					APIName:                  "groupId", // APIName preserved
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
				{
					TFSchemaName:             "results",
					TFModelName:              "Results",
					APIName:                  "results",
					ComputedOptionalRequired: codespec.Computed,
					CustomType:               codespec.NewCustomNestedListType("Results"),
					ListNested: &codespec.ListNestedAttribute{
						NestedObject: codespec.NestedAttributeObject{
							Attributes: codespec.Attributes{
								{
									TFSchemaName:             "name",
									TFModelName:              "Name",
									APIName:                  "name",
									ComputedOptionalRequired: codespec.Computed,
									String:                   &codespec.StringAttribute{},
									ReqBodyUsage:             codespec.OmitAlways,
								},
							},
						},
					},
					ReqBodyUsage: codespec.OmitAlways,
				},
			},
			expectedListPath: "/api/atlas/v2/groups/{projectId}/resources",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := codespec.ApplyTransformationsToDataSources(tc.inputConfig, tc.inputDataSources)
			require.NoError(t, err)

			// Check if test is for singular or plural data sources
			if tc.inputDataSources.Schema.SingularDSAttributes != nil {
				assert.Equal(t, tc.expectedAttributes, *tc.inputDataSources.Schema.SingularDSAttributes)
			}
			if tc.inputDataSources.Schema.PluralDSAttributes != nil {
				assert.Equal(t, tc.expectedAttributes, *tc.inputDataSources.Schema.PluralDSAttributes)
			}

			if tc.expectedReadPath != "" {
				assert.Equal(t, tc.expectedReadPath, tc.inputDataSources.Operations.Read.Path)
			}
			if tc.expectedListPath != "" {
				assert.Equal(t, tc.expectedListPath, tc.inputDataSources.Operations.List.Path)
			}
		})
	}
}

func TestApplyTransformationsToDataSources_OverrideTransformation(t *testing.T) {
	tests := map[string]struct {
		inputDataSources   *codespec.DataSources
		inputConfig        *config.DataSources
		expectedAttributes codespec.Attributes
	}{
		"Override description": {
			inputDataSources: &codespec.DataSources{
				Schema: &codespec.DataSourceSchema{
					SingularDSAttributes: &codespec.Attributes{
						{
							TFSchemaName:             "name",
							TFModelName:              "Name",
							APIName:                  "name",
							ComputedOptionalRequired: codespec.Computed,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Original description"),
							ReqBodyUsage:             codespec.OmitAlways,
						},
					},
				},
				Operations: codespec.APIOperations{},
			},
			inputConfig: &config.DataSources{
				SchemaOptions: config.SchemaOptions{
					Overrides: map[string]config.Override{
						"name": {
							Description: "Overridden description",
						},
					},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "name",
					TFModelName:              "Name",
					APIName:                  "name",
					ComputedOptionalRequired: codespec.Computed,
					String:                   &codespec.StringAttribute{},
					Description:              conversion.StringPtr("Overridden description"),
					ReqBodyUsage:             codespec.OmitAlways,
				},
			},
		},
		"Override computability": {
			inputDataSources: &codespec.DataSources{
				Schema: &codespec.DataSourceSchema{
					SingularDSAttributes: &codespec.Attributes{
						{
							TFSchemaName:             "optional_attr",
							TFModelName:              "OptionalAttr",
							APIName:                  "optionalAttr",
							ComputedOptionalRequired: codespec.Computed,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitAlways,
						},
					},
				},
				Operations: codespec.APIOperations{},
			},
			inputConfig: &config.DataSources{
				SchemaOptions: config.SchemaOptions{
					Overrides: map[string]config.Override{
						"optional_attr": {
							Computability: &config.Computability{Optional: true},
						},
					},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "optional_attr",
					TFModelName:              "OptionalAttr",
					APIName:                  "optionalAttr",
					ComputedOptionalRequired: codespec.Optional,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
			},
		},
		"Override description in plural data source": {
			inputDataSources: &codespec.DataSources{
				Schema: &codespec.DataSourceSchema{
					PluralDSAttributes: &codespec.Attributes{
						{
							TFSchemaName:             "project_id",
							TFModelName:              "ProjectId",
							APIName:                  "projectId",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							Description:              conversion.StringPtr("Original description"),
							ReqBodyUsage:             codespec.OmitAlways,
						},
					},
				},
				Operations: codespec.APIOperations{},
			},
			inputConfig: &config.DataSources{
				SchemaOptions: config.SchemaOptions{
					Overrides: map[string]config.Override{
						"project_id": {
							Description: "Overridden description for plural",
						},
					},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "project_id",
					TFModelName:              "ProjectId",
					APIName:                  "projectId",
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					Description:              conversion.StringPtr("Overridden description for plural"),
					ReqBodyUsage:             codespec.OmitAlways,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := codespec.ApplyTransformationsToDataSources(tc.inputConfig, tc.inputDataSources)
			require.NoError(t, err)

			// Check if test is for singular or plural data sources
			if tc.inputDataSources.Schema.SingularDSAttributes != nil {
				assert.Equal(t, tc.expectedAttributes, *tc.inputDataSources.Schema.SingularDSAttributes)
			}
			if tc.inputDataSources.Schema.PluralDSAttributes != nil {
				assert.Equal(t, tc.expectedAttributes, *tc.inputDataSources.Schema.PluralDSAttributes)
			}
		})
	}
}

func TestApplyTransformationsToDataSources_IgnoreTransformation(t *testing.T) {
	inputDataSources := &codespec.DataSources{
		Schema: &codespec.DataSourceSchema{
			SingularDSAttributes: &codespec.Attributes{
				{
					TFSchemaName:             "keep_attr",
					TFModelName:              "KeepAttr",
					APIName:                  "keepAttr",
					ComputedOptionalRequired: codespec.Computed,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
				{
					TFSchemaName:             "ignore_attr",
					TFModelName:              "IgnoreAttr",
					APIName:                  "ignoreAttr",
					ComputedOptionalRequired: codespec.Computed,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
			},
		},
		Operations: codespec.APIOperations{},
	}

	inputConfig := &config.DataSources{
		SchemaOptions: config.SchemaOptions{
			Ignores: []string{"ignore_attr"},
		},
	}

	err := codespec.ApplyTransformationsToDataSources(inputConfig, inputDataSources)
	require.NoError(t, err)

	// Only keep_attr should remain
	expectedAttributes := codespec.Attributes{
		{
			TFSchemaName:             "keep_attr",
			TFModelName:              "KeepAttr",
			APIName:                  "keepAttr",
			ComputedOptionalRequired: codespec.Computed,
			String:                   &codespec.StringAttribute{},
			ReqBodyUsage:             codespec.OmitAlways,
		},
	}

	assert.Equal(t, expectedAttributes, *inputDataSources.Schema.SingularDSAttributes)
}

func TestApplyTransformationsToDataSources_IgnorePluralDataSource(t *testing.T) {
	inputDataSources := &codespec.DataSources{
		Schema: &codespec.DataSourceSchema{
			PluralDSAttributes: &codespec.Attributes{
				{
					TFSchemaName:             "project_id",
					TFModelName:              "ProjectId",
					APIName:                  "projectId",
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
				{
					TFSchemaName:             "results",
					TFModelName:              "Results",
					APIName:                  "results",
					ComputedOptionalRequired: codespec.Computed,
					CustomType:               codespec.NewCustomNestedListType("Results"),
					ListNested: &codespec.ListNestedAttribute{
						NestedObject: codespec.NestedAttributeObject{
							Attributes: codespec.Attributes{
								{
									TFSchemaName:             "keep_attr",
									TFModelName:              "KeepAttr",
									APIName:                  "keepAttr",
									ComputedOptionalRequired: codespec.Computed,
									String:                   &codespec.StringAttribute{},
									ReqBodyUsage:             codespec.OmitAlways,
								},
								{
									TFSchemaName:             "ignore_attr",
									TFModelName:              "IgnoreAttr",
									APIName:                  "ignoreAttr",
									ComputedOptionalRequired: codespec.Computed,
									String:                   &codespec.StringAttribute{},
									ReqBodyUsage:             codespec.OmitAlways,
								},
							},
						},
					},
					ReqBodyUsage: codespec.OmitAlways,
				},
			},
		},
		Operations: codespec.APIOperations{},
	}

	inputConfig := &config.DataSources{
		SchemaOptions: config.SchemaOptions{
			Ignores: []string{"ignore_attr"},
		},
	}

	err := codespec.ApplyTransformationsToDataSources(inputConfig, inputDataSources)
	require.NoError(t, err)

	// Verify project_id is still present
	require.Len(t, *inputDataSources.Schema.PluralDSAttributes, 2, "Should have project_id and results")

	// Verify results array with only keep_attr
	var resultsAttr *codespec.Attribute
	for i := range *inputDataSources.Schema.PluralDSAttributes {
		if (*inputDataSources.Schema.PluralDSAttributes)[i].TFSchemaName == "results" {
			resultsAttr = &(*inputDataSources.Schema.PluralDSAttributes)[i]
			break
		}
	}
	require.NotNil(t, resultsAttr, "results attribute should exist")
	require.NotNil(t, resultsAttr.ListNested, "results should have ListNested")

	// Only keep_attr should remain in nested attributes
	expectedNestedAttrs := codespec.Attributes{
		{
			TFSchemaName:             "keep_attr",
			TFModelName:              "KeepAttr",
			APIName:                  "keepAttr",
			ComputedOptionalRequired: codespec.Computed,
			String:                   &codespec.StringAttribute{},
			ReqBodyUsage:             codespec.OmitAlways,
		},
	}

	assert.Equal(t, expectedNestedAttrs, resultsAttr.ListNested.NestedObject.Attributes)
}

func TestApplyTransformationsToDataSources_NilInputs(t *testing.T) {
	// Test nil DataSources
	err := codespec.ApplyTransformationsToDataSources(&config.DataSources{}, nil)
	require.NoError(t, err)

	// Test nil Schema
	err = codespec.ApplyTransformationsToDataSources(&config.DataSources{}, &codespec.DataSources{})
	require.NoError(t, err)
}
