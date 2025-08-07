package advancedcluster_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/mockadmin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

var (
	dummyProjectID   = "projectId"
	errGeneric       = errors.New("generic")
	advancedClusters = []admin.ClusterDescription20240805{{StateName: conversion.StringPtr("NOT IDLE")}}
)

func TestFlattenAdvancedReplicationSpecsOldShardingConfig(t *testing.T) {
	var (
		regionName         = "EU_WEST_1"
		providerName       = "AWS"
		expectedID         = "id1"
		unexpectedID       = "id2"
		expectedZoneName   = "z1"
		unexpectedZoneName = "z2"
		regionConfigAdmin  = []admin.CloudRegionConfig20240805{{
			ProviderName: &providerName,
			RegionName:   &regionName,
		}}
		regionConfigTfSameZone = map[string]any{
			"provider_name": "AWS",
			"region_name":   regionName,
		}
		regionConfigTfDiffZone = map[string]any{
			"provider_name": "AWS",
			"region_name":   regionName,
			"zone_name":     unexpectedZoneName,
		}
		apiSpecExpected  = admin.ReplicationSpec20240805{Id: &expectedID, ZoneName: &expectedZoneName, RegionConfigs: &regionConfigAdmin}
		apiSpecDifferent = admin.ReplicationSpec20240805{Id: &unexpectedID, ZoneName: &unexpectedZoneName, RegionConfigs: &regionConfigAdmin}
		testSchema       = map[string]*schema.Schema{
			"project_id": {Type: schema.TypeString},
		}
		tfSameIDSameZone = map[string]any{
			"id":             expectedID,
			"num_shards":     1,
			"region_configs": []any{regionConfigTfSameZone},
			"zone_name":      expectedZoneName,
		}
		tfNoIDSameZone = map[string]any{
			"id":             nil,
			"num_shards":     1,
			"region_configs": []any{regionConfigTfSameZone},
			"zone_name":      expectedZoneName,
		}
		tfNoIDDiffZone = map[string]any{
			"id":             nil,
			"num_shards":     1,
			"region_configs": []any{regionConfigTfDiffZone},
			"zone_name":      unexpectedZoneName,
		}
		tfdiffIDDiffZone = map[string]any{
			"id":             "unique",
			"num_shards":     1,
			"region_configs": []any{regionConfigTfDiffZone},
			"zone_name":      unexpectedZoneName,
		}
	)
	testCases := map[string]struct {
		adminSpecs                       []admin.ReplicationSpec20240805
		zoneNameToOldReplicationSpecMeta map[string]advancedcluster.OldShardConfigMeta
		tfInputSpecs                     []any
		expectedLen                      int
	}{
		"empty admin spec should return empty list": {
			[]admin.ReplicationSpec20240805{},
			map[string]advancedcluster.OldShardConfigMeta{},
			[]any{tfSameIDSameZone},
			0,
		},
		"existing id, should match admin": {
			[]admin.ReplicationSpec20240805{apiSpecExpected},
			map[string]advancedcluster.OldShardConfigMeta{expectedZoneName: {expectedID, 1}},
			[]any{tfSameIDSameZone},
			1,
		},
		"existing different id, should change to admin spec": {
			[]admin.ReplicationSpec20240805{apiSpecExpected},
			map[string]advancedcluster.OldShardConfigMeta{expectedZoneName: {expectedID, 1}},
			[]any{tfdiffIDDiffZone},
			1,
		},
		"missing id, should be set when zone_name matches": {
			[]admin.ReplicationSpec20240805{apiSpecExpected},
			map[string]advancedcluster.OldShardConfigMeta{expectedZoneName: {expectedID, 1}},
			[]any{tfNoIDSameZone},
			1,
		},
		"missing id and diff zone, should change to admin spec": {
			[]admin.ReplicationSpec20240805{apiSpecExpected},
			map[string]advancedcluster.OldShardConfigMeta{expectedZoneName: {expectedID, 1}},
			[]any{tfNoIDDiffZone},
			1,
		},
		"existing id, should match correct api spec using `id` and extra api spec added": {
			[]admin.ReplicationSpec20240805{apiSpecDifferent, apiSpecExpected},
			map[string]advancedcluster.OldShardConfigMeta{unexpectedZoneName: {unexpectedID, 1}, expectedZoneName: {expectedID, 1}},
			[]any{tfSameIDSameZone},
			2,
		},
		"missing id, should match correct api spec using `zone_name` and extra api spec added": {
			[]admin.ReplicationSpec20240805{apiSpecDifferent, apiSpecExpected},
			map[string]advancedcluster.OldShardConfigMeta{unexpectedZoneName: {unexpectedID, 1}, expectedZoneName: {expectedID, 1}},
			[]any{tfNoIDSameZone},
			2,
		},
		"two matching specs should be set to api specs": {
			[]admin.ReplicationSpec20240805{apiSpecExpected, apiSpecDifferent},
			map[string]advancedcluster.OldShardConfigMeta{expectedZoneName: {expectedID, 1}, unexpectedZoneName: {unexpectedID, 1}},
			[]any{tfSameIDSameZone, tfdiffIDDiffZone},
			2,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			peeringAPI := mockadmin.NetworkPeeringApi{}

			peeringAPI.EXPECT().ListPeeringContainerByCloudProviderWithParams(mock.Anything, mock.Anything).Return(admin.ListPeeringContainerByCloudProviderApiRequest{ApiService: &peeringAPI})
			containerResult := []admin.CloudProviderContainer{{Id: conversion.StringPtr("c1"), RegionName: &regionName, ProviderName: &providerName}}
			peeringAPI.EXPECT().ListPeeringContainerByCloudProviderExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{Results: &containerResult}, nil, nil)

			client := &admin.APIClient{
				NetworkPeeringApi: &peeringAPI,
			}
			resourceData := schema.TestResourceDataRaw(t, testSchema, map[string]any{"project_id": "p1"})

			tfOutputSpecs, err := advancedcluster.FlattenAdvancedReplicationSpecsOldShardingConfig(t.Context(), tc.adminSpecs, tc.zoneNameToOldReplicationSpecMeta, tc.tfInputSpecs, resourceData, client)

			require.NoError(t, err)
			assert.Len(t, tfOutputSpecs, tc.expectedLen)
			if tc.expectedLen != 0 {
				assert.Equal(t, expectedID, tfOutputSpecs[0]["id"])
				assert.Equal(t, expectedZoneName, tfOutputSpecs[0]["zone_name"])
			}
		})
	}
}

