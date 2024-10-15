package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type CodeStatement struct {
	Result  string
	Imports []string
}

func GenerateSchemaAttributes(attrs codespec.Attributes) []CodeStatement {
	result := []CodeStatement{}
	for i := range attrs {
		result = append(result, attribute(&attrs[i]))
	}
	return result
}

func attribute(attr *codespec.Attribute) CodeStatement {
	generator := typeGenerator(attr)

	typeDefinition := generator.TypeDefinition()
	additionalPropertyStatements := generator.TypeSpecificProperties()

	properties := commonProperties(attr)
	imports := []string{"github.com/hashicorp/terraform-plugin-framework/resource/schema"}
	for i := range additionalPropertyStatements {
		properties = append(properties, additionalPropertyStatements[i].Result)
		imports = append(imports, additionalPropertyStatements[i].Imports...)
	}

	name := attr.Name
	propsResultString := strings.Join(properties, ",\n") + ","
	code := fmt.Sprintf(`
	"%s": %s{
		%s
	},`, name, typeDefinition, propsResultString)
	return CodeStatement{
		Result:  code,
		Imports: imports,
	}
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

type schemaAttributeGenerator interface {
	TypeDefinition() string
	TypeSpecificProperties() []CodeStatement
}

func typeGenerator(attr *codespec.Attribute) schemaAttributeGenerator {
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
