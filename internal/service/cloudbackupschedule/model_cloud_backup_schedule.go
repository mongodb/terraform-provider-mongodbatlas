package cloudbackupschedule

import (
	admin20231115 "go.mongodb.org/atlas-sdk/v20231115014/admin"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

func FlattenPolicyItem(items []admin.DiskBackupApiPolicyItem, frequencyType string) []map[string]any {
	policyItems := make([]map[string]any, 0)
	for _, v := range items {
		if frequencyType == v.GetFrequencyType() {
			policyItems = append(policyItems, map[string]any{
				"id":                 v.GetId(),
				"frequency_interval": v.GetFrequencyInterval(),
				"frequency_type":     v.GetFrequencyType(),
				"retention_unit":     v.GetRetentionUnit(),
				"retention_value":    v.GetRetentionValue(),
			})
		}
	}
	return policyItems
}

func FlattenExport(roles *admin.DiskBackupSnapshotSchedule20250101) []map[string]any {
	exportList := make([]map[string]any, 0)
	emptyStruct := admin.DiskBackupSnapshotSchedule20250101{}
	if emptyStruct.GetExport() != roles.GetExport() {
		exportList = append(exportList, map[string]any{
			"frequency_type":   roles.Export.GetFrequencyType(),
			"export_bucket_id": roles.Export.GetExportBucketId(),
		})
	}
	return exportList
}

func flattenCopySettingsOldSDK(copySettingList []admin20231115.DiskBackupCopySetting) []map[string]any {
	copySettings := make([]map[string]any, 0)
	for _, v := range copySettingList {
		copySettings = append(copySettings, map[string]any{
			"cloud_provider":      v.GetCloudProvider(),
			"frequencies":         v.GetFrequencies(),
			"region_name":         v.GetRegionName(),
			"replication_spec_id": v.GetReplicationSpecId(),
			"should_copy_oplogs":  v.GetShouldCopyOplogs(),
		})
	}
	return copySettings
}

func FlattenCopySettings(copySettingList []admin.DiskBackupCopySetting20250101) []map[string]any {
	copySettings := make([]map[string]any, 0)
	for _, v := range copySettingList {
		copySettings = append(copySettings, map[string]any{
			"cloud_provider":     v.GetCloudProvider(),
			"frequencies":        v.GetFrequencies(),
			"region_name":        v.GetRegionName(),
			"zone_id":            v.GetZoneId(),
			"should_copy_oplogs": v.GetShouldCopyOplogs(),
		})
	}
	return copySettings
}
