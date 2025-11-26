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
	var attributes codespec.Attributes
	if input.Schema != nil {
		attributes = input.Schema.Attributes
	}

	tmplInputs := codetemplate.ResourceFileInputs{
		PackageName:  input.PackageName,
		ResourceName: input.Name,
		APIOperations: codetemplate.APIOperations{
			VersionHeader: input.Operations.VersionHeader,
			Create:        *toCodeTemplateOpModel(&input.Operations.Create),
			Update:        toCodeTemplateOpModel(input.Operations.Update),
			Read:          *toCodeTemplateOpModel(&input.Operations.Read),
			Delete:        toCodeTemplateOpModel(input.Operations.Delete),
		},
		MoveState:    toCodeTemplateMoveStateModel(input.MoveState),
		IDAttributes: getIDAttributes(input.Operations.Read.Path, attributes),
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
		PathParams:        getPathParams(op.Path),
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
			PascalCaseName: stringcase.Capitalize(paramName),
		})
	}
	return params
}

func getIDAttributes(readPath string, attributes codespec.Attributes) []string {
	params := getPathParams(readPath)
	result := make([]string, len(params))

	for i, param := range params {
		snakeCaseName := stringcase.ToSnakeCase(param.CamelCaseName)

		// Find the matching attribute in the schema (which has aliases already applied)
		// Path params are marked as Required with OmitAlways req body usage
		found := false
		for j := range attributes {
			attr := &attributes[j]
			if attr.ReqBodyUsage == codespec.OmitAlways && attr.ComputedOptionalRequired == codespec.Required {
				// Check if this attribute's TFModelName matches the param name
				if stringcase.ToSnakeCase(attr.TFModelName) == snakeCaseName {
					result[i] = attr.TFSchemaName // Use the aliased TF schema name
					found = true
					break
				}
			}
		}
		if !found {
			result[i] = snakeCaseName // Fallback to snake_case if not found
		}
	}
	return result
}
