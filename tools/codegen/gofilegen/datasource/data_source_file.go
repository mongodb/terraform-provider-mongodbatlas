package datasource

import (
	"fmt"
	"go/format"

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
