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
	// Start with empty paths for both schemaPath (snake_case) and apiPath (camelCase)
	if err := applyAttributeTransformations(resourceConfig.SchemaOptions, &resource.Schema.Attributes, &attrPaths{schemaPath: "", apiPath: ""}); err != nil {
		return fmt.Errorf("failed to apply attribute transformations: %w", err)
	}
	applyAliasToPathParams(&resource.Operations, resourceConfig.SchemaOptions.Aliases)
	ApplyDeleteOnCreateTimeoutTransformation(resource)
	ApplyTimeoutTransformation(resource)
	return nil
}

// ApplyTransformationsToDataSources applies schema transformations and path param aliasing to data sources.
// This mirrors applyTransformationsWithConfigOpts for resources, without timeout-related and create-only transformations.
// Exported for testing purposes.
func ApplyTransformationsToDataSources(dsConfig *config.DataSources, ds *DataSources) error {
	if ds == nil || ds.Schema == nil {
		return nil
	}

	// Apply attribute-level transformations (aliases, overrides, ignores) - excludes create-only for data sources
	if err := applyDataSourceAttributeTransformations(dsConfig.SchemaOptions, &ds.Schema.Attributes, &attrPaths{schemaPath: "", apiPath: ""}); err != nil {
		return fmt.Errorf("failed to apply attribute transformations: %w", err)
	}

	// Alias placeholders in operation paths after attribute transformations
	applyAliasToPathParams(&ds.Operations, dsConfig.SchemaOptions.Aliases)
	return nil
}

// applyAliasToPathParams replaces path parameter placeholders with their aliased names in all operation paths.
// Works for both resources (Create, Read, Update, Delete) and data sources (Read, List).
func applyAliasToPathParams(operations *APIOperations, aliases map[string]string) {
	if operations == nil {
		return
	}

	for original, alias := range aliases {
		placeholder := fmt.Sprintf("{%s}", original)
		aliasedPlaceholder := fmt.Sprintf("{%s}", alias)

		if operations.Create != nil {
			operations.Create.Path = strings.ReplaceAll(operations.Create.Path, placeholder, aliasedPlaceholder)
		}
		if operations.Read != nil {
			operations.Read.Path = strings.ReplaceAll(operations.Read.Path, placeholder, aliasedPlaceholder)
		}
		if operations.Update != nil {
			operations.Update.Path = strings.ReplaceAll(operations.Update.Path, placeholder, aliasedPlaceholder)
		}
		if operations.Delete != nil {
			operations.Delete.Path = strings.ReplaceAll(operations.Delete.Path, placeholder, aliasedPlaceholder)
		}
		if operations.List != nil {
			operations.List.Path = strings.ReplaceAll(operations.List.Path, placeholder, aliasedPlaceholder)
		}
	}
}

// attrPaths holds both the snake_case path (for overrides) and the camelCase path (for aliases).
// This avoids the lossy conversion from snake to camel case which can cause mismatches with API properties
// that have multiple consecutive uppercase letters (e.g., MongoDBMajorVersion).
type attrPaths struct {
	schemaPath string // snake_case path for overrides (e.g., "replication_specs.region_configs")
	apiPath    string // camelCase path for aliases (e.g., "replicationSpecs.regionConfigs")
}

// AttributeTransformation represents a operation applied to an attribute during traversal.
// Implementations may mutate the attribute in-place.
type AttributeTransformation func(attr *Attribute, paths *attrPaths, schemaOptions config.SchemaOptions) error

var transformations = []AttributeTransformation{
	aliasTransformation,
	overridesTransformation,
	createOnlyTransformation,
}

var dataSourceTransformations = []AttributeTransformation{
	aliasTransformation,
	overridesTransformation,
	// Note: createOnlyTransformation is excluded for data sources (read-only, no create operation)
}

func applyAttributeTransformations(schemaOptions config.SchemaOptions, attributes *Attributes, parentPaths *attrPaths) error {
	return applyAttributeTransformationsList(schemaOptions, attributes, parentPaths, transformations)
}

func applyDataSourceAttributeTransformations(schemaOptions config.SchemaOptions, attributes *Attributes, parentPaths *attrPaths) error {
	return applyAttributeTransformationsList(schemaOptions, attributes, parentPaths, dataSourceTransformations)
}

