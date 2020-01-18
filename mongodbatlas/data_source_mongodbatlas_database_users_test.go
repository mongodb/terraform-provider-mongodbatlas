package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasDatabaseUsers_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_database_users.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	username := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	roleName := "atlasAdmin"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUsersDataSourceConfig(projectID, roleName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.db_user", "id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.db_user_1", "id"),
					resource.TestCheckResourceAttr("mongodbatlas_database_user.db_user_1", "labels.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUsersDataSourceConfigWithDS(projectID, roleName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.roles.#"),
				),
			},
		},
	})

}

func testAccMongoDBAtlasDatabaseUsersDataSourceConfig(projectID, roleName, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "db_user" {
			username      = "%[3]s"
			password      = "test-acc-password"
			project_id    = "%[1]s"
			database_name = "admin"
			
			roles {
				role_name     = "%[2]s"
				database_name = "admin"
			}
		}

		resource "mongodbatlas_database_user" "db_user_1" {
			username      = "%[3]s-1"
			password      = "test-acc-password-1"
			project_id    = "%[1]s"
			database_name = "admin"
			
			roles {
				role_name     = "%[2]s"
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
	`, projectID, roleName, username)
}

func testAccMongoDBAtlasDatabaseUsersDataSourceConfigWithDS(projectID, roleName, username string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_database_users" "test" {
			project_id = "%s"
		}
	`, testAccMongoDBAtlasDatabaseUsersDataSourceConfig(projectID, roleName, username), projectID)
}
