package advancedclustertpf

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312003/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

var (
	RetryMinTimeout   = 1 * time.Minute
	RetryDelay        = 30 * time.Second
	RetryPollInterval = 30 * time.Second
)

type ClusterWaitParams struct {
	ProjectID   string
	ClusterName string
	Timeout     time.Duration
	IsDelete    bool
}

func AwaitChangesUpgrade(ctx context.Context, client *config.MongoDBClient, waitParams *ClusterWaitParams, errorLocator string, diags *diag.Diagnostics) *admin.ClusterDescription20240805 {
	upgraded := AwaitChanges(ctx, client, waitParams, errorLocator, diags)
	if diags.HasError() || upgraded == nil {
		return nil
	}
	providerName := getProviderName(upgraded.ReplicationSpecs)
	if slices.Contains([]string{flexcluster.FlexClusterType, constant.TENANT}, providerName) {
		tflog.Warn(ctx, fmt.Sprintf("cluster upgrade unexpected provider %s, retrying", providerName))
		return AwaitChanges(ctx, client, waitParams, errorLocator, diags)
	}
	return upgraded
}

func AwaitChanges(ctx context.Context, client *config.MongoDBClient, waitParams *ClusterWaitParams, errorLocator string, diags *diag.Diagnostics) *admin.ClusterDescription20240805 {
	api := client.AtlasV2.ClustersApi
	targetState := retrystrategy.RetryStrategyIdleState
	extraPending := []string{}
	isDelete := waitParams.IsDelete
	if isDelete {
		targetState = retrystrategy.RetryStrategyDeletedState
		extraPending = append(extraPending, retrystrategy.RetryStrategyIdleState)
	}
	clusterName := waitParams.ClusterName
	stateConf := createStateChangeConfig(ctx, api, waitParams.ProjectID, clusterName, targetState, waitParams.Timeout, extraPending...)
	clusterAny, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) && isDelete {
			return nil
		}
		addErrorDiag(diags, errorLocator, fmt.Sprintf("cluster=%s didn't reach desired state: %s, error: %s", clusterName, targetState, err))
		return nil
	}
	if isDelete {
		return nil
	}
	cluster, ok := clusterAny.(*admin.ClusterDescription20240805)
	if !ok {
		addErrorDiag(diags, errorLocator, fmt.Sprintf("cluster=%s, got unexpected type: %T", clusterName, clusterAny))
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
		Refresh:      ResourceRefreshFunc(ctx, name, projectID, api),
		Timeout:      timeout,
		MinTimeout:   RetryMinTimeout,
		Delay:        RetryDelay,
		PollInterval: RetryPollInterval,
	}
}

func ResourceRefreshFunc(ctx context.Context, name, projectID string, api admin.ClustersApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		cluster, resp, err := api.GetCluster(ctx, projectID, name).Execute()
		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, retrystrategy.RetryStrategyRepeatingState, nil
		}

		if err != nil && cluster == nil && resp == nil {
			return nil, "", err
		}

		if err != nil {
			if admin.IsErrorCode(err, "CANNOT_USE_FLEX_CLUSTER_IN_CLUSTER_API") {
				return nil, retrystrategy.RetryStrategyUpdatingState, nil
			}
			if validate.StatusNotFound(resp) {
				return "", retrystrategy.RetryStrategyDeletedState, nil
			}
			if validate.StatusServiceUnavailable(resp) {
				return "", retrystrategy.RetryStrategyPendingState, nil
			}
			return nil, "", err
		}

		state := cluster.GetStateName()
		return cluster, state, nil
	}
}
