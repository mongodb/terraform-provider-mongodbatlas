package advancedclustertpf_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/stretchr/testify/assert"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"
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
