package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasMaintenanceWindow_basic(t *testing.T) {
	var maintenance matlas.MaintenanceWindow

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	dayOfWeek := 7
	hourOfDay := 3

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceMaintenanceWindowConfig(projectID, dayOfWeek, hourOfDay),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasMaintenanceWindowExists("mongodbatlas_maintenance_window.test", &maintenance),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_maintenance_window.test", "project_id"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_maintenance_window.test", "day_of_week"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_maintenance_window.test", "hour_of_day"),
				),
			},
		},
	})
}

func TestAccDataSourceMongoDBAtlasMaintenanceWindow_basicWithStartASAP(t *testing.T) {
	var maintenance matlas.MaintenanceWindow

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceMaintenanceWindowConfigWithStartASAP(projectID, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasMaintenanceWindowExists("mongodbatlas_maintenance_window.test", &maintenance),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_maintenance_window.test", "start_asap"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceMaintenanceWindowConfig(projectID string, dayOfWeek, hourOfDay int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_maintenance_window" "test" {
			project_id  = "%s"
			day_of_week = %d
			hour_of_day = %d
		}

		data "mongodbatlas_maintenance_window" "test" {
			project_id = "${mongodbatlas_maintenance_window.test.id}"
		}
	`, projectID, dayOfWeek, hourOfDay)
}

func testAccMongoDBAtlasDataSourceMaintenanceWindowConfigWithStartASAP(projectID string, startAsap bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_maintenance_window" "test" {
			project_id  = "%s"
			start_asap  = %t
		}
		
		data "mongodbatlas_maintenance_window" "test" {
			project_id  = "${mongodbatlas_maintenance_window.test.id}"
		}
	`, projectID, startAsap)
}
