package schema

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

var elementTypeToString = map[codespec.ElemType]string{
	codespec.Bool:           "types.BoolType",
	codespec.Float64:        "types.Float64Type",
	codespec.Int64:          "types.Int64Type",
	codespec.Number:         "types.NumberType",
	codespec.String:         "types.StringType",
	codespec.CustomTypeJSON: codespec.CustomTypeJSONVar.Schema,
}

const typesImportStatement = "github.com/hashicorp/terraform-plugin-framework/types"

func ElementTypeProperty(elementType codespec.ElemType) CodeStatement {
	result := elementTypeToString[elementType]
	return CodeStatement{
		Code:    fmt.Sprintf("ElementType: %s", result),
		Imports: []string{typesImportStatement},
	}
}
