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

// HasUnknowns uses reflection to check if the object has any unknown fields
// Pass &TFModel{}
// Will only check the root level attributes
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
		if isUnknown(field) {
			return true
		}
	}
	return false
}

// CopyUnknowns use reflection to copy unknown fields from src to dest.
// The implementation is similar to internal/common/conversion/model_generation.go#CopyModel
// keepUnknown is a list of fields that should not be copied, should always use the TF config name (snake_case)
// nestedStructMapping is a map of field names to their type: object, list. (`set` not implemented yet)
func CopyUnknowns(ctx context.Context, src, dest any, keepUnknown []string, keepUnknownCall ...func(string, attr.Value) bool) {
	validateKeepUnknown(keepUnknown)
	slicesContains := func(name string, value attr.Value) bool {
		return slices.Contains(keepUnknown, name)
	}
	call := combineKeepUnknownCalls(
		combineKeepUnknownCalls(keepUnknownCall...),
		slicesContains,
	)
	copyUnknowns(ctx, src, dest, call)
}

func combineKeepUnknownCalls(calls ...func(string, attr.Value) bool) func(string, attr.Value) bool {
	return func(name string, value attr.Value) bool {
		for _, call := range calls {
			if call(name, value) {
				return true
			}
		}
		return false
	}
}

func copyUnknowns(ctx context.Context, src, dest any, keepUnknownCall func(string, attr.Value) bool) {
	valSrc, valDest := validateStructPointers(src, dest)
	typeSrc := valSrc.Type()
	typeDest := valDest.Type()
	for i := range typeDest.NumField() {
		fieldDest := typeDest.Field(i)
		name, tfName := fieldNameTFName(&fieldDest)
		srcValue := valSrc.FieldByName(name).Interface()
		if keepUnknownCall(tfName, srcValue.(attr.Value)) {
			continue
		}
		_, found := typeSrc.FieldByName(name)
		if !found || !valDest.Field(i).CanSet() {
			continue
		}
		nestedDest := valDest.FieldByName(name).Interface()
		objValueSrc, okSrc := srcValue.(types.Object)
		objValueDest, okDest := nestedDest.(types.Object)
		if okSrc && okDest {
			objValueNew := copyUnknownsFromObject(ctx, objValueSrc, objValueDest, keepUnknownCall)
			valDest.Field(i).Set(reflect.ValueOf(objValueNew))
			continue
		}
		listValueSrc, okSrc := srcValue.(types.List)
		listValueDest, okDest := nestedDest.(types.List)
		if okSrc && okDest {
			listValueNew := copyUnknownsFromList(ctx, listValueSrc, listValueDest, keepUnknownCall)
			valDest.Field(i).Set(reflect.ValueOf(listValueNew))
			continue
		}
		if isUnknown(valDest.Field(i)) {
			tflog.Info(ctx, fmt.Sprintf("Copying unknown field: %s\n", name))
			valDest.Field(i).Set(valSrc.FieldByName(name))
			continue
		}
	}
}

func fieldNameTFName(fieldDest *reflect.StructField) (name, tfName string) {
	name = fieldDest.Name
	tfName = fieldDest.Tag.Get("tfsdk")
	if tfName == "" {
		panic(fmt.Sprintf("field %s has no tfsdk tag", name))
	}
	return name, tfName
}

func validateStructPointers(src, dest any) (reflectSrc, reflectDest reflect.Value) {
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
	return valSrc, valDest
}

