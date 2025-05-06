package cloudbackupschedule_test

import (
	"reflect"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

func TestFlattenPolicyItem(t *testing.T) {
	testCases := []struct {
		name          string
		items         []admin.DiskBackupApiPolicyItem
		frequencyType string
		expected      []map[string]any
	}{
		{
			name: "Matching Frequency Type",
			items: []admin.DiskBackupApiPolicyItem{
				{Id: conversion.StringPtr("1"), FrequencyType: "daily", FrequencyInterval: 1, RetentionUnit: "days", RetentionValue: 30},
				{Id: conversion.StringPtr("2"), FrequencyType: "weekly", FrequencyInterval: 1, RetentionUnit: "weeks", RetentionValue: 52},
				{Id: conversion.StringPtr("3"), FrequencyType: "daily", FrequencyInterval: 2, RetentionUnit: "days", RetentionValue: 60},
			},
			frequencyType: "daily",
			expected: []map[string]any{
				{"id": "1", "frequency_interval": 1, "frequency_type": "daily", "retention_unit": "days", "retention_value": 30},
				{"id": "3", "frequency_interval": 2, "frequency_type": "daily", "retention_unit": "days", "retention_value": 60},
			},
		},
		{
			name: "No Matching Frequency Type",
			items: []admin.DiskBackupApiPolicyItem{
				{Id: conversion.StringPtr("1"), FrequencyType: "weekly", FrequencyInterval: 1, RetentionUnit: "weeks", RetentionValue: 52},
				{Id: conversion.StringPtr("2"), FrequencyType: "monthly", FrequencyInterval: 1, RetentionUnit: "months", RetentionValue: 12},
			},
			frequencyType: "daily",
			expected:      []map[string]any{},
		},
		{
			name:          "Empty input",
			items:         []admin.DiskBackupApiPolicyItem{},
			frequencyType: "daily",
			expected:      []map[string]any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cloudbackupschedule.FlattenPolicyItem(tc.items, tc.frequencyType)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Test %s failed: expected %+v, got %+v", tc.name, tc.expected, result)
			}
		})
	}
}

func TestFlattenExport(t *testing.T) {
	testCases := []struct {
		name     string
		roles    *admin.DiskBackupSnapshotSchedule20240805
		expected []map[string]any
	}{
		{
			name: "Non-empty Export",
			roles: &admin.DiskBackupSnapshotSchedule20240805{
				Export: &admin.AutoExportPolicy{
					FrequencyType:  conversion.StringPtr("daily"),
					ExportBucketId: conversion.StringPtr("bucket123"),
				},
			},
			expected: []map[string]any{
				{"frequency_type": "daily", "export_bucket_id": "bucket123"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cloudbackupschedule.FlattenExport(tc.roles)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Test %s failed: expected %+v, got %+v", tc.name, tc.expected, result)
			}
		})
	}
}

func TestFlattenCopySettings(t *testing.T) {
	testCases := []struct {
		name     string
		settings []admin.DiskBackupCopySetting20240805
		expected []map[string]any
	}{
		{
			name: "Multiple Copy Settings",
			settings: []admin.DiskBackupCopySetting20240805{
				{
					CloudProvider:    conversion.StringPtr("AWS"),
					Frequencies:      &[]string{"daily", "weekly"},
					RegionName:       conversion.StringPtr("US_WEST_1"),
					ZoneId:           "12345",
					ShouldCopyOplogs: conversion.Pointer(true),
				},
				{
					CloudProvider:    conversion.StringPtr("Azure"),
					Frequencies:      &[]string{"monthly"},
					RegionName:       conversion.StringPtr("EAST_US"),
					ZoneId:           "67895",
					ShouldCopyOplogs: conversion.Pointer(false),
				},
			},
			expected: []map[string]any{
				{"cloud_provider": "AWS", "frequencies": []string{"daily", "weekly"}, "region_name": "US_WEST_1", "zone_id": "12345", "should_copy_oplogs": true},
				{"cloud_provider": "Azure", "frequencies": []string{"monthly"}, "region_name": "EAST_US", "zone_id": "67895", "should_copy_oplogs": false},
			},
		},
		{
			name:     "Empty Copy Settings List",
			settings: []admin.DiskBackupCopySetting20240805{},
			expected: []map[string]any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cloudbackupschedule.FlattenCopySettings(tc.settings)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Test %s failed: expected %+v, got %+v", tc.name, tc.expected, result)
			}
		})
	}
}

func TestExpandPolicyItems(t *testing.T) {
	testCases := []struct {
		expected      *[]admin.DiskBackupApiPolicyItem
		name          string
		frequencyType string
		items         []any
	}{
		{
			name: "Valid Input",
			items: []any{
				map[string]any{"id": "123", "retention_unit": "days", "retention_value": 30, "frequency_interval": 1},
				map[string]any{"id": "456", "retention_unit": "weeks", "retention_value": 52, "frequency_interval": 1},
			},
			frequencyType: "monthly",
			expected: &[]admin.DiskBackupApiPolicyItem{
				{Id: conversion.StringPtr("123"), RetentionUnit: "days", RetentionValue: 30, FrequencyInterval: 1, FrequencyType: "monthly"},
				{Id: conversion.StringPtr("456"), RetentionUnit: "weeks", RetentionValue: 52, FrequencyInterval: 1, FrequencyType: "monthly"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cloudbackupschedule.ExpandPolicyItems(tc.items, tc.frequencyType)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Test %s failed: expected %+v, got %+v", tc.name, *tc.expected, *result)
			}
		})
	}
}
