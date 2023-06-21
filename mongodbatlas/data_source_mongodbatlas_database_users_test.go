package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSDatabaseUsers_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_database_users.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")

	username := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	roleName := "atlasAdmin"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUsersDataSourceConfig(orgID, projectName, roleName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.db_user", "id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.db_user_1", "id"),
					resource.TestCheckResourceAttr("mongodbatlas_database_user.db_user_1", "labels.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUsersDataSourceConfigWithDS(orgID, projectName, roleName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.x509_type"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.roles.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.scopes.#"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDatabaseUsersDataSourceConfig(orgID, projectName, roleName, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_database_user" "db_user" {
			username           = %[4]q
			password           = "test-acc-password"
			project_id         = mongodbatlas_project.test.id
			auth_database_name = "admin"

			roles {
				role_name     = %[3]q
				database_name = "admin"
			}
		}

		resource "mongodbatlas_database_user" "db_user_1" {
			username           = "%[4]s-1"
			password           = "test-acc-password-1"
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
	`, orgID, projectName, roleName, username)
}

func testAccMongoDBAtlasDatabaseUsersDataSourceConfigWithDS(orgID, projectName, roleName, username string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_database_users" "test" {
			project_id = mongodbatlas_project.test.id
		}
	`, testAccMongoDBAtlasDatabaseUsersDataSourceConfig(orgID, projectName, roleName, username))
}
