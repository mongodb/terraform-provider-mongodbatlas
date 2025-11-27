package advancedcluster

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

const (
	errorZoneNameNotSet = "zoneName is required for legacy schema"
)

func newTFModel(ctx context.Context, input *admin.ClusterDescription20240805, diags *diag.Diagnostics, containerIDs map[string]string) *TFModel {
	biConnector := newBiConnectorConfigObjType(ctx, input.BiConnector, diags)
	connectionStrings := newConnectionStringsObjType(ctx, input.ConnectionStrings, diags)
	labels := newLabelsObjType(ctx, diags, input.Labels)
	replicationSpecs := newReplicationSpecsObjType(ctx, input.ReplicationSpecs, diags, containerIDs)
	tags := newTagsObjType(ctx, diags, input.Tags)
	pinnedFCV := newPinnedFCVObjType(ctx, input, diags)
	if diags.HasError() {
		return nil
	}
	return &TFModel{
		AcceptDataRisksAndForceReplicaSetReconfig: types.StringPointerValue(conversion.TimePtrToStringPtr(input.AcceptDataRisksAndForceReplicaSetReconfig)),
		BackupEnabled:                    types.BoolValue(conversion.SafeValue(input.BackupEnabled)),
		BiConnectorConfig:                biConnector,
		ClusterType:                      types.StringValue(conversion.SafeValue(input.ClusterType)),
		ConfigServerManagementMode:       types.StringValue(conversion.SafeValue(input.ConfigServerManagementMode)),
		ConfigServerType:                 types.StringValue(conversion.SafeValue(input.ConfigServerType)),
		ConnectionStrings:                connectionStrings,
		CreateDate:                       types.StringValue(conversion.SafeValue(conversion.TimePtrToStringPtr(input.CreateDate))),
		EncryptionAtRestProvider:         types.StringValue(conversion.SafeValue(input.EncryptionAtRestProvider)),
		GlobalClusterSelfManagedSharding: types.BoolValue(conversion.SafeValue(input.GlobalClusterSelfManagedSharding)),
		ProjectID:                        types.StringValue(conversion.SafeValue(input.GroupId)),
		ClusterID:                        types.StringValue(conversion.SafeValue(input.Id)),
		Labels:                           labels,
		MongoDBMajorVersion:              types.StringValue(conversion.SafeValue(input.MongoDBMajorVersion)),
		MongoDBVersion:                   types.StringValue(conversion.SafeValue(input.MongoDBVersion)),
		Name:                             types.StringValue(conversion.SafeValue(input.Name)),
		Paused:                           types.BoolValue(conversion.SafeValue(input.Paused)),
		PitEnabled:                       types.BoolValue(conversion.SafeValue(input.PitEnabled)),
		RedactClientLogData:              types.BoolValue(conversion.SafeValue(input.RedactClientLogData)),
		ReplicaSetScalingStrategy:        types.StringValue(conversion.SafeValue(input.ReplicaSetScalingStrategy)),
		ReplicationSpecs:                 replicationSpecs,
		RootCertType:                     types.StringValue(conversion.SafeValue(input.RootCertType)),
		StateName:                        types.StringValue(conversion.SafeValue(input.StateName)),
		Tags:                             tags,
		TerminationProtectionEnabled:     types.BoolValue(conversion.SafeValue(input.TerminationProtectionEnabled)),
		VersionReleaseSystem:             types.StringValue(conversion.SafeValue(input.VersionReleaseSystem)),
		PinnedFCV:                        pinnedFCV,
	}
}

func newTFModelDS(ctx context.Context, input *admin.ClusterDescription20240805, diags *diag.Diagnostics, containerIDs map[string]string) *TFModelDS {
	resourceModel := newTFModel(ctx, input, diags, containerIDs)
	if diags.HasError() {
		return nil
	}
	dsModel := conversion.CopyModel[TFModelDS](resourceModel)
	dsModel.ReplicationSpecs = newReplicationSpecsDSObjType(ctx, input.ReplicationSpecs, diags, containerIDs)
	return dsModel
}

