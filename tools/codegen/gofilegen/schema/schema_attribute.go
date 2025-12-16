package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

// GenerateSchemaAttributes generates schema attributes for resource schemas.
func GenerateSchemaAttributes(attrs codespec.Attributes) CodeStatement {
	return generateSchemaAttributesWithGenerator(attrs, generator)
}

// GenerateDataSourceSchemaAttributes generates schema attributes for data source schemas.
// Data source attributes use dsschema types instead of resource schema types.
func GenerateDataSourceSchemaAttributes(attrs codespec.Attributes) CodeStatement {
	return generateSchemaAttributesForDataSource(attrs, "DS")
}

// GeneratePluralDataSourceSchemaAttributes generates schema attributes for plural data source schemas.
// Plural data source attributes use dsschema types and PluralDS prefix for nested models.
func GeneratePluralDataSourceSchemaAttributes(attrs codespec.Attributes) CodeStatement {
	return generateSchemaAttributesForDataSource(attrs, "PluralDS")
}

func generateSchemaAttributesForDataSource(attrs codespec.Attributes, dsPrefix string) CodeStatement {
	return generateSchemaAttributesWithGenerator(attrs, func(attr *codespec.Attribute) attributeGenerator {
		return dataSourceAttrGeneratorWithPrefix(attr, dsPrefix)
	})
}

// generateSchemaAttributesWithGenerator is the shared implementation for schema attribute generation.
func generateSchemaAttributesWithGenerator(attrs codespec.Attributes, genFunc func(*codespec.Attribute) attributeGenerator) CodeStatement {
	attrsCode := []string{}
	imports := []string{}
	for i := range attrs {
		result := genFunc(&attrs[i]).AttributeCode()
		attrsCode = append(attrsCode, result.Code)
		imports = append(imports, result.Imports...)
	}
	finalAttrs := strings.Join(attrsCode, ",\n") + ","
	return CodeStatement{
		Code:    finalAttrs,
		Imports: imports,
	}
}

// dataSourceAttrGenerator wraps resource generators to produce data source schema code.
func dataSourceAttrGenerator(attr *codespec.Attribute) attributeGenerator {
	return dataSourceAttrGeneratorWithPrefix(attr, "DS")
}

// dataSourceAttrGeneratorWithPrefix wraps resource generators to produce data source schema code with a specific prefix.
func dataSourceAttrGeneratorWithPrefix(attr *codespec.Attribute, dsPrefix string) attributeGenerator {
	return &dsAttrGeneratorWrapper{inner: generator(attr), attr: attr, dsPrefix: dsPrefix}
}

// dsAttrGeneratorWrapper wraps resource attribute generators to produce data source schema code.
// It replaces "schema." with "dsschema." in the generated code and filters out resource-specific imports.
type dsAttrGeneratorWrapper struct {
	inner    attributeGenerator
	attr     *codespec.Attribute
	dsPrefix string // "DS" for singular, "PluralDS" for plural data sources
}

func (g *dsAttrGeneratorWrapper) AttributeCode() CodeStatement {
	result := g.inner.AttributeCode()
	// Replace schema. with dsschema. for data source schemas
	result.Code = strings.ReplaceAll(result.Code, "schema.", "dsschema.")
	// Add DS prefix to nested model references in CustomType (e.g., TFNestedObjectAttrModel -> TFDSNestedObjectAttrModel or TFPluralDSNestedObjectAttrModel)
	// This ensures data source schemas reference their own nested models instead of resource models.
	result.Code = strings.ReplaceAll(result.Code, "[TF", "[TF"+g.dsPrefix)
	// Filter out resource-specific imports (data sources don't need plan modifiers)
	var filteredImports []string
	for _, imp := range result.Imports {
		if imp == "github.com/hashicorp/terraform-plugin-framework/resource/schema" {
			continue
		}
		if strings.Contains(imp, "planmodifier") || strings.Contains(imp, "customplanmodifier") {
			continue
		}
		filteredImports = append(filteredImports, imp)
	}
	result.Imports = filteredImports
	return result
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
		var imports []string
		switch attr.CustomType.Package {
		case codespec.JSONTypesPkg:
			imports = append(imports, "github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes")
		case codespec.CustomTypesPkg:
			imports = append(imports, "github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes")
		}

		result = append(result, CodeStatement{
			Code:    fmt.Sprintf("CustomType: %s", attr.CustomType.Schema),
			Imports: imports,
		})
	}
	if attr.CreateOnly { // As of now this is the only property which implies defining plan modifiers.
		planModifierImports := []string{
			"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier",
			"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier",
		}
		code := fmt.Sprintf("PlanModifiers: []%s{customplanmodifier.CreateOnly()}", planModifierType)

		// For bool attributes with create-only and default value, use CreateOnlyBoolWithDefault
		if attr.Bool != nil && attr.Bool.Default != nil {
			code = fmt.Sprintf("PlanModifiers: []%s{customplanmodifier.CreateOnlyBoolWithDefault(%t)}", planModifierType, *attr.Bool.Default)
		}

		result = append(result, CodeStatement{
			Code:    code,
			Imports: planModifierImports,
		})
	}
	return result
}
