package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigDSDatabaseUser_basic(t *testing.T) {
	var dbUser matlas.DatabaseUser

	resourceName := "data.mongodbatlas_database_user.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	roleName := "atlasAdmin"
	username := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserDataSourceConfig(orgID, projectName, roleName, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "x509_type", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "roles.0.database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDatabaseUserDataSourceConfig(orgID, projectName, roleName, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_database_user" "test" {
			username           = %[4]q
			password           = "test-acc-password"
			project_id         = mongodbatlas_project.test.id
			auth_database_name = "admin"

			roles {
				role_name     = %[3]q
				database_name = "admin"
			}

			labels {
				key   = "key 1"
				value = "value 1"
			}
			labels {
				key   = "key 2"
				value = "value 2"
			}
		}

		data "mongodbatlas_database_user" "test" {
			username           = mongodbatlas_database_user.test.username
			project_id         = mongodbatlas_database_user.test.project_id
			auth_database_name = "admin"
		}
	`, orgID, projectName, roleName, username)
}
