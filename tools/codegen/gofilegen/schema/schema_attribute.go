package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

// GenerateSchemaAttributes generates schema attributes for resource schemas.
// The disc parameter provides the discriminator for the current schema level (root or nested),
// allowing the generator to attach a ValidateDiscriminator validator to the discriminator property.
func GenerateSchemaAttributes(attrs codespec.Attributes, disc *codespec.Discriminator) (CodeStatement, error) {
	return generateSchemaAttributesWithGenerator(attrs, disc, generator)
}

// GenerateDataSourceSchemaAttributes generates schema attributes for data source schemas.
// Data source attributes use dsschema types instead of resource schema types.
func GenerateDataSourceSchemaAttributes(attrs codespec.Attributes) (CodeStatement, error) {
	return generateSchemaAttributesForDataSource(attrs, "DS")
}

// GeneratePluralDataSourceSchemaAttributes generates schema attributes for plural data source schemas.
// Plural data source attributes use dsschema types and PluralDS prefix for nested models.
func GeneratePluralDataSourceSchemaAttributes(attrs codespec.Attributes) (CodeStatement, error) {
	return generateSchemaAttributesForDataSource(attrs, "PluralDS")
}

func generateSchemaAttributesForDataSource(attrs codespec.Attributes, dsPrefix string) (CodeStatement, error) {
	return generateSchemaAttributesWithGenerator(attrs, nil, func(attr *codespec.Attribute, _ *codespec.Discriminator) attributeGenerator {
		return dataSourceAttrGeneratorWithPrefix(attr, dsPrefix)
	})
}

// generateSchemaAttributesWithGenerator is the shared implementation for schema attribute generation.
// When disc is non-nil and an attribute matches the discriminator property name, the discriminator
// is passed to that attribute's generator so it can emit a ValidateDiscriminator validator.
func generateSchemaAttributesWithGenerator(attrs codespec.Attributes, disc *codespec.Discriminator, genFunc func(*codespec.Attribute, *codespec.Discriminator) attributeGenerator) (CodeStatement, error) {
	attrsCode := []string{}
	imports := []string{}
	for i := range attrs {
		var attrDisc *codespec.Discriminator
		if disc != nil && attrs[i].TFSchemaName == disc.PropertyName.TFSchemaName {
			attrDisc = disc
		}
		result, err := genFunc(&attrs[i], attrDisc).AttributeCode()
		if err != nil {
			return CodeStatement{}, fmt.Errorf("failed to generate attribute '%s': %w", attrs[i].TFSchemaName, err)
		}
		attrsCode = append(attrsCode, result.Code)
		imports = append(imports, result.Imports...)
	}
	finalAttrs := strings.Join(attrsCode, ",\n") + ","
	return CodeStatement{
		Code:    finalAttrs,
		Imports: imports,
	}, nil
}

// dataSourceAttrGeneratorWithPrefix wraps resource generators to produce data source schema code with a specific prefix.
// Data sources never have validators, so the discriminator is always nil.
func dataSourceAttrGeneratorWithPrefix(attr *codespec.Attribute, dsPrefix string) attributeGenerator {
	return &dsAttrGeneratorWrapper{inner: generator(attr, nil), attr: attr, dsPrefix: dsPrefix}
}

// dsAttrGeneratorWrapper wraps resource attribute generators to produce data source schema code.
// It replaces "schema." with "dsschema." in the generated code and filters out resource-specific imports.
type dsAttrGeneratorWrapper struct {
	inner    attributeGenerator
	attr     *codespec.Attribute
	dsPrefix string // "DS" for singular, "PluralDS" for plural data sources
}

func (g *dsAttrGeneratorWrapper) AttributeCode() (CodeStatement, error) {
	result, err := g.inner.AttributeCode()
	if err != nil {
		return CodeStatement{}, err
	}
	// Replace schema. with dsschema. for data source schemas
	result.Code = strings.ReplaceAll(result.Code, "schema.", "dsschema.")
	// Add DS prefix to nested model references in CustomType (e.g., TFNestedObjectAttrModel -> TFDSNestedObjectAttrModel or TFPluralDSNestedObjectAttrModel)
	// This ensures data source schemas reference their own nested models instead of resource models.
	result.Code = strings.ReplaceAll(result.Code, "[TF", "[TF"+g.dsPrefix)
	// Filter out resource-specific imports (data sources don't need plan modifiers or validators)
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
	return result, nil
}

type attributeGenerator interface {
	AttributeCode() (CodeStatement, error)
}

