package advancedcluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

const (
	errorUpdate                    = "error updating advanced cluster (%s): %s"
	ErrorClusterSetting            = "error setting `%s` for MongoDB Cluster (%s): %s"
	ErrorAdvancedConfRead          = "error reading Advanced Configuration Option %s for MongoDB Cluster (%s): %s"
	ErrorClusterAdvancedSetting    = "error setting `%s` for MongoDB ClusterAdvanced (%s): %s"
	ErrorFlexClusterSetting        = "error setting `%s` for MongoDB Flex Cluster (%s): %s"
	ErrorAdvancedClusterListStatus = "error awaiting MongoDB ClusterAdvanced List IDLE: %s"
	ErrorOperationNotPermitted     = "error operation not permitted"
	ErrorDefaultMaxTimeMinVersion  = "`advanced_configuration.default_max_time_ms` can only be configured if the mongo_db_major_version is 8.0 or higher"
	DeprecationOldSchemaAction     = "Please refer to our examples, documentation, and 1.18.0 migration guide for more details at https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide"
	V20240530                      = "(v20240530)"
)

var DeprecationMsgOldSchema = fmt.Sprintf("%s %s", constant.DeprecationParam, DeprecationOldSchemaAction)

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

// GetReplicationSpecAttributesFromOldAPI returns the id and num shard values of replication specs coming from old API. This is used to populate replication_specs.*.id and replication_specs.*.num_shard attributes for old sharding confirgurations.
// In the old API (2023-02-01), each replications spec has a 1:1 relation with each zone, so ids and num shards are stored in a struct oldShardConfigMeta and are returned in a map from zoneName to oldShardConfigMeta.
func GetReplicationSpecAttributesFromOldAPI(ctx context.Context, projectID, clusterName string, client20240530 admin20240530.ClustersApi) (map[string]OldShardConfigMeta, error) {
	clusterOldAPI, _, err := client20240530.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, err
	}
	specs := clusterOldAPI.GetReplicationSpecs()
	result := make(map[string]OldShardConfigMeta, len(specs))
	for _, spec := range specs {
		result[spec.GetZoneName()] = OldShardConfigMeta{spec.GetId(), spec.GetNumShards()}
	}
	return result, nil
}

func WarningIfFCVExpiredOrUnpinnedExternally(d *schema.ResourceData, cluster *admin.ClusterDescription20240805) diag.Diagnostics {
	pinnedFCVBlock, _ := d.Get("pinned_fcv").([]any)
	fcvPresentInState := len(pinnedFCVBlock) > 0
	diagsTpf := advancedclustertpf.GenerateFCVPinningWarningForRead(fcvPresentInState, cluster.FeatureCompatibilityVersionExpirationDate)
	return conversion.FromTPFDiagsToSDKV2Diags(diagsTpf)
}

func GetUpgradeToFlexClusterRequest(d *schema.ResourceData) *admin.LegacyAtlasTenantClusterUpgradeRequest {
	return &admin.LegacyAtlasTenantClusterUpgradeRequest{
		ProviderSettings: &admin.ClusterProviderSettings{
			ProviderName:        flexcluster.FlexClusterType,
			BackingProviderName: conversion.StringPtr(d.Get("replication_specs.0.region_configs.0.backing_provider_name").(string)),
			InstanceSizeName:    conversion.StringPtr(flexcluster.FlexClusterType),
			RegionName:          conversion.StringPtr(d.Get("replication_specs.0.region_configs.0.region_name").(string)),
		},
	}
}

func HandlePinnedFCVUpdate(ctx context.Context, connV2 *admin.APIClient, projectID, clusterName string, d *schema.ResourceData, timeout time.Duration) diag.Diagnostics {
	if d.HasChange("pinned_fcv") {
		pinnedFCVBlock, _ := d.Get("pinned_fcv").([]any)
		isFCVPresentInConfig := len(pinnedFCVBlock) > 0
		if isFCVPresentInConfig {
			// pinned_fcv has been defined or updated expiration date
			nestedObj := pinnedFCVBlock[0].(map[string]any)
			expDateStr := cast.ToString(nestedObj["expiration_date"])
			if err := advancedclustertpf.PinFCV(ctx, connV2.ClustersApi, projectID, clusterName, expDateStr); err != nil {
				return diag.FromErr(fmt.Errorf(errorUpdate, clusterName, err))
			}
		} else {
			// pinned_fcv has been removed from the config so unpin method is called
			if _, err := connV2.ClustersApi.UnpinFeatureCompatibilityVersion(ctx, projectID, clusterName).Execute(); err != nil {
				return diag.FromErr(fmt.Errorf(errorUpdate, clusterName, err))
			}
		}
		// ensures cluster is in IDLE state before continuing with other changes
		if err := waitForUpdateToFinish(ctx, connV2, projectID, clusterName, timeout); err != nil {
			return diag.FromErr(fmt.Errorf(errorUpdate, clusterName, err))
		}
	}
	return nil
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

func waitForUpdateToFinish(ctx context.Context, connV2 *admin.APIClient, projectID, name string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "PENDING", "REPEATING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceRefreshFunc(ctx, name, projectID, connV2),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
