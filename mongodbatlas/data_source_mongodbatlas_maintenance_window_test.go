package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigDSMaintenanceWindow_basic(t *testing.T) {
	var maintenance matlas.MaintenanceWindow

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	dayOfWeek := 7
	hourOfDay := 3

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceMaintenanceWindowConfig(orgID, projectName, dayOfWeek, hourOfDay),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasMaintenanceWindowExists("mongodbatlas_maintenance_window.test", &maintenance),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_maintenance_window.test", "project_id"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_maintenance_window.test", "day_of_week"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_maintenance_window.test", "hour_of_day"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_maintenance_window.test", "auto_defer_once_enabled"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceMaintenanceWindowConfig(orgID, projectName string, dayOfWeek, hourOfDay int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
			resource "mongodbatlas_maintenance_window" "test" {
			project_id  = mongodbatlas_project.project.id
			day_of_week = %[3]d
			hour_of_day = %[4]d
			auto_defer_once_enabled = true
		}

		data "mongodbatlas_maintenance_window" "test" {
			project_id = "${mongodbatlas_maintenance_window.test.id}"
		}
	`, orgID, projectName, dayOfWeek, hourOfDay)
}
