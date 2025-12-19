package codetemplate

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed schema-file.go.tmpl
var schemaFileTemplate string

type SchemaFileInputs struct {
	DeprecationMessage *string
	PackageName        string
	SchemaAttributes   string
	Models             string
	Imports            []string
}

//go:embed resource-file.go.tmpl
var resourceFileTemplate string

//go:embed data-source-file.go.tmpl
var dataSourceFileTemplate string

//go:embed data-source-schema-file.go.tmpl
var dataSourceSchemaFileTemplate string

//go:embed plural-data-source-file.go.tmpl
var pluralDataSourceFileTemplate string

//go:embed plural-data-source-schema-file.go.tmpl
var pluralDataSourceSchemaFileTemplate string

type ResourceFileInputs struct {
	PackageName   string
	ResourceName  string
	APIOperations APIOperations
	MoveState     *MoveState
	IDAttributes  []string // e.g. ["project_id", "name"]
}

type APIOperations struct {
	Delete        *Operation
	Update        *Operation
	VersionHeader string
	Create        Operation
	Read          Operation
}

type Operation struct {
	Wait              *Wait
	Path              string
	HTTPMethod        string
	StaticRequestBody string
	PathParams        []Param
}

type Wait struct {
	StateProperty     string
	PendingStates     []string
	TargetStates      []string
	TimeoutSeconds    int
	MinTimeoutSeconds int
	DelaySeconds      int
}
type Param struct {
	PascalCaseName string
	CamelCaseName  string
}

type MoveState struct {
	SourceResources []string
}

type DataSourceFileInputs struct {
	PackageName    string
	DataSourceName string
	VersionHeader  string
	ReadPath       string
	ReadMethod     string
	PathParams     []Param
}

type DataSourceSchemaFileInputs struct {
	DeprecationMessage *string
	PackageName        string
	SchemaAttributes   string
	DSModel            string
	Imports            []string
}

type PluralDataSourceFileInputs struct {
	PackageName    string
	DataSourceName string
	VersionHeader  string
	ReadPath       string
	ReadMethod     string
	PathParams     []Param
	QueryParams    []Param
}

type PluralDataSourceSchemaFileInputs struct {
	DeprecationMessage *string
	PackageName        string
	SchemaAttributes   string
	PluralDSModel      string
	Imports            []string
}

func ApplySchemaFileTemplate(inputs *SchemaFileInputs) bytes.Buffer {
	return applyTemplate(schemaFileTemplate, inputs)
}

func ApplyResourceFileTemplate(inputs *ResourceFileInputs) bytes.Buffer {
	return applyTemplate(resourceFileTemplate, inputs)
}

func ApplyDataSourceFileTemplate(inputs *DataSourceFileInputs) bytes.Buffer {
	return applyTemplate(dataSourceFileTemplate, inputs)
}

func ApplyDataSourceSchemaFileTemplate(inputs *DataSourceSchemaFileInputs) bytes.Buffer {
	return applyTemplate(dataSourceSchemaFileTemplate, inputs)
}

func ApplyPluralDataSourceFileTemplate(inputs *PluralDataSourceFileInputs) bytes.Buffer {
	return applyTemplate(pluralDataSourceFileTemplate, inputs)
}

func ApplyPluralDataSourceSchemaFileTemplate(inputs *PluralDataSourceSchemaFileInputs) bytes.Buffer {
	return applyTemplate(pluralDataSourceSchemaFileTemplate, inputs)
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
