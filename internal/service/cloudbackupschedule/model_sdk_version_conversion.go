package cloudbackupschedule

import (
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

// Conversions from one SDK model version to another are used to avoid duplicating our flatten/expand conversion functions.
// - These functions must not contain any business logic.
// - All will be removed once we rely on a single API version.

func convertPolicyItemsToOldSDK(slice *[]admin.DiskBackupApiPolicyItem) []admin20240530.DiskBackupApiPolicyItem {
	if slice == nil {
		return nil
	}
	policyItemsSlice := *slice
	results := make([]admin20240530.DiskBackupApiPolicyItem, len(policyItemsSlice))
	for i := range len(policyItemsSlice) {
		policyItem := policyItemsSlice[i]
		results[i] = admin20240530.DiskBackupApiPolicyItem{
			FrequencyInterval: policyItem.FrequencyInterval,
			FrequencyType:     policyItem.FrequencyType,
			Id:                policyItem.Id,
			RetentionUnit:     policyItem.RetentionUnit,
			RetentionValue:    policyItem.RetentionValue,
		}
	}
	return results
}

func convertPoliciesToLatest(slice *[]admin20240530.AdvancedDiskBackupSnapshotSchedulePolicy) *[]admin.AdvancedDiskBackupSnapshotSchedulePolicy {
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

func convertPolicyItemsToLatest(slice *[]admin20240530.DiskBackupApiPolicyItem) *[]admin.DiskBackupApiPolicyItem {
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

func convertAutoExportPolicyToOldSDK(exportPolicy *admin.AutoExportPolicy) *admin20240530.AutoExportPolicy {
	if exportPolicy == nil {
		return nil
	}

	return &admin20240530.AutoExportPolicy{
		ExportBucketId: exportPolicy.ExportBucketId,
		FrequencyType:  exportPolicy.FrequencyType,
	}
}

func convertAutoExportPolicyToLatest(exportPolicy *admin20240530.AutoExportPolicy) *admin.AutoExportPolicy {
	if exportPolicy == nil {
		return nil
	}

	return &admin.AutoExportPolicy{
		ExportBucketId: exportPolicy.ExportBucketId,
		FrequencyType:  exportPolicy.FrequencyType,
	}
}

func convertBackupScheduleReqToOldSDK(req *admin.DiskBackupSnapshotSchedule20240805,
	copySettingsOldSDK *[]admin20240530.DiskBackupCopySetting,
	policiesOldSDK *[]admin20240530.AdvancedDiskBackupSnapshotSchedulePolicy) *admin20240530.DiskBackupSnapshotSchedule {
	return &admin20240530.DiskBackupSnapshotSchedule{
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

func convertBackupScheduleToLatestExcludeCopySettings(backupSchedule *admin20240530.DiskBackupSnapshotSchedule) *admin.DiskBackupSnapshotSchedule20240805 {
	return &admin.DiskBackupSnapshotSchedule20240805{
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
