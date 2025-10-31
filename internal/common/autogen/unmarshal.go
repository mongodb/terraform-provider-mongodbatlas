package autogen

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
)

// Unmarshal gets a JSON (e.g. from an Atlas response) and unmarshals it into a Terraform model.
// It supports the following Terraform model types: String, Bool, Int64, Float64, Object, List, Set.
// Map is not supported yet, will be done in CLOUDP-312797.
// Attributes that are in JSON but not in the model are ignored, no error is returned.
func Unmarshal(raw []byte, model any) error {
	if isEmptyJSON(raw) {
		return nil // Some operations return an empty response body, in that case there is no need to update the model.
	}
	var objJSON map[string]any
	if err := json.Unmarshal(raw, &objJSON); err != nil {
		return err
	}
	return unmarshalAttrs(objJSON, model)
}

func unmarshalAttrs(objJSON map[string]any, model any) error {
	valModel := reflect.ValueOf(model)
	if valModel.Kind() != reflect.Ptr {
		panic("model must be pointer")
	}
	valModel = valModel.Elem()
	if valModel.Kind() != reflect.Struct {
		panic("model must be pointer to struct")
	}
	for attrNameJSON, attrObjJSON := range objJSON {
		if err := unmarshalAttr(attrNameJSON, attrObjJSON, valModel); err != nil {
			return err
		}
	}
	return nil
}

func unmarshalAttr(attrNameJSON string, attrObjJSON any, valModel reflect.Value) error {
	attrNameModel := stringcase.Capitalize(attrNameJSON)
	fieldModel := valModel.FieldByName(attrNameModel)
	if !fieldModel.CanSet() {
		return nil // skip fields that cannot be set, are invalid or not found
	}
	if attrObjJSON == nil {
		return nil // skip nil values, no need to set anything
	}
	oldVal, ok := fieldModel.Interface().(attr.Value)
	if !ok {
		return fmt.Errorf("unmarshal trying to set non-Terraform attribute %s", attrNameModel)
	}
	valueType := oldVal.Type(context.Background())
	newValue, err := getTfAttr(attrObjJSON, valueType, oldVal, attrNameModel)
	if err != nil {
		return err
	}
	return setAttrTfModel(attrNameModel, fieldModel, newValue)
}

func setAttrTfModel(name string, field reflect.Value, val attr.Value) error {
	obj := reflect.ValueOf(val)
	if !field.Type().AssignableTo(obj.Type()) {
		return fmt.Errorf("unmarshal can't assign value to model field %s", name)
	}
	field.Set(obj)
	return nil
}

func setObjElmAttrModel(name string, value any, mapAttrs map[string]attr.Value, mapTypes map[string]attr.Type) error {
	nameChildTf := stringcase.ToSnakeCase(name)
	valueType, found := mapTypes[nameChildTf]
	if !found {
		return nil // skip attributes that are not in the model
	}
	newValue, err := getTfAttr(value, valueType, mapAttrs[nameChildTf], nameChildTf)
	if err != nil {
		return err
	}
	if newValue != nil {
		mapAttrs[nameChildTf] = newValue
	}
	return nil
}

