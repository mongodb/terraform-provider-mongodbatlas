package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccMongoDBAtlasUser_basic(t *testing.T) {
	var user matlas.AtlasUser

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resourceName := "mongodbatlas_user.test"

	username := fmt.Sprintf("john.doe%d@example.com", acctest.RandIntRange(0, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasUserConfig(projectID, orgID, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasUserExists(resourceName, &user),
					testAccCheckMongoDBAtlasUserAttributes(&user, username),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasUser_importBasic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasUserConfig(projectID, orgID, fmt.Sprintf("john.doe%d@example.com", acctest.RandIntRange(0, 255))),
			},
			{
				ResourceName:      "mongodbatlas_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasUserExists(resourceName string, user *matlas.AtlasUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if userResp, _, err := conn.AtlasUsers.Get(context.Background(), rs.Primary.ID); err == nil {
			*user = *userResp
			return nil
		}
		return fmt.Errorf("user(%s) does not exist", rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasUserAttributes(user *matlas.AtlasUser, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if user.Username != name {
			return fmt.Errorf("bad name: %s", user.Username)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasUserDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_user" {
			continue
		}

		// Try to find the user
		_, _, err := conn.Teams.Get(context.Background(), rs.Primary.Attributes["org_id"], rs.Primary.Attributes["user_id"])

		if err == nil {
			return fmt.Errorf("user (%s) still exists", rs.Primary.Attributes["user_id"])
		}
	}
	return nil
}

func testAccMongoDBAtlasUserConfig(projectID, orgID, username string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_user" "test" {
		username     = "%[1]s"
		password     = "myPassword1@"
		email_address = "%[1]s"
		mobile_number = "2125550198"
		first_name    = "John"
		last_name     = "Does"
		
		roles {
		  org_id    = "%[2]s"
		  role_name = "ORG_MEMBER"
	  
		}
	  
		roles {
		  project_id  = "%[3]s"
		  role_name   = "GROUP_READ_ONLY"
	  
		}
	  
		country = "US"
	}
	`, username, orgID, projectID)
}
