package autogeneration

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/huandu/xstrings"
)

const (
	tagKey               = "autogeneration"
	tagValOmitJSON       = "omitjson"
	tagValOmitJSONUpdate = "omitjsonupdate"
)

// Marshal gets a Terraform model and marshals it into JSON (e.g. for an Atlas request).
// It supports the following Terraform model types: String, Bool, Int64, Float64, Object, Map, List, Set.
// Attributes that are null or unknown are not marshaled.
// Attributes with autogeneration tag `omitjson` are never marshaled, this only applies to the root model.
// Attributes with autogeneration tag `omitjsonupdate` are not marshaled if isUpdate is true, this only applies to the root model.
func Marshal(model any, isUpdate bool) ([]byte, error) {
	valModel := reflect.ValueOf(model)
	if valModel.Kind() != reflect.Ptr {
		panic("model must be pointer")
	}
	valModel = valModel.Elem()
	if valModel.Kind() != reflect.Struct {
		panic("model must be pointer to struct")
	}
	objJSON, err := marshalAttrs(valModel, isUpdate)
	if err != nil {
		return nil, err
	}
	return json.Marshal(objJSON)
}

// Unmarshal gets a JSON (e.g. from an Atlas response) and unmarshals it into a Terraform model.
// It supports the following Terraform model types: String, Bool, Int64, Float64.
// Attributes that are in JSON but not in the model are ignored, no error is returned.
// Object attributes that are unknown are converted to null as all values must be known in the response state.
func Unmarshal(raw []byte, model any) error {
	var objJSON map[string]any
	if err := json.Unmarshal(raw, &objJSON); err != nil {
		return err
	}
	return unmarshalAttrs(objJSON, model)
}

func marshalAttrs(valModel reflect.Value, isUpdate bool) (map[string]any, error) {
	objJSON := make(map[string]any)
	for i := 0; i < valModel.NumField(); i++ {
		attrTypeModel := valModel.Type().Field(i)
		tag := attrTypeModel.Tag.Get(tagKey)
		if tag == tagValOmitJSON {
			continue // skip fields with tag `omitjson`
		}
		if isUpdate && tag == tagValOmitJSONUpdate {
			continue // skip fields with tag `omitjsonupdate` if in update mode
		}
		attrNameModel := attrTypeModel.Name
		attrValModel := valModel.Field(i)
		if err := marshalAttr(attrNameModel, attrValModel, objJSON); err != nil {
			return nil, err
		}
	}
	return objJSON, nil
}

func marshalAttr(attrNameModel string, attrValModel reflect.Value, objJSON map[string]any) error {
	attrNameJSON := xstrings.ToCamelCase(attrNameModel)
	obj, ok := attrValModel.Interface().(attr.Value)
	if !ok {
		panic("marshal expects only Terraform types in the model")
	}
	val, err := getAttr(obj)
	if err != nil {
		return err
	}
	if val != nil {
		objJSON[attrNameJSON] = val
	}
	return nil
}

func getAttr(val attr.Value) (any, error) {
	if val.IsNull() || val.IsUnknown() {
		return nil, nil // skip null or unknown values
	}
	switch v := val.(type) {
	case types.String:
		return v.ValueString(), nil
	case types.Bool:
		return v.ValueBool(), nil
	case types.Int64:
		return v.ValueInt64(), nil
	case types.Float64:
		return v.ValueFloat64(), nil
	case types.Object:
		return getMapAttr(v.Attributes(), false)
	case types.Map:
		return getMapAttr(v.Elements(), true)
	case types.List:
		return getListAttr(v.Elements())
	case types.Set:
		return getListAttr(v.Elements())
	default:
		return nil, fmt.Errorf("unmarshal not supported yet for type %T", v)
	}
}

func getListAttr(elms []attr.Value) (any, error) {
	arr := make([]any, 0)
	for _, attr := range elms {
		valChild, err := getAttr(attr)
		if err != nil {
			return nil, err
		}
		if valChild != nil {
			arr = append(arr, valChild)
		}
	}
	return arr, nil
}

// getMapAttr gets a map of attributes and returns a map of JSON attributes.
// keepKeyCase is used for types.Map to keep key case. However, we want to use JSON key case for types.Object
func getMapAttr(elms map[string]attr.Value, keepKeyCase bool) (any, error) {
	objJSON := make(map[string]any)
	for name, attr := range elms {
		valChild, err := getAttr(attr)
		if err != nil {
			return nil, err
		}
		if valChild != nil {
			nameJSON := xstrings.ToCamelCase(name)
			if keepKeyCase {
				nameJSON = name
			}
			objJSON[nameJSON] = valChild
		}
	}
	return objJSON, nil
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
	convertUnknownToNull(valModel)
	return nil
}

