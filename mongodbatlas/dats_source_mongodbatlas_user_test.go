package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasUser_basic(t *testing.T) {
	var user matlas.AtlasUser

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resourceName := "data.mongodbatlas_user.test"

	username := fmt.Sprintf("john.doe%d@example.com", acctest.RandIntRange(0, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasUserConfig(projectID, orgID, username),
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

func testAccDataSourceMongoDBAtlasUserConfig(projectID, orgID, username string) string {
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

	data "mongodbatlas_user" "test" {
		user_id = mongodbatlas_user.test.user_id
	}

	`, username, orgID, projectID)
}
