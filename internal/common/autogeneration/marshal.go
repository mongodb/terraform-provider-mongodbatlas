package autogeneration

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/huandu/xstrings"
)

// Unmarshal gets a JSON (e.g. from an Atlas response) and unmarshals it into a Terraform model.
func Unmarshal(raw []byte, dest any) error {
	var fields map[string]any
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	return mapFields(fields, dest)
}

func mapFields(fields map[string]any, dest any) error {
	valDest := reflect.ValueOf(dest)
	if valDest.Kind() != reflect.Ptr {
		panic("dest must be pointer")
	}
	valDest = valDest.Elem()
	if valDest.Kind() != reflect.Struct {
		panic("dest must be pointer to struct")
	}
	for nameSrc, valueSrc := range fields {
		nameDest := xstrings.ToPascalCase(nameSrc)
		fieldDest := valDest.FieldByName(nameDest)
		if !fieldDest.CanSet() {
			continue // skip fields that cannot be set, are invalid or not found
		}
		switch v := valueSrc.(type) {
		case bool:
			return fmt.Errorf("not supported yet type %T for field %s", v, nameSrc)
		case float64:
			return fmt.Errorf("not supported yet type %T for field %s", v, nameSrc)
		case string:
			valObj := reflect.ValueOf(types.StringValue(v))
			fieldDest.Set(valObj)
		case nil:
			return fmt.Errorf("not supported yet type %T for field %s", v, nameSrc)
		default:
			return fmt.Errorf("not supported yet type %T for field %s", v, nameSrc)
		}
	}
	return nil
}
