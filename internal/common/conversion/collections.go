package conversion

import "reflect"

// HasElementsSliceOrMap checks if param is a non-empty slice or map
func HasElementsSliceOrMap(value any) bool {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Map {
		return v.Len() > 0
	}
	return false
}

// ToAnySlicePointer converts to a slice pointer of any as needed in some Atlas SDK Go structs
func ToAnySlicePointer(value *[]map[string]any) *[]any {
	if value == nil {
		return nil
	}
	ret := make([]any, len(*value))
	for i, item := range *value {
		ret[i] = item
	}
	return &ret
}
