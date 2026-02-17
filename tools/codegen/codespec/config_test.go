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

func TestApplyTransformationsToResource_AliasAttributeTransformation(t *testing.T) {
	tests := map[string]struct {
		inputResource      *codespec.Resource
		inputConfig        *config.Resource
		expectedAttributes codespec.Attributes
	}{
		"Root-level simple alias renames attribute preserving APIName": {
			inputConfig: &config.Resource{
				SchemaOptions: config.SchemaOptions{
					Aliases: map[string]string{
						"groupId": "projectId",
					},
				},
			},
			inputResource: &codespec.Resource{
				Schema: &codespec.Schema{
					Attributes: codespec.Attributes{
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
					Create: &codespec.APIOperation{},
					Read:   &codespec.APIOperation{},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "project_id",
					TFModelName:              "ProjectId",
					APIName:                  "groupId", // preserved for apiname tag
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
					CreateOnly:               true, // OmitAlways + Required triggers createOnly
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
		"Path-based nested alias only renames targeted nested attribute": {
			inputConfig: &config.Resource{
				SchemaOptions: config.SchemaOptions{
					Aliases: map[string]string{
						"nestedObjA.innerAttr": "renamedAttr",
					},
				},
			},
			inputResource: &codespec.Resource{
				Schema: &codespec.Schema{
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "nested_obj_a",
							TFModelName:              "NestedObjA",
							APIName:                  "nestedObjA",
							ComputedOptionalRequired: codespec.Computed,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_attr",
											TFModelName:              "InnerAttr",
											APIName:                  "innerAttr",
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
							TFSchemaName:             "nested_obj_b",
							TFModelName:              "NestedObjB",
							APIName:                  "nestedObjB",
							ComputedOptionalRequired: codespec.Computed,
							ListNested: &codespec.ListNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_attr",
											TFModelName:              "InnerAttr",
											APIName:                  "innerAttr",
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
					Create: &codespec.APIOperation{},
					Read:   &codespec.APIOperation{},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "nested_obj_a",
					TFModelName:              "NestedObjA",
					APIName:                  "nestedObjA",
					ComputedOptionalRequired: codespec.Computed,
					ListNested: &codespec.ListNestedAttribute{
						NestedObject: codespec.NestedAttributeObject{
							Attributes: codespec.Attributes{
								{
									// Renamed via path-based alias
									TFSchemaName:             "renamed_attr",
									TFModelName:              "RenamedAttr",
									APIName:                  "innerAttr", // preserved
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
					TFSchemaName:             "nested_obj_b",
					TFModelName:              "NestedObjB",
					APIName:                  "nestedObjB",
					ComputedOptionalRequired: codespec.Computed,
					ListNested: &codespec.ListNestedAttribute{
						NestedObject: codespec.NestedAttributeObject{
							Attributes: codespec.Attributes{
								{
									// NOT renamed - path-based alias is scoped to nestedObjA only
									TFSchemaName:             "inner_attr",
									TFModelName:              "InnerAttr",
									APIName:                  "innerAttr",
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
		"Non-dotted alias does NOT rename same-named nested attribute": {
			inputConfig: &config.Resource{
				SchemaOptions: config.SchemaOptions{
					Aliases: map[string]string{
						"innerAttr": "renamedAttr",
					},
				},
			},
			inputResource: &codespec.Resource{
				Schema: &codespec.Schema{
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "inner_attr",
							TFModelName:              "InnerAttr",
							APIName:                  "innerAttr",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitAlways,
						},
						{
							TFSchemaName:             "nested_obj",
							TFModelName:              "NestedObj",
							APIName:                  "nestedObj",
							ComputedOptionalRequired: codespec.Computed,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_attr",
											TFModelName:              "InnerAttr",
											APIName:                  "innerAttr",
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
					Create: &codespec.APIOperation{},
					Read:   &codespec.APIOperation{},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "renamed_attr",
					TFModelName:              "RenamedAttr",
					APIName:                  "innerAttr", // preserved
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
					CreateOnly:               true,
				},
				{
					TFSchemaName:             "nested_obj",
					TFModelName:              "NestedObj",
					APIName:                  "nestedObj",
					ComputedOptionalRequired: codespec.Computed,
					SingleNested: &codespec.SingleNestedAttribute{
						NestedObject: codespec.NestedAttributeObject{
							Attributes: codespec.Attributes{
								{
									// NOT renamed - non-dotted alias only applies at root level
									TFSchemaName:             "inner_attr",
									TFModelName:              "InnerAttr",
									APIName:                  "innerAttr",
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := codespec.ApplyTransformationsToResource(tc.inputConfig, tc.inputResource)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedAttributes, tc.inputResource.Schema.Attributes)
		})
	}
}

func TestApplyTransformationsToResource_AliasDiscriminatorTransformation(t *testing.T) {
	tests := map[string]struct {
		inputResource      *codespec.Resource
		inputConfig        *config.Resource
		expectedDiscrim    *codespec.Discriminator
		expectedAttributes codespec.Attributes
	}{
		"Root-level discriminator property and variant lists renamed by alias": {
			inputConfig: &config.Resource{
				SchemaOptions: config.SchemaOptions{
					Aliases: map[string]string{
						"typeField": "kind",
					},
				},
			},
			inputResource: &codespec.Resource{
				Schema: &codespec.Schema{
					Discriminator: &codespec.Discriminator{
						PropertyName: "type_field",
						Mapping: map[string]codespec.DiscriminatorType{
							"VariantA": {
								Allowed:  []string{"type_field", "variant_a_attr"},
								Required: []string{"type_field"},
							},
							"VariantB": {
								Allowed:  []string{"type_field", "variant_b_attr"},
								Required: []string{"type_field"},
							},
						},
					},
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "type_field",
							TFModelName:              "TypeField",
							APIName:                  "typeField",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "variant_a_attr",
							TFModelName:              "VariantAAttr",
							APIName:                  "variantAAttr",
							ComputedOptionalRequired: codespec.Computed,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.OmitAlways,
						},
					},
				},
				Operations: codespec.APIOperations{
					Create: &codespec.APIOperation{},
					Read:   &codespec.APIOperation{},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "kind",
					TFModelName:              "Kind",
					APIName:                  "typeField", // preserved
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.AllRequestBodies,
				},
				{
					TFSchemaName:             "variant_a_attr",
					TFModelName:              "VariantAAttr",
					APIName:                  "variantAAttr",
					ComputedOptionalRequired: codespec.Computed,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.OmitAlways,
				},
			},
			expectedDiscrim: &codespec.Discriminator{
				PropertyName: "kind",
				Mapping: map[string]codespec.DiscriminatorType{
					"VariantA": {
						Allowed:  []string{"kind", "variant_a_attr"},
						Required: []string{"kind"},
					},
					"VariantB": {
						Allowed:  []string{"kind", "variant_b_attr"},
						Required: []string{"kind"},
					},
				},
			},
		},
		"Nested discriminator with path-based alias renames variant entries": {
			inputConfig: &config.Resource{
				SchemaOptions: config.SchemaOptions{
					Aliases: map[string]string{
						"nestedObj.innerType": "innerKind",
					},
				},
			},
			inputResource: &codespec.Resource{
				Schema: &codespec.Schema{
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "nested_obj",
							TFModelName:              "NestedObj",
							APIName:                  "nestedObj",
							ComputedOptionalRequired: codespec.Computed,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Discriminator: &codespec.Discriminator{
										PropertyName: "inner_type",
										Mapping: map[string]codespec.DiscriminatorType{
											"TypeA": {
												Allowed:  []string{"inner_type", "attr_a"},
												Required: []string{"inner_type"},
											},
										},
									},
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "inner_type",
											TFModelName:              "InnerType",
											APIName:                  "innerType",
											ComputedOptionalRequired: codespec.Computed,
											String:                   &codespec.StringAttribute{},
											ReqBodyUsage:             codespec.OmitAlways,
										},
										{
											TFSchemaName:             "attr_a",
											TFModelName:              "AttrA",
											APIName:                  "attrA",
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
					Create: &codespec.APIOperation{},
					Read:   &codespec.APIOperation{},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "nested_obj",
					TFModelName:              "NestedObj",
					APIName:                  "nestedObj",
					ComputedOptionalRequired: codespec.Computed,
					SingleNested: &codespec.SingleNestedAttribute{
						NestedObject: codespec.NestedAttributeObject{
							Discriminator: &codespec.Discriminator{
								PropertyName: "inner_kind",
								Mapping: map[string]codespec.DiscriminatorType{
									"TypeA": {
										Allowed:  []string{"attr_a", "inner_kind"},
										Required: []string{"inner_kind"},
									},
								},
							},
							Attributes: codespec.Attributes{
								{
									TFSchemaName:             "inner_kind",
									TFModelName:              "InnerKind",
									APIName:                  "innerType", // preserved
									ComputedOptionalRequired: codespec.Computed,
									String:                   &codespec.StringAttribute{},
									ReqBodyUsage:             codespec.OmitAlways,
								},
								{
									TFSchemaName:             "attr_a",
									TFModelName:              "AttrA",
									APIName:                  "attrA",
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
		"Non-dotted alias does NOT rename nested discriminator property": {
			inputConfig: &config.Resource{
				SchemaOptions: config.SchemaOptions{
					Aliases: map[string]string{
						"typeField": "kind",
					},
				},
			},
			inputResource: &codespec.Resource{
				Schema: &codespec.Schema{
					Discriminator: &codespec.Discriminator{
						PropertyName: "type_field",
						Mapping: map[string]codespec.DiscriminatorType{
							"VariantA": {
								Allowed:  []string{"type_field", "root_attr"},
								Required: []string{"type_field"},
							},
						},
					},
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "type_field",
							TFModelName:              "TypeField",
							APIName:                  "typeField",
							ComputedOptionalRequired: codespec.Required,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
						{
							TFSchemaName:             "nested_obj",
							TFModelName:              "NestedObj",
							APIName:                  "nestedObj",
							ComputedOptionalRequired: codespec.Computed,
							SingleNested: &codespec.SingleNestedAttribute{
								NestedObject: codespec.NestedAttributeObject{
									Discriminator: &codespec.Discriminator{
										PropertyName: "type_field",
										Mapping: map[string]codespec.DiscriminatorType{
											"InnerA": {
												Allowed:  []string{"inner_attr", "type_field"},
												Required: []string{"type_field"},
											},
										},
									},
									Attributes: codespec.Attributes{
										{
											TFSchemaName:             "type_field",
											TFModelName:              "TypeField",
											APIName:                  "typeField",
											ComputedOptionalRequired: codespec.Computed,
											String:                   &codespec.StringAttribute{},
											ReqBodyUsage:             codespec.OmitAlways,
										},
										{
											TFSchemaName:             "inner_attr",
											TFModelName:              "InnerAttr",
											APIName:                  "innerAttr",
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
					Create: &codespec.APIOperation{},
					Read:   &codespec.APIOperation{},
				},
			},
			expectedDiscrim: &codespec.Discriminator{
				PropertyName: "kind",
				Mapping: map[string]codespec.DiscriminatorType{
					"VariantA": {
						Allowed:  []string{"kind", "root_attr"},
						Required: []string{"kind"},
					},
				},
			},
			expectedAttributes: codespec.Attributes{
				{
					TFSchemaName:             "kind",
					TFModelName:              "Kind",
					APIName:                  "typeField",
					ComputedOptionalRequired: codespec.Required,
					String:                   &codespec.StringAttribute{},
					ReqBodyUsage:             codespec.AllRequestBodies,
				},
				{
					TFSchemaName:             "nested_obj",
					TFModelName:              "NestedObj",
					APIName:                  "nestedObj",
					ComputedOptionalRequired: codespec.Computed,
					SingleNested: &codespec.SingleNestedAttribute{
						NestedObject: codespec.NestedAttributeObject{
							Discriminator: &codespec.Discriminator{
								// NOT renamed - non-dotted alias only applies at root level
								PropertyName: "type_field",
								Mapping: map[string]codespec.DiscriminatorType{
									"InnerA": {
										Allowed:  []string{"inner_attr", "type_field"},
										Required: []string{"type_field"},
									},
								},
							},
							Attributes: codespec.Attributes{
								{
									// NOT renamed - non-dotted alias only applies at root level
									TFSchemaName:             "type_field",
									TFModelName:              "TypeField",
									APIName:                  "typeField",
									ComputedOptionalRequired: codespec.Computed,
									String:                   &codespec.StringAttribute{},
									ReqBodyUsage:             codespec.OmitAlways,
								},
								{
									TFSchemaName:             "inner_attr",
									TFModelName:              "InnerAttr",
									APIName:                  "innerAttr",
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := codespec.ApplyTransformationsToResource(tc.inputConfig, tc.inputResource)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedAttributes, tc.inputResource.Schema.Attributes)
			if tc.expectedDiscrim != nil {
				assert.Equal(t, tc.expectedDiscrim, tc.inputResource.Schema.Discriminator)
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
			Ignores: []string{"results.ignore_attr"},
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

func TestApplyTransformationsToDataSources_TypeOverride(t *testing.T) {
	inputDataSources := &codespec.DataSources{
		Schema: &codespec.DataSourceSchema{
			SingularDSAttributes: &codespec.Attributes{
				// List to Set
				{
					TFSchemaName:             "list_attr",
					TFModelName:              "ListAttr",
					APIName:                  "listAttr",
					ComputedOptionalRequired: codespec.ComputedOptional,
					CustomType:               codespec.NewCustomListType(codespec.String),
					List: &codespec.ListAttribute{
						ElementType: codespec.String,
					},
					ReqBodyUsage: codespec.AllRequestBodies,
				},
				// Set to List
				{
					TFSchemaName:             "set_attr",
					TFModelName:              "SetAttr",
					APIName:                  "setAttr",
					ComputedOptionalRequired: codespec.ComputedOptional,
					CustomType:               codespec.NewCustomSetType(codespec.String),
					Set: &codespec.SetAttribute{
						ElementType: codespec.String,
					},
					ReqBodyUsage: codespec.AllRequestBodies,
				},
				// ListNested to SetNested
				{
					TFSchemaName:             "nested_list_attr",
					TFModelName:              "NestedListAttr",
					APIName:                  "nestedListAttr",
					ComputedOptionalRequired: codespec.ComputedOptional,
					CustomType:               codespec.NewCustomNestedListType("NestedListAttr"),
					ListNested: &codespec.ListNestedAttribute{
						NestedObject: codespec.NestedAttributeObject{
							Attributes: codespec.Attributes{
								{
									TFSchemaName:             "name",
									TFModelName:              "Name",
									APIName:                  "name",
									ComputedOptionalRequired: codespec.Optional,
									String:                   &codespec.StringAttribute{},
									ReqBodyUsage:             codespec.AllRequestBodies,
								},
							},
						},
					},
					ReqBodyUsage: codespec.AllRequestBodies,
				},
				// SetNested to ListNested
				{
					TFSchemaName:             "nested_set_attr",
					TFModelName:              "NestedSetAttr",
					APIName:                  "nestedSetAttr",
					ComputedOptionalRequired: codespec.ComputedOptional,
					CustomType:               codespec.NewCustomNestedSetType("NestedSetAttr"),
					SetNested: &codespec.SetNestedAttribute{
						NestedObject: codespec.NestedAttributeObject{
							Attributes: codespec.Attributes{
								{
									TFSchemaName:             "value",
									TFModelName:              "Value",
									APIName:                  "value",
									ComputedOptionalRequired: codespec.Optional,
									String:                   &codespec.StringAttribute{},
									ReqBodyUsage:             codespec.AllRequestBodies,
								},
							},
						},
					},
					ReqBodyUsage: codespec.AllRequestBodies,
				},
			},
		},
		Operations: codespec.APIOperations{},
	}

	inputConfig := &config.DataSources{
		SchemaOptions: config.SchemaOptions{
			Overrides: map[string]config.Override{
				"list_attr":        {Type: conversion.Pointer(config.Set)},
				"set_attr":         {Type: conversion.Pointer(config.List)},
				"nested_list_attr": {Type: conversion.Pointer(config.Set)},
				"nested_set_attr":  {Type: conversion.Pointer(config.List)},
			},
		},
	}

	expectedAttributes := codespec.Attributes{
		// List to Set
		{
			TFSchemaName:             "list_attr",
			TFModelName:              "ListAttr",
			APIName:                  "listAttr",
			ComputedOptionalRequired: codespec.ComputedOptional,
			CustomType:               codespec.NewCustomSetType(codespec.String),
			Set: &codespec.SetAttribute{
				ElementType: codespec.String,
			},
			ReqBodyUsage: codespec.AllRequestBodies,
		},
		// Set to List
		{
			TFSchemaName:             "set_attr",
			TFModelName:              "SetAttr",
			APIName:                  "setAttr",
			ComputedOptionalRequired: codespec.ComputedOptional,
			CustomType:               codespec.NewCustomListType(codespec.String),
			List: &codespec.ListAttribute{
				ElementType: codespec.String,
			},
			ReqBodyUsage: codespec.AllRequestBodies,
		},
		// ListNested to SetNested
		{
			TFSchemaName:             "nested_list_attr",
			TFModelName:              "NestedListAttr",
			APIName:                  "nestedListAttr",
			ComputedOptionalRequired: codespec.ComputedOptional,
			CustomType:               codespec.NewCustomNestedSetType("NestedListAttr"),
			SetNested: &codespec.SetNestedAttribute{
				NestedObject: codespec.NestedAttributeObject{
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "name",
							TFModelName:              "Name",
							APIName:                  "name",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
					},
				},
			},
			ReqBodyUsage: codespec.AllRequestBodies,
		},
		// SetNested to ListNested
		{
			TFSchemaName:             "nested_set_attr",
			TFModelName:              "NestedSetAttr",
			APIName:                  "nestedSetAttr",
			ComputedOptionalRequired: codespec.ComputedOptional,
			CustomType:               codespec.NewCustomNestedListType("NestedSetAttr"),
			ListNested: &codespec.ListNestedAttribute{
				NestedObject: codespec.NestedAttributeObject{
					Attributes: codespec.Attributes{
						{
							TFSchemaName:             "value",
							TFModelName:              "Value",
							APIName:                  "value",
							ComputedOptionalRequired: codespec.Optional,
							String:                   &codespec.StringAttribute{},
							ReqBodyUsage:             codespec.AllRequestBodies,
						},
					},
				},
			},
			ReqBodyUsage: codespec.AllRequestBodies,
		},
	}

	err := codespec.ApplyTransformationsToDataSources(inputConfig, inputDataSources)
	require.NoError(t, err)
	assert.Equal(t, expectedAttributes, *inputDataSources.Schema.SingularDSAttributes)
}
