package searchindex_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchindex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
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

func TestSplitImportID(t *testing.T) {
	tests := map[string]struct {
		importID    string
		projectID   string
		clusterName string
		indexID     string
		errorString string
	}{
		"valid input": {
			importID:    "664619d870c247237f4b86a6/my-cluster/664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "my-cluster",
			indexID:     "664619d870c247237f4b86a7",
		},
		"valid input with double hyphen in cluster name": {
			importID:    "664619d870c247237f4b86a6/my--cluster/664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "my--cluster",
			indexID:     "664619d870c247237f4b86a7",
		},
		"valid input legacy": {
			importID:    "664619d870c247237f4b86a6--my-cluster--664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "my-cluster",
			indexID:     "664619d870c247237f4b86a7",
		},
		"valid input legacy with double hyphen in cluster name": {
			importID:    "664619d870c247237f4b86a6--my--cluster--664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "my--cluster",
			indexID:     "664619d870c247237f4b86a7",
		},
		"invalid input with more parts": {
			importID:    "part1/part2/part3/part4",
			errorString: searchindex.ErrorImportFormat,
		},
		"invalid input with less parts": {
			importID:    "part1/part2",
			errorString: searchindex.ErrorImportFormat,
		},
		"invalid input legacy with less parts": {
			importID:    "664619d870c247237f4b86a6--my-cluster",
			errorString: searchindex.ErrorImportFormat,
		},
		"invalid input legacy with empty part": {
			importID:    "664619d870c247237f4b86a6--my-cluster--",
			errorString: searchindex.ErrorImportFormat,
		},
		"empty input": {
			importID:    "",
			errorString: searchindex.ErrorImportFormat,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			projectID, clusterName, indexID, err := searchindex.SplitImportID(tc.importID)
			if tc.errorString != "" {
				require.EqualError(t, err, tc.errorString)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.projectID, projectID)
			assert.Equal(t, tc.clusterName, clusterName)
			assert.Equal(t, tc.indexID, indexID)
		})
	}
}
