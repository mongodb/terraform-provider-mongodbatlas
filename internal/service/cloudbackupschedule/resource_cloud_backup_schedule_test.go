package cloudbackupschedule_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115014/admin"
)

var (
	resourceName   = "mongodbatlas_cloud_backup_schedule.schedule_test"
	dataSourceName = "data.mongodbatlas_cloud_backup_schedule.schedule_test"
)

func TestAccBackupRSCloudBackupSchedule_basic(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configNoPolicies(&clusterInfo, &admin.DiskBackupSnapshotSchedule{
					ReferenceHourOfDay:    conversion.Pointer(3),
					ReferenceMinuteOfHour: conversion.Pointer(45),
					RestoreWindowDays:     conversion.Pointer(4),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "3"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "45"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttrSet(dataSourceName, "reference_hour_of_day"),
					resource.TestCheckResourceAttrSet(dataSourceName, "reference_minute_of_hour"),
					resource.TestCheckResourceAttrSet(dataSourceName, "restore_window_days"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_item_hourly.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_item_daily.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_item_weekly.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_item_monthly.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_item_yearly.#"),
				),
			},
			{
				Config: configNewPolicies(&clusterInfo, &admin.DiskBackupSnapshotSchedule{
					ReferenceHourOfDay:    conversion.Pointer(0),
					ReferenceMinuteOfHour: conversion.Pointer(0),
					RestoreWindowDays:     conversion.Pointer(7),
				}, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "0"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "0"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_hourly.0.id"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_daily.0.id"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_weekly.0.id"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.frequency_interval", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_value", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_monthly.0.id"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_value", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_yearly.0.id"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_unit", "years"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttrSet(dataSourceName, "reference_hour_of_day"),
					resource.TestCheckResourceAttrSet(dataSourceName, "reference_minute_of_hour"),
					resource.TestCheckResourceAttrSet(dataSourceName, "restore_window_days"),
				),
			},
			{
				Config: configAdvancedPolicies(&clusterInfo, &admin.DiskBackupSnapshotSchedule{
					ReferenceHourOfDay:    conversion.Pointer(0),
					ReferenceMinuteOfHour: conversion.Pointer(0),
					RestoreWindowDays:     conversion.Pointer(7),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "auto_export_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "0"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "0"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_hourly.0.id"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_daily.0.id"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_weekly.0.id"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.frequency_interval", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_value", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_monthly.0.id"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_value", "3"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.1.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.1.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.1.retention_value", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.1.frequency_interval", "6"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.1.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.1.retention_value", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_unit", "years"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_value", "1"),
				),
			},
		},
	})
}

func TestAccBackupRSCloudBackupSchedule_export(t *testing.T) {
	var (
		// A snapshot export bucket can't be deleted it there exist a cluster that is still using it. So the cluster resource needs to depend on it
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true, ResourceDependencyName: "mongodbatlas_cloud_backup_snapshot_export_bucket.test"})
		policyName  = acc.RandomName()
		roleName    = acc.RandomIAMRole()
		bucketName  = acc.RandomS3BucketName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,

		Steps: []resource.TestStep{
			{
				Config: configExportPolicies(&clusterInfo, policyName, roleName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "auto_export_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "20"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "5"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "4"),
				),
			},
		},
	})
}
func TestAccBackupRSCloudBackupSchedule_onePolicy(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDefault(&clusterInfo, &admin.DiskBackupSnapshotSchedule{
					ReferenceHourOfDay:    conversion.Pointer(3),
					ReferenceMinuteOfHour: conversion.Pointer(45),
					RestoreWindowDays:     conversion.Pointer(4),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "3"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "45"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.frequency_interval", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_value", "3"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_value", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_unit", "years"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_value", "1"),
				),
			},
			{
				Config: configOnePolicy(&clusterInfo, &admin.DiskBackupSnapshotSchedule{
					ReferenceHourOfDay:    conversion.Pointer(0),
					ReferenceMinuteOfHour: conversion.Pointer(0),
					RestoreWindowDays:     conversion.Pointer(7),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "0"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "0"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.#", "0"),
				),
			},
		},
	})
}

