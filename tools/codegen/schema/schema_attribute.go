package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func GenerateSchemaAttributes(attrs codespec.Attributes) CodeStatement {
	attrsCode := []string{}
	imports := []string{}
	for i := range attrs {
		result := generator(&attrs[i]).AttributeCode()
		attrsCode = append(attrsCode, result.Code)
		imports = append(imports, result.Imports...)
	}
	finalAttrs := strings.Join(attrsCode, ",\n") + ","
	return CodeStatement{
		Code:    finalAttrs,
		Imports: imports,
	}
}

type attributeGenerator interface {
	AttributeCode() CodeStatement
}

type TimeoutAttributeGenerator struct {
	timeouts codespec.TimeoutsAttribute
}

func (s *TimeoutAttributeGenerator) AttributeCode() CodeStatement {
	var optionProperties string
	for op := range s.timeouts.ConfigurableTimeouts {
		switch op {
		case int(codespec.Create):
			optionProperties += "Create: true,"
		case int(codespec.Update):
			optionProperties += "Update: true,"
		case int(codespec.Delete):
			optionProperties += "Delete: true,"
		case int(codespec.Read):
			optionProperties += "Read: true,"
		}
	}
	return CodeStatement{
		Code: fmt.Sprintf(`"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
			%s
		})`, optionProperties),
		Imports: []string{"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"},
	}
}

type ConventionalAttributeGenerator struct {
	typeSpecificCode convetionalTypeSpecificCodeGenerator
	attribute        codespec.Attribute
}

func (s *ConventionalAttributeGenerator) AttributeCode() CodeStatement {
	typeDefinition := s.typeSpecificCode.TypeDefinition()
	additionalPropertyStatements := s.typeSpecificCode.TypeSpecificProperties()

	properties := commonProperties(&s.attribute)
	imports := []string{}
	for i := range additionalPropertyStatements {
		properties = append(properties, additionalPropertyStatements[i].Code)
		imports = append(imports, additionalPropertyStatements[i].Imports...)
	}

	name := s.attribute.Name
	propsResultString := strings.Join(properties, ",\n") + ","
	code := fmt.Sprintf(`"%s": %s{
		%s
	}`, name, typeDefinition, propsResultString)
	return CodeStatement{
		Code:    code,
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

type convetionalTypeSpecificCodeGenerator interface {
	TypeDefinition() string
	TypeSpecificProperties() []CodeStatement
}

func generator(attr *codespec.Attribute) attributeGenerator {
	if attr.Int64 != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &Int64AttrGenerator{model: *attr.Int64},
			attribute:        *attr,
		}
	}
	if attr.Float64 != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &Float64AttrGenerator{model: *attr.Float64},
			attribute:        *attr,
		}
	}
	if attr.String != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &StringAttrGenerator{model: *attr.String},
			attribute:        *attr,
		}
	}
	if attr.Bool != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &BoolAttrGenerator{model: *attr.Bool},
			attribute:        *attr,
		}
	}
	if attr.List != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &ListAttrGenerator{model: *attr.List},
			attribute:        *attr,
		}
	}
	if attr.ListNested != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &ListNestedAttrGenerator{model: *attr.ListNested},
			attribute:        *attr,
		}
	}
	if attr.Map != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &MapAttrGenerator{model: *attr.Map},
			attribute:        *attr,
		}
	}
	if attr.MapNested != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &MapNestedAttrGenerator{model: *attr.MapNested},
			attribute:        *attr,
		}
	}
	if attr.Number != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &NumberAttrGenerator{model: *attr.Number},
			attribute:        *attr,
		}
	}
	if attr.Set != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &SetAttrGenerator{model: *attr.Set},
			attribute:        *attr,
		}
	}
	if attr.SetNested != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &SetNestedGenerator{model: *attr.SetNested},
			attribute:        *attr,
		}
	}
	if attr.SingleNested != nil {
		return &ConventionalAttributeGenerator{
			typeSpecificCode: &SingleNestedAttrGenerator{model: *attr.SingleNested},
			attribute:        *attr,
		}
	}
	if attr.Timeouts != nil {
		return &TimeoutAttributeGenerator{}
	}
	panic("Attribute with unknown type defined when generating schema attribute")
}
