package codespec

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
)

const DeleteOnCreateTimeoutDescription = "Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. " +
	"When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for " +
	"deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a " +
	"transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`."

func applyTransformationsWithConfigOpts(resourceConfig *config.Resource, resource *Resource) error {
	if err := applyAttributeTransformations(resourceConfig.SchemaOptions, &resource.Schema.Attributes, ""); err != nil {
		return fmt.Errorf("failed to apply attribute transformations: %w", err)
	}
	applyAliasToPathParams(resource, resourceConfig.SchemaOptions.Aliases)
	ApplyDeleteOnCreateTimeoutTransformation(resource)
	ApplyTimeoutTransformation(resource)
	return nil
}

// AttributeTransformation represents a operation applied to an attribute during traversal.
// Implementations may mutate the attribute in-place.
type AttributeTransformation func(attr *Attribute, attrPathName *string, schemaOptions config.SchemaOptions) error

var transformations = []AttributeTransformation{
	aliasTransformation,
	overridesTransformation,
	createOnlyTransformation,
}

func applyAttributeTransformations(schemaOptions config.SchemaOptions, attributes *Attributes, parentName string) error {
	ignoredAttrs := getIgnoredAttributesMap(schemaOptions.Ignores)

	var finalAttributes Attributes

	for i := range *attributes {
		attr := &(*attributes)[i]
		attrPathName := getAttributePathName(attr.TFSchemaName, parentName)

		if shouldIgnoreAttribute(attrPathName, ignoredAttrs) {
			continue
		}

		for _, t := range transformations {
			if err := t(attr, &attrPathName, schemaOptions); err != nil {
				return err
			}
		}

		// apply transformations to nested attributes
		switch {
		case attr.ListNested != nil:
			if err := applyAttributeTransformations(schemaOptions, &attr.ListNested.NestedObject.Attributes, attrPathName); err != nil {
				return err
			}
		case attr.SingleNested != nil:
			if err := applyAttributeTransformations(schemaOptions, &attr.SingleNested.NestedObject.Attributes, attrPathName); err != nil {
				return err
			}
		case attr.SetNested != nil:
			if err := applyAttributeTransformations(schemaOptions, &attr.SetNested.NestedObject.Attributes, attrPathName); err != nil {
				return err
			}
		case attr.MapNested != nil:
			if err := applyAttributeTransformations(schemaOptions, &attr.MapNested.NestedObject.Attributes, attrPathName); err != nil {
				return err
			}
		}

		finalAttributes = append(finalAttributes, *attr)
	}

	*attributes = finalAttributes
	return nil
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
	if newModelName, ok := schemaOptions.Aliases[attr.TFModelName]; ok {
		attr.TFModelName = newModelName
		attr.TFSchemaName = stringcase.ToSnakeCase(newModelName)
		parts := strings.Split(*attrPathName, ".")
		if len(parts) > 0 {
			parts[len(parts)-1] = attr.TFSchemaName
			*attrPathName = strings.Join(parts, ".")
		}
	}
}

func applyAliasToPathParams(resource *Resource, aliases map[string]string) {
	for original, alias := range aliases {
		originalCamel := stringcase.Uncapitalize(original)
		aliasCamel := stringcase.Uncapitalize(alias)
		resource.Operations.Create.Path = strings.ReplaceAll(resource.Operations.Create.Path, fmt.Sprintf("{%s}", originalCamel), fmt.Sprintf("{%s}", aliasCamel))
		resource.Operations.Read.Path = strings.ReplaceAll(resource.Operations.Read.Path, fmt.Sprintf("{%s}", originalCamel), fmt.Sprintf("{%s}", aliasCamel))
		if resource.Operations.Update != nil {
			resource.Operations.Update.Path = strings.ReplaceAll(resource.Operations.Update.Path, fmt.Sprintf("{%s}", originalCamel), fmt.Sprintf("{%s}", aliasCamel))
		}
		if resource.Operations.Delete != nil {
			resource.Operations.Delete.Path = strings.ReplaceAll(resource.Operations.Delete.Path, fmt.Sprintf("{%s}", originalCamel), fmt.Sprintf("{%s}", aliasCamel))
		}
	}
}

// Transformations
func aliasTransformation(attr *Attribute, attrPathName *string, schemaOptions config.SchemaOptions) error {
	// the config is expected to use alias name for defining any subsequent overrides (description, etc)
	applyAliasToAttribute(attr, attrPathName, schemaOptions)
	return nil
}

func overridesTransformation(attr *Attribute, attrPathName *string, schemaOptions config.SchemaOptions) error {
	return applyOverrides(attr, *attrPathName, schemaOptions)
}

func createOnlyTransformation(attr *Attribute, _ *string, _ config.SchemaOptions) error {
	setCreateOnlyValue(attr)
	return nil
}

func applyOverrides(attr *Attribute, attrPathName string, schemaOptions config.SchemaOptions) error {
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
			if err := applyTypeOverride(&override, attr); err != nil {
				return err
			}
		}
	}
	return nil
}

func applyTypeOverride(override *config.Override, attr *Attribute) error {
	switch *override.Type {
	case config.Set:
		if attr.List != nil {
			if attr.CustomType != nil {
				attr.CustomType = NewCustomSetType(attr.List.ElementType)
			}
			attr.Set = &SetAttribute{ElementType: attr.List.ElementType}
			attr.List = nil
			return nil
		}
	case config.List:
		if attr.Set != nil {
			if attr.CustomType != nil {
				attr.CustomType = NewCustomListType(attr.Set.ElementType)
			}
			attr.List = &ListAttribute{ElementType: attr.Set.ElementType}
			attr.Set = nil
			return nil
		}
	default:
		return fmt.Errorf("unsupported type override defined in configuration: %s for attribute %s", *override.Type, attr.TFSchemaName)
	}
	return fmt.Errorf("unsupported override from original type to %s for attribute %s", *override.Type, attr.TFSchemaName)
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
	// Update operation is optional
	if ops.Update != nil && ops.Update.Wait != nil {
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
