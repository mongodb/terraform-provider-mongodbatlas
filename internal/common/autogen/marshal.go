package autogen

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
)

// Marshal gets a Terraform model and marshals it into JSON (e.g. for an Atlas request).
// It supports the following types:
//   - Terraform types: String, Bool, Int64, Float64.
//   - Custom types: Object, Map, List, Set & jsontypes.Normalized.
//
// Attributes that are null or unknown are not marshaled by default.
// This behavior can be controlled via autogen tags (tags are exclusive):
//   - `omitjson`: Attribute is never marshaled.
//   - `omitjsonupdate`: Attribute is not marshaled if isUpdate is true.
//   - `sendnullasnullonupdate`: Attribute is marshaled as null if isUpdate is true.
//   - `sendnullasemptyonupdate`: Attribute is marshaled as empty if isUpdate is true (collections only).
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
		tags := GetPropertyTags(&attrTypeModel)
		if tags.OmitJSON {
			continue // skip fields with tag `omitjson`
		}
		if isUpdate && tags.OmitJSONUpdate {
			continue // skip fields with tag `omitjsonupdate` if in update mode
		}
		apiName := getAPINameFromTag(attrTypeModel.Name, tags)
		attrValModel := valModel.Field(i)
		if err := marshalAttr(apiName, attrValModel, objJSON, isUpdate, tags); err != nil {
			return nil, err
		}
	}
	return objJSON, nil
}

// getAPINameFromTag extracts the API name from the apiname tag if present (e.g., apiname:"groupId"),
// otherwise returns the model name uncapitalized as the default JSON name.
func getAPINameFromTag(modelName string, propertyTags PropertyTags) string {
	if propertyTags.APIName != nil {
		return *propertyTags.APIName
	}
	return stringcase.Uncapitalize(modelName)
}

func marshalAttr(attrNameJSON string, attrValModel reflect.Value, objJSON map[string]any, isUpdate bool, tags PropertyTags) error {
	obj, ok := attrValModel.Interface().(attr.Value)
	if !ok {
		panic("marshal expects only Terraform types in the model")
	}
	val, err := getModelAttr(obj, isUpdate)
	if err != nil {
		return err
	}

	// Emit empty collection on update for null list/set/map attributes when configured via sendNullAsEmptyOnUpdate
	if val == nil && isUpdate && tags.SendNullAsEmptyOnUpdate {
		switch obj.(type) {
		case customtypes.ListValueInterface, customtypes.NestedListValueInterface, customtypes.SetValueInterface, customtypes.NestedSetValueInterface:
			val = []any{}
		case customtypes.MapValueInterface, customtypes.NestedMapValueInterface:
			val = map[string]any{}
		}
	}

	// Emit value if non-nil, or emit null on update when configured via sendNullAsNullOnUpdate
	if val != nil || (isUpdate && tags.SendNullAsNullOnUpdate) {
		if tags.ListAsMap {
			val = ModifyJSONFromMapToList(val)
		}
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
