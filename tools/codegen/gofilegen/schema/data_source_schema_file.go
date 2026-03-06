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
	if input.DataSources == nil || input.DataSources.Singular == nil {
		return nil, fmt.Errorf("singular data source schema is required for %s", input.Name)
	}

	singular := input.DataSources.Singular
	schemaAttrs, err := GenerateDataSourceSchemaAttributes(singular.Attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate data source schema attributes: %w", err)
	}
	dsModel := GenerateDataSourceTypedModels(singular.Attributes, false)

	// Collect imports (dsschema is hardcoded in the template)
	var imports []string
	imports = append(imports, schemaAttrs.Imports...)
	imports = append(imports, dsModel.Imports...)

	tmplInputs := codetemplate.DataSourceSchemaFileInputs{
		PackageName:        input.PackageName,
		Imports:            imports,
		SchemaAttributes:   schemaAttrs.Code,
		DSModel:            dsModel.Code,
		DeprecationMessage: singular.DeprecationMessage,
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
	if input.DataSources == nil || input.DataSources.Plural == nil {
		return nil, fmt.Errorf("plural data source schema is required for %s", input.Name)
	}

	plural := input.DataSources.Plural
	schemaAttrs, err := GeneratePluralDataSourceSchemaAttributes(plural.Attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plural data source schema attributes: %w", err)
	}

	// Generate TFPluralDSModel and nested TFResultsModel using the reusable function
	pluralDSModel := GenerateDataSourceTypedModels(plural.Attributes, true)

	// Collect imports (dsschema is hardcoded in the template)
	var imports []string
	imports = append(imports, schemaAttrs.Imports...)
	imports = append(imports, pluralDSModel.Imports...)

	tmplInputs := codetemplate.PluralDataSourceSchemaFileInputs{
		PackageName:        input.PackageName,
		Imports:            imports,
		SchemaAttributes:   schemaAttrs.Code,
		PluralDSModel:      pluralDSModel.Code,
		DeprecationMessage: plural.DeprecationMessage,
	}
	result := codetemplate.ApplyPluralDataSourceSchemaFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code (plural data source schema): %w", err)
	}
	return formattedResult, nil
}
