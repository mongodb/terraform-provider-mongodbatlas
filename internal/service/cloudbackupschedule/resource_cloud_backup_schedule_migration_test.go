package cloudbackupschedule_test

import (
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupRSCloudBackupSchedule_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.29.0") // version when advanced cluster TPF was introduced
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true})
		useYearly   = mig.IsProviderVersionAtLeast("1.16.0") // attribute introduced in this version
		config      = configNewPolicies(&clusterInfo, &admin.DiskBackupSnapshotSchedule20240805{
			ReferenceHourOfDay:    conversion.Pointer(0),
			ReferenceMinuteOfHour: conversion.Pointer(0),
			RestoreWindowDays:     conversion.Pointer(7),
		}, useYearly)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasicSleep(t); mig.PreCheckOldPreviewEnv(t) },
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

func TestMigBackupRSCloudBackupSchedule_export(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0") // in 2.0.0 we made auto_export_enabled and export fields optional only
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true, ResourceDependencyName: "mongodbatlas_cloud_backup_snapshot_export_bucket.test"})
		policyName  = acc.RandomName()
		roleName    = acc.RandomIAMRole()
		bucketName  = acc.RandomS3BucketName()

		configWithExport    = configExportPolicies(&clusterInfo, policyName, roleName, bucketName, true, true)
		configWithoutExport = configExportPolicies(&clusterInfo, policyName, roleName, bucketName, false, false)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     mig.PreCheckBasicSleep(t),
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			// Step 1: Apply config with export and auto_export_enabled (old provider)
			{
				ExternalProviders: mig.ExternalProvidersWithAWS(),
				Config:            configWithExport,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.Name),
					resource.TestCheckResourceAttr(resourceName, "auto_export_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "export.#", "1"),
				),
			},
			// Step 2: Remove export and auto_export_enabled, expect empty plan (old provider)
			{
				ExternalProviders: mig.ExternalProvidersWithAWS(),
				Config:            configWithExport,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Apply config without export and auto_export_enabled (new provider)
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configWithoutExport,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.Name),
					resource.TestCheckResourceAttr(resourceName, "auto_export_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "export.#", "0"),
				),
			},
		},
	})
}
