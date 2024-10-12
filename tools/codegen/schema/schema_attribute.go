package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func GenerateSchemaAttributes(attrs codespec.Attributes) []string {
	result := []string{}
	for i := range attrs {
		result = append(result, attribute(&attrs[i]))
	}
	return result
}

func attribute(attr *codespec.Attribute) string {
	generator := typeGenerator(attr)

	typeDefinition := generator.TypeDefinition()
	typeSpecificProps := generator.TypeSpecificProperties()
	generalProps := commonProperties(attr)
	properties := strings.Join(append(generalProps, typeSpecificProps...), ",\n") + ","

	name := attr.Name
	return fmt.Sprintf(`
	"%s": %s{
		%s
	},`, name, typeDefinition, properties)
}

func commonProperties(attr *codespec.Attribute) []string {
	var result []string
	if attr.ComputedOptionalRequired == codespec.Required {
		result = append(result, "Required: true")
	}
	if attr.ComputedOptionalRequired == codespec.Computed || attr.ComputedOptionalRequired == codespec.ComputedOptional {
		result = append(result, "Computed: true")
	}
	if attr.ComputedOptionalRequired == codespec.Optional || attr.ComputedOptionalRequired == codespec.ComputedOptional {
		result = append(result, "Optional: true")
	}
	if attr.Description != nil {
		result = append(result, fmt.Sprintf("MarkdownDescription: %q", *attr.Description))
	}
	if attr.Sensitive != nil && *attr.Sensitive {
		result = append(result, "Sensitive: true")
	}
	return result
}

type schemaAttrGenerator interface {
	TypeDefinition() string
	TypeSpecificProperties() []string
}

func typeGenerator(attr *codespec.Attribute) schemaAttrGenerator {
	if attr.Int64 != nil {
		return &Int64AttrGenerator{model: *attr.Int64}
	}
	if attr.Float64 != nil {
		return &Float64AttrGenerator{model: *attr.Float64}
	}
	if attr.String != nil {
		return &StringAttrGenerator{model: *attr.String}
	}
	if attr.Bool != nil {
		return &BoolAttrGenerator{model: *attr.Bool}
	}
	if attr.List != nil {
		return &ListAttrGenerator{model: *attr.List}
	}
	if attr.ListNested != nil {
		return &ListNestedAttrGenerator{model: *attr.ListNested}
	}
	if attr.Map != nil {
		return &MapAttrGenerator{model: *attr.Map}
	}
	if attr.MapNested != nil {
		return &MapNestedAttrGenerator{model: *attr.MapNested}
	}
	if attr.Number != nil {
		return &NumberAttrGenerator{model: *attr.Number}
	}
	if attr.Object != nil {
		return &ObjectAttrGenerator{model: *attr.Object}
	}
	if attr.Set != nil {
		return &SetAttrGenerator{model: *attr.Set}
	}
	if attr.SetNested != nil {
		return &SetNestedGenerator{model: *attr.SetNested}
	}
	if attr.SingleNested != nil {
		return &SingleNestedAttrGenerator{model: *attr.SingleNested}
	}
	panic("Attribute with unknown type defined")
}

type Int64AttrGenerator struct {
	model codespec.Int64Attribute
}

func (i *Int64AttrGenerator) TypeDefinition() string {
	return "schema.Int64Attribute"
}

func (i *Int64AttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type Float64AttrGenerator struct {
	model codespec.Float64Attribute
}

func (f *Float64AttrGenerator) TypeDefinition() string {
	return "schema.Float64Attribute"
}

func (f *Float64AttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type StringAttrGenerator struct {
	model codespec.StringAttribute
}

func (s *StringAttrGenerator) TypeDefinition() string {
	return "schema.StringAttribute"
}

func (s *StringAttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type BoolAttrGenerator struct {
	model codespec.BoolAttribute
}

func (b *BoolAttrGenerator) TypeDefinition() string {
	return "schema.BoolAttribute"
}

func (b *BoolAttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type ListAttrGenerator struct {
	model codespec.ListAttribute
}

func (l *ListAttrGenerator) TypeDefinition() string {
	return "schema.ListAttribute"
}

func (l *ListAttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type ListNestedAttrGenerator struct {
	model codespec.ListNestedAttribute
}

func (l *ListNestedAttrGenerator) TypeDefinition() string {
	return "schema.ListNestedAttribute"
}

func (l *ListNestedAttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type MapAttrGenerator struct {
	model codespec.MapAttribute
}

func (m *MapAttrGenerator) TypeDefinition() string {
	return "schema.MapAttribute"
}

func (m *MapAttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type MapNestedAttrGenerator struct {
	model codespec.MapNestedAttribute
}

func (m *MapNestedAttrGenerator) TypeDefinition() string {
	return "schema.MapNestedAttribute"
}

func (m *MapNestedAttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type NumberAttrGenerator struct {
	model codespec.NumberAttribute
}

func (n *NumberAttrGenerator) TypeDefinition() string {
	return "schema.NumberAttribute"
}

func (n *NumberAttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type ObjectAttrGenerator struct {
	model codespec.ObjectAttribute
}

func (o *ObjectAttrGenerator) TypeDefinition() string {
	return "schema.ObjectAttribute"
}

func (o *ObjectAttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type SetAttrGenerator struct {
	model codespec.SetAttribute
}

func (s *SetAttrGenerator) TypeDefinition() string {
	return "schema.SetAttribute"
}

func (s *SetAttrGenerator) TypeSpecificProperties() []string {
	return nil
}

type SetNestedGenerator struct {
	model codespec.SetNestedAttribute
}

func (s *SetNestedGenerator) TypeDefinition() string {
	return "schema.SetNestedAttribute"
}

func (s *SetNestedGenerator) TypeSpecificProperties() []string {
	return nil
}

type SingleNestedAttrGenerator struct {
	model codespec.SingleNestedAttribute
}

func (s *SingleNestedAttrGenerator) TypeDefinition() string {
	return "schema.SingleNestedAttribute"
}

func (s *SingleNestedAttrGenerator) TypeSpecificProperties() []string {
	return nil
}
