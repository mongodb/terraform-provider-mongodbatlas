package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasDatabaseUser_basic(t *testing.T) {
	var dbUser matlas.DatabaseUser

	resourceName := "mongodbatlas_database_user.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	username := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserConfig(projectID, "atlasAdmin", username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.role_name", "atlasAdmin"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUserConfig(projectID, "read", username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.0.role_name", "read"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasDatabaseUser_importBasic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	username := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	importStateID := fmt.Sprintf("%s-%s", projectID, username)

	resourceName := "mongodbatlas_database_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserConfig(projectID, "read", username),
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           importStateID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasDatabaseUserExists(resourceName string, dbUser *matlas.DatabaseUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])

		if dbUserResp, _, err := conn.DatabaseUsers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["username"]); err == nil {
			*dbUser = *dbUserResp
			return nil
		}
		return fmt.Errorf("database user(%s) does not exist", rs.Primary.Attributes["project_id"])
	}
}

func testAccCheckMongoDBAtlasDatabaseUserAttributes(dbUser *matlas.DatabaseUser, username string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if dbUser.Username != username {
			return fmt.Errorf("bad username: %s", dbUser.Username)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasDatabaseUserDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_database_user" {
			continue
		}

		// Try to find the database user
		_, _, err := conn.DatabaseUsers.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["username"])

		if err == nil {
			return fmt.Errorf("database user (%s) still exists", rs.Primary.Attributes["project_id"])
		}
	}
	return nil
}

func testAccMongoDBAtlasDatabaseUserConfig(projectID, roleName, username string) string {
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
	`, projectID, roleName, username)
}
