package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasAuditing_basic(t *testing.T) {
	var auditing matlas.Auditing
	resourceName := "mongodbatlas_auditing.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	auditAuth := true
	auditFilter := "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
	enabled := true

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAuditingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAuditingConfig(projectID, auditFilter, auditAuth, enabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAuditingExists(resourceName, &auditing),

					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_filter"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_authorization_success"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "audit_filter", auditFilter),
					resource.TestCheckResourceAttr(resourceName, "audit_authorization_success", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configuration_type", "FILTER_JSON"),
				),
			},
			{
				Config: testAccMongoDBAtlasAuditingConfig(projectID, "{}", false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAuditingExists(resourceName, &auditing),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_filter"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_authorization_success"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "audit_filter", "{}"),
					resource.TestCheckResourceAttr(resourceName, "audit_authorization_success", "false"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configuration_type", "FILTER_JSON"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasAuditing_importBasic(t *testing.T) {
	auditing := &matlas.Auditing{}
	resourceName := "mongodbatlas_auditing.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	auditAuth := true
	auditFilter := "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
	enabled := true

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAuditingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAuditingConfig(projectID, auditFilter, auditAuth, enabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAuditingExists(resourceName, auditing),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_filter"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_authorization_success"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "audit_filter", auditFilter),
					resource.TestCheckResourceAttr(resourceName, "audit_authorization_success", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configuration_type", "FILTER_JSON"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasAuditingImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasAuditingExists(resourceName string, auditing *matlas.Auditing) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] projectID: %s", rs.Primary.ID)

		auditingRes, _, err := conn.Auditing.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Auditing (%s) does not exist", rs.Primary.ID)
		}
		auditing = auditingRes
		return nil
	}
}

func testAccCheckMongoDBAtlasAuditingDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_auditing" {
			continue
		}

		_, _, err := conn.Auditing.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Auditing (%s) does not exist", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasAuditingImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccMongoDBAtlasAuditingConfig(projectID, auditFilter string, auditAuth, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_auditing" "test" {
			project_id                  = "%s"
			audit_filter                = "%s"
			audit_authorization_success = %t
			enabled                     = %t
		}`, projectID, auditFilter, auditAuth, enabled)
}
