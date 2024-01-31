package maintenancewindow_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccMigrationConfigMaintenanceWindow_basic(t *testing.T) {
	var (
		maintenance matlas.MaintenanceWindow
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		dayOfWeek   = 7
		hourOfDay   = 3
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configBasic(orgID, projectName, dayOfWeek, hourOfDay),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "day_of_week"),
					resource.TestCheckResourceAttrSet(resourceName, "hour_of_day"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
					checkAttributes("day_of_week", dayOfWeek, &maintenance.DayOfWeek),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configBasic(orgID, projectName, dayOfWeek, hourOfDay),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
