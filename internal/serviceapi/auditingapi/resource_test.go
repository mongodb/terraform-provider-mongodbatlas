package auditingapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_auditing_api.test"

func TestAccAuditingAPI_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, true, false, `{"atype": ["authenticate"]}`),
				Check:  checkExists(resourceName),
			},
			{
				Config: configBasic(projectID, true, true, `{"atype": ["authenticate"]}`),
				Check:  checkExists(resourceName),
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

func configBasic(projectID string, enabled, auditSuccess bool, auditFilter string) string {
	auditSuccessField := ""
	if auditSuccess {
		auditSuccessField = "audit_authorization_success = true"
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_auditing_api" "test" {
			group_id     = %[1]q
			enabled      = %[2]t
			audit_filter = %[3]q
			%[4]s
		}
	`, projectID, enabled, auditFilter, auditSuccessField)
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
		if _, _, err := acc.ConnV2().AuditingApi.GetGroupAuditLog(context.Background(), groupID).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("auditing configuration for project(%s) does not exist", groupID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_auditing_api" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		auditingConfig, _, _ := acc.ConnV2().AuditingApi.GetGroupAuditLog(context.Background(), groupID).Execute()
		// Check if it's back to default settings (enabled = false means it's been reset)
		if auditingConfig.GetEnabled() {
			return fmt.Errorf("auditing configuration for project (%s) was not properly reset to defaults", groupID)
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
