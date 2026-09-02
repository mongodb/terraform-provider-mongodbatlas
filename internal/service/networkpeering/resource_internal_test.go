package networkpeering

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
)

func TestFailedStatusDiagnostics(t *testing.T) {
	testCases := map[string]struct {
		peer         *admin.BaseNetworkPeeringConnectionSettings
		errDetail    string
		expectsDiags bool
	}{
		"AWS FAILED status": {
			peer: &admin.BaseNetworkPeeringConnectionSettings{
				StatusName:     new("FAILED"),
				ErrorStateName: new("REJECTED"),
			},
			errDetail:    "REJECTED",
			expectsDiags: true,
		},
		"AWS FAILED status with error message": {
			peer: &admin.BaseNetworkPeeringConnectionSettings{
				StatusName:     new("FAILED"),
				ErrorStateName: new("REJECTED"),
				ErrorMessage:   new("vpc peering connection request was rejected by AWS"),
			},
			errDetail:    "REJECTED: vpc peering connection request was rejected by AWS",
			expectsDiags: true,
		},
		"Azure/GCP FAILED status": {
			peer: &admin.BaseNetworkPeeringConnectionSettings{
				Status:     new("FAILED"),
				ErrorState: new("CIDR_BLOCK_CONFLICT"),
			},
			errDetail:    "CIDR_BLOCK_CONFLICT",
			expectsDiags: true,
		},
		"AVAILABLE status": {
			peer: &admin.BaseNetworkPeeringConnectionSettings{
				Status:     new("AVAILABLE"),
				StatusName: new("AVAILABLE"),
			},
			expectsDiags: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			errorDiags := errorIfFailedStatusIsPresent(tc.peer)
			warningDiags := warnIfFailedStatusIsPresent(tc.peer)
			if !tc.expectsDiags {
				assert.Empty(t, errorDiags)
				assert.Empty(t, warningDiags)
				return
			}
			assert.Equal(t, diag.FromErr(fmt.Errorf("peer networking is in a failed state: %s", tc.errDetail)), errorDiags)
			expectedWarning := diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "Network peering connection is in FAILED status",
				Detail:   fmt.Sprintf("Peer networking is in a failed state: %s. The resource is kept in the Terraform state. Fix the reported issue and recreate the resource if needed.", tc.errDetail),
			}}
			assert.Equal(t, expectedWarning, warningDiags)
		})
	}
}
