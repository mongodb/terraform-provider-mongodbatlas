package codespec

import (
	"fmt"
	"log"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
)

func applyTransformationsWithConfigOpts(resourceConfig *config.Resource, resource *Resource) {
	applyAttributeTransformations(resourceConfig.SchemaOptions, &resource.Schema.Attributes, "")

	applyAliasToPathParams(resource, resourceConfig.SchemaOptions.Aliases)
}

// AttributeTransformation represents a operation applied to an attribute during traversal.
// Implementations may mutate the attribute in-place.
type AttributeTransformation func(attr *Attribute, attrPathName *string, schemaOptions config.SchemaOptions)

var transformations = []AttributeTransformation{
	aliasTransformation,
	overridesTransformation,
	createOnlyTransformation,
}

func applyAttributeTransformations(schemaOptions config.SchemaOptions, attributes *Attributes, parentName string) {
	ignoredAttrs := getIgnoredAttributesMap(schemaOptions.Ignores)

	var finalAttributes Attributes

	for i := range *attributes {
		attr := &(*attributes)[i]
		attrPathName := getAttributePathName(string(attr.TFSchemaName), parentName)

		if shouldIgnoreAttribute(attrPathName, ignoredAttrs) {
			continue
		}

		for _, t := range transformations {
			t(attr, &attrPathName, schemaOptions)
		}

		// apply transformations to nested attributes
		switch {
		case attr.ListNested != nil:
			applyAttributeTransformations(schemaOptions, &attr.ListNested.NestedObject.Attributes, attrPathName)
		case attr.SingleNested != nil:
			applyAttributeTransformations(schemaOptions, &attr.SingleNested.NestedObject.Attributes, attrPathName)
		case attr.SetNested != nil:
			applyAttributeTransformations(schemaOptions, &attr.SetNested.NestedObject.Attributes, attrPathName)
		case attr.MapNested != nil:
			applyAttributeTransformations(schemaOptions, &attr.MapNested.NestedObject.Attributes, attrPathName)
		}

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
				attr.TFSchemaName = stringcase.SnakeCaseString(newName)
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
		if resource.Operations.Delete != nil {
			resource.Operations.Delete.Path = strings.ReplaceAll(resource.Operations.Delete.Path, fmt.Sprintf("{%s}", originalCamel), fmt.Sprintf("{%s}", aliasCamel))
		}
	}
}

// Transformations
func aliasTransformation(attr *Attribute, attrPathName *string, schemaOptions config.SchemaOptions) {
	// the config is expected to use alias name for defining any subsequent overrides (description, etc)
	applyAliasToAttribute(attr, attrPathName, schemaOptions)
}

func overridesTransformation(attr *Attribute, attrPathName *string, schemaOptions config.SchemaOptions) {
	applyOverrides(attr, *attrPathName, schemaOptions)
}

func createOnlyTransformation(attr *Attribute, _ *string, _ config.SchemaOptions) {
	setCreateOnlyValue(attr)
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
		if override.IncludeNullOnUpdate != nil && *override.IncludeNullOnUpdate {
			attr.ReqBodyUsage = IncludeNullOnUpdate
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
			TFSchemaName: "timeouts",
			Timeouts:     &TimeoutsAttribute{ConfigurableTimeouts: result},
			ReqBodyUsage: OmitAlways,
		}
	}
	return nil
}

func setCreateOnlyValue(attr *Attribute) {
	// CreateOnly plan modifier will not be applied for computed attributes
	if attr.ComputedOptionalRequired == Computed || attr.ComputedOptionalRequired == ComputedOptional {
		return
	}

	// captures case of path param attributes (no present in request body) and properties which are only present in post request
	if attr.ReqBodyUsage == OmitAlways || attr.ReqBodyUsage == OmitInUpdateBody {
		attr.CreateOnly = true
	}
}
