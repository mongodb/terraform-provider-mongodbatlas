package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241023002/admin"
)

func NewAtlasReq(ctx context.Context, input *TFModel, diags *diag.Diagnostics) *admin.ClusterDescription20240805 {
	acceptDataRisksAndForceReplicaSetReconfig, ok := conversion.StringPtrToTimePtr(input.AcceptDataRisksAndForceReplicaSetReconfig.ValueStringPointer())
	if !ok {
		diags.AddError("error converting AcceptDataRisksAndForceReplicaSetReconfig", fmt.Sprintf("not a valid time: %s", input.AcceptDataRisksAndForceReplicaSetReconfig.ValueString()))
	}
	return &admin.ClusterDescription20240805{
		AcceptDataRisksAndForceReplicaSetReconfig: acceptDataRisksAndForceReplicaSetReconfig,
		BackupEnabled:                    input.BackupEnabled.ValueBoolPointer(),
		BiConnector:                      newBiConnector(ctx, input.BiConnectorConfig, diags),
		ClusterType:                      input.ClusterType.ValueStringPointer(),
		ConfigServerManagementMode:       input.ConfigServerManagementMode.ValueStringPointer(),
		EncryptionAtRestProvider:         input.EncryptionAtRestProvider.ValueStringPointer(),
		GlobalClusterSelfManagedSharding: input.GlobalClusterSelfManagedSharding.ValueBoolPointer(),
		GroupId:                          input.ProjectID.ValueStringPointer(),
		Labels:                           newComponentLabel(ctx, input.Labels, diags),
		MongoDBMajorVersion:              input.MongoDBMajorVersion.ValueStringPointer(),
		Name:                             input.Name.ValueStringPointer(),
		Paused:                           input.Paused.ValueBoolPointer(),
		PitEnabled:                       input.PitEnabled.ValueBoolPointer(),
		RedactClientLogData:              input.RedactClientLogData.ValueBoolPointer(),
		ReplicaSetScalingStrategy:        input.ReplicaSetScalingStrategy.ValueStringPointer(),
		ReplicationSpecs:                 newReplicationSpec20240805(ctx, input.ReplicationSpecs, diags),
		RootCertType:                     input.RootCertType.ValueStringPointer(),
		Tags:                             newResourceTag(ctx, input.Tags, diags),
		TerminationProtectionEnabled:     input.TerminationProtectionEnabled.ValueBoolPointer(),
		VersionReleaseSystem:             input.VersionReleaseSystem.ValueStringPointer(),
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
		Enabled:        item.Enabled.ValueBoolPointer(),
		ReadPreference: item.ReadPreference.ValueStringPointer(),
	}
}
func newComponentLabel(ctx context.Context, input types.Set, diags *diag.Diagnostics) *[]admin.ComponentLabel {
	if input.IsUnknown() || input.IsNull() {
		return nil
	}
	elements := make([]TFLabelsModel, 0, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	resp := make([]admin.ComponentLabel, 0, len(input.Elements()))
	for i := range elements {
		item := &elements[i]
		resp[i] = admin.ComponentLabel{
			Key:   item.Key.ValueStringPointer(),
			Value: item.Value.ValueStringPointer(),
		}
	}
	return &resp
}
func newReplicationSpec20240805(ctx context.Context, input types.List, diags *diag.Diagnostics) *[]admin.ReplicationSpec20240805 {
	if input.IsUnknown() || input.IsNull() {
		return nil
	}
	elements := make([]TFReplicationSpecsModel, 0, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	resp := make([]admin.ReplicationSpec20240805, 0, len(input.Elements()))
	for i := range elements {
		item := &elements[i]
		resp[i] = admin.ReplicationSpec20240805{
			RegionConfigs: newCloudRegionConfig20240805(ctx, item.RegionConfigs, diags),
			ZoneName:      item.ZoneName.ValueStringPointer(),
		}
	}
	return &resp
}
func newResourceTag(ctx context.Context, input types.Set, diags *diag.Diagnostics) *[]admin.ResourceTag {
	if input.IsUnknown() || input.IsNull() {
		return nil
	}
	elements := make([]TFTagsModel, 0, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	resp := make([]admin.ResourceTag, 0, len(input.Elements()))
	for i := range elements {
		item := &elements[i]
		resp[i] = admin.ResourceTag{
			Key:   item.Key.ValueString(),
			Value: item.Value.ValueString(),
		}
	}
	return &resp
}
func newCloudRegionConfig20240805(ctx context.Context, input types.List, diags *diag.Diagnostics) *[]admin.CloudRegionConfig20240805 {
	if input.IsUnknown() || input.IsNull() {
		return nil
	}
	elements := make([]TFRegionConfigsModel, 0, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	resp := make([]admin.CloudRegionConfig20240805, 0, len(input.Elements()))
	for i := range elements {
		item := &elements[i]
		resp[i] = admin.CloudRegionConfig20240805{
			AnalyticsAutoScaling: newAdvancedAutoScalingSettings(ctx, item.AnalyticsAutoScaling, diags),
			AnalyticsSpecs:       newDedicatedHardwareSpec20240805(ctx, item.AnalyticsSpecs, diags),
			AutoScaling:          newAdvancedAutoScalingSettings(ctx, item.AutoScaling, diags),
			BackingProviderName:  item.BackingProviderName.ValueStringPointer(),
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
		DiskIOPS:      conversion.Int64PtrToIntPtr(item.DiskIops.ValueInt64Pointer()),
		DiskSizeGB:    item.DiskSizeGb.ValueFloat64Pointer(),
		EbsVolumeType: item.EbsVolumeType.ValueStringPointer(),
		InstanceSize:  item.InstanceSize.ValueStringPointer(),
		NodeCount:     conversion.Int64PtrToIntPtr(item.NodeCount.ValueInt64Pointer()),
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
		DiskIOPS:      conversion.Int64PtrToIntPtr(item.DiskIops.ValueInt64Pointer()),
		DiskSizeGB:    item.DiskSizeGb.ValueFloat64Pointer(),
		EbsVolumeType: item.EbsVolumeType.ValueStringPointer(),
		InstanceSize:  item.InstanceSize.ValueStringPointer(),
		NodeCount:     conversion.Int64PtrToIntPtr(item.NodeCount.ValueInt64Pointer()),
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
		Enabled:          item.ComputeEnabled.ValueBoolPointer(),
		MaxInstanceSize:  item.ComputeMaxInstanceSize.ValueStringPointer(),
		MinInstanceSize:  item.ComputeMinInstanceSize.ValueStringPointer(),
		ScaleDownEnabled: item.ComputeScaleDownEnabled.ValueBoolPointer(),
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
		Enabled: item.DiskGBEnabled.ValueBoolPointer(),
	}
}
