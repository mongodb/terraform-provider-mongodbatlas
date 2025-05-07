package prompts

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed generatemainhcl.user.md
var generateMainHCLUserPromptTemplate string

//go:embed generatemainhcl.system.md
var GenerateMainHCLSystemPrompt string

//go:embed generatevarsdefhcl.user.md
var generateVarsDefHCLUserPromptTemplate string

//go:embed generatevarsdefhcl.system.md
var GenerateVarsDefHCLSystemPrompt string

type MainHCLUserPromptInputs struct {
	ResourceName                  string
	ResourceImplementationSchema  string
	ResourceAPISpecResponseSchema string
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
