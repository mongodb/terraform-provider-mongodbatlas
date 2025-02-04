package schemafunc

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"

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

func validateKeepUnknown(keepUnknown []string) {
	invalidNames := []string{}
	for _, name := range keepUnknown {
		if strings.ToLower(name) != name {
			invalidNames = append(invalidNames, name)
		}
	}
	if len(invalidNames) > 0 {
		panic(fmt.Sprintf("keepUnknown names must be lowercase and use TF config format: %v", invalidNames))
	}
}

// CopyUnknowns use reflection to copy unknown fields from src to dest.
// The implementation is similar to internal/common/conversion/model_generation.go#CopyModel
// keepUnknown is a list of fields that should not be copied, should always use the TF config name
// nestedStructMapping is a map of field names to their type: object, list. (set not implemented yet)
func CopyUnknowns(ctx context.Context, src, dest any, keepUnknown []string, nestedStructMapping map[string]string) {
	validateKeepUnknown(keepUnknown)
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
		tfName := fieldDest.Tag.Get("tfsdk")
		if tfName == "" {
			panic(fmt.Sprintf("field %s has no tfsdk tag", name))
		}
		if slices.Contains(keepUnknown, tfName) {
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
			tflog.Info(ctx, fmt.Sprintf("Processing nested field: %s with type %s\n", name, nestedType))
			nestedSrc := valSrc.FieldByName(name).Interface()
			nestedDest := valDest.FieldByName(name).Interface()
			if nestedType == "object" {
				objValueSrc := nestedSrc.(types.Object)
				objValueDest := nestedDest.(types.Object)
				objValueNew := copyUnknownsFromObject(ctx, objValueSrc, objValueDest, keepUnknown)
				valDest.Field(i).Set(reflect.ValueOf(objValueNew))
			} else if nestedType == "list" {
				listValueSrc := nestedSrc.(types.List)
				listValueDest := nestedDest.(types.List)
				listValueNew := copyUnknownsFromList(ctx, listValueSrc, listValueDest, keepUnknown)
				valDest.Field(i).Set(reflect.ValueOf(listValueNew))
			} else {
				panic(fmt.Sprintf("nested type not supported yet: %s", nestedType))
			}
		}
	}
}

func copyUnknownsFromObject(ctx context.Context, src, dest types.Object, keepUnknown []string) types.Object {
	attributesSrc := src.Attributes()
	attributesDest := dest.Attributes()
	newAttributes := map[string]attr.Value{}
	for name, attr := range attributesDest {
		if attr.IsUnknown() {
			newAttributes[name] = attributesSrc[name]
			tflog.Info(ctx, fmt.Sprintf("Copying unknown field: %s\n", name))
		} else {
			tfListDest, ok := attr.(types.List)
			if ok {
				tfListSrc := attributesSrc[name].(types.List)
				attr = copyUnknownsFromList(ctx, tfListSrc, tfListDest, keepUnknown)
			}
			tfObjectDest, ok := attr.(types.Object)
			if ok {
				tfObjectSrc := attributesSrc[name].(types.Object)
				newObject := copyUnknownsFromObject(ctx, tfObjectSrc, tfObjectDest, keepUnknown)
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

func copyUnknownsFromList(ctx context.Context, src, dest types.List, keepUnknown []string) types.List {
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
		newObj := copyUnknownsFromObject(ctx, srcObj, destObj, keepUnknown)
		new[i] = newObj
	}
	return types.ListValueMust(dest.ElementType(ctx), new)
}