func getTfAttr(value any, valueType attr.Type, oldVal attr.Value, name string) (attr.Value, error) {
	nameErr := stringcase.ToSnakeCase(name)
	switch v := value.(type) {
	case string:
		if valueType == types.StringType {
			return types.StringValue(v), nil
		}
		return nil, errUnmarshal(value, valueType, "String", nameErr)
	case bool:
		if valueType == types.BoolType {
			return types.BoolValue(v), nil
		}
		return nil, errUnmarshal(value, valueType, "Bool", nameErr)
	case float64:
		switch valueType {
		case types.Int64Type:
			return types.Int64Value(int64(v)), nil
		case types.Float64Type:
			return types.Float64Value(v), nil
		}
		return nil, errUnmarshal(value, valueType, "Number", nameErr)
	case map[string]any:
		if _, ok := valueType.(jsontypes.NormalizedType); ok {
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal object to JSON for attribute %s: %v", nameErr, err)
			}
			return jsontypes.NewNormalizedValue(string(jsonBytes)), nil
		}

		switch oldVal := oldVal.(type) {
		case types.Object:
			objNew, err := setObjAttrModel(oldVal, v)
			if err != nil {
				return nil, err
			}
			return objNew, nil
		case types.Map:
			mapNew, err := setMapAttrModel(oldVal, v)
			if err != nil {
				return nil, err
			}
			return mapNew, nil
		case customtypes.ObjectValueInterface:
			ctx := context.Background()
			valuePtr, diags := oldVal.ValuePtrAsAny(ctx)
			if diags.HasError() {
				return nil, fmt.Errorf("unmarshal failed to convert object: %v", diags)
			}

			err := unmarshalAttrs(v, valuePtr)
			if err != nil {
				return nil, err
			}

			objNew := oldVal.NewObjectValue(ctx, valuePtr)
			return objNew, nil
		case customtypes.NestedMapValueInterface:
			return getNestedMapValueTFAttr(context.Background(), v, oldVal)
		}
		return nil, errUnmarshal(value, valueType, "Object", nameErr)
	case []any:
		if _, ok := valueType.(jsontypes.NormalizedType); ok {
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal array to JSON for attribute %s: %v", nameErr, err)
			}
			return jsontypes.NewNormalizedValue(string(jsonBytes)), nil
		}

		switch oldVal := oldVal.(type) {
		case types.List:
			return setListAttrModel(oldVal, v, nameErr)
		case types.Set:
			return setSetAttrModel(oldVal, v, nameErr)
		case customtypes.ListValueInterface:
			return getListValueTFAttr(context.Background(), v, oldVal, nameErr)
		case customtypes.NestedListValueInterface:
			return getNestedListValueTFAttr(context.Background(), v, oldVal)
		case customtypes.SetValueInterface:
			return getSetValueTFAttr(context.Background(), v, oldVal, nameErr)
		case customtypes.NestedSetValueInterface:
			return getNestedSetValueTFAttr(context.Background(), v, oldVal)
		}
		return nil, errUnmarshal(value, valueType, "Array", nameErr)
	case nil:
		return nil, nil // skip nil values, no need to set anything
	}
	return nil, fmt.Errorf("unmarshal not supported yet for type %T for attribute %s", value, nameErr)
}

func errUnmarshal(value any, valueType attr.Type, typeReceived, name string) error {
	nameErr := stringcase.ToSnakeCase(name)
	parts := strings.Split(reflect.TypeOf(valueType).String(), ".")
	typeErr := parts[len(parts)-1]
	return fmt.Errorf("unmarshal of attribute %s expects type %s but got %s with value: %v", nameErr, typeErr, typeReceived, value)
}

func setObjAttrModel(obj types.Object, objJSON map[string]any) (attr.Value, error) {
	mapAttrs, mapTypes, err := getObjAttrsAndTypes(obj)
	if err != nil {
		return nil, err
	}
	for nameChild, valueChild := range objJSON {
		if err := setObjElmAttrModel(nameChild, valueChild, mapAttrs, mapTypes); err != nil {
			return nil, err
		}
	}
	objNew, diags := types.ObjectValue(obj.AttributeTypes(context.Background()), mapAttrs)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert JSON map to object: %v", diags)
	}
	return objNew, nil
}

func setMapAttrModel(m types.Map, objJSON map[string]any) (attr.Value, error) {
	mapAttrs := m.Elements()
	valueType := m.ElementType(context.Background())
	nullVal, err := getNullAttr(valueType)
	if err != nil {
		return nil, err
	}
	for nameChild, valueChild := range objJSON {
		oldVal, found := mapAttrs[nameChild]
		if !found {
			oldVal = nullVal
		}
		newValue, err := getTfAttr(valueChild, valueType, oldVal, nameChild)
		if err != nil {
			return nil, err
		}
		if newValue != nil {
			mapAttrs[nameChild] = newValue
		}
	}
	mapNew, diags := types.MapValue(valueType, mapAttrs)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert JSON map to map: %v", diags)
	}
	return mapNew, nil
}

