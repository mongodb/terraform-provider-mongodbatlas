package codespec

import "sort"

func mergeAttributes(pathParams, createRequest, createResponse, readResponse Attributes) Attributes {
	merged := make(map[string]*Attribute)

	addOrUpdate := func(attr *Attribute, computability ComputedOptionalRequired) {
		if existingAttr, exists := merged[attr.Name.SnakeCase()]; exists {
			if existingAttr.Description == nil || *existingAttr.Description == "" {
				existingAttr.Description = attr.Description
			}

			if existingAttr.ComputedOptionalRequired == Required {
				return
			}

			existingAttr.ComputedOptionalRequired = computability
		} else {
			newAttr := *attr
			newAttr.ComputedOptionalRequired = computability
			merged[attr.Name.SnakeCase()] = &newAttr
		}
	}

	// Path parameters: all attributes will be "required"
	for i := range pathParams {
		addOrUpdate(&pathParams[i], Required)
	}

	// POST request body: optional/required is as defined
	for i := range createRequest {
		addOrUpdate(&createRequest[i], createRequest[i].ComputedOptionalRequired)
	}

	// POST/GET response body: properties not in the request body are "computed" or "computed_optional" (if a default is present)
	for i := range createResponse {
		if _, exists := merged[createResponse[i].Name.SnakeCase()]; !exists {
			if isOptional(&createResponse[i]) {
				addOrUpdate(&createResponse[i], ComputedOptional)
			} else {
				addOrUpdate(&createResponse[i], Computed)
			}
		}
	}

	for i := range readResponse {
		if _, exists := merged[readResponse[i].Name.SnakeCase()]; !exists {
			if isOptional(&readResponse[i]) {
				addOrUpdate(&readResponse[i], ComputedOptional)
			} else {
				addOrUpdate(&readResponse[i], Computed)
			}
		}
	}

	resourceAttributes := make(Attributes, 0, len(merged))
	for _, attr := range merged {
		resourceAttributes = append(resourceAttributes, *attr)
	}

	sortAttributes(resourceAttributes)

	return resourceAttributes
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
