package autogen

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
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
	for i := range valModel.NumField() {
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

	switch v := value.(type) {
	case customtypes.ObjectValueInterface:
		if v.IsUnknown() {
			return v.NewObjectValueNull(ctx), nil
		}

		valuePtr, diags := v.ValuePtrAsAny(ctx)
		if diags.HasError() {
			return nil, fmt.Errorf("unmarshal failed to convert object: %v", diags)
		}

		err := ResolveUnknowns(valuePtr)
		if err != nil {
			return nil, err
		}

		objNew := v.NewObjectValue(ctx, valuePtr)
		return objNew, nil
	case customtypes.ListValueInterface:
		if v.IsUnknown() {
			return v.NewListValueNull(ctx), nil
		}
		// If known, no need to process each list item since unmarshal does not generate unknown attributes.
		return v, nil
	case customtypes.NestedListValueInterface:
		if v.IsUnknown() {
			return v.NewNestedListValueNull(ctx), nil
		}

		slicePtr, diags := v.SlicePtrAsAny(ctx)
		if diags.HasError() {
			return nil, fmt.Errorf("unmarshal failed to convert list: %v", diags)
		}

		sliceVal := reflect.ValueOf(slicePtr).Elem()
		for i := range sliceVal.Len() {
			elementPtr := sliceVal.Index(i).Addr().Interface()
			err := ResolveUnknowns(elementPtr)
			if err != nil {
				return nil, err
			}
		}

		return v.NewNestedListValue(ctx, slicePtr), nil
	case customtypes.SetValueInterface:
		if v.IsUnknown() {
			return v.NewSetValueNull(ctx), nil
		}
		// If known, no need to process each set item since unmarshal does not generate unknown attributes.
		return v, nil
	case customtypes.NestedSetValueInterface:
		if v.IsUnknown() {
			return v.NewNestedSetValueNull(ctx), nil
		}
		// If known, no need to process each set item since unmarshal does not generate unknown attributes.
		return v, nil
	case customtypes.MapValueInterface:
		if v.IsUnknown() {
			return v.NewMapValueNull(ctx), nil
		}
		// If known, no need to process each map entry since unmarshal does not generate unknown attributes.
		return v, nil
	case customtypes.NestedMapValueInterface:
		if v.IsUnknown() {
			return v.NewNestedMapValueNull(ctx), nil
		}

		/*
			Resolving nested map unknowns is a pretty expensive operation if the nested struct is big since it requires
			copying every map element (reflection limitation) even when no unknows are present.
			Since nested maps are very rare, it is unlikely that this will become a problem in practice.
			If it does, we could either:
				- Check whether there are any unknowns in the nested object before creating a copy.
				- Resolve unknowns during unmarshal.
			Not worth the complexity as of now.
		*/

		mapPtr, diags := v.MapPtrAsAny(ctx)
		if diags.HasError() {
			return nil, fmt.Errorf("unmarshal failed to convert map: %v", diags)
		}

		mapVal := reflect.ValueOf(mapPtr).Elem()
		mapElemType := mapVal.Type().Elem()

		iter := mapVal.MapRange()
		for iter.Next() {
			elementVal := reflect.New(mapElemType)
			elementVal.Elem().Set(iter.Value())

			err := ResolveUnknowns(elementVal.Interface())
			if err != nil {
				return nil, err
			}

			mapVal.SetMapIndex(iter.Key(), elementVal.Elem())
		}

		return v.NewNestedMapValue(ctx, mapPtr), nil
	}

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
		if _, ok := attrType.(jsontypes.NormalizedType); ok {
			return jsontypes.NewNormalizedNull(), nil
		}
		return nil, fmt.Errorf("unmarshal to get null value not supported yet for type %T", attrType)
	}
}
