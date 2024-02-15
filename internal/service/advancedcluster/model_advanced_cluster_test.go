package advancedcluster_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mocksvc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

var (
	dummyClusterName = "clusterName"
	dummyProjectID   = "projectId"
	genericError     = matlas.NewArgError("error", "generic")
	advancedClusters = []*matlas.AdvancedCluster{{StateName: "NOT IDLE"}}
)

type Result struct {
	response any
	error    error
	state    string
}

func TestResourceClusterRefreshFunc(t *testing.T) {
	testCases := []struct {
		mockCluster    *matlas.Cluster
		mockResponse   *matlas.Response
		expectedResult Result
		mockError      error
		name           string
		expectedError  bool
	}{
		{
			name:          "Error in the API call: reset by peer",
			mockError:     matlas.NewArgError("error", "reset by peer"),
			expectedError: false,
			expectedResult: Result{
				response: nil,
				state:    "REPEATING",
				error:    nil,
			},
		},
		{
			name:          "Generic error in the API call",
			mockError:     genericError,
			expectedError: true,
			expectedResult: Result{
				response: nil,
				state:    "",
				error:    genericError,
			},
		},
		{
			name:          "Error in the API call: HTTP 404",
			mockError:     genericError,
			mockResponse:  &matlas.Response{Response: &http.Response{StatusCode: 404}, Links: nil, Raw: nil},
			expectedError: false,
			expectedResult: Result{
				response: "",
				state:    "DELETED",
				error:    nil,
			},
		},
		{
			name:          "Error in the API call: HTTP 503",
			mockError:     genericError,
			mockResponse:  &matlas.Response{Response: &http.Response{StatusCode: 503}, Links: nil, Raw: nil},
			expectedError: false,
			expectedResult: Result{
				response: "",
				state:    "PENDING",
				error:    nil,
			},
		},
		{
			name:          "Error in the API call: Neither HTTP 503 or 404",
			mockError:     genericError,
			mockResponse:  &matlas.Response{Response: &http.Response{StatusCode: 400}, Links: nil, Raw: nil},
			expectedError: true,
			expectedResult: Result{
				response: nil,
				state:    "",
				error:    genericError,
			},
		},
		{
			name:          "Successful",
			mockCluster:   &matlas.Cluster{StateName: "stateName"},
			mockResponse:  &matlas.Response{Response: &http.Response{StatusCode: 200}, Links: nil, Raw: nil},
			expectedError: false,
			expectedResult: Result{
				response: &matlas.Cluster{StateName: "stateName"},
				state:    "stateName",
				error:    nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mocksvc.NewClusterService(t)

			testObject.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockCluster, tc.mockResponse, tc.mockError)

			result, stateName, err := advancedcluster.ResourceClusterRefreshFunc(context.Background(), dummyClusterName, dummyProjectID, testObject)()
			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}

			assert.Equal(t, tc.expectedResult.error, err)
			assert.Equal(t, tc.expectedResult.response, result)
			assert.Equal(t, tc.expectedResult.state, stateName)
		})
	}
}

func TestStringIsUppercase(t *testing.T) {
	testCases := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "AWS",
			expectedError: false,
		},
		{
			name:          "aws",
			expectedError: true,
		},
		{
			name:          "",
			expectedError: false,
		},
		{
			name:          "AwS",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			diag := advancedcluster.StringIsUppercase()(tc.name, nil)
			if diag.HasError() != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, diag[0].Summary)
			}
		})
	}
}

func TestResourceListAdvancedRefreshFunc(t *testing.T) {
	testCases := []struct {
		mockCluster    *matlas.AdvancedClustersResponse
		mockResponse   *matlas.Response
		expectedResult Result
		mockError      error
		name           string
		expectedError  bool
	}{
		{
			name:          "Error in the API call: reset by peer",
			mockError:     matlas.NewArgError("error", "reset by peer"),
			expectedError: false,
			expectedResult: Result{
				response: nil,
				state:    "REPEATING",
				error:    nil,
			},
		},
		{
			name:          "Generic error in the API call",
			mockError:     genericError,
			expectedError: true,
			expectedResult: Result{
				response: nil,
				state:    "",
				error:    genericError,
			},
		},
		{
			name:          "Error in the API call: HTTP 404",
			mockError:     genericError,
			mockResponse:  &matlas.Response{Response: &http.Response{StatusCode: 404}, Links: nil, Raw: nil},
			expectedError: false,
			expectedResult: Result{
				response: "",
				state:    "DELETED",
				error:    nil,
			},
		},
		{
			name:          "Error in the API call: HTTP 503",
			mockError:     genericError,
			mockResponse:  &matlas.Response{Response: &http.Response{StatusCode: 503}, Links: nil, Raw: nil},
			expectedError: false,
			expectedResult: Result{
				response: "",
				state:    "PENDING",
				error:    nil,
			},
		},
		{
			name:          "Error in the API call: Neither HTTP 503 or 404",
			mockError:     genericError,
			mockResponse:  &matlas.Response{Response: &http.Response{StatusCode: 400}, Links: nil, Raw: nil},
			expectedError: true,
			expectedResult: Result{
				response: nil,
				state:    "",
				error:    genericError,
			},
		},
		{
			name:          "Successful but with at least one cluster not idle",
			mockCluster:   &matlas.AdvancedClustersResponse{Results: advancedClusters},
			mockResponse:  &matlas.Response{Response: &http.Response{StatusCode: 200}, Links: nil, Raw: nil},
			expectedError: false,
			expectedResult: Result{
				response: advancedClusters[0],
				state:    "PENDING",
				error:    nil,
			},
		},
		{
			name:          "Successful",
			mockCluster:   &matlas.AdvancedClustersResponse{},
			mockResponse:  &matlas.Response{Response: &http.Response{StatusCode: 200}, Links: nil, Raw: nil},
			expectedError: false,
			expectedResult: Result{
				response: &matlas.AdvancedClustersResponse{},
				state:    "IDLE",
				error:    nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mocksvc.NewClusterService(t)

			testObject.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockCluster, tc.mockResponse, tc.mockError)

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
