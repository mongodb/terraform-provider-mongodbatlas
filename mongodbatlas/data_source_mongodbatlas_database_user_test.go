package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasDatabaseUser_basic(t *testing.T) {
	var dbUser matlas.DatabaseUser

	resourceName := "data.mongodbatlas_database_user.test"
	groupID := "5cf5a45a9ccf6400e60981b6" // Modify until project data source is created.

	roleName := "atlasAdmin"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserDataSourceConfig(groupID, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, "test-acc-username"),
					resource.TestCheckResourceAttrSet(resourceName, "group_id"),
					resource.TestCheckResourceAttr(resourceName, "username", "test-acc-username"),
					resource.TestCheckResourceAttr(resourceName, "database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "roles.0.database_name", "admin"),
				),
			},
		},
	})

}

func testAccMongoDBAtlasDatabaseUserDataSourceConfig(groupID, roleName string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_database_user" "test" {
	username      = "test-acc-username"
	password      = "test-acc-password"
	group_id      = "%s"
	database_name = "admin"
	
	roles {
		role_name     = "%s"
		database_name = "admin"
	}
}

data "mongodbatlas_database_user" "test" {
	group_id = mongodbatlas_database_user.test.group_id
	username = mongodbatlas_database_user.test.username
}
`, groupID, roleName)
}
