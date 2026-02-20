package schema

import (
	"fmt"
	"go/format"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/codetemplate"
)

func GenerateGoCode(input *codespec.Resource) ([]byte, error) {
	allSchemaAttrs := append(append(codespec.Attributes{}, input.Schema.Attributes...), input.Schema.CraftedAttributes...)

	schemaAttrs, err := GenerateSchemaAttributes(allSchemaAttrs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate schema attributes: %w", err)
	}
	models := GenerateTypedModels(allSchemaAttrs)

	imports := []string{"github.com/hashicorp/terraform-plugin-framework/resource/schema"}
	imports = append(imports, schemaAttrs.Imports...)
	imports = append(imports, models.Imports...)

	tmplInputs := codetemplate.SchemaFileInputs{
		PackageName:        input.PackageName,
		Imports:            imports,
		SchemaAttributes:   schemaAttrs.Code,
		Models:             models.Code,
		DeprecationMessage: input.Schema.DeprecationMessage,
	}
	result := codetemplate.ApplySchemaFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code (schema): %w", err)
	}
	return formattedResult, nil
}
