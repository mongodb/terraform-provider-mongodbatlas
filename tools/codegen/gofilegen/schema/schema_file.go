package schema

import (
	"fmt"
	"go/format"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/codetemplate"
)

func GenerateGoCode(input *codespec.Resource) ([]byte, error) {
	schemaAttrs := GenerateSchemaAttributes(input.Schema.Attributes)
	models := GenerateTypedModels(input.Schema.Attributes)

	imports := []string{"github.com/hashicorp/terraform-plugin-framework/resource/schema"}
	imports = append(imports, schemaAttrs.Imports...)
	imports = append(imports, models.Imports...)

	// Generate DS model if data sources are defined
	var dsModel CodeStatement
	if input.DataSources != nil && input.DataSources.Schema != nil {
		dsModel = GenerateDataSourceTypedModels(input.DataSources.Schema.Attributes)
		imports = append(imports, dsModel.Imports...)
	}

	tmplInputs := codetemplate.SchemaFileInputs{
		PackageName:        input.PackageName,
		Imports:            imports,
		SchemaAttributes:   schemaAttrs.Code,
		Models:             models.Code,
		DSModel:            dsModel.Code,
		DeprecationMessage: input.Schema.DeprecationMessage,
	}
	result := codetemplate.ApplySchemaFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code (schema): %w", err)
	}
	return formattedResult, nil
}
