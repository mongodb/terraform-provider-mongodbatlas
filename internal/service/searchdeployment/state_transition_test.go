package searchdeployment_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchdeployment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/mockadmin"
)

var (
	updating = "UPDATING"
	idle     = "IDLE"
	unknown  = ""
	sc404    = conversion.IntPtr(404)
	sc500    = conversion.IntPtr(500)
	sc503    = conversion.IntPtr(503)
)

type testCase struct {
	expectedState *string
	name          string
	mockResponses []response
	expectedError bool
}

func TestSearchDeploymentStateTransition(t *testing.T) {
	testCases := []testCase{
		{
			name: "Successful transition to IDLE",
			mockResponses: []response{
				{state: &updating},
				{state: &idle},
			},
			expectedState: &idle,
			expectedError: false,
		},
		{
			name: "Successful transition to IDLE with 503 error in between",
			mockResponses: []response{
				{state: &updating},
				{statusCode: sc503, err: errors.New("Service Unavailable")},
				{state: &idle},
			},
			expectedState: &idle,
			expectedError: false,
		},
		{
			name: "Error when transitioning to an unknown state",
			mockResponses: []response{
				{state: &updating},
				{state: &unknown},
			},
			expectedState: nil,
			expectedError: true,
		},
		{
			name: "Error when API responds with error",
			mockResponses: []response{
				{statusCode: sc500, err: errors.New("Internal server error")},
			},
			expectedState: nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := mockadmin.NewAtlasSearchApi(t)
			m.EXPECT().GetAtlasSearchDeployment(mock.Anything, mock.Anything, mock.Anything).Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: m})

			for _, resp := range tc.mockResponses {
				modelResp, httpResp, err := resp.get()
				m.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).Return(modelResp, httpResp, err).Once()
			}
			resp, err := searchdeployment.WaitSearchNodeStateTransition(t.Context(), dummyProjectID, "Cluster0", m, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, responseWithState(tc.expectedState), resp)
		})
	}
}

func TestSearchDeploymentStateTransitionForDelete(t *testing.T) {
	testCases := []testCase{
		{
			name: "Regular transition to DELETED",
			mockResponses: []response{
				{state: &updating},
				{statusCode: sc404, err: errors.New(searchdeployment.SearchDeploymentDoesNotExistsError)},
			},
			expectedError: false,
		},
		{
			name: "Error when API responds with error",
			mockResponses: []response{
				{statusCode: sc500, err: errors.New("Internal server error")},
			},
			expectedError: true,
		},
		{
			name: "Failed delete when responding with unknown state",
			mockResponses: []response{
				{state: &updating},
				{state: &unknown},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := mockadmin.NewAtlasSearchApi(t)
			m.EXPECT().GetAtlasSearchDeployment(mock.Anything, mock.Anything, mock.Anything).Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: m})

			for _, resp := range tc.mockResponses {
				modelResp, httpResp, err := resp.get()
				m.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).Return(modelResp, httpResp, err).Once()
			}
			err := searchdeployment.WaitSearchNodeDelete(t.Context(), dummyProjectID, clusterName, m, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
		})
	}
}

var testTimeoutConfig = retrystrategy.TimeConfig{
	Timeout:    30 * time.Second,
	MinTimeout: 100 * time.Millisecond,
	Delay:      0,
}

func responseWithState(state *string) *admin.ApiSearchDeploymentResponse {
	if state == nil {
		return nil
	}
	return &admin.ApiSearchDeploymentResponse{
		GroupId: admin.PtrString(dummyProjectID),
		Id:      admin.PtrString(dummyDeploymentID),
		Specs: &[]admin.ApiSearchDeploymentSpec{
			{
				InstanceSize: instanceSize,
				NodeCount:    nodeCount,
			},
		},
		StateName: state,
	}
}

type response struct {
	state      *string
	statusCode *int
	err        error
}

func (r *response) get() (*admin.ApiSearchDeploymentResponse, *http.Response, error) {
	var httpResp *http.Response
	if r.statusCode != nil {
		httpResp = &http.Response{StatusCode: *r.statusCode}
	}
	return responseWithState(r.state), httpResp, r.err
}
