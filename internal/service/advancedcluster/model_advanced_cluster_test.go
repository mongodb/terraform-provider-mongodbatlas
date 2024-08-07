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
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	mockadmin20240530 "go.mongodb.org/atlas-sdk/v20240530005/mockadmin"
)

var (
	dummyClusterName = "clusterName"
	dummyProjectID   = "projectId"
	errGeneric       = errors.New("generic")
	advancedClusters = []admin20240530.AdvancedClusterDescription{{StateName: conversion.StringPtr("NOT IDLE")}}
)

func TestFlattenReplicationSpecs(t *testing.T) {
	var (
		regionName         = "EU_WEST_1"
		providerName       = "AWS"
		expectedID         = "id1"
		unexpectedID       = "id2"
		expectedZoneName   = "z1"
		unexpectedZoneName = "z2"
		regionConfigAdmin  = []admin20240530.CloudRegionConfig{{
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
		apiSpecExpected  = admin20240530.ReplicationSpec{Id: &expectedID, ZoneName: &expectedZoneName, RegionConfigs: &regionConfigAdmin}
		apiSpecDifferent = admin20240530.ReplicationSpec{Id: &unexpectedID, ZoneName: &unexpectedZoneName, RegionConfigs: &regionConfigAdmin}
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
		adminSpecs   []admin20240530.ReplicationSpec
		tfInputSpecs []any
		expectedLen  int
	}{
		"empty admin spec should return empty list": {
			[]admin20240530.ReplicationSpec{},
			[]any{tfSameIDSameZone},
			0,
		},
		"existing id, should match admin": {
			[]admin20240530.ReplicationSpec{apiSpecExpected},
			[]any{tfSameIDSameZone},
			1,
		},
		"existing different id, should change to admin spec": {
			[]admin20240530.ReplicationSpec{apiSpecExpected},
			[]any{tfdiffIDDiffZone},
			1,
		},
		"missing id, should be set when zone_name matches": {
			[]admin20240530.ReplicationSpec{apiSpecExpected},
			[]any{tfNoIDSameZone},
			1,
		},
		"missing id and diff zone, should change to admin spec": {
			[]admin20240530.ReplicationSpec{apiSpecExpected},
			[]any{tfNoIDDiffZone},
			1,
		},
		"existing id, should match correct api spec using `id` and extra api spec added": {
			[]admin20240530.ReplicationSpec{apiSpecDifferent, apiSpecExpected},
			[]any{tfSameIDSameZone},
			2,
		},
		"missing id, should match correct api spec using `zone_name` and extra api spec added": {
			[]admin20240530.ReplicationSpec{apiSpecDifferent, apiSpecExpected},
			[]any{tfNoIDSameZone},
			2,
		},
		"two matching specs should be set to api specs": {
			[]admin20240530.ReplicationSpec{apiSpecExpected, apiSpecDifferent},
			[]any{tfSameIDSameZone, tfdiffIDDiffZone},
			2,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			peeringAPI := mockadmin20240530.NetworkPeeringApi{}

			peeringAPI.EXPECT().ListPeeringContainerByCloudProviderWithParams(mock.Anything, mock.Anything).Return(admin20240530.ListPeeringContainerByCloudProviderApiRequest{ApiService: &peeringAPI})
			containerResult := []admin20240530.CloudProviderContainer{{Id: conversion.StringPtr("c1"), RegionName: &regionName, ProviderName: &providerName}}
			peeringAPI.EXPECT().ListPeeringContainerByCloudProviderExecute(mock.Anything).Return(&admin20240530.PaginatedCloudProviderContainer{Results: &containerResult}, nil, nil)

			client := &admin20240530.APIClient{
				NetworkPeeringApi: &peeringAPI,
			}
			resourceData := schema.TestResourceDataRaw(t, testSchema, map[string]any{"project_id": "p1"})

			tfOutputSpecs, err := advancedcluster.FlattenAdvancedReplicationSpecs(context.Background(), tc.adminSpecs, tc.tfInputSpecs, resourceData, client)

			require.NoError(t, err)
			assert.Len(t, tfOutputSpecs, tc.expectedLen)
			if tc.expectedLen != 0 {
				assert.Equal(t, expectedID, tfOutputSpecs[0]["id"])
				assert.Equal(t, expectedZoneName, tfOutputSpecs[0]["zone_name"])
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
		mockCluster    *admin20240530.AdvancedClusterDescription
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
			mockCluster:   &admin20240530.AdvancedClusterDescription{StateName: conversion.StringPtr("stateName")},
			mockResponse:  &http.Response{StatusCode: 200},
			expectedError: false,
			expectedResult: Result{
				response: &admin20240530.AdvancedClusterDescription{StateName: conversion.StringPtr("stateName")},
				state:    "stateName",
				error:    nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mockadmin20240530.NewClustersApi(t)

			testObject.EXPECT().GetCluster(mock.Anything, mock.Anything, mock.Anything).Return(admin20240530.GetClusterApiRequest{ApiService: testObject}).Once()
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
		mockCluster    *admin20240530.PaginatedAdvancedClusterDescription
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
			mockCluster:   &admin20240530.PaginatedAdvancedClusterDescription{Results: &advancedClusters},
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
			mockCluster:   &admin20240530.PaginatedAdvancedClusterDescription{},
			mockResponse:  &http.Response{StatusCode: 200},
			expectedError: false,
			expectedResult: Result{
				response: &admin20240530.PaginatedAdvancedClusterDescription{},
				state:    "IDLE",
				error:    nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mockadmin20240530.NewClustersApi(t)

			testObject.EXPECT().ListClusters(mock.Anything, mock.Anything).Return(admin20240530.ListClustersApiRequest{ApiService: testObject}).Once()
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
