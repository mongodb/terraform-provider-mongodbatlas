package schema

import "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"

type Int64AttrGenerator struct {
	intModel codespec.Int64Attribute
	attr     codespec.Attribute
}

func (i *Int64AttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&i.attr, "schema.Int64Attribute", []CodeStatement{})
}

type Float64AttrGenerator struct {
	floatModel codespec.Float64Attribute
	attr       codespec.Attribute
}

func (f *Float64AttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&f.attr, "schema.Float64Attribute", []CodeStatement{})
}

type StringAttrGenerator struct {
	stringModel codespec.StringAttribute
	attr        codespec.Attribute
}

func (s *StringAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&s.attr, "schema.StringAttribute", []CodeStatement{})
}

type BoolAttrGenerator struct {
	boolModel codespec.BoolAttribute
	attr      codespec.Attribute
}

func (s *BoolAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&s.attr, "schema.BoolAttribute", []CodeStatement{})
}

type NumberAttrGenerator struct {
	numberModel codespec.NumberAttribute
	attr        codespec.Attribute
}

func (s *NumberAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&s.attr, "schema.NumberAttribute", []CodeStatement{})
}

type ListAttrGenerator struct {
	listModel codespec.ListAttribute
	attr      codespec.Attribute
}

func (l *ListAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&l.attr, "schema.ListAttribute", []CodeStatement{ElementTypeProperty(l.listModel.ElementType)})
}

type MapAttrGenerator struct {
	mapModel codespec.MapAttribute
	attr     codespec.Attribute
}

func (m *MapAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&m.attr, "schema.MapAttribute", []CodeStatement{ElementTypeProperty(m.mapModel.ElementType)})
}

type SetAttrGenerator struct {
	setModel codespec.SetAttribute
	attr     codespec.Attribute
}

func (s *SetAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&s.attr, "schema.SetAttribute", []CodeStatement{ElementTypeProperty(s.setModel.ElementType)})
}
