package streamconnection

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
)

// Connection state constants
const (
	StatePending  = "PENDING"
	StateReady    = "READY"
	StateDeleting = "DELETING"
	StateFailed   = "FAILED"
	StateDeleted  = "DELETED" // Virtual state used when resource is not found (404)
)

func DeleteStreamConnection(ctx context.Context, api admin.StreamsApi, projectID, workspaceName, connectionName string, timeout time.Duration) error {
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		resp, err := api.DeleteStreamConnection(ctx, projectID, workspaceName, connectionName).Execute()
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
	if err != nil {
		return err
	}

	// Wait for delete to complete - some connections (e.g., Kafka VPC) are deleted asynchronously
	// and go through a DELETING state before being fully removed
	pendingStates := []string{StateDeleting}
	targetStates := []string{StateDeleted, StateFailed}
	model, err := WaitStateTransition(ctx, projectID, workspaceName, connectionName, api, timeout, pendingStates, targetStates)
	if err != nil {
		return err
	}
	if model != nil && model.GetState() == StateFailed {
		return fmt.Errorf("stream connection deletion failed for connection '%s' in workspace '%s' (project: %s)", connectionName, workspaceName, projectID)
	}
	return nil
}

// WaitStateTransition waits for a stream connection to reach the specified target states.
func WaitStateTransition(ctx context.Context, projectID, workspaceName, connectionName string, client admin.StreamsApi, timeout time.Duration, pendingStates, targetStates []string) (*admin.StreamsConnection, error) {
	stateConf := &retry.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: refreshFunc(ctx, projectID, workspaceName, connectionName, client),
		Timeout: timeout,
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

// refreshFunc returns a function that polls the stream connection state.
// Returns StateDeleted when resource is not found (404).
func refreshFunc(ctx context.Context, projectID, workspaceName, connectionName string, client admin.StreamsApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		model, resp, err := client.GetStreamConnection(ctx, projectID, workspaceName, connectionName).Execute()
		if err != nil && model == nil && resp == nil {
			return nil, "", err
		}
		if err != nil {
			if validate.StatusNotFound(resp) {
				return &admin.StreamsConnection{}, StateDeleted, nil
			}
			return nil, "", err
		}
		state := model.GetState()
		return model, state, nil
	}
}
