package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSMaintenanceWindow_basic(t *testing.T) {
	var (
		maintenance      matlas.MaintenanceWindow
		resourceName     = "mongodbatlas_maintenance_window.test"
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acctest.RandomWithPrefix("test-acc")
		dayOfWeek        = 7
		hourOfDay        = 3
		dayOfWeekUpdated = 4
		hourOfDayUpdated = 5
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasMaintenanceWindowConfig(orgID, projectName, dayOfWeek, hourOfDay),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasMaintenanceWindowExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "day_of_week"),
					resource.TestCheckResourceAttrSet(resourceName, "hour_of_day"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
					testAccCheckMongoDBAtlasMaintenanceWindowAttributes("day_of_week", dayOfWeek, &maintenance.DayOfWeek),
				),
			},
			{
				Config: testAccMongoDBAtlasMaintenanceWindowConfig(orgID, projectName, dayOfWeekUpdated, hourOfDayUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasMaintenanceWindowExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "day_of_week"),
					resource.TestCheckResourceAttrSet(resourceName, "hour_of_day"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeekUpdated)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDayUpdated)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
					testAccCheckMongoDBAtlasMaintenanceWindowAttributes("day_of_week", dayOfWeekUpdated, &maintenance.DayOfWeek),
				),
			},
		},
	})
}

func TestAccConfigRSMaintenanceWindow_importBasic(t *testing.T) {
	var (
		maintenance  matlas.MaintenanceWindow
		resourceName = "mongodbatlas_maintenance_window.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		dayOfWeek    = 1
		hourOfDay    = 3
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasMaintenanceWindowConfig(orgID, projectName, dayOfWeek, hourOfDay),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasMaintenanceWindowExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "day_of_week"),
					resource.TestCheckResourceAttrSet(resourceName, "hour_of_day"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
					resource.TestCheckResourceAttr(resourceName, "start_asap", "false"),
					testAccCheckMongoDBAtlasMaintenanceWindowAttributes("day_of_week", dayOfWeek, &maintenance.DayOfWeek),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasMaintenanceWindowImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConfigRSMaintenanceWindow_autoDeferActivated(t *testing.T) {
	var (
		maintenance  matlas.MaintenanceWindow
		resourceName = "mongodbatlas_maintenance_window.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		dayOfWeek    = 7
		hourOfDay    = 3
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasMaintenanceWindowConfigAutoDeferEnabled(orgID, projectName, dayOfWeek, hourOfDay),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasMaintenanceWindowExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "day_of_week"),
					resource.TestCheckResourceAttrSet(resourceName, "hour_of_day"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_defer_once_enabled", "true"),
					testAccCheckMongoDBAtlasMaintenanceWindowAttributes("day_of_week", dayOfWeek, &maintenance.DayOfWeek),
				),
			},
		},
	})
}

func testAccMongoDBAtlasMaintenanceWindowConfigAutoDeferEnabled(orgID, projectName string, dayOfWeek, hourOfDay int) string {
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
		}`, orgID, projectName, dayOfWeek, hourOfDay)
}

func testAccCheckMongoDBAtlasMaintenanceWindowExists(resourceName string, maintenance *matlas.MaintenanceWindow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] projectID: %s", rs.Primary.ID)

		maintenanceWindow, _, err := conn.MaintenanceWindows.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("maintenance Window (%s) does not exist", rs.Primary.ID)
		}

		*maintenance = *maintenanceWindow

		return nil
	}
}

func testAccCheckMongoDBAtlasMaintenanceWindowAttributes(attr string, expected int, got *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if diff := deep.Equal(expected, *got); diff != nil {
			return fmt.Errorf("bad %s \n got = %#v\nwant = %#v \ndiff = %#v", attr, expected, *got, diff)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasMaintenanceWindowImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func testAccMongoDBAtlasMaintenanceWindowConfig(orgID, projectName string, dayOfWeek, hourOfDay int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_maintenance_window" "test" {
			project_id  = mongodbatlas_project.test.id
			day_of_week = %[3]d
			hour_of_day = %[4]d
		}`, orgID, projectName, dayOfWeek, hourOfDay)
}
