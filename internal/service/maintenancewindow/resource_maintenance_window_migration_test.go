package maintenancewindow_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	"github.com/spf13/cast"
)

func TestMigConfigMaintenanceWindow_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		dayOfWeek   = 7
		hourOfDay   = 3
		config      = configBasic(orgID, projectName, dayOfWeek, conversion.Pointer(hourOfDay))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
