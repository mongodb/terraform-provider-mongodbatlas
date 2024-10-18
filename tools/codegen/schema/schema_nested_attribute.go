package schema

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type ListNestedAttrGenerator struct {
	model codespec.ListNestedAttribute
}

func (l *ListNestedAttrGenerator) TypeDefinition() string {
	return "schema.ListNestedAttribute"
}

func (l *ListNestedAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return []CodeStatement{getPropertyForNestedObj(l.model.NestedObject)}
}

type SetNestedGenerator struct {
	model codespec.SetNestedAttribute
}

func (s *SetNestedGenerator) TypeDefinition() string {
	return "schema.SetNestedAttribute"
}

func (s *SetNestedGenerator) TypeSpecificProperties() []CodeStatement {
	return []CodeStatement{getPropertyForNestedObj(s.model.NestedObject)}
}

type MapNestedAttrGenerator struct {
	model codespec.MapNestedAttribute
}

func (m *MapNestedAttrGenerator) TypeDefinition() string {
	return "schema.MapNestedAttribute"
}

func (m *MapNestedAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return []CodeStatement{getPropertyForNestedObj(m.model.NestedObject)}
}

type SingleNestedAttrGenerator struct {
	model codespec.SingleNestedAttribute
}

func (s *SingleNestedAttrGenerator) TypeDefinition() string {
	return "schema.SingleNestedAttribute"
}

func (s *SingleNestedAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return []CodeStatement{getPropertyForNestedObj(s.model.NestedObject)}
}

func getPropertyForNestedObj(nested codespec.NestedAttributeObject) CodeStatement {
	attrs := GenerateSchemaAttributes(nested.Attributes)
	attributeProperty := fmt.Sprintf(`Attributes: map[string]schema.Attribute{ 
		%s 
	}`, attrs.Code)
	return CodeStatement{
		Code:    attributeProperty,
		Imports: attrs.Imports,
	}
}
