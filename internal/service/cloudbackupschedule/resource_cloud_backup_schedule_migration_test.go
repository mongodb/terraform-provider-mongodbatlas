package cloudbackupschedule_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	"go.mongodb.org/atlas-sdk/v20231115006/admin"
)

func TestAccMigrationBackupRSCloudBackupSchedule_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
		config      = configBasic(orgID, projectName, clusterName, &admin.DiskBackupApiPolicyItem{
			FrequencyInterval: 1,
			RetentionUnit:     "days",
			RetentionValue:    1,
		})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "reference_hour_of_day"),
					resource.TestCheckResourceAttrSet(resourceName, "reference_minute_of_hour"),
					resource.TestCheckResourceAttrSet(resourceName, "restore_window_days"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_hourly.#"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_daily.#"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_weekly.#"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_item_monthly.#"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
