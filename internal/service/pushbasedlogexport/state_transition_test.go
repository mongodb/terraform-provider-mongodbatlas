package pushbasedlogexport_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312003/admin"
	"go.mongodb.org/atlas-sdk/v20250312003/mockadmin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/pushbasedlogexport"
)

var (
	activeState           = "ACTIVE"
	unconfiguredState     = "UNCONFIGURED"
	initiatingState       = "INITIATING"
	bucketVerifiedState   = "BUCKET_VERIFIED"
	assumeRoleFailedState = "ASSUME_ROLE_FAILED"
	unknown               = ""
	sc500                 = conversion.IntPtr(500)
	currentTime           = time.Now()
)

var testTimeoutConfig = retrystrategy.TimeConfig{
	Timeout:    30 * time.Second,
	MinTimeout: 100 * time.Millisecond,
	Delay:      0,
}

type testCase struct {
	expectedState *string
	name          string
	mockResponses []response
	expectedError bool
}

func TestPushBasedLogExportStateTransition(t *testing.T) {
	testCases := []testCase{
		{
			name: "Successful transition to ACTIVE",
			mockResponses: []response{
				{state: &initiatingState},
				{state: &bucketVerifiedState},
				{state: &activeState},
			},
			expectedState: &activeState,
			expectedError: false,
		},
		{
			name: "Error when transitioning to an unknown state",
			mockResponses: []response{
				{state: &initiatingState},
				{state: &assumeRoleFailedState},
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
			m := mockadmin.NewPushBasedLogExportApi(t)
			m.EXPECT().GetPushBasedLogConfiguration(mock.Anything, mock.Anything).Return(admin.GetPushBasedLogConfigurationApiRequest{ApiService: m})

			for _, resp := range tc.mockResponses {
				modelResp, httpResp, err := resp.get()
				m.EXPECT().GetPushBasedLogConfigurationExecute(mock.Anything).Return(modelResp, httpResp, err).Once()
			}
			resp, err := pushbasedlogexport.WaitStateTransition(t.Context(), testProjectID, m, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, responseWithState(tc.expectedState), resp)
		})
	}
}

func TestPushBasedLogExportStateTransitionForDelete(t *testing.T) {
	testCases := []testCase{
		{
			name: "Successful transition to UNCONFIGURED from ACTIVE",
			mockResponses: []response{
				{state: &activeState},
				{state: &unconfiguredState},
			},
			expectedError: false,
		},
		{
			name: "Successful transition to UNCONFIGURED",
			mockResponses: []response{
				{state: &unconfiguredState},
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
				{state: &activeState},
				{state: &unknown},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := mockadmin.NewPushBasedLogExportApi(t)
			m.EXPECT().GetPushBasedLogConfiguration(mock.Anything, mock.Anything).Return(admin.GetPushBasedLogConfigurationApiRequest{ApiService: m})

			for _, resp := range tc.mockResponses {
				modelResp, httpResp, err := resp.get()
				m.EXPECT().GetPushBasedLogConfigurationExecute(mock.Anything).Return(modelResp, httpResp, err).Once()
			}
			err := pushbasedlogexport.WaitResourceDelete(t.Context(), testProjectID, m, testTimeoutConfig)
			assert.Equal(t, tc.expectedError, err != nil)
		})
	}
}

type response struct {
	state      *string
	statusCode *int
	err        error
}

func (r *response) get() (*admin.PushBasedLogExportProject, *http.Response, error) {
	var httpResp *http.Response
	if r.statusCode != nil {
		httpResp = &http.Response{StatusCode: *r.statusCode}
	}
	return responseWithState(r.state), httpResp, r.err
}

func responseWithState(state *string) *admin.PushBasedLogExportProject {
	if state == nil {
		return nil
	}
	return &admin.PushBasedLogExportProject{
		BucketName: admin.PtrString(testBucketName),
		CreateDate: admin.PtrTime(currentTime),
		IamRoleId:  admin.PtrString(testIAMRoleID),
		PrefixPath: admin.PtrString(testPrefixPath),
		State:      state,
	}
}
