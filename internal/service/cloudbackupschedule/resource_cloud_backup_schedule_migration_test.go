package cloudbackupschedule_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
)

func TestMigBackupRSCloudBackupSchedule_basic(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true})
		useYearly   = mig.IsProviderVersionAtLeast("1.16.0") // attribute introduced in this version
		config      = configNewPolicies(&clusterInfo, &admin20240530.DiskBackupSnapshotSchedule{
			ReferenceHourOfDay:    conversion.Pointer(0),
			ReferenceMinuteOfHour: conversion.Pointer(0),
			RestoreWindowDays:     conversion.Pointer(7),
		}, useYearly)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     mig.PreCheckBasicSleep(t),
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.Name),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "0"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "0"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigBackupRSCloudBackupSchedule_copySettings(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0") // yearly policy item introduced in this version
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{
			CloudBackup: true,
			ReplicationSpecs: []acc.ReplicationSpecRequest{
				{Region: "US_EAST_2"},
			},
			PitEnabled: true, // you cannot copy oplogs when pit is not enabled
		})
		clusterName                     = clusterInfo.Name
		terraformStr                    = clusterInfo.TerraformStr
		clusterResourceName             = clusterInfo.ResourceName
		projectID                       = clusterInfo.ProjectID
		copySettingsConfigWithRepSpecID = configCopySettings(terraformStr, projectID, clusterResourceName, false, true, &admin20240530.DiskBackupSnapshotSchedule{
			ReferenceHourOfDay:    conversion.Pointer(3),
			ReferenceMinuteOfHour: conversion.Pointer(45),
			RestoreWindowDays:     conversion.Pointer(1),
		})
		copySettingsConfigWithZoneID = configCopySettings(terraformStr, projectID, clusterResourceName, false, false, &admin20240530.DiskBackupSnapshotSchedule{
			ReferenceHourOfDay:    conversion.Pointer(3),
			ReferenceMinuteOfHour: conversion.Pointer(45),
			RestoreWindowDays:     conversion.Pointer(1),
		})
		checkMap = map[string]string{
			"cluster_name":                             clusterName,
			"reference_hour_of_day":                    "3",
			"reference_minute_of_hour":                 "45",
			"restore_window_days":                      "1",
			"policy_item_hourly.#":                     "1",
			"policy_item_daily.#":                      "1",
			"policy_item_weekly.#":                     "1",
			"policy_item_monthly.#":                    "1",
			"policy_item_yearly.#":                     "1",
			"policy_item_hourly.0.frequency_interval":  "1",
			"policy_item_hourly.0.retention_unit":      "days",
			"policy_item_hourly.0.retention_value":     "1",
			"policy_item_daily.0.frequency_interval":   "1",
			"policy_item_daily.0.retention_unit":       "days",
			"policy_item_daily.0.retention_value":      "2",
			"policy_item_weekly.0.frequency_interval":  "4",
			"policy_item_weekly.0.retention_unit":      "weeks",
			"policy_item_weekly.0.retention_value":     "3",
			"policy_item_monthly.0.frequency_interval": "5",
			"policy_item_monthly.0.retention_unit":     "months",
			"policy_item_monthly.0.retention_value":    "4",
			"policy_item_yearly.0.frequency_interval":  "1",
			"policy_item_yearly.0.retention_unit":      "years",
			"policy_item_yearly.0.retention_value":     "1",
		}
		copySettingsChecks = map[string]string{
			"copy_settings.#":                    "1",
			"copy_settings.0.cloud_provider":     "AWS",
			"copy_settings.0.region_name":        "US_EAST_1",
			"copy_settings.0.should_copy_oplogs": "true",
		}
	)

	checksDefault := acc.AddAttrChecks(resourceName, []resource.TestCheckFunc{checkExists(resourceName)}, checkMap)
	checksCreate := acc.AddAttrChecks(resourceName, checksDefault, copySettingsChecks)
	checksCreateWithReplicationSpecID := acc.AddAttrSetChecks(resourceName, checksCreate, "copy_settings.0.replication_spec_id")
	checksUpdateWithZoneID := acc.AddAttrSetChecks(resourceName, checksCreate, "copy_settings.0.zone_id")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     mig.PreCheckBasicSleep(t),
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            copySettingsConfigWithRepSpecID,
				Check:             resource.ComposeAggregateTestCheckFunc(checksCreateWithReplicationSpecID...),
			},
			mig.TestStepCheckEmptyPlan(copySettingsConfigWithRepSpecID),
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   copySettingsConfigWithZoneID,
				Check:                    resource.ComposeAggregateTestCheckFunc(checksUpdateWithZoneID...),
			},
			mig.TestStepCheckEmptyPlan(copySettingsConfigWithZoneID),
		},
	})
}
