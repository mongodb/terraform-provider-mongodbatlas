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

	resourceName := "data.mongodbatlas_teams.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	username := "mongodbatlas.testing@gmail.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasTeamConfig(orgID, projectID, name, username, "GROUP_READ_ONLY"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasTeamExists(resourceName, &team),
					testAccCheckMongoDBAtlasTeamAttributes(&team, name),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "team_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "team_roles.#", "1"),
				),
			},
		},
	})

}

func testAccDataSourceMongoDBAtlasTeamConfig(orgID, projectID, name, username, teamRoles string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_teams" "test" {
			org_id     = "%s"
			project_id = "%s"
			name       = "%s"
			usernames  = ["%s"]
			team_roles = ["%s"]
		}
		
		data "mongodbatlas_teams" "test" {
			org_id     = mongodbatlas_teams.test.org_id
			team_id    = mongodbatlas_teams.test.team_id
			project_id = mongodbatlas_teams.test.project_id
		}	
	`, orgID, projectID, name, username, teamRoles)
}
