package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasDatabaseUser_basic(t *testing.T) {
	var dbUser matlas.DatabaseUser

	resourceName := "data.mongodbatlas_database_user.test"
	projectID := "5cf5a45a9ccf6400e60981b6" // Modify until project data source is created.

	roleName := "atlasAdmin"

	username := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserDataSourceConfig(projectID, roleName, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "roles.0.database_name", "admin"),
				),
			},
		},
	})

}

func testAccMongoDBAtlasDatabaseUserDataSourceConfig(projectID, roleName, username string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_database_user" "test" {
	username      = "%[3]s"
	password      = "test-acc-password"
	project_id    = "%[1]s"
	database_name = "admin"
	
	roles {
		role_name     = "%[2]s"
		database_name = "admin"
	}
}

data "mongodbatlas_database_user" "test" {
	username   = mongodbatlas_database_user.test.username
	project_id = mongodbatlas_database_user.test.project_id
}
`, projectID, roleName, username)
}
