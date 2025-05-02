package resource

import (
	"go/format"
	"regexp"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
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
			Create:        toCodeTemplateOpModel(input.Operations.Create, false),
			Update:        toCodeTemplateOpModel(input.Operations.Update, false),
			Read:          toCodeTemplateOpModel(input.Operations.Read, false),
			Delete:        toCodeTemplateOpModel(input.Operations.Delete, true),
		},
		ImportIDAttributes: getIDAttributes(input.Operations.Read.Path),
	}
	result := codetemplate.ApplyResourceFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		panic(err)
	}
	return string(formattedResult)
}

func toCodeTemplateOpModel(op codespec.APIOperation, isDelete bool) codetemplate.Operation {
	return codetemplate.Operation{
		Path:       op.Path,
		HTTPMethod: op.HTTPMethod,
		PathParams: getPathParams(op.Path),
		Wait:       getWaitValues(op.Wait, isDelete),
	}
}

func getWaitValues(wait *codespec.Wait, isDelete bool) *codetemplate.Wait {
	if wait == nil {
		return nil
	}
	targetStates := wait.TargetStates

	if isDelete {
		// DELETED is a special state value returned in refresh function when API returns 404 or empty object
		targetStates = append(targetStates, retrystrategy.RetryStrategyDeletedState)
	}
	return &codetemplate.Wait{
		StateProperty:     wait.StateProperty,
		PendingStates:     wait.PendingStates,
		TargetStates:      targetStates,
		TimeoutSeconds:    wait.TimeoutSeconds,
		MinTimeoutSeconds: wait.MinTimeoutSeconds,
		DelaySeconds:      wait.DelaySeconds,
	}
}

// obtains path parameters for URL, this can evetually be explicitly defined in the intermediate model if additional information is required
func getPathParams(s string) []codetemplate.Param {
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

func getIDAttributes(readPath string) []string {
	params := getPathParams(readPath)
	result := make([]string, len(params))
	for i, param := range params {
		result[i] = stringcase.FromCamelCase(param.CamelCaseName).SnakeCase()
	}
	return result
}
