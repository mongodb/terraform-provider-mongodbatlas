package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasTeam_basic(t *testing.T) {
	var team matlas.Team

	resourceName := "data.mongodbatlas_teams.test"
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
					resource.TestCheckResourceAttrSet(resourceName, "team_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
		},
	})

}

func testAccDataSourceMongoDBAtlasTeamConfig(orgID, name, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_teams" "test" {
			org_id     = "%s"
			name       = "%s"
			usernames  = ["%s"]
		}

		data "mongodbatlas_teams" "test" {
			org_id     = mongodbatlas_teams.test.org_id
			team_id    = mongodbatlas_teams.test.team_id
		}

		data "mongodbatlas_teams" "test2" {
			org_id     = mongodbatlas_teams.test.org_id
			name    = mongodbatlas_teams.test.name
		}
	`, orgID, name, username)
}
