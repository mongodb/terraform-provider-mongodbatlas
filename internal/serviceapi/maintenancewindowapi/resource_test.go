package maintenancewindowapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_maintenance_window_api.test"

func TestAccMaintenanceWindowAPI_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, 3, 4, false),
				Check:  checkExists(resourceName),
			},
			{
				Config: configBasic(projectID, 7, 2, false),
				Check:  checkExists(resourceName),
			},
			{
				Config: configBasic(projectID, 7, 2, true),
				Check:  checkExists(resourceName),
			},
			{
				Config: configBasic(projectID, 7, 2, false),
				Check:  checkExists(resourceName),
			},
			{
				Config:   configBasic(projectID, 7, 2, false),
				PlanOnly: true, // Check that nested attribute protected_hours was correctly unset in API.
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "group_id",
				// Patch doesn't return a response body, but these attributes are changed when Get is called.
				ImportStateVerifyIgnore: []string{"auto_defer_once_enabled", "number_of_deferrals", "time_zone_id", "start_asap"},
			},
		},
	})
}

func configBasic(projectID string, dayOfWeek, hourOfDay int, useProtectedHours bool) string {
	protectedHours := ""
	if useProtectedHours {
		protectedHours = `
			protected_hours = {
				start_hour_of_day = 18
				end_hour_of_day   = 23
			}`
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_maintenance_window_api" "test" {
			group_id    = %[1]q
			day_of_week = %[2]d
			hour_of_day = %[3]d
			
			%[4]s
		}
	`, projectID, dayOfWeek, hourOfDay, protectedHours)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().MaintenanceWindowsApi.GetMaintenanceWindow(context.Background(), groupID).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("maintenance window for project(%s) does not exist", groupID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_maintenance_window_api" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		maintenanceWindow, _, err := acc.ConnV2().MaintenanceWindowsApi.GetMaintenanceWindow(context.Background(), groupID).Execute()
		if err != nil {
			return fmt.Errorf("maintenance window for project (%s) still exists", groupID)
		}
		// Check if it's back to default settings (day_of_week = 0 means it's been reset)
		if maintenanceWindow.GetDayOfWeek() != 0 {
			return fmt.Errorf("maintenance window for project (%s) was not properly reset to defaults", groupID)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return groupID, nil
	}
}
