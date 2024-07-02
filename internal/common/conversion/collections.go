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
