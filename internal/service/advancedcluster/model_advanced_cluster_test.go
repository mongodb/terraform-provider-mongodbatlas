package advancedcluster_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115014/admin"
	"go.mongodb.org/atlas-sdk/v20231115014/mockadmin"
)

var (
	dummyClusterName = "clusterName"
	dummyProjectID   = "projectId"
	errGeneric       = errors.New("generic")
	advancedClusters = []admin.AdvancedClusterDescription{{StateName: conversion.StringPtr("NOT IDLE")}}
)

func TestFlattenReplicationSpecs(t *testing.T) {
	var (
		regionName         = "EU_WEST_1"
		providerName       = "AWS"
		expectedID         = "id1"
		unexpectedID       = "id2"
		zoneName           = "z1"
		unexpectedZoneName = "z2"
		admin1             = admin.ReplicationSpec{Id: &expectedID, ZoneName: &zoneName, RegionConfigs: &[]admin.CloudRegionConfig{{
			ProviderName: &providerName,
			RegionName:   conversion.StringPtr(regionName),
		}}}
		admin2 = admin.ReplicationSpec{Id: &unexpectedID, ZoneName: &unexpectedZoneName, RegionConfigs: &[]admin.CloudRegionConfig{{
			ProviderName: &providerName,
			RegionName:   conversion.StringPtr(regionName),
		}}}
		testSchema = map[string]*schema.Schema{
			"project_id": {Type: schema.TypeString},
		}
		tf1SameIDSameZone = map[string]any{
			"id":         expectedID,
			"num_shards": 1,
			"region_configs": []any{
				map[string]any{
					"provider_name": "AWS",
					"region_name":   regionName,
					"zone_name":     zoneName,
				},
			},
		}
		tf2NoIDSameZone = map[string]any{
			"id":         nil,
			"num_shards": 1,
			"region_configs": []any{
				map[string]any{
					"provider_name": "AWS",
					"region_name":   regionName,
					"zone_name":     zoneName,
				},
			},
		}
		tf3NoIDDiffZone = map[string]any{
			"id":         nil,
			"num_shards": 1,
			"region_configs": []any{
				map[string]any{
					"provider_name": "AWS",
					"region_name":   regionName,
					"zone_name":     "differentZone",
				},
			},
		}
		tf4diffIDDiffZone = map[string]any{
			"id":         "unique",
			"num_shards": 1,
			"region_configs": []any{
				map[string]any{
					"provider_name": "AWS",
					"region_name":   regionName,
					"zone_name":     "uniqueZone",
				},
			},
		}
	)
	type expectFlags struct {
		expectedLen int
		differentID bool
	}
	testCases := map[string]struct {
		adminSpecs   []admin.ReplicationSpec
		tfInputSpecs []any
		expectFlags  expectFlags
	}{
		"existing id, should match admin": {
			[]admin.ReplicationSpec{admin1},
			[]any{tf1SameIDSameZone},
			expectFlags{},
		},
		"missing id, should be set when zone_name matches": {
			[]admin.ReplicationSpec{admin1},
			[]any{tf2NoIDSameZone},
			expectFlags{},
		},
		"missing id, should be set when there is one admin spec": {
			[]admin.ReplicationSpec{admin1},
			[]any{tf3NoIDDiffZone},
			expectFlags{},
		},
		"existing different id, should change to the admin spec 1": {
			[]admin.ReplicationSpec{admin1},
			[]any{tf4diffIDDiffZone},
			expectFlags{},
		},
		"existing different id, should change to the admin spec 2": {
			[]admin.ReplicationSpec{admin2},
			[]any{tf1SameIDSameZone},
			expectFlags{
				differentID: true,
			},
		},
		"existing id, should match correct api spec and extra api spec added": {
			[]admin.ReplicationSpec{admin2, admin1},
			[]any{tf1SameIDSameZone},
			expectFlags{
				expectedLen: 2,
			},
		},
		"existing different id and existing same id, only api spec kept": {
			[]admin.ReplicationSpec{admin1},
			[]any{tf4diffIDDiffZone, tf1SameIDSameZone},
			expectFlags{
				expectedLen: 1,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			peeringAPI := mockadmin.NetworkPeeringApi{}

			peeringAPI.EXPECT().ListPeeringContainerByCloudProviderWithParams(mock.Anything, mock.Anything).Return(admin.ListPeeringContainerByCloudProviderApiRequest{ApiService: &peeringAPI})
			peeringAPI.EXPECT().ListPeeringContainerByCloudProviderExecute(mock.Anything).Return(&admin.PaginatedCloudProviderContainer{Results: &[]admin.CloudProviderContainer{{Id: conversion.StringPtr("c1"), RegionName: &regionName, ProviderName: &providerName}}}, nil, nil)

			client := &admin.APIClient{
				NetworkPeeringApi: &peeringAPI,
			}
			resourceData := schema.TestResourceDataRaw(t, testSchema, map[string]any{"project_id": "p1"})

			tfOutputSpecs, err := advancedcluster.FlattenAdvancedReplicationSpecs(context.Background(), tc.adminSpecs, tc.tfInputSpecs, resourceData, client)

			asserter := assert.New(t)
			require.NoError(t, err)
			flags := tc.expectFlags
			var expectedLen int
			if flags.expectedLen == 0 {
				expectedLen = 1
			} else {
				expectedLen = flags.expectedLen
			}
			asserter.Len(tfOutputSpecs, expectedLen)
			if flags.differentID {
				asserter.NotEqual(expectedID, tfOutputSpecs[0]["id"])
			} else {
				asserter.Equal(expectedID, tfOutputSpecs[0]["id"])
			}
		})
	}
}

