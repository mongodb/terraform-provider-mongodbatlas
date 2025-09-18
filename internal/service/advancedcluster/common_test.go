package advancedcluster_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
	"go.mongodb.org/atlas-sdk/v20250312007/mockadmin"
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
			mockCluster:   &admin.PaginatedClusterDescription20240805{Results: &advancedClusters},
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
			resultSpecs := advancedcluster.AddIDsToReplicationSpecs(tc.ReplicationSpecs, tc.ZoneToReplicationSpecsIDs)
			assert.Equal(t, tc.ExpectedReplicationSpecs, resultSpecs)
		})
	}
}
