package resource

import (
	"fmt"
	"go/format"
	"regexp"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/codetemplate"
)

func GenerateGoCode(input *codespec.Resource) ([]byte, error) {
	// Resources require Create and Read operations - fail fast if missing
	if input.Operations.Create == nil {
		return nil, fmt.Errorf("resource %s is missing required Create operation", input.Name)
	}
	if input.Operations.Read == nil {
		return nil, fmt.Errorf("resource %s is missing required Read operation", input.Name)
	}

	idAttributes := GetIDAttributes(input.Operations.Read.Path)

	tmplInputs := codetemplate.ResourceFileInputs{
		PackageName:  input.PackageName,
		ResourceName: input.Name,
		APIOperations: codetemplate.APIOperations{
			VersionHeader: input.Operations.VersionHeader,
			Create:        *toCodeTemplateOpModel(input.Operations.Create),
			Update:        toCodeTemplateOpModel(input.Operations.Update),
			Read:          *toCodeTemplateOpModel(input.Operations.Read),
			Delete:        toCodeTemplateOpModel(input.Operations.Delete),
		},
		MoveState:    toCodeTemplateMoveStateModel(input.MoveState),
		IDAttributes: idAttributes,
	}
	result := codetemplate.ApplyResourceFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code (resource): %w", err)
	}
	return formattedResult, nil
}

func toCodeTemplateMoveStateModel(moveState *codespec.MoveState) *codetemplate.MoveState {
	if moveState == nil {
		return nil
	}
	return &codetemplate.MoveState{SourceResources: moveState.SourceResources}
}

func toCodeTemplateOpModel(op *codespec.APIOperation) *codetemplate.Operation {
	if op == nil {
		return nil
	}
	return &codetemplate.Operation{
		Path:              op.Path,
		HTTPMethod:        op.HTTPMethod,
		PathParams:        GetPathParams(op.Path),
		Wait:              getWaitValues(op.Wait),
		StaticRequestBody: op.StaticRequestBody,
	}
}

func getWaitValues(wait *codespec.Wait) *codetemplate.Wait {
	if wait == nil {
		return nil
	}
	return &codetemplate.Wait{
		StateProperty:     wait.StateProperty,
		PendingStates:     wait.PendingStates,
		TargetStates:      wait.TargetStates,
		TimeoutSeconds:    wait.TimeoutSeconds,
		MinTimeoutSeconds: wait.MinTimeoutSeconds,
		DelaySeconds:      wait.DelaySeconds,
	}
}

// GetPathParams extracts path parameters from a URL path and returns them as Param structs.
// This can eventually be explicitly defined in the intermediate model if additional information is required.
func GetPathParams(s string) []codetemplate.Param {
	params := []codetemplate.Param{}

	// Use regex to find all {paramName} patterns
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(s, -1)

	for _, match := range matches {
		paramName := match[1]
		params = append(params, codetemplate.Param{
			CamelCaseName:  paramName,
			PascalCaseName: stringcase.Capitalize(paramName),
		})
	}
	return params
}

// GetIDAttributes converts path params to snake_case attribute names.
// Used for both resource ID attributes and data source required fields.
func GetIDAttributes(readPath string) []string {
	params := GetPathParams(readPath)
	result := make([]string, len(params))
	for i, param := range params {
		result[i] = stringcase.ToSnakeCase(param.PascalCaseName)
	}
	return result
}
