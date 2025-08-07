package encryptionatrestprivateendpoint_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/mockadmin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/encryptionatrestprivateendpoint"
)

type testCase struct {
	expectedState *string
	mockResponses []response
	expectedError bool
}

func TestStateTransition(t *testing.T) {
	testCases := map[string]testCase{
		"Successful transitioning to PENDING_ACCEPTANCE": {
			mockResponses: []response{
				{state: conversion.StringPtr(retrystrategy.RetryStrategyInitiatingState)},
				{state: conversion.StringPtr(retrystrategy.RetryStrategyPendingAcceptanceState)},
			},
			expectedState: conversion.StringPtr(retrystrategy.RetryStrategyPendingAcceptanceState),
			expectedError: false,
		},
		"Successful transitioning to ACTIVE": {
			mockResponses: []response{
				{state: conversion.StringPtr(retrystrategy.RetryStrategyInitiatingState)},
				{state: conversion.StringPtr(retrystrategy.RetryStrategyActiveState)},
			},
			expectedState: conversion.StringPtr(retrystrategy.RetryStrategyActiveState),
			expectedError: false,
		},
		"Return model without error when transitioning to FAILED state": {
			mockResponses: []response{
				{state: conversion.StringPtr(retrystrategy.RetryStrategyInitiatingState)},
				{state: conversion.StringPtr(retrystrategy.RetryStrategyFailedState)},
			},
			expectedState: conversion.StringPtr(retrystrategy.RetryStrategyFailedState),
			expectedError: false,
		},
		"Error when API responds with error": {
			mockResponses: []response{
				{statusCode: admin.PtrInt(http.StatusInternalServerError), err: errors.New("Internal server error")},
			},
			expectedState: nil,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			m := mockadmin.NewEncryptionAtRestUsingCustomerKeyManagementApi(t)
			m.EXPECT().GetEncryptionAtRestPrivateEndpoint(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(admin.GetEncryptionAtRestPrivateEndpointApiRequest{ApiService: m})

			for _, resp := range tc.mockResponses {
				modelResp, httpResp, err := resp.get()
				m.EXPECT().GetEncryptionAtRestPrivateEndpointExecute(mock.Anything).Return(modelResp, httpResp, err).Once()
			}
			resp, err := encryptionatrestprivateendpoint.WaitStateTransitionWithMinTimeout(t.Context(), 1*time.Second, "project-id", "cloud-provider", "endpoint-id", m)
			assert.Equal(t, tc.expectedError, err != nil)
			if resp != nil {
				assert.Equal(t, tc.expectedState, resp.Status)
			}
		})
	}
}

func TestDeleteStateTransition(t *testing.T) {
	testCases := map[string]testCase{
		"Successful transitioning from DELELTING to deleted": {
			mockResponses: []response{
				{state: conversion.StringPtr(retrystrategy.RetryStrategyDeletingState)},
				{statusCode: admin.PtrInt(http.StatusNotFound), err: errors.New("does not exist")},
			},
			expectedError: false,
		},
		"Return model without error when transitioning to FAILED state": {
			mockResponses: []response{
				{state: conversion.StringPtr(retrystrategy.RetryStrategyDeletingState)},
				{state: conversion.StringPtr(retrystrategy.RetryStrategyFailedState)},
			},
			expectedError: false,
			expectedState: conversion.StringPtr(retrystrategy.RetryStrategyFailedState),
		},
		"Error when API responds with error": {
			mockResponses: []response{
				{statusCode: admin.PtrInt(http.StatusInternalServerError), err: errors.New("Internal server error")},
			},
			expectedState: nil,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			m := mockadmin.NewEncryptionAtRestUsingCustomerKeyManagementApi(t)
			m.EXPECT().GetEncryptionAtRestPrivateEndpoint(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(admin.GetEncryptionAtRestPrivateEndpointApiRequest{ApiService: m})

			for _, resp := range tc.mockResponses {
				modelResp, httpResp, err := resp.get()
				m.EXPECT().GetEncryptionAtRestPrivateEndpointExecute(mock.Anything).Return(modelResp, httpResp, err).Once()
			}
			resp, err := encryptionatrestprivateendpoint.WaitDeleteStateTransitionWithMinTimeout(t.Context(), 1*time.Second, "project-id", "cloud-provider", "endpoint-id", m)
			assert.Equal(t, tc.expectedError, err != nil)
			if resp != nil {
				assert.Equal(t, tc.expectedState, resp.Status)
			}
		})
	}
}

type response struct {
	state      *string
	statusCode *int
	err        error
}

func (r *response) get() (*admin.EARPrivateEndpoint, *http.Response, error) {
	var httpResp *http.Response
	if r.statusCode != nil {
		httpResp = &http.Response{StatusCode: *r.statusCode}
	}
	return &admin.EARPrivateEndpoint{Status: r.state}, httpResp, r.err
}
