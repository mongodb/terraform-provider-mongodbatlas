package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func GenerateTypedModels(attributes codespec.Attributes, withObjTypes bool) CodeStatement {
	return generateTypedModels(attributes, "", false, withObjTypes)
}

func generateTypedModels(attributes codespec.Attributes, name string, isNested, withObjTypes bool) CodeStatement {
	models := []CodeStatement{generateStructOfTypedModel(attributes, name)}

	if isNested && withObjTypes {
		models = append(models, generateModelObjType(attributes, name))
	}

	for i := range attributes {
		additionalModel := getNestedModel(&attributes[i], name, withObjTypes)
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
		comment := limitationComment(&attrs[i])
		prop := fmt.Sprintf(`%q: %sType,%s`, attrs[i].Name.SnakeCase(), propType, comment)
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

func getNestedModel(attribute *codespec.Attribute, ancestorsName string, withObjTypes bool) *CodeStatement {
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
	res := generateTypedModels(nested.Attributes, ancestorsName+attribute.Name.PascalCase(), true, withObjTypes)
	return &res
}

func generateStructOfTypedModel(attributes codespec.Attributes, name string) CodeStatement {
	structProperties := []string{}
	for i := range attributes {
		structProperties = append(structProperties, typedModelProperty(&attributes[i]))
	}
	structPropsCode := strings.Join(structProperties, "\n")
	tfModelName := "TF" + name + "Model"
	discriminatorAttrFunc := generateDiscriminatorAttrFunc(tfModelName, attributes)
	return CodeStatement{
		Code: fmt.Sprintf(`type %s struct {
			%s
		}`, tfModelName, structPropsCode) + discriminatorAttrFunc,
		Imports: []string{"github.com/hashicorp/terraform-plugin-framework/types"},
	}
}

func generateDiscriminatorAttrFunc(tfModelName string, attributes codespec.Attributes) string {
	/*
		Example of stream_connection:
			func (s *streamConn) DiscriminatorAttr(objJSON map[string]any) string {
				// Probably can return a single attribute
				switch objJSON["type"] {
				case "Cluster":
					return "TypeCluster"
				case "Https":
					return "TypeHttps"
				}
				return ""
			}
	*/
	discriminators := []codespec.DiscriminatorMapping{}
	for i := range attributes {
		attr := &attributes[i]
		if attr.Discriminator != nil {
			discriminators = append(discriminators, *attr.Discriminator)
		}
	}
	if len(discriminators) == 0 {
		return ""
	}
	discriminatorCases := []string{}
	discriminatorPropName := ""
	for _, disc := range discriminators {
		discriminatorCases = append(discriminatorCases, fmt.Sprintf(`case %q: return %q`, disc.DiscriminatorValue, disc.AttributeName()))
		if discriminatorPropName == "" {
			discriminatorPropName = disc.DiscriminatorProperty
		}
		if discriminatorPropName != disc.DiscriminatorProperty {
			panic(fmt.Sprintf("Discriminator mappings must have the same DiscriminatorProperty, found %s and %s for model %s", discriminatorPropName, disc.DiscriminatorProperty, tfModelName))
		}
	}
	return fmt.Sprintf(`
func (t *%s) DiscriminatorAttr(objJSON map[string]any) string {
	switch objJSON["%s"] {
	%s
	}
	return ""
}`, tfModelName, discriminatorPropName, strings.Join(discriminatorCases, "\n\t"))
}

func typedModelProperty(attr *codespec.Attribute) string {
	var (
		namePascalCase = attr.Name.PascalCase()
		propType       = attrModelType(attr)
		autogenTag     string
	)
	switch attr.ReqBodyUsage {
	case codespec.OmitAlways:
		autogenTag = ` autogen:"omitjson"`
	case codespec.OmitInUpdateBody:
		autogenTag = ` autogen:"omitjsonupdate"`
	}
	if discriminator := attr.Discriminator; discriminator != nil {
		autogenTag = fmt.Sprintf(` autogen:"discriminator:%s=%s"`, discriminator.DiscriminatorProperty, discriminator.DiscriminatorValue)
	}
	return fmt.Sprintf("%s %s", namePascalCase, propType) + " `" + fmt.Sprintf("tfsdk:%q", attr.Name.SnakeCase()) + autogenTag + "`"
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

func limitationComment(attr *codespec.Attribute) string {
	if attr.List == nil && attr.ListNested == nil &&
		attr.Map == nil && attr.MapNested == nil &&
		attr.Set == nil && attr.SetNested == nil {
		return ""
	}
	return " // TODO: missing ElemType, codegen limitation tracked in CLOUDP-311105"
}
