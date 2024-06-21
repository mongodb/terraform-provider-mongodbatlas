package advancedcluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	admin20231115 "go.mongodb.org/atlas-sdk/v20231115014/admin"
	"go.mongodb.org/atlas-sdk/v20240530001/admin"
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
