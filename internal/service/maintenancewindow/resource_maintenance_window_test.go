package maintenancewindow_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const resourceName = "mongodbatlas_maintenance_window.test"

func TestAccConfigRSMaintenanceWindow_basic(t *testing.T) {
	var (
		maintenance      matlas.MaintenanceWindow
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acctest.RandomWithPrefix("test-acc")
		dayOfWeek        = 7
		hourOfDay        = 3
		dayOfWeekUpdated = 4
		hourOfDayUpdated = 5
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, dayOfWeek, hourOfDay),
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
				Config: configBasic(orgID, projectName, dayOfWeekUpdated, hourOfDayUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "day_of_week"),
					resource.TestCheckResourceAttrSet(resourceName, "hour_of_day"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeekUpdated)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDayUpdated)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
					checkAttributes("day_of_week", dayOfWeekUpdated, &maintenance.DayOfWeek),
				),
			},
		},
	})
}

func TestAccConfigRSMaintenanceWindow_importBasic(t *testing.T) {
	var (
		maintenance matlas.MaintenanceWindow
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		dayOfWeek   = 1
		hourOfDay   = 3
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, dayOfWeek, hourOfDay),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "day_of_week"),
					resource.TestCheckResourceAttrSet(resourceName, "hour_of_day"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
					resource.TestCheckResourceAttr(resourceName, "start_asap", "false"),
					checkAttributes("day_of_week", dayOfWeek, &maintenance.DayOfWeek),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConfigRSMaintenanceWindow_autoDeferActivated(t *testing.T) {
	var (
		maintenance matlas.MaintenanceWindow
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		dayOfWeek   = 7
		hourOfDay   = 3
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithAutoDeferEnabled(orgID, projectName, dayOfWeek, hourOfDay),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &maintenance),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "day_of_week"),
					resource.TestCheckResourceAttrSet(resourceName, "hour_of_day"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", cast.ToString(dayOfWeek)),
					resource.TestCheckResourceAttr(resourceName, "hour_of_day", cast.ToString(hourOfDay)),
					resource.TestCheckResourceAttr(resourceName, "number_of_deferrals", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_defer_once_enabled", "true"),
					checkAttributes("day_of_week", dayOfWeek, &maintenance.DayOfWeek),
				),
			},
		},
	})
}

func checkExists(resourceName string, maintenance *matlas.MaintenanceWindow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		log.Printf("[DEBUG] projectID: %s", rs.Primary.ID)
		maintenanceWindow, _, err := acc.Conn().MaintenanceWindows.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("maintenance Window (%s) does not exist", rs.Primary.ID)
		}
		*maintenance = *maintenanceWindow
		return nil
	}
}

func checkAttributes(attr string, expected int, got *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if diff := deep.Equal(expected, *got); diff != nil {
			return fmt.Errorf("bad %s \n got = %#v\nwant = %#v \ndiff = %#v", attr, expected, *got, diff)
		}

		return nil
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func configBasic(orgID, projectName string, dayOfWeek, hourOfDay int) string {
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

func configWithAutoDeferEnabled(orgID, projectName string, dayOfWeek, hourOfDay int) string {
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
