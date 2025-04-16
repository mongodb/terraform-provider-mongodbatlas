package autogen

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/huandu/xstrings"
)

// Unmarshal gets a JSON (e.g. from an Atlas response) and unmarshals it into a Terraform model.
// It supports the following Terraform model types: String, Bool, Int64, Float64, Object, List, Set.
// Map is not supported yet, will be done in CLOUDP-312797.
// Attributes that are in JSON but not in the model are ignored, no error is returned.
func Unmarshal(raw []byte, model any) error {
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
	attrNameModel := xstrings.ToPascalCase(attrNameJSON)
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
	newValue, err := getTfAttr(attrObjJSON, valueType, oldVal)
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

func setMapAttrModel(name string, value any, mapAttrs map[string]attr.Value, mapTypes map[string]attr.Type) error {
	nameChildTf := xstrings.ToSnakeCase(name)
	valueType, found := mapTypes[nameChildTf]
	if !found {
		return nil // skip attributes that are not in the model
	}
	newValue, err := getTfAttr(value, valueType, mapAttrs[nameChildTf])
	if err != nil {
		return err
	}
	if newValue != nil {
		mapAttrs[nameChildTf] = newValue
	}
	return nil
}

func getTfAttr(value any, valueType attr.Type, oldVal attr.Value) (attr.Value, error) {
	switch v := value.(type) {
	case string:
		if valueType == types.StringType {
			return types.StringValue(v), nil
		}
		return nil, fmt.Errorf("unmarshal gets incorrect string for value: %v", v)
	case bool:
		if valueType == types.BoolType {
			return types.BoolValue(v), nil
		}
		return nil, fmt.Errorf("unmarshal gets incorrect bool for value: %v", v)
	case float64:
		switch valueType {
		case types.Int64Type:
			return types.Int64Value(int64(v)), nil
		case types.Float64Type:
			return types.Float64Value(v), nil
		}
		return nil, fmt.Errorf("unmarshal gets incorrect number for value: %v", v)
	case map[string]any:
		obj, ok := oldVal.(types.Object)
		if !ok {
			return nil, fmt.Errorf("unmarshal gets incorrect object for value: %v", v)
		}
		objNew, err := setObjAttrModel(obj, v)
		if err != nil {
			return nil, err
		}
		return objNew, nil
	case []any:
		if list, ok := oldVal.(types.List); ok {
			listNew, err := setListAttrModel(list, v)
			if err != nil {
				return nil, err
			}
			return listNew, nil
		}
		if set, ok := oldVal.(types.Set); ok {
			setNew, err := setSetAttrModel(set, v)
			if err != nil {
				return nil, err
			}
			return setNew, nil
		}
		return nil, fmt.Errorf("unmarshal gets incorrect array for value: %v", v)
	case nil:
		return nil, nil // skip nil values, no need to set anything
	}
	return nil, fmt.Errorf("unmarshal not supported yet for type %T", value)
}

func setObjAttrModel(obj types.Object, objJSON map[string]any) (attr.Value, error) {
	mapAttrs, mapTypes, err := getObjAttrsAndTypes(obj)
	if err != nil {
		return nil, err
	}
	for nameChild, valueChild := range objJSON {
		if err := setMapAttrModel(nameChild, valueChild, mapAttrs, mapTypes); err != nil {
			return nil, err
		}
	}
	objNew, diags := types.ObjectValue(obj.AttributeTypes(context.Background()), mapAttrs)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert map to object: %v", diags)
	}
	return objNew, nil
}

func setListAttrModel(list types.List, arrayJSON []any) (attr.Value, error) {
	elmType := list.ElementType(context.Background())
	elms, err := getCollectionElements(arrayJSON, elmType, list.Elements())
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

func setSetAttrModel(set types.Set, arrayJSON []any) (attr.Value, error) {
	elmType := set.ElementType(context.Background())
	elms, err := getCollectionElements(arrayJSON, elmType, set.Elements())
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

func getCollectionElements(arrayJSON []any, valueType attr.Type, oldVals []attr.Value) ([]attr.Value, error) {
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
		newValue, err := getTfAttr(item, valueType, oldVal)
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
