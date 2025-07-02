package auditingapi_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_auditing_api.test"

func TestAccAuditingAPI_basic(t *testing.T) {
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
				Config: configBasic(orgID, projectName, true, "{}"),
				Check:  checkBasic(),
			},
			{
				Config: configBasic(orgID, projectName, false, `{"atype":"authenticate"}`),
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

func configBasic(orgID, projectName string, enabled bool, auditFilter string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %q
			org_id = %q
		}

		resource "mongodbatlas_auditing_api" "test" {
			project_id = mongodbatlas_project.test.id
			enabled    = %t
			audit_filter = %q
		}
	`, projectName, orgID, enabled, auditFilter)
}

func checkBasic() resource.TestCheckFunc {
	// adds checks for computed attributes not defined in config
	setAttrsChecks := []string{"configuration_type"}
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
		if _, _, err := acc.ConnV2().AuditingApi.GetAuditingConfiguration(context.Background(), projectID).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("auditing configuration for project(%s) does not exist", projectID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_auditing_api" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		if projectID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		auditingConfig, _, err := acc.ConnV2().AuditingApi.GetAuditingConfiguration(context.Background(), projectID).Execute()
		if err != nil {
			return fmt.Errorf("auditing configuration for project (%s) still exists", projectID)
		}
		// Check if it's back to default settings (enabled = false means it's been reset)
		if auditingConfig.GetEnabled() {
			return fmt.Errorf("auditing configuration for project (%s) was not properly reset to defaults", projectID)
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
