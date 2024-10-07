package schema

import (
	"bytes"
	"go/format"
	"text/template"

	genconfigmapper "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type TemplateInputs struct {
	PackageName       string
	AdditionalImports []string
	Attributes        []string
}

func GenerateGoCode(input genconfigmapper.Resource) string {
	const tmpl = `package {{ .PackageName }}

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	{{range .AdditionalImports }}
	"{{ . }}"
	{{- end }}
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			{{ range  .Attributes }}{{ . }}{{- end }}
		},
	}
}
`
	// TODO get attributes

	tmplInputs := TemplateInputs{
		PackageName:       input.Name, // TODO adjust format
		AdditionalImports: []string{},
		Attributes:        RenderAttributes(input.Schema.Attributes),
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

	formattedResult, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	return string(formattedResult[:])
}
