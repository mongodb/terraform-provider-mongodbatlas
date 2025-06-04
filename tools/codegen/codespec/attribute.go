package codespec

import (
	"context"
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/stringcase"
	"github.com/pb33f/libopenapi/orderedmap"
)

func buildResourceAttrs(s *APISpecSchema, isFromRequest bool) (Attributes, error) {
	objectAttributes := Attributes{}

	sortedProperties := orderedmap.SortAlpha(s.Schema.Properties)
	for pair := range orderedmap.Iterate(context.Background(), sortedProperties) {
		name := pair.Key()
		proxy := pair.Value()

		schema, err := BuildSchema(proxy)
		if err != nil {
			return nil, err
		}

		// ignores properties defined in request which are defined with readOnly (common in Atlas API Spec)
		if schema.Schema.ReadOnly != nil && *schema.Schema.ReadOnly && isFromRequest {
			continue
		}

		attribute, err := schema.buildResourceAttr(name, s.GetComputability(name), isFromRequest)
		if err != nil {
			return nil, err
		}

		if attribute != nil {
			objectAttributes = append(objectAttributes, *attribute)
		}
	}
	oneOfRefs := s.Schema.OneOf
	if discriminator := s.Schema.Discriminator; discriminator != nil && orderedmap.Len(discriminator.Mapping) == len(oneOfRefs) {
		discriminatorProperty := discriminator.PropertyName
		ensureComputed(objectAttributes, discriminatorProperty) // the discriminator property should always be computed
		counter := -1
		for pair := range orderedmap.Iterate(context.Background(), discriminator.Mapping) {
			counter++
			name := pair.Key()
			ref := oneOfRefs[counter]
			if ref == nil {
				return nil, fmt.Errorf("discriminator mapping '%s' has no reference", name)
			}
			oneOfSchema, err := BuildSchema(ref)
			if err != nil {
				return nil, err
			}
			allOfRefs := oneOfSchema.Schema.AllOf
			// the first schema in allOf is the discriminator schema, the second is the actual schema
			if len(allOfRefs) != 2 {
				return nil, fmt.Errorf("discriminator mapping '%s' has invalid schema, expected 2 schemas in allOf, got %d", name, len(allOfRefs))
			}
			schema, err := BuildSchema(allOfRefs[1])
			if err != nil {
				return nil, err
			}
			attribute, err := schema.buildResourceAttr(discriminatorProperty+"_"+name, Optional, isFromRequest)
			if err != nil {
				return nil, err
			}
			attribute.Discriminator = &DiscriminatorMapping{
				DiscriminatorProperty: discriminatorProperty,
				DiscriminatorValue:    name,
			}
			objectAttributes = append(objectAttributes, *attribute)
		}
	}

	return objectAttributes, nil
}

func ensureComputed(objectAttributes Attributes, discriminatorProperty string) {
	for i := range objectAttributes {
		attr := &objectAttributes[i]
		if attr.Name == stringcase.SnakeCaseString(discriminatorProperty) {
			attr.ComputedOptionalRequired = Computed
		}
	}
}

