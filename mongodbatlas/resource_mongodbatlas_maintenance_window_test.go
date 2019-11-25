package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
	"github.com/spf13/cast"
)

func TestAccResourceMongoDBAtlasMaintenanceWindow_basic(t *testing.T) {
	var maintenance matlas.MaintenanceWindow
	resourceName := "mongodbatlas_maintenance_window.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	dayOfWeek := 7
	hourOfDay := 3

	dayOfWeekUpdated := 4
	hourOfDayUpdated := 5

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasMaintenanceWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasMaintenanceWindowConfig(projectID, dayOfWeek, hourOfDay),
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
				Config: testAccMongoDBAtlasMaintenanceWindowConfig(projectID, dayOfWeekUpdated, hourOfDayUpdated),
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

func TestAccResourceMongoDBAtlasMaintenanceWindow_WithStartASAP(t *testing.T) {
	var maintenance matlas.MaintenanceWindow
	resourceName := "mongodbatlas_maintenance_window.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	dayOfWeek := 7
	hourOfDay := 3

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasMaintenanceWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasMaintenanceWindowConfig(projectID, dayOfWeek, hourOfDay),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasMaintenanceWindowExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "day_of_week"),
					resource.TestCheckResourceAttrSet(resourceName, "hour_of_day"),

					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
				),
			},
			{
				Config: testAccMongoDBAtlasMaintenanceWindowConfigWithStartASAP(projectID, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasMaintenanceWindowExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "start_asap", "true"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasMaintenanceWindow_importBasic(t *testing.T) {
	var maintenance matlas.MaintenanceWindow
	resourceName := "mongodbatlas_maintenance_window.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	dayOfWeek := 1
	hourOfDay := 3

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasMaintenanceWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasMaintenanceWindowConfig(projectID, dayOfWeek, hourOfDay),
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
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasMaintenanceWindowImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasMaintenanceWindowExists(resourceName string, maintenance *matlas.MaintenanceWindow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

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
			return fmt.Errorf("Maintenance Window (%s) does not exist", rs.Primary.ID)
		}
		*maintenance = *maintenanceWindow
		return nil
	}
}

func testAccCheckMongoDBAtlasMaintenanceWindowAttributes(attr string, expected int, got *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if diff := deep.Equal(expected, *got); diff != nil {
			return fmt.Errorf("Bad %s \n got = %#v\nwant = %#v \ndiff = %#v", attr, expected, *got, diff)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasMaintenanceWindowDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_maintenance_window" {
			continue
		}

		_, _, err := conn.MaintenanceWindows.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Maintenance Window (%s) does not exist", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasMaintenanceWindowImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccMongoDBAtlasMaintenanceWindowConfig(projectID string, dayOfWeek, hourOfDay int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_maintenance_window" "test" {
			project_id  = "%s"
			day_of_week = %d
			hour_of_day = %d
		}`, projectID, dayOfWeek, hourOfDay)
}

func testAccMongoDBAtlasMaintenanceWindowConfigWithStartASAP(projectID string, startAsap bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_maintenance_window" "test" {
			project_id  = "%s"
			start_asap  = %t
		}`, projectID, startAsap)
}
