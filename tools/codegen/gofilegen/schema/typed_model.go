package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func GenerateTypedModels(attributes codespec.Attributes) CodeStatement {
	return generateTypedModels(attributes, "", false)
}

// GenerateDataSourceTypedModels generates the TFDSModel struct for data sources.
// DS models are simpler: no autogen tags (no request body marshaling).
func GenerateDataSourceTypedModels(attributes codespec.Attributes) CodeStatement {
	return generateTypedModels(attributes, "DS", true)
}

func generateTypedModels(attributes codespec.Attributes, name string, isDataSource bool) CodeStatement {
	models := []CodeStatement{generateStructOfTypedModel(attributes, name, isDataSource)}

	// Generate nested models for all attributes (both resource and data source)
	// Data source nested models use the DS prefix to avoid naming clashes with resource models
	for i := range attributes {
		additionalModel := getNestedModel(&attributes[i], name, isDataSource)
		if additionalModel != nil {
			models = append(models, *additionalModel)
		}
	}

	return GroupCodeStatements(models, func(list []string) string { return strings.Join(list, "\n") })
}

func getNestedModel(attribute *codespec.Attribute, ancestorsName string, isDataSource bool) *CodeStatement {
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
	res := generateTypedModels(nested.Attributes, ancestorsName+attribute.TFModelName, isDataSource)
	return &res
}

func generateStructOfTypedModel(attributes codespec.Attributes, name string, isDataSource bool) CodeStatement {
	structProperties := []string{}
	for i := range attributes {
		structProperties = append(structProperties, typedModelProperty(&attributes[i], isDataSource))
	}
	structPropsCode := strings.Join(structProperties, "\n")
	return CodeStatement{
		Code: fmt.Sprintf(`type TF%sModel struct {
			%s
		}`, name, structPropsCode),
		Imports: []string{"github.com/hashicorp/terraform-plugin-framework/types"},
	}
}

func typedModelProperty(attr *codespec.Attribute, isDataSource bool) string {
	propType := attrModelType(attr, isDataSource)

	// Resource models need additional tags for marshaling
	var (
		tagsStr     = ""
		autogenTags = make([]string, 0)
		apinameTag  = ""
	)

	// Add apiname tag if the API name is different from the uncapitalized model name
	// This ensures correct marshaling/unmarshaling when TFModelName doesn't derive to the correct API property name
	if attr.APIName != "" && attr.APIName != stringcase.Uncapitalize(attr.TFModelName) {
		apinameTag = fmt.Sprintf(" apiname:%q", attr.APIName)
	}

	if attr.Sensitive {
		autogenTags = append(autogenTags, "sensitive")
	}

	switch attr.ReqBodyUsage {
	case codespec.AllRequestBodies:
	case codespec.OmitAlways:
		autogenTags = append(autogenTags, "omitjson")
	case codespec.OmitInUpdateBody:
		autogenTags = append(autogenTags, "omitjsonupdate")
	case codespec.IncludeNullOnUpdate:
		autogenTags = append(autogenTags, "includenullonupdate")
	}

	if len(autogenTags) > 0 {
		tagsStr = fmt.Sprintf(" autogen:%q", strings.Join(autogenTags, ","))
	}
	return fmt.Sprintf("%s %s", attr.TFModelName, propType) + " `" + fmt.Sprintf("tfsdk:%q", attr.TFSchemaName) + apinameTag + tagsStr + "`"
}

func attrModelType(attr *codespec.Attribute, isDataSource bool) string {
	switch {
	case attr.CustomType != nil:
		model := attr.CustomType.Model
		if isDataSource {
			model = addDSPrefixToNestedModels(model)
		}
		return model
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

// addDSPrefixToNestedModels transforms nested model references by adding DS prefix.
// e.g., "customtypes.ObjectValue[TFNestedObjectAttrModel]" -> "customtypes.ObjectValue[TFDSNestedObjectAttrModel]"
// This only applies to nested models (TF*Model pattern), not primitive types like types.String.
func addDSPrefixToNestedModels(s string) string {
	// Pattern: [TF followed by word characters and ending with Model]
	// We only want to add DS prefix to nested models, not change types.String etc.
	return strings.ReplaceAll(s, "[TF", "[TFDS")
}
