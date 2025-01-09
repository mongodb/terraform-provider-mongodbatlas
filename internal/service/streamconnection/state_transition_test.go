package streamconnection_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20241113004/admin"
	"go.mongodb.org/atlas-sdk/v20241113004/mockadmin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil, &genericErr)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil, &genericErr)
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, nil, nil)
	err := streamconnection.DeleteStreamConnection(context.Background(), m, projectID, instanceName, connectionName, time.Minute)
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
	m.EXPECT().DeleteStreamConnectionExecute(mock.Anything).Once().Return(nil, &http.Response{StatusCode: 404}, nil)
	err := streamconnection.DeleteStreamConnection(context.Background(), m, projectID, instanceName, connectionName, time.Minute)
	assert.NoError(t, err)
}
