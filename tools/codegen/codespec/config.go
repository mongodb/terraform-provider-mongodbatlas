package codespec

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
)

var commonIgnoredAttributes = []string{"total_count", "envelope", "items_per_page", "page_num", "links", "pretty", "include_count"}

const DeleteOnCreateTimeoutDescription = "Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. " +
	"When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for " +
	"deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a " +
	"transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`."

func ApplyTransformationsToResource(resourceConfig *config.Resource, resource *Resource) error {
	if resource == nil || resource.Schema == nil {
		return nil
	}
	parentPaths := &attrPaths{schemaPath: "", apiPath: ""}
	if err := applyAttributeTransformationsList(resourceConfig.SchemaOptions, &resource.Schema.Attributes, parentPaths, resourceTransformations); err != nil {
		return fmt.Errorf("failed to apply attribute transformations: %w", err)
	}
	applyCommonAliasTransformations(&resource.Operations, resourceConfig.SchemaOptions.Aliases, resource.Schema.Discriminator, &resource.Schema.Attributes)
	ApplyDeleteOnCreateTimeoutTransformation(resource)
	ApplyTimeoutTransformation(resource)
	return nil
}

// ApplyTransformationsToDataSources applies schema transformations and path param aliasing to data sources.
// This mirrors ApplyTransformationsToResource for resources, without timeout-related and create-only transformations.
// Exported for testing purposes.
func ApplyTransformationsToDataSources(dsConfig *config.DataSources, ds *DataSources) error {
	if ds == nil || ds.Schema == nil {
		return nil
	}

	parentPaths := &attrPaths{schemaPath: "", apiPath: ""}
	if err := applyAttributeTransformationsList(dsConfig.SchemaOptions, ds.Schema.SingularDSAttributes, parentPaths, dataSourceTransformations); err != nil {
		return fmt.Errorf("failed to apply attribute transformations for singular data source: %w", err)
	}
	if err := applyAttributeTransformationsList(dsConfig.SchemaOptions, ds.Schema.PluralDSAttributes, parentPaths, dataSourceTransformations); err != nil {
		return fmt.Errorf("failed to apply attribute transformations for plural data source: %w", err)
	}

	applyCommonAliasTransformations(
		&ds.Operations,
		dsConfig.SchemaOptions.Aliases,
		nil,
		ds.Schema.SingularDSAttributes,
		ds.Schema.PluralDSAttributes,
	)
	return nil
}

