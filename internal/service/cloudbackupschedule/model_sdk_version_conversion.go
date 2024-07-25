package cloudbackupschedule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20231115 "go.mongodb.org/atlas-sdk/v20231115014/admin"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

// Conversions from one SDK model version to another are used to avoid duplicating our flatten/expand conversion functions.
// - These functions must not contain any business logic.
// - All will be removed once we rely on a single API version.

func convertPolicyItemsToOldSDK(slice *[]admin.DiskBackupApiPolicyItem) []admin20231115.DiskBackupApiPolicyItem {
	if slice == nil {
		return nil
	}
	policyItemsSlice := *slice
	results := make([]admin20231115.DiskBackupApiPolicyItem, len(policyItemsSlice))
	for i := range len(policyItemsSlice) {
		policyItem := policyItemsSlice[i]
		results[i] = admin20231115.DiskBackupApiPolicyItem{
			FrequencyInterval: policyItem.FrequencyInterval,
			FrequencyType:     policyItem.FrequencyType,
			Id:                policyItem.Id,
			RetentionUnit:     policyItem.RetentionUnit,
			RetentionValue:    policyItem.RetentionValue,
		}
	}
	return results
}

func convertPoliciesToLatest(slice *[]admin20231115.AdvancedDiskBackupSnapshotSchedulePolicy) *[]admin.AdvancedDiskBackupSnapshotSchedulePolicy {
	if slice == nil {
		return nil
	}

	policySlice := *slice
	results := make([]admin.AdvancedDiskBackupSnapshotSchedulePolicy, len(policySlice))
	for i := range len(policySlice) {
		policyItem := policySlice[i]
		results[i] = admin.AdvancedDiskBackupSnapshotSchedulePolicy{
			Id:          policyItem.Id,
			PolicyItems: convertPolicyItemsToLatest(policyItem.PolicyItems),
		}
	}
	return &results
}

func convertPolicyItemsToLatest(slice *[]admin20231115.DiskBackupApiPolicyItem) *[]admin.DiskBackupApiPolicyItem {
	if slice == nil {
		return nil
	}
	policyItemsSlice := *slice
	results := make([]admin.DiskBackupApiPolicyItem, len(policyItemsSlice))
	for i := range len(policyItemsSlice) {
		policyItem := policyItemsSlice[i]
		results[i] = admin.DiskBackupApiPolicyItem{
			FrequencyInterval: policyItem.FrequencyInterval,
			FrequencyType:     policyItem.FrequencyType,
			Id:                policyItem.Id,
			RetentionUnit:     policyItem.RetentionUnit,
			RetentionValue:    policyItem.RetentionValue,
		}
	}
	return &results
}

func expandAutoExportPolicy(items []any, d *schema.ResourceData) *admin.AutoExportPolicy {
	itemObj := items[0].(map[string]any)

	if autoExportEnabled := d.Get("auto_export_enabled"); autoExportEnabled != nil && autoExportEnabled.(bool) {
		return &admin.AutoExportPolicy{
			ExportBucketId: conversion.StringPtr(itemObj["export_bucket_id"].(string)),
			FrequencyType:  conversion.StringPtr(itemObj["frequency_type"].(string)),
		}
	}
	return nil
}

func convertAutoExportPolicyToOldSDK(exportPolicy *admin.AutoExportPolicy) *admin20231115.AutoExportPolicy {
	if exportPolicy == nil {
		return nil
	}

	return &admin20231115.AutoExportPolicy{
		ExportBucketId: exportPolicy.ExportBucketId,
		FrequencyType:  exportPolicy.FrequencyType,
	}
}

func convertAutoExportPolicyToLatest(exportPolicy *admin20231115.AutoExportPolicy) *admin.AutoExportPolicy {
	if exportPolicy == nil {
		return nil
	}

	return &admin.AutoExportPolicy{
		ExportBucketId: exportPolicy.ExportBucketId,
		FrequencyType:  exportPolicy.FrequencyType,
	}
}

func getRequestPoliciesOldSDK(policiesItem []admin20231115.DiskBackupApiPolicyItem, respPolicies []admin20231115.AdvancedDiskBackupSnapshotSchedulePolicy) *[]admin20231115.AdvancedDiskBackupSnapshotSchedulePolicy {
	if len(policiesItem) > 0 {
		policy := admin20231115.AdvancedDiskBackupSnapshotSchedulePolicy{
			PolicyItems: &policiesItem,
		}
		if len(respPolicies) == 1 {
			policy.Id = respPolicies[0].Id
		}
		return &[]admin20231115.AdvancedDiskBackupSnapshotSchedulePolicy{policy}
	}
	return nil
}

func getRequestPolicies(policiesItem []admin.DiskBackupApiPolicyItem, respPolicies []admin.AdvancedDiskBackupSnapshotSchedulePolicy) *[]admin.AdvancedDiskBackupSnapshotSchedulePolicy {
	if len(policiesItem) > 0 {
		policy := admin.AdvancedDiskBackupSnapshotSchedulePolicy{
			PolicyItems: &policiesItem,
		}
		if len(respPolicies) == 1 {
			policy.Id = respPolicies[0].Id
		}
		return &[]admin.AdvancedDiskBackupSnapshotSchedulePolicy{policy}
	}
	return nil
}

func convertBackupScheduleReqToOldSDK(req *admin.DiskBackupSnapshotSchedule20250101,
	copySettingsOldSDK *[]admin20231115.DiskBackupCopySetting,
	policiesOldSDK *[]admin20231115.AdvancedDiskBackupSnapshotSchedulePolicy) *admin20231115.DiskBackupSnapshotSchedule {
	return &admin20231115.DiskBackupSnapshotSchedule{
		CopySettings:                      copySettingsOldSDK,
		Policies:                          policiesOldSDK,
		AutoExportEnabled:                 req.AutoExportEnabled,
		Export:                            convertAutoExportPolicyToOldSDK(req.Export),
		UseOrgAndGroupNamesInExportPrefix: req.UseOrgAndGroupNamesInExportPrefix,
		ReferenceHourOfDay:                req.ReferenceHourOfDay,
		ReferenceMinuteOfHour:             req.ReferenceMinuteOfHour,
		RestoreWindowDays:                 req.RestoreWindowDays,
		UpdateSnapshots:                   req.UpdateSnapshots,
	}
}

func convertBackupScheduleToLatestExcludeCopySettings(backupSchedule *admin20231115.DiskBackupSnapshotSchedule) *admin.DiskBackupSnapshotSchedule20250101 {
	return &admin.DiskBackupSnapshotSchedule20250101{
		// CopySettings:                      copySettingsOldSDK,
		Policies:                          convertPoliciesToLatest(backupSchedule.Policies),
		AutoExportEnabled:                 backupSchedule.AutoExportEnabled,
		Export:                            convertAutoExportPolicyToLatest(backupSchedule.Export),
		UseOrgAndGroupNamesInExportPrefix: backupSchedule.UseOrgAndGroupNamesInExportPrefix,
		ReferenceHourOfDay:                backupSchedule.ReferenceHourOfDay,
		ReferenceMinuteOfHour:             backupSchedule.ReferenceMinuteOfHour,
		RestoreWindowDays:                 backupSchedule.RestoreWindowDays,
		UpdateSnapshots:                   backupSchedule.UpdateSnapshots,
	}
}
