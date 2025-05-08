package searchindex_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchindex"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

func TestUnmarshalSearchIndexAnalyzersFields(t *testing.T) {
	tc := map[string]struct {
		input             string
		expected          []admin.AtlasSearchAnalyzer
		expectedHasErrors bool
	}{
		"empty string returns nil not empty slice": {
			input:    "",
			expected: nil,
		},
		"valid input": {
			input: `
				[{
					"name": "index_analyzer_test_name",
					"charFilters": [{
						"type": "mapping",
						"mappings": {"\\" : "/"}
					}],
					"tokenizer": {
						"type": "nGram",
						"minGram": 2,
						"maxGram": 5
					},
					"tokenFilters": [{
						"type": "length",
						"min": 20,
						"max": 33
					}]
				}]
			`,
			expected: []admin.AtlasSearchAnalyzer{
				{
					Name: "index_analyzer_test_name",
					CharFilters: &[]any{
						map[string]any{
							"type": "mapping",
							"mappings": map[string]any{
								"\\": "/",
							},
						},
					},
					Tokenizer: map[string]any{
						"type":    "nGram",
						"minGram": float64(2),
						"maxGram": float64(5),
					},
					TokenFilters: &[]any{
						map[string]any{
							"type": "length",
							"min":  float64(20),
							"max":  float64(33),
						},
					},
				},
			},
		},
		"invalid input": {
			input:             "{bad json format",
			expectedHasErrors: true,
		},
	}
	for name, tc := range tc {
		t.Run(name, func(t *testing.T) {
			actual, diags := searchindex.UnmarshalSearchIndexAnalyzersFields(tc.input)
			assert.Equal(t, tc.expectedHasErrors, diags.HasError())
			if !diags.HasError() {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
