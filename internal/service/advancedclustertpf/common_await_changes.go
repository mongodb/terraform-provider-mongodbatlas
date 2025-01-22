package advancedclustertpf

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

var (
	RetryMinTimeout      = 1 * time.Minute
	RetryDelay           = 30 * time.Second
	RetryPollInterval    = 30 * time.Second
	AwaitDeleteOperation = operationDelete
)

type ClusterReader struct {
	ProjectID   string
	ClusterName string
	Timeout     time.Duration
	IsDelete    bool
}

func AwaitChanges(ctx context.Context, client *config.MongoDBClient, reader *ClusterReader, lastOperation string, diags *diag.Diagnostics) *admin.ClusterDescription20240805 {
	api := client.AtlasV2.ClustersApi
	targetState := retrystrategy.RetryStrategyIdleState
	extraPending := []string{}
	isDelete := reader.IsDelete
	if isDelete {
		targetState = retrystrategy.RetryStrategyDeletedState
		extraPending = append(extraPending, retrystrategy.RetryStrategyIdleState)
	}
	clusterName := reader.ClusterName
	stateConf := createStateChangeConfig(ctx, api, reader.ProjectID, clusterName, targetState, reader.Timeout, extraPending...)
	clusterAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) && isDelete {
			return nil
		}
		diags.AddError(errorAwaitState, fmt.Sprintf("cluster=%s didn't reach desired state: %s, last operation: %s, error: %s", clusterName, targetState, lastOperation, err))
		return nil
	}
	if isDelete {
		return nil
	}
	cluster, ok := clusterAny.(*admin.ClusterDescription20240805)
	if !ok {
		diags.AddError(errorAwaitStateResultType, fmt.Sprintf("cluster=%s, got unexpected type: %T, last operation: %s", clusterName, clusterAny, lastOperation))
		return nil
	}
	return cluster
}

func createStateChangeConfig(ctx context.Context, api admin.ClustersApi, projectID, name, targetState string, timeout time.Duration, extraPending ...string) retry.StateChangeConf {
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
