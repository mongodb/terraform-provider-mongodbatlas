package streamprocessor_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitStreamProcessorImportID(t *testing.T) {
	const projectID = "664619d870c247237f4b86a6"
	tests := map[string]struct {
		importID      string
		projectID     string
		workspaceName string
		processorName string
		legacy        bool
		expectError   bool
	}{
		"slash format": {
			importID:      projectID + "/workspace/processor",
			projectID:     projectID,
			workspaceName: "workspace",
			processorName: "processor",
		},
		"legacy format": {
			importID:      "workspace-" + projectID + "-processor",
			projectID:     projectID,
			workspaceName: "workspace",
			processorName: "processor",
			legacy:        true,
		},
		"invalid format": {
			importID:    "invalid",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotProjectID, workspaceName, processorName, legacy, err := streamprocessor.SplitStreamProcessorImportID(tc.importID)
			if tc.expectError {
				require.Error(t, err)
				assert.False(t, legacy)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotProjectID)
			require.NotNil(t, workspaceName)
			require.NotNil(t, processorName)
			assert.Equal(t, tc.projectID, *gotProjectID)
			assert.Equal(t, tc.workspaceName, *workspaceName)
			assert.Equal(t, tc.processorName, *processorName)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