func TestAccBackupRSCloudBackupSchedule_copySettings(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configCopySettings(projectID, clusterName, &admin.DiskBackupSnapshotSchedule{
					ReferenceHourOfDay:    conversion.Pointer(3),
					ReferenceMinuteOfHour: conversion.Pointer(45),
					RestoreWindowDays:     conversion.Pointer(1),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "3"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "45"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.frequency_interval", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_value", "3"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_value", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_unit", "years"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "copy_settings.0.cloud_provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "copy_settings.0.region_name", "US_EAST_1"),
					resource.TestCheckResourceAttr(resourceName, "copy_settings.0.should_copy_oplogs", "true"),
				),
			},
		},
	})
}
func TestAccBackupRSCloudBackupScheduleImport_basic(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDefault(&clusterInfo, &admin.DiskBackupSnapshotSchedule{
					ReferenceHourOfDay:    conversion.Pointer(3),
					ReferenceMinuteOfHour: conversion.Pointer(45),
					RestoreWindowDays:     conversion.Pointer(4),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "3"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "45"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.frequency_interval", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_value", "3"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_value", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_unit", "years"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_yearly.0.retention_value", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBackupRSCloudBackupSchedule_azure(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true, ProviderName: constant.AZURE})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAzure(&clusterInfo, &admin.DiskBackupApiPolicyItem{
					FrequencyInterval: 1,
					RetentionUnit:     "days",
					RetentionValue:    1,
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1")),
			},
			{
				Config: configAzure(&clusterInfo, &admin.DiskBackupApiPolicyItem{
					FrequencyInterval: 2,
					RetentionUnit:     "days",
					RetentionValue:    3,
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]
		clusterName := ids["cluster_name"]
		_, _, err := acc.ConnV2().CloudBackupsApi.GetBackupSchedule(context.Background(), projectID, clusterName).Execute()
		if err != nil {
			return fmt.Errorf("cloud Provider Snapshot Schedule (%s) does not exist: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	if acc.ExistingClusterUsed() {
		return nil
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_schedule" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]
		clusterName := ids["cluster_name"]
		_, _, err := acc.ConnV2().CloudBackupsApi.GetBackupSchedule(context.Background(), projectID, clusterName).Execute()
		if err == nil {
			return fmt.Errorf("cloud Provider Snapshot Schedule (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func configNoPolicies(info *acc.ClusterInfo, p *admin.DiskBackupSnapshotSchedule) string {
	return info.ClusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s

			reference_hour_of_day    = %[3]d
			reference_minute_of_hour = %[4]d
			restore_window_days      = %[5]d
		}

		data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
		}	
	`, info.ClusterNameStr, info.ProjectIDStr, p.GetReferenceHourOfDay(), p.GetReferenceMinuteOfHour(), p.GetRestoreWindowDays())
}

func configDefault(info *acc.ClusterInfo, p *admin.DiskBackupSnapshotSchedule) string {
	return info.ClusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s

			reference_hour_of_day    = %[3]d
			reference_minute_of_hour = %[4]d
			restore_window_days      = %[5]d

			policy_item_hourly {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 1
			}
			policy_item_daily {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 2
			}
			policy_item_weekly {
				frequency_interval = 4
				retention_unit     = "weeks"
				retention_value    = 3
			}
			policy_item_monthly {
				frequency_interval = 5
				retention_unit     = "months"
				retention_value    = 4
			}
			policy_item_yearly {
				frequency_interval = 1
				retention_unit     = "years"
				retention_value    = 1
			}
		}

		data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
		 }	
	`, info.ClusterNameStr, info.ProjectIDStr, p.GetReferenceHourOfDay(), p.GetReferenceMinuteOfHour(), p.GetRestoreWindowDays())
}

func configCopySettings(projectID, clusterName string, p *admin.DiskBackupSnapshotSchedule) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = %[1]q
			name         = %[2]q
			
			cluster_type = "REPLICASET"
						replication_specs {
						num_shards = 1
						regions_config {
							region_name     = "US_EAST_2"
							electable_nodes = 3
							priority        = 7
							read_only_nodes = 0
							}
						}
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_2"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
			pit_enabled = true // enable point in time restore. you cannot copy oplogs when pit is not enabled.
		}
		
		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = %[1]q
			cluster_name     = %[2]q

			reference_hour_of_day    = %[3]d
			reference_minute_of_hour = %[4]d
			restore_window_days      = %[5]d

			policy_item_hourly {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 1
			}
			policy_item_daily {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 2
			}
			policy_item_weekly {
				frequency_interval = 4
				retention_unit     = "weeks"
				retention_value    = 3
			}
			policy_item_monthly {
				frequency_interval = 5
				retention_unit     = "months"
				retention_value    = 4
			}
			policy_item_yearly {
				frequency_interval = 1
				retention_unit     = "years"
				retention_value    = 1
			}
			copy_settings {
				cloud_provider = "AWS"
				frequencies = ["HOURLY",
							"DAILY",
							"WEEKLY",
							"MONTHLY",
							"YEARLY",
							"ON_DEMAND"]
				region_name = "US_EAST_1"
				replication_spec_id = mongodbatlas_cluster.my_cluster.replication_specs.*.id[0]
				should_copy_oplogs = true
			}
		}
	`, projectID, clusterName, p.GetReferenceHourOfDay(), p.GetReferenceMinuteOfHour(), p.GetRestoreWindowDays())
}

func configOnePolicy(info *acc.ClusterInfo, p *admin.DiskBackupSnapshotSchedule) string {
	return info.ClusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s

			reference_hour_of_day    = %[3]d
			reference_minute_of_hour = %[4]d
			restore_window_days      = %[5]d

			policy_item_hourly {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 1
			}
		}
	`, info.ClusterNameStr, info.ProjectIDStr, p.GetReferenceHourOfDay(), p.GetReferenceMinuteOfHour(), p.GetRestoreWindowDays())
}

