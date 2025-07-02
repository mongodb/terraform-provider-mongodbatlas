package maintenancewindowapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_maintenance_window_api.test"

func TestAccMaintenanceWindowAPI_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, 3, 4),
				Check:  checkBasic(),
			},
			{
				Config: configBasic(orgID, projectName, 7, 2),
				Check:  checkBasic(),
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

func configBasic(orgID, projectName string, dayOfWeek, hourOfDay int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %q
			org_id = %q
		}

		resource "mongodbatlas_maintenance_window_api" "test" {
			project_id  = mongodbatlas_project.test.id
			day_of_week = %d
			hour_of_day = %d
		}
	`, projectName, orgID, dayOfWeek, hourOfDay)
}

func checkBasic() resource.TestCheckFunc {
	// adds checks for computed attributes not defined in config
	setAttrsChecks := []string{"number_of_deferrals", "time_zone_id"}
	checks := acc.AddAttrSetChecks(resourceName, nil, setAttrsChecks...)
	checks = append(checks, checkExists(resourceName))
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		if projectID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().MaintenanceWindowsApi.GetMaintenanceWindow(context.Background(), projectID).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("maintenance window for project(%s) does not exist", projectID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_maintenance_window_api" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		if projectID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		maintenanceWindow, _, err := acc.ConnV2().MaintenanceWindowsApi.GetMaintenanceWindow(context.Background(), projectID).Execute()
		if err != nil {
			return fmt.Errorf("maintenance window for project (%s) still exists", projectID)
		}
		// Check if it's back to default settings (day_of_week = 0 means it's been reset)
		if maintenanceWindow.GetDayOfWeek() != 0 {
			return fmt.Errorf("maintenance window for project (%s) was not properly reset to defaults", projectID)
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
		projectID := rs.Primary.Attributes["project_id"]
		if projectID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return projectID, nil
	}
}
