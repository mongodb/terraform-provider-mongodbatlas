package datasource

import (
	"fmt"
	"go/format"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/codetemplate"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/resource"
)

func GenerateGoCode(input *codespec.Resource) ([]byte, error) {
	if input.DataSources == nil || input.DataSources.Operations.Read == nil {
		return nil, fmt.Errorf("data source read operation is required for %s", input.Name)
	}

	readOp := input.DataSources.Operations.Read
	pathParams := resource.GetPathParams(readOp.Path)

	tmplInputs := codetemplate.DataSourceFileInputs{
		PackageName:    input.PackageName,
		DataSourceName: input.Name,
		VersionHeader:  input.DataSources.Operations.VersionHeader,
		ReadPath:       readOp.Path,
		ReadMethod:     readOp.HTTPMethod,
		PathParams:     pathParams,
	}
	result := codetemplate.ApplyDataSourceFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code (data source): %w", err)
	}
	return formattedResult, nil
}

// GeneratePluralGoCode generates the plural_data_source.go file for list-based data sources
func GeneratePluralGoCode(input *codespec.Resource) ([]byte, error) {
	if input.DataSources == nil || input.DataSources.Operations.List == nil {
		return nil, fmt.Errorf("data source list operation is required for plural data source %s", input.Name)
	}

	listOp := input.DataSources.Operations.List
	pathParams := resource.GetPathParams(listOp.Path)
	queryParams := getQueryParams(*input.DataSources.Schema.PluralDSAttributes)

	tmplInputs := codetemplate.PluralDataSourceFileInputs{
		PackageName:    input.PackageName,
		DataSourceName: input.Name + "_list", // pluralize name
		VersionHeader:  input.DataSources.Operations.VersionHeader,
		ReadPath:       listOp.Path,
		ReadMethod:     listOp.HTTPMethod,
		PathParams:     pathParams,
		QueryParams:    queryParams,
	}
	result := codetemplate.ApplyPluralDataSourceFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code (plural data source): %w", err)
	}
	return formattedResult, nil
}

// getQueryParams extracts query parameters from plural data source attributes.
// Query parameters are optional top-level attributes
func getQueryParams(attributes codespec.Attributes) []codetemplate.Param {
	var queryParams []codetemplate.Param

	for i := range attributes {
		// Only consider optional attributes as query parameters
		if attributes[i].ComputedOptionalRequired == codespec.Optional {
			attrType := getAttributeType(&attributes[i])
			param := codetemplate.Param{
				PascalCaseName: stringcase.Capitalize(attributes[i].TFModelName),
				CamelCaseName:  stringcase.Uncapitalize(attributes[i].TFModelName),
				Type:           attrType,
			}
			queryParams = append(queryParams, param)
		}
	}

	return queryParams
}

// getAttributeType returns the type of an attribute for query parameter handling
func getAttributeType(attr *codespec.Attribute) string {
	switch {
	case attr.String != nil:
		return "string"
	case attr.Int64 != nil:
		return "int64"
	case attr.Bool != nil:
		return "bool"
	case attr.List != nil:
		return "list"
	case attr.Set != nil:
		return "set"
	default:
		return "unknown"
	}
}
