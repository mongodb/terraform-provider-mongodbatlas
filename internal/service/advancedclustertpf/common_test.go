package advancedclustertpf_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.mongodb.org/atlas-sdk/v20250312003/admin"
)

func TestAddIDsToReplicationSpecs(t *testing.T) {
	testCases := map[string]struct {
		ReplicationSpecs          []admin.ReplicationSpec20240805
		ZoneToReplicationSpecsIDs map[string][]string
		ExpectedReplicationSpecs  []admin.ReplicationSpec20240805
	}{
		"two zones with same amount of available ids and replication specs to populate": {
			ReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: admin.PtrString("Zone 1"),
				},
				{
					ZoneName: admin.PtrString("Zone 2"),
				},
				{
					ZoneName: admin.PtrString("Zone 1"),
				},
				{
					ZoneName: admin.PtrString("Zone 2"),
				},
			},
			ZoneToReplicationSpecsIDs: map[string][]string{
				"Zone 1": {"zone1-id1", "zone1-id2"},
				"Zone 2": {"zone2-id1", "zone2-id2"},
			},
			ExpectedReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: admin.PtrString("Zone 1"),
					Id:       admin.PtrString("zone1-id1"),
				},
				{
					ZoneName: admin.PtrString("Zone 2"),
					Id:       admin.PtrString("zone2-id1"),
				},
				{
					ZoneName: admin.PtrString("Zone 1"),
					Id:       admin.PtrString("zone1-id2"),
				},
				{
					ZoneName: admin.PtrString("Zone 2"),
					Id:       admin.PtrString("zone2-id2"),
				},
			},
		},
		"less available ids than replication specs to populate": {
			ReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: admin.PtrString("Zone 1"),
				},
				{
					ZoneName: admin.PtrString("Zone 1"),
				},
				{
					ZoneName: admin.PtrString("Zone 1"),
				},
				{
					ZoneName: admin.PtrString("Zone 2"),
				},
			},
			ZoneToReplicationSpecsIDs: map[string][]string{
				"Zone 1": {"zone1-id1"},
				"Zone 2": {"zone2-id1"},
			},
			ExpectedReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: admin.PtrString("Zone 1"),
					Id:       admin.PtrString("zone1-id1"),
				},
				{
					ZoneName: admin.PtrString("Zone 1"),
					Id:       nil,
				},
				{
					ZoneName: admin.PtrString("Zone 1"),
					Id:       nil,
				},
				{
					ZoneName: admin.PtrString("Zone 2"),
					Id:       admin.PtrString("zone2-id1"),
				},
			},
		},
		"more available ids than replication specs to populate": {
			ReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: admin.PtrString("Zone 1"),
				},
				{
					ZoneName: admin.PtrString("Zone 2"),
				},
			},
			ZoneToReplicationSpecsIDs: map[string][]string{
				"Zone 1": {"zone1-id1", "zone1-id2"},
				"Zone 2": {"zone2-id1", "zone2-id2"},
			},
			ExpectedReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: admin.PtrString("Zone 1"),
					Id:       admin.PtrString("zone1-id1"),
				},
				{
					ZoneName: admin.PtrString("Zone 2"),
					Id:       admin.PtrString("zone2-id1"),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			resultSpecs := advancedclustertpf.AddIDsToReplicationSpecs(tc.ReplicationSpecs, tc.ZoneToReplicationSpecsIDs)
			assert.Equal(t, tc.ExpectedReplicationSpecs, resultSpecs)
		})
	}
}

func TestCleanupOnErrorSkippedWhenNoError(t *testing.T) {
	cleanupCalled := false
	cleanup := func(ctx context.Context) error {
		cleanupCalled = true
		return nil
	}
	advancedclustertpf.CleanupOnError(t.Context(), &diag.Diagnostics{}, "warning detail", cleanup, nil)
	assert.False(t, cleanupCalled, "cleanup should not be called")
}

func TestCleanupOnErrorCalledForAnError(t *testing.T) {
	cleanupCalled := false
	cleanupFailed := func(ctx context.Context) error {
		cleanupCalled = true
		return errors.New("cleanup failed")
	}
	diags := diag.Diagnostics{}
	diags.AddError("error", "handler error")
	advancedclustertpf.CleanupOnError(t.Context(), &diags, "warning detail", cleanupFailed, nil)
	assert.True(t, cleanupCalled, "cleanup should be called")
	assert.Len(t, diags, 3)
	assert.Equal(t, "Failed to create, will perform cleanup due to error", diags[1].Summary())
	assert.Equal(t, "warning detail", diags[1].Detail())
	assert.Equal(t, "Error during cleanup", diags[2].Summary())
	assert.Equal(t, "warning detail error=cleanup failed", diags[2].Detail())
}

func TestCleanupOnErrorCalledForATimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Millisecond)
	defer cancel()
	time.Sleep(2 * time.Millisecond)
	cleanupCalled := false
	finalContext := ctx
	cleanup := func(callbackCtx context.Context) error {
		cleanupCalled = true
		finalContext = callbackCtx
		return nil
	}
	diags := diag.Diagnostics{}
	diags.AddError("error", "timeout")
	advancedclustertpf.CleanupOnError(ctx, &diags, "warning detail", cleanup, nil)
	assert.True(t, cleanupCalled, "cleanup should be called")
	assert.NotEqual(t, finalContext, ctx, "cleanup should be called with a new context")
	require.NoError(t, finalContext.Err(), "cleanup should be called with a new context that hasn't been cancelled")
	assert.Len(t, diags, 2)
	assert.Equal(t, "Failed to create, will perform cleanup due to timeout reached", diags[1].Summary())
	assert.Equal(t, "warning detail", diags[1].Detail())
}
