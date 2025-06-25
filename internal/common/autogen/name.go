package autogen

import (
	"github.com/huandu/xstrings"
)

func toJSONName(name string) string {
	return xstrings.ToCamelCase(name)
}

func toTfSchemaName(name string) string {
	return xstrings.ToSnakeCase(name)
}

func toTfModelName(name string) string {
	return xstrings.ToPascalCase(name)
}
