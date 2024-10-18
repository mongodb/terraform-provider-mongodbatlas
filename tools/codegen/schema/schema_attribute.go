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

type conventionalAttributeGenerator struct {
	typeSpecificCode convetionalTypeSpecificCodeGenerator
	attribute        codespec.Attribute
}

type convetionalTypeSpecificCodeGenerator interface {
	TypeDefinition() string
	TypeSpecificProperties() []CodeStatement
}

func (s *conventionalAttributeGenerator) AttributeCode() CodeStatement {
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

func generator(attr *codespec.Attribute) attributeGenerator {
	if attr.Int64 != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &Int64AttrGenerator{model: *attr.Int64},
			attribute:        *attr,
		}
	}
	if attr.Float64 != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &Float64AttrGenerator{model: *attr.Float64},
			attribute:        *attr,
		}
	}
	if attr.String != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &StringAttrGenerator{model: *attr.String},
			attribute:        *attr,
		}
	}
	if attr.Bool != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &BoolAttrGenerator{model: *attr.Bool},
			attribute:        *attr,
		}
	}
	if attr.List != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &ListAttrGenerator{model: *attr.List},
			attribute:        *attr,
		}
	}
	if attr.ListNested != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &ListNestedAttrGenerator{model: *attr.ListNested},
			attribute:        *attr,
		}
	}
	if attr.Map != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &MapAttrGenerator{model: *attr.Map},
			attribute:        *attr,
		}
	}
	if attr.MapNested != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &MapNestedAttrGenerator{model: *attr.MapNested},
			attribute:        *attr,
		}
	}
	if attr.Number != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &NumberAttrGenerator{model: *attr.Number},
			attribute:        *attr,
		}
	}
	if attr.Set != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &SetAttrGenerator{model: *attr.Set},
			attribute:        *attr,
		}
	}
	if attr.SetNested != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &SetNestedGenerator{model: *attr.SetNested},
			attribute:        *attr,
		}
	}
	if attr.SingleNested != nil {
		return &conventionalAttributeGenerator{
			typeSpecificCode: &SingleNestedAttrGenerator{model: *attr.SingleNested},
			attribute:        *attr,
		}
	}
	if attr.Timeouts != nil {
		return &timeoutAttributeGenerator{
			timeouts: *attr.Timeouts,
		}
	}
	panic("Attribute with unknown type defined when generating schema attribute")
}
