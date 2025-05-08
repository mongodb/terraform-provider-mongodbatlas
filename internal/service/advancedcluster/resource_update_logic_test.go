package advancedcluster_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

func TestSyncAutoScalingConfigs(t *testing.T) {
	testCases := map[string]struct {
		ReplicationSpecs         []admin.ReplicationSpec20240805
		ExpectedReplicationSpecs []admin.ReplicationSpec20240805
	}{
		"apply same autoscaling options for new replication spec which does not have autoscaling defined": {
			ReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					Id: admin.PtrString("id-1"),
					RegionConfigs: &[]admin.CloudRegionConfig20240805{
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
					RegionConfigs: &[]admin.CloudRegionConfig20240805{
						{
							AutoScaling:          nil,
							AnalyticsAutoScaling: nil,
						},
					},
				},
			},
			ExpectedReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					Id: admin.PtrString("id-1"),
					RegionConfigs: &[]admin.CloudRegionConfig20240805{
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
					RegionConfigs: &[]admin.CloudRegionConfig20240805{
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
			ReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					Id: admin.PtrString("id-1"),
					RegionConfigs: &[]admin.CloudRegionConfig20240805{
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
					RegionConfigs: &[]admin.CloudRegionConfig20240805{
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
			ExpectedReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					Id: admin.PtrString("id-1"),
					RegionConfigs: &[]admin.CloudRegionConfig20240805{
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
					RegionConfigs: &[]admin.CloudRegionConfig20240805{
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