func setListAttrModel(list types.List, arrayJSON []any, listName string) (attr.Value, error) {
	elmType := list.ElementType(context.Background())
	elms, err := getCollectionElements(arrayJSON, elmType, list.Elements(), listName)
	if err != nil {
		return nil, err
	}
	if len(elms) == 0 && len(list.Elements()) == 0 {
		// Keep current list if both model and JSON lists are zero-len (empty or null) so config is preserved.
		// It avoids inconsistent result after apply when user explicitly sets an empty list in config.
		return list, nil
	}
	listNew, diags := types.ListValue(elmType, elms)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert list to object: %v", diags)
	}
	return listNew, nil
}

func setSetAttrModel(set types.Set, arrayJSON []any, setName string) (attr.Value, error) {
	elmType := set.ElementType(context.Background())
	elms, err := getCollectionElements(arrayJSON, elmType, set.Elements(), setName)
	if err != nil {
		return nil, err
	}
	if len(elms) == 0 && len(set.Elements()) == 0 {
		// Keep current set if both model and JSON sets are zero-len (empty or null) so config is preserved.
		// It avoids inconsistent result after apply when user explicitly sets an empty set in config.
		return set, nil
	}
	setNew, diags := types.SetValue(elmType, elms)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert set to object: %v", diags)
	}
	return setNew, nil
}

func getCollectionElements(arrayJSON []any, valueType attr.Type, oldVals []attr.Value, collectionName string) ([]attr.Value, error) {
	elms := make([]attr.Value, len(arrayJSON))
	nullVal, err := getNullAttr(valueType)
	if err != nil {
		return nil, err
	}
	for i, item := range arrayJSON {
		oldVal := nullVal
		if i < len(oldVals) {
			oldVal = oldVals[i]
		}
		newValue, err := getTfAttr(item, valueType, oldVal, collectionName)
		if err != nil {
			return nil, err
		}
		if newValue != nil {
			elms[i] = newValue
		}
	}
	return elms, nil
}

func getObjAttrsAndTypes(obj types.Object) (mapAttrs map[string]attr.Value, mapTypes map[string]attr.Type, err error) {
	// mapTypes has all attributes, mapAttrs might not have them, e.g. in null or unknown objects
	mapAttrs = obj.Attributes()
	mapTypes = obj.AttributeTypes(context.Background())
	for attrName, attrType := range mapTypes {
		if _, found := mapAttrs[attrName]; found {
			continue // skip attributes that are already set
		}
		nullVal, err := getNullAttr(attrType)
		if err != nil {
			return nil, nil, err
		}
		mapAttrs[attrName] = nullVal
	}
	return mapAttrs, mapTypes, nil
}

func getNestedMapValueTFAttr(ctx context.Context, mapJSON map[string]any, m customtypes.NestedMapValueInterface) (attr.Value, error) {
	oldMapPtr, diags := m.MapPtrAsAny(ctx)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert map: %v", diags)
	}

	oldMapVal := reflect.ValueOf(oldMapPtr).Elem()

	mapPtr := m.NewEmptyMapPtr()
	mapVal := reflect.ValueOf(mapPtr).Elem()
	mapVal.Set(reflect.MakeMap(mapVal.Type()))

	mapElemType := mapVal.Type().Elem()
	for key, item := range mapJSON {
		keyVal := reflect.ValueOf(key)

		elementVal := reflect.New(mapElemType)

		if oldValue := oldMapVal.MapIndex(keyVal); oldValue.IsValid() {
			elementVal.Elem().Set(oldValue)
		}

		objJSON, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unmarshal of map item failed to convert to object: %v", item)
		}

		err := unmarshalAttrs(objJSON, elementVal.Interface())
		if err != nil {
			return nil, err
		}

		mapVal.SetMapIndex(keyVal, elementVal.Elem())
	}

	return m.NewNestedMapValue(ctx, mapPtr), nil
}

