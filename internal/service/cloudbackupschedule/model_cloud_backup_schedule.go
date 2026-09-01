package cloudbackupschedule

import (
	"go.mongodb.org/atlas-sdk/v20250312024/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	export := roles.GetExport()
	if export.FrequencyType != nil || export.ExportBucketId != nil {
		exportList = append(exportList, map[string]any{
			"frequency_type":   export.GetFrequencyType(),
			"export_bucket_id": export.GetExportBucketId(),
		})
	}
	return exportList
}

func FlattenCopySettings(copySettingList []admin.DiskBackupCopySetting20240805) []map[string]any {
	copySettings := make([]map[string]any, 0)
	for _, v := range copySettingList {
		copySetting := map[string]any{
			"cloud_provider":     v.GetCloudProvider(),
			"frequencies":        flattenFrequencies(v.GetFrequencies()),
			"region_name":        v.GetRegionName(),
			"zone_id":            v.GetZoneId(),
			"should_copy_oplogs": v.GetShouldCopyOplogs(),
		}
		if items, ok := v.GetCopyPolicyItemsOk(); ok && items != nil && len(*items) > 0 {
			copySetting["copy_policy_items"] = flattenCopyPolicyItems(*items)
		}
		if lastN, ok := v.GetLastNumberOfSnapshotsOk(); ok && lastN != nil {
			copySetting["last_number_of_snapshots"] = *lastN
		}
		copySettings = append(copySettings, copySetting)
	}
	return copySettings
}

// flattenFrequencies returns a non-nil empty slice when GET omits frequencies so d.Set
// clears Optional+Computed leftover after a switch to copy_policy_items or last-N.
func flattenFrequencies(freqs []string) []string {
	if freqs == nil {
		return []string{}
	}
	return freqs
}

func flattenCopyPolicyItems(items []admin.DiskBackupCopyPolicyItem) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]any{
			"id":              item.GetId(),
			"frequency_type":  item.GetFrequencyType(),
			"retention_unit":  item.GetRetentionUnit(),
			"retention_value": item.GetRetentionValue(),
		})
	}
	return result
}

func ExpandCopyPolicyItems(items []any) *[]admin.DiskBackupCopyPolicyItem {
	if len(items) == 0 {
		return nil
	}
	results := make([]admin.DiskBackupCopyPolicyItem, 0, len(items))
	for _, raw := range items {
		item, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		results = append(results, expandCopyPolicyItem(item))
	}
	if len(results) == 0 {
		return nil
	}
	return &results
}

func expandCopyPolicyItem(item map[string]any) admin.DiskBackupCopyPolicyItem {
	result := admin.DiskBackupCopyPolicyItem{
		FrequencyType: item["frequency_type"].(string),
		Id:            policyItemID(item),
	}
	if unit, ok := item["retention_unit"].(string); ok && unit != "" {
		result.RetentionUnit = new(unit)
	}
	if value, ok := item["retention_value"].(int); ok && value != 0 {
		result.RetentionValue = new(value)
	}
	return result
}

func ExpandCopySetting(tfMap map[string]any) *admin.DiskBackupCopySetting20240805 {
	if tfMap == nil {
		return nil
	}

	var frequencies []string
	if set, ok := tfMap["frequencies"].(*schema.Set); ok && set != nil {
		frequencies = conversion.ExpandStringList(set.List())
	}
	var copyPolicyItems []any
	if items, ok := tfMap["copy_policy_items"].([]any); ok {
		copyPolicyItems = items
	}
	lastN, _ := tfMap["last_number_of_snapshots"].(int)

	copySetting := &admin.DiskBackupCopySetting20240805{
		CloudProvider:    new(tfMap["cloud_provider"].(string)),
		RegionName:       new(tfMap["region_name"].(string)),
		ZoneId:           tfMap["zone_id"].(string),
		ShouldCopyOplogs: new(tfMap["should_copy_oplogs"].(bool)),
	}

	expandedItems := ExpandCopyPolicyItems(copyPolicyItems)
	switch {
	case expandedItems != nil:
		copySetting.CopyPolicyItems = expandedItems
	case lastN > 0:
		copySetting.LastNumberOfSnapshots = new(lastN)
	default:
		copySetting.Frequencies = &frequencies
	}
	return copySetting
}

func ExpandCopySettings(tfList []any) *[]admin.DiskBackupCopySetting20240805 {
	copySettings := make([]admin.DiskBackupCopySetting20240805, 0)

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]any)
		if !ok {
			continue
		}
		apiObject := ExpandCopySetting(tfMap)
		copySettings = append(copySettings, *apiObject)
	}
	return &copySettings
}
