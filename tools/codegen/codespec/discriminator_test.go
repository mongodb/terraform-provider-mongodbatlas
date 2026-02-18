package codespec_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/stretchr/testify/assert"
)

func TestAllVariantsEmpty_AllEmpty(t *testing.T) {
	mapping := map[string]codespec.DiscriminatorType{
		"METRIC_A": {Allowed: nil},
		"METRIC_B": {Allowed: []codespec.AttributeName{}},
		"METRIC_C": {},
	}
	assert.True(t, codespec.AllVariantsEmpty(mapping))
}

func TestMergeDiscriminators_BothNil(t *testing.T) {
	result := codespec.MergeDiscriminators(nil, nil, false)
	assert.Nil(t, result)
}

func TestMergeDiscriminators_ExistingNilIncomingRequest(t *testing.T) {
	incoming := &codespec.Discriminator{
		PropertyName: codespec.NewAttributeName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []codespec.AttributeName{codespec.NewAttributeName("attrA")},
				Required: []codespec.AttributeName{codespec.NewAttributeName("attrA")},
			},
		},
	}
	result := codespec.MergeDiscriminators(nil, incoming, false)
	assert.Equal(t, incoming, result)
}

func TestMergeDiscriminators_ExistingNilIncomingResponse(t *testing.T) {
	incoming := &codespec.Discriminator{
		PropertyName: codespec.NewAttributeName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []codespec.AttributeName{codespec.NewAttributeName("attrA"), codespec.NewAttributeName("computedAttr")},
				Required: []codespec.AttributeName{codespec.NewAttributeName("attrA")},
			},
		},
	}
	result := codespec.MergeDiscriminators(nil, incoming, true)
	assert.Equal(t, codespec.NewAttributeName("type").TFSchemaName, result.PropertyName.TFSchemaName)
	assert.Equal(t, []codespec.AttributeName{codespec.NewAttributeName("attrA"), codespec.NewAttributeName("computedAttr")}, result.Mapping["TypeA"].Allowed)
	assert.Empty(t, result.Mapping["TypeA"].Required)
}

func TestMergeDiscriminators_IncomingNil(t *testing.T) {
	existing := &codespec.Discriminator{
		PropertyName: codespec.NewAttributeName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []codespec.AttributeName{codespec.NewAttributeName("attrA")},
				Required: []codespec.AttributeName{codespec.NewAttributeName("attrA")},
			},
		},
	}
	result := codespec.MergeDiscriminators(existing, nil, false)
	assert.Equal(t, existing, result)
}

func TestMergeDiscriminators_MergeAllowedUnion(t *testing.T) {
	existing := &codespec.Discriminator{
		PropertyName: codespec.NewAttributeName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []codespec.AttributeName{codespec.NewAttributeName("attrA"), codespec.NewAttributeName("attrB")},
				Required: []codespec.AttributeName{codespec.NewAttributeName("attrA")},
			},
		},
	}
	incoming := &codespec.Discriminator{
		PropertyName: codespec.NewAttributeName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []codespec.AttributeName{codespec.NewAttributeName("attrB"), codespec.NewAttributeName("attrC")},
				Required: []codespec.AttributeName{codespec.NewAttributeName("attrB")},
			},
		},
	}

	result := codespec.MergeDiscriminators(existing, incoming, true)
	assert.Equal(t, []codespec.AttributeName{codespec.NewAttributeName("attrA"), codespec.NewAttributeName("attrB"), codespec.NewAttributeName("attrC")}, result.Mapping["TypeA"].Allowed)
	assert.Equal(t, []codespec.AttributeName{codespec.NewAttributeName("attrA")}, result.Mapping["TypeA"].Required)
}

func TestMergeDiscriminators_MergeRequiredFromRequest(t *testing.T) {
	existing := &codespec.Discriminator{
		PropertyName: codespec.NewAttributeName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []codespec.AttributeName{codespec.NewAttributeName("attrA")},
				Required: []codespec.AttributeName{codespec.NewAttributeName("attrA")},
			},
		},
	}
	incoming := &codespec.Discriminator{
		PropertyName: codespec.NewAttributeName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []codespec.AttributeName{codespec.NewAttributeName("attrA"), codespec.NewAttributeName("attrB")},
				Required: []codespec.AttributeName{codespec.NewAttributeName("attrA"), codespec.NewAttributeName("attrB")},
			},
		},
	}

	result := codespec.MergeDiscriminators(existing, incoming, false)
	assert.Equal(t, []codespec.AttributeName{codespec.NewAttributeName("attrA"), codespec.NewAttributeName("attrB")}, result.Mapping["TypeA"].Allowed)
	assert.Equal(t, []codespec.AttributeName{codespec.NewAttributeName("attrA"), codespec.NewAttributeName("attrB")}, result.Mapping["TypeA"].Required)
}

func TestMergeDiscriminators_NewVariantFromResponse(t *testing.T) {
	existing := &codespec.Discriminator{
		PropertyName: codespec.NewAttributeName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []codespec.AttributeName{codespec.NewAttributeName("attrA")},
				Required: []codespec.AttributeName{codespec.NewAttributeName("attrA")},
			},
		},
	}
	incoming := &codespec.Discriminator{
		PropertyName: codespec.NewAttributeName("type"),
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed: []codespec.AttributeName{codespec.NewAttributeName("attrA")},
			},
			"TypeB": {
				Allowed:  []codespec.AttributeName{codespec.NewAttributeName("attrB")},
				Required: []codespec.AttributeName{codespec.NewAttributeName("attrB")},
			},
		},
	}

	result := codespec.MergeDiscriminators(existing, incoming, true)
	assert.Equal(t, []codespec.AttributeName{codespec.NewAttributeName("attrA")}, result.Mapping["TypeA"].Required)
	assert.Equal(t, []codespec.AttributeName{codespec.NewAttributeName("attrB")}, result.Mapping["TypeB"].Allowed)
	assert.Empty(t, result.Mapping["TypeB"].Required)
}
