package streamconnection

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
)

const (
	defaultCreateUpdateTimeout = 20 * time.Minute // The amount of time to wait before timeout for create/update
)

// Connection state constants
const (
	StatePending  = "PENDING"
	StateReady    = "READY"
	StateDeleting = "DELETING"
	StateFailed   = "FAILED"
)

func DeleteStreamConnection(ctx context.Context, api admin.StreamsApi, projectID, instanceName, connectionName string, timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		resp, err := api.DeleteStreamConnection(ctx, projectID, instanceName, connectionName).Execute()
		if err == nil {
			return nil
		}
		if admin.IsErrorCode(err, "STREAM_KAFKA_CONNECTION_IS_DEPLOYING") {
			return retry.RetryableError(err)
		}
		if validate.StatusNotFound(resp) {
			return nil
		}
		return retry.NonRetryableError(err)
	})
}

// WaitStateTransition waits for a stream connection to reach a READY or FAILED state after create or update operations.
// It polls the GET endpoint until the connection is no longer in a PENDING state.
// Some connections may be READY immediately, so there is no minimum wait time.
func WaitStateTransition(ctx context.Context, projectID, workspaceName, connectionName string, client admin.StreamsApi) (*admin.StreamsConnection, error) {
	return WaitStateTransitionWithTimeout(ctx, projectID, workspaceName, connectionName, client, defaultCreateUpdateTimeout)
}

// WaitStateTransitionWithTimeout waits for a stream connection to reach a READY or FAILED state with a custom timeout.
func WaitStateTransitionWithTimeout(ctx context.Context, projectID, workspaceName, connectionName string, client admin.StreamsApi, timeout time.Duration) (*admin.StreamsConnection, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{StatePending},
		Target:  []string{StateReady, StateFailed},
		Refresh: streamConnectionRefreshFunc(ctx, projectID, workspaceName, connectionName, client),
		Timeout: timeout,
		Delay:   0,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if connectionResp, ok := result.(*admin.StreamsConnection); ok {
		return connectionResp, nil
	}
	return nil, errors.New("did not obtain valid result when waiting for stream connection state transition")
}

// streamConnectionRefreshFunc returns a function that polls the stream connection state.
func streamConnectionRefreshFunc(ctx context.Context, projectID, workspaceName, connectionName string, client admin.StreamsApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		model, resp, err := client.GetStreamConnection(ctx, projectID, workspaceName, connectionName).Execute()
		if err != nil && model == nil && resp == nil {
			return nil, "", err
		}
		if err != nil {
			if validate.StatusNotFound(resp) {
				return nil, "", fmt.Errorf("stream connection %s was not found", connectionName)
			}
			return nil, "", err
		}
		state := model.GetState()
		if state == "" {
			// If state is not present in the response, assume the connection is ready
			return model, StateReady, nil
		}
		return model, state, nil
	}
}
