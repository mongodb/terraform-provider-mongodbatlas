package advancedcluster

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

const (
	errorZoneNameNotSet = "zoneName is required for legacy schema"
)

func newTFModel(ctx context.Context, input *admin.ClusterDescription20240805, diags *diag.Diagnostics) *TFModel {
	biConnector := newBiConnectorConfigObjType(ctx, input.BiConnector, diags)
	connectionStrings := newConnectionStringsObjType(ctx, input.ConnectionStrings, diags)
	labels := newLabelsObjType(ctx, diags, input.Labels)
	replicationSpecs := newReplicationSpecsObjType(ctx, input.ReplicationSpecs, diags)
	tags := newTagsObjType(ctx, diags, input.Tags)
	pinnedFCV := newPinnedFCVObjType(ctx, input, diags)
	if diags.HasError() {
		return nil
	}
	return &TFModel{
		AcceptDataRisksAndForceReplicaSetReconfig: types.StringPointerValue(conversion.TimePtrToStringPtr(input.AcceptDataRisksAndForceReplicaSetReconfig)),
		AdvancedConfiguration:                     types.ObjectNull(advancedConfigurationObjType.AttrTypes),
		BackupEnabled:                             types.BoolPointerValue(input.BackupEnabled),
		BiConnectorConfig:                         biConnector,
		ClusterType:                               types.StringPointerValue(input.ClusterType),
		ConfigServerManagementMode:                types.StringPointerValue(input.ConfigServerManagementMode),
		ConfigServerType:                          types.StringPointerValue(input.ConfigServerType),
		ConnectionStrings:                         connectionStrings,
		CreateDate:                                types.StringPointerValue(conversion.TimePtrToStringPtr(input.CreateDate)),
		EncryptionAtRestProvider:                  types.StringPointerValue(input.EncryptionAtRestProvider),
		GlobalClusterSelfManagedSharding:          types.BoolPointerValue(input.GlobalClusterSelfManagedSharding),
		ProjectID:                                 types.StringPointerValue(input.GroupId),
		ClusterID:                                 types.StringPointerValue(input.Id),
		Labels:                                    labels,
		MongoDBMajorVersion:                       types.StringPointerValue(input.MongoDBMajorVersion),
		MongoDBVersion:                            types.StringPointerValue(input.MongoDBVersion),
		Name:                                      types.StringPointerValue(input.Name),
		Paused:                                    types.BoolPointerValue(input.Paused),
		PitEnabled:                                types.BoolPointerValue(input.PitEnabled),
		RedactClientLogData:                       types.BoolPointerValue(input.RedactClientLogData),
		ReplicaSetScalingStrategy:                 types.StringPointerValue(input.ReplicaSetScalingStrategy),
		ReplicationSpecs:                          replicationSpecs,
		RootCertType:                              types.StringPointerValue(input.RootCertType),
		StateName:                                 types.StringPointerValue(input.StateName),
		Tags:                                      tags,
		TerminationProtectionEnabled:              types.BoolPointerValue(input.TerminationProtectionEnabled),
		VersionReleaseSystem:                      types.StringPointerValue(input.VersionReleaseSystem),
		PinnedFCV:                                 pinnedFCV,
	}
}

