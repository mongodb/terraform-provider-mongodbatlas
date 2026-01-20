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

	// Iterate over model fields and look up corresponding JSON properties
	structType := valModel.Type()
	for i := range structType.NumField() {
		field := structType.Field(i)
		fieldModel := valModel.Field(i)

		if !fieldModel.CanSet() {
			continue // skip fields that cannot be set
		}

		tags := GetPropertyTags(&field)
		apiName := getAPINameFromTag(field.Name, tags)

		// Look up the JSON property
		attrObjJSON, ok := objJSON[apiName]
		if !ok {
			continue // skip fields not found in JSON (attributes in JSON but not in model are ignored)
		}

		if attrObjJSON == nil {
			continue // skip nil values, no need to set anything
		}

		if err := unmarshalAttr(attrObjJSON, fieldModel, &field); err != nil {
			return err
		}
	}
	return nil
}

func unmarshalAttr(attrObjJSON any, fieldModel reflect.Value, structField *reflect.StructField) error {
	attrNameModel := structField.Name
	tags := GetPropertyTags(structField)

	oldVal, ok := fieldModel.Interface().(attr.Value)
	if !ok {
		return fmt.Errorf("unmarshal trying to set non-Terraform attribute %s", attrNameModel)
	}

	if !oldVal.IsNull() && !oldVal.IsUnknown() { // Check if oldVal is a known value
		if tags.Sensitive {
			return nil // skip sensitive fields that are already set in the plan/state to avoid overwriting with redacted values
		}
	}

	if tags.ListAsMap {
		attrObjJSON = ModifyJSONFromListToMap(attrObjJSON)
	}

	valueType := oldVal.Type(context.Background())
	newValue, err := getTfAttr(attrObjJSON, valueType, oldVal, attrNameModel, tags.SkipListMerge)
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

func getTfAttr(value any, valueType attr.Type, oldVal attr.Value, name string, skipListMerge bool) (attr.Value, error) {
	nameErr := stringcase.ToSnakeCase(name)
	if _, ok := valueType.(jsontypes.NormalizedType); ok {
		return getNormalizedJSONAttrValue(value, nameErr)
	}
	switch v := value.(type) {
	case string:
		if valueType == types.StringType {
			return types.StringValue(v), nil
		}
		return nil, errUnmarshal(valueType, "String", nameErr)
	case bool:
		if valueType == types.BoolType {
			return types.BoolValue(v), nil
		}
		return nil, errUnmarshal(valueType, "Bool", nameErr)
	case float64:
		switch valueType {
		case types.Int64Type:
			return types.Int64Value(int64(v)), nil
		case types.Float64Type:
			return types.Float64Value(v), nil
		}
		return nil, errUnmarshal(valueType, "Number", nameErr)
	case map[string]any:
		switch oldVal := oldVal.(type) {
		case customtypes.ObjectValueInterface:
			return getObjectValueTFAttr(context.Background(), v, oldVal)
		case customtypes.MapValueInterface:
			return getMapValueTFAttr(context.Background(), v, oldVal)
		case customtypes.NestedMapValueInterface:
			return getNestedMapValueTFAttr(context.Background(), v, oldVal)
		}
		return nil, errUnmarshal(valueType, "Object", nameErr)
	case []any:
		switch oldVal := oldVal.(type) {
		case customtypes.ListValueInterface:
			return getListValueTFAttr(context.Background(), v, oldVal, nameErr)
		case customtypes.NestedListValueInterface:
			return getNestedListValueTFAttr(context.Background(), v, oldVal, skipListMerge)
		case customtypes.SetValueInterface:
			return getSetValueTFAttr(context.Background(), v, oldVal, nameErr)
		case customtypes.NestedSetValueInterface:
			return getNestedSetValueTFAttr(context.Background(), v, oldVal)
		}
		return nil, errUnmarshal(valueType, "Array", nameErr)
	case nil:
		return nil, nil // skip nil values, no need to set anything
	}
	return nil, fmt.Errorf("unmarshal not supported yet for type %T for attribute %s", value, nameErr)
}

func getNormalizedJSONAttrValue(value any, nameErr string) (attr.Value, error) {
	// Marshal the value as a JSON string and return a jsontypes.NormalizedValue.
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value to JSON for attribute %s", nameErr)
	}
	return jsontypes.NewNormalizedValue(string(jsonBytes)), nil
}

func errUnmarshal(valueType attr.Type, typeReceived, name string) error {
	nameErr := stringcase.ToSnakeCase(name)
	parts := strings.Split(reflect.TypeOf(valueType).String(), ".")
	typeErr := parts[len(parts)-1]
	return fmt.Errorf("unmarshal of attribute %s expects type %s but got %s", nameErr, typeErr, typeReceived)
}

func getObjectValueTFAttr(ctx context.Context, objJSON map[string]any, obj customtypes.ObjectValueInterface) (attr.Value, error) {
	valuePtr, diags := obj.ValuePtrAsAny(ctx)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert object: %v", diags)
	}

	err := unmarshalAttrs(objJSON, valuePtr)
	if err != nil {
		return nil, err
	}

	return obj.NewObjectValue(ctx, valuePtr), nil
}

func getMapValueTFAttr(ctx context.Context, mapJSON map[string]any, m customtypes.MapValueInterface) (attr.Value, error) {
	if len(mapJSON) == 0 && len(m.Elements()) == 0 {
		// Keep current map if both model and JSON maps are zero-len (empty or null) so config is preserved.
		// It avoids inconsistent result after apply when user explicitly sets an empty map in config.
		return m, nil
	}

	mapAttrs := make(map[string]attr.Value, len(mapJSON))
	elemType := m.ElementType(ctx)

	for key, item := range mapJSON {
		value, err := getTfAttr(item, elemType, nil, key, false)
		if err != nil {
			return nil, err
		}
		if value != nil {
			mapAttrs[key] = value
		}
	}

	return m.NewMapValue(ctx, mapAttrs), nil
}

func getNestedMapValueTFAttr(ctx context.Context, mapJSON map[string]any, m customtypes.NestedMapValueInterface) (attr.Value, error) {
	oldMapPtr, diags := m.MapPtrAsAny(ctx)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert map: %v", diags)
	}

	oldMapVal := reflect.ValueOf(oldMapPtr).Elem()
	oldMapLen := oldMapVal.Len()

	if len(mapJSON) == 0 && oldMapLen == 0 {
		// Keep current map if both model and JSON map are zero-len (empty or null) so config is preserved.
		// It avoids inconsistent result after apply when user explicitly sets an empty map in config.
		return m, nil
	}

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
			return nil, fmt.Errorf("failed to unmarshal map item, expected object but got %T", item)
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
		value, err := getTfAttr(item, elemType, nil, nameErr, false)
		if err != nil {
			return nil, err
		}
		slice[i] = value
	}

	return slice, nil
}

func getNestedListValueTFAttr(ctx context.Context, arrayJSON []any, list customtypes.NestedListValueInterface, skipListMerge bool) (attr.Value, error) {
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
		if !skipListMerge && i < oldSliceLen {
			elementVal.Set(oldSliceVal.Index(i))
		}
		elementPtr := elementVal.Addr().Interface()
		objJSON, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed to unmarshal array item, expected object but got %T", item)
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
			return nil, fmt.Errorf("failed to unmarshal set item, expected object but got %T", item)
		}
		err := unmarshalAttrs(objJSON, elementPtr)
		if err != nil {
			return nil, err
		}
	}

	return set.NewNestedSetValue(ctx, slicePtr), nil
}
