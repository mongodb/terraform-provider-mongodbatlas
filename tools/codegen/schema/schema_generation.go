package schema

import (
	"bytes"
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
		{{ range  .Attributes }}
		{{ . }},
		{{- end }}
	}
}
`
	// TODO get attributes

	tmplInputs := TemplateInputs{
		PackageName:       input.Name, // TODO adjust format
		AdditionalImports: []string{},
		Attributes: []string{`"str_attr": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the search deployment.",
			}`, `"other": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the search deployment.",
			}`},
	}
	// Parse the template
	t, err := template.New("code").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	// Execute the template with the input data
	var buf bytes.Buffer
	err = t.Execute(&buf, tmplInputs)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
