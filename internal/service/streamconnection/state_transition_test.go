package streamconnection_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
	"go.mongodb.org/atlas-sdk/v20250312014/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamconnection"
)

func TestStreamConnectionDeletion(t *testing.T) {
	var (
		m                   = mockadmin.NewStreamsApi(t)
		projectID           = "projectID"
		instanceName        = "instanceName"
		connectionName      = "connectionName"
		errDeleteInProgress = admin.ApiError{
			ErrorCode: "STREAM_KAFKA_CONNECTION_IS_DEPLOYING",
			Error:     409,
		}
		genericErr  = admin.GenericOpenAPIError{}
		notFoundErr = admin.GenericOpenAPIError{}
	)
	genericErr.SetError("error")
	genericErr.SetModel(errDeleteInProgress)
	notFoundErr.SetError("not found")

	// Delete retries until success
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Times(3)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, &genericErr)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, &genericErr)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil)

	// After delete succeeds, wait for resource to be deleted (404 response)
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(nil, &http.Response{StatusCode: http.StatusNotFound}, &notFoundErr)

	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, instanceName, connectionName, time.Minute)
	assert.NoError(t, err)
}

func TestStreamConnectionDeletion404(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		instanceName   = "instanceName"
		connectionName = "connectionName"
		notFoundErr    = admin.GenericOpenAPIError{}
	)
	notFoundErr.SetError("not found")

	// Delete returns 404 immediately (resource doesn't exist)
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(&http.Response{StatusCode: http.StatusNotFound}, nil)

	// Wait for delete still checks, gets 404 confirming resource is gone
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(nil, &http.Response{StatusCode: http.StatusNotFound}, &notFoundErr)

	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, instanceName, connectionName, time.Minute)
	assert.NoError(t, err)
}

// TestStreamConnectionDeletionWithDeletingState tests the async delete case where the
// connection goes through a DELETING state before being fully removed.
func TestStreamConnectionDeletionWithDeletingState(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		instanceName   = "instanceName"
		connectionName = "connectionName"
		notFoundErr    = admin.GenericOpenAPIError{}
	)
	notFoundErr.SetError("not found")

	deletingConnection := &admin.StreamsConnection{
		Name:  admin.PtrString(connectionName),
		Type:  admin.PtrString("Kafka"),
		State: admin.PtrString("DELETING"),
	}

	// Delete call succeeds immediately
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil)

	// Wait for delete polls: first returns DELETING, second returns 404
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Times(2)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(deletingConnection, nil, nil)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(nil, &http.Response{StatusCode: http.StatusNotFound}, &notFoundErr)

	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, instanceName, connectionName, time.Minute)
	assert.NoError(t, err)
}

// TestStreamConnectionDeletionFailed tests that deletion returns an error when the connection
// transitions to FAILED state during deletion.
func TestStreamConnectionDeletionFailed(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		instanceName   = "instanceName"
		connectionName = "connectionName"
	)

	failedConnection := &admin.StreamsConnection{
		Name:  admin.PtrString(connectionName),
		Type:  admin.PtrString("Kafka"),
		State: admin.PtrString("FAILED"),
	}

	// Delete call succeeds immediately
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil)

	// Wait for delete polls: returns FAILED state
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(failedConnection, nil, nil)

	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, instanceName, connectionName, time.Minute)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "deletion failed")
}

func TestWaitStateTransitionSuccess(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
	)
	expectedConnection := &admin.StreamsConnection{
		Name:  new(connectionName),
		Type:  new("Kafka"),
		State: new("READY"),
	}
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(expectedConnection, nil, nil)

	result, err := streamconnection.WaitStateTransitionWithTimeout(t.Context(), projectID, workspaceName, connectionName, m, 30*time.Second)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, connectionName, *result.Name)
}

