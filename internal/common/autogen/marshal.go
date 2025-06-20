package autogen

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/huandu/xstrings"
)

const (
	tagKey               = "autogen"
	tagValOmitJSON       = "omitjson"
	tagValOmitJSONUpdate = "omitjsonupdate"
)

// Marshal gets a Terraform model and marshals it into JSON (e.g. for an Atlas request).
// It supports the following Terraform model types: String, Bool, Int64, Float64, Object, Map, List, Set.
// Attributes that are null or unknown are not marshaled.
// Attributes with autogen tag `omitjson` are never marshaled, this only applies to the root model.
// Attributes with autogen tag `omitjsonupdate` are not marshaled if isUpdate is true, this only applies to the root model.
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
	for i := 0; i < valModel.NumField(); i++ {
		attrTypeModel := valModel.Type().Field(i)
		tag := attrTypeModel.Tag.Get(tagKey)
		if tag == tagValOmitJSON {
			continue // skip fields with tag `omitjson`
		}
		if isUpdate && tag == tagValOmitJSONUpdate {
			continue // skip fields with tag `omitjsonupdate` if in update mode
		}
		attrNameModel := attrTypeModel.Name
		attrValModel := valModel.Field(i)
		if err := marshalAttr(attrNameModel, attrValModel, objJSON, isUpdate); err != nil {
			return nil, err
		}
	}
	return objJSON, nil
}

func marshalAttr(attrNameModel string, attrValModel reflect.Value, objJSON map[string]any, isUpdate bool) error {
	attrNameJSON := xstrings.ToCamelCase(attrNameModel)
	obj, ok := attrValModel.Interface().(attr.Value)
	if !ok {
		panic("marshal expects only Terraform types in the model")
	}
	val, err := getModelAttr(obj)
	if err != nil {
		return err
	}

	// What to send in update request body if the root field is null, nothing is sent by default
	if val == nil && isUpdate {
		switch obj.(type) {
		case types.List, types.Set:
			val = []any{} // Send an empty array if it's a list or set
		}
	}

	if val != nil {
		objJSON[attrNameJSON] = val
	}
	return nil
}

func getModelAttr(val attr.Value) (any, error) {
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
	case types.Object:
		return getMapAttr(v.Attributes(), false)
	case types.Map:
		return getMapAttr(v.Elements(), true)
	case types.List:
		return getListAttr(v.Elements())
	case types.Set:
		return getListAttr(v.Elements())
	default:
		return nil, fmt.Errorf("unmarshal not supported yet for type %T", v)
	}
}

func getListAttr(elms []attr.Value) (any, error) {
	arr := make([]any, 0)
	for _, attr := range elms {
		valChild, err := getModelAttr(attr)
		if err != nil {
			return nil, err
		}
		if valChild != nil {
			arr = append(arr, valChild)
		}
	}
	return arr, nil
}

// getMapAttr gets a map of attributes and returns a map of JSON attributes.
// keepKeyCase is used for types.Map to keep key case. However, we want to use JSON key case for types.Object
func getMapAttr(elms map[string]attr.Value, keepKeyCase bool) (any, error) {
	objJSON := make(map[string]any)
	for name, attr := range elms {
		valChild, err := getModelAttr(attr)
		if err != nil {
			return nil, err
		}
		if valChild != nil {
			nameJSON := xstrings.ToCamelCase(name)
			if keepKeyCase {
				nameJSON = name
			}
			objJSON[nameJSON] = valChild
		}
	}
	return objJSON, nil
}
