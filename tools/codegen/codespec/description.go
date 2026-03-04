package codespec

import (
	"fmt"
	"slices"
	"strings"
)

// EnhanceDescriptionsWithDiscriminator prepends polymorphic type context to attribute descriptions.
// useRequiredPrefix=true (resources) distinguishes "Required when" vs "Applicable when";
// useRequiredPrefix=false (data sources) always uses "Applicable when" since all attributes are read-only.
func EnhanceDescriptionsWithDiscriminator(attrs Attributes, disc *Discriminator, useRequiredPrefix bool) {
	enhanceCurrentLevel(attrs, disc, useRequiredPrefix)

	for i := range attrs {
		attr := &attrs[i]
		if nested := nestedObjectFromAttr(attr); nested != nil {
			EnhanceDescriptionsWithDiscriminator(nested.Attributes, nested.Discriminator, useRequiredPrefix)
		}
	}
}

type attrTypeInfo struct {
	requiredTypes []string
	optionalTypes []string
}

func enhanceCurrentLevel(attrs Attributes, disc *Discriminator, useRequiredPrefix bool) {
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

		prefix := buildPrefix(info, discriminatorPropName, useRequiredPrefix)
		if prefix != "" {
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

func buildPrefix(info *attrTypeInfo, discriminatorName string, useRequiredPrefix bool) string {
	if !useRequiredPrefix {
		allTypes := append(slices.Clone(info.requiredTypes), info.optionalTypes...)
		slices.Sort(allTypes)
		return fmt.Sprintf("Applicable when %s is: %s.", discriminatorName, strings.Join(allTypes, ", "))
	}

	var parts []string
	if len(info.requiredTypes) > 0 {
		parts = append(parts, fmt.Sprintf("Required when %s is: %s.", discriminatorName, strings.Join(info.requiredTypes, ", ")))
	}
	if len(info.optionalTypes) > 0 {
		parts = append(parts, fmt.Sprintf("Applicable when %s is: %s.", discriminatorName, strings.Join(info.optionalTypes, ", ")))
	}
	return strings.Join(parts, " ")
}

func nestedObjectFromAttr(attr *Attribute) *NestedAttributeObject {
	switch {
	case attr.ListNested != nil:
		return &attr.ListNested.NestedObject
	case attr.SetNested != nil:
		return &attr.SetNested.NestedObject
	case attr.SingleNested != nil:
		return &attr.SingleNested.NestedObject
	case attr.MapNested != nil:
		return &attr.MapNested.NestedObject
	default:
		return nil
	}
}
