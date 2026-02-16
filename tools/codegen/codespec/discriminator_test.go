package codespec_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/stretchr/testify/assert"
)

func TestAllVariantsEmpty_AllEmpty(t *testing.T) {
	mapping := map[string]codespec.DiscriminatorType{
		"METRIC_A": {Allowed: nil},
		"METRIC_B": {Allowed: []string{}},
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
		PropertyName: "type",
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []string{"attr_a"},
				Required: []string{"attr_a"},
			},
		},
	}
	result := codespec.MergeDiscriminators(nil, incoming, false)
	assert.Equal(t, incoming, result)
}

func TestMergeDiscriminators_ExistingNilIncomingResponse(t *testing.T) {
	incoming := &codespec.Discriminator{
		PropertyName: "type",
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []string{"attr_a", "computed_attr"},
				Required: []string{"attr_a"},
			},
		},
	}
	result := codespec.MergeDiscriminators(nil, incoming, true)
	// Response-only: required should be cleared
	assert.Equal(t, "type", result.PropertyName)
	assert.Equal(t, []string{"attr_a", "computed_attr"}, result.Mapping["TypeA"].Allowed)
	assert.Empty(t, result.Mapping["TypeA"].Required)
}

func TestMergeDiscriminators_IncomingNil(t *testing.T) {
	existing := &codespec.Discriminator{
		PropertyName: "type",
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []string{"attr_a"},
				Required: []string{"attr_a"},
			},
		},
	}
	result := codespec.MergeDiscriminators(existing, nil, false)
	assert.Equal(t, existing, result)
}

func TestMergeDiscriminators_MergeAllowedUnion(t *testing.T) {
	existing := &codespec.Discriminator{
		PropertyName: "type",
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []string{"attr_a", "attr_b"},
				Required: []string{"attr_a"},
			},
		},
	}
	incoming := &codespec.Discriminator{
		PropertyName: "type",
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []string{"attr_b", "attr_c"},
				Required: []string{"attr_b"},
			},
		},
	}

	// Incoming is from response: allowed = union, required = existing only
	result := codespec.MergeDiscriminators(existing, incoming, true)
	assert.Equal(t, []string{"attr_a", "attr_b", "attr_c"}, result.Mapping["TypeA"].Allowed)
	assert.Equal(t, []string{"attr_a"}, result.Mapping["TypeA"].Required)
}

func TestMergeDiscriminators_MergeRequiredFromRequest(t *testing.T) {
	existing := &codespec.Discriminator{
		PropertyName: "type",
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []string{"attr_a"},
				Required: []string{"attr_a"},
			},
		},
	}
	incoming := &codespec.Discriminator{
		PropertyName: "type",
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []string{"attr_a", "attr_b"},
				Required: []string{"attr_a", "attr_b"},
			},
		},
	}

	// Incoming is from request: required = union of both
	result := codespec.MergeDiscriminators(existing, incoming, false)
	assert.Equal(t, []string{"attr_a", "attr_b"}, result.Mapping["TypeA"].Allowed)
	assert.Equal(t, []string{"attr_a", "attr_b"}, result.Mapping["TypeA"].Required)
}

func TestMergeDiscriminators_NewVariantFromResponse(t *testing.T) {
	existing := &codespec.Discriminator{
		PropertyName: "type",
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed:  []string{"attr_a"},
				Required: []string{"attr_a"},
			},
		},
	}
	incoming := &codespec.Discriminator{
		PropertyName: "type",
		Mapping: map[string]codespec.DiscriminatorType{
			"TypeA": {
				Allowed: []string{"attr_a"},
			},
			"TypeB": {
				Allowed:  []string{"attr_b"},
				Required: []string{"attr_b"},
			},
		},
	}

	result := codespec.MergeDiscriminators(existing, incoming, true)
	// TypeA: required preserved from existing
	assert.Equal(t, []string{"attr_a"}, result.Mapping["TypeA"].Required)
	// TypeB: new from response, required should be empty
	assert.Equal(t, []string{"attr_b"}, result.Mapping["TypeB"].Allowed)
	assert.Empty(t, result.Mapping["TypeB"].Required)
}
