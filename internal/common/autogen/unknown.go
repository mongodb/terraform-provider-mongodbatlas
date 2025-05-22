package autogen

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResolveUnknowns converts unknown attributes to null.
func ResolveUnknowns(model any) error {
	valModel := reflect.ValueOf(model)
	if valModel.Kind() != reflect.Ptr {
		panic("model must be pointer")
	}
	valModel = valModel.Elem()
	if valModel.Kind() != reflect.Struct {
		panic("model must be pointer to struct")
	}
	for i := 0; i < valModel.NumField(); i++ {
		field := valModel.Field(i)
		value, ok := field.Interface().(attr.Value)
		if !ok || !field.CanSet() {
			continue // skip attributes that are not Terraform or not settable
		}
		valNew, err := prepareAttr(value)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(valNew))
	}
	return nil
}

func prepareAttr(value attr.Value) (attr.Value, error) {
	if value.IsNull() { // null values are not converted
		return value, nil
	}
	ctx := context.Background()
	if value.IsUnknown() { // unknown values are converted to null
		return getNullAttr(value.Type(ctx))
	}
	switch v := value.(type) {
	case types.Object:
		mapAttrs := make(map[string]attr.Value)
		for nameAttr, valAttr := range v.Attributes() {
			valNew, err := prepareAttr(valAttr)
			if err != nil {
				return nil, err
			}
			mapAttrs[nameAttr] = valNew
		}
		objNew, diags := types.ObjectValue(v.AttributeTypes(ctx), mapAttrs)
		if diags.HasError() {
			return nil, fmt.Errorf("unmarshal failed to convert object: %v", diags)
		}
		return objNew, nil
	case types.List:
		elems, err := getPreparedCollection(v.Elements())
		if err != nil {
			return nil, err
		}
		listNew, diags := types.ListValue(v.ElementType(ctx), elems)
		if diags.HasError() {
			return nil, fmt.Errorf("unmarshal failed to convert list: %v", diags)
		}
		return listNew, nil
	case types.Set:
		elems, err := getPreparedCollection(v.Elements())
		if err != nil {
			return nil, err
		}
		setNew, diags := types.SetValue(v.ElementType(ctx), elems)
		if diags.HasError() {
			return nil, fmt.Errorf("unmarshal failed to convert set: %v", diags)
		}
		return setNew, nil
	}
	return value, nil
}

func getPreparedCollection(elems []attr.Value) ([]attr.Value, error) {
	arrayAttrs := make([]attr.Value, len(elems))
	for i, elm := range elems {
		valNew, err := prepareAttr(elm)
		if err != nil {
			return nil, err
		}
		arrayAttrs[i] = valNew
	}
	return arrayAttrs, nil
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
		if mapType, ok := attrType.(types.MapType); ok {
			return types.MapNull(mapType.ElemType), nil
		}
		return nil, fmt.Errorf("unmarshal to get null value not supported yet for type %T", attrType)
	}
}