func TestGetDiskSizeGBFromReplicationSpec(t *testing.T) {
	diskSizeGBValue := 40.0

	testCases := map[string]struct {
		clusterDescription     admin.ClusterDescription20240805
		expectedDiskSizeResult float64
	}{
		"cluster description with disk size gb value at electable spec": {
			clusterDescription: admin.ClusterDescription20240805{
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{{
					RegionConfigs: &[]admin.CloudRegionConfig20240805{{
						ElectableSpecs: &admin.HardwareSpec20240805{
							DiskSizeGB: admin.PtrFloat64(diskSizeGBValue),
						},
					}},
				}},
			},
			expectedDiskSizeResult: diskSizeGBValue,
		},
		"cluster description with no electable spec": {
			clusterDescription: admin.ClusterDescription20240805{
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{RegionConfigs: &[]admin.CloudRegionConfig20240805{{}}},
				},
			},
			expectedDiskSizeResult: 0,
		},
		"cluster description with no replication spec": {
			clusterDescription:     admin.ClusterDescription20240805{},
			expectedDiskSizeResult: 0,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := advancedcluster.GetDiskSizeGBFromReplicationSpec(&tc.clusterDescription)
			assert.Equal(t, fmt.Sprintf("%.f", tc.expectedDiskSizeResult), fmt.Sprintf("%.f", result)) // formatting to string to avoid float comparison
		})
	}
}

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
			mockResponse:  &http.Response{StatusCode: 404},
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
			mockResponse:  &http.Response{StatusCode: 503},
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
			mockResponse:  &http.Response{StatusCode: 400},
			expectedError: true,
			expectedResult: Result{
				response: nil,
				state:    "",
				error:    errGeneric,
			},
		},
		{
			name:          "Successful but with at least one cluster not idle",
			mockCluster:   &admin.PaginatedClusterDescription20240805{Results: &advancedClusters},
			mockResponse:  &http.Response{StatusCode: 200},
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
			mockResponse:  &http.Response{StatusCode: 200},
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
			testObject := mockadmin.NewClustersApi(t)

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

func TestIsChangeStreamOptionsMinRequiredMajorVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty input", "", true},
		{"Valid input equal to 6", "6", true},
		{"Valid input greater than 6", "7.0", true},
		{"Valid input less than 6", "5", false},
		{"Valid float input greater", "6.5", true},
		{"Valid float input less", "5.9", false},
		{"Valid float complete semantic version", "6.0.2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := advancedcluster.IsChangeStreamOptionsMinRequiredMajorVersion(&tt.input); got != tt.want {
				t.Errorf("abc(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCheckRegionConfigsPriorityOrder(t *testing.T) {
	testCases := map[string]struct {
		priorities    []int
		errorExpected bool
	}{
		"Priority order 3 entries": {
			priorities: []int{7, 6, 5},
		},
		"Priority order 2 entries": {
			priorities: []int{7, 6},
		},
		"Only 1 entry": {
			priorities: []int{7},
		},
		"Same order 3 entries": {
			priorities: []int{7, 0, 0},
		},
		"Same order 2 entries": {
			priorities: []int{0, 0},
		},
		"Invalid priority order 2 entries": {
			priorities:    []int{6, 7},
			errorExpected: true,
		},
		"Invalid priority order 3 entries": {
			priorities:    []int{7, 5, 6},
			errorExpected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			configs := make([]admin.CloudRegionConfig20240805, len(tc.priorities))
			configsOld := make([]admin20240530.CloudRegionConfig, len(tc.priorities))
			for i, priority := range tc.priorities {
				configs[i].Priority = conversion.IntPtr(priority)
				configsOld[i].Priority = conversion.IntPtr(priority)
			}
			err := advancedcluster.CheckRegionConfigsPriorityOrder([]admin.ReplicationSpec20240805{{RegionConfigs: &configs}})
			assert.Equal(t, tc.errorExpected, err != nil)
			err = advancedcluster.CheckRegionConfigsPriorityOrderOld([]admin20240530.ReplicationSpec{{RegionConfigs: &configsOld}})
			assert.Equal(t, tc.errorExpected, err != nil)
		})
	}
}
