package schema

import (
	"go/format"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/schema/codetemplate"
)

func GenerateGoCode(input codespec.Resource) string {
	schemaAttrs := GenerateSchemaAttributes(input.Schema.Attributes)
	models := GenerateTypedModels(input.Schema.Attributes)

	imports := []string{"github.com/hashicorp/terraform-plugin-framework/resource/schema"}
	imports = append(imports, schemaAttrs.Imports...)
	imports = append(imports, models.Imports...)

	tmplInputs := codetemplate.SchemaFileInputs{
		PackageName:      input.Name.LowerCaseNoUnderscore(),
		Imports:          imports,
		SchemaAttributes: schemaAttrs.Code,
		Models:           models.Code,
	}
	result := codetemplate.ApplySchemaFileTemplate(tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		panic(err)
	}
	return string(formattedResult)
}
