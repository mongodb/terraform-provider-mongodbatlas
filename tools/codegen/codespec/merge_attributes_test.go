package codespec_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

func TestMergeAttributes_requestPrimitiveListKeepsShapeOverResponseNestedList(t *testing.T) {
	request := codespec.Attribute{
		TFSchemaName:             "roles",
		ComputedOptionalRequired: codespec.Required,
		List:                     &codespec.ListAttribute{ElementType: codespec.String},
	}
	response := codespec.Attribute{
		TFSchemaName:             "roles",
		ComputedOptionalRequired: codespec.Computed,
		ListNested: &codespec.ListNestedAttribute{NestedObject: codespec.NestedAttributeObject{Attributes: codespec.Attributes{
			{TFSchemaName: "role_name", ComputedOptionalRequired: codespec.Computed, String: &codespec.StringAttribute{}},
		}}},
	}

	merged := codespec.MergeRequestAndResponseAttributesForTest(codespec.Attributes{request}, codespec.Attributes{response})

	assert.Len(t, merged, 1)
	assert.NotNil(t, merged[0].List)
	assert.Nil(t, merged[0].ListNested)
	assert.Equal(t, codespec.Required, merged[0].ComputedOptionalRequired)
	assert.True(t, merged[0].PresentInAnyResponse)
}
