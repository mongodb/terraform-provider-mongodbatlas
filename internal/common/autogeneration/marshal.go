package autogeneration

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/huandu/xstrings"
)

// Marshal gets Terraform model and marshals it into JSON (e.g. for an Atlas request).
func Marshal(model any) ([]byte, error) {
	valModel := reflect.ValueOf(model)
	if valModel.Kind() != reflect.Ptr {
		panic("Marshal expects a pointer")
	}
	valModel = valModel.Elem()
	if valModel.Kind() != reflect.Struct {
		panic("Marshal expects a pointer to struct")
	}
	objJSON, err := marshalAttrs(valModel)
	if err != nil {
		return nil, err
	}
	return json.Marshal(objJSON)
}

func marshalAttrs(valModel reflect.Value) (map[string]any, error) {
	objJSON := make(map[string]any)
	for i := 0; i < valModel.NumField(); i++ {
		attrNameModel := valModel.Type().Field(i).Name
		attrValModel := valModel.Field(i)
		if err := marshalAttr(attrNameModel, attrValModel, objJSON); err != nil {
			return nil, err
		}
	}
	return objJSON, nil
}

func marshalAttr(attrNameModel string, attrValModel reflect.Value, objJSON map[string]any) error {
	attrNameJSON := xstrings.ToSnakeCase(attrNameModel)
	if v, ok := attrValModel.Interface().(attr.Value); ok {
		if v.IsNull() || v.IsUnknown() {
			return nil // skip nil or unknown values
		}
	}
	switch v := attrValModel.Interface().(type) {
	case types.String:
		objJSON[attrNameJSON] = v.ValueString()
	case types.Int64:
		objJSON[attrNameJSON] = v.ValueInt64()
	default:
		return fmt.Errorf("marshal not supported yet for type %T for field %s", v, attrNameJSON)
	}
	return nil
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

func unmarshalAttrs(objJSON map[string]any, model any) error {
	valModel := reflect.ValueOf(model)
	if valModel.Kind() != reflect.Ptr {
		panic("dest must be pointer")
	}
	valModel = valModel.Elem()
	if valModel.Kind() != reflect.Struct {
		panic("dest must be pointer to struct")
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
