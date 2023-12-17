package searchdeployment_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchdeployment"
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
			expectedError:  true,
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
			mockService := MockSearchDeploymentService{
				MockResponses: tc.mockResponses,
			}

			resp, err := searchdeployment.WaitSearchNodeStateTransition(context.Background(), dummyProjectID, "Cluster0", &mockService, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, tc.expectedResult, resp)
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
			mockService := MockSearchDeploymentService{
				MockResponses: tc.mockResponses,
			}

			err := searchdeployment.WaitSearchNodeDelete(context.Background(), dummyProjectID, clusterName, &mockService, testTimeoutConfig)

			assert.Equal(t, tc.expectedError, err != nil)
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

type MockSearchDeploymentService struct {
	MockResponses []SearchDeploymentResponse
	index         int
}

func (a *MockSearchDeploymentService) GetAtlasSearchDeployment(ctx context.Context, groupID, clusterName string) (*admin.ApiSearchDeploymentResponse, *http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return resp.DeploymentResp, resp.HTTPResponse, resp.Err
}

type SearchDeploymentResponse struct {
	DeploymentResp *admin.ApiSearchDeploymentResponse
	HTTPResponse   *http.Response
	Err            error
}
