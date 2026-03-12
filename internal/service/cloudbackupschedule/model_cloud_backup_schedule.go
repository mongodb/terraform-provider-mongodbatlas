package cloudbackupschedule

import (
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
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

func FlattenExport(roles *admin.DiskBackupSnapshotSchedule20240805) []map[string]any {
	exportList := make([]map[string]any, 0)
	emptyStruct := admin.DiskBackupSnapshotSchedule20240805{}
	if emptyStruct.GetExport() != roles.GetExport() {
		exportList = append(exportList, map[string]any{
			"frequency_type":   roles.Export.GetFrequencyType(),
			"export_bucket_id": roles.Export.GetExportBucketId(),
		})
	}
	return exportList
}

func FlattenCopySettings(copySettingList []admin.DiskBackupCopySetting20240805) []map[string]any {
	copySettings := make([]map[string]any, 0)
	for _, v := range copySettingList {
		setting := map[string]any{
			"cloud_provider":     v.GetCloudProvider(),
			"frequencies":        v.GetFrequencies(),
			"region_name":        v.GetRegionName(),
			"zone_id":            v.GetZoneId(),
			"should_copy_oplogs": v.GetShouldCopyOplogs(),
		}
		if copyPolicyItems, ok := v.GetCopyPolicyItemsOk(); ok && copyPolicyItems != nil {
			setting["copy_policy_items"] = FlattenCopyPolicyItems(*copyPolicyItems)
		}
		if lastNumberOfSnapshots, ok := v.GetLastNumberOfSnapshotsOk(); ok && lastNumberOfSnapshots != nil {
			setting["last_number_of_snapshots"] = *lastNumberOfSnapshots
		}
		copySettings = append(copySettings, setting)
	}
	return copySettings
}

func FlattenCopyPolicyItems(items []admin.DiskBackupCopyPolicyItem) []map[string]any {
	policyItems := make([]map[string]any, 0)
	for _, item := range items {
		policyItem := map[string]any{
			"frequency_type": item.GetFrequencyType(),
		}
		if id, ok := item.GetIdOk(); ok && id != nil {
			policyItem["id"] = *id
		}
		if retentionUnit, ok := item.GetRetentionUnitOk(); ok && retentionUnit != nil {
			policyItem["retention_unit"] = *retentionUnit
		}
		if retentionValue, ok := item.GetRetentionValueOk(); ok && retentionValue != nil {
			policyItem["retention_value"] = *retentionValue
		}
		policyItems = append(policyItems, policyItem)
	}
	return policyItems
}

func ExpandCopyPolicyItems(tfList []any) *[]admin.DiskBackupCopyPolicyItem {
	if len(tfList) == 0 {
		return nil
	}

	policyItems := make([]admin.DiskBackupCopyPolicyItem, 0, len(tfList))
	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]any)
		if !ok {
			continue
		}

		policyItem := admin.DiskBackupCopyPolicyItem{
			FrequencyType: tfMap["frequency_type"].(string),
		}

		if retentionUnit, ok := tfMap["retention_unit"]; ok && retentionUnit.(string) != "" {
			ru := retentionUnit.(string)
			policyItem.RetentionUnit = &ru
		}

		if retentionValue, ok := tfMap["retention_value"]; ok && retentionValue.(int) > 0 {
			rv := retentionValue.(int)
			policyItem.RetentionValue = &rv
		}

		policyItems = append(policyItems, policyItem)
	}

	return &policyItems
}
