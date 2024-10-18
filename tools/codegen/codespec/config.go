package codespec

import (
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
)

func applyConfigSchemaOptions(resourceConfig *config.Resource, resource *Resource) {
	applySchemaOptions(resourceConfig.SchemaOptions, &resource.Schema.Attributes, "")
}

func applySchemaOptions(schemaOptions config.SchemaOptions, attributes *Attributes, parentName string) {
	ignoredAttrs := getIgnoredAttributesMap(schemaOptions.Ignores)

	var finalAttributes Attributes

	for i := range *attributes {
		attr := &(*attributes)[i]
		attrPathName := getAttributePathName(string(attr.Name), parentName)

		if shouldIgnoreAttribute(attrPathName, ignoredAttrs) {
			continue
		}

		// the config is expected to use alias name for defining any subsequent overrides (description, etc)
		applyAlias(attr, &attrPathName, schemaOptions)

		applyOverrides(attr, attrPathName, schemaOptions)

		processNestedAttributes(attr, schemaOptions, attrPathName)

		finalAttributes = append(finalAttributes, *attr)
	}

	*attributes = finalAttributes
}

func getAttributePathName(attrName, parentName string) string {
	if parentName == "" {
		return attrName
	}
	return parentName + "." + attrName
}

func getIgnoredAttributesMap(ignores []string) map[string]bool {
	ignoredAttrs := make(map[string]bool)
	for _, ignoredAttr := range ignores {
		ignoredAttrs[ignoredAttr] = true
	}
	return ignoredAttrs
}

func shouldIgnoreAttribute(attrName string, ignoredAttrs map[string]bool) bool {
	return ignoredAttrs[attrName]
}

func applyAlias(attr *Attribute, attrPathName *string, schemaOptions config.SchemaOptions) {
	parts := strings.Split(*attrPathName, ".")

	for i := range parts {
		currentPath := strings.Join(parts[:i+1], ".")

		if newName, ok := schemaOptions.Aliases[currentPath]; ok {
			parts[i] = newName

			if i == len(parts)-1 {
				attr.Name = SnakeCaseString(newName)
			}
		}
	}

	*attrPathName = strings.Join(parts, ".")
}

func applyOverrides(attr *Attribute, attrPathName string, schemaOptions config.SchemaOptions) {
	if override, ok := schemaOptions.Overrides[attrPathName]; ok {
		attr.Description = &override.Description
	}
}

func processNestedAttributes(attr *Attribute, schemaOptions config.SchemaOptions, attrPathName string) {
	switch {
	case attr.ListNested != nil:
		applySchemaOptions(schemaOptions, &attr.ListNested.NestedObject.Attributes, attrPathName)
	case attr.SingleNested != nil:
		applySchemaOptions(schemaOptions, &attr.SingleNested.NestedObject.Attributes, attrPathName)
	case attr.SetNested != nil:
		applySchemaOptions(schemaOptions, &attr.SetNested.NestedObject.Attributes, attrPathName)
	case attr.MapNested != nil:
		applySchemaOptions(schemaOptions, &attr.MapNested.NestedObject.Attributes, attrPathName)
	}
}
