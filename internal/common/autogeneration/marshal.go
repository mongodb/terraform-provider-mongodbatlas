package autogeneration

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/huandu/xstrings"
)

// Unmarshal gets a JSON (e.g. from an Atlas response) and unmarshals it into a Terraform model.
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
		return assignField(fieldDest, types.StringValue(v))
	case bool:
		return assignField(fieldDest, types.BoolValue(v))
	case float64: // number: try int or float
		if assignField(fieldDest, types.Float64Value(v)) == nil {
			return nil
		}
		return assignField(fieldDest, types.Int64Value(int64(v)))
	case nil:
		return nil // skip nil values, no need to set anything
	default:
		return fmt.Errorf("not supported yet type %T for field %s", v, nameAttrSrc)
	}
}

func assignField(fieldDest reflect.Value, valueDest attr.Value) error {
	valObj := reflect.ValueOf(valueDest)
	if !fieldDest.Type().AssignableTo(valObj.Type()) {
		return fmt.Errorf("can't assign value to model field %s", fieldDest.Type().Name())
	}
	fieldDest.Set(valObj)
	return nil
}
