package schemafunc

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func IsUnknown(obj reflect.Value) bool {
	method := obj.MethodByName("IsUnknown")
	if !method.IsValid() {
		panic(fmt.Sprintf("IsUnknown method not found for %v", obj))
	}
	results := method.Call([]reflect.Value{})
	if len(results) != 1 {
		panic(fmt.Sprintf("IsUnknown method must return a single value, got %v", results))
	}
	result := results[0]
	response, ok := result.Interface().(bool)
	if !ok {
		panic(fmt.Sprintf("IsUnknown method must return a bool, got %v", result))
	}
	return response
}

func HasUnknowns(obj any) bool {
	valObj := reflect.ValueOf(obj)
	if valObj.Kind() != reflect.Ptr {
		panic("params must be pointer")
	}
	valObj = valObj.Elem()
	if valObj.Kind() != reflect.Struct {
		panic("params must be pointer to struct")
	}
	typeObj := valObj.Type()
	for i := range typeObj.NumField() {
		field := valObj.Field(i)
		if IsUnknown(field) {
			return true
		}
	}
	return false
}

// CopyUnknowns use reflection to copy unknown fields from src to dest.
// The alternative without reflection would need to pass every field in a struct.
// The implementation is similar to internal/common/conversion/model_generation.go#CopyModel
func CopyUnknowns(ctx context.Context, src, dest any, keepUnknown []string) {
	valSrc := reflect.ValueOf(src)
	valDest := reflect.ValueOf(dest)
	if valSrc.Kind() != reflect.Ptr || valDest.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("params must be pointers %v %v\n", src, dest))
	}
	valSrc = valSrc.Elem()
	valDest = valDest.Elem()
	if valSrc.Kind() != reflect.Struct || valDest.Kind() != reflect.Struct {
		panic(fmt.Sprintf("params must be pointers to structs: %v %v\n", src, dest))
	}
	typeSrc := valSrc.Type()
	typeDest := valDest.Type()
	for i := range typeDest.NumField() {
		fieldDest := typeDest.Field(i)
		name := fieldDest.Name
		if slices.Contains(keepUnknown, name) {
			continue
		}
		_, found := typeSrc.FieldByName(name)
		if !found || !IsUnknown(valDest.Field(i)) || !valDest.Field(i).CanSet() {
			continue
		}
		tflog.Info(ctx, fmt.Sprintf("Copying unknown field: %s", name))
		valDest.Field(i).Set(valSrc.FieldByName(name))
	}
}
