package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasAuditing_basic(t *testing.T) {
	var auditing matlas.Auditing
	dataSourceName := "data.mongodbatlas_auditing.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	auditAuth := true
	auditFilter := "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
	enabled := true

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceAuditingConfig(projectID, auditFilter, auditAuth, enabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAuditingExists("mongodbatlas_auditing.test", &auditing),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),

					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "audit_filter", auditFilter),
					resource.TestCheckResourceAttr(dataSourceName, "audit_authorization_success", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "configuration_type", "FILTER_JSON"),
				),
			},
			{
				Config: testAccMongoDBAtlasDataSourceAuditingConfig(projectID, "{}", false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAuditingExists("mongodbatlas_auditing.test", &auditing),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),

					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "audit_filter", "{}"),
					resource.TestCheckResourceAttr(dataSourceName, "audit_authorization_success", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "configuration_type", "FILTER_JSON"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceAuditingConfig(projectID, auditFilter string, auditAuth, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_auditing" "test" {
			project_id                  = "%s"
			audit_filter                = "%s"
			audit_authorization_success = %t
			enabled                     = %t
		}

		data "mongodbatlas_auditing" "test" {
			project_id = "${mongodbatlas_auditing.test.id}"
		}
	`, projectID, auditFilter, auditAuth, enabled)
}
