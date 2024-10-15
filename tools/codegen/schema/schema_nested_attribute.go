package schema

import "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"

type ListNestedAttrGenerator struct {
	model codespec.ListNestedAttribute
}

func (l *ListNestedAttrGenerator) TypeDefinition() string {
	return "schema.ListNestedAttribute"
}

func (l *ListNestedAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return nil
}

type SetNestedGenerator struct {
	model codespec.SetNestedAttribute
}

func (s *SetNestedGenerator) TypeDefinition() string {
	return "schema.SetNestedAttribute"
}

func (s *SetNestedGenerator) TypeSpecificProperties() []CodeStatement {
	return nil
}

type MapNestedAttrGenerator struct {
	model codespec.MapNestedAttribute
}

func (m *MapNestedAttrGenerator) TypeDefinition() string {
	return "schema.MapNestedAttribute"
}

func (m *MapNestedAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return nil
}

type SingleNestedAttrGenerator struct {
	model codespec.SingleNestedAttribute
}

func (s *SingleNestedAttrGenerator) TypeDefinition() string {
	return "schema.SingleNestedAttribute"
}

func (s *SingleNestedAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return nil
}