func newTFModelDS(ctx context.Context, input *admin.ClusterDescription20240805, diags *diag.Diagnostics, containerIDs map[string]string) *TFModelDS {
	biConnector := newBiConnectorConfigObjType(ctx, input.BiConnector, diags)
	connectionStrings := newConnectionStringsObjType(ctx, input.ConnectionStrings, diags)
	labels := newLabelsObjType(ctx, diags, input.Labels)
	replicationSpecs := newReplicationSpecsDSObjType(ctx, input.ReplicationSpecs, diags)
	effectiveReplicationSpecs := newEffectiveReplicationSpecsObjType(ctx, input.EffectiveReplicationSpecs, diags, containerIDs)
	tags := newTagsObjType(ctx, diags, input.Tags)
	pinnedFCV := newPinnedFCVObjType(ctx, input, diags)
	if diags.HasError() {
		return nil
	}
	return &TFModelDS{
		BackupEnabled:                    types.BoolPointerValue(input.BackupEnabled),
		BiConnectorConfig:                biConnector,
		ClusterType:                      types.StringPointerValue(input.ClusterType),
		ConfigServerManagementMode:       types.StringPointerValue(input.ConfigServerManagementMode),
		ConfigServerType:                 types.StringPointerValue(input.ConfigServerType),
		ConnectionStrings:                connectionStrings,
		CreateDate:                       types.StringPointerValue(conversion.TimePtrToStringPtr(input.CreateDate)),
		EncryptionAtRestProvider:         types.StringPointerValue(input.EncryptionAtRestProvider),
		GlobalClusterSelfManagedSharding: types.BoolPointerValue(input.GlobalClusterSelfManagedSharding),
		ProjectID:                        types.StringPointerValue(input.GroupId),
		ClusterID:                        types.StringPointerValue(input.Id),
		Labels:                           labels,
		MongoDBMajorVersion:              types.StringPointerValue(input.MongoDBMajorVersion),
		MongoDBVersion:                   types.StringPointerValue(input.MongoDBVersion),
		Name:                             types.StringPointerValue(input.Name),
		Paused:                           types.BoolPointerValue(input.Paused),
		PitEnabled:                       types.BoolPointerValue(input.PitEnabled),
		RedactClientLogData:              types.BoolPointerValue(input.RedactClientLogData),
		ReplicaSetScalingStrategy:        types.StringPointerValue(input.ReplicaSetScalingStrategy),
		ReplicationSpecs:                 replicationSpecs,
		EffectiveReplicationSpecs:        effectiveReplicationSpecs,
		RootCertType:                     types.StringPointerValue(input.RootCertType),
		StateName:                        types.StringPointerValue(input.StateName),
		Tags:                             tags,
		TerminationProtectionEnabled:     types.BoolPointerValue(input.TerminationProtectionEnabled),
		UseAwsTimeBasedSnapshotCopyForFastInitialSync: types.BoolPointerValue(input.UseAwsTimeBasedSnapshotCopyForFastInitialSync),
		VersionReleaseSystem:                          types.StringPointerValue(input.VersionReleaseSystem),
		PinnedFCV:                                     pinnedFCV,
	}
}

