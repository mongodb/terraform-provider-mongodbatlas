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

func generator(attr *codespec.Attribute) attributeGenerator {
	if attr.Int64 != nil {
		return &Int64AttrGenerator{intModel: *attr.Int64, attr: *attr}
	}
	if attr.Float64 != nil {
		return &Float64AttrGenerator{
			floatModel: *attr.Float64,
			attr:       *attr,
		}
	}
	if attr.String != nil {
		return &StringAttrGenerator{
			stringModel: *attr.String,
			attr:        *attr,
		}
	}
	if attr.Bool != nil {
		return &BoolAttrGenerator{
			boolModel: *attr.Bool,
			attr:      *attr,
		}
	}
	if attr.List != nil {
		return &ListAttrGenerator{
			listModel: *attr.List,
			attr:      *attr,
		}
	}
	if attr.ListNested != nil {
		return &ListNestedAttrGenerator{
			listNestedModel: *attr.ListNested,
			attr:            *attr,
		}
	}
	if attr.Map != nil {
		return &MapAttrGenerator{
			mapModel: *attr.Map,
			attr:     *attr,
		}
	}
	if attr.MapNested != nil {
		return &MapNestedAttrGenerator{
			mapNestedModel: *attr.MapNested,
			attr:           *attr,
		}
	}
	if attr.Number != nil {
		return &NumberAttrGenerator{
			numberModel: *attr.Number,
			attr:        *attr,
		}
	}
	if attr.Set != nil {
		return &SetAttrGenerator{
			setModel: *attr.Set,
			attr:     *attr,
		}
	}
	if attr.SetNested != nil {
		return &SetNestedGenerator{
			setNestedModel: *attr.SetNested,
			attr:           *attr,
		}
	}
	if attr.SingleNested != nil {
		return &SingleNestedAttrGenerator{
			singleNestedModel: *attr.SingleNested,
			attr:              *attr,
		}
	}
	if attr.Timeouts != nil {
		return &timeoutAttributeGenerator{
			timeouts: *attr.Timeouts,
		}
	}
	panic("Attribute with unknown type defined when generating schema attribute")
}

// generation of conventional attribute types which have common properties like MarkdownDescription, Computed/Optional/Required, Sensitive
func commonAttrStructure(attr *codespec.Attribute, typeDef string, specificProperties []CodeStatement) CodeStatement {
	properties := commonProperties(attr)
	imports := []string{}
	for i := range specificProperties {
		properties = append(properties, specificProperties[i].Code)
		imports = append(imports, specificProperties[i].Imports...)
	}

	name := attr.Name
	propsResultString := strings.Join(properties, ",\n") + ","
	code := fmt.Sprintf(`"%s": %s{
		%s
	}`, name, typeDef, propsResultString)
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
	if attr.Sensitive {
		result = append(result, "Sensitive: true")
	}
	if attr.CustomType != nil {
		result = append(result, fmt.Sprintf("CustomType: %s", attr.CustomType.Schema))
	}
	return result
}