func unmarshalAttr(attrNameJSON string, attrObjJSON any, valModel reflect.Value) error {
	attrNameModel := xstrings.ToPascalCase(attrNameJSON)
	fieldModel := valModel.FieldByName(attrNameModel)
	if !fieldModel.CanSet() {
		return nil // skip fields that cannot be set, are invalid or not found
	}
	switch v := attrObjJSON.(type) {
	case string:
		return setAttrTfModel(attrNameModel, fieldModel, types.StringValue(v))
	case bool:
		return setAttrTfModel(attrNameModel, fieldModel, types.BoolValue(v))
	case float64: // number: try int or float
		if setAttrTfModel(attrNameModel, fieldModel, types.Float64Value(v)) == nil {
			return nil
		}
		return setAttrTfModel(attrNameModel, fieldModel, types.Int64Value(int64(v)))
	case nil:
		return nil // skip nil values, no need to set anything
	case map[string]any:
		obj, ok := fieldModel.Interface().(types.Object)
		if !ok {
			return fmt.Errorf("unmarshal expects object for field %s", attrNameJSON)
		}
		objNew, err := setObjAttrModel(obj, v)
		if err != nil {
			return err
		}
		return setAttrTfModel(attrNameModel, fieldModel, objNew)
	default:
		return fmt.Errorf("unmarshal not supported yet for type %T for field %s", v, attrNameJSON)
	}
}

func setAttrTfModel(name string, field reflect.Value, val attr.Value) error {
	obj := reflect.ValueOf(val)
	if !field.Type().AssignableTo(obj.Type()) {
		return fmt.Errorf("unmarshal can't assign value to model field %s", name)
	}
	field.Set(obj)
	return nil
}

func setAttrModel(name string, value any, mapAttrs map[string]attr.Value, mapTypes map[string]attr.Type) error {
	nameChildTf := xstrings.ToSnakeCase(name)
	valueType, found := mapTypes[nameChildTf]
	if !found {
		return nil // skip attributes that are not in the model
	}
	switch v := value.(type) {
	case string:
		if valueType == types.StringType {
			mapAttrs[nameChildTf] = types.StringValue(v)
			return nil
		}
		return fmt.Errorf("unmarshal gets incorrect string for field %s, value: %v", nameChildTf, v)
	case bool:
		if valueType == types.BoolType {
			mapAttrs[nameChildTf] = types.BoolValue(v)
			return nil
		}
		return fmt.Errorf("unmarshal gets incorrect bool for field %s, value: %v", nameChildTf, v)
	case float64:
		switch valueType {
		case types.Int64Type:
			mapAttrs[nameChildTf] = types.Int64Value(int64(v))
			return nil
		case types.Float64Type:
			mapAttrs[nameChildTf] = types.Float64Value(v)
			return nil
		}
		return fmt.Errorf("unmarshal gets incorrect number for field %s, value: %v", nameChildTf, v)
	case map[string]any:
		mapAttr, found := mapAttrs[nameChildTf]
		obj, ok := mapAttr.(types.Object)
		if !found || !ok {
			return fmt.Errorf("unmarshal gets incorrect object for field %s, value: %v", nameChildTf, v)
		}
		objNew, err := setObjAttrModel(obj, v)
		if err != nil {
			return err
		}
		mapAttrs[nameChildTf] = objNew
		return nil
	case nil:
		return nil // skip nil values, no need to set anything
	}
	return fmt.Errorf("unmarshal not supported yet for type %T for field %s", value, nameChildTf)
}

func setObjAttrModel(obj types.Object, objJSON map[string]any) (attr.Value, error) {
	mapAttrs, mapTypes, err := getObjAttrsAndTypes(obj)
	if err != nil {
		return nil, err
	}
	for nameChild, valueChild := range objJSON {
		if err := setAttrModel(nameChild, valueChild, mapAttrs, mapTypes); err != nil {
			return nil, err
		}
	}
	objNew, diags := types.ObjectValue(obj.AttributeTypes(context.Background()), mapAttrs)
	if diags.HasError() {
		return nil, fmt.Errorf("unmarshal failed to convert map to object: %v", diags)
	}
	return objNew, nil
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
		return nil, fmt.Errorf("unmarshal to get null value not supported yet for type %T", attrType)
	}
}

func convertUnknownToNull(valModel reflect.Value) {
	for i := 0; i < valModel.NumField(); i++ {
		field := valModel.Field(i)
		if obj, ok := field.Interface().(types.Object); ok {
			if obj.IsUnknown() && field.CanSet() {
				field.Set(reflect.ValueOf(types.ObjectNull(obj.AttributeTypes(context.Background()))))
			}
		}
	}
}
