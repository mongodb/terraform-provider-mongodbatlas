package schema

import (
	"go/format"

	genconfigmapper "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/schema/codetemplate"
)

func GenerateGoCode(input genconfigmapper.Resource) string {
	schemaAttrs := GenerateSchemaAttributes(input.Schema.Attributes)

	tmplInputs := codetemplate.SchemaFileInputs{
		PackageName:      input.Name,
		Imports:          schemaAttrs.Imports,
		SchemaAttributes: schemaAttrs.Result,
	}
	result := codetemplate.ApplySchemaFileTemplate(tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		panic(err)
	}
	return string(formattedResult)
}
