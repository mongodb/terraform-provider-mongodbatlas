package schema

import "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"

type Int64AttrGenerator struct {
	model codespec.Int64Attribute
}

func (i *Int64AttrGenerator) TypeDefinition() string {
	return "schema.Int64Attribute"
}

func (i *Int64AttrGenerator) TypeSpecificProperties() []CodeStatement {
	return nil
}

type Float64AttrGenerator struct {
	model codespec.Float64Attribute
}

func (f *Float64AttrGenerator) TypeDefinition() string {
	return "schema.Float64Attribute"
}

func (f *Float64AttrGenerator) TypeSpecificProperties() []CodeStatement {
	return nil
}

type StringAttrGenerator struct {
	model codespec.StringAttribute
}

func (s *StringAttrGenerator) TypeDefinition() string {
	return "schema.StringAttribute"
}

func (s *StringAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return nil
}

type BoolAttrGenerator struct {
	model codespec.BoolAttribute
}

func (b *BoolAttrGenerator) TypeDefinition() string {
	return "schema.BoolAttribute"
}

func (b *BoolAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return nil
}

type NumberAttrGenerator struct {
	model codespec.NumberAttribute
}

func (n *NumberAttrGenerator) TypeDefinition() string {
	return "schema.NumberAttribute"
}

func (n *NumberAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return nil
}

type ListAttrGenerator struct {
	model codespec.ListAttribute
}

func (l *ListAttrGenerator) TypeDefinition() string {
	return "schema.ListAttribute"
}

func (l *ListAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return []CodeStatement{ElementTypeProperty(l.model.ElementType)}
}

type MapAttrGenerator struct {
	model codespec.MapAttribute
}

func (m *MapAttrGenerator) TypeDefinition() string {
	return "schema.MapAttribute"
}

func (m *MapAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return []CodeStatement{ElementTypeProperty(m.model.ElementType)}
}

type SetAttrGenerator struct {
	model codespec.SetAttribute
}

func (s *SetAttrGenerator) TypeDefinition() string {
	return "schema.SetAttribute"
}

func (s *SetAttrGenerator) TypeSpecificProperties() []CodeStatement {
	return []CodeStatement{ElementTypeProperty(s.model.ElementType)}
}
