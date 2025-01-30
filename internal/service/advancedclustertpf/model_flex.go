package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func NewFlexCreateReq(clusterName string, terminationProtectionEnabled bool, tags *[]admin.ResourceTag, replicationSpecs *[]admin.ReplicationSpec20240805) *admin.FlexClusterDescriptionCreate20241113 {
	if replicationSpecs == nil || len(*replicationSpecs) == 0 {
		return nil
	}
	regionConfigs := (*replicationSpecs)[0].GetRegionConfigs()[0]
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

func NewTFModelFlex(ctx context.Context, diags *diag.Diagnostics, flexCluster *admin.FlexClusterDescription20241113) *TFModel {
	tags := NewTagsObjType(ctx, diags, flexCluster.Tags)
	replicationSpecs := NewReplicationSpecsObjType(ctx, NewReplicationSpecsFromFlexDescription(flexCluster), diags, &ExtraAPIInfo{UsingLegacySchema: false})
	connectionStrings := NewConnectionStringsObjType(ctx, NewClusterConnectionStringsFromFlex(flexCluster.ConnectionStrings), diags)
	if diags.HasError() {
		return nil
	}
	return &TFModel{
		ClusterType:                  types.StringPointerValue(flexCluster.ClusterType),
		BackupEnabled:                types.BoolPointerValue(flexCluster.BackupSettings.Enabled),
		ConnectionStrings:            connectionStrings,
		CreateDate:                   types.StringValue(conversion.SafeValue(conversion.TimePtrToStringPtr(flexCluster.CreateDate))),
		MongoDBVersion:               types.StringValue(conversion.SafeValue(flexCluster.MongoDBVersion)),
		ReplicationSpecs:             replicationSpecs,
		Name:                         types.StringValue(conversion.SafeValue(flexCluster.Name)),
		ProjectID:                    types.StringValue(conversion.SafeValue(flexCluster.GroupId)),
		StateName:                    types.StringValue(conversion.SafeValue(flexCluster.StateName)),
		Tags:                         tags,
		TerminationProtectionEnabled: types.BoolPointerValue(flexCluster.TerminationProtectionEnabled),
		VersionReleaseSystem:         types.StringValue(conversion.SafeValue(flexCluster.VersionReleaseSystem)),
	}
}

func NewReplicationSpecsFromFlexDescription(input *admin.FlexClusterDescription20241113) *[]admin.ReplicationSpec20240805 {
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
				},
			},
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

// TODO: TFMOdel
func isValidUpgradeToFlex(state, plan TFModel) bool {
	// if d.HasChange("replication_specs.0.region_configs.0") {
	// 	oldProviderName, newProviderName := d.GetChange("replication_specs.0.region_configs.0.provider_name")
	// 	oldInstanceSize, newInstanceSize := d.GetChange("replication_specs.0.region_configs.0.electable_specs.instance_size")
	// 	if oldProviderName == constant.TENANT && newProviderName == flexcluster.FlexClusterType && oldInstanceSize != nil && newInstanceSize == nil {
	// 		return true
	// 	}
	// }
	return false
}

func isValidUpdateOfFlex(state, plan TFModel) bool {
	// updatableAttrHaveBeenUpdated := d.HasChange("tags") || d.HasChange("termination_protection_enabled")
	// nonUpdatableAttrHaveNotBeenUpdated := !d.HasChange("cluster_type") && !d.HasChange("replication_specs") && !d.HasChange("project_id") && !d.HasChange("name")
	// if updatableAttrHaveBeenUpdated && nonUpdatableAttrHaveNotBeenUpdated {
	// 	return true
	// }
	return false
}
