package schema

import (
	"fmt"
	"go/format"

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
	schemaAttrs := GenerateDataSourceSchemaAttributes(*dsSchema.SingularDSAttributes)
	dsModel := GenerateDataSourceTypedModels(*dsSchema.SingularDSAttributes)

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

// GeneratePluralDataSourceSchemaGoCode generates the plural_data_source_schema.go file containing
// PluralDataSourceSchema() and TFPluralDSModel for plural data sources.
func GeneratePluralDataSourceSchemaGoCode(input *codespec.Resource) ([]byte, error) {
	if input.DataSources == nil || input.DataSources.Schema == nil {
		return nil, fmt.Errorf("data source schema is required for %s", input.Name)
	}

	dsSchema := input.DataSources.Schema
	if dsSchema.PluralDSAttributes == nil {
		return nil, fmt.Errorf("plural data source attributes are required for %s", input.Name)
	}

	schemaAttrs := GenerateDataSourceSchemaAttributes(*dsSchema.PluralDSAttributes)

	// Generate TFPluralDSModel and nested TFResultsModel using the reusable function
	pluralDSModel := GenerateDataSourceTypedModelsWithName(*dsSchema.PluralDSAttributes, "PluralDS")

	// Collect imports (dsschema is hardcoded in the template)
	var imports []string
	imports = append(imports, schemaAttrs.Imports...)
	imports = append(imports, pluralDSModel.Imports...)

	tmplInputs := codetemplate.PluralDataSourceSchemaFileInputs{
		PackageName:        input.PackageName,
		Imports:            imports,
		SchemaAttributes:   schemaAttrs.Code,
		PluralDSModel:      pluralDSModel.Code,
		DeprecationMessage: dsSchema.DeprecationMessage,
	}
	result := codetemplate.ApplyPluralDataSourceSchemaFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code (plural data source schema): %w", err)
	}
	return formattedResult, nil
}
