package advancedcluster_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
	"go.mongodb.org/atlas-sdk/v20250312024/mockadmin"
)

var (
	dummyProjectID   = "projectId"
	errGeneric       = errors.New("generic")
	advancedClusters = []admin.ClusterDescription20240805{{StateName: conversion.StringPtr("NOT IDLE")}}
)

type Result struct {
	response any
	error    error
	state    string
}

func TestResourceListAdvancedRefreshFunc(t *testing.T) {
	testCases := []struct {
		mockCluster    *admin.PaginatedClusterDescription20240805
		mockResponse   *http.Response
		expectedResult Result
		mockError      error
		name           string
		expectedError  bool
	}{
		{
			name:          "Error in the API call: reset by peer",
			mockError:     errors.New("reset by peer"),
			expectedError: false,
			expectedResult: Result{
				response: nil,
				state:    "REPEATING",
				error:    nil,
			},
		},
		{
			name:          "Generic error in the API call",
			mockError:     errGeneric,
			expectedError: true,
			expectedResult: Result{
				response: nil,
				state:    "",
				error:    errGeneric,
			},
		},
		{
			name:          "Error in the API call: HTTP 404",
			mockError:     errGeneric,
			mockResponse:  &http.Response{StatusCode: http.StatusNotFound},
			expectedError: false,
			expectedResult: Result{
				response: "",
				state:    "DELETED",
				error:    nil,
			},
		},
		{
			name:          "Error in the API call: HTTP 503",
			mockError:     errGeneric,
			mockResponse:  &http.Response{StatusCode: http.StatusServiceUnavailable},
			expectedError: false,
			expectedResult: Result{
				response: "",
				state:    "PENDING",
				error:    nil,
			},
		},
		{
			name:          "Error in the API call: Neither HTTP 503 or 404",
			mockError:     errGeneric,
			mockResponse:  &http.Response{StatusCode: http.StatusBadRequest},
			expectedError: true,
			expectedResult: Result{
				response: nil,
				state:    "",
				error:    errGeneric,
			},
		},
		{
			name:          "Successful but with at least one cluster not idle",
			mockCluster:   &admin.PaginatedClusterDescription20240805{Results: advancedClusters},
			mockResponse:  &http.Response{StatusCode: http.StatusOK},
			expectedError: false,
			expectedResult: Result{
				response: advancedClusters[0],
				state:    "PENDING",
				error:    nil,
			},
		},
		{
			name:          "Successful",
			mockCluster:   &admin.PaginatedClusterDescription20240805{},
			mockResponse:  &http.Response{StatusCode: http.StatusOK},
			expectedError: false,
			expectedResult: Result{
				response: &admin.PaginatedClusterDescription20240805{},
				state:    "IDLE",
				error:    nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mockadmin.NewClustersAPI(t)

			testObject.EXPECT().ListClusters(mock.Anything, mock.Anything).Return(admin.ListClustersApiRequest{ApiService: testObject}).Once()
			testObject.EXPECT().ListClustersExecute(mock.Anything).Return(tc.mockCluster, tc.mockResponse, tc.mockError).Once()

			result, stateName, err := advancedcluster.ResourceClusterListAdvancedRefreshFunc(t.Context(), dummyProjectID, testObject)()
			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}

			assert.Equal(t, tc.expectedResult.error, err)
			assert.Equal(t, tc.expectedResult.response, result)
			assert.Equal(t, tc.expectedResult.state, stateName)
		})
	}
}

