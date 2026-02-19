package codespec

import (
	"context"
	"slices"
	"sort"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/pb33f/libopenapi/orderedmap"
)

// NewDiscriminatorAttrName creates a DiscriminatorAttrName from a camelCase API property name,
// deriving the TF schema name via snake_case conversion.
func NewDiscriminatorAttrName(apiName string) DiscriminatorAttrName {
	return DiscriminatorAttrName{
		APIName:      apiName,
		TFSchemaName: stringcase.ToSnakeCase(apiName),
	}
}

// extractDiscriminator converts a raw XGenDiscriminator into an intermediate model Discriminator.
// It converts API property names (camelCase) to TF schema names (snake_case),
// excludes the discriminator property itself from variant mappings,
// and filters readOnly properties from the required list (keeping them in allowed).
//
// Returns nil when all variant mappings have empty allowed lists. This happens when
// polymorphism is at the value level (different enum values for a specific property
// like `units`) rather than at the structural level (different properties per variant).
// In such cases the discriminator acts as a pure enum constraint on the discriminator
// property itself and carries no actionable per-variant metadata for Terraform.
func extractDiscriminator(schema *APISpecSchema) *Discriminator {
	if schema == nil {
		return nil
	}

	xGenDisc := schema.GetXGenDiscriminator()
	if xGenDisc == nil {
		return nil
	}

	propertyName := NewDiscriminatorAttrName(xGenDisc.PropertyName)

	mapping := make(map[string]DiscriminatorType, len(xGenDisc.Mapping))
	for discriminatorValue, variantType := range xGenDisc.Mapping {
		var allowed []DiscriminatorAttrName
		var required []DiscriminatorAttrName

		for _, prop := range variantType.Properties {
			if prop == xGenDisc.PropertyName {
				continue
			}
			allowed = append(allowed, NewDiscriminatorAttrName(prop))
		}

		for _, prop := range variantType.Required {
			if prop == xGenDisc.PropertyName {
				continue
			}
			if isReadOnlyProperty(schema, prop) {
				continue
			}
			required = append(required, NewDiscriminatorAttrName(prop))
		}

		sortDiscriminatorAttrNames(allowed)
		sortDiscriminatorAttrNames(required)

		mapping[discriminatorValue] = DiscriminatorType{
			Allowed:  allowed,
			Required: required,
		}
	}

	if AllVariantsEmpty(mapping) {
		return nil
	}

	return &Discriminator{
		PropertyName: propertyName,
		Mapping:      mapping,
	}
}

// MergeDiscriminators merges discriminator data from two sources (typically request and response schemas).
// allowed = union of both sources (captures response-only computed properties).
// required = union of request sources only; response data never contributes to required.
// If either argument is nil, the other is returned as-is.
func MergeDiscriminators(existing, incoming *Discriminator, incomingIsFromResponse bool) *Discriminator {
	if existing == nil && incoming == nil {
		return nil
	}
	if existing == nil {
		if incomingIsFromResponse {
			return clearRequired(incoming)
		}
		return incoming
	}
	if incoming == nil {
		return existing
	}

	merged := &Discriminator{
		PropertyName: existing.PropertyName,
		Mapping:      make(map[string]DiscriminatorType, len(existing.Mapping)),
	}

	for key, variant := range existing.Mapping {
		merged.Mapping[key] = DiscriminatorType{
			Allowed:  slices.Clone(variant.Allowed),
			Required: slices.Clone(variant.Required),
		}
	}

	for key, incomingVariant := range incoming.Mapping {
		if existingVariant, found := merged.Mapping[key]; found {
			mergedAllowed := unionDiscriminatorAttrNames(existingVariant.Allowed, incomingVariant.Allowed)

			mergedRequired := existingVariant.Required
			if !incomingIsFromResponse {
				mergedRequired = unionDiscriminatorAttrNames(existingVariant.Required, incomingVariant.Required)
			}

			merged.Mapping[key] = DiscriminatorType{
				Allowed:  mergedAllowed,
				Required: mergedRequired,
			}
		} else {
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
		}
	}
	return result
}

func AllVariantsEmpty(mapping map[string]DiscriminatorType) bool {
	for _, variant := range mapping {
		if len(variant.Allowed) > 0 {
			return false
		}
	}
	return true
}

// unionDiscriminatorAttrNames returns the sorted union of two DiscriminatorAttrName slices, deduplicating by APIName.
func unionDiscriminatorAttrNames(a, b []DiscriminatorAttrName) []DiscriminatorAttrName {
	seen := make(map[string]DiscriminatorAttrName, len(a)+len(b))
	for _, n := range a {
		seen[n.APIName] = n
	}
	for _, n := range b {
		seen[n.APIName] = n
	}
	result := make([]DiscriminatorAttrName, 0, len(seen))
	for _, n := range seen {
		result = append(result, n)
	}
	sortDiscriminatorAttrNames(result)
	return result
}

func sortDiscriminatorAttrNames(names []DiscriminatorAttrName) {
	sort.Slice(names, func(i, j int) bool {
		return names[i].TFSchemaName < names[j].TFSchemaName
	})
}