func newBiConnectorConfigObjType(ctx context.Context, input *admin.BiConnector, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(biConnectorConfigObjType.AttrTypes)
	}
	tfModel := TFBiConnectorModel{
		Enabled:        types.BoolPointerValue(input.Enabled),
		ReadPreference: types.StringPointerValue(input.ReadPreference),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, biConnectorConfigObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

func newConnectionStringsObjType(ctx context.Context, input *admin.ClusterConnectionStrings, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(connectionStringsObjType.AttrTypes)
	}
	privateEndpoint := newPrivateEndpointObjType(ctx, input.PrivateEndpoint, diags)
	tfModel := TFConnectionStringsModel{
		Private:         types.StringPointerValue(input.Private),
		PrivateEndpoint: privateEndpoint,
		PrivateSrv:      types.StringPointerValue(input.PrivateSrv),
		Standard:        types.StringPointerValue(input.Standard),
		StandardSrv:     types.StringPointerValue(input.StandardSrv),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, connectionStringsObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

func newLabelsObjType(ctx context.Context, diags *diag.Diagnostics, input *[]admin.ComponentLabel) types.Map {
	elms := make(map[string]string)
	if input != nil {
		for _, item := range *input {
			key := item.GetKey()
			value := item.GetValue()
			if key == LegacyIgnoredLabelKey {
				continue
			}
			elms[key] = value
		}
	}
	return conversion.ToTFMapOfString(ctx, diags, elms)
}

func newReplicationSpecsObjType(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics) types.List {
	if input == nil {
		return types.ListNull(replicationSpecsObjType)
	}
	tfModels := convertReplicationSpecs(ctx, input, diags, newRegionConfigsObjType)
	if diags.HasError() {
		return types.ListNull(replicationSpecsObjType)
	}
	listType, diagsLocal := types.ListValueFrom(ctx, replicationSpecsObjType, *tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func newReplicationSpecsDSObjType(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics) types.List {
	if input == nil {
		return types.ListNull(replicationSpecsDSObjType)
	}
	tfModels := convertReplicationSpecs(ctx, input, diags, newRegionConfigsDSObjType)
	if diags.HasError() {
		return types.ListNull(replicationSpecsDSObjType)
	}
	listType, diagsLocal := types.ListValueFrom(ctx, replicationSpecsDSObjType, *tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func newEffectiveReplicationSpecsObjType(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, containerIDs map[string]string) types.List {
	if input == nil {
		return types.ListNull(effectiveReplicationSpecsObjType)
	}
	tfModels := convertEffectiveReplicationSpecs(ctx, input, diags, containerIDs, newRegionConfigsDSObjType)
	if diags.HasError() {
		return types.ListNull(effectiveReplicationSpecsObjType)
	}
	listType, diagsLocal := types.ListValueFrom(ctx, effectiveReplicationSpecsObjType, *tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func newPinnedFCVObjType(ctx context.Context, cluster *admin.ClusterDescription20240805, diags *diag.Diagnostics) types.Object {
	if cluster.FeatureCompatibilityVersionExpirationDate == nil {
		return types.ObjectNull(pinnedFCVObjType.AttrTypes)
	}
	tfModel := TFPinnedFCVModel{
		Version:        types.StringValue(cluster.GetFeatureCompatibilityVersion()),
		ExpirationDate: types.StringValue(conversion.TimeToString(cluster.GetFeatureCompatibilityVersionExpirationDate())),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, pinnedFCVObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

// regionConfigsConverter is a function type for converting region configs
type regionConfigsConverter func(context.Context, *[]admin.CloudRegionConfig20240805, *diag.Diagnostics) types.List

func convertReplicationSpecs(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, regionConfigsConv regionConfigsConverter) *[]TFReplicationSpecsModel {
	tfModels := make([]TFReplicationSpecsModel, len(*input))
	for i, item := range *input {
		regionConfigs := regionConfigsConv(ctx, item.RegionConfigs, diags)
		zoneName := item.GetZoneName()
		if zoneName == "" {
			diags.AddError(errorZoneNameNotSet, errorZoneNameNotSet)
			return &tfModels
		}
		tfModels[i] = TFReplicationSpecsModel{
			RegionConfigs: regionConfigs,
			ZoneName:      types.StringPointerValue(item.ZoneName),
		}
	}
	return &tfModels
}

func convertEffectiveReplicationSpecs(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, containerIDs map[string]string, regionConfigsConv regionConfigsConverter) *[]TFEffectiveReplicationSpecsModel {
	tfModels := make([]TFEffectiveReplicationSpecsModel, len(*input))
	for i, item := range *input {
		regionConfigs := regionConfigsConv(ctx, item.RegionConfigs, diags)
		zoneName := item.GetZoneName()
		if zoneName == "" {
			diags.AddError(errorZoneNameNotSet, errorZoneNameNotSet)
			return &tfModels
		}
		specContainerIDs := selectContainerIDs(&item, containerIDs)
		tfModels[i] = TFEffectiveReplicationSpecsModel{
			ExternalId:    types.StringPointerValue(item.Id),
			ContainerId:   conversion.ToTFMapOfString(ctx, diags, specContainerIDs),
			RegionConfigs: regionConfigs,
			ZoneId:        types.StringPointerValue(item.ZoneId),
			ZoneName:      types.StringPointerValue(item.ZoneName),
		}
	}
	return &tfModels
}

func selectContainerIDs(spec *admin.ReplicationSpec20240805, allIDs map[string]string) map[string]string {
	containerIDs := map[string]string{}
	if allIDs == nil {
		return containerIDs
	}

	regions := spec.GetRegionConfigs()
	for i := range regions {
		regionConfig := regions[i]
		providerName := regionConfig.GetProviderName()
		key := containerIDKey(providerName, regionConfig.GetRegionName())
		value := allIDs[key]
		// Should be no hard failure if not found, as it is not required for TENANT, error responsibility in resolveContainerIDs
		if value == "" {
			continue
		}
		containerIDs[key] = value
	}
	return containerIDs
}

func newTagsObjType(ctx context.Context, diags *diag.Diagnostics, input *[]admin.ResourceTag) types.Map {
	elms := make(map[string]string)
	if input != nil {
		for _, item := range *input {
			elms[item.GetKey()] = item.GetValue()
		}
	}
	return conversion.ToTFMapOfString(ctx, diags, elms)
}

func newPrivateEndpointObjType(ctx context.Context, input *[]admin.ClusterDescriptionConnectionStringsPrivateEndpoint, diags *diag.Diagnostics) types.List {
	if input == nil {
		return types.ListNull(privateEndpointObjType)
	}
	tfModels := make([]TFPrivateEndpointModel, len(*input))
	for i, item := range *input {
		endpoints := newEndpointsObjType(ctx, item.Endpoints, diags)
		tfModels[i] = TFPrivateEndpointModel{
			ConnectionString:                  types.StringPointerValue(item.ConnectionString),
			Endpoints:                         endpoints,
			SrvConnectionString:               types.StringPointerValue(item.SrvConnectionString),
			SrvShardOptimizedConnectionString: types.StringPointerValue(item.SrvShardOptimizedConnectionString),
			Type:                              types.StringPointerValue(item.Type),
		}
	}
	listType, diagsLocal := types.ListValueFrom(ctx, privateEndpointObjType, tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func newRegionConfigModel(ctx context.Context, item *admin.CloudRegionConfig20240805, diags *diag.Diagnostics) TFRegionConfigsModel {
	return TFRegionConfigsModel{
		AnalyticsAutoScaling: newAutoScalingObjType(ctx, item.AnalyticsAutoScaling, diags),
		AnalyticsSpecs:       newSpecsObjType(ctx, item.AnalyticsSpecs, diags),
		AutoScaling:          newAutoScalingObjType(ctx, item.AutoScaling, diags),
		BackingProviderName:  types.StringPointerValue(item.BackingProviderName),
		ElectableSpecs:       newSpecsFromHwObjType(ctx, item.ElectableSpecs, diags),
		Priority:             types.Int64PointerValue(conversion.IntPtrToInt64Ptr(item.Priority)),
		ProviderName:         types.StringPointerValue(item.ProviderName),
		ReadOnlySpecs:        newSpecsObjType(ctx, item.ReadOnlySpecs, diags),
		RegionName:           types.StringPointerValue(item.RegionName),
	}
}

func newRegionConfigsObjType(ctx context.Context, input *[]admin.CloudRegionConfig20240805, diags *diag.Diagnostics) types.List {
	if input == nil {
		return types.ListNull(regionConfigsObjType)
	}
	tfModels := make([]TFRegionConfigsModel, len(*input))
	for i := range *input {
		tfModels[i] = newRegionConfigModel(ctx, &(*input)[i], diags)
	}
	listType, diagsLocal := types.ListValueFrom(ctx, regionConfigsObjType, tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func newRegionConfigsDSObjType(ctx context.Context, input *[]admin.CloudRegionConfig20240805, diags *diag.Diagnostics) types.List {
	if input == nil {
		return types.ListNull(regionConfigsDSObjType)
	}
	tfModels := make([]TFRegionConfigsDSModel, len(*input))
	for i := range *input {
		item := &(*input)[i]
		baseModel := newRegionConfigModel(ctx, item, diags)
		dsModel := *conversion.CopyModel[TFRegionConfigsDSModel](&baseModel)
		tfModels[i] = dsModel
	}
	listType, diagsLocal := types.ListValueFrom(ctx, regionConfigsDSObjType, tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func newEndpointsObjType(ctx context.Context, input *[]admin.ClusterDescriptionConnectionStringsPrivateEndpointEndpoint, diags *diag.Diagnostics) types.List {
	if input == nil {
		return types.ListNull(endpointsObjType)
	}
	tfModels := make([]TFEndpointsModel, len(*input))
	for i, item := range *input {
		tfModels[i] = TFEndpointsModel{
			EndpointId:   types.StringPointerValue(item.EndpointId),
			ProviderName: types.StringPointerValue(item.ProviderName),
			Region:       types.StringPointerValue(item.Region),
		}
	}
	listType, diagsLocal := types.ListValueFrom(ctx, endpointsObjType, tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func newSpecsObjType(ctx context.Context, input *admin.DedicatedHardwareSpec20240805, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(specsObjType.AttrTypes)
	}
	tfModel := TFSpecsModel{
		DiskIops:      types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.DiskIOPS)),
		DiskSizeGb:    types.Float64PointerValue(input.DiskSizeGB),
		EbsVolumeType: types.StringPointerValue(input.EbsVolumeType),
		InstanceSize:  types.StringPointerValue(input.InstanceSize),
		NodeCount:     types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.NodeCount)),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

func newSpecsFromHwObjType(ctx context.Context, input *admin.HardwareSpec20240805, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(specsObjType.AttrTypes)
	}
	tfModel := TFSpecsModel{
		DiskIops:      types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.DiskIOPS)),
		DiskSizeGb:    types.Float64PointerValue(input.DiskSizeGB),
		EbsVolumeType: types.StringPointerValue(input.EbsVolumeType),
		InstanceSize:  types.StringPointerValue(input.InstanceSize),
		NodeCount:     types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.NodeCount)),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

func newAutoScalingObjType(ctx context.Context, input *admin.AdvancedAutoScalingSettings, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(autoScalingObjType.AttrTypes)
	}
	compute := input.Compute
	tfModel := TFAutoScalingModel{}
	if compute != nil {
		tfModel.ComputeMaxInstanceSize = types.StringPointerValue(compute.MaxInstanceSize)
		tfModel.ComputeMinInstanceSize = types.StringPointerValue(compute.MinInstanceSize)
		tfModel.ComputeEnabled = types.BoolPointerValue(compute.Enabled)
		tfModel.ComputeScaleDownEnabled = types.BoolPointerValue(compute.ScaleDownEnabled)
	}
	diskGB := input.DiskGB
	if diskGB != nil {
		tfModel.DiskGBEnabled = types.BoolPointerValue(diskGB.Enabled)
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, autoScalingObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}
