package advancedclustertpf

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
)

var (
	RetryMinTimeout   = 1 * time.Minute
	RetryDelay        = 30 * time.Second
	RetryPollInterval = 30 * time.Second
)

func AwaitChanges(ctx context.Context, api admin.ClustersApi, t *timeouts.Value, diags *diag.Diagnostics, projectID, clusterName, changeReason string) (cluster *admin.ClusterDescription20240805) {
	var (
		timeoutDuration time.Duration
		localDiags      diag.Diagnostics
		targetState     = retrystrategy.RetryStrategyIdleState
		extraPending    = []string{}
	)
	switch changeReason {
	case changeReasonCreate:
		timeoutDuration, localDiags = t.Create(ctx, defaultTimeout)
		diags.Append(localDiags...)
	case changeReasonUpdate:
		timeoutDuration, localDiags = t.Update(ctx, defaultTimeout)
		diags.Append(localDiags...)
	case changeReasonDelete:
		timeoutDuration, localDiags = t.Delete(ctx, defaultTimeout)
		diags.Append(localDiags...)
		targetState = retrystrategy.RetryStrategyDeletedState
		extraPending = append(extraPending, retrystrategy.RetryStrategyIdleState)
	default:
		diags.AddError("errorAwaitingChanges", "unknown change reason "+changeReason)
	}
	if diags.HasError() {
		return nil
	}
	stateConf := CreateStateChangeConfig(ctx, api, projectID, clusterName, targetState, timeoutDuration, extraPending...)
	clusterAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) && changeReason == "delete" {
			return nil
		}
		diags.AddError("errorAwaitingCluster", fmt.Sprintf(errorCreate, err))
		return nil
	}
	if targetState == retrystrategy.RetryStrategyDeletedState {
		return nil
	}
	cluster, ok := clusterAny.(*admin.ClusterDescription20240805)
	if !ok {
		diags.AddError("errorAwaitingCluster", fmt.Sprintf(errorCreate, "unexpected type from WaitForStateContext"))
		return nil
	}
	return cluster
}

func CreateStateChangeConfig(ctx context.Context, api admin.ClustersApi, projectID, name, targetState string, timeout time.Duration, extraPending ...string) retry.StateChangeConf {
	return retry.StateChangeConf{
		Pending: slices.Concat([]string{
			retrystrategy.RetryStrategyCreatingState,
			retrystrategy.RetryStrategyUpdatingState,
			retrystrategy.RetryStrategyRepairingState,
			retrystrategy.RetryStrategyRepeatingState,
			retrystrategy.RetryStrategyPendingState,
			retrystrategy.RetryStrategyDeletingState,
		}, extraPending),
		Target:       []string{targetState},
		Refresh:      resourceRefreshFunc(ctx, name, projectID, api),
		Timeout:      timeout,
		MinTimeout:   RetryMinTimeout,
		Delay:        RetryDelay,
		PollInterval: RetryPollInterval,
	}
}

func resourceRefreshFunc(ctx context.Context, name, projectID string, api admin.ClustersApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		cluster, resp, err := api.GetCluster(ctx, projectID, name).Execute()
		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, retrystrategy.RetryStrategyRepeatingState, nil
		}

		if err != nil && cluster == nil && resp == nil {
			return nil, "", err
		}

		if err != nil {
			if resp.StatusCode == 404 {
				return "", retrystrategy.RetryStrategyDeletedState, nil
			}
			if resp.StatusCode == 503 {
				return "", retrystrategy.RetryStrategyPendingState, nil
			}
			return nil, "", err
		}

		state := cluster.GetStateName()
		return cluster, state, nil
	}
}
