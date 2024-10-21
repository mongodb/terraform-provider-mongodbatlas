package flexcluster_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
	"go.mongodb.org/atlas-sdk/v20240805004/mockadmin"
)

var (
	IdleState      = "IDLE"
	CreatingState  = "CREATING"
	UpdatingState  = "UPDATING"
	DeletingState  = "DELETING"
	RepairingState = "REPAIRING"
	sc500          = conversion.IntPtr(500)
	sc200          = conversion.IntPtr(200)
	sc404          = conversion.IntPtr(404)
	clusterName    = "clusterName"
	requestParams  = &admin.GetFlexClusterApiParams{
		GroupId: "groupId",
		Name:    clusterName,
	}
)

type testCase struct {
	expectedState *string
	name          string
	mockResponses []response
	desiredStates []string
	pendingStates []string
	expectedError bool
}

func TestFlexClusterStateTransition(t *testing.T) {
	testCases := []testCase{
		{
			name: "Successful transition to IDLE",
			mockResponses: []response{
				{state: &CreatingState, statusCode: sc200},
				{state: &IdleState, statusCode: sc200},
			},
			expectedState: &IdleState,
			expectedError: false,
			desiredStates: []string{IdleState},
			pendingStates: []string{CreatingState},
		},
		{
			name: "Error when API returns 5XX",
			mockResponses: []response{
				{statusCode: sc500, err: errors.New("Internal server error")},
			},
			expectedState: nil,
			expectedError: true,
			desiredStates: []string{IdleState},
			pendingStates: []string{CreatingState},
		},
		{
			name: "Error when API returns 404",
			mockResponses: []response{
				{statusCode: sc404, err: errors.New("Not found")},
			},
			expectedState: nil,
			expectedError: true,
			desiredStates: []string{IdleState},
			pendingStates: []string{UpdatingState},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := mockadmin.NewFlexClustersApi(t)
			m.EXPECT().GetFlexClusterWithParams(mock.Anything, mock.Anything).Return(admin.GetFlexClusterApiRequest{ApiService: m})

			for _, resp := range tc.mockResponses {
				modelResp, httpResp, err := resp.get()
				m.EXPECT().GetFlexClusterExecute(mock.Anything).Return(modelResp, httpResp, err).Once()
			}
			resp, err := flexcluster.WaitStateTransition(context.Background(), requestParams, m, tc.pendingStates, tc.desiredStates)
			assert.Equal(t, tc.expectedError, err != nil)
			if resp != nil {
				assert.Equal(t, *tc.expectedState, *resp.StateName)
			}
		})
	}
}

type response struct {
	state      *string
	statusCode *int
	err        error
}

func (r *response) get() (*admin.FlexClusterDescription20250101, *http.Response, error) {
	var httpResp *http.Response
	if r.statusCode != nil {
		httpResp = &http.Response{
			StatusCode: *r.statusCode,
		}
	}
	return responseWithState(r.state), httpResp, r.err
}

func responseWithState(state *string) *admin.FlexClusterDescription20250101 {
	if state == nil {
		return nil
	}
	return &admin.FlexClusterDescription20250101{
		Name:      &clusterName,
		StateName: state,
	}
}
