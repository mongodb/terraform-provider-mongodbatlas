package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

func CheckDestroyProject(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project" {
			continue
		}

		projectRes, _, _ := Conn().Projects.GetOneProjectByName(context.Background(), rs.Primary.ID)
		if projectRes != nil {
			return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func ConfigProjectWithUpdatedRole(projectName, orgID, teamID, roleName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"

			teams {
				team_id = "%s"
				role_names = ["%s"]
			}
		}
	`, projectName, orgID, teamID, roleName)
}

func ConfigProjectWithOwner(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%[1]s"
			org_id 			 = "%[2]s"
		    project_owner_id = "%[3]s"
		}
	`, projectName, orgID, projectOwnerID)
}

func ConfigProjectGovWithOwner(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%[1]s"
			org_id 			 = "%[2]s"
		    project_owner_id = "%[3]s"
			region_usage_restrictions = "GOV_REGIONS_ONLY"
		}
	`, projectName, orgID, projectOwnerID)
}

func ConfigProjectWithFalseDefaultSettings(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%[1]s"
			org_id 			 = "%[2]s"
			project_owner_id = "%[3]s"
			with_default_alerts_settings = false
		}
	`, projectName, orgID, projectOwnerID)
}

func ConfigProjectWithSettings(projectName, orgID, projectOwnerID string, value bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = %[1]q
			org_id 			 = %[2]q
			project_owner_id = %[3]q
			with_default_alerts_settings = %[4]t
			is_collect_database_specifics_statistics_enabled = %[4]t
			is_data_explorer_enabled = %[4]t
			is_extended_storage_sizes_enabled = %[4]t
			is_performance_advisor_enabled = %[4]t
			is_realtime_performance_panel_enabled = %[4]t
			is_schema_advisor_enabled = %[4]t
		}
	`, projectName, orgID, projectOwnerID, value)
}

func ConfigProjectWithLimits(projectName, orgID string, limits []*admin.DataFederationLimit) string {
	var limitsString string

	for _, limit := range limits {
		limitsString += fmt.Sprintf(`
		limits {
			name = "%s"
			value = %d
		}
		`, limit.Name, limit.Value)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%s"
			org_id 			 = "%s"

			%s
		}
	`, projectName, orgID, limitsString)
}

func ImportStateProjectIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
