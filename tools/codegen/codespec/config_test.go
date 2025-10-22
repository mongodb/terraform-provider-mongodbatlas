package codespec_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/stretchr/testify/assert"
)

func TestApplyTimeoutTransformation(t *testing.T) {
	tests := map[string]struct {
		inputOperations  codespec.APIOperations
		expectedTimeouts []codespec.Operation
	}{
		"No wait blocks - no timeout attribute added": {
			inputOperations: codespec.APIOperations{
				Create: codespec.APIOperation{},
				Read:   codespec.APIOperation{},
				Update: codespec.APIOperation{},
			},
			expectedTimeouts: nil,
		},
		"Create wait only": {
			inputOperations: codespec.APIOperations{
				Create: codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Read:   codespec.APIOperation{},
				Update: codespec.APIOperation{},
			},
			expectedTimeouts: []codespec.Operation{codespec.Create},
		},
		"Create, Update, Delete waits": {
			inputOperations: codespec.APIOperations{
				Create: codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Read: codespec.APIOperation{},
				Update: codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Delete: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
			},
			expectedTimeouts: []codespec.Operation{codespec.Create, codespec.Update, codespec.Delete},
		},
		"All operations with waits": {
			inputOperations: codespec.APIOperations{
				Create: codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Read: codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Update: codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
				Delete: &codespec.APIOperation{
					Wait: &codespec.Wait{},
				},
			},
			expectedTimeouts: []codespec.Operation{codespec.Create, codespec.Update, codespec.Read, codespec.Delete},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resource := &codespec.Resource{
				Name: "test_resource",
				Schema: &codespec.Schema{
					Attributes: []codespec.Attribute{},
				},
				Operations: tc.inputOperations,
			}
			codespec.ApplyTimeoutTransformation(resource)
			if tc.expectedTimeouts == nil {
				assert.Empty(t, resource.Schema.Attributes)
			} else {
				assert.Len(t, resource.Schema.Attributes, 1)
				expectedAttr := codespec.Attribute{
					TFSchemaName: "timeouts",
					TFModelName:  "Timeouts",
					ReqBodyUsage: codespec.OmitAlways,
					Timeouts:     &codespec.TimeoutsAttribute{ConfigurableTimeouts: tc.expectedTimeouts},
				}
				assert.Equal(t, expectedAttr, resource.Schema.Attributes[0])
			}
		})
	}
}
