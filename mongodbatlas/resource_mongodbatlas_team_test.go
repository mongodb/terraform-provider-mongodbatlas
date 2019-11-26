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

func TestAccResourceMongoDBAtlasTeam_basic(t *testing.T) {
	var team matlas.Team

	resourceName := "mongodbatlas_team.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	username := "mongodbatlas.testing@gmail.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, name, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasTeamExists(resourceName, &team),
					testAccCheckMongoDBAtlasTeamAttributes(&team, name),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "usernames.0", username),
				),
			},
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, name, "marin.salinas@digitalonus.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasTeamExists(resourceName, &team),
					testAccCheckMongoDBAtlasTeamAttributes(&team, name),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "usernames.0", "marin.salinas@digitalonus.com"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasTeam_importBasic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resourceName := "mongodbatlas_team.test"

	username := "mongodbatlas.testing@gmail.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, name, username),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasTeamStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasTeamExists(resourceName string, team *matlas.Team) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["org_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] orgID: %s", rs.Primary.Attributes["org_id"])

		if teamResp, _, err := conn.Teams.Get(context.Background(), rs.Primary.Attributes["org_id"], rs.Primary.Attributes["team_id"]); err == nil {
			*team = *teamResp
			return nil
		}
		return fmt.Errorf("team(%s) does not exist", rs.Primary.Attributes["team_id"])
	}
}

func testAccCheckMongoDBAtlasTeamAttributes(team *matlas.Team, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if team.Name != name {
			return fmt.Errorf("bad name: %s", team.Name)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasTeamDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_team" {
			continue
		}

		// Try to find the team
		_, _, err := conn.Teams.Get(context.Background(), rs.Primary.Attributes["org_id"], rs.Primary.Attributes["team_id"])

		if err == nil {
			return fmt.Errorf("team (%s) still exists", rs.Primary.Attributes["team_id"])
		}
	}
	return nil
}

func testAccMongoDBAtlasTeamConfig(orgID, name, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_team" "test" {
			org_id    = "%s"
			name      = "%s"
			usernames = ["%s"]
		}
	`, orgID, name, username)
}

func testAccCheckMongoDBAtlasTeamStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["org_id"], rs.Primary.Attributes["team_id"]), nil
	}
}
