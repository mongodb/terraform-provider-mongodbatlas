package auditing_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestMain(m *testing.M) {
	acc.SetupSharedResources()
	exitCode := m.Run()
	acc.CleanupSharedResources()
	os.Exit(exitCode)
}

func TestAccGenericAuditing_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_auditing.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		auditAuth    = true
		auditFilter  = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		enabled      = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, auditFilter, auditAuth, enabled),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_filter"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_authorization_success"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "audit_filter", auditFilter),
					resource.TestCheckResourceAttr(resourceName, "audit_authorization_success", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configuration_type", "FILTER_JSON"),
				),
			},
			{
				Config: configBasic(orgID, projectName, "{}", false, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_filter"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_authorization_success"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "audit_filter", "{}"),
					resource.TestCheckResourceAttr(resourceName, "audit_authorization_success", "false"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configuration_type", "FILTER_JSON"),
				),
			},
		},
	})
}

func TestAccGenericAuditing_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_auditing.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		auditAuth    = true
		auditFilter  = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		enabled      = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, auditFilter, auditAuth, enabled),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_filter"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_authorization_success"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "audit_filter", auditFilter),
					resource.TestCheckResourceAttr(resourceName, "audit_authorization_success", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configuration_type", "FILTER_JSON"),
				),
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
		if auditingRes != nil {
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

func configBasic(orgID, projectName, auditFilter string, auditAuth, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_auditing" "test" {
			project_id                  = mongodbatlas_project.test.id
			audit_filter                = %[3]q
			audit_authorization_success = %[4]t
			enabled                     = %[5]t
		}
		
		data "mongodbatlas_auditing" "test" {
			project_id = mongodbatlas_auditing.test.id
		}		
	`, orgID, projectName, auditFilter, auditAuth, enabled)
}
