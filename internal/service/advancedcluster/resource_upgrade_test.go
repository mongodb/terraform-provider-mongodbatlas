package advancedcluster_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

func clusterDescription(providerName, majorVersion string) *admin.ClusterDescription20240805 {
	return &admin.ClusterDescription20240805{
		ClusterType:         new("REPLICASET"),
		MongoDBMajorVersion: optionalString(majorVersion),
		Name:                new("cluster"),
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{{
			RegionConfigs: &[]admin.CloudRegionConfig20240805{{
				ProviderName: new(providerName),
				RegionName:   new("US_EAST_1"),
				ElectableSpecs: &admin.HardwareSpec20240805{
					InstanceSize: new("M30"),
				},
			}},
		}},
	}
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return new(value)
}

func TestGetUpgradeFlexToDedicatedRequest(t *testing.T) {
	testCases := map[string]struct {
		stateProvider   string
		stateVersion    string
		planVersion     string
		expectedVersion string
		expectRequest   bool
	}{
		"sends version set in config": {
			expectRequest:   true,
			stateProvider:   flexcluster.FlexClusterType,
			planVersion:     "8.0",
			expectedVersion: "8.0",
		},
		"omits version unset in config": {
			expectRequest: true,
			stateProvider: flexcluster.FlexClusterType,
		},
		"sends version from config even when it matches a version in state": {
			expectRequest:   true,
			stateProvider:   flexcluster.FlexClusterType,
			stateVersion:    "8.0",
			planVersion:     "8.0",
			expectedVersion: "8.0",
		},
		"omits version when state has a version but config does not": {
			expectRequest: true,
			stateProvider: flexcluster.FlexClusterType,
			stateVersion:  "8.0",
		},
		"returns nil when not upgrading from flex": {
			stateProvider: constant.AWS,
			stateVersion:  "8.0",
			planVersion:   "8.0",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			state := clusterDescription(tc.stateProvider, tc.stateVersion)
			plan := clusterDescription(constant.AWS, tc.planVersion)
			// PatchPayload omits fields unchanged between state and plan, so the version is usually absent from the patch.
			patch := clusterDescription(constant.AWS, "")

			req := advancedcluster.GetUpgradeFlexToDedicatedRequest(state, plan, patch)

			if !tc.expectRequest {
				assert.Nil(t, req)
				return
			}
			require.NotNil(t, req)
			assert.Equal(t, optionalString(tc.expectedVersion), req.MongoDBMajorVersion)
		})
	}
}

func TestGetUpgradeTenantRequest(t *testing.T) {
	testCases := map[string]struct {
		stateProvider   string
		stateVersion    string
		planVersion     string
		expectedVersion string
		expectRequest   bool
	}{
		"sends version set in config": {
			expectRequest:   true,
			stateProvider:   constant.TENANT,
			planVersion:     "8.0",
			expectedVersion: "8.0",
		},
		"omits version unset in config": {
			expectRequest: true,
			stateProvider: constant.TENANT,
		},
		"sends version from config even when it matches a version in state": {
			expectRequest:   true,
			stateProvider:   constant.TENANT,
			stateVersion:    "8.0",
			planVersion:     "8.0",
			expectedVersion: "8.0",
		},
		"omits version when state has a version but config does not": {
			expectRequest: true,
			stateProvider: constant.TENANT,
			stateVersion:  "8.0",
		},
		"returns nil when not upgrading from tenant": {
			stateProvider: constant.AWS,
			stateVersion:  "8.0",
			planVersion:   "8.0",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			state := clusterDescription(tc.stateProvider, tc.stateVersion)
			plan := clusterDescription(constant.AWS, tc.planVersion)
			// PatchPayload omits fields unchanged between state and plan, so the version is usually absent from the patch.
			patch := clusterDescription(constant.AWS, "")

			req := advancedcluster.GetUpgradeTenantRequest(state, plan, patch)

			if !tc.expectRequest {
				assert.Nil(t, req)
				return
			}
			require.NotNil(t, req)
			assert.Equal(t, optionalString(tc.expectedVersion), req.MongoDBMajorVersion)
		})
	}
}