func applyAttributeTransformationsList(schemaOptions config.SchemaOptions, attributes *Attributes, parentPaths *attrPaths, transformationList []AttributeTransformation) error {
	ignoredAttrs := getIgnoredAttributesMap(schemaOptions.Ignores)

	var finalAttributes Attributes

	for i := range *attributes {
		attr := &(*attributes)[i]
		paths := attrPaths{
			schemaPath: buildPath(parentPaths.schemaPath, attr.TFSchemaName),
			apiPath:    buildPath(parentPaths.apiPath, attr.APIName),
		}

		if shouldIgnoreAttribute(paths.schemaPath, ignoredAttrs) {
			continue
		}

		for _, t := range transformationList {
			if err := t(attr, &paths, schemaOptions); err != nil {
				return err
			}
		}

		// apply transformations to nested attributes recursively with the same transformation list
		switch {
		case attr.ListNested != nil:
			if err := applyAttributeTransformationsList(schemaOptions, &attr.ListNested.NestedObject.Attributes, &paths, transformationList); err != nil {
				return err
			}
		case attr.SingleNested != nil:
			if err := applyAttributeTransformationsList(schemaOptions, &attr.SingleNested.NestedObject.Attributes, &paths, transformationList); err != nil {
				return err
			}
		case attr.SetNested != nil:
			if err := applyAttributeTransformationsList(schemaOptions, &attr.SetNested.NestedObject.Attributes, &paths, transformationList); err != nil {
				return err
			}
		case attr.MapNested != nil:
			if err := applyAttributeTransformationsList(schemaOptions, &attr.MapNested.NestedObject.Attributes, &paths, transformationList); err != nil {
				return err
			}
		}

		finalAttributes = append(finalAttributes, *attr)
	}

	*attributes = finalAttributes
	return nil
}

func buildPath(parentPath, attrName string) string {
	if parentPath == "" {
		return attrName
	}
	return parentPath + "." + attrName
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

func applyAliasToAttribute(attr *Attribute, paths *attrPaths, schemaOptions config.SchemaOptions) {
	// Config uses full camelCase for aliases (e.g., groupId: projectId, nestedObject.innerAttr: renamedAttr)
	// The apiPath is built from APIName values which preserve the original casing (e.g., "MongoDBMajorVersion")
	// This avoids the lossy conversion from snake to camel case.

	var aliasCamel string
	var found bool

	// First try path-based lookup for targeted aliasing of nested attributes
	// apiPath already contains the correct camelCase path built from APIName values
	if aliasCamel, found = schemaOptions.Aliases[paths.apiPath]; !found {
		// Fall back to attribute-name only lookup (for path params like groupId: projectId)
		aliasCamel, found = schemaOptions.Aliases[attr.APIName]
	}

	if found {
		// Change both TFSchemaName and TFModelName to the aliased name.
		// The APIName field preserves the original API property name, and the apiname tag
		// will be generated when TFModelName doesn't derive to the correct API name.
		attr.TFSchemaName = stringcase.ToSnakeCase(aliasCamel)
		attr.TFModelName = stringcase.Capitalize(aliasCamel)
		// Update the schema path to reflect the new schema name
		parts := strings.Split(paths.schemaPath, ".")
		if len(parts) > 0 {
			parts[len(parts)-1] = attr.TFSchemaName
			paths.schemaPath = strings.Join(parts, ".")
		}
		// Note: apiPath is not updated because it's only used for alias lookup,
		// and aliases are defined using original API names
	}
}

// Transformations
func aliasTransformation(attr *Attribute, paths *attrPaths, schemaOptions config.SchemaOptions) error {
	// Alias transformation runs first, updating TFSchemaName and paths.schemaPath.
	// Subsequent overrides should use the aliased snake_case path (e.g., nested_list_array_attr.inner_num_attr_alias)
	applyAliasToAttribute(attr, paths, schemaOptions)
	return nil
}

func overridesTransformation(attr *Attribute, paths *attrPaths, schemaOptions config.SchemaOptions) error {
	// Overrides use the snake_case schemaPath
	return applyOverrides(attr, paths.schemaPath, schemaOptions)
}

func createOnlyTransformation(attr *Attribute, _ *attrPaths, _ config.SchemaOptions) error {
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
