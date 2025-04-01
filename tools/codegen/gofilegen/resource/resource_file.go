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
			Create: codetemplate.Operation{
				Path:       input.Operations.CreatePath,
				PathParams: obtainPathParams(input.Operations.CreatePath),
			},
			Update: codetemplate.Operation{
				Path:       input.Operations.UpdatePath,
				PathParams: obtainPathParams(input.Operations.UpdatePath),
			},
			Read: codetemplate.Operation{
				Path:       input.Operations.ReadPath,
				PathParams: obtainPathParams(input.Operations.ReadPath),
			},
			Delete: codetemplate.Operation{
				Path:       input.Operations.DeletePath,
				PathParams: obtainPathParams(input.Operations.DeletePath),
			},
		},
	}
	result := codetemplate.ApplyResourceFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		panic(err)
	}
	return string(formattedResult)
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