func isUnknown(obj reflect.Value) bool {
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

func copyUnknownsFromObject(ctx context.Context, src, dest types.Object, keepUnknownCall func(string, attr.Value) bool) types.Object {
	// if something is null in the state and unknown in plan, we expect it to remain null
	if src.IsNull() && dest.IsUnknown() {
		return src
	}
	// if state is null we have nothing to copy, if plan is null we shouldn't copy anything
	if src.IsNull() || dest.IsNull() {
		return dest
	}
	attributesSrc := src.Attributes()
	attributesDest := dest.Attributes()
	attributesMerged := map[string]attr.Value{}
	if dest.IsUnknown() {
		// an unknown object will have emptyAttributes, to support keep unknowns on unknown objects we use fillUnknowns to populate the attributes
		attributesDest = fillUnknowns(ctx, attributesSrc)
	}
	for name, attr := range attributesDest {
		replacement := attributesSrc[name]
		if keepUnknownCall(name, replacement) {
			attributesMerged[name] = attr
			continue
		}
		tfListDest, isList := attr.(types.List)
		tfObjectDest, isObject := attr.(types.Object)
		if attr.IsUnknown() {
			tflog.Info(ctx, fmt.Sprintf("Copying unknown field: %s\n", name))
			switch {
			case isObject:
				attr = copyUnknownsFromObject(ctx, attributesSrc[name].(types.Object), tfObjectDest, keepUnknownCall)
			case isList:
				attr = copyUnknownsFromList(ctx, attributesSrc[name].(types.List), tfListDest, keepUnknownCall)
			default:
				attr = attributesSrc[name]
			}
			attributesMerged[name] = attr
			continue
		}
		if isList {
			tfListSrc := attributesSrc[name].(types.List)
			attr = copyUnknownsFromList(ctx, tfListSrc, tfListDest, keepUnknownCall)
		}
		if isObject {
			tfObjectSrc := attributesSrc[name].(types.Object)
			newObject := copyUnknownsFromObject(ctx, tfObjectSrc, tfObjectDest, keepUnknownCall)
			attr = newObject
		}
		attributesMerged[name] = attr
	}
	merged, diags := types.ObjectValue(src.AttributeTypes(ctx), attributesMerged)
	if diags.HasError() {
		panic(fmt.Sprintf("Error converting object to model: %v", diags))
	}
	return merged
}

// fillUnknowns creates a new map with all the attributes as unknown
func fillUnknowns(ctx context.Context, attributesSrc map[string]attr.Value) map[string]attr.Value {
	unknownAttributes := map[string]attr.Value{}
	for name, attrSrc := range attributesSrc {
		unknownAttributes[name] = asUnknownValue(ctx, attrSrc)
	}
	return unknownAttributes
}

func copyUnknownsFromList(ctx context.Context, src, dest types.List, keepUnknownCall func(string, attr.Value) bool) types.List {
	srcElements := src.Elements()
	destElements := dest.Elements()
	count := len(srcElements)
	if count != len(destElements) || src.IsNull() || dest.IsNull() {
		return dest
	}
	merged := make([]attr.Value, count)
	for i := range count {
		srcObj := srcElements[i].(types.Object)
		destObj := destElements[i].(types.Object)
		newObj := copyUnknownsFromObject(ctx, srcObj, destObj, keepUnknownCall)
		merged[i] = newObj
	}
	return types.ListValueMust(dest.ElementType(ctx), merged)
}

// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/types@v1.13.0
func asUnknownValue(ctx context.Context, value attr.Value) attr.Value {
	switch v := value.(type) {
	case types.List:
		return types.ListUnknown(v.ElementType(ctx))
	case types.Object:
		return types.ObjectUnknown(v.AttributeTypes(ctx))
	case types.Map:
		return types.MapUnknown(v.ElementType(ctx))
	case types.Set:
		return types.SetUnknown(v.ElementType(ctx))
	case types.Tuple:
		return types.TupleUnknown(v.ElementTypes(ctx))
	case types.String:
		return types.StringUnknown()
	case types.Bool:
		return types.BoolUnknown()
	case types.Int64:
		return types.Int64Unknown()
	case types.Int32:
		return types.Int32Unknown()
	case types.Float64:
		return types.Float64Unknown()
	case types.Float32:
		return types.Float32Unknown()
	case types.Number:
		return types.NumberUnknown()
	case types.Dynamic:
		return types.DynamicUnknown()
	}
	panic(fmt.Sprintf("Unknown value to create unknown: %v", value))
}
