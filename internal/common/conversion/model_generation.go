package conversion

import (
	"fmt"
	"reflect"
)

// CopyModel creates a new struct with the same values as the source struct. Fields in destination struct that are not in source are left with zero value.
// It panics if there are some structural problems so it should only happen during development.
func CopyModel[T any](src any) *T {
	dest := new(T)
	valSrc := reflect.ValueOf(src)
	valDest := reflect.ValueOf(dest)
	if valSrc.Kind() != reflect.Ptr || valDest.Kind() != reflect.Ptr {
		panic("params must be pointers")
	}
	valSrc = valSrc.Elem()
	valDest = valDest.Elem()
	if valSrc.Kind() != reflect.Struct || valDest.Kind() != reflect.Struct {
		panic("params must be pointers to structs")
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
				panic(fmt.Sprintf("field has different type: %s", name))
			}
		}
		if !valDest.Field(i).CanSet() {
			panic(fmt.Sprintf("field can't be set, probably unexported: %s", name))
		}
		valDest.Field(i).Set(valSrc.FieldByName(name))
	}
	return dest
}
