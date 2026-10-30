package codespec

// Test helpers exported only for package codespec_test (see merge_attributes_test.go).

func MergeRequestAndResponseAttributesForTest(createRequest, readResponse Attributes) Attributes {
	return mergeAttributes(&attributeDefinitionSources{createRequest: createRequest, readResponse: readResponse})
}
