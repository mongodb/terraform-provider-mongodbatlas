package resource

import (
	"go/format"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/codetemplate"
)

func GenerateGoCode(input *codespec.Resource) string {
	tmplInputs := codetemplate.ResourceFileInputs{
		PackageName:  input.Name.LowerCaseNoUnderscore(),
		ResourceName: input.Name.SnakeCase(),
		APIOperations: codetemplate.APIOperations{
			VersionHeader: "application/vnd.atlas.2023-01-01+json",
			Create: codetemplate.Operation{
				Path: "/api/atlas/v2/groups/{groupId}/pushBasedLogExport",
				PathParams: []codetemplate.Param{
					{
						PascalCaseName: "GroupId",
						CamelCaseName:  "groupId",
					},
				},
			},
			Update: codetemplate.Operation{
				Path: "/api/atlas/v2/groups/{groupId}/pushBasedLogExport",
				PathParams: []codetemplate.Param{
					{
						PascalCaseName: "GroupId",
						CamelCaseName:  "groupId",
					},
				},
			},
			Read: codetemplate.Operation{
				Path: "/api/atlas/v2/groups/{groupId}/pushBasedLogExport",
				PathParams: []codetemplate.Param{
					{
						PascalCaseName: "GroupId",
						CamelCaseName:  "groupId",
					},
				},
			},
			Delete: codetemplate.Operation{
				Path: "/api/atlas/v2/groups/{groupId}/pushBasedLogExport",
				PathParams: []codetemplate.Param{
					{
						PascalCaseName: "GroupId",
						CamelCaseName:  "groupId",
					},
				},
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
