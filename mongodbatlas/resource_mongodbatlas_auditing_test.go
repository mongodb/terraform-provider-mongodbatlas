package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccAdvRSAuditing_basic(t *testing.T) {
	var (
		auditing     matlas.Auditing
		resourceName = "mongodbatlas_auditing.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		auditAuth    = true
		auditFilter  = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		enabled      = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAuditingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAuditingConfig(orgID, projectName, auditFilter, auditAuth, enabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAuditingExists(resourceName, &auditing),

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
				Config: testAccMongoDBAtlasAuditingConfig(orgID, projectName, "{}", false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAuditingExists(resourceName, &auditing),
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

func TestAccAdvRSAuditing_importBasic(t *testing.T) {
	var (
		auditing     = &matlas.Auditing{}
		resourceName = "mongodbatlas_auditing.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		auditAuth    = true
		auditFilter  = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		enabled      = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAuditingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAuditingConfig(orgID, projectName, auditFilter, auditAuth, enabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAuditingExists(resourceName, auditing),
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
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		auditingRes, _, err := conn.Auditing.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("auditing (%s) does not exist", rs.Primary.ID)
		}

		auditing = auditingRes

		return nil
	}
}

func testAccCheckMongoDBAtlasAuditingDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_auditing" {
			continue
		}

		_, _, err := conn.Auditing.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("auditing (%s) does not exist", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasAuditingImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func testAccMongoDBAtlasAuditingConfig(orgID, projectName, auditFilter string, auditAuth, enabled bool) string {
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
		}`, orgID, projectName, auditFilter, auditAuth, enabled)
}
