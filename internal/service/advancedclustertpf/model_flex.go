package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

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
			ZoneName: conversion.StringPtr("ZoneName managed by Terraform"),
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

func isValidUpgradeToFlex(stateCluster, planCluster *admin.ClusterDescription20240805) bool {
	if stateCluster.ReplicationSpecs == nil || planCluster.ReplicationSpecs == nil {
		return false
	}
	oldRegion := getRegionConfig(stateCluster.ReplicationSpecs)
	newRegion := getRegionConfig(planCluster.ReplicationSpecs)
	if oldRegion.ElectableSpecs == nil || newRegion.ElectableSpecs == nil {
		return false
	}
	return oldRegion != newRegion &&
		oldRegion.GetProviderName() == constant.TENANT &&
		newRegion.GetProviderName() == flexcluster.FlexClusterType &&
		oldRegion.ElectableSpecs.InstanceSize != nil &&
		newRegion.ElectableSpecs.InstanceSize == nil
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
	return &admin.ClusterDescription20240805{
		ClusterType:                  flexCluster.ClusterType,
		BackupEnabled:                flexCluster.BackupSettings.Enabled,
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

func NewTFModelFlex(ctx context.Context, diags *diag.Diagnostics, flexCluster *admin.FlexClusterDescription20241113, priority *int, modelIn *TFModel) *TFModel {
	modelOut := NewTFModel(ctx, FlexDescriptionToClusterDescription(flexCluster, priority), diags, ExtraAPIInfo{UseNewShardingConfig: true})
	if diags.HasError() {
		return nil
	}
	modelOut.AdvancedConfiguration = types.ObjectNull(AdvancedConfigurationObjType.AttrTypes)
	overrideAttributesWithPrevStateValue(modelIn, modelOut)
	modelOut.Timeouts = modelIn.Timeouts
	return modelOut
}

func FlexUpgrade(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, waitParams *ClusterWaitParams, req *admin.LegacyAtlasTenantClusterUpgradeRequest) *admin.FlexClusterDescription20241113 {
	//TODO: CLOUDP-296220
	return nil
}

func GetUpgradeToFlexClusterRequest() *admin.LegacyAtlasTenantClusterUpgradeRequest {
	// WIP: will be finished as part of CLOUDP-296220
	return &admin.LegacyAtlasTenantClusterUpgradeRequest{
		ProviderSettings: &admin.ClusterProviderSettings{
			ProviderName: flexcluster.FlexClusterType,
		},
	}
}
