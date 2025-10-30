package codespec_test

import (
	"fmt"
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

const openAPIDocTemplate = `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
components:
  schemas:
    TestSchema:
%s`

func TestBuildSchema(t *testing.T) {
	tests := map[string]struct {
		schemaDefinition string
		errorContains    string
		expectedType     string
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			openAPIDoc := fmt.Sprintf(openAPIDocTemplate, tc.schemaDefinition)
			doc, err := libopenapi.NewDocument([]byte(openAPIDoc))
			require.NoError(t, err)
			model, errs := doc.BuildV3Model()
			require.Empty(t, errs)
			schemaProxy := model.Model.Components.Schemas.GetOrZero("TestSchema")
			require.NotNil(t, schemaProxy)
			result, err := codespec.BuildSchema(schemaProxy)
			if tc.errorContains == "" {
				require.NoError(t, err)
				assert.NotNil(t, result.Schema)
				assert.Equal(t, tc.expectedType, result.Type)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
				assert.Nil(t, result)
			}
		})
	}
}
