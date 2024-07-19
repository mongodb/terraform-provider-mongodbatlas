package cloudbackupschedule_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	admin20231115 "go.mongodb.org/atlas-sdk/v20231115014/admin"
)

func TestMigBackupRSCloudBackupSchedule_basic(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true})
		useYearly   = mig.IsProviderVersionAtLeast("1.16.0") // attribute introduced in this version
		config      = configNewPolicies(&clusterInfo, &admin20231115.DiskBackupSnapshotSchedule{
			ReferenceHourOfDay:    conversion.Pointer(0),
			ReferenceMinuteOfHour: conversion.Pointer(0),
			RestoreWindowDays:     conversion.Pointer(7),
		}, useYearly)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.Name),
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
