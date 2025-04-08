package autogeneration

import (
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
	case types.Int64:
		return v.ValueInt64(), nil
	case types.Float64:
		return v.ValueFloat64(), nil
	case types.Object:
		return getMapAttr(v.Attributes())
	case types.Map:
		return getMapAttr(v.Elements())
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

func getMapAttr(elms map[string]attr.Value) (any, error) {
	obj := make(map[string]any)
	for name, attr := range elms {
		valChild, err := getAttr(attr)
		if err != nil {
			return nil, err
		}
		if valChild != nil {
			obj[name] = valChild
		}
	}
	return obj, nil
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
	switch v := attrObjJSON.(type) {
	case string:
		return setAttrModel(attrNameModel, fieldModel, types.StringValue(v))
	case bool:
		return setAttrModel(attrNameModel, fieldModel, types.BoolValue(v))
	case float64: // number: try int or float
		if setAttrModel(attrNameModel, fieldModel, types.Float64Value(v)) == nil {
			return nil
		}
		return setAttrModel(attrNameModel, fieldModel, types.Int64Value(int64(v)))
	case nil:
		return nil // skip nil values, no need to set anything
	default:
		return fmt.Errorf("unmarshal not supported yet for type %T for field %s", v, attrNameJSON)
	}
}

func setAttrModel(name string, field reflect.Value, val attr.Value) error {
	obj := reflect.ValueOf(val)
	if !field.Type().AssignableTo(obj.Type()) {
		return fmt.Errorf("unmarshal can't assign value to model field %s", name)
	}
	field.Set(obj)
	return nil
}