func generator(attr *codespec.Attribute, disc *codespec.Discriminator) attributeGenerator {
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
			stringModel:   *attr.String,
			attr:          *attr,
			discriminator: disc,
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

// TypeSpecificStatements bundles the type-varying inputs that each attribute generator provides
// to commonAttrStructure. This keeps the function signature stable as new per-type concerns are added.
type TypeSpecificStatements struct {
	AttrType         string          // schema attribute type, e.g. "schema.StringAttribute"
	PlanModifierType string          // plan modifier slice element type, e.g. "planmodifier.String"
	ValidatorType    string          // validator slice element type, e.g. "validator.String"
	Validators       []CodeStatement // individual validator entries; commonProperties wraps them in Validators: []validatorType{...}
	Properties       []CodeStatement // other type-specific properties (ElementType, NestedObject, etc.)
}

// commonAttrStructure generates the full attribute declaration by combining common properties
// (Computed/Optional/Required, description, plan modifiers, validators) with type-specific statements.
func commonAttrStructure(attr *codespec.Attribute, typeStatements *TypeSpecificStatements) (CodeStatement, error) {
	properties, err := commonProperties(attr, typeStatements)
	if err != nil {
		return CodeStatement{}, err
	}
	properties = append(properties, typeStatements.Properties...)

	name := attr.TFSchemaName
	propsStmts := GroupCodeStatements(properties, func(properties []string) string {
		return strings.Join(properties, ",\n") + ","
	})
	code := fmt.Sprintf(`"%s": %s{
		%s
	}`, name, typeStatements.AttrType, propsStmts.Code)
	return CodeStatement{
		Code:    code,
		Imports: propsStmts.Imports,
	}, nil
}

func commonProperties(attr *codespec.Attribute, typeStatements *TypeSpecificStatements) ([]CodeStatement, error) {
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

	planModStmt, err := planModifierStatement(attr, typeStatements.PlanModifierType)
	if err != nil {
		return nil, err
	}
	if planModStmt != nil {
		result = append(result, *planModStmt)
	}

	if validatorStmt := validatorStatement(typeStatements.Validators, typeStatements.ValidatorType); validatorStmt != nil {
		result = append(result, *validatorStmt)
	}

	return result, nil
}

func planModifierStatement(attr *codespec.Attribute, planModifierType string) (*CodeStatement, error) {
	const (
		importCustomPlanModifier = "github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
		importPlanModifier       = "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
		importStringPlanModifier = "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	)
	var modifiers []string
	imports := make(map[string]struct{})

	if attr.CreateOnly {
		if attr.Bool != nil && attr.Bool.Default != nil {
			modifiers = append(modifiers, fmt.Sprintf("customplanmodifier.CreateOnlyBoolWithDefault(%t)", *attr.Bool.Default))
		} else {
			modifiers = append(modifiers, "customplanmodifier.CreateOnly()")
		}
		imports[importCustomPlanModifier] = struct{}{}
	}
	if attr.RequestOnlyRequiredOnCreate {
		modifiers = append(modifiers, "customplanmodifier.RequestOnlyRequiredOnCreate()")
		imports[importCustomPlanModifier] = struct{}{}
	}
	if attr.ImmutableComputed {
		if attr.String == nil {
			return nil, fmt.Errorf("immutableComputed is only supported for string attributes, found non-string type for attribute '%s'", attr.TFSchemaName)
		}
		modifiers = append(modifiers, "stringplanmodifier.UseStateForUnknown()")
		imports[importStringPlanModifier] = struct{}{}
	}

	if len(modifiers) == 0 {
		return nil, nil
	}

	imports[importPlanModifier] = struct{}{}
	var importList []string
	for imp := range imports {
		importList = append(importList, imp)
	}
	return &CodeStatement{
		Code:    fmt.Sprintf("PlanModifiers: []%s{%s}", planModifierType, strings.Join(modifiers, ", ")),
		Imports: importList,
	}, nil
}

func validatorStatement(validators []CodeStatement, validatorType string) *CodeStatement {
	if len(validators) == 0 {
		return nil
	}
	const importSchemaValidator = "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	imports := map[string]struct{}{
		importSchemaValidator: {},
	}
	var entries []string
	for _, v := range validators {
		entries = append(entries, v.Code)
		for _, imp := range v.Imports {
			imports[imp] = struct{}{}
		}
	}
	var importList []string
	for imp := range imports {
		importList = append(importList, imp)
	}
	return &CodeStatement{
		Code:    fmt.Sprintf("Validators: []%s{\n\t\t%s,\n\t}", validatorType, strings.Join(entries, ",\n\t\t")),
		Imports: importList,
	}
}
