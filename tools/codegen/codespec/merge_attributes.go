package codespec

import (
	"sort"
)

// mergeNestedAttributes recursively merges nested attributes
func mergeNestedAttributes(existingAttrs *Attributes, newAttrs Attributes, reqBodyUsage AttributeReqBodyUsage, isFromResponse bool) {
	mergedMap := make(map[string]*Attribute)
	if existingAttrs != nil {
		for i := range *existingAttrs {
			mergedMap[(*existingAttrs)[i].Name.SnakeCase()] = &(*existingAttrs)[i]
		}
	}

	// add new attributes and merge when necessary
	for i := range newAttrs {
		newAttr := &newAttrs[i]

		if existingAttr, found := mergedMap[newAttr.Name.SnakeCase()]; found {
			// merge computability of nested property using most restrictive value
			newComputability := mergeComputability(newAttr.ComputedOptionalRequired, existingAttr.ComputedOptionalRequired)
			addOrUpdate(newAttr, newComputability, reqBodyUsage, mergedMap, isFromResponse)
		} else {
			// setting as computed as nested attribute was defined in response
			if isFromResponse {
				newAttr.ComputedOptionalRequired = Computed
			}
			newAttr.ReqBodyUsage = reqBodyUsage
			mergedMap[newAttr.Name.SnakeCase()] = newAttr
		}
	}

	// update original existingAttrs with the merged result
	*existingAttrs = make(Attributes, 0, len(mergedMap))
	for _, attr := range mergedMap {
		*existingAttrs = append(*existingAttrs, *attr)
	}

	sortAttributes(*existingAttrs)
}

// mergeComputability merges two ComputedOptionalRequired values and returns the most restrictive one
func mergeComputability(first, second ComputedOptionalRequired) ComputedOptionalRequired {
	if first == Required || second == Required {
		return Required
	}
	if first == ComputedOptional || second == ComputedOptional {
		return ComputedOptional
	}
	if first == Optional || second == Optional {
		return Optional
	}
	return Computed
}

// addOrUpdate adds or updates an attribute in the merged map, including nested attributes
func addOrUpdate(newAttr *Attribute, targetComputability ComputedOptionalRequired, reqBodyUsage AttributeReqBodyUsage, merged map[string]*Attribute, isFromResponse bool) {
	if existingAttr, found := merged[newAttr.Name.SnakeCase()]; found {
		if existingAttr.Description == nil || *existingAttr.Description == "" {
			existingAttr.Description = newAttr.Description
		}

		// when property is in both request and response values computablity and reqBodyUsage will ignore information from response
		if !isFromResponse {
			existingAttr.ReqBodyUsage = reqBodyUsage
			// merging ensures if property is defined in POST and PATCH it will have the most restrictive computability
			existingAttr.ComputedOptionalRequired = mergeComputability(newAttr.ComputedOptionalRequired, existingAttr.ComputedOptionalRequired)
		}

		// handle nested attributes
		if existingAttr.ListNested != nil && newAttr.ListNested != nil {
			mergeNestedAttributes(&existingAttr.ListNested.NestedObject.Attributes, newAttr.ListNested.NestedObject.Attributes, reqBodyUsage, isFromResponse)
		} else if newAttr.ListNested != nil {
			existingAttr.ListNested = newAttr.ListNested
		}

		if existingAttr.SingleNested != nil && newAttr.SingleNested != nil {
			mergeNestedAttributes(&existingAttr.SingleNested.NestedObject.Attributes, newAttr.SingleNested.NestedObject.Attributes, reqBodyUsage, isFromResponse)
		} else if newAttr.SingleNested != nil {
			existingAttr.SingleNested = newAttr.SingleNested
		}

		if existingAttr.SetNested != nil && newAttr.SetNested != nil {
			mergeNestedAttributes(&existingAttr.SetNested.NestedObject.Attributes, newAttr.SetNested.NestedObject.Attributes, reqBodyUsage, isFromResponse)
		} else if newAttr.SetNested != nil {
			existingAttr.SetNested = newAttr.SetNested
		}

		if existingAttr.MapNested != nil && newAttr.MapNested != nil {
			mergeNestedAttributes(&existingAttr.MapNested.NestedObject.Attributes, newAttr.MapNested.NestedObject.Attributes, reqBodyUsage, isFromResponse)
		} else if newAttr.MapNested != nil {
			existingAttr.MapNested = newAttr.MapNested
		}
	} else {
		// add new attribute with the given computability
		newAttr.ComputedOptionalRequired = targetComputability
		newAttr.ReqBodyUsage = reqBodyUsage
		merged[newAttr.Name.SnakeCase()] = newAttr
	}
}

