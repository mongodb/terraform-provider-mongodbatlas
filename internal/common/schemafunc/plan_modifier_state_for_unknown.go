package schemafunc

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
func CopyUnknowns(ctx context.Context, src, dest any, keepUnknown []string, nestedStructMapping map[string]string) {
	valSrc := reflect.ValueOf(src)
	valDest := reflect.ValueOf(dest)
	if valSrc.Kind() != reflect.Ptr || valDest.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("params must be pointers %T %T\n", src, dest))
	}
	valSrc = valSrc.Elem()
	valDest = valDest.Elem()
	if valSrc.Kind() != reflect.Struct || valDest.Kind() != reflect.Struct {
		panic(fmt.Sprintf("params must be pointers to structs: %T, %T and not nil: (%v, %v)\n", src, dest, src, dest))
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
		if !found || !valDest.Field(i).CanSet() {
			continue
		}
		if IsUnknown(valDest.Field(i)) {
			tflog.Info(ctx, fmt.Sprintf("Copying unknown field: %s\n", name))
			valDest.Field(i).Set(valSrc.FieldByName(name))
			continue
		}
		nestedType := nestedStructMapping[name]
		if nestedType != "" {
			fmt.Printf("Processing nested field: %s with type %s\n", name, nestedType)
			if nestedType == "object" {
				objValueSrc := valSrc.FieldByName(name).Interface().(types.Object)
				objValueDest := valDest.FieldByName(name).Interface().(types.Object)
				objValueNew := CopyUnknownsFromObject(ctx, objValueSrc, objValueDest)
				valDest.Field(i).Set(reflect.ValueOf(objValueNew))
			} else if nestedType == "list" {
				listValueSrc := valSrc.FieldByName(name).Interface().(types.List)
				listValueDest := valDest.FieldByName(name).Interface().(types.List)
				listValueNew := CopyUnknownsFromList(ctx, listValueSrc, listValueDest)
				valDest.Field(i).Set(reflect.ValueOf(listValueNew))
			} else {
				panic(fmt.Sprintf("nested type not supported yet: %s", nestedType))
			}
		}
	}
}

func CopyUnknownsFromObject(ctx context.Context, src, dest types.Object) types.Object {
	attributesSrc := src.Attributes()
	attributesDest := dest.Attributes()
	newAttributes := map[string]attr.Value{}
	for name, attr := range attributesDest {
		if attr.IsUnknown() {
			newAttributes[name] = attributesSrc[name]
			fmt.Printf("Copying unknown field: %s\n", name)
		} else {
			tfListDest, ok := attr.(types.List)
			if ok {
				tfListSrc := attributesSrc[name].(types.List)
				attr = CopyUnknownsFromList(ctx, tfListSrc, tfListDest)
			}
			tfObjectDest, ok := attr.(types.Object)
			if ok {
				tfObjectSrc := attributesSrc[name].(types.Object)
				newObject := CopyUnknownsFromObject(ctx, tfObjectSrc, tfObjectDest)
				attr = newObject
			}
			newAttributes[name] = attr
		}
	}
	new, diags := types.ObjectValue(src.AttributeTypes(ctx), newAttributes)
	if diags.HasError() {
		panic(fmt.Sprintf("Error converting object to model: %v", diags))
	}
	return new
}

func CopyUnknownsFromList(ctx context.Context, src, dest types.List) types.List {
	srcElements := src.Elements()
	count := len(srcElements)
	destElements := dest.Elements()
	if count != len(destElements) {
		return dest
	}
	new := make([]attr.Value, count)
	for i := range count {
		srcObj := srcElements[i].(types.Object)
		destObj := destElements[i].(types.Object)
		newObj := CopyUnknownsFromObject(ctx, srcObj, destObj)
		new[i] = newObj
	}
	return types.ListValueMust(dest.ElementType(ctx), new)
}
