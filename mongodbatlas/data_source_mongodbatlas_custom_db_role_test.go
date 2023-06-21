package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSCustomDBRole_basic(t *testing.T) {
	resourceName := "mongodbatlas_custom_db_role.test"
	dataSourceName := "data.mongodbatlas_custom_db_role.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	roleName := fmt.Sprintf("test-acc-custom_role-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasCustomDBRoleConfig(orgID, projectName, roleName, "INSERT", fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
				Check: resource.ComposeTestCheckFunc(
					// Test for Resource
					testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_name"),
					resource.TestCheckResourceAttrSet(resourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "INSERT"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.resources.#", "1"),

					// Test for Data source
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role_name"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasCustomDBRoleConfig(orgID, projectName, roleName, action, databaseName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_custom_db_role" "test" {
			project_id = mongodbatlas_project.test.id
			role_name  = %[3]q

			actions {
				action = %[4]q
				resources {
					collection_name = ""
					database_name   = %[5]q
				}
			}
		}

		data "mongodbatlas_custom_db_role" "test" {
			project_id = "${mongodbatlas_custom_db_role.test.project_id}"
			role_name  = "${mongodbatlas_custom_db_role.test.role_name}"
		}
	`, orgID, projectName, roleName, action, databaseName)
}
