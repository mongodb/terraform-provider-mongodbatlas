package advancedcluster

import (
	admin20231115 "go.mongodb.org/atlas-sdk/v20231115014/admin"
	"go.mongodb.org/atlas-sdk/v20240530001/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Conversions from one SDK model version to another are used to avoid duplicating our flatten/expand conversion functions.
// - These functions must not contain any business logic.
// - All will be removed once we rely on a single API version.

func convertTagsToLatest(tags *[]admin20231115.ResourceTag) *[]admin.ResourceTag {
	if tags == nil {
		return nil
	}
	tagSlice := *tags
	results := make([]admin.ResourceTag, len(tagSlice))
	for i := range len(tagSlice) {
		tag := tagSlice[i]
		results[i] = admin.ResourceTag{
			Key:   tag.Key,
			Value: tag.Value,
		}
	}
	return &results
}

func convertBiConnectToOldSDK(biconnector *admin.BiConnector) *admin20231115.BiConnector {
	if biconnector == nil {
		return nil
	}
	return &admin20231115.BiConnector{
		Enabled:        biconnector.Enabled,
		ReadPreference: biconnector.ReadPreference,
	}
}

func convertBiConnectToLatest(biconnector *admin20231115.BiConnector) *admin.BiConnector {
	return &admin.BiConnector{
		Enabled:        biconnector.Enabled,
		ReadPreference: biconnector.ReadPreference,
	}
}

func convertConnectionStringToLatest(connStrings *admin20231115.ClusterConnectionStrings) *admin.ClusterConnectionStrings {
	return &admin.ClusterConnectionStrings{
		AwsPrivateLink:    connStrings.AwsPrivateLink,
		AwsPrivateLinkSrv: connStrings.AwsPrivateLinkSrv,
		Private:           connStrings.Private,
		PrivateEndpoint:   convertPrivateEndpointToLatest(connStrings.PrivateEndpoint),
		PrivateSrv:        connStrings.PrivateSrv,
		Standard:          connStrings.Standard,
		StandardSrv:       connStrings.StandardSrv,
	}
}

func convertPrivateEndpointToLatest(privateEndpoints *[]admin20231115.ClusterDescriptionConnectionStringsPrivateEndpoint) *[]admin.ClusterDescriptionConnectionStringsPrivateEndpoint {
	if privateEndpoints == nil {
		return nil
	}
	peSlice := *privateEndpoints
	results := make([]admin.ClusterDescriptionConnectionStringsPrivateEndpoint, len(peSlice))
	for i := range len(peSlice) {
		pe := peSlice[i]
		results[i] = admin.ClusterDescriptionConnectionStringsPrivateEndpoint{
			ConnectionString:                  pe.ConnectionString,
			Endpoints:                         convertEndpointsToLatest(pe.Endpoints),
			SrvConnectionString:               pe.SrvConnectionString,
			SrvShardOptimizedConnectionString: pe.SrvShardOptimizedConnectionString,
			Type:                              pe.Type,
		}
	}
	return &results
}

func convertEndpointsToLatest(privateEndpoints *[]admin20231115.ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) *[]admin.ClusterDescriptionConnectionStringsPrivateEndpointEndpoint {
	if privateEndpoints == nil {
		return nil
	}
	peSlice := *privateEndpoints
	results := make([]admin.ClusterDescriptionConnectionStringsPrivateEndpointEndpoint, len(peSlice))
	for i := range len(peSlice) {
		pe := peSlice[i]
		results[i] = admin.ClusterDescriptionConnectionStringsPrivateEndpointEndpoint{
			EndpointId:   pe.EndpointId,
			ProviderName: pe.ProviderName,
			Region:       pe.Region,
		}
	}
	return &results
}

func convertLabelsToLatest(labels *[]admin20231115.ComponentLabel) *[]admin.ComponentLabel {
	labelSlice := *labels
	results := make([]admin.ComponentLabel, len(labelSlice))
	for i := range len(labelSlice) {
		label := labelSlice[i]
		results[i] = admin.ComponentLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	return &results
}

func convertLabelSliceToOldSDK(slice []admin.ComponentLabel, err diag.Diagnostics) ([]admin20231115.ComponentLabel, diag.Diagnostics) {
	if err != nil {
		return nil, err
	}
	results := make([]admin20231115.ComponentLabel, len(slice))
	for i := range len(slice) {
		label := slice[i]
		results[i] = admin20231115.ComponentLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	return results, nil
}

func convertRegionConfigSliceToOldSDK(slice *[]admin.CloudRegionConfig) *[]admin20231115.CloudRegionConfig {
	if slice == nil {
		return nil
	}
	cloudRegionSlice := *slice
	results := make([]admin20231115.CloudRegionConfig, len(cloudRegionSlice))
	for i := range len(cloudRegionSlice) {
		cloudRegion := cloudRegionSlice[i]
		results[i] = admin20231115.CloudRegionConfig{
			ElectableSpecs:       convertHardwareSpecToOldSDK(cloudRegion.ElectableSpecs),
			Priority:             cloudRegion.Priority,
			ProviderName:         cloudRegion.ProviderName,
			RegionName:           cloudRegion.RegionName,
			AnalyticsAutoScaling: convertAdvancedAutoScalingSettingsToOldSDK(cloudRegion.AnalyticsAutoScaling),
			AnalyticsSpecs:       convertDedicatedHardwareSpecToOldSDK(cloudRegion.AnalyticsSpecs),
			AutoScaling:          convertAdvancedAutoScalingSettingsToOldSDK(cloudRegion.AutoScaling),
			ReadOnlySpecs:        convertDedicatedHardwareSpecToOldSDK(cloudRegion.ReadOnlySpecs),
			BackingProviderName:  cloudRegion.BackingProviderName,
		}
	}
	return &results
}

func convertHardwareSpecToOldSDK(hwspec *admin.HardwareSpec) *admin20231115.HardwareSpec {
	if hwspec == nil {
		return nil
	}
	return &admin20231115.HardwareSpec{
		DiskIOPS:      hwspec.DiskIOPS,
		EbsVolumeType: hwspec.EbsVolumeType,
		InstanceSize:  hwspec.InstanceSize,
		NodeCount:     hwspec.NodeCount,
	}
}

func convertAdvancedAutoScalingSettingsToOldSDK(settings *admin.AdvancedAutoScalingSettings) *admin20231115.AdvancedAutoScalingSettings {
	if settings == nil {
		return nil
	}
	return &admin20231115.AdvancedAutoScalingSettings{
		Compute: convertAdvancedComputeAutoScalingToOldSDK(settings.Compute),
		DiskGB:  convertDiskGBAutoScalingToOldSDK(settings.DiskGB),
	}
}

func convertAdvancedComputeAutoScalingToOldSDK(settings *admin.AdvancedComputeAutoScaling) *admin20231115.AdvancedComputeAutoScaling {
	if settings == nil {
		return nil
	}
	return &admin20231115.AdvancedComputeAutoScaling{
		Enabled:          settings.Enabled,
		MaxInstanceSize:  settings.MaxInstanceSize,
		MinInstanceSize:  settings.MinInstanceSize,
		ScaleDownEnabled: settings.ScaleDownEnabled,
	}
}

func convertDiskGBAutoScalingToOldSDK(settings *admin.DiskGBAutoScaling) *admin20231115.DiskGBAutoScaling {
	if settings == nil {
		return nil
	}
	return &admin20231115.DiskGBAutoScaling{
		Enabled: settings.Enabled,
	}
}

func convertDedicatedHardwareSpecToOldSDK(spec *admin.DedicatedHardwareSpec) *admin20231115.DedicatedHardwareSpec {
	if spec == nil {
		return nil
	}
	return &admin20231115.DedicatedHardwareSpec{
		NodeCount:     spec.NodeCount,
		DiskIOPS:      spec.DiskIOPS,
		EbsVolumeType: spec.EbsVolumeType,
		InstanceSize:  spec.InstanceSize,
	}
}

func convertDedicatedHwSpecToLatest(spec *admin20231115.DedicatedHardwareSpec, rootDiskSizeGB *float64) *admin.DedicatedHardwareSpec {
	if spec == nil {
		return nil
	}
	return &admin.DedicatedHardwareSpec{
		NodeCount:     spec.NodeCount,
		DiskIOPS:      spec.DiskIOPS,
		EbsVolumeType: spec.EbsVolumeType,
		InstanceSize:  spec.InstanceSize,
		DiskSizeGB:    rootDiskSizeGB,
	}
}

func convertAdvancedAutoScalingSettingsToLatest(settings *admin20231115.AdvancedAutoScalingSettings) *admin.AdvancedAutoScalingSettings {
	if settings == nil {
		return nil
	}
	return &admin.AdvancedAutoScalingSettings{
		Compute: convertAdvancedComputeAutoScalingToLatest(settings.Compute),
		DiskGB:  convertDiskGBAutoScalingToLatest(settings.DiskGB),
	}
}

func convertAdvancedComputeAutoScalingToLatest(settings *admin20231115.AdvancedComputeAutoScaling) *admin.AdvancedComputeAutoScaling {
	if settings == nil {
		return nil
	}
	return &admin.AdvancedComputeAutoScaling{
		Enabled:          settings.Enabled,
		MaxInstanceSize:  settings.MaxInstanceSize,
		MinInstanceSize:  settings.MinInstanceSize,
		ScaleDownEnabled: settings.ScaleDownEnabled,
	}
}

func convertDiskGBAutoScalingToLatest(settings *admin20231115.DiskGBAutoScaling) *admin.DiskGBAutoScaling {
	if settings == nil {
		return nil
	}
	return &admin.DiskGBAutoScaling{
		Enabled: settings.Enabled,
	}
}

func convertHardwareSpecToLatest(hwspec *admin20231115.HardwareSpec, rootDiskSizeGB *float64) *admin.HardwareSpec {
	if hwspec == nil {
		return nil
	}
	return &admin.HardwareSpec{
		DiskIOPS:      hwspec.DiskIOPS,
		EbsVolumeType: hwspec.EbsVolumeType,
		InstanceSize:  hwspec.InstanceSize,
		NodeCount:     hwspec.NodeCount,
		DiskSizeGB:    rootDiskSizeGB,
	}
}

func convertRegionConfigSliceToLatest(slice *[]admin20231115.CloudRegionConfig, rootDiskSizeGB *float64) *[]admin.CloudRegionConfig {
	if slice == nil {
		return nil
	}
	cloudRegionSlice := *slice
	results := make([]admin.CloudRegionConfig, len(cloudRegionSlice))
	for i := range len(cloudRegionSlice) {
		cloudRegion := cloudRegionSlice[i]
		results[i] = admin.CloudRegionConfig{
			ElectableSpecs:       convertHardwareSpecToLatest(cloudRegion.ElectableSpecs, rootDiskSizeGB),
			Priority:             cloudRegion.Priority,
			ProviderName:         cloudRegion.ProviderName,
			RegionName:           cloudRegion.RegionName,
			AnalyticsAutoScaling: convertAdvancedAutoScalingSettingsToLatest(cloudRegion.AnalyticsAutoScaling),
			AnalyticsSpecs:       convertDedicatedHwSpecToLatest(cloudRegion.AnalyticsSpecs, rootDiskSizeGB),
			AutoScaling:          convertAdvancedAutoScalingSettingsToLatest(cloudRegion.AutoScaling),
			ReadOnlySpecs:        convertDedicatedHwSpecToLatest(cloudRegion.ReadOnlySpecs, rootDiskSizeGB),
			BackingProviderName:  cloudRegion.BackingProviderName,
		}
	}
	return &results
}

func convertClusterDescToLatestExcludeRepSpecs(oldClusterDesc *admin20231115.AdvancedClusterDescription) *admin.ClusterDescription20240710 {
	return &admin.ClusterDescription20240710{
		BackupEnabled: oldClusterDesc.BackupEnabled,
		AcceptDataRisksAndForceReplicaSetReconfig: oldClusterDesc.AcceptDataRisksAndForceReplicaSetReconfig,
		ClusterType:                      oldClusterDesc.ClusterType,
		CreateDate:                       oldClusterDesc.CreateDate,
		DiskWarmingMode:                  oldClusterDesc.DiskWarmingMode,
		EncryptionAtRestProvider:         oldClusterDesc.EncryptionAtRestProvider,
		GlobalClusterSelfManagedSharding: oldClusterDesc.GlobalClusterSelfManagedSharding,
		GroupId:                          oldClusterDesc.GroupId,
		Id:                               oldClusterDesc.Id,
		MongoDBMajorVersion:              oldClusterDesc.MongoDBMajorVersion,
		MongoDBVersion:                   oldClusterDesc.MongoDBVersion,
		Name:                             oldClusterDesc.Name,
		Paused:                           oldClusterDesc.Paused,
		PitEnabled:                       oldClusterDesc.PitEnabled,
		RootCertType:                     oldClusterDesc.RootCertType,
		StateName:                        oldClusterDesc.StateName,
		TerminationProtectionEnabled:     oldClusterDesc.TerminationProtectionEnabled,
		VersionReleaseSystem:             oldClusterDesc.VersionReleaseSystem,
		Tags:                             convertTagsToLatest(oldClusterDesc.Tags),
		BiConnector:                      convertBiConnectToLatest(oldClusterDesc.BiConnector),
		ConnectionStrings:                convertConnectionStringToLatest(oldClusterDesc.ConnectionStrings),
		Labels:                           convertLabelsToLatest(oldClusterDesc.Labels),
	}
}
