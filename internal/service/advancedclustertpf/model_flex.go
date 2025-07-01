package advancedclustertpf

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312004/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

const defaultPriority int = 7

func NewFlexCreateReq(clusterName string, terminationProtectionEnabled bool, tags *[]admin.ResourceTag, replicationSpecs *[]admin.ReplicationSpec20240805) *admin.FlexClusterDescriptionCreate20241113 {
	if replicationSpecs == nil || len(*replicationSpecs) == 0 {
		return nil
	}
	regionConfigs := getRegionConfig(replicationSpecs)
	return &admin.FlexClusterDescriptionCreate20241113{
		Name: clusterName,
		ProviderSettings: admin.FlexProviderSettingsCreate20241113{
			BackingProviderName: regionConfigs.GetBackingProviderName(),
			RegionName:          regionConfigs.GetRegionName(),
		},
		TerminationProtectionEnabled: conversion.Pointer(terminationProtectionEnabled),
		Tags:                         tags,
	}
}

func NewReplicationSpecsFromFlexDescription(input *admin.FlexClusterDescription20241113, priority *int) *[]admin.ReplicationSpec20240805 {
	if input == nil {
		return nil
	}
	return &[]admin.ReplicationSpec20240805{
		{
			RegionConfigs: &[]admin.CloudRegionConfig20240805{
				{
					BackingProviderName: input.ProviderSettings.BackingProviderName,
					RegionName:          input.ProviderSettings.RegionName,
					ProviderName:        input.ProviderSettings.ProviderName,
					Priority:            priority,
				},
			},
			ZoneName: conversion.StringPtr("Zone 1"),
		},
	}
}

func NewClusterConnectionStringsFromFlex(connectionStrings *admin.FlexConnectionStrings20241113) *admin.ClusterConnectionStrings {
	if connectionStrings == nil {
		return nil
	}
	return &admin.ClusterConnectionStrings{
		Standard:    connectionStrings.Standard,
		StandardSrv: connectionStrings.StandardSrv,
	}
}

func isValidUpgradeTenantToFlex(stateCluster, planCluster *admin.ClusterDescription20240805) bool {
	if stateCluster.ReplicationSpecs == nil || planCluster.ReplicationSpecs == nil {
		return false
	}
	oldRegion := getRegionConfig(stateCluster.ReplicationSpecs)
	newRegion := getRegionConfig(planCluster.ReplicationSpecs)

	return oldRegion != newRegion &&
		oldRegion.GetProviderName() == constant.TENANT &&
		newRegion.GetProviderName() == flexcluster.FlexClusterType
}

func isValidUpdateOfFlex(stateCluster, planCluster *admin.ClusterDescription20240805) bool {
	patchFlex, err := update.PatchPayload(stateCluster, planCluster)
	if err != nil {
		return false
	}
	if update.IsZeroValues(patchFlex) { // No updates
		return false
	}
	okUpdatesChanged := patchFlex.Tags != nil || patchFlex.TerminationProtectionEnabled != nil
	notOkUpdatesChanged := patchFlex.ClusterType != nil && patchFlex.ReplicationSpecs != nil
	return okUpdatesChanged && !notOkUpdatesChanged
}

func GetFlexClusterUpdateRequest(tags *[]admin.ResourceTag, terminationProtectionEnabled *bool) *admin.FlexClusterDescriptionUpdate20241113 {
	return &admin.FlexClusterDescriptionUpdate20241113{
		Tags:                         tags,
		TerminationProtectionEnabled: terminationProtectionEnabled,
	}
}

func FlexDescriptionToClusterDescription(flexCluster *admin.FlexClusterDescription20241113, priority *int) *admin.ClusterDescription20240805 {
	if flexCluster == nil {
		return nil
	}
	var backupEnabled *bool
	if flexCluster.BackupSettings != nil {
		backupEnabled = flexCluster.BackupSettings.Enabled
	}
	return &admin.ClusterDescription20240805{
		ClusterType:                  flexCluster.ClusterType,
		BackupEnabled:                backupEnabled,
		CreateDate:                   flexCluster.CreateDate,
		MongoDBVersion:               flexCluster.MongoDBVersion,
		ReplicationSpecs:             NewReplicationSpecsFromFlexDescription(flexCluster, priority),
		Name:                         flexCluster.Name,
		GroupId:                      flexCluster.GroupId,
		StateName:                    flexCluster.StateName,
		Tags:                         flexCluster.Tags,
		TerminationProtectionEnabled: flexCluster.TerminationProtectionEnabled,
		VersionReleaseSystem:         flexCluster.VersionReleaseSystem,
		ConnectionStrings:            NewClusterConnectionStringsFromFlex(flexCluster.ConnectionStrings),
	}
}

func NewTFModelFlexResource(ctx context.Context, diags *diag.Diagnostics, flexCluster *admin.FlexClusterDescription20241113, priority *int, modelIn *TFModel) *TFModel {
	modelOut := NewTFModelFlex(ctx, diags, flexCluster, priority)
	if modelOut != nil {
		overrideAttributesWithPrevStateValue(modelIn, modelOut)
		modelOut.Timeouts = modelIn.Timeouts
	}
	return modelOut
}

func NewTFModelFlex(ctx context.Context, diags *diag.Diagnostics, flexCluster *admin.FlexClusterDescription20241113, priority *int) *TFModel {
	if priority == nil {
		priority = conversion.Pointer(defaultPriority)
	}
	modelOut := NewTFModel(ctx, FlexDescriptionToClusterDescription(flexCluster, priority), diags, ExtraAPIInfo{UseNewShardingConfig: true})
	if diags.HasError() {
		return nil
	}
	modelOut.AdvancedConfiguration = types.ObjectNull(AdvancedConfigurationObjType.AttrTypes)
	return modelOut
}

func FlexUpgrade(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, req *admin.LegacyAtlasTenantClusterUpgradeRequest) *admin.FlexClusterDescription20241113 {
	if _, _, err := client.AtlasV2.ClustersApi.UpgradeSharedCluster(ctx, waitParams.ProjectID, req).Execute(); err != nil {
		diags.AddError(fmt.Sprintf(flexcluster.ErrorUpgradeFlex, req.Name), err.Error())
		return nil
	}

	flexClusterParams := &admin.GetFlexClusterApiParams{
		GroupId: waitParams.ProjectID,
		Name:    waitParams.ClusterName,
	}

	flexClusterResp, err := flexcluster.WaitStateTransition(ctx, flexClusterParams, client.AtlasV2.FlexClustersApi, []string{retrystrategy.RetryStrategyUpdatingState}, []string{retrystrategy.RetryStrategyIdleState}, true, &waitParams.Timeout)
	if err != nil {
		diags.AddError(fmt.Sprintf(flexcluster.ErrorUpgradeFlex, req.Name), err.Error())
		return nil
	}
	return flexClusterResp
}

func GetUpgradeToFlexClusterRequest(planReq *admin.ClusterDescription20240805) *admin.LegacyAtlasTenantClusterUpgradeRequest {
	regionConfig := getRegionConfig(planReq.ReplicationSpecs)

	return &admin.LegacyAtlasTenantClusterUpgradeRequest{
		Name: planReq.GetName(),
		ProviderSettings: &admin.ClusterProviderSettings{
			ProviderName:        flexcluster.FlexClusterType,
			BackingProviderName: regionConfig.BackingProviderName,
			InstanceSizeName:    conversion.StringPtr(flexcluster.FlexClusterType),
			RegionName:          regionConfig.RegionName,
		},
	}
}
