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
		workspaceName       = "workspaceName"
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
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Times(3)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, &genericErr)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, &genericErr)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil)

	// After delete succeeds, wait for resource to be deleted (404 response)
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(nil, &http.Response{StatusCode: http.StatusNotFound}, &notFoundErr)

	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, workspaceName, connectionName, time.Minute)
	assert.NoError(t, err)
}

func TestStreamConnectionDeletionWithDeletingState(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
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
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil)

	// Wait for delete polls: first returns DELETING, second returns 404
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Times(2)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(deletingConnection, nil, nil)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(nil, &http.Response{StatusCode: http.StatusNotFound}, &notFoundErr)

	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, workspaceName, connectionName, time.Minute)
	assert.NoError(t, err)
}

func TestStreamConnectionDeletionNonRetryableError(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
		genericErr     = admin.GenericOpenAPIError{}
	)
	genericErr.SetError("internal server error")

	// Delete returns non-retryable error
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(&http.Response{StatusCode: http.StatusInternalServerError}, &genericErr)

	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, workspaceName, connectionName, time.Minute)
	require.Error(t, err)
}

func TestStreamConnectionDeletionFailed(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
	)

	failedConnection := &admin.StreamsConnection{
		Name:  admin.PtrString(connectionName),
		Type:  admin.PtrString("Kafka"),
		State: admin.PtrString("FAILED"),
	}

	// Delete call succeeds immediately
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil)

	// Wait for delete polls: returns FAILED state
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(failedConnection, nil, nil)

	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, workspaceName, connectionName, time.Minute)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "stream connection deletion failed for connection 'connectionName' in workspace 'workspaceName' (project: projectID)")
}

func TestWaitStateTransition(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
	)

	pendingConnection := &admin.StreamsConnection{
		Name:  admin.PtrString(connectionName),
		Type:  admin.PtrString("Kafka"),
		State: admin.PtrString("PENDING"),
	}
	readyConnection := &admin.StreamsConnection{
		Name:  admin.PtrString(connectionName),
		Type:  admin.PtrString("Kafka"),
		State: admin.PtrString("READY"),
	}

	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Times(2)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(pendingConnection, nil, nil)
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(readyConnection, nil, nil)

	pendingStates := []string{streamconnection.StatePending}
	targetStates := []string{streamconnection.StateReady, streamconnection.StateFailed}
	result, err := streamconnection.WaitStateTransition(t.Context(), projectID, workspaceName, connectionName, m, 30*time.Second, pendingStates, targetStates)
	require.NoError(t, err)
	assert.Equal(t, "READY", result.GetState())
}

func TestWaitStateTransitionFailed(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
	)

	failedConnection := &admin.StreamsConnection{
		Name:  admin.PtrString(connectionName),
		Type:  admin.PtrString("Kafka"),
		State: admin.PtrString("FAILED"),
	}

	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(failedConnection, nil, nil)

	// WaitStateTransition returns the model without error - caller must check the state
	pendingStates := []string{streamconnection.StatePending}
	targetStates := []string{streamconnection.StateReady, streamconnection.StateFailed}
	result, err := streamconnection.WaitStateTransition(t.Context(), projectID, workspaceName, connectionName, m, 30*time.Second, pendingStates, targetStates)
	require.NoError(t, err)
	assert.Equal(t, "FAILED", result.GetState())
}
