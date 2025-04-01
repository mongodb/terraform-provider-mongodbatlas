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
	SchemaAttributes string
	Models           string
	Imports          []string
}

//go:embed resource-file.go.tmpl
var resourceFileTemplate string

type ResourceFileInputs struct {
	PackageName  string
	ResourceName string
}

func ApplySchemaFileTemplate(inputs SchemaFileInputs) bytes.Buffer {
	return applyTemplate(schemaFileTemplate, inputs)
}

func ApplyResourceFileTemplate(inputs ResourceFileInputs) bytes.Buffer {
	return applyTemplate(resourceFileTemplate, inputs)
}

func applyTemplate[T any](templateStr string, inputs T) bytes.Buffer {
	t, err := template.New("template").Parse(templateStr)
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
