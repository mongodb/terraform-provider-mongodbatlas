package advancedcluster

import (
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Conversions from one SDK model version to another are used to avoid duplicating our flatten/expand conversion functions.
// - These functions must not contain any business logic.
// - All will be removed once we rely on a single API version.

func convertTagsPtrToLatest(tags *[]admin20240530.ResourceTag) *[]admin.ResourceTag {
	if tags == nil {
		return nil
	}
	result := convertTagsToLatest(*tags)
	return &result
}

func convertTagsPtrToOldSDK(tags *[]admin.ResourceTag) *[]admin20240530.ResourceTag {
	if tags == nil {
		return nil
	}
	tagsSlice := *tags
	results := make([]admin20240530.ResourceTag, len(tagsSlice))
	for i := range len(tagsSlice) {
		tag := tagsSlice[i]
		results[i] = admin20240530.ResourceTag{
			Key:   tag.Key,
			Value: tag.Value,
		}
	}
	return &results
}

func convertTagsToLatest(tags []admin20240530.ResourceTag) []admin.ResourceTag {
	results := make([]admin.ResourceTag, len(tags))
	for i := range len(tags) {
		tag := tags[i]
		results[i] = admin.ResourceTag{
			Key:   tag.Key,
			Value: tag.Value,
		}
	}
	return results
}

func convertBiConnectToOldSDK(biconnector *admin.BiConnector) *admin20240530.BiConnector {
	if biconnector == nil {
		return nil
	}
	return &admin20240530.BiConnector{
		Enabled:        biconnector.Enabled,
		ReadPreference: biconnector.ReadPreference,
	}
}

func convertBiConnectToLatest(biconnector *admin20240530.BiConnector) *admin.BiConnector {
	return &admin.BiConnector{
		Enabled:        biconnector.Enabled,
		ReadPreference: biconnector.ReadPreference,
	}
}

func convertConnectionStringToLatest(connStrings *admin20240530.ClusterConnectionStrings) *admin.ClusterConnectionStrings {
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

func convertPrivateEndpointToLatest(privateEndpoints *[]admin20240530.ClusterDescriptionConnectionStringsPrivateEndpoint) *[]admin.ClusterDescriptionConnectionStringsPrivateEndpoint {
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

func convertEndpointsToLatest(privateEndpoints *[]admin20240530.ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) *[]admin.ClusterDescriptionConnectionStringsPrivateEndpointEndpoint {
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

func convertLabelsToLatest(labels *[]admin20240530.ComponentLabel) *[]admin.ComponentLabel {
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

func convertLabelSliceToOldSDK(slice []admin.ComponentLabel, err diag.Diagnostics) ([]admin20240530.ComponentLabel, diag.Diagnostics) {
	if err != nil {
		return nil, err
	}
	results := make([]admin20240530.ComponentLabel, len(slice))
	for i := range len(slice) {
		label := slice[i]
		results[i] = admin20240530.ComponentLabel{
			Key:   label.Key,
			Value: label.Value,
		}
	}
	return results, nil
}

func convertRegionConfigSliceToOldSDK(slice *[]admin.CloudRegionConfig20240805) *[]admin20240530.CloudRegionConfig {
	if slice == nil {
		return nil
	}
	cloudRegionSlice := *slice
	results := make([]admin20240530.CloudRegionConfig, len(cloudRegionSlice))
	for i := range len(cloudRegionSlice) {
		cloudRegion := cloudRegionSlice[i]
		results[i] = admin20240530.CloudRegionConfig{
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

func convertHardwareSpecToOldSDK(hwspec *admin.HardwareSpec20240805) *admin20240530.HardwareSpec {
	if hwspec == nil {
		return nil
	}
	return &admin20240530.HardwareSpec{
		DiskIOPS:      hwspec.DiskIOPS,
		EbsVolumeType: hwspec.EbsVolumeType,
		InstanceSize:  hwspec.InstanceSize,
		NodeCount:     hwspec.NodeCount,
	}
}

func convertAdvancedAutoScalingSettingsToOldSDK(settings *admin.AdvancedAutoScalingSettings) *admin20240530.AdvancedAutoScalingSettings {
	if settings == nil {
		return nil
	}
	return &admin20240530.AdvancedAutoScalingSettings{
		Compute: convertAdvancedComputeAutoScalingToOldSDK(settings.Compute),
		DiskGB:  convertDiskGBAutoScalingToOldSDK(settings.DiskGB),
	}
}

func convertAdvancedComputeAutoScalingToOldSDK(settings *admin.AdvancedComputeAutoScaling) *admin20240530.AdvancedComputeAutoScaling {
	if settings == nil {
		return nil
	}
	return &admin20240530.AdvancedComputeAutoScaling{
		Enabled:          settings.Enabled,
		MaxInstanceSize:  settings.MaxInstanceSize,
		MinInstanceSize:  settings.MinInstanceSize,
		ScaleDownEnabled: settings.ScaleDownEnabled,
	}
}

func convertDiskGBAutoScalingToOldSDK(settings *admin.DiskGBAutoScaling) *admin20240530.DiskGBAutoScaling {
	if settings == nil {
		return nil
	}
	return &admin20240530.DiskGBAutoScaling{
		Enabled: settings.Enabled,
	}
}

func convertDedicatedHardwareSpecToOldSDK(spec *admin.DedicatedHardwareSpec20240805) *admin20240530.DedicatedHardwareSpec {
	if spec == nil {
		return nil
	}
	return &admin20240530.DedicatedHardwareSpec{
		NodeCount:     spec.NodeCount,
		DiskIOPS:      spec.DiskIOPS,
		EbsVolumeType: spec.EbsVolumeType,
		InstanceSize:  spec.InstanceSize,
	}
}

func convertDedicatedHwSpecToLatest(spec *admin20240530.DedicatedHardwareSpec, rootDiskSizeGB float64) *admin.DedicatedHardwareSpec20240805 {
	if spec == nil {
		return nil
	}
	return &admin.DedicatedHardwareSpec20240805{
		NodeCount:     spec.NodeCount,
		DiskIOPS:      spec.DiskIOPS,
		EbsVolumeType: spec.EbsVolumeType,
		InstanceSize:  spec.InstanceSize,
		DiskSizeGB:    &rootDiskSizeGB,
	}
}

func convertAdvancedAutoScalingSettingsToLatest(settings *admin20240530.AdvancedAutoScalingSettings) *admin.AdvancedAutoScalingSettings {
	if settings == nil {
		return nil
	}
	return &admin.AdvancedAutoScalingSettings{
		Compute: convertAdvancedComputeAutoScalingToLatest(settings.Compute),
		DiskGB:  convertDiskGBAutoScalingToLatest(settings.DiskGB),
	}
}

func convertAdvancedComputeAutoScalingToLatest(settings *admin20240530.AdvancedComputeAutoScaling) *admin.AdvancedComputeAutoScaling {
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

func convertDiskGBAutoScalingToLatest(settings *admin20240530.DiskGBAutoScaling) *admin.DiskGBAutoScaling {
	if settings == nil {
		return nil
	}
	return &admin.DiskGBAutoScaling{
		Enabled: settings.Enabled,
	}
}

func convertHardwareSpecToLatest(hwspec *admin20240530.HardwareSpec, rootDiskSizeGB float64) *admin.HardwareSpec20240805 {
	if hwspec == nil {
		return nil
	}
	return &admin.HardwareSpec20240805{
		DiskIOPS:      hwspec.DiskIOPS,
		EbsVolumeType: hwspec.EbsVolumeType,
		InstanceSize:  hwspec.InstanceSize,
		NodeCount:     hwspec.NodeCount,
		DiskSizeGB:    &rootDiskSizeGB,
	}
}

func convertRegionConfigSliceToLatest(slice *[]admin20240530.CloudRegionConfig, rootDiskSizeGB float64) *[]admin.CloudRegionConfig20240805 {
	if slice == nil {
		return nil
	}
	cloudRegionSlice := *slice
	results := make([]admin.CloudRegionConfig20240805, len(cloudRegionSlice))
	for i := range len(cloudRegionSlice) {
		cloudRegion := cloudRegionSlice[i]
		results[i] = admin.CloudRegionConfig20240805{
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

func convertClusterDescToLatestExcludeRepSpecs(oldClusterDesc *admin20240530.AdvancedClusterDescription) *admin.ClusterDescription20240805 {
	return &admin.ClusterDescription20240805{
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
		Tags:                             convertTagsPtrToLatest(oldClusterDesc.Tags),
		BiConnector:                      convertBiConnectToLatest(oldClusterDesc.BiConnector),
		ConnectionStrings:                convertConnectionStringToLatest(oldClusterDesc.ConnectionStrings),
		Labels:                           convertLabelsToLatest(oldClusterDesc.Labels),
	}
}
