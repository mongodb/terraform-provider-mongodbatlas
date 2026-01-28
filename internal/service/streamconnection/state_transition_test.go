package streamconnection_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"

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
		genericErr = admin.GenericOpenAPIError{}
	)
	genericErr.SetError("error")
	genericErr.SetModel(errDeleteInProgress)
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Times(3)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, &genericErr)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, &genericErr)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil)
	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, instanceName, connectionName, time.Minute)
	assert.NoError(t, err)
}

func TestStreamConnectionDeletion404(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		instanceName   = "instanceName"
		connectionName = "connectionName"
	)
	m.EXPECT().DeleteStreamConnection(mock.Anything, projectID, instanceName, connectionName).Return(admin.DeleteStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(&http.Response{StatusCode: http.StatusNotFound}, nil)
	err := streamconnection.DeleteStreamConnection(t.Context(), m, projectID, instanceName, connectionName, time.Minute)
	assert.NoError(t, err)
}

func TestWaitStateTransitionSuccess(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
	)
	expectedConnection := &admin.StreamsConnection{
		Name:  admin.PtrString(connectionName),
		Type:  admin.PtrString("Kafka"),
		State: admin.PtrString("READY"),
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
		Name:  admin.PtrString(connectionName),
		Type:  admin.PtrString("Kafka"),
		State: admin.PtrString("PENDING"),
	}
	readyConnection := &admin.StreamsConnection{
		Name:  admin.PtrString(connectionName),
		Type:  admin.PtrString("Kafka"),
		State: admin.PtrString("READY"),
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

func TestWaitStateTransition404(t *testing.T) {
	var (
		m              = mockadmin.NewStreamsApi(t)
		projectID      = "projectID"
		workspaceName  = "workspaceName"
		connectionName = "connectionName"
		genericErr     = admin.GenericOpenAPIError{}
	)
	genericErr.SetError("not found")
	m.EXPECT().GetStreamConnection(mock.Anything, projectID, workspaceName, connectionName).Return(admin.GetStreamConnectionApiRequest{ApiService: m}).Once()
	m.EXPECT().GetStreamConnectionExecute(mock.Anything).Once().Return(nil, &http.Response{StatusCode: http.StatusNotFound}, &genericErr)

	// When connection is not found, WaitStateTransition should return an error immediately
	result, err := streamconnection.WaitStateTransitionWithTimeout(t.Context(), projectID, workspaceName, connectionName, m, 1*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "was not found")
	assert.Nil(t, result)
}
