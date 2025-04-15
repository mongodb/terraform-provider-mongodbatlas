package autogen

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResolveUnknowns changes unknown attributes to null so they can be used in the response state.
func ResolveUnknowns(model any) error {
	valModel := reflect.ValueOf(model)
	if valModel.Kind() != reflect.Ptr {
		panic("model must be pointer")
	}
	valModel = valModel.Elem()
	if valModel.Kind() != reflect.Struct {
		panic("model must be pointer to struct")
	}
	ctx := context.Background()
	for i := 0; i < valModel.NumField(); i++ {
		field := valModel.Field(i)
		value, ok := field.Interface().(attr.Value)
		if ok && value.IsUnknown() && field.CanSet() {
			nullVal, err := getNullAttr(value.Type(ctx))
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(nullVal))
		}
	}
	return nil
}

func getNullAttr(attrType attr.Type) (attr.Value, error) {
	switch attrType {
	case types.StringType:
		return types.StringNull(), nil
	case types.BoolType:
		return types.BoolNull(), nil
	case types.Int64Type:
		return types.Int64Null(), nil
	case types.Float64Type:
		return types.Float64Null(), nil
	default:
		if objType, ok := attrType.(types.ObjectType); ok {
			return types.ObjectNull(objType.AttributeTypes()), nil
		}
		if listType, ok := attrType.(types.ListType); ok {
			return types.ListNull(listType.ElemType), nil
		}
		if setType, ok := attrType.(types.SetType); ok {
			return types.SetNull(setType.ElemType), nil
		}
		return nil, fmt.Errorf("unmarshal to get null value not supported yet for type %T", attrType)
	}
}