type Result struct {
	response any
	error    error
	state    string
}

func TestUpgradeRefreshFunc(t *testing.T) {
	testCases := []struct {
		mockCluster    *admin.AdvancedClusterDescription
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
			name:          "Successful",
			mockCluster:   &admin.AdvancedClusterDescription{StateName: conversion.StringPtr("stateName")},
			mockResponse:  &http.Response{StatusCode: 200},
			expectedError: false,
			expectedResult: Result{
				response: &admin.AdvancedClusterDescription{StateName: conversion.StringPtr("stateName")},
				state:    "stateName",
				error:    nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mockadmin.NewClustersApi(t)

			testObject.EXPECT().GetCluster(mock.Anything, mock.Anything, mock.Anything).Return(admin.GetClusterApiRequest{ApiService: testObject}).Once()
			testObject.EXPECT().GetClusterExecute(mock.Anything).Return(tc.mockCluster, tc.mockResponse, tc.mockError).Once()

			result, stateName, err := advancedcluster.UpgradeRefreshFunc(context.Background(), dummyClusterName, dummyProjectID, testObject)()
			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}

			assert.Equal(t, tc.expectedResult.error, err)
			assert.Equal(t, tc.expectedResult.response, result)
			assert.Equal(t, tc.expectedResult.state, stateName)
		})
	}
}

func TestResourceListAdvancedRefreshFunc(t *testing.T) {
	testCases := []struct {
		mockCluster    *admin.PaginatedAdvancedClusterDescription
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
			mockCluster:   &admin.PaginatedAdvancedClusterDescription{Results: &advancedClusters},
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
			mockCluster:   &admin.PaginatedAdvancedClusterDescription{},
			mockResponse:  &http.Response{StatusCode: 200},
			expectedError: false,
			expectedResult: Result{
				response: &admin.PaginatedAdvancedClusterDescription{},
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

			result, stateName, err := advancedcluster.ResourceClusterListAdvancedRefreshFunc(context.Background(), dummyProjectID, testObject)()
			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}

			assert.Equal(t, tc.expectedResult.error, err)
			assert.Equal(t, tc.expectedResult.response, result)
			assert.Equal(t, tc.expectedResult.state, stateName)
		})
	}
}
