package acc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.mongodb.org/atlas-sdk/v20250312021/admin"
)

func CheckDestroyProject(s *terraform.State) error {
	return checkDestroyProject(ConnV2(), s)
}

func CheckDestroyProjectGov(s *terraform.State) error {
	return checkDestroyProject(ConnV2UsingGov(), s)
}

func checkDestroyProject(conn *admin.APIClient, s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project" {
			continue
		}
		projectRes, _, _ := conn.ProjectsApi.GetGroupByName(context.Background(), rs.Primary.ID).Execute()
		if projectRes != nil {
			return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

const dataSourceProjectByID = `
	data "mongodbatlas_project" "test" {
		project_id = mongodbatlas_project.test.id
	}`

func ConfigProjectWithSettings(projectName, orgID, projectOwnerID string, value, withDS bool) string {
	ds := ""
	if withDS {
		ds = dataSourceProjectByID
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name             = %[1]q
			org_id           = %[2]q
			project_owner_id = %[3]q
			is_collect_database_specifics_statistics_enabled = %[4]t
			is_data_explorer_enabled = %[4]t
			is_extended_storage_sizes_enabled = %[4]t
			is_performance_advisor_enabled = %[4]t
			is_realtime_performance_panel_enabled = %[4]t
			is_schema_advisor_enabled = %[4]t
			is_cluster_ai_assistant_enabled = %[4]t
			is_data_explorer_gen_ai_features_enabled = %[4]t
			is_data_explorer_gen_ai_sample_document_passing_enabled = %[4]t
		}
	`, projectName, orgID, projectOwnerID, value) + ds
}

func ConfigProjectWithoutSettings(projectName, orgID, projectOwnerID string, withDS bool) string {
	ds := ""
	if withDS {
		ds = dataSourceProjectByID
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name             = %[1]q
			org_id           = %[2]q
			project_owner_id = %[3]q
		}
	`, projectName, orgID, projectOwnerID) + ds
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
