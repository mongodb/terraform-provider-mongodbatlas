package projectsettingsapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_project_settings_api.test"

func TestAccProjectSettingsAPI_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: config(projectID, false),
				Check:  check(projectID, false),
			},
			{
				Config: config(projectID, true),
				Check:  check(projectID, true),
			},
			{
				Config: config(projectID, false),
				Check:  check(projectID, false),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "group_id",
			},
		},
	})
}

func config(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_settings_api" "test" {
			group_id = %[1]q
			is_collect_database_specifics_statistics_enabled = %[2]t
			is_data_explorer_enabled = %[2]t
			is_data_explorer_gen_ai_features_enabled = %[2]t
			is_data_explorer_gen_ai_sample_document_passing_enabled = %[2]t
			is_extended_storage_sizes_enabled = %[2]t
			is_performance_advisor_enabled = %[2]t
			is_realtime_performance_panel_enabled = %[2]t
			is_schema_advisor_enabled = %[2]t
		}
	`, projectID, enabled)
}

func check(projectID string, enabled bool) resource.TestCheckFunc {
	expected := fmt.Sprintf("%t", enabled)

	check := resource.ComposeAggregateTestCheckFunc(
		checkExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
		resource.TestCheckResourceAttr(resourceName, "is_collect_database_specifics_statistics_enabled", expected),
		resource.TestCheckResourceAttr(resourceName, "is_data_explorer_enabled", expected),
		resource.TestCheckResourceAttr(resourceName, "is_data_explorer_gen_ai_features_enabled", expected),
		resource.TestCheckResourceAttr(resourceName, "is_data_explorer_gen_ai_sample_document_passing_enabled", expected),
		resource.TestCheckResourceAttr(resourceName, "is_extended_storage_sizes_enabled", expected),
		resource.TestCheckResourceAttr(resourceName, "is_performance_advisor_enabled", expected),
		resource.TestCheckResourceAttr(resourceName, "is_realtime_performance_panel_enabled", expected),
		resource.TestCheckResourceAttr(resourceName, "is_schema_advisor_enabled", expected),
	)

	return check
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return fmt.Errorf("no group_id is set")
		}
		if _, _, err := acc.ConnV2().ProjectsApi.GetGroupSettings(context.Background(), groupID).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("project settings for project(%s) do not exist", groupID)
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		groupID := rs.Primary.Attributes["group_id"]
		return groupID, nil
	}
}
