package autogen

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
)

const (
	tagKey                    = "autogen"
	tagSensitive              = "sensitive"
	tagValOmitJSON            = "omitjson"
	tagValOmitJSONUpdate      = "omitjsonupdate"
	tagValIncludeNullOnUpdate = "includenullonupdate"
)

// Marshal gets a Terraform model and marshals it into JSON (e.g. for an Atlas request).
// It supports the following types:
//   - Terraform types: String, Bool, Int64, Float64.
//   - Custom types: Object, Map, List, Set & jsontypes.Normalized.
//
// Attributes that are null or unknown are not marshaled.
// Attributes with autogen tag `omitjson` are never marshaled, this only applies to the root model.
// Attributes with autogen tag `omitjsonupdate` are not marshaled if isUpdate is true, this only applies to the root model.
// Attributes with autogen tag `includenullonupdate` are marshaled if isUpdate is true (even if null), this only applies to the root model.
// Null list or set root elements are sent as empty arrays if isUpdate is true.
func Marshal(model any, isUpdate bool) ([]byte, error) {
	valModel := reflect.ValueOf(model)
	if valModel.Kind() != reflect.Ptr {
		panic("model must be pointer")
	}
	valModel = valModel.Elem()
	if valModel.Kind() != reflect.Struct {
		panic("model must be pointer to struct")
	}
	objJSON, err := marshalAttrs(valModel, isUpdate)
	if err != nil {
		return nil, err
	}
	return json.Marshal(objJSON)
}

func marshalAttrs(valModel reflect.Value, isUpdate bool) (map[string]any, error) {
	objJSON := make(map[string]any)
	for i := range valModel.NumField() {
		attrTypeModel := valModel.Type().Field(i)
		tags := strings.Split(attrTypeModel.Tag.Get(tagKey), ",")
		if slices.Contains(tags, tagValOmitJSON) {
			continue // skip fields with tag `omitjson`
		}
		if isUpdate && slices.Contains(tags, tagValOmitJSONUpdate) {
			continue // skip fields with tag `omitjsonupdate` if in update mode
		}
		attrNameModel := attrTypeModel.Name
		attrValModel := valModel.Field(i)
		includeNullOnUpdate := slices.Contains(tags, tagValIncludeNullOnUpdate)
		if err := marshalAttr(attrNameModel, attrValModel, objJSON, isUpdate, includeNullOnUpdate); err != nil {
			return nil, err
		}
	}
	return objJSON, nil
}

func marshalAttr(attrNameModel string, attrValModel reflect.Value, objJSON map[string]any, isUpdate, includeNullOnUpdate bool) error {
	attrNameJSON := stringcase.Uncapitalize(attrNameModel)
	obj, ok := attrValModel.Interface().(attr.Value)
	if !ok {
		panic("marshal expects only Terraform types in the model")
	}
	val, err := getModelAttr(obj, isUpdate)
	if err != nil {
		return err
	}

	if val == nil && isUpdate {
		switch obj.(type) {
		case customtypes.ListValueInterface, customtypes.NestedListValueInterface, customtypes.SetValueInterface, customtypes.NestedSetValueInterface:
			val = []any{} // Send an empty array if it's a null root list or set
		}
	}

	if val != nil || (isUpdate && includeNullOnUpdate) {
		objJSON[attrNameJSON] = val
	}
	return nil
}

func getModelAttr(val attr.Value, isUpdate bool) (any, error) {
	if val.IsNull() || val.IsUnknown() {
		return nil, nil // skip null or unknown values
	}
	switch v := val.(type) {
	case types.String:
		return v.ValueString(), nil
	case types.Bool:
		return v.ValueBool(), nil
	case types.Int64:
		return v.ValueInt64(), nil
	case types.Float64:
		return v.ValueFloat64(), nil
	case jsontypes.Normalized:
		var valueJSON any
		if err := json.Unmarshal([]byte(v.ValueString()), &valueJSON); err != nil {
			return nil, fmt.Errorf("marshal failed for JSON custom type: %v", err)
		}
		return valueJSON, nil
	case customtypes.ListValueInterface:
		return getListAttr(v.Elements(), isUpdate)
	case customtypes.SetValueInterface:
		return getListAttr(v.Elements(), isUpdate)
	case customtypes.MapValueInterface:
		return getMapAttr(v.Elements(), isUpdate)
	case customtypes.ObjectValueInterface:
		valuePtr, diags := v.ValuePtrAsAny(context.Background())
		if diags.HasError() {
			return nil, fmt.Errorf("marshal failed for type: %v", diags)
		}

		result, err := marshalAttrs(reflect.ValueOf(valuePtr).Elem(), isUpdate)
		return result, err
	case customtypes.NestedListValueInterface:
		slicePtr, diags := v.SlicePtrAsAny(context.Background())
		if diags.HasError() {
			return nil, fmt.Errorf("marshal failed for type: %v", diags)
		}

		return getNestedSliceAttr(slicePtr, isUpdate)
	case customtypes.NestedSetValueInterface:
		slicePtr, diags := v.SlicePtrAsAny(context.Background())
		if diags.HasError() {
			return nil, fmt.Errorf("marshal failed for type: %v", diags)
		}

		return getNestedSliceAttr(slicePtr, isUpdate)
	case customtypes.NestedMapValueInterface:
		mapPtr, diags := v.MapPtrAsAny(context.Background())
		if diags.HasError() {
			return nil, fmt.Errorf("marshal failed for type: %v", diags)
		}

		return getNestedMapAttr(mapPtr, isUpdate)
	default:
		return nil, fmt.Errorf("marshal not supported yet for type %T", v)
	}
}

func getListAttr(elms []attr.Value, isUpdate bool) (any, error) {
	slice := make([]any, 0, len(elms))
	for _, attr := range elms {
		value, err := getModelAttr(attr, isUpdate)
		if err != nil {
			return nil, err
		}
		if value != nil {
			slice = append(slice, value)
		}
	}
	return slice, nil
}

func getMapAttr(elms map[string]attr.Value, isUpdate bool) (any, error) {
	objJSON := make(map[string]any)
	for name, attr := range elms {
		value, err := getModelAttr(attr, isUpdate)
		if err != nil {
			return nil, err
		}
		if value != nil {
			objJSON[name] = value
		}
	}
	return objJSON, nil
}

func getNestedSliceAttr(slicePtr any, isUpdate bool) (any, error) {
	sliceValue := reflect.ValueOf(slicePtr).Elem()
	length := sliceValue.Len()

	result := make([]any, 0, length)
	for i := range length {
		value, err := marshalAttrs(sliceValue.Index(i), isUpdate)
		if err != nil {
			return nil, err
		}
		if value != nil {
			result = append(result, value)
		}
	}

	return result, nil
}

func getNestedMapAttr(mapPtr any, isUpdate bool) (any, error) {
	mapValue := reflect.ValueOf(mapPtr).Elem()

	result := make(map[string]any, mapValue.Len())
	iter := mapValue.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		value, err := marshalAttrs(iter.Value(), isUpdate)
		if err != nil {
			return nil, err
		}
		if value != nil {
			result[key] = value
		}
	}

	return result, nil
}
