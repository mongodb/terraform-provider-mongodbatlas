package codespec

import (
	"context"
	"slices"
	"sort"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/pb33f/libopenapi/orderedmap"
)

// extractDiscriminator converts a raw XGenDiscriminator into an intermediate model Discriminator.
// It converts API property names (camelCase) to TF schema names (snake_case),
// excludes the discriminator property itself from variant mappings,
// and filters readOnly properties from the required list (keeping them in allowed).
func extractDiscriminator(schema *APISpecSchema) *Discriminator {
	if schema == nil {
		return nil
	}

	xGenDisc := schema.GetXGenDiscriminator()
	if xGenDisc == nil {
		return nil
	}

	propertyNameSnake := stringcase.ToSnakeCase(xGenDisc.PropertyName)

	mapping := make(map[string]DiscriminatorType, len(xGenDisc.Mapping))
	for discriminatorValue, variantType := range xGenDisc.Mapping {
		var allowed []string
		var required []string

		for _, prop := range variantType.Properties {
			snakeProp := stringcase.ToSnakeCase(prop)
			// Exclude the discriminator property itself from variant mappings
			if snakeProp == propertyNameSnake {
				continue
			}
			allowed = append(allowed, snakeProp)
		}

		for _, prop := range variantType.Required {
			snakeProp := stringcase.ToSnakeCase(prop)
			// Exclude the discriminator property itself
			if snakeProp == propertyNameSnake {
				continue
			}
			// Filter readOnly properties from required (they map to Computed in TF)
			if isReadOnlyProperty(schema, prop) {
				continue
			}
			required = append(required, snakeProp)
		}

		sort.Strings(allowed)
		sort.Strings(required)

		mapping[discriminatorValue] = DiscriminatorType{
			Allowed:  allowed,
			Required: required,
		}
	}

	return &Discriminator{
		PropertyName: propertyNameSnake,
		Mapping:      mapping,
	}
}

// mergeDiscriminators merges discriminator data from two sources (typically request and response schemas).
// allowed = union of both sources (captures response-only computed properties).
// required = union of request sources only; response data never contributes to required.
// If either argument is nil, the other is returned as-is.
func MergeDiscriminators(existing, incoming *Discriminator, incomingIsFromResponse bool) *Discriminator {
	if existing == nil && incoming == nil {
		return nil
	}
	if existing == nil {
		if incomingIsFromResponse {
			// Response-only discriminator: clear all required lists since response properties are Computed
			return clearRequired(incoming)
		}
		return incoming
	}
	if incoming == nil {
		return existing
	}

	// Both are non-nil; merge variant mappings
	merged := &Discriminator{
		PropertyName: existing.PropertyName,
		Mapping:      make(map[string]DiscriminatorType, len(existing.Mapping)),
	}

	// Start with existing variants
	for key, variant := range existing.Mapping {
		merged.Mapping[key] = DiscriminatorType{
			Allowed:  slices.Clone(variant.Allowed),
			Required: slices.Clone(variant.Required),
		}
	}

	// Merge incoming variants
	for key, incomingVariant := range incoming.Mapping {
		if existingVariant, found := merged.Mapping[key]; found {
			// Union allowed lists
			mergedAllowed := unionStrings(existingVariant.Allowed, incomingVariant.Allowed)

			// required comes only from request sources
			mergedRequired := existingVariant.Required
			if !incomingIsFromResponse {
				mergedRequired = unionStrings(existingVariant.Required, incomingVariant.Required)
			}

			merged.Mapping[key] = DiscriminatorType{
				Allowed:  mergedAllowed,
				Required: mergedRequired,
			}
		} else {
			// New variant from incoming source
			newVariant := DiscriminatorType{
				Allowed: slices.Clone(incomingVariant.Allowed),
			}
			if !incomingIsFromResponse {
				newVariant.Required = slices.Clone(incomingVariant.Required)
			}
			merged.Mapping[key] = newVariant
		}
	}

	return merged
}

// isReadOnlyProperty checks whether a property (by its camelCase API name) is readOnly in the schema.
func isReadOnlyProperty(schema *APISpecSchema, apiPropertyName string) bool {
	if schema == nil || schema.Schema == nil || schema.Schema.Properties == nil {
		return false
	}

	for pair := range orderedmap.Iterate(context.Background(), schema.Schema.Properties) {
		if pair.Key() == apiPropertyName {
			propSchema, err := BuildSchema(pair.Value())
			if err != nil {
				return false
			}
			return propSchema.Schema.ReadOnly != nil && *propSchema.Schema.ReadOnly
		}
	}

	return false
}

// clearRequired returns a copy of the discriminator with all required lists emptied.
// Used when a discriminator is sourced exclusively from a response schema.
func clearRequired(disc *Discriminator) *Discriminator {
	if disc == nil {
		return nil
	}
	result := &Discriminator{
		PropertyName: disc.PropertyName,
		Mapping:      make(map[string]DiscriminatorType, len(disc.Mapping)),
	}
	for key, variant := range disc.Mapping {
		result.Mapping[key] = DiscriminatorType{
			Allowed: slices.Clone(variant.Allowed),
			// Required intentionally omitted (nil/empty)
		}
	}
	return result
}

// unionStrings returns the sorted union of two string slices with no duplicates.
func unionStrings(a, b []string) []string {
	seen := make(map[string]bool, len(a)+len(b))
	for _, s := range a {
		seen[s] = true
	}
	for _, s := range b {
		seen[s] = true
	}
	result := make([]string, 0, len(seen))
	for s := range seen {
		result = append(result, s)
	}
	sort.Strings(result)
	return result
}
