package schemafunc

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// HasUnknowns uses reflection to check if the object has any unknown fields
// Pass &TFModel{}
// Will only check the root level attributes
func HasUnknowns(obj any) bool {
	valObj := reflect.ValueOf(obj)
	if valObj.Kind() != reflect.Pointer {
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

// UnknownInConfig lists the top-level attributes that the configuration leaves unknown, to be passed as
// keepUnknown to CopyUnknowns. Planning their state value hides a change that the plan Terraform re-runs
// during apply then reveals, so the apply fails with "Provider produced inconsistent final plan". This is
// not covered by CopyUnknowns keeping the Optional-without-Computed attributes unknown, which leaves out the
// Optional+Computed ones. The name is top-level because an unknown nested in a list also hides whether that
// list changes.
func UnknownInConfig(config tftypes.Value) []string {
	attributes := map[string]tftypes.Value{}
	if config.IsNull() || !config.Type().Is(tftypes.Object{}) || config.As(&attributes) != nil {
		return nil
	}
	var names []string
	for name, value := range attributes {
		if !value.IsFullyKnown() {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	return names
}

// CopyUnknowns use reflection to copy unknown fields from src to dest.
// The implementation is similar to internal/common/conversion/model_generation.go#CopyModel
// attributes are the schema attributes of src and dest, see shouldKeepUnknown. Blocks are not traversed.
// keepUnknown is a list of fields that should not be copied, should always use the TF config name (snake_case)
// nestedStructMapping is a map of field names to their type: object, list. (`set` not implemented yet)
func CopyUnknowns(ctx context.Context, src, dest any, attributes map[string]schema.Attribute, keepUnknown []string) {
	validateKeepUnknown(keepUnknown)
	copyUnknowns(ctx, src, dest, attributes, keepUnknown)
}

// shouldKeepUnknown reports whether the destination value must stay unknown: it is in the keepUnknown list,
// or it is an unknown value of an attribute that is not Computed, whose value can only come from the
// configuration, so the state value would contradict the expression it resolves to. Known values are not
// kept, so an object is still traversed to reach its Computed children, such as pinned_fcv.version.
func shouldKeepUnknown(attributes map[string]schema.Attribute, keepUnknown []string, name string, destUnknown bool) bool {
	if slices.Contains(keepUnknown, name) {
		return true
	}
	attribute, found := attributes[name]
	return found && !attribute.IsComputed() && destUnknown
}

// nestedAttributes returns the attributes of a nested object, nil if the attribute has none.
// Only object and list are supported, the same nested types copyUnknowns traverses.
func nestedAttributes(attributes map[string]schema.Attribute, name string) map[string]schema.Attribute {
	switch attribute := attributes[name].(type) {
	case schema.SingleNestedAttribute:
		return attribute.Attributes
	case schema.ListNestedAttribute:
		return attribute.NestedObject.Attributes
	}
	return nil
}

func copyUnknowns(ctx context.Context, src, dest any, attributes map[string]schema.Attribute, keepUnknown []string) {
	valSrc, valDest := validateStructPointers(src, dest)
	typeSrc := valSrc.Type()
	typeDest := valDest.Type()
	for i := range typeDest.NumField() {
		fieldDest := typeDest.Field(i)
		name, tfName := fieldNameTFName(&fieldDest)
		srcValue := valSrc.FieldByName(name).Interface()
		if shouldKeepUnknown(attributes, keepUnknown, tfName, isUnknown(valDest.Field(i))) {
			continue
		}
		_, found := typeSrc.FieldByName(name)
		if !found || !valDest.Field(i).CanSet() {
			continue
		}
		nestedAttrs := nestedAttributes(attributes, tfName)
		nestedDest := valDest.FieldByName(name).Interface()
		objValueSrc, okSrc := srcValue.(types.Object)
		objValueDest, okDest := nestedDest.(types.Object)
		if okSrc && okDest {
			objValueNew := copyUnknownsFromObject(ctx, objValueSrc, objValueDest, nestedAttrs, keepUnknown)
			valDest.Field(i).Set(reflect.ValueOf(objValueNew))
			continue
		}
		listValueSrc, okSrc := srcValue.(types.List)
		listValueDest, okDest := nestedDest.(types.List)
		if okSrc && okDest {
			listValueNew := copyUnknownsFromList(ctx, listValueSrc, listValueDest, nestedAttrs, keepUnknown)
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
	if valSrc.Kind() != reflect.Pointer || valDest.Kind() != reflect.Pointer {
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

func copyUnknownsFromObject(ctx context.Context, src, dest types.Object, attributes map[string]schema.Attribute, keepUnknown []string) types.Object {
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
		if shouldKeepUnknown(attributes, keepUnknown, name, attr.IsUnknown()) {
			attributesMerged[name] = attr
			continue
		}
		nestedAttrs := nestedAttributes(attributes, name)
		tfListDest, isList := attr.(types.List)
		tfObjectDest, isObject := attr.(types.Object)
		if attr.IsUnknown() {
			tflog.Info(ctx, fmt.Sprintf("Copying unknown field: %s\n", name))
			switch {
			case isObject:
				attr = copyUnknownsFromObject(ctx, attributesSrc[name].(types.Object), tfObjectDest, nestedAttrs, keepUnknown)
			case isList:
				attr = copyUnknownsFromList(ctx, attributesSrc[name].(types.List), tfListDest, nestedAttrs, keepUnknown)
			default:
				attr = attributesSrc[name]
			}
			attributesMerged[name] = attr
			continue
		}
		if isList {
			tfListSrc := attributesSrc[name].(types.List)
			attr = copyUnknownsFromList(ctx, tfListSrc, tfListDest, nestedAttrs, keepUnknown)
		}
		if isObject {
			tfObjectSrc := attributesSrc[name].(types.Object)
			newObject := copyUnknownsFromObject(ctx, tfObjectSrc, tfObjectDest, nestedAttrs, keepUnknown)
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

func copyUnknownsFromList(ctx context.Context, src, dest types.List, attributes map[string]schema.Attribute, keepUnknown []string) types.List {
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
		newObj := copyUnknownsFromObject(ctx, srcObj, destObj, attributes, keepUnknown)
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
