package streamprocessor_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.mongodb.org/atlas-sdk/v20250312009/mockadmin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
)

var (
	InitiatingState     = "INIT"
	CreatingState       = "CREATING"
	CreatedState        = "CREATED"
	StartedState        = "STARTED"
	StoppedState        = "STOPPED"
	DroppedState        = "DROPPED"
	FailedState         = "FAILED"
	sc500               = conversion.IntPtr(500)
	sc200               = conversion.IntPtr(200)
	sc404               = conversion.IntPtr(404)
	streamProcessorName = "processorName"
	requestParams       = &admin.GetStreamProcessorApiParams{
		GroupId:       "groupId",
		TenantName:    "tenantName",
		ProcessorName: streamProcessorName,
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

func TestStreamProcessorStateTransition(t *testing.T) {
	testCases := []testCase{
		{
			name: "Successful transition to CREATED",
			mockResponses: []response{
				{state: &InitiatingState, statusCode: sc200},
				{state: &CreatingState, statusCode: sc200},
				{state: &CreatedState, statusCode: sc200},
			},
			expectedState: &CreatedState,
			expectedError: false,
			desiredStates: []string{CreatedState},
			pendingStates: []string{InitiatingState, CreatingState},
		},
		{
			name: "Successful transition to STARTED",
			mockResponses: []response{
				{state: &CreatedState, statusCode: sc200},
				{state: &StartedState, statusCode: sc200},
			},
			expectedState: &StartedState,
			expectedError: false,
			desiredStates: []string{StartedState},
			pendingStates: []string{CreatedState, StoppedState},
		},
		{
			name: "Successful transition to STOPPED",
			mockResponses: []response{
				{state: &StartedState, statusCode: sc200},
				{state: &StoppedState, statusCode: sc200},
			},
			expectedState: &StoppedState,
			expectedError: false,
			desiredStates: []string{StoppedState},
			pendingStates: []string{StartedState},
		},
		{
			name: "Error when transitioning to FAILED state",
			mockResponses: []response{
				{state: &InitiatingState, statusCode: sc200},
				{state: &FailedState, statusCode: sc200},
			},
			expectedState: nil,
			expectedError: true,
			desiredStates: []string{CreatedState},
			pendingStates: []string{InitiatingState, CreatingState},
		},
		{
			name: "Error when API responds with error",
			mockResponses: []response{
				{statusCode: sc500, err: errors.New("Internal server error")},
			},
			expectedState: nil,
			expectedError: true,
			desiredStates: []string{CreatedState, FailedState},
			pendingStates: []string{InitiatingState, CreatingState},
		},
		{
			name: "Dropped state when 404 is returned",
			mockResponses: []response{
				{statusCode: sc404, err: errors.New("Not found")},
			},
			expectedState: &DroppedState,
			expectedError: true,
			desiredStates: []string{CreatedState, FailedState},
			pendingStates: []string{InitiatingState, CreatingState},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := mockadmin.NewStreamsApi(t)
			m.EXPECT().GetStreamProcessorWithParams(mock.Anything, mock.Anything).Return(admin.GetStreamProcessorApiRequest{ApiService: m})

			for _, resp := range tc.mockResponses {
				modelResp, httpResp, err := resp.get()
				m.EXPECT().GetStreamProcessorExecute(mock.Anything).Return(modelResp, httpResp, err).Once()
			}
			resp, err := streamprocessor.WaitStateTransition(t.Context(), requestParams, m, tc.pendingStates, tc.desiredStates)
			assert.Equal(t, tc.expectedError, err != nil)
			if resp != nil {
				assert.Equal(t, *tc.expectedState, resp.State)
			}
		})
	}
}

type response struct {
	state      *string
	statusCode *int
	err        error
}

func (r *response) get() (*admin.StreamsProcessorWithStats, *http.Response, error) {
	var httpResp *http.Response
	if r.statusCode != nil {
		httpResp = &http.Response{StatusCode: *r.statusCode}
	}
	return responseWithState(r.state), httpResp, r.err
}

func responseWithState(state *string) *admin.StreamsProcessorWithStats {
	if state == nil {
		return nil
	}
	return &admin.StreamsProcessorWithStats{
		Name:  streamProcessorName,
		State: *state,
	}
}

func TestValidateUpdateStateTransition(t *testing.T) {
	testCases := []struct {
		name           string
		currentState   string
		plannedState   string
		wantErrMsg     string
		wantValidation bool
	}{
		// Same state transitions (always valid)
		{
			name:           "Same state - CREATED to CREATED",
			currentState:   CreatedState,
			plannedState:   CreatedState,
			wantErrMsg:     "",
			wantValidation: true,
		},
		{
			name:           "Same state - STARTED to STARTED",
			currentState:   StartedState,
			plannedState:   StartedState,
			wantErrMsg:     "",
			wantValidation: true,
		},
		{
			name:           "Same state - STOPPED to STOPPED",
			currentState:   StoppedState,
			plannedState:   StoppedState,
			wantErrMsg:     "",
			wantValidation: true,
		},

		// Valid transitions between different states
		{
			name:           "Valid transition - CREATED to STARTED",
			currentState:   CreatedState,
			plannedState:   StartedState,
			wantErrMsg:     "",
			wantValidation: true,
		},
		{
			name:           "Valid transition - STARTED to STOPPED",
			currentState:   StartedState,
			plannedState:   StoppedState,
			wantErrMsg:     "",
			wantValidation: true,
		},
		{
			name:           "Valid transition - STOPPED to STARTED",
			currentState:   StoppedState,
			plannedState:   StartedState,
			wantErrMsg:     "",
			wantValidation: true,
		},

		// Invalid transitions
		{
			name:           "Invalid transition - CREATED to STOPPED",
			currentState:   CreatedState,
			plannedState:   StoppedState,
			wantErrMsg:     fmt.Sprintf(streamprocessor.ErrorUpdateStateTransition, StartedState, StoppedState),
			wantValidation: false,
		},
		{
			name:           "Invalid transition - STARTED to CREATED",
			currentState:   StartedState,
			plannedState:   CreatedState,
			wantErrMsg:     fmt.Sprintf(streamprocessor.ErrorUpdateToCreatedState, StartedState),
			wantValidation: false,
		},
		{
			name:           "Invalid transition - STOPPED to CREATED",
			currentState:   StoppedState,
			plannedState:   CreatedState,
			wantErrMsg:     fmt.Sprintf(streamprocessor.ErrorUpdateToCreatedState, StoppedState),
			wantValidation: false,
		},

		// Edge cases with other states
		{
			name:           "Edge case - FAILED to any state",
			currentState:   FailedState,
			plannedState:   StartedState,
			wantErrMsg:     "",
			wantValidation: true,
		},
		{
			name:           "Edge case - INIT to any state",
			currentState:   InitiatingState,
			plannedState:   CreatedState,
			wantErrMsg:     fmt.Sprintf(streamprocessor.ErrorUpdateToCreatedState, InitiatingState),
			wantValidation: false,
		},
		{
			name:           "Edge case - any state to FAILED",
			currentState:   StoppedState,
			plannedState:   FailedState,
			wantErrMsg:     "",
			wantValidation: true,
		},
		{
			name:           "Edge case - empty state values",
			currentState:   "",
			plannedState:   "",
			wantErrMsg:     "",
			wantValidation: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrMsg, gotValidation := streamprocessor.ValidateUpdateStateTransition(tc.currentState, tc.plannedState)

			if gotErrMsg != tc.wantErrMsg {
				t.Errorf("validateUpdateStateTransition() error message = %q, want %q",
					gotErrMsg, tc.wantErrMsg)
			}

			if gotValidation != tc.wantValidation {
				t.Errorf("validateUpdateStateTransition() validation = %v, want %v",
					gotValidation, tc.wantValidation)
			}
		})
	}
}
