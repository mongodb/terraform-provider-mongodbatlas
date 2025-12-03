package datasource

import (
	"fmt"
	"go/format"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/codetemplate"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/resource"
)

func GenerateGoCode(input *codespec.DataSource) ([]byte, error) {
	pathParams := resource.GetPathParams(input.ReadOperation.Path)
	requiredFields := resource.GetIDAttributes(input.ReadOperation.Path)

	tmplInputs := codetemplate.DataSourceFileInputs{
		PackageName:    input.PackageName,
		DataSourceName: input.Name,
		VersionHeader:  input.VersionHeader,
		ReadPath:       input.ReadOperation.Path,
		ReadMethod:     input.ReadOperation.HTTPMethod,
		RequiredFields: requiredFields,
		PathParams:     pathParams,
	}
	result := codetemplate.ApplyDataSourceFileTemplate(&tmplInputs)

	formattedResult, err := format.Source(result.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated Go code (data source): %w", err)
	}
	return formattedResult, nil
}