func getListValueTFAttr(ctx context.Context, arrayJSON []any, list customtypes.ListValueInterface, nameErr string) (attr.Value, error) {
	if len(arrayJSON) == 0 && len(list.Elements()) == 0 {
		// Keep current list if both model and JSON lists are zero-len (empty or null) so config is preserved.
		// It avoids inconsistent result after apply when user explicitly sets an empty list in config.
		return list, nil
	}

	slice, err := getArrayTFAttr(arrayJSON, list.ElementType(ctx), nameErr)
	if err != nil {
		return nil, err
	}

	listNew := list.NewListValue(ctx, slice)
	return listNew, nil
}

func getSetValueTFAttr(ctx context.Context, arrayJSON []any, set customtypes.SetValueInterface, nameErr string) (attr.Value, error) {
	if len(arrayJSON) == 0 && len(set.Elements()) == 0 {
		// Keep current set if both model and JSON lists are zero-len (empty or null) so config is preserved.
		// It avoids inconsistent result after apply when user explicitly sets an empty set in config.
		return set, nil
	}

	slice, err := getArrayTFAttr(arrayJSON, set.ElementType(ctx), nameErr)
	if err != nil {
		return nil, err
	}

	return set.NewSetValue(ctx, slice), nil
}

func getArrayTFAttr(arrayJSON []any, elemType attr.Type, nameErr string) ([]attr.Value, error) {
	slice := make([]attr.Value, len(arrayJSON))

	for i, item := range arrayJSON {
		newValue, err := getTfAttr(item, elemType, nil, nameErr)
		if err != nil {
			return nil, err
		}
		slice[i] = newValue
	}

	return slice, nil
}

func getNestedListValueTFAttr(ctx context.Context, arrayJSON []any, list customtypes.NestedListValueInterface) (attr.Value, error) {
	oldSlicePtr, diags := list.SlicePtrAsAny(ctx)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert list: %v", diags)
	}
	oldSliceVal := reflect.ValueOf(oldSlicePtr).Elem()
	oldSliceLen := oldSliceVal.Len()

	if len(arrayJSON) == 0 && oldSliceLen == 0 {
		// Keep current list if both model and JSON lists are zero-len (empty or null) so config is preserved.
		// It avoids inconsistent result after apply when user explicitly sets an empty list in config.
		return list, nil
	}

	slicePtr := list.NewEmptySlicePtr()
	sliceVal := reflect.ValueOf(slicePtr).Elem()
	sliceVal.Set(reflect.MakeSlice(sliceVal.Type(), len(arrayJSON), len(arrayJSON)))

	for i, item := range arrayJSON {
		elementVal := sliceVal.Index(i)
		if i < oldSliceLen {
			elementVal.Set(oldSliceVal.Index(i))
		}
		elementPtr := elementVal.Addr().Interface()
		objJSON, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unmarshal of array item failed to convert to object: %v", item)
		}
		err := unmarshalAttrs(objJSON, elementPtr)
		if err != nil {
			return nil, err
		}
	}

	return list.NewNestedListValue(ctx, slicePtr), nil
}

func getNestedSetValueTFAttr(ctx context.Context, arrayJSON []any, set customtypes.NestedSetValueInterface) (attr.Value, error) {
	if len(arrayJSON) == 0 && set.Len() == 0 {
		// Keep current set if both model and JSON lists are zero-len (empty or null) so config is preserved.
		// It avoids inconsistent result after apply when user explicitly sets an empty set in config.
		return set, nil
	}

	slicePtr := set.NewEmptySlicePtr()
	sliceVal := reflect.ValueOf(slicePtr).Elem()
	sliceVal.Set(reflect.MakeSlice(sliceVal.Type(), len(arrayJSON), len(arrayJSON)))

	for i, item := range arrayJSON {
		elementPtr := sliceVal.Index(i).Addr().Interface()
		objJSON, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unmarshal of array item failed to convert to object: %v", item)
		}
		err := unmarshalAttrs(objJSON, elementPtr)
		if err != nil {
			return nil, err
		}
	}

	return set.NewNestedSetValue(ctx, slicePtr), nil
}
