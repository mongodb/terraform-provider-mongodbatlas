package codespec

import (
	"fmt"
	"log"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/stringcase"
)

func applyConfigSchemaOptions(resourceConfig *config.Resource, resource *Resource) {
	applySchemaOptions(resourceConfig.SchemaOptions, &resource.Schema.Attributes, "")
	applyAliasToPathParams(resource, resourceConfig.SchemaOptions.Aliases)
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
		applyAliasToAttribute(attr, &attrPathName, schemaOptions)

		applyOverrides(attr, attrPathName, schemaOptions)

		processNestedAttributes(attr, schemaOptions, attrPathName)

		finalAttributes = append(finalAttributes, *attr)
	}

	if timeoutAttr := applyTimeoutConfig(schemaOptions); parentName == "" && timeoutAttr != nil { // will not run for nested attributes
		finalAttributes = append(finalAttributes, *timeoutAttr)
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

func applyAliasToAttribute(attr *Attribute, attrPathName *string, schemaOptions config.SchemaOptions) {
	parts := strings.Split(*attrPathName, ".")

	for i := range parts {
		currentPath := strings.Join(parts[:i+1], ".")

		if newName, ok := schemaOptions.Aliases[currentPath]; ok {
			parts[i] = newName

			if i == len(parts)-1 {
				attr.Name = stringcase.SnakeCaseString(newName)
			}
		}
	}

	*attrPathName = strings.Join(parts, ".")
}

func applyAliasToPathParams(resource *Resource, aliases map[string]string) {
	for original, alias := range aliases {
		originalCamel := stringcase.SnakeCaseString(original).CamelCase()
		aliasCamel := stringcase.SnakeCaseString(alias).CamelCase()
		resource.Operations.Create.Path = strings.ReplaceAll(resource.Operations.Create.Path, fmt.Sprintf("{%s}", originalCamel), fmt.Sprintf("{%s}", aliasCamel))
		resource.Operations.Read.Path = strings.ReplaceAll(resource.Operations.Read.Path, fmt.Sprintf("{%s}", originalCamel), fmt.Sprintf("{%s}", aliasCamel))
		resource.Operations.Update.Path = strings.ReplaceAll(resource.Operations.Update.Path, fmt.Sprintf("{%s}", originalCamel), fmt.Sprintf("{%s}", aliasCamel))
		resource.Operations.Delete.Path = strings.ReplaceAll(resource.Operations.Delete.Path, fmt.Sprintf("{%s}", originalCamel), fmt.Sprintf("{%s}", aliasCamel))
	}
}

func applyOverrides(attr *Attribute, attrPathName string, schemaOptions config.SchemaOptions) {
	if override, ok := schemaOptions.Overrides[attrPathName]; ok {
		if override.Description != "" {
			attr.Description = &override.Description
		}
		if override.Computability != nil {
			attr.ComputedOptionalRequired = getComputabilityFromConfig(*override.Computability)
		}
		if override.Sensitive != nil {
			attr.Sensitive = *override.Sensitive
		}
		if override.IncludeJSONUpdate != nil && *override.IncludeJSONUpdate {
			attr.ReqBodyUsage = IncludeInUpdateBody
		}
	}
}

func getComputabilityFromConfig(computability config.Computability) ComputedOptionalRequired {
	if computability.Computed && computability.Optional {
		return ComputedOptional
	}
	if computability.Computed {
		return Computed
	}
	if computability.Optional {
		return Optional
	}
	return Required
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

func applyTimeoutConfig(options config.SchemaOptions) *Attribute {
	var result []Operation
	for _, op := range options.Timeouts {
		switch op {
		case "create":
			result = append(result, Create)
		case "read":
			result = append(result, Read)
		case "delete":
			result = append(result, Delete)
		case "update":
			result = append(result, Update)
		default:
			log.Printf("[WARN] Unknown operation type defined in timeout configuration: %s", op)
		}
	}
	if result != nil {
		return &Attribute{
			Name:         "timeouts",
			Timeouts:     &TimeoutsAttribute{ConfigurableTimeouts: result},
			ReqBodyUsage: OmitAlways,
		}
	}
	return nil
}