func configNewPolicies(info *acc.ClusterInfo, p *admin.DiskBackupSnapshotSchedule, useYearly bool) string {
	var strYearly string
	if useYearly {
		strYearly = `
			policy_item_yearly {
				frequency_interval = 1
				retention_unit     = "years"
				retention_value    = 1
			}
		`
	}

	return info.ClusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s

			reference_hour_of_day    = %[3]d
			reference_minute_of_hour = %[4]d
			restore_window_days      = %[5]d

			policy_item_hourly {
				frequency_interval = 2
				retention_unit     = "days"
				retention_value    = 1
			}
			policy_item_daily {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 4
			}
			policy_item_weekly {
				frequency_interval = 4
				retention_unit     = "weeks"
				retention_value    = 2
			}
			policy_item_monthly {
				frequency_interval = 5
				retention_unit     = "months"
				retention_value    = 3
			}
			%[6]s
		}

		data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
		 }	
	`, info.ClusterNameStr, info.ProjectIDStr, p.GetReferenceHourOfDay(), p.GetReferenceMinuteOfHour(), p.GetRestoreWindowDays(), strYearly)
}

func configAzure(info *acc.ClusterInfo, policy *admin.DiskBackupApiPolicyItem) string {
	return info.ClusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s

			policy_item_hourly {
				frequency_interval = %[3]d
				retention_unit     = %[4]q
				retention_value    = %[5]d
			}
		}

		data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
		}	
	`, info.ClusterNameStr, info.ProjectIDStr, policy.GetFrequencyInterval(), policy.GetRetentionUnit(), policy.GetRetentionValue())
}

