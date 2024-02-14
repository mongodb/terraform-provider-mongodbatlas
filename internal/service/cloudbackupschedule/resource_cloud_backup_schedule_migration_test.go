package cloudbackupschedule_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	"go.mongodb.org/atlas-sdk/v20231115006/admin"
)

func TestAccMigrationBackupRSCloudBackupSchedule_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		config       = configBasic(orgID, projectName, clusterName, &admin.DiskBackupApiPolicyItem{
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
			mig.TestStep(config),
		},
	})
}
