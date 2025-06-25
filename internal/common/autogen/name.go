package autogen

import (
	"strings"

	"github.com/huandu/xstrings"
)

// SanitizeTfAttrName returns a valid Terraform attribute name.
func SanitizeTfAttrName(name string) string {
	// Names starting with _ can't be used in Terraform attributes, e.g. _id in search_index_api is converted to id.
	return strings.TrimPrefix(name, "_")
}

func ToJSONName(name string) string {
	return xstrings.ToCamelCase(name)
}

func ToTfSchemaName(name string) string {
	return SanitizeTfAttrName(xstrings.ToSnakeCase(name))
}

func ToTfModelName(name string) string {
	return SanitizeTfAttrName(xstrings.ToPascalCase(name))
}
