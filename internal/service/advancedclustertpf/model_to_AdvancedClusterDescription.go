package advancedclustertpf

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20241113002/admin"
)

func newLegacyModel20240530ReplicationSpecsAndDiskGBOnly(specs *[]admin.ReplicationSpec20240805, zoneNameNumShards map[string]int64, oldDiskGB *float64) *admin20240530.AdvancedClusterDescription {
	newDiskGB := findRegionRootDiskSize(specs)
	if oldDiskGB != nil && newDiskGB != nil && (*newDiskGB-*oldDiskGB) < 0.01 {
		newDiskGB = nil
	}
	return &admin20240530.AdvancedClusterDescription{
		DiskSizeGB:       newDiskGB,
		ReplicationSpecs: convertReplicationSpecs20240805to20240530(specs, zoneNameNumShards),
	}
}

func convertReplicationSpecs20240805to20240530(replicationSpecs *[]admin.ReplicationSpec20240805, zoneNameNumShards map[string]int64) *[]admin20240530.ReplicationSpec {
	if replicationSpecs == nil {
		return nil
	}
	result := make([]admin20240530.ReplicationSpec, len(*replicationSpecs))
	for i, replicationSpec := range *replicationSpecs {
		numShards, ok := zoneNameNumShards[replicationSpec.GetZoneName()]
		if !ok {
			numShards = 1
		}
		result[i] = admin20240530.ReplicationSpec{
			NumShards:     conversion.Int64PtrToIntPtr(&numShards),
			Id:            replicationSpec.Id,
			ZoneName:      replicationSpec.ZoneName,
			RegionConfigs: convertCloudRegionConfig20240805to20240530(replicationSpec.RegionConfigs),
		}
	}
	return &result
}

func convertCloudRegionConfig20240805to20240530(cloudRegionConfig *[]admin.CloudRegionConfig20240805) *[]admin20240530.CloudRegionConfig {
	if cloudRegionConfig == nil {
		return nil
	}
	result := make([]admin20240530.CloudRegionConfig, len(*cloudRegionConfig))
	for i, regionConfig := range *cloudRegionConfig {
		result[i] = admin20240530.CloudRegionConfig{
			ProviderName:         regionConfig.ProviderName,
			RegionName:           regionConfig.RegionName,
			BackingProviderName:  regionConfig.BackingProviderName,
			Priority:             regionConfig.Priority,
			ElectableSpecs:       convertHardwareSpecToOldSDK(regionConfig.ElectableSpecs),
			ReadOnlySpecs:        convertDedicatedHardwareSpecToOldSDK(regionConfig.ReadOnlySpecs),
			AnalyticsSpecs:       convertDedicatedHardwareSpecToOldSDK(regionConfig.AnalyticsSpecs),
			AutoScaling:          convertAdvancedAutoScalingSettingsToOldSDK(regionConfig.AutoScaling),
			AnalyticsAutoScaling: convertAdvancedAutoScalingSettingsToOldSDK(regionConfig.AnalyticsAutoScaling),
		}
	}
	return &result
}
