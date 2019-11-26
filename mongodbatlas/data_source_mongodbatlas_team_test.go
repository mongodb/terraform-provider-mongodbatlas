package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasTeam_basic(t *testing.T) {
	var team matlas.Team

	resourceName := "data.mongodbatlas_team.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	username := "mongodbatlas.testing@gmail.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasTeamConfig(orgID, name, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasTeamExists(resourceName, &team),
					testAccCheckMongoDBAtlasTeamAttributes(&team, name),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users.0.username", username),
				),
			},
		},
	})

}

func testAccDataSourceMongoDBAtlasTeamConfig(orgID, name, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_team" "test" {
			org_id    = "%s"
			name      = "%s"
			usernames = ["%s"]
		}

		data "mongodbatlas_team" "test" {
			org_id    = mongodbatlas_team.test.org_id
			team_id   = mongodbatlas_team.test.team_id
		}
	`, orgID, name, username)
}