func (s *APISpecSchema) buildResourceAttr(name string, computability ComputedOptionalRequired, isFromRequest bool) (*Attribute, error) {
	switch s.Type {
	case OASTypeString:
		return s.buildStringAttr(name, computability)
	case OASTypeInteger:
		return s.buildIntegerAttr(name, computability)
	case OASTypeNumber:
		return s.buildNumberAttr(name, computability)
	case OASTypeBoolean:
		return s.buildBoolAttr(name, computability)
	case OASTypeArray:
		return s.buildArrayAttr(name, computability, isFromRequest)
	case OASTypeObject:
		if s.Schema.AdditionalProperties != nil && s.Schema.AdditionalProperties.IsA() {
			return s.buildMapAttr(name, computability, isFromRequest)
		}
		return s.buildSingleNestedAttr(name, computability, isFromRequest)
	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}

func (s *APISpecSchema) buildStringAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	result := &Attribute{
		Name:                     stringcase.FromCamelCase(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
		Sensitive:                s.IsSensitive(),
		String:                   &StringAttribute{},
	}

	if s.Schema.Default != nil {
		var staticDefault string
		if err := s.Schema.Default.Decode(&staticDefault); err == nil {
			result.ComputedOptionalRequired = ComputedOptional

			result.String.Default = &staticDefault
		}
	}

	return result, nil
}

func (s *APISpecSchema) buildIntegerAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	result := &Attribute{
		Name:                     stringcase.FromCamelCase(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
		Int64:                    &Int64Attribute{},
	}

	if s.Schema.Default != nil {
		var staticDefault int64
		if err := s.Schema.Default.Decode(&staticDefault); err == nil {
			result.ComputedOptionalRequired = ComputedOptional

			result.Int64.Default = &staticDefault
		}
	}

	return result, nil
}

func (s *APISpecSchema) buildNumberAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	if s.Schema.Format == OASFormatDouble || s.Schema.Format == OASFormatFloat {
		result := &Attribute{
			Name:                     stringcase.FromCamelCase(name),
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
			Float64:                  &Float64Attribute{},
		}

		if s.Schema.Default != nil {
			var staticDefault float64
			if err := s.Schema.Default.Decode(&staticDefault); err == nil {
				result.ComputedOptionalRequired = ComputedOptional

				result.Float64.Default = &staticDefault
			}
		}

		return result, nil
	}

	return &Attribute{
		Name:                     stringcase.FromCamelCase(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
		Number:                   &NumberAttribute{},
	}, nil
}

func (s *APISpecSchema) buildBoolAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	result := &Attribute{
		Name:                     stringcase.FromCamelCase(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
		Bool:                     &BoolAttribute{},
	}

	if s.Schema.Default != nil {
		var staticDefault bool
		if err := s.Schema.Default.Decode(&staticDefault); err == nil {
			result.ComputedOptionalRequired = ComputedOptional
			result.Bool.Default = &staticDefault
		}
	}

	return result, nil
}

func (s *APISpecSchema) buildArrayAttr(name string, computability ComputedOptionalRequired, isFromRequest bool) (*Attribute, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array items property, schema doesn't exist: %s", name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A)
	if err != nil {
		return nil, fmt.Errorf("error while building nested schema: %s", name)
	}

	isSet := s.Schema.Format == OASFormatSet || (s.Schema.UniqueItems != nil && *s.Schema.UniqueItems)

	createAttribute := func(nestedObject *NestedAttributeObject, elemType ElemType) *Attribute {
		attr := &Attribute{
			Name:                     stringcase.FromCamelCase(name),
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		}

		if nestedObject != nil {
			if isSet {
				attr.SetNested = &SetNestedAttribute{NestedObject: *nestedObject}
			} else {
				attr.ListNested = &ListNestedAttribute{NestedObject: *nestedObject}
			}
		} else {
			if isSet {
				attr.Set = &SetAttribute{ElementType: elemType}
			} else {
				attr.List = &ListAttribute{ElementType: elemType}
			}
		}

		return attr
	}

	if itemSchema.Type == OASTypeObject {
		objectAttributes, err := buildResourceAttrs(itemSchema, isFromRequest)
		if err != nil {
			return nil, fmt.Errorf("error while building nested schema: %s", name)
		}
		nestedObject := &NestedAttributeObject{Attributes: objectAttributes}

		return createAttribute(nestedObject, Unknown), nil // Using Unknown ElemType as a placeholder for no ElemType
	}

	elemType, err := itemSchema.buildElementType()
	if err != nil {
		return nil, fmt.Errorf("error while building nested schema: %s", name)
	}

	result := createAttribute(nil, elemType)

	if s.Schema.Default != nil {
		var staticDefault bool
		if err := s.Schema.Default.Decode(&staticDefault); err == nil {
			result.ComputedOptionalRequired = ComputedOptional
			result.Bool.Default = &staticDefault
		}
	}

	return result, nil
}

func (s *APISpecSchema) buildMapAttr(name string, computability ComputedOptionalRequired, isFromRequest bool) (*Attribute, error) {
	mapSchema, err := BuildSchema(s.Schema.AdditionalProperties.A)
	if err != nil {
		return nil, err
	}

	result := &Attribute{
		Name:                     stringcase.FromCamelCase(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
	}

	if mapSchema.Type == OASTypeObject {
		mapAttributes, err := buildResourceAttrs(mapSchema, isFromRequest)
		if err != nil {
			return nil, err
		}

		result.MapNested = &MapNestedAttribute{
			NestedObject: NestedAttributeObject{
				Attributes: mapAttributes,
			},
		}
	} else {
		elemType, err := mapSchema.buildElementType()
		if err != nil {
			return nil, err
		}

		result.Map = &MapAttribute{
			ElementType: elemType,
		}
	}

	return result, nil
}

func (s *APISpecSchema) buildSingleNestedAttr(name string, computability ComputedOptionalRequired, isFromRequest bool) (*Attribute, error) {
	objectAttributes, err := buildResourceAttrs(s, isFromRequest)
	if err != nil {
		return nil, err
	}

	return &Attribute{
		Name:                     stringcase.FromCamelCase(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
		SingleNested: &SingleNestedAttribute{
			NestedObject: NestedAttributeObject{
				Attributes: objectAttributes,
			},
		},
	}, nil
}