func newBiConnectorConfigObjType(ctx context.Context, input *admin.BiConnector, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(biConnectorConfigObjType.AttrTypes)
	}
	tfModel := TFBiConnectorModel{
		Enabled:        types.BoolValue(conversion.SafeValue(input.Enabled)),
		ReadPreference: types.StringValue(conversion.SafeValue(input.ReadPreference)),
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
		Private:         types.StringValue(conversion.SafeValue(input.Private)),
		PrivateEndpoint: privateEndpoint,
		PrivateSrv:      types.StringValue(conversion.SafeValue(input.PrivateSrv)),
		Standard:        types.StringValue(conversion.SafeValue(input.Standard)),
		StandardSrv:     types.StringValue(conversion.SafeValue(input.StandardSrv)),
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

func newReplicationSpecsObjType(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, containerIDs map[string]string) types.List {
	if input == nil {
		return types.ListNull(replicationSpecsObjType)
	}
	tfModels := convertReplicationSpecs(ctx, input, diags, containerIDs, newRegionConfigsObjType)
	if diags.HasError() {
		return types.ListNull(replicationSpecsObjType)
	}
	listType, diagsLocal := types.ListValueFrom(ctx, replicationSpecsObjType, *tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func newReplicationSpecsDSObjType(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, containerIDs map[string]string) types.List {
	if input == nil {
		return types.ListNull(replicationSpecsDSObjType)
	}
	tfModels := convertReplicationSpecs(ctx, input, diags, containerIDs, newRegionConfigsDSObjType)
	if diags.HasError() {
		return types.ListNull(replicationSpecsDSObjType)
	}
	listType, diagsLocal := types.ListValueFrom(ctx, replicationSpecsDSObjType, *tfModels)
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

func convertReplicationSpecs(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, containerIDs map[string]string, regionConfigsConv regionConfigsConverter) *[]TFReplicationSpecsModel {
	tfModels := make([]TFReplicationSpecsModel, len(*input))
	for i, item := range *input {
		regionConfigs := regionConfigsConv(ctx, item.RegionConfigs, diags)
		zoneName := item.GetZoneName()
		if zoneName == "" {
			diags.AddError(errorZoneNameNotSet, errorZoneNameNotSet)
			return &tfModels
		}
		containerIDs := selectContainerIDs(&item, containerIDs)
		tfModels[i] = TFReplicationSpecsModel{
			ExternalId:    types.StringValue(conversion.SafeValue(item.Id)),
			ContainerId:   conversion.ToTFMapOfString(ctx, diags, containerIDs),
			RegionConfigs: regionConfigs,
			ZoneId:        types.StringValue(conversion.SafeValue(item.ZoneId)),
			ZoneName:      types.StringValue(conversion.SafeValue(item.ZoneName)),
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
			ConnectionString:                  types.StringValue(conversion.SafeValue(item.ConnectionString)),
			Endpoints:                         endpoints,
			SrvConnectionString:               types.StringValue(conversion.SafeValue(item.SrvConnectionString)),
			SrvShardOptimizedConnectionString: types.StringValue(conversion.SafeValue(item.SrvShardOptimizedConnectionString)),
			Type:                              types.StringValue(conversion.SafeValue(item.Type)),
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
		ProviderName:         types.StringValue(conversion.SafeValue(item.ProviderName)),
		ReadOnlySpecs:        newSpecsObjType(ctx, item.ReadOnlySpecs, diags),
		RegionName:           types.StringValue(conversion.SafeValue(item.RegionName)),
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
		dsModel.EffectiveAnalyticsSpecs = newSpecsObjType(ctx, item.EffectiveAnalyticsSpecs, diags)
		dsModel.EffectiveElectableSpecs = newSpecsObjType(ctx, item.EffectiveElectableSpecs, diags)
		dsModel.EffectiveReadOnlySpecs = newSpecsObjType(ctx, item.EffectiveReadOnlySpecs, diags)
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
			EndpointId:   types.StringValue(conversion.SafeValue(item.EndpointId)),
			ProviderName: types.StringValue(conversion.SafeValue(item.ProviderName)),
			Region:       types.StringValue(conversion.SafeValue(item.Region)),
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
		EbsVolumeType: types.StringValue(conversion.SafeValue(input.EbsVolumeType)),
		InstanceSize:  types.StringValue(conversion.SafeValue(input.InstanceSize)),
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
		EbsVolumeType: types.StringValue(conversion.SafeValue(input.EbsVolumeType)),
		InstanceSize:  types.StringValue(conversion.SafeValue(input.InstanceSize)),
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
		tfModel.ComputeMaxInstanceSize = types.StringValue(conversion.SafeValue(compute.MaxInstanceSize))
		tfModel.ComputeMinInstanceSize = types.StringValue(conversion.SafeValue(compute.MinInstanceSize))
		tfModel.ComputeEnabled = types.BoolValue(conversion.SafeValue(compute.Enabled))
		tfModel.ComputeScaleDownEnabled = types.BoolValue(conversion.SafeValue(compute.ScaleDownEnabled))
	}
	diskGB := input.DiskGB
	if diskGB != nil {
		tfModel.DiskGBEnabled = types.BoolValue(conversion.SafeValue(diskGB.Enabled))
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, autoScalingObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}
