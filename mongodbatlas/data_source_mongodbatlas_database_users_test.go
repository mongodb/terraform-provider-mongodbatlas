package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasDatabaseUsers_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_database_users.test"
	groupID := "5cf5a45a9ccf6400e60981b6" // Modify until project data source is created.

	roleName := "atlasAdmin"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUsersDataSourceConfig(groupID, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.db_user", "id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.db_user_1", "id"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUsersDataSourceConfigWithDS(groupID, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.roles.#"),
				),
			},
		},
	})

}

func testAccMongoDBAtlasDatabaseUsersDataSourceConfig(groupID, roleName string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_database_user" "db_user" {
	username      = "test-acc-username"
	password      = "test-acc-password"
	group_id      = "%[1]s"
	database_name = "admin"
	
	roles {
		role_name     = "%[2]s"
		database_name = "admin"
	}
}

resource "mongodbatlas_database_user" "db_user_1" {
	username      = "test-acc-username-1"
	password      = "test-acc-password-1"
	group_id      = "%[1]s"
	database_name = "admin"
	
	roles {
		role_name     = "%[2]s"
		database_name = "admin"
	}
}

`, groupID, roleName)
}

func testAccMongoDBAtlasDatabaseUsersDataSourceConfigWithDS(groupID, roleName string) string {
	return fmt.Sprintf(`
%s

data "mongodbatlas_database_users" "test" {
	group_id = "%s"
}
`, testAccMongoDBAtlasDatabaseUsersDataSourceConfig(groupID, roleName), groupID)
}
