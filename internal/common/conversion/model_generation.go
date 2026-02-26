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
	if valSrc.Kind() != reflect.Pointer || valDest.Kind() != reflect.Pointer {
		panic("params must be pointers")
	}
	valSrc = valSrc.Elem()
	valDest = valDest.Elem()
	if valSrc.Kind() != reflect.Struct || valDest.Kind() != reflect.Struct {
		panic("params must be pointers to structs")
	}
	typeSrc := valSrc.Type()
	for fieldDest := range valDest.Type().Fields() {
		fieldSrc, found := typeSrc.FieldByName(fieldDest.Name)
		if !found {
			continue
		}
		if fieldDest.Type != fieldSrc.Type {
			panic(fmt.Sprintf("field has different type: %s", fieldDest.Name))
		}
		destVal := valDest.FieldByIndex(fieldDest.Index)
		if !destVal.CanSet() {
			panic(fmt.Sprintf("field can't be set, probably unexported: %s", fieldDest.Name))
		}
		destVal.Set(valSrc.FieldByName(fieldDest.Name))
	}
	return dest
}
