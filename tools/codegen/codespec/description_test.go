package codespec_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/stretchr/testify/assert"
)

func TestEnhanceDescriptions_Resource(t *testing.T) {
	attrs := codespec.Attributes{
		{TFSchemaName: "discriminator_prop", Description: conversion.StringPtr("The type.")},
		{TFSchemaName: "base_attr", Description: conversion.StringPtr("Base attr description.")},
		{TFSchemaName: "required_in_all_types", Description: conversion.StringPtr("Required in all types description.")},
		{TFSchemaName: "mixed_required_and_optional", Description: conversion.StringPtr("Mixed required and optional description.")},
		{TFSchemaName: "required_single_type", Description: conversion.StringPtr("Required single type description.")},
		{TFSchemaName: "nil_description", Description: nil},
	}
	disc := &codespec.Discriminator{
		PropertyName: codespec.DiscriminatorAttrName{APIName: "discriminatorProp", TFSchemaName: "discriminator_prop"},
		Mapping: map[string]codespec.DiscriminatorType{
			"ALPHA": {
				Allowed:  []codespec.DiscriminatorAttrName{{APIName: "requiredInAllTypes", TFSchemaName: "required_in_all_types"}, {APIName: "mixedRequiredAndOptional", TFSchemaName: "mixed_required_and_optional"}},
				Required: []codespec.DiscriminatorAttrName{{APIName: "requiredInAllTypes", TFSchemaName: "required_in_all_types"}},
			},
			"BETA": {
				Allowed:  []codespec.DiscriminatorAttrName{{APIName: "requiredInAllTypes", TFSchemaName: "required_in_all_types"}, {APIName: "mixedRequiredAndOptional", TFSchemaName: "mixed_required_and_optional"}},
				Required: []codespec.DiscriminatorAttrName{{APIName: "requiredInAllTypes", TFSchemaName: "required_in_all_types"}, {APIName: "mixedRequiredAndOptional", TFSchemaName: "mixed_required_and_optional"}},
			},
			"GAMMA": {
				Allowed:  []codespec.DiscriminatorAttrName{{APIName: "requiredSingleType", TFSchemaName: "required_single_type"}},
				Required: []codespec.DiscriminatorAttrName{{APIName: "requiredSingleType", TFSchemaName: "required_single_type"}},
			},
		},
	}

	codespec.EnhanceDescriptionsWithDiscriminator(attrs, disc, true)

	assert.Equal(t, "The type.", *attrs[0].Description, "discriminator property itself is skipped")
	assert.Equal(t, "Base attr description.", *attrs[1].Description, "base/common attribute is untouched")
	assert.Equal(t, "Required for discriminator_prop: ALPHA, BETA. Required in all types description.", *attrs[2].Description)
	assert.Equal(t, "Required for discriminator_prop: BETA. Applies to discriminator_prop: ALPHA. Mixed required and optional description.", *attrs[3].Description)
	assert.Equal(t, "Required for discriminator_prop: GAMMA. Required single type description.", *attrs[4].Description)
	assert.Nil(t, attrs[5].Description, "nil description stays nil")
}

func TestEnhanceDescriptions_DataSource(t *testing.T) {
	attrs := codespec.Attributes{
		{TFSchemaName: "type", Description: conversion.StringPtr("The type.")},
		{TFSchemaName: "bucket_name", Description: conversion.StringPtr("Name of the bucket.")},
	}
	disc := &codespec.Discriminator{
		PropertyName: codespec.NewDiscriminatorAttrName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"GCS_LOG_EXPORT": {
				Allowed:  []codespec.DiscriminatorAttrName{codespec.NewDiscriminatorAttrName("bucketName")},
				Required: []codespec.DiscriminatorAttrName{codespec.NewDiscriminatorAttrName("bucketName")},
			},
			"S3_LOG_EXPORT": {
				Allowed:  []codespec.DiscriminatorAttrName{codespec.NewDiscriminatorAttrName("bucketName")},
				Required: []codespec.DiscriminatorAttrName{codespec.NewDiscriminatorAttrName("bucketName")},
			},
		},
	}

	codespec.EnhanceDescriptionsWithDiscriminator(attrs, disc, false)

	assert.Equal(t, "Applies to type: GCS_LOG_EXPORT, S3_LOG_EXPORT. Name of the bucket.", *attrs[1].Description, "data sources always use 'Applies to'")
}

func TestEnhanceDescriptions_NestedDiscriminator(t *testing.T) {
	nestedDisc := &codespec.Discriminator{
		PropertyName: codespec.NewDiscriminatorAttrName("authType"),
		Mapping: map[string]codespec.DiscriminatorType{
			"OAUTH": {
				Allowed:  []codespec.DiscriminatorAttrName{codespec.NewDiscriminatorAttrName("clientId")},
				Required: []codespec.DiscriminatorAttrName{codespec.NewDiscriminatorAttrName("clientId")},
			},
		},
	}

	attrs := codespec.Attributes{
		{
			TFSchemaName: "authentication",
			Description:  conversion.StringPtr("Auth config."),
			SingleNested: &codespec.SingleNestedAttribute{
				NestedObject: codespec.NestedAttributeObject{
					Discriminator: nestedDisc,
					Attributes: codespec.Attributes{
						{TFSchemaName: "auth_type", Description: conversion.StringPtr("The auth type.")},
						{TFSchemaName: "client_id", Description: conversion.StringPtr("OAuth client ID.")},
					},
				},
			},
		},
	}

	codespec.EnhanceDescriptionsWithDiscriminator(attrs, nil, true)

	assert.Equal(t, "Auth config.", *attrs[0].Description, "parent not modified when root disc is nil")
	nestedAttrs := attrs[0].SingleNested.NestedObject.Attributes
	assert.Equal(t, "The auth type.", *nestedAttrs[0].Description, "nested discriminator property is skipped")
	assert.Equal(t, "Required for auth_type: OAUTH. OAuth client ID.", *nestedAttrs[1].Description)
}
