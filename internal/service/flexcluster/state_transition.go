package flexcluster

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func WaitStateTransition(ctx context.Context, requestParams *admin.GetFlexClusterApiParams, client admin.FlexClustersApi, pendingStates, desiredStates []string, isUpgradeFromM0 bool, timeout *time.Duration) (*admin.FlexClusterDescription20241113, error) {
	if timeout == nil {
		timeout = conversion.Pointer(constant.DefaultTimeout)
	}
	stateConf := &retry.StateChangeConf{
		Pending:    pendingStates,
		Target:     desiredStates,
		Refresh:    refreshFunc(ctx, requestParams, client, isUpgradeFromM0),
		Timeout:    *timeout,
		MinTimeout: 3 * time.Second,
		Delay:      0,
	}

	flexClusterResp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	if flexCluster, ok := flexClusterResp.(*admin.FlexClusterDescription20241113); ok && flexCluster != nil {
		return flexCluster, nil
	}

	return nil, errors.New("did not obtain valid result when waiting for flex cluster state transition")
}

func WaitStateTransitionDelete(ctx context.Context, requestParams *admin.GetFlexClusterApiParams, client admin.FlexClustersApi) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{retrystrategy.RetryStrategyDeletingState},
		Target:     []string{retrystrategy.RetryStrategyDeletedState},
		Refresh:    refreshFunc(ctx, requestParams, client, false),
		Timeout:    3 * time.Hour,
		MinTimeout: 3 * time.Second,
		Delay:      0,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshFunc(ctx context.Context, requestParams *admin.GetFlexClusterApiParams, client admin.FlexClustersApi, isUpgradeFromM0 bool) retry.StateRefreshFunc {
	return func() (any, string, error) {
		flexCluster, resp, err := client.GetFlexClusterWithParams(ctx, requestParams).Execute()
		if err != nil {
			if isUpgradeFromM0 && (validate.StatusNotFound(resp) || admin.IsErrorCode(err, "CANNOT_USE_NON_FLEX_CLUSTER_IN_FLEX_API")) {
				return "", retrystrategy.RetryStrategyUpdatingState, nil
			}
			if validate.StatusNotFound(resp) {
				return "", retrystrategy.RetryStrategyDeletedState, nil
			}
			return nil, "", err
		}
		state := flexCluster.GetStateName()
		return flexCluster, state, nil
	}
}
