package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func GenerateTypedModels(attributes codespec.Attributes) CodeStatement {
	return generateTypedModels(attributes, "")
}

func generateTypedModels(attributes codespec.Attributes, name string) CodeStatement {
	models := []CodeStatement{generateStructOfTypedModel(attributes, name)}

	for i := range attributes {
		additionalModel := getNestedModel(&attributes[i], name)
		if additionalModel != nil {
			models = append(models, *additionalModel)
		}
	}

	return GroupCodeStatements(models, func(list []string) string { return strings.Join(list, "\n") })
}

func getNestedModel(attribute *codespec.Attribute, ancestorsName string) *CodeStatement {
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
	res := generateTypedModels(nested.Attributes, ancestorsName+attribute.TFModelName)
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
	var (
		propType   = attrModelType(attr)
		autogenTag string
	)
	switch attr.ReqBodyUsage {
	case codespec.AllRequestBodies:
		autogenTag = ""
	case codespec.OmitAlways:
		autogenTag = ` autogen:"omitjson"`
	case codespec.OmitInUpdateBody:
		autogenTag = ` autogen:"omitjsonupdate"`
	case codespec.IncludeNullOnUpdate:
		autogenTag = ` autogen:"includenullonupdate"`
	}
	return fmt.Sprintf("%s %s", attr.TFModelName, propType) + " `" + fmt.Sprintf("tfsdk:%q", attr.TFSchemaName.SnakeCase()) + autogenTag + "`"
}

func attrModelType(attr *codespec.Attribute) string {
	switch {
	case attr.CustomType != nil:
		return attr.CustomType.Model
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
	default:
		panic("Attribute with unknown type defined when generating typed model")
	}
}
