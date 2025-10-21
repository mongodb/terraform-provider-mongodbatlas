package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func GenerateSchemaAttributes(resource *codespec.Resource) CodeStatement {
	generateTimeoutAttributes(resource)
	return generateSchemaAttributesFromList(resource.Schema.Attributes)
}

// generateTimeoutAttributes adds timeouts attribute based on wait blocks.
func generateTimeoutAttributes(resource *codespec.Resource) {
	ops := &resource.Operations
	var result []codespec.Operation
	if ops.Create.Wait != nil {
		result = append(result, codespec.Create)
	}
	if ops.Update.Wait != nil {
		result = append(result, codespec.Update)
	}
	if ops.Read.Wait != nil {
		result = append(result, codespec.Read)
	}
	if ops.Delete != nil && ops.Delete.Wait != nil {
		result = append(result, codespec.Delete)
	}
	if len(result) > 0 {
		resource.Schema.Attributes = append(resource.Schema.Attributes, codespec.Attribute{
			TFSchemaName: "timeouts",
			TFModelName:  "Timeouts",
			Timeouts:     &codespec.TimeoutsAttribute{ConfigurableTimeouts: result},
			ReqBodyUsage: codespec.OmitAlways,
		})
	}
}

func generateSchemaAttributesFromList(attrs codespec.Attributes) CodeStatement {
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
func commonAttrStructure(attr *codespec.Attribute, attrDefType, planModifierType string, specificProperties []CodeStatement) CodeStatement {
	properties := commonProperties(attr, planModifierType)
	properties = append(properties, specificProperties...)

	name := attr.TFSchemaName
	propsStmts := GroupCodeStatements(properties, func(properties []string) string {
		return strings.Join(properties, ",\n") + ","
	})
	code := fmt.Sprintf(`"%s": %s{
		%s
	}`, name, attrDefType, propsStmts.Code)
	return CodeStatement{
		Code:    code,
		Imports: propsStmts.Imports,
	}
}

func commonProperties(attr *codespec.Attribute, planModifierType string) []CodeStatement {
	var result []CodeStatement
	if attr.ComputedOptionalRequired == codespec.Required {
		result = append(result, CodeStatement{Code: "Required: true"})
	}
	if attr.ComputedOptionalRequired == codespec.Computed || attr.ComputedOptionalRequired == codespec.ComputedOptional {
		result = append(result, CodeStatement{Code: "Computed: true"})
	}
	if attr.ComputedOptionalRequired == codespec.Optional || attr.ComputedOptionalRequired == codespec.ComputedOptional {
		result = append(result, CodeStatement{Code: "Optional: true"})
	}
	if attr.Description != nil {
		result = append(result, CodeStatement{Code: fmt.Sprintf("MarkdownDescription: %q", *attr.Description)})
	}
	if attr.Sensitive {
		result = append(result, CodeStatement{Code: "Sensitive: true"})
	}
	if attr.CustomType != nil {
		result = append(result, CodeStatement{
			Code:    fmt.Sprintf("CustomType: %s", attr.CustomType.Schema),
			Imports: []string{"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"},
		})
	}
	if attr.CreateOnly { // as of now this is the only property which implies defining plan modifiers
		result = append(result, CodeStatement{
			Code: fmt.Sprintf("PlanModifiers: []%s{customplanmodifier.CreateOnly()}", planModifierType),
			Imports: []string{
				"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier",
				"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier",
			},
		})
	}
	return result
}
