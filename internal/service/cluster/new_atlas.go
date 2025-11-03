package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

func newAtlasUpdate(ctx context.Context, timeout time.Duration, connV2 *admin.APIClient, projectID, clusterName string, redactClientLogData bool) error {
	current, err := newAtlasGet(ctx, connV2, projectID, clusterName)
	if err != nil {
		return err
	}
	if current.GetRedactClientLogData() == redactClientLogData {
		return nil
	}
	req := &admin.ClusterDescription20240805{
		RedactClientLogData: &redactClientLogData,
	}
	// can call latest API (2024-10-23 or newer) as replications specs (with nested autoscaling property) is not specified
	if _, _, err = connV2.ClustersApi.UpdateCluster(ctx, projectID, clusterName, req).Execute(); err != nil {
		return err
	}
	stateConf := CreateStateChangeConfig(ctx, connV2, projectID, clusterName, timeout)
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func newAtlasGet(ctx context.Context, connV2 *admin.APIClient, projectID, clusterName string) (*admin.ClusterDescription20240805, error) {
	cluster, _, err := connV2.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	return cluster, err
}

func newAtlasList(ctx context.Context, connV2 *admin.APIClient, projectID string) (map[string]*admin.ClusterDescription20240805, error) {
	clusters, _, err := connV2.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		return nil, err
	}
	results := clusters.GetResults()
	list := make(map[string]*admin.ClusterDescription20240805)
	for i := range results {
		list[results[i].GetName()] = &results[i]
	}
	return list, nil
}

func CreateStateChangeConfig(ctx context.Context, connV2 *admin.APIClient, projectID, name string, timeout time.Duration) retry.StateChangeConf {
	return retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING", "PENDING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceRefreshFunc(ctx, name, projectID, connV2),
		Timeout:    timeout,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}
}

func DeleteStateChangeConfig(ctx context.Context, connV2 *admin.APIClient, projectID, name string, timeout time.Duration) retry.StateChangeConf {
	return retry.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING", "PENDING", "REPEATING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceRefreshFunc(ctx, name, projectID, connV2),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}
}

func resourceRefreshFunc(ctx context.Context, name, projectID string, connV2 *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		cluster, resp, err := connV2.ClustersApi.GetCluster(ctx, projectID, name).Execute()
		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && cluster == nil && resp == nil {
			return nil, "", err
		}

		if err != nil {
			if validate.StatusNotFound(resp) {
				return "", "DELETED", nil
			}
			if validate.StatusServiceUnavailable(resp) {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		state := cluster.GetStateName()
		return cluster, state, nil
	}
}

func handlePinnedFCVUpdate(ctx context.Context, connV2 *admin.APIClient, projectID, clusterName string, d *schema.ResourceData, timeout time.Duration) diag.Diagnostics {
	if d.HasChange("pinned_fcv") {
		pinnedFCVBlock, _ := d.Get("pinned_fcv").([]any)
		isFCVPresentInConfig := len(pinnedFCVBlock) > 0
		if isFCVPresentInConfig {
			// pinned_fcv has been defined or updated expiration date
			nestedObj := pinnedFCVBlock[0].(map[string]any)
			expDateStr := cast.ToString(nestedObj["expiration_date"])
			if err := advancedcluster.PinFCV(ctx, connV2.ClustersApi, projectID, clusterName, expDateStr); err != nil {
				return diag.FromErr(fmt.Errorf(errorClusterUpdate, clusterName, err))
			}
		} else {
			// pinned_fcv has been removed from the config so unpin method is called
			if _, err := connV2.ClustersApi.UnpinFeatureCompatibilityVersion(ctx, projectID, clusterName).Execute(); err != nil {
				return diag.FromErr(fmt.Errorf(errorClusterUpdate, clusterName, err))
			}
		}
		// ensures cluster is in IDLE state before continuing with other changes
		if err := waitForUpdateToFinish(ctx, connV2, projectID, clusterName, timeout); err != nil {
			return diag.FromErr(fmt.Errorf(errorClusterUpdate, clusterName, err))
		}
	}
	return nil
}

func waitForUpdateToFinish(ctx context.Context, connV2 *admin.APIClient, projectID, clusterName string, timeout time.Duration) error {
	stateConf := CreateStateChangeConfig(ctx, connV2, projectID, clusterName, timeout)
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
