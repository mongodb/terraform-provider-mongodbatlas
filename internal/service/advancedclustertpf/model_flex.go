package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	if planCluster.ReplicationSpecs == nil {
		return false
	}
	if stateCluster.ReplicationSpecs == nil {
		return false
	}
	oldRegion := stateCluster.GetReplicationSpecs()[0].GetRegionConfigs()[0]
	oldProviderName := oldRegion.GetProviderName()
	oldInstanceSize := oldRegion.ElectableSpecs.InstanceSize
	newRegion := planCluster.GetReplicationSpecs()[0].GetRegionConfigs()[0]
	newProviderName := newRegion.GetProviderName()
	newInstanceSize := newRegion.ElectableSpecs.InstanceSize
	if oldRegion != newRegion {
		if oldProviderName == constant.TENANT && newProviderName == flexcluster.FlexClusterType && oldInstanceSize != nil && newInstanceSize == nil {
			return true
		}
	}
	return false
}

func isValidUpdateOfFlex(stateCluster, planCluster *admin.ClusterDescription20240805) bool {
	updatableAttrHaveBeenUpdated := stateCluster.Tags == planCluster.Tags || stateCluster.TerminationProtectionEnabled == planCluster.TerminationProtectionEnabled
	nonUpdatableAttrHaveNotBeenUpdated := stateCluster.ClusterType == planCluster.ClusterType && stateCluster.ReplicationSpecs == planCluster.ReplicationSpecs && stateCluster.GroupId == planCluster.GroupId && stateCluster.Name == planCluster.Name
	if updatableAttrHaveBeenUpdated && nonUpdatableAttrHaveNotBeenUpdated {
		return true
	}
	return false
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

func NewTFModelFlex(ctx context.Context, diags *diag.Diagnostics, flexCluster *admin.FlexClusterDescription20241113, priority *int, timeout timeouts.Value) *TFModel {
	model := NewTFModel(ctx, FlexDescriptionToClusterDescription(flexCluster, priority), timeout, diags, ExtraAPIInfo{UsingLegacySchema: false})
	AddAdvancedConfig(ctx, model, nil, nil, diags)
	return model
}
