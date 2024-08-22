package encryptionatrestprivateendpoint

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

func waitStateTransition(ctx context.Context, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) (*admin.EARPrivateEndpoint, error) {
	return WaitStateTransitionWithMinTimeout(ctx, 1*time.Minute, projectID, cloudProvider, endpointID, client)
}

func WaitStateTransitionWithMinTimeout(ctx context.Context, minTimeout time.Duration, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) (*admin.EARPrivateEndpoint, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{retrystrategy.RetryStrategyInitiatingState},
		Target:     []string{retrystrategy.RetryStrategyPendingAcceptanceState, retrystrategy.RetryStrategyActiveState},
		Refresh:    refreshFunc(ctx, projectID, cloudProvider, endpointID, client),
		Timeout:    20 * time.Minute,
		MinTimeout: minTimeout,
		Delay:      0,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if privateEndpointResp, ok := result.(*admin.EARPrivateEndpoint); ok && privateEndpointResp != nil {
		return privateEndpointResp, nil
	}
	return nil, errors.New("did not obtain valid result when waiting for state transition")
}

func WaitDeleteStateTransition(ctx context.Context, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) error {
	return WaitDeleteStateTransitionWithMinTimeout(ctx, 1*time.Minute, projectID, cloudProvider, endpointID, client)
}

func WaitDeleteStateTransitionWithMinTimeout(ctx context.Context, minTimeout time.Duration, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{retrystrategy.RetryStrategyPendingAcceptanceState, retrystrategy.RetryStrategyActiveState, retrystrategy.RetryStrategyPendingRecreationState},
		Target:     []string{retrystrategy.RetryStrategyDeletedState},
		Refresh:    refreshFunc(ctx, projectID, cloudProvider, endpointID, client),
		Timeout:    20 * time.Minute,
		MinTimeout: minTimeout,
		Delay:      0,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshFunc(ctx context.Context, projectID, cloudProvider, endpointID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		model, resp, err := client.GetEncryptionAtRestPrivateEndpoint(ctx, projectID, cloudProvider, endpointID).Execute()
		if err != nil && model == nil && resp == nil {
			return nil, "", err
		}
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return "", retrystrategy.RetryStrategyDeletedState, nil
			}
			return nil, "", err
		}
		status := model.GetStatus()
		return model, status, nil
	}
}
