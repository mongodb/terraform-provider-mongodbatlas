package searchdeployment_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchdeployment"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mocksvc"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

var (
	updating = "UPDATING"
	idle     = "IDLE"
	unknown  = ""
	sc400    = conversion.IntPtr(400)
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
			svc := mocksvc.NewDeploymentService(t)
			ctx := context.Background()
			for _, resp := range tc.mockResponses {
				svc.On("GetAtlasSearchDeployment", ctx, dummyProjectID, clusterName).Return(resp.get()...).Once()
			}
			resp, err := searchdeployment.WaitSearchNodeStateTransition(ctx, dummyProjectID, "Cluster0", svc, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, responseWithState(tc.expectedState), resp)
			svc.AssertExpectations(t)
		})
	}
}

func TestSearchDeploymentStateTransitionForDelete(t *testing.T) {
	testCases := []testCase{
		{
			name: "Regular transition to DELETED",
			mockResponses: []response{
				{state: &updating},
				{statusCode: sc400, err: errors.New(searchdeployment.SearchDeploymentDoesNotExistsError)},
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
			svc := mocksvc.NewDeploymentService(t)
			ctx := context.Background()
			for _, resp := range tc.mockResponses {
				svc.On("GetAtlasSearchDeployment", ctx, dummyProjectID, clusterName).Return(resp.get()...).Once()
			}
			err := searchdeployment.WaitSearchNodeDelete(ctx, dummyProjectID, clusterName, svc, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
			svc.AssertExpectations(t)
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

func (r *response) get() []interface{} {
	var httpResp *http.Response
	if r.statusCode != nil {
		httpResp = &http.Response{StatusCode: *r.statusCode}
	}
	return []interface{}{responseWithState(r.state), httpResp, r.err}
}