func TestWaitStateTransitionPendingToReady(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
	)
	pendingConnection := &admin.StreamsConnection{
		Name:  new(connectionName),
		Type:  new("Kafka"),
		State: new("PENDING"),
	}
	readyConnection := &admin.StreamsConnection{
		Name:  new(connectionName),
		Type:  new("Kafka"),
		State: new("READY"),
	}
	// First call returns PENDING, second call returns READY
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Times(2)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(pendingConnection, nil, nil)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(readyConnection, nil, nil)

	result, err := streamconnection.WaitStateTransitionWithTimeout(t.Context(), projectID, workspaceName, connectionName, m, 30*time.Second)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "READY", *result.State)
}

func TestWaitStateTransition404ThenReady(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
		genericErr     = admin.GenericOpenAPIError{}
	)
	genericErr.SetError("not found")
	readyConnection := &admin.StreamsConnection{
		Name:  new(connectionName),
		Type:  new("Kafka"),
		State: new("READY"),
	}
	// First call returns 404 (eventual consistency), second call returns READY
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Times(2)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(nil, &http.Response{StatusCode: http.StatusNotFound}, &genericErr)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(readyConnection, nil, nil)

	// 404 should be treated as PENDING state to handle eventual consistency after creation
	result, err := streamconnection.WaitStateTransitionWithTimeout(t.Context(), projectID, workspaceName, connectionName, m, 30*time.Second)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "READY", *result.State)
}

func TestWaitStateTransitionNotFoundExceedsLimit(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
		genericErr     = admin.GenericOpenAPIError{}
	)
	genericErr.SetError("not found")
	// Return 404 more times than NotFoundChecks allows (3 checks + 1 to trigger error = 4 calls)
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Times(4)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Return(nil, &http.Response{StatusCode: http.StatusNotFound}, &genericErr).Times(4)

	// After exceeding NotFoundChecks, should return an error
	result, err := streamconnection.WaitStateTransitionWithTimeout(t.Context(), projectID, workspaceName, connectionName, m, 30*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "couldn't find resource")
	assert.Nil(t, result)
}

// TestWaitStateTransitionFailed verifies that when a connection reaches the FAILED state,
// the function returns the connection (so Terraform can show the failure details).
func TestWaitStateTransitionFailed(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
	)
	failedConnection := &admin.StreamsConnection{
		Name:  new(connectionName),
		Type:  new("Kafka"),
		State: new("FAILED"),
	}
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(failedConnection, nil, nil)

	// FAILED is a target state, so it should return successfully with the failed connection
	result, err := streamconnection.WaitStateTransitionWithTimeout(t.Context(), projectID, workspaceName, connectionName, m, 30*time.Second)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "FAILED", *result.State)
}

// TestWaitStateTransitionEmptyStateAssumesReady verifies backward compatibility:
// when the API response doesn't include a state field, we assume the connection is READY.
func TestWaitStateTransitionEmptyStateAssumesReady(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
	)
	// Connection without State field (backward compatibility with older API versions)
	connectionWithoutState := &admin.StreamsConnection{
		Name: new(connectionName),
		Type: new("Kafka"),
		// State is not set
	}
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(connectionWithoutState, nil, nil)

	// Empty state should be treated as READY
	result, err := streamconnection.WaitStateTransitionWithTimeout(t.Context(), projectID, workspaceName, connectionName, m, 30*time.Second)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, connectionName, *result.Name)
}

// TestWaitStateTransitionAPIError verifies that non-404 API errors fail immediately
// rather than retrying indefinitely.
func TestWaitStateTransitionAPIError(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
		genericErr     = admin.GenericOpenAPIError{}
	)
	genericErr.SetError("internal server error")
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(nil, &http.Response{StatusCode: http.StatusInternalServerError}, &genericErr)

	// Non-404 errors should fail immediately
	result, err := streamconnection.WaitStateTransitionWithTimeout(t.Context(), projectID, workspaceName, connectionName, m, 30*time.Second)
	require.Error(t, err)
	assert.Nil(t, result)
}
