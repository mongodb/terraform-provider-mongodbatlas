package resource

import (
	"go/format"
	"regexp"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/codetemplate"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/stringcase"
)

func GenerateGoCode(input *codespec.Resource) string {
	tmplInputs := codetemplate.ResourceFileInputs{
		PackageName:  input.Name.LowerCaseNoUnderscore(),
		ResourceName: input.Name.SnakeCase(),
		APIOperations: codetemplate.APIOperations{
			VersionHeader: input.Operations.VersionHeader,
			Create:        toCodeTemplateOpModel(input.Operations.Create),
			Update:        toCodeTemplateOpModel(input.Operations.Update),
			Read:          toCodeTemplateOpModel(input.Operations.Read),
			Delete:        toCodeTemplateOpModel(input.Operations.Delete),
		},
	}
	result := codetemplate.ApplyResourceFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		panic(err)
	}
	return string(formattedResult)
}

func toCodeTemplateOpModel(op codespec.APIOperation) codetemplate.Operation {
	return codetemplate.Operation{
		Path:       op.Path,
		HTTPMethod: op.HTTPMethod,
		PathParams: obtainPathParams(op.Path),
	}
}

// obtains path parameters for URL, this can evetually be explicitly defined in the intermediate model if additional information is required
func obtainPathParams(s string) []codetemplate.Param {
	params := []codetemplate.Param{}

	// Use regex to find all {paramName} patterns
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(s, -1)

	for _, match := range matches {
		paramName := match[1]
		params = append(params, codetemplate.Param{
			CamelCaseName:  paramName,
			PascalCaseName: stringcase.FromCamelCase(paramName).PascalCase(),
		})
	}

	return params
}
