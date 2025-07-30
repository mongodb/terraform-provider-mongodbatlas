package flexcluster_test

import (
	"errors"
	"net/http"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"
	"go.mongodb.org/atlas-sdk/v20250312005/mockadmin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

var (
	IdleState     = "IDLE"
	CreatingState = "CREATING"
	UpdatingState = "UPDATING"
	DeletingState = "DELETING"
	DeletedState  = "DELETED"
	UnknownState  = ""
	sc500         = conversion.IntPtr(500)
	sc200         = conversion.IntPtr(200)
	sc404         = conversion.IntPtr(404)
	clusterName   = "clusterName"
	requestParams = &admin.GetFlexClusterApiParams{
		GroupId: "groupId",
		Name:    clusterName,
	}
)

type testCase struct {
	expectedState   *string
	name            string
	mockResponses   []response
	desiredStates   []string
	pendingStates   []string
	expectedError   bool
	isUpgradeFromM0 bool
}

func TestFlexClusterStateTransition(t *testing.T) {
	testCases := []testCase{
		{
			name: "Successful transition to IDLE",
			mockResponses: []response{
				{state: &CreatingState, statusCode: sc200},
				{state: &IdleState, statusCode: sc200},
			},
			expectedState:   &IdleState,
			expectedError:   false,
			desiredStates:   []string{IdleState},
			pendingStates:   []string{CreatingState},
			isUpgradeFromM0: false,
		},
		{
			name: "Successful transition to IDLE during cluster (M0) upgrade to Flex",
			mockResponses: []response{
				{state: &UpdatingState, statusCode: sc200},
				{state: &IdleState, statusCode: sc200},
			},
			expectedState:   &IdleState,
			expectedError:   false,
			desiredStates:   []string{IdleState},
			pendingStates:   []string{UpdatingState},
			isUpgradeFromM0: true,
		},
		{
			name: "Error when API returns 5XX",
			mockResponses: []response{
				{statusCode: sc500, err: errors.New("Internal server error")},
			},
			expectedState:   nil,
			expectedError:   true,
			desiredStates:   []string{IdleState},
			pendingStates:   []string{CreatingState},
			isUpgradeFromM0: false,
		},
		{
			name: "Deleted state when API returns 404",
			mockResponses: []response{
				{state: &DeletingState, statusCode: sc200},
				{statusCode: sc404, err: errors.New("Not found")},
			},
			expectedState:   nil,
			expectedError:   true,
			desiredStates:   []string{IdleState},
			pendingStates:   []string{DeletingState},
			isUpgradeFromM0: false,
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
			resp, err := flexcluster.WaitStateTransition(t.Context(), requestParams, m, tc.pendingStates, tc.desiredStates, tc.isUpgradeFromM0, nil)
			assert.Equal(t, tc.expectedError, err != nil)
			if resp != nil {
				assert.Equal(t, *tc.expectedState, *resp.StateName)
			}
		})
	}
}

func TestFlexClusterStateTransitionForDelete(t *testing.T) {
	testCases := []testCase{
		{
			name: "Successful transition to DELETED",
			mockResponses: []response{
				{state: &DeletingState, statusCode: sc200},
				{statusCode: sc404, err: errors.New("Not found")},
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
				{state: &DeletingState},
				{state: &UnknownState},
			},
			expectedError: true,
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
			err := flexcluster.WaitStateTransitionDelete(t.Context(), requestParams, m, nil)
			assert.Equal(t, tc.expectedError, err != nil)
		})
	}
}

type response struct {
	state      *string
	statusCode *int
	err        error
}

func (r *response) get() (*admin.FlexClusterDescription20241113, *http.Response, error) {
	var httpResp *http.Response
	if r.statusCode != nil {
		httpResp = &http.Response{
			StatusCode: *r.statusCode,
		}
	}
	return responseWithState(r.state), httpResp, r.err
}

func responseWithState(state *string) *admin.FlexClusterDescription20241113 {
	if state == nil {
		return nil
	}
	return &admin.FlexClusterDescription20241113{
		Name:      &clusterName,
		StateName: state,
	}
}
