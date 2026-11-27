package privatelinkendpointservice_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpointservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitImportStateID(t *testing.T) {
	const (
		projectID = "664619d870c247237f4b86a6"
		linkID    = "3242342343112"
		azureID   = "/subscriptions/123/resourceGroups/rg/providers/Microsoft.Network/privateEndpoints/my-pe"
	)
	tests := map[string]struct {
		importID          string
		projectID         string
		privateLinkID     string
		endpointServiceID string
		providerName      string
		legacy            bool
		expectError       bool
	}{
		"slash format aws": {
			importID:          projectID + "/" + linkID + "/vpce-123/AWS",
			projectID:         projectID,
			privateLinkID:     linkID,
			endpointServiceID: "vpce-123",
			providerName:      "AWS",
		},
		"slash format azure with slashes in endpoint service id": {
			importID:          projectID + "/" + linkID + "/" + azureID + "/AZURE",
			projectID:         projectID,
			privateLinkID:     linkID,
			endpointServiceID: azureID,
			providerName:      "AZURE",
		},
		"legacy format aws": {
			importID:          projectID + "--" + linkID + "--vpce-123--AWS",
			projectID:         projectID,
			privateLinkID:     linkID,
			endpointServiceID: "vpce-123",
			providerName:      "AWS",
			legacy:            true,
		},
		"legacy format azure with slashes in endpoint service id": {
			importID:          projectID + "--" + linkID + "--" + azureID + "--AZURE",
			projectID:         projectID,
			privateLinkID:     linkID,
			endpointServiceID: azureID,
			providerName:      "AZURE",
			legacy:            true,
		},
		"invalid slash format without provider name": {
			importID:    projectID + "/" + linkID + "/vpce-123",
			expectError: true,
		},
		"invalid format": {
			importID:    "invalid",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			importID, err := privatelinkendpointservice.SplitImportStateID(tc.importID)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.projectID, importID.ProjectID)
			assert.Equal(t, tc.privateLinkID, importID.PrivateLinkID)
			assert.Equal(t, tc.endpointServiceID, importID.EndpointServiceID)
			assert.Equal(t, tc.providerName, importID.ProviderName)
			assert.Equal(t, tc.legacy, importID.Legacy)
		})
	}
}
