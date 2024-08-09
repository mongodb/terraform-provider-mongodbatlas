package advancedcluster_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

func TestAddIDsToReplicationSpecs(t *testing.T) {
	testCases := map[string]struct {
		ReplicationSpecs          []admin.ReplicationSpec20250101
		ZoneToReplicationSpecsIDs map[string][]string
		ExpectedReplicationSpecs  []admin.ReplicationSpec20250101
	}{
		"two zones with same amount of available ids and replication specs to populate": {
			ReplicationSpecs: []admin.ReplicationSpec20250101{
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
			ExpectedReplicationSpecs: []admin.ReplicationSpec20250101{
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
			ReplicationSpecs: []admin.ReplicationSpec20250101{
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
			ExpectedReplicationSpecs: []admin.ReplicationSpec20250101{
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
			ReplicationSpecs: []admin.ReplicationSpec20250101{
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
			ExpectedReplicationSpecs: []admin.ReplicationSpec20250101{
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
			resultSpecs := advancedcluster.AddIDsToReplicationSpecs(tc.ReplicationSpecs, tc.ZoneToReplicationSpecsIDs)
			assert.Equal(t, tc.ExpectedReplicationSpecs, resultSpecs)
		})
	}
}

func TestSyncAutoScalingConfigs(t *testing.T) {
	testCases := map[string]struct {
		ReplicationSpecs         []admin.ReplicationSpec20250101
		ExpectedReplicationSpecs []admin.ReplicationSpec20250101
	}{
		"apply same autoscaling options for new replication spec which does not have autoscaling defined": {
			ReplicationSpecs: []admin.ReplicationSpec20250101{
				{
					Id: admin.PtrString("id-1"),
					RegionConfigs: &[]admin.CloudRegionConfig20250101{
						{
							AutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(false),
									ScaleDownEnabled: admin.PtrBool(false),
								},
							},
							AnalyticsAutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(false),
									ScaleDownEnabled: admin.PtrBool(false),
								},
							},
						},
					},
				},
				{
					Id: admin.PtrString("id-2"),
					RegionConfigs: &[]admin.CloudRegionConfig20250101{
						{
							AutoScaling:          nil,
							AnalyticsAutoScaling: nil,
						},
					},
				},
			},
			ExpectedReplicationSpecs: []admin.ReplicationSpec20250101{
				{
					Id: admin.PtrString("id-1"),
					RegionConfigs: &[]admin.CloudRegionConfig20250101{
						{
							AutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(false),
									ScaleDownEnabled: admin.PtrBool(false),
								},
							},
							AnalyticsAutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(false),
									ScaleDownEnabled: admin.PtrBool(false),
								},
							},
						},
					},
				},
				{
					Id: admin.PtrString("id-2"),
					RegionConfigs: &[]admin.CloudRegionConfig20250101{
						{
							AutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(false),
									ScaleDownEnabled: admin.PtrBool(false),
								},
							},
							AnalyticsAutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(false),
									ScaleDownEnabled: admin.PtrBool(false),
								},
							},
						},
					},
				},
			},
		},
		// for this case the API will respond with an error and guide the user to align autoscaling options cross all nodes
		"when different autoscaling options are defined values will not be changed": {
			ReplicationSpecs: []admin.ReplicationSpec20250101{
				{
					Id: admin.PtrString("id-1"),
					RegionConfigs: &[]admin.CloudRegionConfig20250101{
						{
							AutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(false),
									ScaleDownEnabled: admin.PtrBool(false),
								},
							},
							AnalyticsAutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(true),
									ScaleDownEnabled: admin.PtrBool(true),
								},
							},
						},
					},
				},
				{
					Id: admin.PtrString("id-2"),
					RegionConfigs: &[]admin.CloudRegionConfig20250101{
						{
							AutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled: admin.PtrBool(true),
								},
							},
							AnalyticsAutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled: admin.PtrBool(false),
								},
							},
						},
					},
				},
			},
			ExpectedReplicationSpecs: []admin.ReplicationSpec20250101{
				{
					Id: admin.PtrString("id-1"),
					RegionConfigs: &[]admin.CloudRegionConfig20250101{
						{
							AutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(false),
									ScaleDownEnabled: admin.PtrBool(false),
								},
							},
							AnalyticsAutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled:          admin.PtrBool(true),
									ScaleDownEnabled: admin.PtrBool(true),
								},
							},
						},
					},
				},
				{
					Id: admin.PtrString("id-2"),
					RegionConfigs: &[]admin.CloudRegionConfig20250101{
						{
							AutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled: admin.PtrBool(true),
								},
							},
							AnalyticsAutoScaling: &admin.AdvancedAutoScalingSettings{
								Compute: &admin.AdvancedComputeAutoScaling{
									Enabled: admin.PtrBool(false),
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			specs := &tc.ReplicationSpecs
			advancedcluster.SyncAutoScalingConfigs(specs)
			assert.Equal(t, tc.ExpectedReplicationSpecs, *specs)
		})
	}
}
