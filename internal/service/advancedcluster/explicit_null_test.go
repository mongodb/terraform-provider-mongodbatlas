package advancedcluster_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312025/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
)

// TestExplicitNullStorageConfig verifies that setting storageConfig to explicit null
// creates a detectable change in the PATCH payload without needing ForceUpdateAttr.
func TestExplicitNullStorageConfig(t *testing.T) {
	// State: has storageConfig with shardSizeLimitGB
	stateReq := &admin.ClusterDescription20240805{
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{
			{
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						ProviderName: new("AWS"),
						RegionName:   new("US_EAST_1"),
						Priority:     new(7),
						AutoScaling: &admin.AdvancedAutoScalingSettings{
							Compute: &admin.AdvancedComputeAutoScaling{
								Enabled: new(false),
							},
							StorageConfig: &admin.StorageConfig{
								ShardSizeLimitGB: new(1024),
							},
						},
					},
				},
			},
		},
	}

	// Plan: has storageConfig explicitly set to null
	planReq := &admin.ClusterDescription20240805{
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{
			{
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						ProviderName: new("AWS"),
						RegionName:   new("US_EAST_1"),
						Priority:     new(7),
						AutoScaling: &admin.AdvancedAutoScalingSettings{
							Compute: &admin.AdvancedComputeAutoScaling{
								Enabled: new(false),
							},
						},
					},
				},
			},
		},
	}

	// Set storageConfig to explicit null in the plan
	planRegion := &(*planReq.ReplicationSpecs)[0].GetRegionConfigs()[0]
	planRegion.AutoScaling.SetStorageConfigNil()

	// Verify the plan serializes with explicit null
	planJSON, err := json.Marshal(planReq)
	require.NoError(t, err)
	require.Contains(t, string(planJSON), `"storageConfig":null`, "plan should serialize with explicit null")

	// Verify the diff detects a change (replacement, not removal)
	patchReq, err := update.PatchPayload(stateReq, planReq)
	require.NoError(t, err)
	require.NotNil(t, patchReq, "patch should not be nil when storageConfig is explicitly set to null")
	require.NotNil(t, patchReq.ReplicationSpecs, "patch should include replicationSpecs")

	// Verify the PATCH request has the correct structure
	patchJSON, err := json.Marshal(patchReq)
	require.NoError(t, err)
	t.Logf("PATCH request: %s", string(patchJSON))

	// The PATCH should include replicationSpecs with autoScaling
	require.Len(t, patchReq.GetReplicationSpecs(), 1)
	require.Len(t, patchReq.GetReplicationSpecs()[0].GetRegionConfigs(), 1)
	patchRegion := patchReq.GetReplicationSpecs()[0].GetRegionConfigs()[0]
	require.NotNil(t, patchRegion.AutoScaling, "PATCH should include autoScaling")

	// Note: The unmarshaling may lose the explicit null marker, but that's OK because
	// Atlas treats omitted storageConfig the same as null when replicationSpecs is present.
	// The important thing is that replicationSpecs is included in the PATCH.
}