type attributeDefinitionSources struct {
	createPathParams, createRequest, updateRequest, createResponse, readResponse Attributes
}

// mergeAttributes merges attributes from different sources (path params, create/get operation bodies) and determines a single merged list of attributes.
// Computability and reqBodyUsage values are determined as part of this process.
// Different sources are applied in a specific order, defining the computability and reqBodyUsage value they have at each step.
func mergeAttributes(sources *attributeDefinitionSources) Attributes {
	merged := make(map[string]*Attribute)

	// create path parameters: all attributes will be "required", reqBodyUsage is defined as omit all at this step
	for i := range sources.createPathParams {
		addOrUpdate(&sources.createPathParams[i], Required, OmitAll, merged, false)
	}

	// POST request body: optional/required is as defined, reqBodyUsage is defined as OmitUpdateBody and will be updated to AllRequestBodies if present in POST request
	for i := range sources.createRequest {
		// for now we do not differentiate AllRequestBodies vs PostBodyOnly as we are not processing update request
		addOrUpdate(&sources.createRequest[i], sources.createRequest[i].ComputedOptionalRequired, OmitUpdateBody, merged, false)
	}

	// PATCH request body: optional/required is as defined, reqBodyUsage is defined as AllRequestBodies
	for i := range sources.updateRequest {
		addOrUpdate(&sources.updateRequest[i], sources.updateRequest[i].ComputedOptionalRequired, AllRequestBodies, merged, false)
	}

	// POST/GET response body: properties not in the request body are "computed" or "computed_optional" (if a default is present), reqBodyUsage will have OmitAll not present in request body
	for i := range sources.createResponse {
		if hasDefault(&sources.createResponse[i]) {
			addOrUpdate(&sources.createResponse[i], ComputedOptional, OmitAll, merged, true)
		} else {
			addOrUpdate(&sources.createResponse[i], Computed, OmitAll, merged, true)
		}
	}

	for i := range sources.readResponse {
		if hasDefault(&sources.readResponse[i]) {
			addOrUpdate(&sources.readResponse[i], ComputedOptional, OmitAll, merged, true)
		} else {
			addOrUpdate(&sources.readResponse[i], Computed, OmitAll, merged, true)
		}
	}

	resourceAttributes := make(Attributes, 0, len(merged))
	for _, attr := range merged {
		resourceAttributes = append(resourceAttributes, *attr)
	}

	sortAttributes(resourceAttributes)

	updateNestedComputabilityAndReqBodyUsage(&resourceAttributes, false, false)

	return resourceAttributes
}

// updateNestedComputabilityAndReqBodyUsage updates the computability and reqBodyUsage of nested attributes based on their parent attributes.
// If the parent is computed, all nested attributes are set to computed.
// If the parent is omitted in the request body, all nested attributes are set to omit all.
func updateNestedComputabilityAndReqBodyUsage(attrs *Attributes, parentIsComputed, parentIsOmittedInReqBody bool) {
	for i := range *attrs {
		attr := &(*attrs)[i]

		if parentIsComputed {
			attr.ComputedOptionalRequired = Computed
		}
		if parentIsOmittedInReqBody {
			attr.ReqBodyUsage = OmitAll
		}

		attrIsComputed := attr.ComputedOptionalRequired == Computed
		attrIsOmittedInReqBody := attr.ReqBodyUsage == OmitAll

		if attr.ListNested != nil {
			updateNestedComputabilityAndReqBodyUsage(&attr.ListNested.NestedObject.Attributes, attrIsComputed, attrIsOmittedInReqBody)
		}
		if attr.SingleNested != nil {
			updateNestedComputabilityAndReqBodyUsage(&attr.SingleNested.NestedObject.Attributes, attrIsComputed, attrIsOmittedInReqBody)
		}
		if attr.SetNested != nil {
			updateNestedComputabilityAndReqBodyUsage(&attr.SetNested.NestedObject.Attributes, attrIsComputed, attrIsOmittedInReqBody)
		}
		if attr.MapNested != nil {
			updateNestedComputabilityAndReqBodyUsage(&attr.MapNested.NestedObject.Attributes, attrIsComputed, attrIsOmittedInReqBody)
		}
	}
}

func hasDefault(attr *Attribute) bool {
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