func applyCommonAliasTransformations(operations *APIOperations, aliases map[string]string, rootDiscriminator *Discriminator, attributeSets ...*Attributes) {
	applyAliasesToDiscriminator(rootDiscriminator, aliases, "")
	for _, attributes := range attributeSets {
		if attributes == nil {
			continue
		}
		applyAliasesToNestedDiscriminators(*attributes, aliases, "")
	}
	applyAliasToPathParams(operations, aliases)
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

var resourceTransformations = []AttributeTransformation{
	aliasTransformation,
	commonRSAndDSOverridesTransformation,
	immutableComputedOverrideTransformation,
	tagsAndLabelsAsMapTypeTransformation,
	createOnlyTransformation,
	requestOnlyRequiredOnCreateTransformation,
}

var dataSourceTransformations = []AttributeTransformation{
	aliasTransformation,
	commonRSAndDSOverridesTransformation,
	tagsAndLabelsAsMapTypeTransformation,
	// Note: resource-specific transformations (createOnly, requestOnlyRequiredOnCreate, immutableComputed) are excluded for data sources
}

func applyAttributeTransformationsList(schemaOptions config.SchemaOptions, attributes *Attributes, parentPaths *attrPaths, transformationList []AttributeTransformation) error {
	if attributes == nil {
		return nil
	}

	schemaOptions.Ignores = append(schemaOptions.Ignores, commonIgnoredAttributes...)
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

func shouldIgnoreAttribute(attrPathName string, ignoredAttrs map[string]bool) bool {
	return ignoredAttrs[attrPathName]
}

func applyAliasToAttribute(attr *Attribute, paths *attrPaths, schemaOptions config.SchemaOptions) {
	// Config uses full camelCase for aliases (e.g., groupId: projectId, nestedObject.innerAttr: renamedAttr)
	// The apiPath is built from APIName values which preserve the original casing (e.g., "MongoDBMajorVersion")
	// This avoids the lossy conversion from snake to camel case.

	// Lookup by full apiPath only. At root level apiPath equals attr.APIName so non-dotted
	// aliases like "groupId: projectId" still match. At nested levels only explicitly
	// path-scoped aliases (e.g., "nestedObj.innerAttr: renamedAttr") apply.
	aliasCamel, found := schemaOptions.Aliases[paths.apiPath]

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

func commonRSAndDSOverridesTransformation(attr *Attribute, paths *attrPaths, schemaOptions config.SchemaOptions) error {
	// Overrides use the snake_case schemaPath
	override, ok := schemaOptions.Overrides[attrPathForOverrides(paths.schemaPath)]
	if !ok {
		return nil
	}
	if override.Description != "" {
		attr.Description = &override.Description
	}
	if override.Computability != nil {
		attr.ComputedOptionalRequired = getComputabilityFromConfig(*override.Computability)
	}
	if override.Sensitive != nil {
		attr.Sensitive = *override.Sensitive
	}
	if override.RequestBodyUsage != nil {
		if err := applyReqBodyUsageOverride(*override.RequestBodyUsage, attr); err != nil {
			return err
		}
	}
	if override.Type != nil {
		if err := applyTypeOverride(&override, attr); err != nil {
			return err
		}
	}
	if override.SkipStateListMerge != nil {
		attr.SkipStateListMerge = *override.SkipStateListMerge
	}
	return nil
}

func immutableComputedOverrideTransformation(attr *Attribute, paths *attrPaths, schemaOptions config.SchemaOptions) error {
	override, ok := schemaOptions.Overrides[attrPathForOverrides(paths.schemaPath)]
	if !ok {
		return nil
	}
	if override.ImmutableComputed != nil {
		attr.ImmutableComputed = *override.ImmutableComputed
	}
	return nil
}

func createOnlyTransformation(attr *Attribute, _ *attrPaths, _ config.SchemaOptions) error {
	setCreateOnlyValue(attr)
	return nil
}

func requestOnlyRequiredOnCreateTransformation(attr *Attribute, _ *attrPaths, _ config.SchemaOptions) error {
	if attr.ComputedOptionalRequired == Required && attr.ReqBodyUsage == OmitInUpdateBody && !attr.PresentInAnyResponse {
		attr.RequestOnlyRequiredOnCreate = true
		attr.ComputedOptionalRequired = Optional
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
		if attr.ListNested != nil {
			if attr.CustomType != nil {
				attr.CustomType = NewCustomNestedSetType(attr.CustomType.Name)
			}
			attr.SetNested = &SetNestedAttribute{NestedObject: attr.ListNested.NestedObject}
			attr.ListNested = nil
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
		if attr.SetNested != nil {
			if attr.CustomType != nil {
				attr.CustomType = NewCustomNestedListType(attr.CustomType.Name)
			}
			attr.ListNested = &ListNestedAttribute{NestedObject: attr.SetNested.NestedObject}
			attr.SetNested = nil
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

func applyReqBodyUsageOverride(reqBodyUsage config.ReqBodyUsage, attr *Attribute) error {
	switch reqBodyUsage {
	case config.SendNullAsNullOnUpdate:
		attr.ReqBodyUsage = SendNullAsNullOnUpdate
		return nil
	case config.SendNullAsEmptyOnUpdate:
		attr.ReqBodyUsage = SendNullAsEmptyOnUpdate
		return nil
	}
	return fmt.Errorf("unsupported request body usage defined in configuration: %s for attribute %s", reqBodyUsage, attr.TFSchemaName)
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

func attrPathForOverrides(attrPathName string) string {
	return strings.TrimPrefix(attrPathName, "results.")
}

// applyAliasesToNestedDiscriminators recursively walks attributes and applies alias renames
// to discriminators found on nested attribute objects.
func applyAliasesToNestedDiscriminators(attributes Attributes, aliases map[string]string, parentAPIPath string) {
	for i := range attributes {
		attr := &attributes[i]
		apiPath := buildPath(parentAPIPath, attr.APIName)

		switch {
		case attr.ListNested != nil:
			applyAliasesToDiscriminator(attr.ListNested.NestedObject.Discriminator, aliases, apiPath)
			applyAliasesToNestedDiscriminators(attr.ListNested.NestedObject.Attributes, aliases, apiPath)
		case attr.SingleNested != nil:
			applyAliasesToDiscriminator(attr.SingleNested.NestedObject.Discriminator, aliases, apiPath)
			applyAliasesToNestedDiscriminators(attr.SingleNested.NestedObject.Attributes, aliases, apiPath)
		case attr.SetNested != nil:
			applyAliasesToDiscriminator(attr.SetNested.NestedObject.Discriminator, aliases, apiPath)
			applyAliasesToNestedDiscriminators(attr.SetNested.NestedObject.Attributes, aliases, apiPath)
		case attr.MapNested != nil:
			applyAliasesToDiscriminator(attr.MapNested.NestedObject.Discriminator, aliases, apiPath)
			applyAliasesToNestedDiscriminators(attr.MapNested.NestedObject.Attributes, aliases, apiPath)
		}
	}
}

// applyAliasesToDiscriminator reconciles discriminator property names and variant attribute names
// when aliases have been applied. It renames PropertyName and all entries in Allowed/Required lists
// according to the aliases applicable at this nesting level.
func applyAliasesToDiscriminator(disc *Discriminator, aliases map[string]string, parentAPIPath string) {
	if disc == nil || len(aliases) == 0 {
		return
	}

	// Build a rename map: old_snake -> new_snake for aliases at this nesting level.
	// Aliases use camelCase API names (e.g., "groupId: projectId" or "nestedObject.innerAttr: renamedAttr").
	// At root level only non-dotted aliases apply; at nested levels only path-scoped aliases
	// whose prefix matches parentAPIPath apply (and only for the immediate child, not deeper).
	renameMap := make(map[string]string)
	for original, alias := range aliases {
		var apiName string
		if parentAPIPath == "" {
			// At root level, only non-dotted aliases apply
			if !strings.Contains(original, ".") {
				apiName = original
			}
		} else {
			// At nested levels, only path-scoped aliases with matching prefix apply
			prefix := parentAPIPath + "."
			if strings.HasPrefix(original, prefix) {
				leafName := strings.TrimPrefix(original, prefix)
				// Only apply if the leaf targets this exact level (no further dots)
				if !strings.Contains(leafName, ".") {
					apiName = leafName
				}
			}
		}
		if apiName != "" {
			oldSnake := stringcase.ToSnakeCase(apiName)
			newSnake := stringcase.ToSnakeCase(alias)
			if oldSnake != newSnake {
				renameMap[oldSnake] = newSnake
			}
		}
	}

	if len(renameMap) == 0 {
		return
	}

	// Rename PropertyName if aliased
	if newName, found := renameMap[disc.PropertyName]; found {
		disc.PropertyName = newName
	}

	// Rename entries in Allowed and Required lists
	for key, variant := range disc.Mapping {
		variant.Allowed = renameStringSlice(variant.Allowed, renameMap)
		variant.Required = renameStringSlice(variant.Required, renameMap)
		disc.Mapping[key] = variant
	}
}

// renameStringSlice applies renames from the map to a string slice, returning a new sorted slice.
func renameStringSlice(items []string, renameMap map[string]string) []string {
	if len(items) == 0 {
		return items
	}
	result := make([]string, len(items))
	for i, item := range items {
		if newName, found := renameMap[item]; found {
			result[i] = newName
		} else {
			result[i] = item
		}
	}
	sort.Strings(result)
	return result
}

// tagsAndLabelsAsMapTypeTransformation transforms attributes that represent collections of key/value pairs (tags and labels) from a nested list of objects into a Map type.
// This makes the Terraform schema expose a Map type while the underlying Atlas API still uses the array of {key, value} objects.
func tagsAndLabelsAsMapTypeTransformation(attr *Attribute, _ *attrPaths, _ config.SchemaOptions) error {
	if attr.TFSchemaName != "tags" && attr.TFSchemaName != "labels" {
		return nil
	}

	// We only transform attributes that are currently modeled as list_nested with a nested object that has exactly "key" and "value" string attributes.
	if attr.ListNested == nil || attr.Map != nil || attr.MapNested != nil {
		return nil
	}
	nestedAttrs := attr.ListNested.NestedObject.Attributes
	if len(nestedAttrs) != 2 {
		return nil
	}
	for i := range nestedAttrs {
		nested := nestedAttrs[i]
		if nested.TFSchemaName != "key" && nested.TFSchemaName != "value" {
			return nil
		}
		if nested.String == nil {
			return nil
		}
	}

	// Rewrite the attribute as a Map of strings with the standard custom map type.
	attr.ListNested = nil
	attr.CustomType = NewCustomMapType(String)
	attr.Map = &MapAttribute{
		ElementType: String,
	}
	attr.ListTypeAsMap = true
	// Send an empty map on updates if the tags/labels attribute is null
	attr.ReqBodyUsage = SendNullAsEmptyOnUpdate
	return nil
}
