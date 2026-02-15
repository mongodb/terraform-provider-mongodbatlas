package advancedcluster

import (
	"context"
	"fmt"
	"sort"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

const defaultZoneName = "ZoneName managed by Terraform"

func newAtlasReq(ctx context.Context, input *TFModel, diags *diag.Diagnostics) *admin.ClusterDescription20240805 {
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
		ReplicationSpecs:                 newReplicationSpec(ctx, input.ReplicationSpecs, diags),
		Tags:                             newResourceTag(ctx, diags, input.Tags),
		TerminationProtectionEnabled:     conversion.NilForUnknown(input.TerminationProtectionEnabled, input.TerminationProtectionEnabled.ValueBoolPointer()),
		VersionReleaseSystem:             conversion.NilForUnknown(input.VersionReleaseSystem, input.VersionReleaseSystem.ValueStringPointer()),
		AdvancedConfiguration:            newClusterAdvancedConfiguration(ctx, &input.AdvancedConfiguration, diags),
	}
}

func newClusterAdvancedConfiguration(ctx context.Context, objInput *types.Object, diags *diag.Diagnostics) *admin.ApiAtlasClusterAdvancedConfiguration {
	if objInput == nil || objInput.IsUnknown() || objInput.IsNull() {
		return nil
	}

	inputAdvConfig := &TFAdvancedConfigurationModel{}
	if localDiags := objInput.As(ctx, inputAdvConfig, basetypes.ObjectAsOptions{}); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}

	return &admin.ApiAtlasClusterAdvancedConfiguration{
		MinimumEnabledTlsProtocol:      conversion.NilForUnknown(inputAdvConfig.MinimumEnabledTlsProtocol, inputAdvConfig.MinimumEnabledTlsProtocol.ValueStringPointer()),
		TlsCipherConfigMode:            conversion.NilForUnknown(inputAdvConfig.TlsCipherConfigMode, inputAdvConfig.TlsCipherConfigMode.ValueStringPointer()),
		CustomOpensslCipherConfigTls12: stringSliceFromSet(ctx, inputAdvConfig.CustomOpensslCipherConfigTls12),
		CustomOpensslCipherConfigTls13: stringSliceFromSet(ctx, inputAdvConfig.CustomOpensslCipherConfigTls13),
	}
}

// stringSliceFromSet returns nil when the set is null or unknown (user didn't configure it), avoiding
// false diffs in PatchPayload. conversion.Pointer(TypesSetToString(...)) always produces &[]string{}
// for null sets, which differs from the API's nil or populated cipher list.
func stringSliceFromSet(ctx context.Context, set types.Set) *[]string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	result := conversion.TypesSetToString(ctx, set)
	return &result
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

// newComponentLabel converts a TF map to an SDK labels slice.
// Results are sorted by key to ensure deterministic ordering, since Go map iteration is random.
// Without sorting, PatchPayload's jsondiff would detect false changes from different orderings
// between state and plan, causing unnecessary cluster PATCH calls.
func newComponentLabel(ctx context.Context, diags *diag.Diagnostics, input types.Map) *[]admin.ComponentLabel {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}
	elms := make(map[string]types.String, len(input.Elements()))
	localDiags := input.ElementsAs(ctx, &elms, false)
	diags.Append(localDiags...)
	if diags.HasError() {
		return nil
	}
	keys := make([]string, 0, len(elms))
	for key := range elms {
		if key == LegacyIgnoredLabelKey {
			diags.AddError(ErrLegacyIgnoreLabel.Error(), ErrLegacyIgnoreLabel.Error())
			return nil
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	ret := make([]admin.ComponentLabel, 0, len(keys))
	for _, key := range keys {
		value := elms[key]
		ret = append(ret, admin.ComponentLabel{
			Key:   conversion.StringPtr(key),
			Value: value.ValueStringPointer(),
		})
	}
	return &ret
}

func newReplicationSpec(ctx context.Context, input types.List, diags *diag.Diagnostics) *[]admin.ReplicationSpec20240805 {
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
			RegionConfigs: newRegionConfig(ctx, item.RegionConfigs, diags),
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

// newResourceTag converts a TF map to an SDK tags slice.
// Results are sorted by key to ensure deterministic ordering, since Go map iteration is random.
// Without sorting, PatchPayload's jsondiff would detect false changes from different orderings
// between state and plan, causing unnecessary cluster PATCH calls.
func newResourceTag(ctx context.Context, diags *diag.Diagnostics, input types.Map) *[]admin.ResourceTag {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}
	elms := make(map[string]types.String, len(input.Elements()))
	localDiags := input.ElementsAs(ctx, &elms, false)
	diags.Append(localDiags...)
	if diags.HasError() {
		return nil
	}
	keys := make([]string, 0, len(elms))
	for key := range elms {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	ret := make([]admin.ResourceTag, 0, len(keys))
	for _, key := range keys {
		value := elms[key]
		ret = append(ret, admin.ResourceTag{
			Key:   key,
			Value: value.ValueString(),
		})
	}
	return &ret
}

func newRegionConfig(ctx context.Context, input types.List, diags *diag.Diagnostics) *[]admin.CloudRegionConfig20240805 {
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
			AnalyticsSpecs:       newDedicatedHardwareSpec(ctx, item.AnalyticsSpecs, diags),
			AutoScaling:          newAdvancedAutoScalingSettings(ctx, item.AutoScaling, diags),
			BackingProviderName:  conversion.NilForUnknown(item.BackingProviderName, item.BackingProviderName.ValueStringPointer()),
			ElectableSpecs:       newHardwareSpec(ctx, item.ElectableSpecs, diags),
			Priority:             conversion.Int64PtrToIntPtr(item.Priority.ValueInt64Pointer()),
			ProviderName:         item.ProviderName.ValueStringPointer(),
			ReadOnlySpecs:        newDedicatedHardwareSpec(ctx, item.ReadOnlySpecs, diags),
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
func newHardwareSpec(ctx context.Context, input types.Object, diags *diag.Diagnostics) *admin.HardwareSpec20240805 {
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
func newDedicatedHardwareSpec(ctx context.Context, input types.Object, diags *diag.Diagnostics) *admin.DedicatedHardwareSpec20240805 {
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
