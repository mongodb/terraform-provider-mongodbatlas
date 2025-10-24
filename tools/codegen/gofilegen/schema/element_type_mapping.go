package schema

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

const typesImportStatement = "github.com/hashicorp/terraform-plugin-framework/types"

func ElementTypeProperty(elementType codespec.ElemType) CodeStatement {
	result := codespec.ElementTypeToSchemaString[elementType]
	return CodeStatement{
		Code:    fmt.Sprintf("ElementType: %s", result),
		Imports: []string{typesImportStatement},
	}
}
