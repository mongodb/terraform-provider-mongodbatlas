package networkpeering_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/networkpeering"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitNetworkPeerImportID(t *testing.T) {
	const (
		projectID = "664619d870c247237f4b86a6"
		peerID    = "664619d870c247237f4b86a7"
	)
	tests := map[string]struct {
		importID     string
		projectID    string
		peerID       string
		providerName string
		legacy       bool
		expectError  bool
	}{
		"slash format": {
			importID:     projectID + "/" + peerID + "/AWS",
			projectID:    projectID,
			peerID:       peerID,
			providerName: "AWS",
		},
		"legacy format": {
			importID:     projectID + "-" + peerID + "-AWS",
			projectID:    projectID,
			peerID:       peerID,
			providerName: "AWS",
			legacy:       true,
		},
		"invalid format": {
			importID:    "invalid",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			projectID, peerID, providerName, legacy, err := networkpeering.SplitNetworkPeerImportID(tc.importID)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.projectID, projectID)
			assert.Equal(t, tc.peerID, peerID)
			assert.Equal(t, tc.providerName, providerName)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
