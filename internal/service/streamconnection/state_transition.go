package streamconnection

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
)

const (
	defaultCreateUpdateTimeout = 20 * time.Minute // The amount of time to wait before timeout for create/update
	notFoundChecks             = 3                // Number of consecutive 404s allowed before failing (~1.4s with exponential backoff)
)

// Connection state constants
const (
	StatePending  = "PENDING"
	StateReady    = "READY"
	StateDeleting = "DELETING"
	StateFailed   = "FAILED"
	StateDeleted  = "DELETED" // Virtual state used when resource is not found (404)
)

func DeleteStreamConnection(ctx context.Context, api admin.StreamsApi, projectID, instanceName, connectionName string, timeout time.Duration) error {
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
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
	if err != nil {
		return err
	}

	// Wait for delete to complete - some connections (e.g., Kafka VPC) are deleted asynchronously
	// and go through a DELETING state before being fully removed
	model, err := WaitDeleteStateTransition(ctx, projectID, instanceName, connectionName, api, timeout)
	if err != nil {
		return err
	}
	if model != nil && model.GetState() == StateFailed {
		return errors.New("stream connection deletion failed")
	}
	return nil
}

// WaitDeleteStateTransition waits for a stream connection to be fully deleted.
// It polls the GET endpoint until the connection returns 404 (not found) or reaches a FAILED state.
// Returns the final model so the caller can inspect the state (e.g., to check for FAILED).
func WaitDeleteStateTransition(ctx context.Context, projectID, workspaceName, connectionName string, client admin.StreamsApi, timeout time.Duration) (*admin.StreamsConnection, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{StateDeleting, StateReady, StatePending},
		Target:     []string{StateDeleted, StateFailed},
		Refresh:    deleteStreamConnectionRefreshFunc(ctx, projectID, workspaceName, connectionName, client),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      0,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if connectionResp, ok := result.(*admin.StreamsConnection); ok {
		return connectionResp, nil
	}
	return nil, nil
}

// deleteStreamConnectionRefreshFunc returns a function that polls the stream connection state during deletion.
func deleteStreamConnectionRefreshFunc(ctx context.Context, projectID, workspaceName, connectionName string, client admin.StreamsApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		model, resp, err := client.GetStreamConnection(ctx, projectID, workspaceName, connectionName).Execute()
		if err != nil && model == nil && resp == nil {
			return nil, "", err
		}
		if err != nil {
			if validate.StatusNotFound(resp) {
				// Resource is deleted
				return &admin.StreamsConnection{}, StateDeleted, nil
			}
			return nil, "", err
		}
		state := model.GetState()
		if state == "" {
			// If state is not present, treat as still existing (will continue polling)
			return model, StateReady, nil
		}
		return model, state, nil
	}
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
		Pending:        []string{StatePending},
		Target:         []string{StateReady, StateFailed},
		Refresh:        streamConnectionRefreshFunc(ctx, projectID, workspaceName, connectionName, client),
		Timeout:        timeout,
		Delay:          0,
		NotFoundChecks: notFoundChecks,
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
				// Return nil to trigger NotFoundChecks for eventual consistency after creation
				// After notFoundChecks consecutive 404s (~1.4s), it will fail fast
				return nil, StatePending, nil
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
