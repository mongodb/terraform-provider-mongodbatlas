package streamconnection

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

const (
	defaultTimeout    = 1 * time.Minute // The amount of time to wait before timeout
	defaultMinTimeout = 5 * time.Second // Smallest time to wait before refreshes
)

func WaitDeletetateTransitionWithMinTimeout(ctx context.Context, minTimeout time.Duration, projectID, instanceName, connectionID string, client admin.StreamsApi) (*admin.StreamsConnection, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{retrystrategy.RetryStrategyDeletingState, retrystrategy.RetryStrategyDeployingState},
		Target:     []string{retrystrategy.RetryStrategyDeletedState, retrystrategy.RetryStrategyFailedState},
		Refresh:    refreshFunc(ctx, projectID, instanceName, connectionID, client),
		Timeout:    defaultTimeout,
		MinTimeout: minTimeout,
		Delay:      0,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if streamConnection, ok := result.(*admin.StreamsConnection); ok {
		return streamConnection, nil
	}
	return nil, errors.New("did not obtain valid result when waiting for state transition")
}

func refreshFunc(ctx context.Context, projectID, instanceName, connectionID string, client admin.StreamsApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		model, resp, err := client.GetStreamConnection(ctx, projectID, instanceName, connectionID).Execute()
		if err != nil && model == nil && resp == nil {
			return nil, "", err
		}
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return &admin.StreamsConnection{}, retrystrategy.RetryStrategyDeletedState, nil
			}
			if admin.IsErrorCode(err, "STREAM_KAFKA_CONNECTION_IS_DEPLOYING") {
				return nil, retrystrategy.RetryStrategyDeployingState, nil
			}
			return nil, "", err
		}
		return model, retrystrategy.RetryStrategyDeletingState, nil
	}
}

func DeleteStreamConnection(ctx context.Context, api admin.StreamsApi, projectID, instanceName, connectionName string, timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, resp, err := api.DeleteStreamConnection(ctx, projectID, instanceName, connectionName).Execute()
		if err == nil {
			return nil
		}
		if admin.IsErrorCode(err, "STREAM_KAFKA_CONNECTION_IS_DEPLOYING") {
			return retry.RetryableError(err)
		}
		if resp != nil && resp.StatusCode == 404 {
			return nil
		}
		return retry.NonRetryableError(err)
	})
}
