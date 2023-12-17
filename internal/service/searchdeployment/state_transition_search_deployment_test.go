package searchdeployment_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchdeployment"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mocksvc"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

var (
	updating = "UPDATING"
	idle     = "IDLE"
	unknown  = ""
)

type stateTransitionTestCase struct {
	expectedState *string
	name          string
	mockResponses []SearchDeploymentResponse
	expectedError bool
}

func TestSearchDeploymentStateTransition(t *testing.T) {
	testCases := []stateTransitionTestCase{
		{
			name: "Successful transition to IDLE",
			mockResponses: []SearchDeploymentResponse{
				{
					state: &updating,
				},
				{
					state: &idle,
				},
			},
			expectedState: &idle,
			expectedError: false,
		},
		{
			name: "Successful transition to IDLE with 503 error in between",
			mockResponses: []SearchDeploymentResponse{
				{
					state: &updating,
				},
				{
					state:        nil,
					HTTPResponse: &http.Response{StatusCode: 503},
					Err:          errors.New("Service Unavailable"),
				},
				{
					state: &idle,
				},
			},
			expectedState: &idle,
			expectedError: false,
		},
		{
			name: "Error when transitioning to an unknown state",
			mockResponses: []SearchDeploymentResponse{
				{
					state: &updating,
				},
				{
					state: &unknown,
				},
			},
			expectedState: nil,
			expectedError: true,
		},
		{
			name: "Error when API responds with error",
			mockResponses: []SearchDeploymentResponse{
				{
					state:        nil,
					HTTPResponse: &http.Response{StatusCode: 500},
					Err:          errors.New("Internal server error"),
				},
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
				svc.On("GetAtlasSearchDeployment", ctx, dummyProjectID, clusterName).Return(responseWithState(resp.state), resp.HTTPResponse, resp.Err).Once()
			}
			resp, err := searchdeployment.WaitSearchNodeStateTransition(ctx, dummyProjectID, "Cluster0", svc, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, responseWithState(tc.expectedState), resp)
			svc.AssertExpectations(t)
		})
	}
}

func TestSearchDeploymentStateTransitionForDelete(t *testing.T) {
	testCases := []stateTransitionTestCase{
		{
			name: "Regular transition to DELETED",
			mockResponses: []SearchDeploymentResponse{
				{
					state: &updating,
				},
				{
					state:        nil,
					HTTPResponse: &http.Response{StatusCode: 400},
					Err:          errors.New(searchdeployment.SearchDeploymentDoesNotExistsError),
				},
			},
			expectedError: false,
		},
		{
			name: "Error when API responds with error",
			mockResponses: []SearchDeploymentResponse{
				{
					state:        nil,
					HTTPResponse: &http.Response{StatusCode: 500},
					Err:          errors.New("Internal server error"),
				},
			},
			expectedError: true,
		},
		{
			name: "Failed delete when responding with unknown state",
			mockResponses: []SearchDeploymentResponse{
				{
					state: &updating,
				},
				{
					state: &unknown,
				},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocksvc.NewDeploymentService(t)
			ctx := context.Background()
			for _, resp := range tc.mockResponses {
				svc.On("GetAtlasSearchDeployment", ctx, dummyProjectID, clusterName).Return(responseWithState(resp.state), resp.HTTPResponse, resp.Err).Once()
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
		Specs: []admin.ApiSearchDeploymentSpec{
			{
				InstanceSize: instanceSize,
				NodeCount:    nodeCount,
			},
		},
		StateName: state,
	}
}

type SearchDeploymentResponse struct {
	state        *string
	HTTPResponse *http.Response
	Err          error
}
