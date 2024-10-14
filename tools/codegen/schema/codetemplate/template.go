package codetemplate

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed schema-file.go.tmpl
var schemaFileTemplate string

type SchemaFileInputs struct {
	PackageName      string
	Imports          []string
	SchemaAttributes []string
}

func ApplySchemaFileTemplate(inputs SchemaFileInputs) bytes.Buffer {
	t, err := template.New("template").Parse(schemaFileTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, inputs)
	if err != nil {
		panic(err)
	}

	return buf
}
