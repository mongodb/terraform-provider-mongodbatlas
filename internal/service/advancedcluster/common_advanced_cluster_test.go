package advancedcluster_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
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

func TestRemoveLabel(t *testing.T) {
	toRemove := matlas.Label{Key: "To Remove", Value: "To remove value"}

	expected := []matlas.Label{
		{Key: "Name", Value: "Test"},
		{Key: "Version", Value: "1.0"},
		{Key: "Type", Value: "testing"},
	}

	labels := []matlas.Label{
		{Key: "Name", Value: "Test"},
		{Key: "Version", Value: "1.0"},
		{Key: "To Remove", Value: "To remove value"},
		{Key: "Type", Value: "testing"},
	}

	got := advancedcluster.RemoveLabel(labels, toRemove)

	if diff := deep.Equal(expected, got); diff != nil {
		t.Fatalf("Bad removeLabel return \n got = %#v\nwant = %#v \ndiff = %#v", got, expected, diff)
	}
}

func TestContainsLabelOrKey(t *testing.T) {
	labelsList := []matlas.Label{
		{Key: "Name", Value: "Test"},
		{Key: "Version", Value: "1.0"},
		{Key: "Type", Value: "testing"},
	}

	labels := []matlas.Label{
		{Key: "Name", Value: "Test"},
		{Key: "Version", Value: "Not same value"},
		{Key: "non existing key", Value: "value"},
	}
	expected := []bool{
		true,
		true,
		false,
	}

	for i, l := range labels {
		result := advancedcluster.ContainsLabelOrKey(labelsList, l)
		assert.Equal(t, expected[i], result)
	}
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
			testObject := new(MockClusterService)

			response := ClusterResponse{
				clusterResponse: tc.mockCluster,
				response:        tc.mockResponse,
				error:           tc.mockError,
			}
			testObject.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(response)

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
			testObject := new(MockClusterService)

			response := ClusterResponse{
				advancedClusterResponse: tc.mockCluster,
				response:                tc.mockResponse,
				error:                   tc.mockError,
			}
			testObject.On("List", mock.Anything, mock.Anything, mock.Anything).Return(response)

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

func TestFlattenLabels(t *testing.T) {
	testCases := []struct {
		name           string
		labels         []matlas.Label
		expectedResult []map[string]any
	}{
		{
			name:           "empty",
			labels:         []matlas.Label{},
			expectedResult: []map[string]any{},
		},
		{
			name: "not empty",
			labels: []matlas.Label{
				{
					Key:   "key",
					Value: "value",
				},
			},
			expectedResult: []map[string]any{
				{
					"key":   "key",
					"value": "value",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := advancedcluster.FlattenLabels(tc.labels)

			assert.Equal(t, len(tc.expectedResult), len(result))
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestFlattenTags(t *testing.T) {
	testCases := []struct {
		name           string
		tags           *[]*matlas.Tag
		expectedResult []map[string]any
	}{
		{
			name:           "empty",
			tags:           &[]*matlas.Tag{},
			expectedResult: []map[string]any{},
		},
		{
			name: "not empty",
			tags: &[]*matlas.Tag{
				{
					Key:   "key",
					Value: "value",
				},
			},
			expectedResult: []map[string]any{
				{
					"key":   "key",
					"value": "value",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := advancedcluster.FlattenTags(tc.tags)

			assert.Equal(t, len(tc.expectedResult), len(result))
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestFlattenConnectionStrings(t *testing.T) {
	testCases := []struct {
		name              string
		connectionStrings *matlas.ConnectionStrings
		expectedResult    []map[string]any
	}{
		{
			name:              "empty",
			connectionStrings: &matlas.ConnectionStrings{},
			expectedResult: []map[string]any{
				{
					"standard":         "",
					"standard_srv":     "",
					"private":          "",
					"private_srv":      "",
					"private_endpoint": []map[string]any{},
				},
			},
		},
		{
			name: "not empty",
			connectionStrings: &matlas.ConnectionStrings{
				Standard:    "Standard",
				StandardSrv: "StandardSrv",
				Private:     "Private",
				PrivateSrv:  "PrivateSrv",
			},

			expectedResult: []map[string]any{
				{
					"standard":         "Standard",
					"standard_srv":     "StandardSrv",
					"private":          "Private",
					"private_srv":      "PrivateSrv",
					"private_endpoint": []map[string]any{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := advancedcluster.FlattenConnectionStrings(tc.connectionStrings)

			assert.Equal(t, len(tc.expectedResult), len(result))
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

type MockClusterService struct {
	mock.Mock
}

func (a *MockClusterService) Get(ctx context.Context, groupID, clusterName string) (*matlas.Cluster, *matlas.Response, error) {
	args := a.Called(ctx, groupID)
	var response = args.Get(0).(ClusterResponse)
	return response.clusterResponse, response.response, response.error
}

func (a *MockClusterService) List(ctx context.Context, groupID string, options *matlas.ListOptions) (*matlas.AdvancedClustersResponse, *matlas.Response, error) {
	args := a.Called(ctx, groupID)
	var response = args.Get(0).(ClusterResponse)
	return response.advancedClusterResponse, response.response, response.error
}

type ClusterResponse struct {
	clusterResponse         *matlas.Cluster
	advancedClusterResponse *matlas.AdvancedClustersResponse
	response                *matlas.Response
	error                   error
}
