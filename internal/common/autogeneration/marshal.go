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
func Marshal(src any) ([]byte, error) {
	valSrc := reflect.ValueOf(src)
	if valSrc.Kind() != reflect.Ptr {
		panic("src must be pointer")
	}
	valSrc = valSrc.Elem()
	if valSrc.Kind() != reflect.Struct {
		panic("src must be pointer to struct")
	}
	dest := make(map[string]any)
	for i := 0; i < valSrc.NumField(); i++ {
		nameAttr := xstrings.ToSnakeCase(valSrc.Type().Field(i).Name)
		valAttr := valSrc.Field(i)
		switch v := valAttr.Interface().(type) {
		case types.String:
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			dest[nameAttr] = v.ValueString()
		default:
			return nil, fmt.Errorf("marshal not supported yet for type %T for field %s", v, nameAttr)
		}
	}
	return json.Marshal(dest)
}

// Unmarshal gets a JSON (e.g. from an Atlas response) and unmarshals it into a Terraform model.
// It supports the following Terraform model types: String, Bool, Int64, Float64.
func Unmarshal(raw []byte, dest any) error {
	var src map[string]any
	if err := json.Unmarshal(raw, &src); err != nil {
		return err
	}
	return mapFields(src, dest)
}

func mapFields(src map[string]any, dest any) error {
	valDest := reflect.ValueOf(dest)
	if valDest.Kind() != reflect.Ptr {
		panic("dest must be pointer")
	}
	valDest = valDest.Elem()
	if valDest.Kind() != reflect.Struct {
		panic("dest must be pointer to struct")
	}
	for nameAttrSrc, valueAttrSrc := range src {
		if err := mapField(nameAttrSrc, valueAttrSrc, valDest); err != nil {
			return err
		}
	}
	return nil
}

func mapField(nameAttrSrc string, valueAttrSrc any, valDest reflect.Value) error {
	nameDest := xstrings.ToPascalCase(nameAttrSrc)
	fieldDest := valDest.FieldByName(nameDest)
	if !fieldDest.CanSet() {
		return nil // skip fields that cannot be set, are invalid or not found
	}
	switch v := valueAttrSrc.(type) {
	case string:
		return assignField(nameDest, fieldDest, types.StringValue(v))
	case bool:
		return assignField(nameDest, fieldDest, types.BoolValue(v))
	case float64: // number: try int or float
		if assignField(nameDest, fieldDest, types.Float64Value(v)) == nil {
			return nil
		}
		return assignField(nameDest, fieldDest, types.Int64Value(int64(v)))
	case nil:
		return nil // skip nil values, no need to set anything
	default:
		return fmt.Errorf("unmarshal not supported yet for type %T for field %s", v, nameAttrSrc)
	}
}

func assignField(nameDest string, fieldDest reflect.Value, valueDest attr.Value) error {
	valObj := reflect.ValueOf(valueDest)
	if !fieldDest.Type().AssignableTo(valObj.Type()) {
		return fmt.Errorf("unmarshal can't assign value to model field %s", nameDest)
	}
	fieldDest.Set(valObj)
	return nil
}
