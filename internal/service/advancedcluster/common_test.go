package advancedcluster_test

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
	"go.mongodb.org/atlas-sdk/v20250312028/mockadmin"
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

func TestClusterStateAfterCreateTimeoutCleanup(t *testing.T) {
	testCases := []struct {
		mockCluster  *admin.ClusterDescription20240805
		mockResponse *http.Response
		mockError    error
		name         string
		expected     string
	}{
		{
			name:         "HTTP 404 reports DELETED",
			mockResponse: &http.Response{StatusCode: http.StatusNotFound},
			mockError:    errGeneric,
			expected:     "DELETED (404)",
		},
		{
			name:      "CLUSTER_NOT_FOUND error code reports DELETED",
			mockError: clusterNotFoundError(),
			expected:  "DELETED (404)",
		},
		{
			name:         "DELETING cluster reports DELETING",
			mockCluster:  &admin.ClusterDescription20240805{StateName: conversion.StringPtr("DELETING")},
			mockResponse: &http.Response{StatusCode: http.StatusOK},
			expected:     "DELETING",
		},
		{
			name:         "IDLE cluster reports IDLE",
			mockCluster:  &admin.ClusterDescription20240805{StateName: conversion.StringPtr("IDLE")},
			mockResponse: &http.Response{StatusCode: http.StatusOK},
			expected:     "IDLE",
		},
		{
			name:         "Empty stateName reports UNKNOWN",
			mockCluster:  &admin.ClusterDescription20240805{},
			mockResponse: &http.Response{StatusCode: http.StatusOK},
			expected:     "UNKNOWN",
		},
		{
			name:         "Non-404 API error reports error text",
			mockResponse: &http.Response{StatusCode: http.StatusInternalServerError},
			mockError:    errGeneric,
			expected:     "error reading cluster: generic",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := mockadmin.NewClustersAPI(t)
			api.EXPECT().GetCluster(mock.Anything, dummyProjectID, "example").Return(admin.GetClusterApiRequest{ApiService: api}).Once()
			api.EXPECT().GetClusterExecute(mock.Anything).Return(tc.mockCluster, tc.mockResponse, tc.mockError).Once()

			got := clusterStateAfterCreateTimeoutCleanup(t.Context(), api, dummyProjectID, "example")
			assert.Equal(t, tc.expected, got)
		})
	}
}

const createTimeoutTestHCL = `
provider "mongodbatlas" {
  public_key  = "test-public-key"
  private_key = "test-private-key"
  base_url    = "https://atlas.invalid/"
}
resource "mongodbatlas_advanced_cluster" "test" {
  name                     = "example"
  project_id               = "111111111111111111111111"
  cluster_type             = "REPLICASET"
  delete_on_create_timeout = true
  timeouts = {
    create = "400ms"
  }
  replication_specs = [{
    region_configs = [{
      provider_name   = "AWS"
      region_name     = "US_EAST_1"
      priority        = 7
      electable_specs = { instance_size = "M10", node_count = 3, disk_size_gb = 10 }
    }]
  }]
}
`

func shortenClusterRetryTimings(t *testing.T) {
	t.Helper()
	origMin, origDelay, origPoll := advancedcluster.RetryMinTimeout, advancedcluster.RetryDelay, advancedcluster.RetryPollInterval
	advancedcluster.RetryMinTimeout = 100 * time.Millisecond
	advancedcluster.RetryDelay = 100 * time.Millisecond
	advancedcluster.RetryPollInterval = 100 * time.Millisecond
	t.Cleanup(func() {
		advancedcluster.RetryMinTimeout = origMin
		advancedcluster.RetryDelay = origDelay
		advancedcluster.RetryPollInterval = origPoll
	})
}

