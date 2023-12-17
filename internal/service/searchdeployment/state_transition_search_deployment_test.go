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

type stateTransitionTestCase struct {
	expectedResult *admin.ApiSearchDeploymentResponse
	name           string
	mockResponses  []SearchDeploymentResponse
	expectedError  bool
}

func TestSearchDeploymentStateTransition(t *testing.T) {
	testCases := []stateTransitionTestCase{
		{
			name: "Successful transition to IDLE",
			mockResponses: []SearchDeploymentResponse{
				{
					DeploymentResp: responseWithState("UPDATING"),
				},
				{
					DeploymentResp: responseWithState("IDLE"),
				},
			},
			expectedResult: responseWithState("IDLE"),
			expectedError:  false,
		},
		{
			name: "Successful transition to IDLE with 503 error in between",
			mockResponses: []SearchDeploymentResponse{
				{
					DeploymentResp: responseWithState("UPDATING"),
				},
				{
					DeploymentResp: nil,
					HTTPResponse:   &http.Response{StatusCode: 503},
					Err:            errors.New("Service Unavailable"),
				},
				{
					DeploymentResp: responseWithState("IDLE"),
				},
			},
			expectedResult: responseWithState("IDLE"),
			expectedError:  false,
		},
		{
			name: "Error when transitioning to an unknown state",
			mockResponses: []SearchDeploymentResponse{
				{
					DeploymentResp: responseWithState("UPDATING"),
				},
				{
					DeploymentResp: responseWithState(""),
				},
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Error when API responds with error",
			mockResponses: []SearchDeploymentResponse{
				{
					DeploymentResp: nil,
					HTTPResponse:   &http.Response{StatusCode: 500},
					Err:            errors.New("Internal server error"),
				},
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := mocksvc.NewDeploymentService(t)
			ctx := context.Background()
			for _, resp := range tc.mockResponses {
				sut.On("GetAtlasSearchDeployment", ctx, dummyProjectID, clusterName).Return(resp.DeploymentResp, resp.HTTPResponse, resp.Err).Once()
			}
			resp, err := searchdeployment.WaitSearchNodeStateTransition(ctx, dummyProjectID, "Cluster0", sut, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, tc.expectedResult, resp)
			sut.AssertExpectations(t)
		})
	}
}

func TestSearchDeploymentStateTransitionForDelete(t *testing.T) {
	testCases := []stateTransitionTestCase{
		{
			name: "Regular transition to DELETED",
			mockResponses: []SearchDeploymentResponse{
				{
					DeploymentResp: responseWithState("UPDATING"),
				},
				{
					DeploymentResp: nil,
					HTTPResponse:   &http.Response{StatusCode: 400},
					Err:            errors.New(searchdeployment.SearchDeploymentDoesNotExistsError),
				},
			},
			expectedError: false,
		},
		{
			name: "Error when API responds with error",
			mockResponses: []SearchDeploymentResponse{
				{
					DeploymentResp: nil,
					HTTPResponse:   &http.Response{StatusCode: 500},
					Err:            errors.New("Internal server error"),
				},
			},
			expectedError: true,
		},
		{
			name: "Failed delete when responding with unknown state",
			mockResponses: []SearchDeploymentResponse{
				{
					DeploymentResp: responseWithState("UPDATING"),
				},
				{
					DeploymentResp: responseWithState(""),
				},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := mocksvc.NewDeploymentService(t)
			ctx := context.Background()
			for _, resp := range tc.mockResponses {
				sut.On("GetAtlasSearchDeployment", ctx, dummyProjectID, clusterName).Return(resp.DeploymentResp, resp.HTTPResponse, resp.Err).Once()
			}
			err := searchdeployment.WaitSearchNodeDelete(ctx, dummyProjectID, clusterName, sut, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
			sut.AssertExpectations(t)
		})
	}
}

var testTimeoutConfig = retrystrategy.TimeConfig{
	Timeout:    30 * time.Second,
	MinTimeout: 100 * time.Millisecond,
	Delay:      0,
}

func responseWithState(state string) *admin.ApiSearchDeploymentResponse {
	return &admin.ApiSearchDeploymentResponse{
		GroupId: admin.PtrString(dummyProjectID),
		Id:      admin.PtrString(dummyDeploymentID),
		Specs: []admin.ApiSearchDeploymentSpec{
			{
				InstanceSize: instanceSize,
				NodeCount:    nodeCount,
			},
		},
		StateName: admin.PtrString(state),
	}
}

type SearchDeploymentResponse struct {
	DeploymentResp *admin.ApiSearchDeploymentResponse
	HTTPResponse   *http.Response
	Err            error
}
