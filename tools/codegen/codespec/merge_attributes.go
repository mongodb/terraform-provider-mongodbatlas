package codespec

import (
	"sort"
)

func mergeNestedAttributes(existingAttrs *Attributes, newAttrs Attributes, computability ComputedOptionalRequired, isFromResponse bool) {
	mergedMap := make(map[string]*Attribute)
	if existingAttrs != nil {
		for i := range *existingAttrs {
			mergedMap[(*existingAttrs)[i].Name.SnakeCase()] = &(*existingAttrs)[i]
		}
	}

	for i := range newAttrs {
		newAttr := &newAttrs[i]

		if _, exists := mergedMap[newAttr.Name.SnakeCase()]; exists {
			addOrUpdate(newAttr, computability, mergedMap, isFromResponse)
		} else {
			newAttr.ComputedOptionalRequired = computability
			mergedMap[newAttr.Name.SnakeCase()] = newAttr
		}
	}

	*existingAttrs = make(Attributes, 0, len(mergedMap))
	for _, attr := range mergedMap {
		*existingAttrs = append(*existingAttrs, *attr)
	}

	sortAttributes(*existingAttrs)
}

func addOrUpdate(attr *Attribute, computability ComputedOptionalRequired, merged map[string]*Attribute, isFromResponse bool) {
	if existingAttr, exists := merged[attr.Name.SnakeCase()]; exists {
		if existingAttr.Description == nil || *existingAttr.Description == "" {
			existingAttr.Description = attr.Description
		}

		if !isFromResponse && existingAttr.ComputedOptionalRequired != Required {
			existingAttr.ComputedOptionalRequired = computability
		}

		if existingAttr.ListNested != nil && attr.ListNested != nil {
			mergeNestedAttributes(&existingAttr.ListNested.NestedObject.Attributes, attr.ListNested.NestedObject.Attributes, computability, isFromResponse)
		} else if attr.ListNested != nil {
			existingAttr.ListNested = attr.ListNested
		}

		if existingAttr.SingleNested != nil && attr.SingleNested != nil {
			mergeNestedAttributes(&existingAttr.SingleNested.NestedObject.Attributes, attr.SingleNested.NestedObject.Attributes, computability, isFromResponse)
		} else if attr.SingleNested != nil {
			existingAttr.SingleNested = attr.SingleNested
		}

		if existingAttr.SetNested != nil && attr.SetNested != nil {
			mergeNestedAttributes(&existingAttr.SetNested.NestedObject.Attributes, attr.SetNested.NestedObject.Attributes, computability, isFromResponse)
		} else if attr.SetNested != nil {
			existingAttr.SetNested = attr.SetNested
		}

		if existingAttr.MapNested != nil && attr.MapNested != nil {
			mergeNestedAttributes(&existingAttr.MapNested.NestedObject.Attributes, attr.MapNested.NestedObject.Attributes, computability, isFromResponse)
		} else if attr.MapNested != nil {
			existingAttr.MapNested = attr.MapNested
		}
	} else {
		newAttr := *attr
		newAttr.ComputedOptionalRequired = computability
		merged[attr.Name.SnakeCase()] = &newAttr
	}
}

func mergeAttributes(pathParams, createRequest, createResponse, readResponse Attributes) Attributes {
	merged := make(map[string]*Attribute)

	// Path parameters: all attributes will be "required"
	for i := range pathParams {
		addOrUpdate(&pathParams[i], Required, merged, false)
	}

	// POST request body: optional/required is as defined
	for i := range createRequest {
		addOrUpdate(&createRequest[i], createRequest[i].ComputedOptionalRequired, merged, false)
	}

	// POST/GET response body: properties not in the request body are "computed" or "computed_optional" (if a default is present)
	for i := range createResponse {
		if isOptional(&createResponse[i]) {
			addOrUpdate(&createResponse[i], ComputedOptional, merged, true)
		} else {
			addOrUpdate(&createResponse[i], Computed, merged, true)
		}
	}

	for i := range readResponse {
		if isOptional(&readResponse[i]) {
			addOrUpdate(&readResponse[i], ComputedOptional, merged, true)
		} else {
			addOrUpdate(&readResponse[i], Computed, merged, true)
		}
	}

	resourceAttributes := make(Attributes, 0, len(merged))
	for _, attr := range merged {
		resourceAttributes = append(resourceAttributes, *attr)
	}

	sortAttributes(resourceAttributes)

	updateNestedComputability(&resourceAttributes, Optional)

	return resourceAttributes
}

func updateNestedComputability(attrs *Attributes, parentComputability ComputedOptionalRequired) {
	for i := range *attrs {
		attr := &(*attrs)[i]

		if parentComputability == Computed {
			attr.ComputedOptionalRequired = Computed
		}

		if attr.ListNested != nil {
			updateNestedComputability(&attr.ListNested.NestedObject.Attributes, attr.ComputedOptionalRequired)
		}
		if attr.SingleNested != nil {
			updateNestedComputability(&attr.SingleNested.NestedObject.Attributes, attr.ComputedOptionalRequired)
		}
		if attr.SetNested != nil {
			updateNestedComputability(&attr.SetNested.NestedObject.Attributes, attr.ComputedOptionalRequired)
		}
		if attr.MapNested != nil {
			updateNestedComputability(&attr.MapNested.NestedObject.Attributes, attr.ComputedOptionalRequired)
		}
	}
}

func isOptional(attr *Attribute) bool {
	return (attr.Bool != nil && attr.Bool.Default != nil) ||
		(attr.Int64 != nil && attr.Int64.Default != nil) ||
		(attr.String != nil && attr.String.Default != nil) ||
		(attr.Float64 != nil && attr.Float64.Default != nil) ||
		(attr.Number != nil && attr.Number.Default != nil)
}

func sortAttributes(attrs Attributes) {
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Name < attrs[j].Name
	})
}
