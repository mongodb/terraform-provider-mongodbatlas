package codespec

import (
	"context"
	"fmt"

	"github.com/pb33f/libopenapi/orderedmap"
)

func buildResourceAttrs(s *APISpecSchema) (Attributes, error) {
	objectAttributes := Attributes{}

	sortedProperties := orderedmap.SortAlpha(s.Schema.Properties)
	for pair := range orderedmap.Iterate(context.Background(), sortedProperties) {
		name := pair.Key()
		proxy := pair.Value()

		schema, err := BuildSchema(proxy)
		if err != nil {
			return nil, err
		}

		attribute, err := schema.buildResourceAttr(name, s.GetComputability(name))
		if err != nil {
			return nil, err
		}

		if attribute != nil {
			objectAttributes = append(objectAttributes, *attribute)
		}
	}

	return objectAttributes, nil
}

func (s *APISpecSchema) buildResourceAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
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
		return s.buildArrayAttr(name, computability)
	case OASTypeObject:
		if s.Schema.AdditionalProperties != nil && s.Schema.AdditionalProperties.IsA() {
			return s.buildMapAttr(name, computability)
		}
		return s.buildSingleNestedAttr(name, computability)
	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}

func (s *APISpecSchema) buildStringAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	result := &Attribute{
		Name:                     terraformAttrName(name),
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
		Name:                     terraformAttrName(name),
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
			Name:                     terraformAttrName(name),
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
		Name:                     terraformAttrName(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
		Number:                   &NumberAttribute{},
	}, nil
}

func (s *APISpecSchema) buildBoolAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	result := &Attribute{
		Name:                     terraformAttrName(name),
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

func (s *APISpecSchema) buildArrayAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
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
			Name:                     terraformAttrName(name),
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
		objectAttributes, err := buildResourceAttrs(itemSchema)
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

func (s *APISpecSchema) buildMapAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	mapSchema, err := BuildSchema(s.Schema.AdditionalProperties.A)
	if err != nil {
		return nil, err
	}

	result := &Attribute{
		Name:                     terraformAttrName(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
	}

	if mapSchema.Type == OASTypeObject {
		mapAttributes, err := buildResourceAttrs(mapSchema)
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

func (s *APISpecSchema) buildSingleNestedAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	objectAttributes, err := buildResourceAttrs(s)
	if err != nil {
		return nil, err
	}

	return &Attribute{
		Name:                     terraformAttrName(name),
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
