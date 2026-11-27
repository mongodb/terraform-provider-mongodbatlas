package flexcluster_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitFlexClusterImportID(t *testing.T) {
	const projectID = "664619d870c247237f4b86a6"
	tests := map[string]struct {
		importID    string
		projectID   string
		name        string
		legacy      bool
		expectError bool
	}{
		"slash format": {
			importID:  projectID + "/my-cluster",
			projectID: projectID,
			name:      "my-cluster",
		},
		"legacy format": {
			importID:  projectID + "-my-cluster",
			projectID: projectID,
			name:      "my-cluster",
			legacy:    true,
		},
		"invalid format": {
			importID:    "invalid",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			projectID, clusterName, legacy, err := flexcluster.SplitFlexClusterImportID(tc.importID)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, projectID)
			require.NotNil(t, clusterName)
			assert.Equal(t, tc.projectID, *projectID)
			assert.Equal(t, tc.name, *clusterName)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
