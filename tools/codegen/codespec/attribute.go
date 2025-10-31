package codespec

import (
	"context"
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/pb33f/libopenapi/orderedmap"
)

func buildResourceAttrs(s *APISpecSchema, ancestorsName string, isFromRequest, useCustomNestedTypes bool) (Attributes, error) {
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

		attribute, err := schema.buildResourceAttr(name, ancestorsName, s.GetComputability(name), isFromRequest, useCustomNestedTypes)
		if err != nil {
			return nil, err
		}

		if attribute != nil {
			objectAttributes = append(objectAttributes, *attribute)
		}
	}

	return objectAttributes, nil
}

func (s *APISpecSchema) buildResourceAttr(name, ancestorsName string, computability ComputedOptionalRequired, isFromRequest, useCustomNestedTypes bool) (*Attribute, error) {
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
		return s.buildArrayAttr(name, ancestorsName, computability, isFromRequest, useCustomNestedTypes)
	case OASTypeObject:
		if s.Schema.AdditionalProperties != nil && s.Schema.AdditionalProperties.IsA() {
			return s.buildMapAttr(name, ancestorsName, computability, isFromRequest, useCustomNestedTypes)
		}
		return s.buildSingleNestedAttr(name, ancestorsName, computability, isFromRequest, useCustomNestedTypes)
	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}

func (s *APISpecSchema) buildStringAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	result := &Attribute{
		TFSchemaName:             stringcase.FromCamelCase(name),
		TFModelName:              stringcase.Capitalize(name),
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
		TFSchemaName:             stringcase.FromCamelCase(name),
		TFModelName:              stringcase.Capitalize(name),
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
			TFSchemaName:             stringcase.FromCamelCase(name),
			TFModelName:              stringcase.Capitalize(name),
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
		TFSchemaName:             stringcase.FromCamelCase(name),
		TFModelName:              stringcase.Capitalize(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
		Number:                   &NumberAttribute{},
	}, nil
}

func (s *APISpecSchema) buildBoolAttr(name string, computability ComputedOptionalRequired) (*Attribute, error) {
	result := &Attribute{
		TFSchemaName:             stringcase.FromCamelCase(name),
		TFModelName:              stringcase.Capitalize(name),
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

func (s *APISpecSchema) buildArrayAttr(name, ancestorsName string, computability ComputedOptionalRequired, isFromRequest, useCustomNestedTypes bool) (*Attribute, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array items property, schema doesn't exist: %s", name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A)
	if err != nil {
		return nil, fmt.Errorf("error while building nested schema: %s", name)
	}

	isSet := s.Schema.Format == OASFormatSet || (s.Schema.UniqueItems != nil && *s.Schema.UniqueItems)

	tfModelName := stringcase.Capitalize(name)
	createAttribute := func(nestedObject *NestedAttributeObject, nestedObjectName *string, elemType ElemType) *Attribute {
		var (
			attr = &Attribute{
				TFSchemaName:             stringcase.FromCamelCase(name),
				TFModelName:              tfModelName,
				ComputedOptionalRequired: computability,
				DeprecationMessage:       s.GetDeprecationMessage(),
				Description:              s.GetDescription(),
			}
			isNested      = nestedObject != nil
			isNestedEmpty = isNested && len(nestedObject.Attributes) == 0
		)

		if isNested && isNestedEmpty { // objects without attributes use JSON custom type
			elemType = CustomTypeJSON
		}

		if isNested && !isNestedEmpty {
			if isSet {
				if useCustomNestedTypes {
					attr.CustomType = NewCustomNestedSetType(*nestedObjectName)
				}
				attr.SetNested = &SetNestedAttribute{NestedObject: *nestedObject}
			} else {
				if useCustomNestedTypes {
					attr.CustomType = NewCustomNestedListType(*nestedObjectName)
				}
				attr.ListNested = &ListNestedAttribute{NestedObject: *nestedObject}
			}
		} else {
			if isSet {
				if useCustomNestedTypes {
					attr.CustomType = NewCustomSetType(elemType)
				}
				attr.Set = &SetAttribute{ElementType: elemType}
			} else {
				if useCustomNestedTypes {
					attr.CustomType = NewCustomListType(elemType)
				}
				attr.List = &ListAttribute{ElementType: elemType}
			}
		}

		return attr
	}

	if itemSchema.Type == OASTypeObject {
		fullName := ancestorsName + tfModelName
		objectAttributes, err := buildResourceAttrs(itemSchema, fullName, isFromRequest, useCustomNestedTypes)
		if err != nil {
			return nil, fmt.Errorf("error while building nested schema: %s", name)
		}

		nestedObject := &NestedAttributeObject{Attributes: objectAttributes}
		return createAttribute(nestedObject, &fullName, Unknown), nil // Using Unknown ElemType as a placeholder for no ElemType
	}

	elemType, err := itemSchema.buildElementType()
	if err != nil {
		return nil, fmt.Errorf("error while building nested schema: %s", name)
	}

	result := createAttribute(nil, nil, elemType)

	if s.Schema.Default != nil {
		var staticDefault bool
		if err := s.Schema.Default.Decode(&staticDefault); err == nil {
			result.ComputedOptionalRequired = ComputedOptional
			result.Bool.Default = &staticDefault
		}
	}

	return result, nil
}

func (s *APISpecSchema) buildMapAttr(name, ancestorsName string, computability ComputedOptionalRequired, isFromRequest, useCustomNestedTypes bool) (*Attribute, error) {
	mapSchema, err := BuildSchema(s.Schema.AdditionalProperties.A)
	if err != nil {
		return nil, err
	}

	result := &Attribute{
		TFSchemaName:             stringcase.FromCamelCase(name),
		TFModelName:              stringcase.Capitalize(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
	}

	if mapSchema.Type == OASTypeObject {
		fullName := ancestorsName + result.TFModelName
		mapAttributes, err := buildResourceAttrs(mapSchema, fullName, isFromRequest, useCustomNestedTypes)
		if err != nil {
			return nil, err
		}

		if useCustomNestedTypes {
			result.CustomType = NewCustomNestedMapType(fullName)
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

func (s *APISpecSchema) buildSingleNestedAttr(name, ancestorsName string, computability ComputedOptionalRequired, isFromRequest, useCustomNestedTypes bool) (*Attribute, error) {
	attr := &Attribute{
		TFSchemaName:             stringcase.FromCamelCase(name),
		TFModelName:              stringcase.Capitalize(name),
		ComputedOptionalRequired: computability,
		DeprecationMessage:       s.GetDeprecationMessage(),
		Description:              s.GetDescription(),
	}
	fullName := ancestorsName + attr.TFModelName
	objectAttributes, err := buildResourceAttrs(s, fullName, isFromRequest, useCustomNestedTypes)
	if err != nil {
		return nil, err
	}
	if len(objectAttributes) > 0 {
		if useCustomNestedTypes {
			attr.CustomType = NewCustomObjectType(fullName)
		}
		attr.SingleNested = &SingleNestedAttribute{
			NestedObject: NestedAttributeObject{
				Attributes: objectAttributes,
			},
		}
	} else { // objects without attributes use JSON custom type
		attr.CustomType = &CustomTypeJSONVar
		attr.String = &StringAttribute{}
	}
	return attr, nil
}
