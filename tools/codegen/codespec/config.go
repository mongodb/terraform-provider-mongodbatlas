package codespec

import (
	"fmt"
	"log"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
)

const DeleteOnCreateTimeoutDescription = "Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. " +
	"When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for " +
	"deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a " +
	"transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`."

func applyTransformationsWithConfigOpts(resourceConfig *config.Resource, resource *Resource) {
	applyAttributeTransformations(resourceConfig.SchemaOptions, &resource.Schema.Attributes, "")
	applyAliasToPathParams(resource, resourceConfig.SchemaOptions.Aliases)
	ApplyDeleteOnCreateTimeoutTransformation(resource)
	ApplyTimeoutTransformation(resource)
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
		if override.Type != nil {
			applyTypeOverride(&override, attr)
		}
	}
}

func applyTypeOverride(override *config.Override, attr *Attribute) {
	switch *override.Type {
	case config.Set:
		if attr.List != nil {
			attr.CustomType = nil
			/* TODO revisit once CustomSetType is supported - CLOUDP-353170
			if attr.CustomType != nil {
				attr.CustomType = NewCustomSetType(attr.List.ElementType)
			}
			*/

			attr.Set = &SetAttribute{ElementType: attr.List.ElementType}
			attr.List = nil
			return
		}
	case config.List:
		if attr.Set != nil {
			/* TODO uncomment once CustomSetType is supported - CLOUDP-353170
			if attr.CustomType != nil {
				attr.CustomType = NewCustomListType(attr.Set.ElementType)
			}
			*/
			attr.List = &ListAttribute{ElementType: attr.Set.ElementType}
			attr.Set = nil
			return
		}
	default:
		log.Printf("[WARN] %s - Unsupported type override defined in configuration: %s", attr.TFSchemaName, *override.Type)
		return
	}
	log.Printf("[WARN] %s - Unsupported override from original type to: %s", attr.TFSchemaName, *override.Type)
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

// ApplyTimeoutTransformation adds a timeout attribute to the resource schema if any operation has wait blocks.
func ApplyTimeoutTransformation(resource *Resource) {
	ops := &resource.Operations
	var configurableTimeouts []Operation

	if ops.Create.Wait != nil {
		configurableTimeouts = append(configurableTimeouts, Create)
	}
	if ops.Update.Wait != nil {
		configurableTimeouts = append(configurableTimeouts, Update)
	}
	if ops.Read.Wait != nil {
		configurableTimeouts = append(configurableTimeouts, Read)
	}
	// Delete operation is optional
	if ops.Delete != nil && ops.Delete.Wait != nil {
		configurableTimeouts = append(configurableTimeouts, Delete)
	}

	if len(configurableTimeouts) > 0 {
		resource.Schema.Attributes = append(resource.Schema.Attributes, Attribute{
			TFSchemaName: "timeouts",
			TFModelName:  "Timeouts",
			Timeouts:     &TimeoutsAttribute{ConfigurableTimeouts: configurableTimeouts},
			ReqBodyUsage: OmitAlways,
		})
	}
}

// ApplyDeleteOnCreateTimeoutTransformation adds a delete_on_create_timeout attribute to the resource schema
// if the Create operation has a wait block and the Delete operation exists.
func ApplyDeleteOnCreateTimeoutTransformation(resource *Resource) {
	if ops := &resource.Operations; ops.Create.Wait == nil || ops.Delete == nil {
		return
	}
	resource.Schema.Attributes = append(resource.Schema.Attributes, Attribute{
		TFSchemaName:             "delete_on_create_timeout",
		TFModelName:              "DeleteOnCreateTimeout",
		Bool:                     &BoolAttribute{Default: conversion.Pointer(true)},
		Description:              conversion.StringPtr(DeleteOnCreateTimeoutDescription),
		ReqBodyUsage:             OmitAlways,
		CreateOnly:               true,
		ComputedOptionalRequired: ComputedOptional,
	})
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
