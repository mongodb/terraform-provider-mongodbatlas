package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func GenerateTypedModels(attributes codespec.Attributes) CodeStatement {
	return generateTypedModels(attributes, "", false) // empty string for root model, results in TFModel
}

func generateTypedModels(attributes codespec.Attributes, name string, isNested bool) CodeStatement {
	models := []CodeStatement{generateStructOfTypedModel(attributes, name)}

	if isNested {
		models = append(models, generateModelObjType(attributes, name))
	}

	for i := range attributes {
		additionalModel := getNestedModel(&attributes[i])
		if additionalModel != nil {
			models = append(models, *additionalModel)
		}
	}

	return GroupCodeStatements(models, func(list []string) string { return strings.Join(list, "\n") })
}

func generateModelObjType(attrs codespec.Attributes, name string) CodeStatement {
	structProperties := []string{}
	for i := range attrs {
		propType := attrModelType(&attrs[i])
		prop := fmt.Sprintf(`%q: %sType,`, attrs[i].Name.SnakeCase(), propType)
		structProperties = append(structProperties, prop)
	}
	structPropsCode := strings.Join(structProperties, "\n")
	return CodeStatement{
		Code: fmt.Sprintf(`var %sObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
    %s
}}`, name, structPropsCode),
		Imports: []string{"github.com/hashicorp/terraform-plugin-framework/types", "github.com/hashicorp/terraform-plugin-framework/attr"},
	}
}

func getNestedModel(attribute *codespec.Attribute) *CodeStatement {
	var nested *codespec.NestedAttributeObject
	if attribute.ListNested != nil {
		nested = &attribute.ListNested.NestedObject
	}
	if attribute.SingleNested != nil {
		nested = &attribute.SingleNested.NestedObject
	}
	if attribute.MapNested != nil {
		nested = &attribute.MapNested.NestedObject
	}
	if attribute.SetNested != nil {
		nested = &attribute.SetNested.NestedObject
	}
	if nested == nil {
		return nil
	}
	res := generateTypedModels(nested.Attributes, attribute.Name.PascalCase(), true)
	return &res
}

func generateStructOfTypedModel(attributes codespec.Attributes, name string) CodeStatement {
	structProperties := []string{}
	for i := range attributes {
		structProperties = append(structProperties, typedModelProperty(&attributes[i]))
	}
	structPropsCode := strings.Join(structProperties, "\n")
	return CodeStatement{
		Code: fmt.Sprintf(`type TF%sModel struct {
			%s
		}`, name, structPropsCode),
		Imports: []string{"github.com/hashicorp/terraform-plugin-framework/types"},
	}
}

func typedModelProperty(attr *codespec.Attribute) string {
	namePascalCase := attr.Name.PascalCase()
	propType := attrModelType(attr)
	autogenerationTag := ""
	if attr.ReqBodyUsage == codespec.OmitAll {
		autogenerationTag = ` autogeneration:"omitjson"`
	} else if attr.ReqBodyUsage == codespec.PostBodyOnly {
		autogenerationTag = ` autogeneration:"createonly"`
	}
	return fmt.Sprintf("%s %s", namePascalCase, propType) + " `" + fmt.Sprintf("tfsdk:%q", attr.Name.SnakeCase()) + autogenerationTag + "`"
}

func attrModelType(attr *codespec.Attribute) string {
	switch {
	case attr.Float64 != nil:
		return "types.Float64"
	case attr.Bool != nil:
		return "types.Bool"
	case attr.String != nil:
		return "types.String"
	case attr.Number != nil:
		return "types.Number"
	case attr.Int64 != nil:
		return "types.Int64"
	case attr.Timeouts != nil:
		return "timeouts.Value"
	// For nested attributes the generic type is used, this is to ensure the model can store all possible values. e.g. nested computed only attributes need to be defined with TPF types to avoid error when getting the config.
	case attr.List != nil || attr.ListNested != nil:
		return "types.List"
	case attr.Set != nil || attr.SetNested != nil:
		return "types.Set"
	case attr.Map != nil || attr.MapNested != nil:
		return "types.Map"
	case attr.SingleNested != nil:
		return "types.Object"
	default:
		panic("Attribute with unknown type defined when generating typed model")
	}
}
