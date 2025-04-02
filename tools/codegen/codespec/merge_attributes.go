package codespec

import (
	"sort"
)

// mergeNestedAttributes recursively merges nested attributes
func mergeNestedAttributes(existingAttrs *Attributes, newAttrs Attributes, computability ComputedOptionalRequired, reqBodyUsage AttributeReqBodyUsage, isFromResponse bool) {
	mergedMap := make(map[string]*Attribute)
	if existingAttrs != nil {
		for i := range *existingAttrs {
			mergedMap[(*existingAttrs)[i].Name.SnakeCase()] = &(*existingAttrs)[i]
		}
	}

	// add new attributes and merge when necessary
	for i := range newAttrs {
		newAttr := &newAttrs[i]

		if _, exists := mergedMap[newAttr.Name.SnakeCase()]; exists {
			addOrUpdate(newAttr, computability, reqBodyUsage, mergedMap, isFromResponse)
		} else {
			newAttr.ComputedOptionalRequired = computability
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

// addOrUpdate adds or updates an attribute in the merged map, including nested attributes
func addOrUpdate(attr *Attribute, computability ComputedOptionalRequired, reqBodyUsage AttributeReqBodyUsage, merged map[string]*Attribute, isFromResponse bool) {
	if existingAttr, exists := merged[attr.Name.SnakeCase()]; exists {
		if existingAttr.Description == nil || *existingAttr.Description == "" {
			existingAttr.Description = attr.Description
		}

		// retain existing computability if already set from request
		if !isFromResponse && existingAttr.ComputedOptionalRequired != Required {
			existingAttr.ComputedOptionalRequired = computability
		}

		// retain existing ReqBodyUsage if already set defined from request information
		if !isFromResponse {
			existingAttr.ReqBodyUsage = reqBodyUsage
		}

		// handle nested attributes
		if existingAttr.ListNested != nil && attr.ListNested != nil {
			mergeNestedAttributes(&existingAttr.ListNested.NestedObject.Attributes, attr.ListNested.NestedObject.Attributes, computability, reqBodyUsage, isFromResponse)
		} else if attr.ListNested != nil {
			existingAttr.ListNested = attr.ListNested
		}

		if existingAttr.SingleNested != nil && attr.SingleNested != nil {
			mergeNestedAttributes(&existingAttr.SingleNested.NestedObject.Attributes, attr.SingleNested.NestedObject.Attributes, computability, reqBodyUsage, isFromResponse)
		} else if attr.SingleNested != nil {
			existingAttr.SingleNested = attr.SingleNested
		}

		if existingAttr.SetNested != nil && attr.SetNested != nil {
			mergeNestedAttributes(&existingAttr.SetNested.NestedObject.Attributes, attr.SetNested.NestedObject.Attributes, computability, reqBodyUsage, isFromResponse)
		} else if attr.SetNested != nil {
			existingAttr.SetNested = attr.SetNested
		}

		if existingAttr.MapNested != nil && attr.MapNested != nil {
			mergeNestedAttributes(&existingAttr.MapNested.NestedObject.Attributes, attr.MapNested.NestedObject.Attributes, computability, reqBodyUsage, isFromResponse)
		} else if attr.MapNested != nil {
			existingAttr.MapNested = attr.MapNested
		}
	} else {
		// add new attribute with the given computability
		newAttr := *attr
		newAttr.ComputedOptionalRequired = computability
		newAttr.ReqBodyUsage = reqBodyUsage
		merged[attr.Name.SnakeCase()] = &newAttr
	}
}

// mergeAttributes merges attributes from different sources (path params, create/get operation bodies) and determines a single merged list of attributes.
// Computability and reqBodyUsage values are determined as part of this process.
func mergeAttributes(createPathParams, createRequest, createResponse, readResponse Attributes) Attributes {
	merged := make(map[string]*Attribute)

	// create path parameters: all attributes will be "required", reqBodyUsage is defined as omit all at this step
	for i := range createPathParams {
		addOrUpdate(&createPathParams[i], Required, OmitAll, merged, false)
	}

	// POST request body: optional/required is as defined, reqBodyUsage is defined as AllRequestBodies
	for i := range createRequest {
		// for now we do not differentiate AllRequestBodies vs PostBodyOnly as we are not processing update request
		addOrUpdate(&createRequest[i], createRequest[i].ComputedOptionalRequired, AllRequestBodies, merged, false)
	}

	// POST/GET response body: properties not in the request body are "computed" or "computed_optional" (if a default is present), reqBodyUsage will have OmitAll not present in request body
	for i := range createResponse {
		if hasDefault(&createResponse[i]) {
			addOrUpdate(&createResponse[i], ComputedOptional, OmitAll, merged, true)
		} else {
			addOrUpdate(&createResponse[i], Computed, OmitAll, merged, true)
		}
	}

	for i := range readResponse {
		if hasDefault(&readResponse[i]) {
			addOrUpdate(&readResponse[i], ComputedOptional, OmitAll, merged, true)
		} else {
			addOrUpdate(&readResponse[i], Computed, OmitAll, merged, true)
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
