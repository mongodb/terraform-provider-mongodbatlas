package eventtrigger_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/eventtrigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitEventTriggerImportID(t *testing.T) {
	const projectID = "664619d870c247237f4b86a6"
	tests := map[string]struct {
		importID    string
		parts       []string
		legacy      bool
		expectError bool
	}{
		"slash format": {
			importID: projectID + "/app-id/trigger-id",
			parts:    []string{projectID, "app-id", "trigger-id"},
		},
		"legacy format": {
			importID: projectID + "--app-id--trigger-id",
			parts:    []string{projectID, "app-id", "trigger-id"},
			legacy:   true,
		},
		"invalid format": {
			importID:    "invalid",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parts, legacy, err := eventtrigger.SplitEventTriggerImportID(tc.importID)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.parts, parts)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
