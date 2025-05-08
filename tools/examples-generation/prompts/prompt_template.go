package prompts

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed generatemainhcl.user.md
var generateMainHCLUserPromptTemplate string

//go:embed generatemainhcl.system.md
var GenerateMainHCLSystemPrompt string

//go:embed generatevarsdefhcl.user.md
var generateVarsDefHCLUserPromptTemplate string

//go:embed generatevarsdefhcl.system.md
var GenerateVarsDefHCLSystemPrompt string

//go:embed generatereadme.user.md
var generateReadmeUserPromptTemplate string

//go:embed generatereadme.system.md
var GenerateReadmeSystemPrompt string

type MainHCLUserPromptInputs struct {
	ResourceName           string
	ResourceImplementation string
	ResourceAPISpec        string
}

func GetMainHCLGenerationUserPrompt(inputs MainHCLUserPromptInputs) string {
	t, err := template.New("template").Parse(generateMainHCLUserPromptTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, inputs)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

type VarsDefHCLUserPromptInputs struct {
	HCLConfig string
}

func GetVarsDefHCLGenerationUserPrompt(inputs VarsDefHCLUserPromptInputs) string {
	t, err := template.New("template").Parse(generateVarsDefHCLUserPromptTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, inputs)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

type ReadmeUserPromptInputs struct {
	HCLConfig       string
	VariablesDefHCL string
	ResourceAPISpec string
}

func GetReadmeGenerationUserPrompt(inputs ReadmeUserPromptInputs) string {
	t, err := template.New("template").Parse(generateReadmeUserPromptTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, inputs)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
