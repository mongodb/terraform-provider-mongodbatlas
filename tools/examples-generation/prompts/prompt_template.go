package prompts

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed generatehcl.user.md
var generateHCLUserPromptTemplate string

//go:embed generatehcl.system.md
var GenerateHCLSystemPrompt string

type UserPromptTemplateInputs struct {
	ResourceName                  string
	ResourceImplementationSchema  string
	ResourceAPISpecResponseSchema string
}

func GetUserPrompt(inputs UserPromptTemplateInputs) string {
	t, err := template.New("template").Parse(generateHCLUserPromptTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, inputs)
	if err != nil {
		panic(err)
	}

	return string(buf.Bytes())
}
