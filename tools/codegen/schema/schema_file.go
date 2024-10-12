package schema

import (
	"bytes"
	"go/format"
	"text/template"

	genconfigmapper "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type TemplateInputs struct {
	PackageName string
	Imports     []string
	Attributes  []string
}

func GenerateGoCode(input genconfigmapper.Resource) string {
	const tmpl = `package {{ .PackageName }}

import (
	"context"
	{{range .Imports }}
	"{{ . }}"
	{{- end }}
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{ {{ range .Attributes }}{{ . }}{{- end }}
		},
	}
}
`
	schemaAttrs := GenerateSchemaAttributes(input.Schema.Attributes)
	attrsCode := []string{}
	imports := []string{}
	for _, attr := range schemaAttrs {
		attrsCode = append(attrsCode, attr.Result)
		imports = append(imports, attr.Imports...)
	}

	tmplInputs := TemplateInputs{
		PackageName: input.Name,
		Imports:     imports,
		Attributes:  attrsCode,
	}
	// Parse the template
	t, err := template.New("schema-template").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	// Execute the template with the input data
	var buf bytes.Buffer
	err = t.Execute(&buf, tmplInputs)
	if err != nil {
		panic(err)
	}

	result := buf.String()

	print(result)

	formattedResult, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	return string(formattedResult)
}
