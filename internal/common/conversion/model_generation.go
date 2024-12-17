package conversion

import (
	"fmt"
	"reflect"
)

// CopyModel creates a new struct with the same values as the source struct. Fields in destination struct that are not in source are left with zero value.
func CopyModel[T any](src any) (*T, error) {
	dest := new(T)
	valSrc := reflect.ValueOf(src)
	valDest := reflect.ValueOf(dest)
	if valSrc.Kind() != reflect.Ptr || valDest.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("params must be pointers")
	}
	valSrc = valSrc.Elem()
	valDest = valDest.Elem()
	if valSrc.Kind() != reflect.Struct || valDest.Kind() != reflect.Struct {
		return nil, fmt.Errorf("params must be pointers to structs")
	}
	typeSrc := valSrc.Type()
	typeDest := valDest.Type()
	for i := range typeDest.NumField() {
		fieldDest := typeDest.Field(i)
		name := fieldDest.Name
		{
			fieldSrc, found := typeSrc.FieldByName(name)
			if !found {
				continue
			}
			if fieldDest.Type != fieldSrc.Type {
				return nil, fmt.Errorf("field has different type: %s", name)
			}
		}
		if !valDest.Field(i).CanSet() {
			return nil, fmt.Errorf("field can't be set, probably unexported: %s", name)
		}
		valDest.Field(i).Set(valSrc.FieldByName(name))
	}
	return dest, nil
}
