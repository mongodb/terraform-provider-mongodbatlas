package advancedclustertpf

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
)

var (
	RetryMinTimeout   = 1 * time.Minute
	RetryDelay        = 30 * time.Second
	RetryPollInterval = 30 * time.Second
)

func CreateStateChangeConfig(ctx context.Context, connV2 *admin20240805.APIClient, projectID, name, targetState string, timeout time.Duration, extraPending ...string) retry.StateChangeConf {
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
		Refresh:      resourceRefreshFunc(ctx, name, projectID, connV2),
		Timeout:      timeout,
		MinTimeout:   RetryMinTimeout,
		Delay:        RetryDelay,
		PollInterval: RetryPollInterval,
	}
}

func resourceRefreshFunc(ctx context.Context, name, projectID string, connV2 *admin20240805.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		cluster, resp, err := connV2.ClustersApi.GetCluster(ctx, projectID, name).Execute()
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
