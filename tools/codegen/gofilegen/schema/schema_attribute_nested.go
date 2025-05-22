package schema

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type ListNestedAttrGenerator struct {
	attr            codespec.Attribute
	listNestedModel codespec.ListNestedAttribute
}

func (l *ListNestedAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&l.attr, "schema.ListNestedAttribute", []CodeStatement{nestedObjectProperty(l.listNestedModel.NestedObject)})
}

type SetNestedGenerator struct {
	attr           codespec.Attribute
	setNestedModel codespec.SetNestedAttribute
}

func (l *SetNestedGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&l.attr, "schema.SetNestedAttribute", []CodeStatement{nestedObjectProperty(l.setNestedModel.NestedObject)})
}

type MapNestedAttrGenerator struct {
	attr           codespec.Attribute
	mapNestedModel codespec.MapNestedAttribute
}

func (m *MapNestedAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&m.attr, "schema.MapNestedAttribute", []CodeStatement{nestedObjectProperty(m.mapNestedModel.NestedObject)})
}

type SingleNestedAttrGenerator struct {
	attr              codespec.Attribute
	singleNestedModel codespec.SingleNestedAttribute
}

func (l *SingleNestedAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&l.attr, "schema.SingleNestedAttribute", []CodeStatement{attributesProperty(l.singleNestedModel.NestedObject)})
}

func attributesProperty(nested codespec.NestedAttributeObject) CodeStatement {
	attrs := GenerateSchemaAttributes(nested.Attributes)
	attributeProperty := fmt.Sprintf(`Attributes: map[string]schema.Attribute{ 
		%s 
	}`, attrs.Code)
	return CodeStatement{
		Code:    attributeProperty,
		Imports: attrs.Imports,
	}
}

func nestedObjectProperty(nested codespec.NestedAttributeObject) CodeStatement {
	result := attributesProperty(nested)
	nestedObj := fmt.Sprintf(`NestedObject: schema.NestedAttributeObject{
		%s,
	}`, result.Code)
	return CodeStatement{
		Code:    nestedObj,
		Imports: result.Imports,
	}
}
