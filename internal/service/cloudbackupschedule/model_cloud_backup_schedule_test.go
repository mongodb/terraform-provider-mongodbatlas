package cloudbackupschedule_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
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
					ShouldCopyOplogs: new(true),
				},
				{
					CloudProvider:    conversion.StringPtr("Azure"),
					Frequencies:      &[]string{"monthly"},
					RegionName:       conversion.StringPtr("EAST_US"),
					ZoneId:           "67895",
					ShouldCopyOplogs: new(false),
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
		{
			name: "Copy policy items with id",
			settings: []admin.DiskBackupCopySetting20240805{
				{
					CloudProvider:    conversion.StringPtr("AWS"),
					RegionName:       conversion.StringPtr("US_WEST_1"),
					ZoneId:           "12345",
					ShouldCopyOplogs: new(true),
					CopyPolicyItems: &[]admin.DiskBackupCopyPolicyItem{
						{
							FrequencyType:  "daily",
							Id:             conversion.StringPtr("item-1"),
							RetentionUnit:  conversion.StringPtr("days"),
							RetentionValue: new(7),
						},
					},
				},
			},
			expected: []map[string]any{
				{
					"cloud_provider":     "AWS",
					"frequencies":        []string(nil),
					"region_name":        "US_WEST_1",
					"zone_id":            "12345",
					"should_copy_oplogs": true,
					"copy_policy_items": []map[string]any{
						{"id": "item-1", "frequency_type": "daily", "retention_unit": "days", "retention_value": 7},
					},
				},
			},
		},
		{
			name: "Copy policy items without id",
			settings: []admin.DiskBackupCopySetting20240805{
				{
					CloudProvider:    conversion.StringPtr("AWS"),
					RegionName:       conversion.StringPtr("US_WEST_1"),
					ZoneId:           "12345",
					ShouldCopyOplogs: new(false),
					CopyPolicyItems: &[]admin.DiskBackupCopyPolicyItem{
						{FrequencyType: "ondemand"},
					},
				},
			},
			expected: []map[string]any{
				{
					"cloud_provider":     "AWS",
					"frequencies":        []string(nil),
					"region_name":        "US_WEST_1",
					"zone_id":            "12345",
					"should_copy_oplogs": false,
					"copy_policy_items": []map[string]any{
						{"id": "", "frequency_type": "ondemand", "retention_unit": "", "retention_value": 0},
					},
				},
			},
		},
		{
			name: "Last number of snapshots",
			settings: []admin.DiskBackupCopySetting20240805{
				{
					CloudProvider:         conversion.StringPtr("AWS"),
					RegionName:            conversion.StringPtr("US_WEST_1"),
					ZoneId:                "12345",
					ShouldCopyOplogs:      new(false),
					LastNumberOfSnapshots: new(5),
				},
			},
			expected: []map[string]any{
				{
					"cloud_provider":           "AWS",
					"frequencies":              []string(nil),
					"region_name":              "US_WEST_1",
					"zone_id":                  "12345",
					"should_copy_oplogs":       false,
					"last_number_of_snapshots": 5,
				},
			},
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

func TestExpandCopySetting(t *testing.T) {
	testCases := []struct {
		expected *admin.DiskBackupCopySetting20240805
		tfMap    map[string]any
		name     string
	}{
		{
			name: "Frequencies only",
			tfMap: map[string]any{
				"cloud_provider":     "AWS",
				"region_name":        "US_WEST_1",
				"zone_id":            "12345",
				"should_copy_oplogs": true,
				"frequencies":        testFreqSet("DAILY"),
			},
			expected: &admin.DiskBackupCopySetting20240805{
				CloudProvider:    new("AWS"),
				RegionName:       new("US_WEST_1"),
				ZoneId:           "12345",
				ShouldCopyOplogs: new(true),
				Frequencies:      &[]string{"DAILY"},
			},
		},
		{
			name: "Copy policy items with id",
			tfMap: map[string]any{
				"cloud_provider":     "AWS",
				"region_name":        "US_WEST_1",
				"zone_id":            "12345",
				"should_copy_oplogs": true,
				"frequencies":        testFreqSet(),
				"copy_policy_items": []any{
					map[string]any{
						"id":              "item-1",
						"frequency_type":  "daily",
						"retention_unit":  "days",
						"retention_value": 7,
					},
				},
			},
			expected: &admin.DiskBackupCopySetting20240805{
				CloudProvider:    new("AWS"),
				RegionName:       new("US_WEST_1"),
				ZoneId:           "12345",
				ShouldCopyOplogs: new(true),
				CopyPolicyItems: &[]admin.DiskBackupCopyPolicyItem{
					{
						Id:             conversion.StringPtr("item-1"),
						FrequencyType:  "daily",
						RetentionUnit:  conversion.StringPtr("days"),
						RetentionValue: new(7),
					},
				},
			},
		},
		{
			name: "Copy policy items without id",
			tfMap: map[string]any{
				"cloud_provider":     "AWS",
				"region_name":        "US_WEST_1",
				"zone_id":            "12345",
				"should_copy_oplogs": false,
				"copy_policy_items": []any{
					map[string]any{"frequency_type": "ondemand"},
				},
			},
			expected: &admin.DiskBackupCopySetting20240805{
				CloudProvider:    new("AWS"),
				RegionName:       new("US_WEST_1"),
				ZoneId:           "12345",
				ShouldCopyOplogs: new(false),
				CopyPolicyItems: &[]admin.DiskBackupCopyPolicyItem{
					{FrequencyType: "ondemand"},
				},
			},
		},
		{
			name: "Last number of snapshots",
			tfMap: map[string]any{
				"cloud_provider":           "AWS",
				"region_name":              "US_WEST_1",
				"zone_id":                  "12345",
				"should_copy_oplogs":       false,
				"last_number_of_snapshots": 5,
				"frequencies":              testFreqSet(),
			},
			expected: &admin.DiskBackupCopySetting20240805{
				CloudProvider:         new("AWS"),
				RegionName:            new("US_WEST_1"),
				ZoneId:                "12345",
				ShouldCopyOplogs:      new(false),
				LastNumberOfSnapshots: new(5),
			},
		},
		{
			name: "Empty copy policy items list uses frequencies",
			tfMap: map[string]any{
				"cloud_provider":     "AWS",
				"region_name":        "US_WEST_1",
				"zone_id":            "12345",
				"should_copy_oplogs": false,
				"frequencies":        testFreqSet("HOURLY"),
				"copy_policy_items":  []any{},
			},
			expected: &admin.DiskBackupCopySetting20240805{
				CloudProvider:    new("AWS"),
				RegionName:       new("US_WEST_1"),
				ZoneId:           "12345",
				ShouldCopyOplogs: new(false),
				Frequencies:      &[]string{"HOURLY"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cloudbackupschedule.ExpandCopySetting(tc.tfMap)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Test %s failed: expected %+v, got %+v", tc.name, tc.expected, result)
			}
			assertUnusedCopySettingModes(t, result)
		})
	}
}

func assertUnusedCopySettingModes(t *testing.T, setting *admin.DiskBackupCopySetting20240805) {
	t.Helper()
	modes := 0
	if setting.Frequencies != nil {
		modes++
	}
	if setting.CopyPolicyItems != nil {
		modes++
	}
	if setting.LastNumberOfSnapshots != nil {
		modes++
	}
	if modes != 1 {
		t.Errorf("expected exactly one copy mode pointer to be set, got frequencies=%v copyPolicyItems=%v lastNumberOfSnapshots=%v", setting.Frequencies != nil, setting.CopyPolicyItems != nil, setting.LastNumberOfSnapshots != nil)
	}
}

func testFreqSet(items ...string) *schema.Set {
	raw := make([]any, len(items))
	for i, item := range items {
		raw[i] = item
	}
	return schema.NewSet(schema.HashString, raw)
}
