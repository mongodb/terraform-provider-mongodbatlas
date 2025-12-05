package schema

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/codetemplate"
)

// GenerateDataSourceSchemaGoCode generates the data_source_schema.go file containing
// DataSourceSchema() and TFDSModel for data sources.
func GenerateDataSourceSchemaGoCode(input *codespec.Resource) ([]byte, error) {
	if input.DataSources == nil || input.DataSources.Schema == nil {
		return nil, fmt.Errorf("data source schema is required for %s", input.Name)
	}

	dsSchema := input.DataSources.Schema
	schemaAttrs := GenerateDataSourceSchemaAttributes(dsSchema.Attributes)
	dsModel := GenerateDataSourceTypedModels(dsSchema.Attributes)

	// Collect imports (dsschema is hardcoded in the template)
	var imports []string
	imports = append(imports, schemaAttrs.Imports...)
	imports = append(imports, dsModel.Imports...)

	tmplInputs := codetemplate.DataSourceSchemaFileInputs{
		PackageName:        input.PackageName,
		Imports:            imports,
		SchemaAttributes:   schemaAttrs.Code,
		DSModel:            dsModel.Code,
		DeprecationMessage: dsSchema.DeprecationMessage,
	}
	result := codetemplate.ApplyDataSourceSchemaFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code (data source schema): %w", err)
	}
	return formattedResult, nil
}

// GenerateDataSourceSchemaAttributes generates schema attributes for data source schema.
// Data source attributes use dsschema types instead of resource schema types.
func GenerateDataSourceSchemaAttributes(attrs codespec.Attributes) CodeStatement {
	attrsCode := []string{}
	imports := []string{}
	for i := range attrs {
		result := dataSourceAttrGenerator(&attrs[i]).AttributeCode()
		attrsCode = append(attrsCode, result.Code)
		imports = append(imports, result.Imports...)
	}
	finalAttrs := strings.Join(attrsCode, ",\n") + ","
	return CodeStatement{
		Code:    finalAttrs,
		Imports: imports,
	}
}

// dataSourceAttrGenerator returns the appropriate generator for data source schema attributes.
// Data source schemas use dsschema types (e.g., dsschema.StringAttribute) instead of schema types.
func dataSourceAttrGenerator(attr *codespec.Attribute) attributeGenerator {
	// Reuse the same generators but they will be wrapped to use dsschema types
	return &dsAttrGeneratorWrapper{inner: generator(attr), attr: attr}
}

// dsAttrGeneratorWrapper wraps resource attribute generators to produce data source schema code.
// It replaces "schema." with "dsschema." in the generated code.
type dsAttrGeneratorWrapper struct {
	inner attributeGenerator
	attr  *codespec.Attribute
}

func (g *dsAttrGeneratorWrapper) AttributeCode() CodeStatement {
	result := g.inner.AttributeCode()
	// Replace schema. with dsschema. for data source schemas
	result.Code = strings.ReplaceAll(result.Code, "schema.", "dsschema.")
	// Filter out resource-specific imports and plan modifiers (data sources don't need them)
	var filteredImports []string
	for _, imp := range result.Imports {
		// Skip resource schema import (we use dsschema instead)
		if imp == "github.com/hashicorp/terraform-plugin-framework/resource/schema" {
			continue
		}
		// Skip plan modifier imports (data sources don't use plan modifiers)
		if strings.Contains(imp, "planmodifier") || strings.Contains(imp, "customplanmodifier") {
			continue
		}
		filteredImports = append(filteredImports, imp)
	}
	result.Imports = filteredImports
	return result
}
