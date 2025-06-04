package autogen

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/huandu/xstrings"
)

const (
	tagKey               = "autogen"
	tagValOmitJSON       = "omitjson"
	tagValOmitJSONUpdate = "omitjsonupdate"
)

type DiscriminatorTag struct {
	DiscriminatorPropName  string
	DiscriminatorPropValue string
}

// / IsDiscriminatorTag decodes a autogen:"discriminator:type=Cluster" tag to a DiscriminatorTag struct.
func IsDiscriminatorTag(tag string) *DiscriminatorTag {
	if tag == "" || !strings.HasPrefix(tag, "discriminator:") {
		return nil
	}
	// decode the tag
	keyValue := strings.TrimPrefix(tag, "discriminator:")
	propName, propValue, found := strings.Cut(keyValue, "=")
	if !found {
		return nil // not a valid discriminator tag
	}
	return &DiscriminatorTag{
		DiscriminatorPropName:  propName,
		DiscriminatorPropValue: propValue,
	}
}

// Marshal gets a Terraform model and marshals it into JSON (e.g. for an Atlas request).
// It supports the following Terraform model types: String, Bool, Int64, Float64, Object, Map, List, Set.
// Attributes that are null or unknown are not marshaled.
// Attributes with autogen tag `omitjson` are never marshaled, this only applies to the root model.
// Attributes with autogen tag `omitjsonupdate` are not marshaled if isUpdate is true, this only applies to the root model.
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
		if err := marshalAttr(attrNameModel, attrValModel, objJSON); err != nil {
			return nil, err
		}
		if err := handleDiscriminator(tag, attrNameModel, objJSON); err != nil {
			return nil, err
		}
	}
	return objJSON, nil
}

func handleDiscriminator(tag string, attrNameModel string, objJSON map[string]any) error {
	discriminatorTag := IsDiscriminatorTag(tag)
	if discriminatorTag == nil {
		return nil // not a discriminator tag, nothing to do
	}
	attrNameJSON := xstrings.ToCamelCase(attrNameModel)
	attrValJSON, ok := objJSON[attrNameJSON]
	if !ok {
		return nil // attribute not found in the JSON, nothing to do (probably Null or Unknown value)
	}
	// if the discriminator is set, we remove the attribute from the JSON
	delete(objJSON, attrNameJSON)
	objJSON[discriminatorTag.DiscriminatorPropName] = discriminatorTag.DiscriminatorPropValue // set the discriminator property in the JSON
	attrValJSONObject, ok := attrValJSON.(map[string]any)
	if !ok {
		return fmt.Errorf("discriminator attribute %s must be an object", attrNameJSON)
	}
	maps.Copy(objJSON, attrValJSONObject) // copy the object attributes to the root JSON
	return nil
}

func marshalAttr(attrNameModel string, attrValModel reflect.Value, objJSON map[string]any) error {
	attrNameJSON := xstrings.ToCamelCase(attrNameModel)
	obj, ok := attrValModel.Interface().(attr.Value)
	if !ok {
		panic("marshal expects only Terraform types in the model")
	}
	val, err := getModelAttr(obj)
	if err != nil {
		return err
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
