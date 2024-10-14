package schema

import (
	"go/format"

	genconfigmapper "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/schema/codetemplate"
)

func GenerateGoCode(input genconfigmapper.Resource) string {
	schemaAttrs := GenerateSchemaAttributes(input.Schema.Attributes)
	attrsCode := []string{}
	imports := []string{}
	for _, attr := range schemaAttrs {
		attrsCode = append(attrsCode, attr.Result)
		imports = append(imports, attr.Imports...)
	}

	tmplInputs := codetemplate.SchemaFileInputs{
		PackageName:      input.Name,
		Imports:          imports,
		SchemaAttributes: attrsCode,
	}
	result := codetemplate.ApplySchemaFileTemplate(tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		panic(err)
	}
	return string(formattedResult)
}
