package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241113005/admin"
)

const defaultZoneName = "ZoneName managed by Terraform"

func NewAtlasReq(ctx context.Context, input *TFModel, diags *diag.Diagnostics) *admin.ClusterDescription20240805 {
	acceptDataRisksAndForceReplicaSetReconfig, ok := conversion.StringPtrToTimePtr(input.AcceptDataRisksAndForceReplicaSetReconfig.ValueStringPointer())
	if !ok {
		diags.AddError("error converting AcceptDataRisksAndForceReplicaSetReconfig", fmt.Sprintf("not a valid time: %s", input.AcceptDataRisksAndForceReplicaSetReconfig.ValueString()))
	}
	majorVersion := conversion.NilForUnknown(input.MongoDBMajorVersion, input.MongoDBMajorVersion.ValueStringPointer())
	if majorVersion != nil {
		majorVersionFormatted := FormatMongoDBMajorVersion(*majorVersion)
		majorVersion = &majorVersionFormatted
	}
	return &admin.ClusterDescription20240805{
		AcceptDataRisksAndForceReplicaSetReconfig: acceptDataRisksAndForceReplicaSetReconfig,
		BackupEnabled:                    conversion.NilForUnknown(input.BackupEnabled, input.BackupEnabled.ValueBoolPointer()),
		BiConnector:                      newBiConnector(ctx, input.BiConnectorConfig, diags),
		ClusterType:                      input.ClusterType.ValueStringPointer(),
		ConfigServerManagementMode:       conversion.NilForUnknown(input.ConfigServerManagementMode, input.ConfigServerManagementMode.ValueStringPointer()),
		EncryptionAtRestProvider:         conversion.NilForUnknown(input.EncryptionAtRestProvider, input.EncryptionAtRestProvider.ValueStringPointer()),
		GlobalClusterSelfManagedSharding: conversion.NilForUnknown(input.GlobalClusterSelfManagedSharding, input.GlobalClusterSelfManagedSharding.ValueBoolPointer()),
		Labels:                           newComponentLabel(ctx, diags, input.Labels),
		MongoDBMajorVersion:              majorVersion,
		Name:                             input.Name.ValueStringPointer(),
		Paused:                           conversion.NilForUnknown(input.Paused, input.Paused.ValueBoolPointer()),
		PitEnabled:                       conversion.NilForUnknown(input.PitEnabled, input.PitEnabled.ValueBoolPointer()),
		RedactClientLogData:              conversion.NilForUnknown(input.RedactClientLogData, input.RedactClientLogData.ValueBoolPointer()),
		ReplicaSetScalingStrategy:        conversion.NilForUnknown(input.ReplicaSetScalingStrategy, input.ReplicaSetScalingStrategy.ValueStringPointer()),
		ReplicationSpecs:                 newReplicationSpec20240805(ctx, input.ReplicationSpecs, diags),
		RootCertType:                     conversion.NilForUnknown(input.RootCertType, input.RootCertType.ValueStringPointer()),
		Tags:                             newResourceTag(ctx, diags, input.Tags),
		TerminationProtectionEnabled:     conversion.NilForUnknown(input.TerminationProtectionEnabled, input.TerminationProtectionEnabled.ValueBoolPointer()),
		VersionReleaseSystem:             conversion.NilForUnknown(input.VersionReleaseSystem, input.VersionReleaseSystem.ValueStringPointer()),
	}
}
func newBiConnector(ctx context.Context, input types.Object, diags *diag.Diagnostics) *admin.BiConnector {
	var resp *admin.BiConnector
	if input.IsUnknown() || input.IsNull() {
		return resp
	}
	item := &TFBiConnectorModel{}
	if localDiags := input.As(ctx, item, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return resp
	}
	return &admin.BiConnector{
		Enabled:        conversion.NilForUnknown(item.Enabled, item.Enabled.ValueBoolPointer()),
		ReadPreference: conversion.NilForUnknown(item.ReadPreference, item.ReadPreference.ValueStringPointer()),
	}
}

