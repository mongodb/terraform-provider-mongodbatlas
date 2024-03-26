package auditing_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_auditing.test"
	dataSourceName = "data.mongodbatlas_auditing.test"
)

func TestAccGenericAuditing_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		auditFilter = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, auditFilter, true, true),
				Check:  resource.ComposeTestCheckFunc(checks(auditFilter, true, true)...),
			},
			{
				Config: configBasic(projectID, "{}", false, false),
				Check:  resource.ComposeTestCheckFunc(checks("{}", false, false)...),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		auditingRes, _, _ := acc.ConnV2().AuditingApi.GetAuditingConfiguration(context.Background(), rs.Primary.ID).Execute()
		if auditingRes == nil {
			return fmt.Errorf("auditing (%s) does not exist", rs.Primary.ID)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_auditing" {
			continue
		}
		auditingRes, _, _ := acc.ConnV2().AuditingApi.GetAuditingConfiguration(context.Background(), rs.Primary.ID).Execute()
		if auditingRes.GetEnabled() {
			return fmt.Errorf("auditing (%s) exists", rs.Primary.ID)
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

		return rs.Primary.ID, nil
	}
}

func configBasic(projectID, auditFilter string, auditAuth, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_auditing" "test" {
			project_id                  = %[1]q
			audit_filter                = %[2]q
			audit_authorization_success = %[3]t
			enabled                     = %[4]t
		}
		
		data "mongodbatlas_auditing" "test" {
			project_id = mongodbatlas_auditing.test.id
		}		
	`, projectID, auditFilter, auditAuth, enabled)
}

func checks(auditFilter string, auditAuth, enabled bool) []resource.TestCheckFunc {
	commonChecks := map[string]string{
		"audit_filter":                auditFilter,
		"audit_authorization_success": strconv.FormatBool(auditAuth),
		"enabled":                     strconv.FormatBool(auditAuth),
		"configuration_type":          "FILTER_JSON",
	}
	checks := acc.AddAttrChecks(resourceName, nil, commonChecks)
	checks = acc.AddAttrChecks(dataSourceName, checks, commonChecks)
	checks = append(checks, checkExists(resourceName), checkExists(dataSourceName))
	return checks
}
