package schema

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type ListNestedAttrGenerator struct {
	listNestedModel codespec.ListNestedAttribute
	attr            codespec.Attribute
}

func (l *ListNestedAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&l.attr, "schema.ListNestedAttribute", "planmodifier.List", []CodeStatement{nestedObjectProperty(l.listNestedModel.NestedObject)})
}

type SetNestedGenerator struct {
	setNestedModel codespec.SetNestedAttribute
	attr           codespec.Attribute
}

func (l *SetNestedGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&l.attr, "schema.SetNestedAttribute", "planmodifier.Set", []CodeStatement{nestedObjectProperty(l.setNestedModel.NestedObject)})
}

type MapNestedAttrGenerator struct {
	mapNestedModel codespec.MapNestedAttribute
	attr           codespec.Attribute
}

func (m *MapNestedAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&m.attr, "schema.MapNestedAttribute", "planmodifier.Map", []CodeStatement{nestedObjectProperty(m.mapNestedModel.NestedObject)})
}

type SingleNestedAttrGenerator struct {
	singleNestedModel codespec.SingleNestedAttribute
	attr              codespec.Attribute
}

func (l *SingleNestedAttrGenerator) AttributeCode() CodeStatement {
	return commonAttrStructure(&l.attr, "schema.SingleNestedAttribute", "planmodifier.Object", []CodeStatement{attributesProperty(l.singleNestedModel.NestedObject)})
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
