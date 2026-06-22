package codespec_test

import (
	"fmt"
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

const openAPIDocTemplate = `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
components:
  schemas:
    Referenced:
      type: string
      description: original description
    Other:
      type: integer
    TestSchema:
%s`

func TestBuildSchema(t *testing.T) {
	tests := map[string]struct {
		schemaDefinition    string
		errorContains       string
		expectedType        string
		expectedDescription string
	}{
		"Explicit object type": {
			schemaDefinition: `      type: object
      properties:
        name:
          type: string`,
			expectedType: "object",
		},
		"Infer object type from properties": {
			schemaDefinition: `      properties:
        name:
          type: string`,
			expectedType: "object",
		},
		"Array without explicit type returns error": {
			schemaDefinition: `      items:
        type: string`,
			errorContains: "type cannot be inferred",
		},
		"Empty schema returns error": {
			schemaDefinition: `      description: A schema with no type or properties`,
			errorContains:    "type cannot be inferred",
		},
		"Nested schema with missing type": {
			schemaDefinition: `      properties:
        name:
          type: string
        mappings:
          type: object`,
			expectedType: "object",
		},
		"String type": {
			schemaDefinition: `      type: string`,
			expectedType:     "string",
		},
		"Integer type": {
			schemaDefinition: `      type: integer`,
			expectedType:     "integer",
		},
		"Explicit array type": {
			schemaDefinition: `      type: array
      items:
        type: string`,
			expectedType: "array",
		},
		"$ref with sibling description is unwrapped to referenced schema": {
			// libopenapi >= v0.36.4 turns `$ref` with a sibling property into a 2-branch allOf
			// wrapper. BuildSchema should unwrap that wrapper and surface the referenced
			// schema's type while honoring the sibling description override.
			schemaDefinition: `      $ref: '#/components/schemas/Referenced'
      description: override description`,
			expectedType:        "string",
			expectedDescription: "override description",
		},
		"allOf with two $refs is not treated as sibling-ref wrapper": {
			// Two $ref branches is a real allOf composition (not the libopenapi sibling-ref
			// shape) and should fall through to the existing error, not be silently unwrapped
			// to one of the refs.
			schemaDefinition: `      allOf:
        - $ref: '#/components/schemas/Referenced'
        - $ref: '#/components/schemas/Other'`,
			errorContains: "type cannot be inferred",
		},
		"$ref with sibling beyond description is not unwrapped": {
			// Only `description` siblings are recognized today. Any other sibling (e.g.
			// `deprecated`) must cause the wrapper to fall through to the existing error so
			// the new pattern is surfaced loudly instead of silently dropping the override.
			schemaDefinition: `      $ref: '#/components/schemas/Referenced'
      description: override description
      deprecated: true`,
			errorContains: "type cannot be inferred",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			openAPIDoc := fmt.Sprintf(openAPIDocTemplate, tc.schemaDefinition)
			doc, err := libopenapi.NewDocument([]byte(openAPIDoc))
			require.NoError(t, err)
			model, err := doc.BuildV3Model()
			require.NoError(t, err)
			schemaProxy := model.Model.Components.Schemas.GetOrZero("TestSchema")
			require.NotNil(t, schemaProxy)
			result, err := codespec.BuildSchema(schemaProxy)
			if tc.errorContains == "" {
				require.NoError(t, err)
				assert.NotNil(t, result.Schema)
				assert.Equal(t, tc.expectedType, result.Type)
				if tc.expectedDescription != "" {
					assert.Equal(t, tc.expectedDescription, result.Schema.Description)
				}
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
				assert.Nil(t, result)
			}
		})
	}
}

func TestGetXGenArraySemantic(t *testing.T) {
	tests := map[string]struct {
		expectedValue    *string
		schemaDefinition string
		expectInvalid    bool
	}{
		"absent extension returns nil": {
			schemaDefinition: `      type: array
      items:
        type: string`,
			expectedValue: nil,
		},
		"set value": {
			schemaDefinition: `      type: array
      x-xgen-array-semantic: set
      items:
        type: string`,
			expectedValue: conversion.StringPtr("set"),
		},
		"list value": {
			schemaDefinition: `      type: array
      x-xgen-array-semantic: list
      items:
        type: string`,
			expectedValue: conversion.StringPtr("list"),
		},
		"invalid value returns sentinel error": {
			schemaDefinition: `      type: array
      x-xgen-array-semantic: unordered
      items:
        type: string`,
			expectInvalid: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			openAPIDoc := fmt.Sprintf(openAPIDocTemplate, tc.schemaDefinition)
			doc, err := libopenapi.NewDocument([]byte(openAPIDoc))
			require.NoError(t, err)
			model, err := doc.BuildV3Model()
			require.NoError(t, err)
			schemaProxy := model.Model.Components.Schemas.GetOrZero("TestSchema")
			require.NotNil(t, schemaProxy)
			result, err := codespec.BuildSchema(schemaProxy)
			require.NoError(t, err)

			semantic, err := result.GetXGenArraySemantic()
			if tc.expectInvalid {
				require.ErrorIs(t, err, codespec.ErrInvalidArraySemantic)
				assert.Nil(t, semantic)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedValue, semantic)
		})
	}
}

func TestGetXGenServerComputedWhenClientOmitted(t *testing.T) {
	tests := map[string]struct {
		schemaDefinition string
		expected         bool
	}{
		"absent extension returns false": {
			schemaDefinition: `      type: string`,
			expected:         false,
		},
		"true value": {
			schemaDefinition: `      type: string
      x-xgen-server-computed-when-client-omitted: true`,
			expected: true,
		},
		"false value": {
			schemaDefinition: `      type: string
      x-xgen-server-computed-when-client-omitted: false`,
			expected: false,
		},
		"non-boolean value returns false": {
			schemaDefinition: `      type: string
      x-xgen-server-computed-when-client-omitted: maybe`,
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			openAPIDoc := fmt.Sprintf(openAPIDocTemplate, tc.schemaDefinition)
			doc, err := libopenapi.NewDocument([]byte(openAPIDoc))
			require.NoError(t, err)
			model, err := doc.BuildV3Model()
			require.NoError(t, err)
			schemaProxy := model.Model.Components.Schemas.GetOrZero("TestSchema")
			require.NotNil(t, schemaProxy)
			result, err := codespec.BuildSchema(schemaProxy)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, result.GetXGenServerComputedWhenClientOmitted())
		})
	}
}
