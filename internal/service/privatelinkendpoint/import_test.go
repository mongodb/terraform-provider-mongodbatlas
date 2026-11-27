package privatelinkendpoint_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpoint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitPrivateLinkEndpointImportID(t *testing.T) {
	tests := map[string]struct {
		importID      string
		projectID     string
		privateLinkID string
		providerName  string
		region        string
		legacy        bool
		expectError   bool
	}{
		"slash format": {
			importID:      "project/link/AWS/us-east-1",
			projectID:     "project",
			privateLinkID: "link",
			providerName:  "AWS",
			region:        "us-east-1",
		},
		"legacy format with hyphen in region": {
			importID:      "project-link-AWS-us-east-1",
			projectID:     "project",
			privateLinkID: "link",
			providerName:  "AWS",
			region:        "us-east-1",
			legacy:        true,
		},
		"legacy format with several hyphens in region": {
			importID:      "project-link-AZURE-east-us-2",
			projectID:     "project",
			privateLinkID: "link",
			providerName:  "AZURE",
			region:        "east-us-2",
			legacy:        true,
		},
		"invalid format": {
			importID:    "project-link-AWS",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := privatelinkendpoint.SplitPrivateLinkEndpointImportID(tc.importID)
			if tc.expectError {
				require.Error(t, err)
				assert.False(t, result.Legacy)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.projectID, result.ProjectID)
			assert.Equal(t, tc.privateLinkID, result.PrivateLinkID)
			assert.Equal(t, tc.providerName, result.ProviderName)
			assert.Equal(t, tc.region, result.Region)
			assert.Equal(t, tc.legacy, result.Legacy)
		})
	}
}
