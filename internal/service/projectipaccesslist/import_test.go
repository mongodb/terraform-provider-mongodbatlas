package projectipaccesslist_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectipaccesslist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitProjectIPAccessListImportID(t *testing.T) {
	tests := map[string]struct {
		importID    string
		projectID   string
		entry       string
		legacy      bool
		expectError bool
	}{
		"slash format with CIDR entry": {
			importID:  "664619d870c247237f4b86a6/10.0.0.0/8",
			projectID: "664619d870c247237f4b86a6",
			entry:     "10.0.0.0/8",
		},
		"legacy format with CIDR entry": {
			importID:  "664619d870c247237f4b86a6-10.0.0.0/8",
			projectID: "664619d870c247237f4b86a6",
			entry:     "10.0.0.0/8",
			legacy:    true,
		},
		"invalid format": {
			importID:    "project",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			projectID, entry, legacy, err := projectipaccesslist.SplitProjectIPAccessListImportID(tc.importID)
			if tc.expectError {
				require.Error(t, err)
				assert.False(t, legacy)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.projectID, projectID)
			assert.Equal(t, tc.entry, entry)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
