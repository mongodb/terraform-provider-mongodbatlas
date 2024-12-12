package conversion

import (
	"fmt"
	"reflect"
)

// TODO: Add comment and tests
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
	for i := 0; i < typeSrc.NumField(); i++ {
		fieldSrc := typeSrc.Field(i)
		fieldDest, found := typeDest.FieldByName(fieldSrc.Name)
		if found {
			if fieldSrc.Type != fieldDest.Type {
				return nil, fmt.Errorf("field %s has different type in source and destination", fieldSrc.Name)
			}
			valDest.FieldByName(fieldDest.Name).Set(valSrc.Field(i))
		}
	}
	return dest, nil
}
