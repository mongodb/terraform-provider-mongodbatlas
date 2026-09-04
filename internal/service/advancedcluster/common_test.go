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

func TestSetShardSizeLimitGBNull(t *testing.T) {
	tags := []admin.ResourceTag{{Key: "environment", Value: "test"}}
	mixedRegionConfig := clusterWithShardSizeLimits(nil, nil)
	mixedRegions := mixedRegionConfig.GetReplicationSpecs()[0].GetRegionConfigs()
	mixedRegions[0].AutoScaling.Compute = &admin.AdvancedComputeAutoScaling{}
	mixedRegions[0].AutoScaling.DiskGB = &admin.DiskGBAutoScaling{}
	mixedRegions[1].AutoScaling = computeAutoScaling()
	mixedRegions[1].AutoScaling.DiskGB = &admin.DiskGBAutoScaling{}
	zeroNodePatch := clusterWithOptionalSpecs(nil, 0)
	zeroNodePatch.Tags = &tags
	testCases := map[string]struct {
		config                *admin.ClusterDescription20240805
		patch                 *admin.ClusterDescription20240805
		expectedJSON          string
		assertConfigUnchanged bool
	}{
		"omits empty auto scaling and preserves configured compute": {
			config:                mixedRegionConfig,
			expectedJSON:          `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}}},{"autoScaling":{"compute":{"enabled":true,"maxInstanceSize":"M30","minInstanceSize":"M10","scaleDownEnabled":true},"storageConfig":{"shardSizeLimitGB":null}}}]}]}`,
			assertConfigUnchanged: true,
		},
		"uses config-shaped specs and preserves unrelated patch fields": {
			config:       clusterWithElectableSpec(),
			patch:        zeroNodePatch,
			expectedJSON: `{"replicationSpecs":[{"regionConfigs":[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":1}}]}],"tags":[{"key":"environment","value":"test"}]}`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var configJSONBefore []byte
			if tc.assertConfigUnchanged {
				var err error
				configJSONBefore, err = json.Marshal(tc.config)
				require.NoError(t, err)
			}
			result := advancedcluster.SetShardSizeLimitGBNull(tc.config.ReplicationSpecs, tc.patch)
			resultJSON, err := json.Marshal(result)
			require.NoError(t, err)
			require.JSONEq(t, tc.expectedJSON, string(resultJSON))
			if tc.assertConfigUnchanged {
				configJSONAfter, err := json.Marshal(tc.config)
				require.NoError(t, err)
				assert.JSONEq(t, string(configJSONBefore), string(configJSONAfter))
			}
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

func computeAutoScaling() *admin.AdvancedAutoScalingSettings {
	return &admin.AdvancedAutoScalingSettings{Compute: &admin.AdvancedComputeAutoScaling{
		Enabled:          new(true),
		MaxInstanceSize:  new("M30"),
		MinInstanceSize:  new("M10"),
		ScaleDownEnabled: new(true),
	}}
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