func newComponentLabel(ctx context.Context, diags *diag.Diagnostics, input types.Map) *[]admin.ComponentLabel {
	elms := make(map[string]types.String, len(input.Elements()))
	localDiags := input.ElementsAs(ctx, &elms, false)
	diags.Append(localDiags...)
	if diags.HasError() {
		return nil
	}
	ret := make([]admin.ComponentLabel, 0, len(input.Elements()))
	for key, value := range elms {
		if key == LegacyIgnoredLabelKey {
			diags.AddError(ErrLegacyIgnoreLabel.Error(), ErrLegacyIgnoreLabel.Error())
			return nil
		}
		ret = append(ret, admin.ComponentLabel{
			Key:   &key,
			Value: value.ValueStringPointer(),
		})
	}
	return &ret
}

func newReplicationSpec20240805(ctx context.Context, input types.List, diags *diag.Diagnostics) *[]admin.ReplicationSpec20240805 {
	if input.IsUnknown() || input.IsNull() {
		return nil
	}
	elements := make([]TFReplicationSpecsModel, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	resp := make([]admin.ReplicationSpec20240805, len(input.Elements()))
	for i := range elements {
		item := &elements[i]
		resp[i] = admin.ReplicationSpec20240805{
			Id:            conversion.NilForUnknown(item.ExternalId, item.ExternalId.ValueStringPointer()),
			ZoneId:        conversion.NilForUnknown(item.ZoneId, item.ZoneId.ValueStringPointer()),
			RegionConfigs: newCloudRegionConfig20240805(ctx, item.RegionConfigs, diags),
			ZoneName:      conversion.StringPtr(resolveZoneNameOrUseDefault(item)),
		}
	}
	return &resp
}

func resolveZoneNameOrUseDefault(item *TFReplicationSpecsModel) string {
	zoneName := conversion.NilForUnknown(item.ZoneName, item.ZoneName.ValueStringPointer())
	if zoneName == nil {
		return defaultZoneName
	}
	return *zoneName
}

func newResourceTag(ctx context.Context, diags *diag.Diagnostics, input types.Map) *[]admin.ResourceTag {
	elms := make(map[string]types.String, len(input.Elements()))
	localDiags := input.ElementsAs(ctx, &elms, false)
	diags.Append(localDiags...)
	if diags.HasError() {
		return nil
	}
	ret := make([]admin.ResourceTag, 0, len(input.Elements()))
	for key, value := range elms {
		ret = append(ret, admin.ResourceTag{
			Key:   key,
			Value: value.ValueString(),
		})
	}
	return &ret
}

func newCloudRegionConfig20240805(ctx context.Context, input types.List, diags *diag.Diagnostics) *[]admin.CloudRegionConfig20240805 {
	if input.IsUnknown() || input.IsNull() {
		return nil
	}
	elements := make([]TFRegionConfigsModel, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	resp := make([]admin.CloudRegionConfig20240805, len(input.Elements()))
	for i := range elements {
		item := &elements[i]
		resp[i] = admin.CloudRegionConfig20240805{
			AnalyticsAutoScaling: newAdvancedAutoScalingSettings(ctx, item.AnalyticsAutoScaling, diags),
			AnalyticsSpecs:       newDedicatedHardwareSpec20240805(ctx, item.AnalyticsSpecs, diags),
			AutoScaling:          newAdvancedAutoScalingSettings(ctx, item.AutoScaling, diags),
			BackingProviderName:  conversion.NilForUnknown(item.BackingProviderName, item.BackingProviderName.ValueStringPointer()),
			ElectableSpecs:       newHardwareSpec20240805(ctx, item.ElectableSpecs, diags),
			Priority:             conversion.Int64PtrToIntPtr(item.Priority.ValueInt64Pointer()),
			ProviderName:         item.ProviderName.ValueStringPointer(),
			ReadOnlySpecs:        newDedicatedHardwareSpec20240805(ctx, item.ReadOnlySpecs, diags),
			RegionName:           item.RegionName.ValueStringPointer(),
		}
	}
	return &resp
}

func newAdvancedAutoScalingSettings(ctx context.Context, input types.Object, diags *diag.Diagnostics) *admin.AdvancedAutoScalingSettings {
	var resp *admin.AdvancedAutoScalingSettings
	if input.IsUnknown() || input.IsNull() {
		return resp
	}
	item := &TFAutoScalingModel{}
	if localDiags := input.As(ctx, item, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return resp
	}
	return &admin.AdvancedAutoScalingSettings{
		Compute: newAdvancedComputeAutoScaling(ctx, input, diags),
		DiskGB:  newDiskGBAutoScaling(ctx, input, diags),
	}
}
func newHardwareSpec20240805(ctx context.Context, input types.Object, diags *diag.Diagnostics) *admin.HardwareSpec20240805 {
	var resp *admin.HardwareSpec20240805
	if input.IsUnknown() || input.IsNull() {
		return resp
	}
	item := &TFSpecsModel{}
	if localDiags := input.As(ctx, item, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return resp
	}
	return &admin.HardwareSpec20240805{
		DiskIOPS:      conversion.NilForUnknown(item.DiskIops, conversion.Int64PtrToIntPtr(item.DiskIops.ValueInt64Pointer())),
		DiskSizeGB:    conversion.NilForUnknown(item.DiskSizeGb, item.DiskSizeGb.ValueFloat64Pointer()),
		EbsVolumeType: conversion.NilForUnknownOrEmptyString(item.EbsVolumeType),
		InstanceSize:  conversion.NilForUnknown(item.InstanceSize, item.InstanceSize.ValueStringPointer()),
		NodeCount:     conversion.NilForUnknown(item.NodeCount, conversion.Int64PtrToIntPtr(item.NodeCount.ValueInt64Pointer())),
	}
}
func newDedicatedHardwareSpec20240805(ctx context.Context, input types.Object, diags *diag.Diagnostics) *admin.DedicatedHardwareSpec20240805 {
	var resp *admin.DedicatedHardwareSpec20240805
	if input.IsUnknown() || input.IsNull() {
		return resp
	}
	item := &TFSpecsModel{}
	if localDiags := input.As(ctx, item, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return resp
	}
	return &admin.DedicatedHardwareSpec20240805{
		DiskIOPS:      conversion.NilForUnknown(item.DiskIops, conversion.Int64PtrToIntPtr(item.DiskIops.ValueInt64Pointer())),
		DiskSizeGB:    conversion.NilForUnknown(item.DiskSizeGb, item.DiskSizeGb.ValueFloat64Pointer()),
		EbsVolumeType: conversion.NilForUnknownOrEmptyString(item.EbsVolumeType),
		InstanceSize:  conversion.NilForUnknownOrEmptyString(item.InstanceSize),
		NodeCount:     conversion.NilForUnknown(item.NodeCount, conversion.Int64PtrToIntPtr(item.NodeCount.ValueInt64Pointer())),
	}
}

func newAdvancedComputeAutoScaling(ctx context.Context, input types.Object, diags *diag.Diagnostics) *admin.AdvancedComputeAutoScaling {
	var resp *admin.AdvancedComputeAutoScaling
	if input.IsUnknown() || input.IsNull() {
		return resp
	}
	item := &TFAutoScalingModel{}
	if localDiags := input.As(ctx, item, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return resp
	}
	return &admin.AdvancedComputeAutoScaling{
		Enabled:          conversion.NilForUnknown(item.ComputeEnabled, item.ComputeEnabled.ValueBoolPointer()),
		ScaleDownEnabled: conversion.NilForUnknown(item.ComputeScaleDownEnabled, item.ComputeScaleDownEnabled.ValueBoolPointer()),
		MaxInstanceSize:  conversion.NilForUnknownOrEmptyString(item.ComputeMaxInstanceSize),
		MinInstanceSize:  conversion.NilForUnknownOrEmptyString(item.ComputeMinInstanceSize),
	}
}
func newDiskGBAutoScaling(ctx context.Context, input types.Object, diags *diag.Diagnostics) *admin.DiskGBAutoScaling {
	var resp *admin.DiskGBAutoScaling
	if input.IsUnknown() || input.IsNull() {
		return resp
	}
	item := &TFAutoScalingModel{}
	if localDiags := input.As(ctx, item, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return resp
	}
	return &admin.DiskGBAutoScaling{
		Enabled: conversion.NilForUnknown(item.DiskGBEnabled, item.DiskGBEnabled.ValueBoolPointer()),
	}
}
