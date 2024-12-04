package advancedclustertpf

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

const (
	errorZoneNameNotSet          = "zoneName is required for legacy schema"
	errorNumShardsNotSet         = "numShards not set for zoneName %s"
	errorReplicationSpecIDNotSet = "replicationSpecID not set for zoneName %s"
)

type LegacySchemaInfo struct {
	ZoneNameNumShards          map[string]int64
	ZoneNameReplicationSpecIDs map[string]string
	RootDiskSize               *float64
}

type ExtraAPIInfo struct {
	ContainerIDs map[string]string
}

func NewTFModel(ctx context.Context, input *admin.ClusterDescription20240805, timeout timeouts.Value, diags *diag.Diagnostics, legacyInfo *LegacySchemaInfo, apiInfo ExtraAPIInfo) *TFModel {
	biConnector := NewBiConnectorConfigObjType(ctx, input.BiConnector, diags)
	connectionStrings := NewConnectionStringsObjType(ctx, input.ConnectionStrings, diags)
	labels := NewLabelsObjType(ctx, input.Labels, diags)
	replicationSpecs := NewReplicationSpecsObjType(ctx, input.ReplicationSpecs, diags, legacyInfo, apiInfo)
	tags := NewTagsObjType(ctx, input.Tags, diags)
	if diags.HasError() {
		return nil
	}
	return &TFModel{
		AcceptDataRisksAndForceReplicaSetReconfig: types.StringPointerValue(conversion.TimePtrToStringPtr(input.AcceptDataRisksAndForceReplicaSetReconfig)),
		BackupEnabled:                    types.BoolPointerValue(input.BackupEnabled),
		BiConnectorConfig:                biConnector,
		ClusterType:                      types.StringPointerValue(input.ClusterType),
		ConfigServerManagementMode:       types.StringPointerValue(input.ConfigServerManagementMode),
		ConfigServerType:                 types.StringPointerValue(input.ConfigServerType),
		ConnectionStrings:                connectionStrings,
		CreateDate:                       types.StringPointerValue(conversion.TimePtrToStringPtr(input.CreateDate)),
		DiskSizeGB:                       types.Float64PointerValue(findRegionRootDiskSize(input.ReplicationSpecs)),
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
		RootCertType:                     types.StringPointerValue(input.RootCertType),
		StateName:                        types.StringPointerValue(input.StateName),
		Tags:                             tags,
		TerminationProtectionEnabled:     types.BoolPointerValue(input.TerminationProtectionEnabled),
		VersionReleaseSystem:             types.StringPointerValue(input.VersionReleaseSystem),
		Timeouts:                         timeout,
	}
}