func TestAddIDsToReplicationSpecs(t *testing.T) {
	testCases := map[string]struct {
		ReplicationSpecs          []admin.ReplicationSpec20240805
		ZoneToReplicationSpecsIDs map[string][]string
		ExpectedReplicationSpecs  []admin.ReplicationSpec20240805
	}{
		"two zones with same amount of available ids and replication specs to populate": {
			ReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: new("Zone 1"),
				},
				{
					ZoneName: new("Zone 2"),
				},
				{
					ZoneName: new("Zone 1"),
				},
				{
					ZoneName: new("Zone 2"),
				},
			},
			ZoneToReplicationSpecsIDs: map[string][]string{
				"Zone 1": {"zone1-id1", "zone1-id2"},
				"Zone 2": {"zone2-id1", "zone2-id2"},
			},
			ExpectedReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: new("Zone 1"),
					Id:       new("zone1-id1"),
				},
				{
					ZoneName: new("Zone 2"),
					Id:       new("zone2-id1"),
				},
				{
					ZoneName: new("Zone 1"),
					Id:       new("zone1-id2"),
				},
				{
					ZoneName: new("Zone 2"),
					Id:       new("zone2-id2"),
				},
			},
		},
		"less available ids than replication specs to populate": {
			ReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: new("Zone 1"),
				},
				{
					ZoneName: new("Zone 1"),
				},
				{
					ZoneName: new("Zone 1"),
				},
				{
					ZoneName: new("Zone 2"),
				},
			},
			ZoneToReplicationSpecsIDs: map[string][]string{
				"Zone 1": {"zone1-id1"},
				"Zone 2": {"zone2-id1"},
			},
			ExpectedReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: new("Zone 1"),
					Id:       new("zone1-id1"),
				},
				{
					ZoneName: new("Zone 1"),
					Id:       nil,
				},
				{
					ZoneName: new("Zone 1"),
					Id:       nil,
				},
				{
					ZoneName: new("Zone 2"),
					Id:       new("zone2-id1"),
				},
			},
		},
		"more available ids than replication specs to populate": {
			ReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: new("Zone 1"),
				},
				{
					ZoneName: new("Zone 2"),
				},
			},
			ZoneToReplicationSpecsIDs: map[string][]string{
				"Zone 1": {"zone1-id1", "zone1-id2"},
				"Zone 2": {"zone2-id1", "zone2-id2"},
			},
			ExpectedReplicationSpecs: []admin.ReplicationSpec20240805{
				{
					ZoneName: new("Zone 1"),
					Id:       new("zone1-id1"),
				},
				{
					ZoneName: new("Zone 2"),
					Id:       new("zone2-id1"),
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

func TestSetShardSizeLimitGBNullOnRemoval(t *testing.T) {
	tags := []admin.ResourceTag{{Key: "environment", Value: "test"}}
	testCases := map[string]struct {
		state        *admin.ClusterDescription20240805
		plan         *admin.ClusterDescription20240805
		config       *admin.ClusterDescription20240805
		patch        *admin.ClusterDescription20240805
		expectedJSON string
		expectedNil  bool
	}{
		"removal only": {
			state:        clusterWithShardSizeLimits(new(1024)),
			config:       clusterWithShardSizeLimits(nil),
			expectedJSON: `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}}}]}]}`,
		},
		"preserves unrelated patch": {
			state:        clusterWithShardSizeLimits(new(1024)),
			config:       clusterWithShardSizeLimits(nil),
			patch:        &admin.ClusterDescription20240805{Tags: &tags},
			expectedJSON: `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}}}]}],"tags":[{"key":"environment","value":"test"}]}`,
		},
		"removes limits from multiple regions": {
			state:        clusterWithShardSizeLimits(new(1024), new(1024)),
			config:       clusterWithShardSizeLimits(nil, nil),
			expectedJSON: `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}}},{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}}}]}]}`,
		},
		"removes limit when config omits auto scaling block": {
			state: clusterWithShardSizeLimits(new(1024)),
			config: &admin.ClusterDescription20240805{ReplicationSpecs: &[]admin.ReplicationSpec20240805{
				{RegionConfigs: &[]admin.CloudRegionConfig20240805{{}}},
			}},
			expectedJSON: `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}}}]}]}`,
		},
		"preserves compute auto scaling while clearing storage config": {
			state:        clusterWithShardSizeLimits(new(1024)),
			config:       clusterWithComputeAutoScaling(),
			expectedJSON: `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"compute":{"enabled":true,"maxInstanceSize":"M30","minInstanceSize":"M10","scaleDownEnabled":true},"storageConfig":{"shardSizeLimitGB":null}}}]}]}`,
		},
		"uses config instead of state-derived zero-node specs": {
			state:        clusterWithOptionalSpecs(new(1024), 0),
			plan:         clusterWithOptionalSpecs(nil, 0),
			config:       clusterWithElectableSpec(),
			expectedJSON: `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":1}}]}]}`,
		},
		"replaces normal replication specs patch with config": {
			state:        clusterWithOptionalSpecs(new(1024), 1),
			plan:         clusterWithOptionalSpecs(nil, 0),
			config:       clusterWithElectableSpec(),
			patch:        clusterWithOptionalSpecs(nil, 0),
			expectedJSON: `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":1}}]}]}`,
		},
		"adds null to existing compute auto scaling patch": {
			state:        clusterWithShardSizeLimits(new(1024)),
			config:       clusterWithComputeAutoScaling(),
			patch:        clusterWithComputeAutoScalingAndStorageConfigNull(),
			expectedJSON: `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"compute":{"enabled":true,"maxInstanceSize":"M30","minInstanceSize":"M10","scaleDownEnabled":true},"storageConfig":{"shardSizeLimitGB":null}}}]}]}`,
		},
		"unchanged limit": {
			state:       clusterWithShardSizeLimits(new(1024)),
			config:      clusterWithShardSizeLimits(new(1024)),
			expectedNil: true,
		},
		"config omission preserved by plan": {
			state:       clusterWithShardSizeLimits(new(1024)),
			plan:        clusterWithShardSizeLimits(new(1024)),
			config:      clusterWithShardSizeLimits(nil),
			expectedNil: true,
		},
		"analytics auto scaling is ignored": {
			state:       clusterWithAnalyticsShardSizeLimit(1024),
			config:      clusterWithAnalyticsShardSizeLimit(0),
			expectedNil: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			plan := tc.plan
			if plan == nil {
				plan = tc.config
			}
			configJSONBefore, err := json.Marshal(tc.config)
			require.NoError(t, err)
			result := advancedcluster.SetShardSizeLimitGBNullOnRemoval(tc.state.ReplicationSpecs, plan.ReplicationSpecs, tc.config.ReplicationSpecs, tc.patch)
			if tc.expectedNil {
				assert.Nil(t, result)
				return
			}
			resultJSON, err := json.Marshal(result)
			require.NoError(t, err)
			require.JSONEq(t, tc.expectedJSON, string(resultJSON))
			assert.NotContains(t, string(resultJSON), `"storageConfig":null`)
			assert.Contains(t, string(resultJSON), `"shardSizeLimitGB":null`)
			configJSONAfter, err := json.Marshal(tc.config)
			require.NoError(t, err)
			assert.JSONEq(t, string(configJSONBefore), string(configJSONAfter))
		})
	}
}

func clusterWithShardSizeLimits(limits ...*int) *admin.ClusterDescription20240805 {
	regions := make([]admin.CloudRegionConfig20240805, len(limits))
	for i, limit := range limits {
		autoScaling := &admin.AdvancedAutoScalingSettings{}
		if limit != nil {
			autoScaling.StorageConfig = &admin.StorageConfig{ShardSizeLimitGB: limit}
		}
		regions[i].AutoScaling = autoScaling
	}
	return &admin.ClusterDescription20240805{ReplicationSpecs: &[]admin.ReplicationSpec20240805{{RegionConfigs: &regions}}}
}

func clusterWithAnalyticsShardSizeLimit(limit int) *admin.ClusterDescription20240805 {
	var storageConfig *admin.StorageConfig
	if limit != 0 {
		storageConfig = &admin.StorageConfig{ShardSizeLimitGB: &limit}
	}
	regions := []admin.CloudRegionConfig20240805{{
		AnalyticsAutoScaling: &admin.AdvancedAutoScalingSettings{StorageConfig: storageConfig},
	}}
	return &admin.ClusterDescription20240805{ReplicationSpecs: &[]admin.ReplicationSpec20240805{{RegionConfigs: &regions}}}
}

func clusterWithComputeAutoScaling() *admin.ClusterDescription20240805 {
	regions := []admin.CloudRegionConfig20240805{{
		AutoScaling: &admin.AdvancedAutoScalingSettings{
			Compute: &admin.AdvancedComputeAutoScaling{
				Enabled:          new(true),
				MaxInstanceSize:  new("M30"),
				MinInstanceSize:  new("M10"),
				ScaleDownEnabled: new(true),
			},
		},
	}}
	return &admin.ClusterDescription20240805{ReplicationSpecs: &[]admin.ReplicationSpec20240805{{RegionConfigs: &regions}}}
}

func clusterWithComputeAutoScalingAndStorageConfigNull() *admin.ClusterDescription20240805 {
	cluster := clusterWithComputeAutoScaling()
	cluster.GetReplicationSpecs()[0].GetRegionConfigs()[0].AutoScaling.SetStorageConfigNil()
	return cluster
}

func clusterWithOptionalSpecs(limit *int, optionalNodeCount int) *admin.ClusterDescription20240805 {
	cluster := clusterWithShardSizeLimits(limit)
	region := &cluster.GetReplicationSpecs()[0].GetRegionConfigs()[0]
	region.AnalyticsSpecs = &admin.DedicatedHardwareSpec20240805{
		InstanceSize: new("M10"),
		NodeCount:    new(optionalNodeCount),
	}
	region.ElectableSpecs = &admin.HardwareSpec20240805{
		InstanceSize: new("M10"),
		NodeCount:    new(1),
	}
	region.ReadOnlySpecs = &admin.DedicatedHardwareSpec20240805{
		InstanceSize: new("M10"),
		NodeCount:    new(optionalNodeCount),
	}
	return cluster
}

func clusterWithElectableSpec() *admin.ClusterDescription20240805 {
	cluster := clusterWithOptionalSpecs(nil, 0)
	region := &cluster.GetReplicationSpecs()[0].GetRegionConfigs()[0]
	region.AnalyticsSpecs = nil
	region.AutoScaling = nil
	region.ReadOnlySpecs = nil
	return cluster
}