func configAdvancedPolicies(info *acc.ClusterInfo, p *admin.DiskBackupSnapshotSchedule) string {
	return info.ClusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			cluster_name     = %[1]s
			project_id       = %[2]s

			auto_export_enabled = false
			reference_hour_of_day    = %[3]d
			reference_minute_of_hour = %[4]d
			restore_window_days      = %[5]d

			policy_item_hourly {
				frequency_interval = 2
				retention_unit     = "days"
				retention_value    = 1
			}
			policy_item_daily {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 4
			}
			policy_item_weekly {
				frequency_interval = 4
				retention_unit     = "weeks"
				retention_value    = 2
			}
			policy_item_weekly {
				frequency_interval = 5
				retention_unit     = "weeks"
				retention_value    = 5
			}
			policy_item_monthly {
				frequency_interval = 5
				retention_unit     = "months"
				retention_value    = 3
			}
			policy_item_monthly {
				frequency_interval = 6
				retention_unit     = "months"
				retention_value    = 4
			}
			policy_item_yearly {
				frequency_interval = 1
				retention_unit     = "years"
				retention_value    = 1
			}
		}
	`, info.ClusterNameStr, info.ProjectIDStr, p.GetReferenceHourOfDay(), p.GetReferenceMinuteOfHour(), p.GetRestoreWindowDays())
}

func configExportPolicies(info *acc.ClusterInfo, policyName, roleName, bucketName string) string {
	return info.ClusterTerraformStr + fmt.Sprintf(`
    resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
        cluster_name             = %[1]s
        project_id               = %[2]s
        auto_export_enabled      = true
        reference_hour_of_day    = 20
        reference_minute_of_hour = "05"
        restore_window_days      = 4

        policy_item_hourly {
            frequency_interval = 1 #accepted values = 1, 2, 4, 6, 8, 12 -> every n hours
            retention_unit     = "days"
            retention_value    = 4
        }		
        policy_item_daily {
            frequency_interval = 1
            retention_unit     = "days"
            retention_value    = 4
        }
        policy_item_weekly {
            frequency_interval = 4        # accepted values = 1 to 7 -> every 1=Monday,2=Tuesday,3=Wednesday,4=Thursday,5=Friday,6=Saturday,7=Sunday day of the week
            retention_unit     = "weeks"
            retention_value    = 4
        }
        policy_item_monthly {
            frequency_interval = 5        # accepted values = 1 to 28 -> 1 to 28 every nth day of the month  
        	                              # accepted values = 40 -> every last day of the month
            retention_unit     = "months"
            retention_value    = 4
        }  		

        export {
            export_bucket_id = mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id
            frequency_type   = "monthly"
        }
    }

    resource "aws_s3_bucket" "backup" {
        bucket          = %[5]q
        force_destroy   = true
    }

    resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
        project_id      = %[2]s
        provider_name   = "AWS"
    }

    resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
        project_id  = %[2]s
        role_id     = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
        aws {
            iam_assumed_role_arn = aws_iam_role.test_role.arn
        }
    }

    resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
        project_id     = %[2]s
        iam_role_id    = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
        bucket_name    = aws_s3_bucket.backup.bucket
        cloud_provider = "AWS"
    }

    resource "aws_iam_role_policy" "test_policy" {
        name = %[3]q
        role = aws_iam_role.test_role.id
        policy = <<-EOF
        {
            "Version": "2012-10-17",
            "Statement": [
            {
                "Effect": "Allow",
                "Action": "s3:GetBucketLocation",
                "Resource": "arn:aws:s3:::%[5]s"
            },
            {
                "Effect": "Allow",
                "Action": "s3:PutObject",
                "Resource": "arn:aws:s3:::%[5]s/*"
            }]
        }
        EOF
    }

    resource "aws_iam_role" "test_role" {
        name = %[4]q
        assume_role_policy = <<EOF
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "AWS": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config.0.atlas_aws_account_arn}"
                    },
                    "Action": "sts:AssumeRole",
                    "Condition": {
                        "StringEquals": {
                            "sts:ExternalId": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config.0.atlas_assumed_role_external_id}"
                        }
                    }
                }
            ]
        }
    EOF
    }
    `, info.ClusterNameStr, info.ProjectIDStr, policyName, roleName, bucketName)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["cluster_name"]), nil
	}
}
