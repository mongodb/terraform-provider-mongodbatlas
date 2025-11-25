package conversion_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateProjectID(t *testing.T) {
	require.NoError(t, conversion.ValidateProjectID("123456789012345678901234"))
	err := conversion.ValidateProjectID("invalid_project_id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid_project_id")
	assert.Contains(t, err.Error(), "project_id must be a 24 character hex string")
}

func TestValidateClusterName(t *testing.T) {
	require.NoError(t, conversion.ValidateClusterName("clusterName"))
	err := conversion.ValidateClusterName("_invalidClusterName")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "_invalidClusterName")
	assert.Contains(t, err.Error(), "cluster_name must be a string with length between 1 and 64, starting and ending with an alphanumeric character, and containing only alphanumeric characters and hyphens")
}

func TestImportSplit(t *testing.T) {
	tests := map[string]struct {
		importRaw     string
		parts         []string
		expectedParts int
		expected      bool
	}{
		"valid input": {
			importRaw:     "part1/part2/part3",
			expectedParts: 3,
			expected:      true,
			parts:         []string{"part1", "part2", "part3"},
		},
		"invalid input with more parts than expected": {
			importRaw:     "part1/part2/part3/part4",
			expectedParts: 3,
			expected:      false,
			parts:         nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ok, parts := conversion.ImportSplit(tc.importRaw, tc.expectedParts)
			if tc.expected {
				assert.Len(t, parts, tc.expectedParts)
			}
			assert.Equal(t, tc.expected, ok)
			assert.Equal(t, tc.parts, parts)
		})
	}
}
