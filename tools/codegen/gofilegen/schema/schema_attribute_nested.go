package schema

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type ListNestedAttrGenerator struct {
	listNestedModel codespec.ListNestedAttribute
	attr            codespec.Attribute
}

func (l *ListNestedAttrGenerator) AttributeCode() (CodeStatement, error) {
	nestedObj, err := nestedObjectProperty(l.listNestedModel.NestedObject)
	if err != nil {
		return CodeStatement{}, err
	}
	return commonAttrStructure(&l.attr, "schema.ListNestedAttribute", "planmodifier.List", []CodeStatement{nestedObj})
}

type SetNestedGenerator struct {
	setNestedModel codespec.SetNestedAttribute
	attr           codespec.Attribute
}

func (l *SetNestedGenerator) AttributeCode() (CodeStatement, error) {
	nestedObj, err := nestedObjectProperty(l.setNestedModel.NestedObject)
	if err != nil {
		return CodeStatement{}, err
	}
	return commonAttrStructure(&l.attr, "schema.SetNestedAttribute", "planmodifier.Set", []CodeStatement{nestedObj})
}

type MapNestedAttrGenerator struct {
	mapNestedModel codespec.MapNestedAttribute
	attr           codespec.Attribute
}

func (m *MapNestedAttrGenerator) AttributeCode() (CodeStatement, error) {
	nestedObj, err := nestedObjectProperty(m.mapNestedModel.NestedObject)
	if err != nil {
		return CodeStatement{}, err
	}
	return commonAttrStructure(&m.attr, "schema.MapNestedAttribute", "planmodifier.Map", []CodeStatement{nestedObj})
}

type SingleNestedAttrGenerator struct {
	singleNestedModel codespec.SingleNestedAttribute
	attr              codespec.Attribute
}

func (l *SingleNestedAttrGenerator) AttributeCode() (CodeStatement, error) {
	attrProp, err := attributesProperty(l.singleNestedModel.NestedObject)
	if err != nil {
		return CodeStatement{}, err
	}
	return commonAttrStructure(&l.attr, "schema.SingleNestedAttribute", "planmodifier.Object", []CodeStatement{attrProp})
}

func attributesProperty(nested codespec.NestedAttributeObject) (CodeStatement, error) {
	attrs, err := GenerateSchemaAttributes(nested.Attributes)
	if err != nil {
		return CodeStatement{}, err
	}
	attributeProperty := fmt.Sprintf(`Attributes: map[string]schema.Attribute{
		%s
	}`, attrs.Code)
	return CodeStatement{
		Code:    attributeProperty,
		Imports: attrs.Imports,
	}, nil
}

func nestedObjectProperty(nested codespec.NestedAttributeObject) (CodeStatement, error) {
	result, err := attributesProperty(nested)
	if err != nil {
		return CodeStatement{}, err
	}
	nestedObj := fmt.Sprintf(`NestedObject: schema.NestedAttributeObject{
		%s,
	}`, result.Code)
	return CodeStatement{
		Code:    nestedObj,
		Imports: result.Imports,
	}, nil
}
