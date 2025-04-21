package encryptionatrestprivateendpoint

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

const (
	defaultTimeout    = 20 * time.Minute // The amount of time to wait before timeout
	defaultMinTimeout = 30 * time.Second // Smallest time to wait before refreshes
)

func waitStateTransition(ctx context.Context, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) (*admin.EARPrivateEndpoint, error) {
	return WaitStateTransitionWithMinTimeout(ctx, defaultMinTimeout, projectID, cloudProvider, endpointID, client)
}

func WaitStateTransitionWithMinTimeout(ctx context.Context, minTimeout time.Duration, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) (*admin.EARPrivateEndpoint, error) {
	return waitStateTransitionForStates(
		ctx,
		[]string{retrystrategy.RetryStrategyInitiatingState},
		[]string{retrystrategy.RetryStrategyPendingAcceptanceState, retrystrategy.RetryStrategyActiveState, retrystrategy.RetryStrategyFailedState},
		minTimeout, projectID, cloudProvider, endpointID, client)
}

func WaitDeleteStateTransition(ctx context.Context, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) (*admin.EARPrivateEndpoint, error) {
	return WaitDeleteStateTransitionWithMinTimeout(ctx, defaultMinTimeout, projectID, cloudProvider, endpointID, client)
}

func WaitDeleteStateTransitionWithMinTimeout(ctx context.Context, minTimeout time.Duration, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) (*admin.EARPrivateEndpoint, error) {
	return waitStateTransitionForStates(
		ctx,
		[]string{retrystrategy.RetryStrategyDeletingState},
		[]string{retrystrategy.RetryStrategyDeletedState, retrystrategy.RetryStrategyFailedState},
		minTimeout, projectID, cloudProvider, endpointID, client)
}

func waitStateTransitionForStates(ctx context.Context, pending, target []string, minTimeout time.Duration, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) (*admin.EARPrivateEndpoint, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    refreshFunc(ctx, projectID, cloudProvider, endpointID, client),
		Timeout:    defaultTimeout,
		MinTimeout: minTimeout,
		Delay:      0,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if privateEndpointResp, ok := result.(*admin.EARPrivateEndpoint); ok {
		return privateEndpointResp, nil
	}
	return nil, errors.New("did not obtain valid result when waiting for state transition")
}

func refreshFunc(ctx context.Context, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		model, resp, err := client.GetEncryptionAtRestPrivateEndpoint(ctx, projectID, cloudProvider, endpointID).Execute()
		if err != nil && model == nil && resp == nil {
			return nil, "", err
		}
		if err != nil {
			if validate.StatusNotFound(resp) {
				return &admin.EARPrivateEndpoint{}, retrystrategy.RetryStrategyDeletedState, nil
			}
			return nil, "", err
		}
		status := model.GetStatus()
		return model, status, nil
	}
}