func NewBiConnectorConfigObjType(ctx context.Context, input *admin.BiConnector, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(BiConnectorConfigObjType.AttrTypes)
	}
	tfModel := TFBiConnectorModel{
		Enabled:        types.BoolPointerValue(input.Enabled),
		ReadPreference: types.StringPointerValue(input.ReadPreference),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, BiConnectorConfigObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

func NewConnectionStringsObjType(ctx context.Context, input *admin.ClusterConnectionStrings, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(ConnectionStringsObjType.AttrTypes)
	}
	privateEndpoint := NewPrivateEndpointObjType(ctx, input.PrivateEndpoint, diags)
	tfModel := TFConnectionStringsModel{
		Private:         types.StringPointerValue(input.Private),
		PrivateEndpoint: privateEndpoint,
		PrivateSrv:      types.StringPointerValue(input.PrivateSrv),
		Standard:        types.StringPointerValue(input.Standard),
		StandardSrv:     types.StringPointerValue(input.StandardSrv),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, ConnectionStringsObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

func NewLabelsObjType(ctx context.Context, input *[]admin.ComponentLabel, diags *diag.Diagnostics) types.Set {
	if input == nil {
		return types.SetNull(LabelsObjType)
	}
	tfModels := make([]TFLabelsModel, len(*input))
	for i, item := range *input {
		tfModels[i] = TFLabelsModel{
			Key:   types.StringPointerValue(item.Key),
			Value: types.StringPointerValue(item.Value),
		}
	}
	setType, diagsLocal := types.SetValueFrom(ctx, LabelsObjType, tfModels)
	diags.Append(diagsLocal...)
	return setType
}

func NewReplicationSpecsObjType(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, legacyInfo *LegacySchemaInfo, apiInfo ExtraAPIInfo) types.List {
	if input == nil {
		return types.ListNull(ReplicationSpecsObjType)
	}
	var tfModels *[]TFReplicationSpecsModel
	if legacyInfo == nil {
		tfModels = convertReplicationSpecs(ctx, input, diags, apiInfo)
	} else {
		tfModels = convertReplicationSpecsLegacy(ctx, input, diags, legacyInfo, apiInfo)
	}
	if diags.HasError() {
		return types.ListNull(ReplicationSpecsObjType)
	}
	listType, diagsLocal := types.ListValueFrom(ctx, ReplicationSpecsObjType, *tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func convertReplicationSpecs(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, apiInfo ExtraAPIInfo) *[]TFReplicationSpecsModel {
	tfModels := make([]TFReplicationSpecsModel, len(*input))
	for i, item := range *input {
		regionConfigs := NewRegionConfigsObjType(ctx, item.RegionConfigs, diags)
		tfModels[i] = TFReplicationSpecsModel{
			Id:            types.StringPointerValue(item.Id),
			ExternalId:    types.StringPointerValue(item.Id),
			NumShards:     types.Int64Value(1), // TODO: Static
			ContainerId:   conversion.ToTFMapOfString(ctx, diags, &apiInfo.ContainerIDs),
			RegionConfigs: regionConfigs,
			ZoneId:        types.StringPointerValue(item.ZoneId),
			ZoneName:      types.StringPointerValue(item.ZoneName),
		}
	}
	return &tfModels
}

func convertReplicationSpecsLegacy(ctx context.Context, input *[]admin.ReplicationSpec20240805, diags *diag.Diagnostics, legacyInfo *LegacySchemaInfo, apiInfo ExtraAPIInfo) *[]TFReplicationSpecsModel {
	tfModels := []TFReplicationSpecsModel{}
	tfModelsSkipIndexes := []int{}
	for i, item := range *input {
		if slices.Contains(tfModelsSkipIndexes, i) {
			continue
		}
		regionConfigs := NewRegionConfigsObjType(ctx, item.RegionConfigs, diags)
		zoneName := item.GetZoneName()
		if zoneName == "" {
			diags.AddError(errorZoneNameNotSet, errorZoneNameNotSet)
			return &tfModels
		}
		numShards, ok := legacyInfo.ZoneNameNumShards[zoneName]
		errMsg := []string{}
		if !ok {
			errMsg = append(errMsg, fmt.Sprintf(errorNumShardsNotSet, zoneName))
		}
		legacyID, ok := legacyInfo.ZoneNameReplicationSpecIDs[zoneName]
		if !ok {
			errMsg = append(errMsg, fmt.Sprintf(errorReplicationSpecIDNotSet, zoneName))
		}
		if len(errMsg) > 0 {
			diags.AddError("replicationSpecsLegacySchema", strings.Join(errMsg, ", "))
			return &tfModels
		}
		if numShards > 1 {
			for j := 1; j < int(numShards); j++ {
				tfModelsSkipIndexes = append(tfModelsSkipIndexes, i+j)
			}
		}
		tfModels = append(tfModels, TFReplicationSpecsModel{
			ContainerId:   conversion.ToTFMapOfString(ctx, diags, &apiInfo.ContainerIDs),
			ExternalId:    types.StringPointerValue(item.Id),
			Id:            types.StringValue(legacyID),
			RegionConfigs: regionConfigs,
			NumShards:     types.Int64Value(numShards),
			ZoneId:        types.StringPointerValue(item.ZoneId),
			ZoneName:      types.StringPointerValue(item.ZoneName),
		})
	}
	return &tfModels
}

func NewTagsObjType(ctx context.Context, input *[]admin.ResourceTag, diags *diag.Diagnostics) types.Set {
	if input == nil {
		// API Response not consistent, even when not set in POST/PATCH `[]` is returned instead of null
		return types.SetValueMust(TagsObjType, nil)
	}
	tfModels := make([]TFTagsModel, len(*input))
	for i, item := range *input {
		tfModels[i] = TFTagsModel{
			Key:   types.StringValue(item.Key),
			Value: types.StringValue(item.Value),
		}
	}
	setType, diagsLocal := types.SetValueFrom(ctx, TagsObjType, tfModels)
	diags.Append(diagsLocal...)
	return setType
}

func NewPrivateEndpointObjType(ctx context.Context, input *[]admin.ClusterDescriptionConnectionStringsPrivateEndpoint, diags *diag.Diagnostics) types.List {
	if input == nil {
		return types.ListNull(PrivateEndpointObjType)
	}
	tfModels := make([]TFPrivateEndpointModel, len(*input))
	for i, item := range *input {
		endpoints := NewEndpointsObjType(ctx, item.Endpoints, diags)
		tfModels[i] = TFPrivateEndpointModel{
			ConnectionString:                  types.StringPointerValue(item.ConnectionString),
			Endpoints:                         endpoints,
			SrvConnectionString:               types.StringPointerValue(item.SrvConnectionString),
			SrvShardOptimizedConnectionString: types.StringPointerValue(item.SrvShardOptimizedConnectionString),
			Type:                              types.StringPointerValue(item.Type),
		}
	}
	listType, diagsLocal := types.ListValueFrom(ctx, PrivateEndpointObjType, tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func NewRegionConfigsObjType(ctx context.Context, input *[]admin.CloudRegionConfig20240805, diags *diag.Diagnostics) types.List {
	if input == nil {
		return types.ListNull(RegionConfigsObjType)
	}
	tfModels := make([]TFRegionConfigsModel, len(*input))
	for i, item := range *input {
		analyticsAutoScaling := NewAutoScalingObjType(ctx, item.AnalyticsAutoScaling, diags)
		analyticsSpecs := NewSpecsObjType(ctx, item.AnalyticsSpecs, diags)
		autoScaling := NewAutoScalingObjType(ctx, item.AutoScaling, diags)
		electableSpecs := NewSpecsFromHwObjType(ctx, item.ElectableSpecs, diags)
		readOnlySpecs := NewSpecsObjType(ctx, item.ReadOnlySpecs, diags)
		tfModels[i] = TFRegionConfigsModel{
			AnalyticsAutoScaling: analyticsAutoScaling,
			AnalyticsSpecs:       analyticsSpecs,
			AutoScaling:          autoScaling,
			BackingProviderName:  types.StringPointerValue(item.BackingProviderName),
			ElectableSpecs:       electableSpecs,
			Priority:             types.Int64PointerValue(conversion.IntPtrToInt64Ptr(item.Priority)),
			ProviderName:         types.StringPointerValue(item.ProviderName),
			ReadOnlySpecs:        readOnlySpecs,
			RegionName:           types.StringPointerValue(item.RegionName),
		}
	}
	listType, diagsLocal := types.ListValueFrom(ctx, RegionConfigsObjType, tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func NewEndpointsObjType(ctx context.Context, input *[]admin.ClusterDescriptionConnectionStringsPrivateEndpointEndpoint, diags *diag.Diagnostics) types.List {
	if input == nil {
		return types.ListNull(EndpointsObjType)
	}
	tfModels := make([]TFEndpointsModel, len(*input))
	for i, item := range *input {
		tfModels[i] = TFEndpointsModel{
			EndpointId:   types.StringPointerValue(item.EndpointId),
			ProviderName: types.StringPointerValue(item.ProviderName),
			Region:       types.StringPointerValue(item.Region),
		}
	}
	listType, diagsLocal := types.ListValueFrom(ctx, EndpointsObjType, tfModels)
	diags.Append(diagsLocal...)
	return listType
}

func NewSpecsObjType(ctx context.Context, input *admin.DedicatedHardwareSpec20240805, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(SpecsObjType.AttrTypes)
	}
	tfModel := TFSpecsModel{
		DiskIops:      types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.DiskIOPS)),
		DiskSizeGb:    types.Float64PointerValue(input.DiskSizeGB),
		EbsVolumeType: types.StringPointerValue(input.EbsVolumeType),
		InstanceSize:  types.StringPointerValue(input.InstanceSize),
		NodeCount:     types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.NodeCount)),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, SpecsObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

func NewSpecsFromHwObjType(ctx context.Context, input *admin.HardwareSpec20240805, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(SpecsObjType.AttrTypes)
	}
	tfModel := TFSpecsModel{
		DiskIops:      types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.DiskIOPS)),
		DiskSizeGb:    types.Float64PointerValue(input.DiskSizeGB),
		EbsVolumeType: types.StringPointerValue(input.EbsVolumeType),
		InstanceSize:  types.StringPointerValue(input.InstanceSize),
		NodeCount:     types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.NodeCount)),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, SpecsObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

func NewAutoScalingObjType(ctx context.Context, input *admin.AdvancedAutoScalingSettings, diags *diag.Diagnostics) types.Object {
	if input == nil {
		return types.ObjectNull(AutoScalingObjType.AttrTypes)
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
	objType, diagsLocal := types.ObjectValueFrom(ctx, AutoScalingObjType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}
