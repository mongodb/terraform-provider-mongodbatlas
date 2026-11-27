package streamconnection_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamconnection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitStreamConnectionImportID(t *testing.T) {
	const projectID = "664619d870c247237f4b86a6"
	tests := map[string]struct {
		importID       string
		workspaceName  string
		projectID      string
		connectionName string
		legacy         bool
		expectError    bool
	}{
		"slash format": {
			importID:       projectID + "/workspace/connection",
			workspaceName:  "workspace",
			projectID:      projectID,
			connectionName: "connection",
		},
		"legacy format": {
			importID:       "workspace-" + projectID + "-connection",
			workspaceName:  "workspace",
			projectID:      projectID,
			connectionName: "connection",
			legacy:         true,
		},
		"invalid format": {
			importID:    "invalid",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			workspaceName, gotProjectID, connectionName, legacy, err := streamconnection.SplitStreamConnectionImportID(tc.importID)
			if tc.expectError {
				require.Error(t, err)
				assert.False(t, legacy)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.workspaceName, workspaceName)
			assert.Equal(t, tc.projectID, gotProjectID)
			assert.Equal(t, tc.connectionName, connectionName)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
