package schema

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

const (
	typesImportStatement     = "github.com/hashicorp/terraform-plugin-framework/types"
	jsontypesImportStatement = "github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

func ElementTypeProperty(elementType codespec.ElemType) CodeStatement {
	result := codespec.ElementTypeToSchemaString[elementType]
	imports := []string{typesImportStatement}
	if elementType == codespec.CustomTypeJSON {
		imports = []string{jsontypesImportStatement}
	}
	return CodeStatement{
		Code:    fmt.Sprintf("ElementType: %s", result),
		Imports: imports,
	}
}

func collectionElemType(attr *codespec.Attribute) (codespec.ElemType, bool) {
	switch {
	case attr.List != nil:
		return attr.List.ElementType, true
	case attr.Map != nil:
		return attr.Map.ElementType, true
	case attr.Set != nil:
		return attr.Set.ElementType, true
	default:
		return codespec.Unknown, false
	}
}