func TestCreateTimeoutWithDeleteOnCreate_waitVsPostTimeout(t *testing.T) {
	waitTimeoutRegex := regexp.MustCompile(`(?s)didn't reach desired state:\s+IDLE`)

	t.Run("wait timeout after POST commits deletes the created cluster", func(t *testing.T) {
		shortenClusterRetryTimings(t)
		storageConfigCredentials(t)
		t.Setenv("MONGODB_ATLAS_ACCESS_TOKEN", "")
		httpMock := &createTimeoutHTTPMock{}
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, httpMock),
			Steps: []resource.TestStep{
				{
					Config:      createTimeoutTestHCL,
					ExpectError: waitTimeoutRegex,
				},
			},
		})
		httpMock.mu.Lock()
		defer httpMock.mu.Unlock()
		assert.Equal(t, 1, httpMock.postCalls)
		assert.Equal(t, 1, httpMock.deleteCalls)
		assert.Equal(t, "DELETING", httpMock.state)
		require.NotEmpty(t, httpMock.events)
		assert.Contains(t, httpMock.events, "POST-201")
		assert.Equal(t, "DELETE-202", httpMock.events[len(httpMock.events)-1])
	})

	t.Run("client POST timeout races cleanup and does not match the IDLE wait error", func(t *testing.T) {
		shortenClusterRetryTimings(t)
		storageConfigCredentials(t)
		t.Setenv("MONGODB_ATLAS_ACCESS_TOKEN", "")
		httpMock := &createTimeoutHTTPMock{stallPostUntilCancel: true}
		postTimeoutRegex := regexp.MustCompile(`(?s)API error details:\s+Post\s+".*/clusters":\s+context deadline exceeded`)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, httpMock),
			Steps: []resource.TestStep{
				{
					Config:      createTimeoutTestHCL,
					ExpectError: postTimeoutRegex,
				},
			},
		})
		httpMock.mu.Lock()
		defer httpMock.mu.Unlock()
		// Current cleanup.OnTimeout behavior wraps POST /clusters as well as AwaitChanges, so when DELETE
		// runs before a stalled server-side POST commits, DELETE sees 404 before any cluster exists.
		// Step 1's IDLE-wait regex rejects this POST-timeout error; update this expectation if
		// cleanup.OnTimeout is later narrowed to wrap only the post-create wait.
		assert.Equal(t, []string{"POST-cancelled", "DELETE-404"}, httpMock.events)
	})
}

type createTimeoutHTTPMock struct {
	state                string
	events               []string
	postCalls            int
	deleteCalls          int
	mu                   sync.Mutex
	stallPostUntilCancel bool
}

func (m *createTimeoutHTTPMock) ModifyHTTPClient(client *http.Client) error {
	client.Transport = m
	return nil
}

func (m *createTimeoutHTTPMock) ResetHTTPClient(*http.Client) {}

func (m *createTimeoutHTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	const groupPath = "/api/atlas/v2/groups/111111111111111111111111"
	const clusterPath = groupPath + "/clusters/example"

	if req.Method == http.MethodPost && req.URL.Path == groupPath+"/clusters" {
		if m.stallPostUntilCancel {
			<-req.Context().Done()
			err := req.Context().Err()
			m.mu.Lock()
			m.postCalls++
			m.events = append(m.events, "POST-cancelled")
			m.mu.Unlock()
			return nil, err
		}
		// Simulate non-zero network latency so outCtx expires ahead of StateChangeConf's internal timer, matching live HTTP calls.
		select {
		case <-time.After(30 * time.Millisecond):
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	status := http.StatusOK
	body := "{}"
	switch {
	case req.Method == http.MethodPost && req.URL.Path == groupPath+"/clusters":
		m.postCalls++
		m.state = "CREATING"
		m.events = append(m.events, "POST-201")
		status = http.StatusCreated
		body = `{"name":"example","stateName":"CREATING"}`
	case req.Method == http.MethodGet && req.URL.Path == clusterPath:
		m.events = append(m.events, "GET")
		if m.state == "" {
			status, body = http.StatusNotFound, `{"errorCode":"CLUSTER_NOT_FOUND","error":404}`
		} else {
			body = `{"name":"example","stateName":"` + m.state + `"}`
		}
	case req.Method == http.MethodDelete && req.URL.Path == clusterPath:
		m.deleteCalls++
		if m.state == "" {
			m.events = append(m.events, "DELETE-404")
			status, body = http.StatusNotFound, `{"errorCode":"CLUSTER_NOT_FOUND","error":404}`
		} else {
			m.state = "DELETING"
			m.events = append(m.events, "DELETE-202")
			status = http.StatusAccepted
		}
	default:
		return nil, errors.New("unexpected mocked Atlas request: " + req.Method + " " + req.URL.Path)
	}
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}, nil
}
