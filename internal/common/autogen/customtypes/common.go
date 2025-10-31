package customtypes

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func getValueElementType[T attr.Value](ctx context.Context) attr.Type {
	var t T
	return t.Type(ctx)
}

func getElementType[T any](ctx context.Context) (attr.Type, diag.Diagnostics) {
	var t T
	attrTypes, diags := valueToAttributeTypes(ctx, reflect.ValueOf(t))
	if diags.HasError() {
		return nil, diags
	}

	return basetypes.ObjectType{AttrTypes: attrTypes}, nil
}

func valueToAttributeTypes(ctx context.Context, value reflect.Value) (map[string]attr.Type, diag.Diagnostics) {
	valueType := value.Type()

	if valueType.Kind() != reflect.Struct {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
			"Error getting value attribute types",
			fmt.Sprintf(`%T has usupported type: %s`, value.Interface(), valueType),
		)}
	}

	attributeTypes := make(map[string]attr.Type)
	for i := range valueType.NumField() {
		typeField := valueType.Field(i)
		valueField := value.Field(i)

		tfName := typeField.Tag.Get(`tfsdk`)
		if tfName == "" {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
				"Error getting value attribute types",
				fmt.Sprintf(`%T has no tfsdk tag on field %s`, value.Interface(), typeField.Name),
			)}
		}

		attrValue, ok := valueField.Interface().(attr.Value)
		if !ok {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
				"Error getting value attribute types",
				fmt.Sprintf(`%T has unsupported type in field %s: %T`, value.Interface(), typeField.Name, valueField.Interface()),
			)}
		}

		attributeTypes[tfName] = attrValue.Type(ctx)
	}

	return attributeTypes, nil
}
