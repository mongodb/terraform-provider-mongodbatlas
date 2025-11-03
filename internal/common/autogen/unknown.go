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
	if value.IsNull() {
		return value, nil
	}
	ctx := context.Background()

	if value.IsUnknown() {
		switch v := value.(type) {
		case types.String:
			return types.StringNull(), nil
		case types.Bool:
			return types.BoolNull(), nil
		case types.Int64:
			return types.Int64Null(), nil
		case types.Float64:
			return types.Float64Null(), nil
		case jsontypes.Normalized:
			return jsontypes.NewNormalizedNull(), nil
		case customtypes.ObjectValueInterface:
			return v.NewObjectValueNull(ctx), nil
		case customtypes.ListValueInterface:
			return v.NewListValueNull(ctx), nil
		case customtypes.SetValueInterface:
			return v.NewSetValueNull(ctx), nil
		case customtypes.MapValueInterface:
			return v.NewMapValueNull(ctx), nil
		case customtypes.NestedListValueInterface:
			return v.NewNestedListValueNull(ctx), nil
		case customtypes.NestedSetValueInterface:
			return v.NewNestedSetValueNull(ctx), nil
		case customtypes.NestedMapValueInterface:
			return v.NewNestedMapValueNull(ctx), nil
		}
	}

	// Resolve nested types that may contain unknowns after unmarshal
	switch v := value.(type) {
	case customtypes.ObjectValueInterface:
		return resolveObjectAttr(ctx, v)
	case customtypes.NestedListValueInterface:
		return resolveNestedListAttr(ctx, v)
	case customtypes.NestedMapValueInterface:
		return resolveNestedMapAttr(ctx, v)
	}

	return value, nil
}

func resolveObjectAttr(ctx context.Context, obj customtypes.ObjectValueInterface) (attr.Value, error) {
	valuePtr, diags := obj.ValuePtrAsAny(ctx)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert object: %v", diags)
	}

	err := ResolveUnknowns(valuePtr)
	if err != nil {
		return nil, err
	}

	return obj.NewObjectValue(ctx, valuePtr), nil
}

func resolveNestedListAttr(ctx context.Context, list customtypes.NestedListValueInterface) (attr.Value, error) {
	slicePtr, diags := list.SlicePtrAsAny(ctx)
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

	return list.NewNestedListValue(ctx, slicePtr), nil
}

func resolveNestedMapAttr(ctx context.Context, m customtypes.NestedMapValueInterface) (attr.Value, error) {
	/*
		Resolving nested map unknowns is a pretty expensive operation if the nested struct is big since it requires
		copying every map element (reflection limitation) even when no unknows are present.
		Since nested maps are very rare, it is unlikely that this will become a problem in practice.
		If it does, we could either:
			- Check whether there are any unknowns in the nested object before creating a copy.
			- Resolve unknowns during unmarshal.
		Not worth the complexity as of now.
	*/

	mapPtr, diags := m.MapPtrAsAny(ctx)
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

	return m.NewNestedMapValue(ctx, mapPtr), nil
}
