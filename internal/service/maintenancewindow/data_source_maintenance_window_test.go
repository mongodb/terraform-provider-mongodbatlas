package maintenancewindow_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const dataSourceName = "mongodbatlas_maintenance_window.test"

func TestAccConfigDSMaintenanceWindow_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		dayOfWeek   = 7
		hourOfDay   = 3
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDS(orgID, projectName, dayOfWeek, hourOfDay),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(dataSourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(dataSourceName, "auto_defer_once_enabled", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "time_zone_id"),
				),
			},
		},
	})
}

func configDS(orgID, projectName string, dayOfWeek, hourOfDay int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
			resource "mongodbatlas_maintenance_window" "test" {
			project_id  = mongodbatlas_project.test.id
			day_of_week = %[3]d
			hour_of_day = %[4]d
			auto_defer_once_enabled = true
		}

		data "mongodbatlas_maintenance_window" "test" {
			project_id = mongodbatlas_maintenance_window.test.id
		}
	`, orgID, projectName, dayOfWeek, hourOfDay)
}
