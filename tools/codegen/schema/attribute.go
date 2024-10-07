package schema

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func RenderAttributes(attrs codespec.Attributes) []string {
	result := []string{}
	for _, attr := range attrs {
		result = append(result, attribute(attr))
	}

	return result
}

func attribute(attr codespec.Attribute) string {
	name := attr.Name
	generalProperties := renderCommonProperties(attr)
	attrType := "schema.string" // TODO fixed to string

	return fmt.Sprintf(`
	"%s": %s{
		%s
	},
	`, name, attrType, generalProperties)
}

func renderCommonProperties(attr codespec.Attribute) string {
	var result string
	if attr.IsComputed != nil && *attr.IsComputed {
		result = result + "Computed: true,\n"
	}
	if attr.IsRequired != nil && *attr.IsRequired {
		result = result + "Required: true,\n"
	}
	if attr.IsOptional != nil && *attr.IsOptional {
		result = result + "Optional: true,\n"
	}
	if attr.Description != nil {
		result = result + fmt.Sprintf("MarkdownDescription: %q,\n", *attr.Description)
	}
	if attr.Sensitive != nil && *attr.Sensitive {
		result = result + "Sensitive: true,\n"
	}
	return result
}
