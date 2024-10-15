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
	case OASTypeArray, OASTypeObject:
		return nil, nil
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
