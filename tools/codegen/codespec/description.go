package codespec

import (
	"fmt"
	"slices"
	"strings"
)

const (
	DescriptionPrefixRequired   = "Required for"
	DescriptionPrefixOptional   = "Optional for"
	DescriptionPrefixApplicable = "Applies to"
)

// EnhanceDescriptionsWithDiscriminator prepends polymorphic type context to attribute descriptions.
// These prefixes serve two purposes: the doc post-processor parses them to generate per-type
// subsections in the published documentation, and they surface discriminator constraints in
// schema-aware tools (IDE autocompletion, LSP hovers) that display attribute descriptions.
// isDataSource=false (resources) distinguishes "Required for" vs "Optional for";
// isDataSource=true (data sources) always uses "Applies to" since all attributes are read-only.
func EnhanceDescriptionsWithDiscriminator(attrs Attributes, disc *Discriminator, isDataSource bool) {
	enhanceCurrentLevel(attrs, disc, isDataSource)

	for i := range attrs {
		attr := &attrs[i]
		if nested := attr.NestedObject(); nested != nil {
			EnhanceDescriptionsWithDiscriminator(nested.Attributes, nested.Discriminator, isDataSource)
		}
	}
}

type attrTypeInfo struct {
	requiredTypes []string
	optionalTypes []string
}

func enhanceCurrentLevel(attrs Attributes, disc *Discriminator, isDataSource bool) {
	if disc == nil {
		return
	}

	reverseIndex := buildReverseIndex(disc)

	discriminatorPropName := disc.PropertyName.TFSchemaName
	for i := range attrs {
		attr := &attrs[i]
		if attr.TFSchemaName == discriminatorPropName {
			continue
		}
		if attr.Description == nil {
			continue
		}

		info, found := reverseIndex[attr.TFSchemaName]
		if !found {
			continue
		}

		if prefix := buildPrefix(info, discriminatorPropName, isDataSource); prefix != "" {
			enhanced := prefix + " " + *attr.Description
			attr.Description = &enhanced
		}
	}
}

func buildReverseIndex(disc *Discriminator) map[string]*attrTypeInfo {
	index := make(map[string]*attrTypeInfo)

	for typeName, variant := range disc.Mapping {
		requiredSet := make(map[string]bool, len(variant.Required))
		for _, r := range variant.Required {
			requiredSet[r.TFSchemaName] = true
		}

		for _, a := range variant.Allowed {
			info, ok := index[a.TFSchemaName]
			if !ok {
				info = &attrTypeInfo{}
				index[a.TFSchemaName] = info
			}
			if requiredSet[a.TFSchemaName] {
				info.requiredTypes = append(info.requiredTypes, typeName)
			} else {
				info.optionalTypes = append(info.optionalTypes, typeName)
			}
		}
	}

	for _, info := range index {
		slices.Sort(info.requiredTypes)
		slices.Sort(info.optionalTypes)
	}

	return index
}

func buildPrefix(info *attrTypeInfo, discriminatorName string, isDataSource bool) string {
	if isDataSource {
		allTypes := append(slices.Clone(info.requiredTypes), info.optionalTypes...)
		slices.Sort(allTypes)
		return fmt.Sprintf("%s %s: %s.", DescriptionPrefixApplicable, discriminatorName, strings.Join(allTypes, ", "))
	}

	var parts []string
	if len(info.requiredTypes) > 0 {
		parts = append(parts, fmt.Sprintf("%s %s: %s.", DescriptionPrefixRequired, discriminatorName, strings.Join(info.requiredTypes, ", ")))
	}
	if len(info.optionalTypes) > 0 {
		parts = append(parts, fmt.Sprintf("%s %s: %s.", DescriptionPrefixOptional, discriminatorName, strings.Join(info.optionalTypes, ", ")))
	}
	return strings.Join(parts, " ")
}
