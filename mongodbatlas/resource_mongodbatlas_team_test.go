package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasTeam_basic(t *testing.T) {
	var team matlas.Team

	resourceName := "mongodbatlas_teams.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	updatedName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, name,
					[]string{
						"mongodbatlas.testing@gmail.com",
						"francisco.preciado@digitalonus.com",
						"antonio.cabrera@digitalonus.com",
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasTeamExists(resourceName, &team),
					testAccCheckMongoDBAtlasTeamAttributes(&team, name),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "3"),
				),
			},
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, updatedName,
					[]string{
						"marin.salinas@digitalonus.com",
						"antonio.cabrera@digitalonus.com",
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasTeamExists(resourceName, &team),
					testAccCheckMongoDBAtlasTeamAttributes(&team, updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, updatedName,
					[]string{
						"marin.salinas@digitalonus.com",
						"mongodbatlas.testing@gmail.com",
						"francisco.preciado@digitalonus.com",
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasTeamExists(resourceName, &team),
					testAccCheckMongoDBAtlasTeamAttributes(&team, updatedName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "3"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasTeam_importBasic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	resourceName := "mongodbatlas_teams.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, name,
					[]string{"mongodbatlas.testing@gmail.com"},
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasTeamStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
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

		orgID := rs.Primary.Attributes["org_id"]
		id := rs.Primary.Attributes["team_id"]

		if orgID == "" && id == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] orgID: %s", orgID)
		log.Printf("[DEBUG] teamID: %s", id)

		teamResp, _, err := conn.Teams.Get(context.Background(), orgID, id)
		if err == nil {
			*team = *teamResp
			return nil
		}
		return fmt.Errorf("team(%s) does not exist", id)
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
		if rs.Type != "mongodbatlas_teams" {
			continue
		}

		orgID := rs.Primary.Attributes["org_id"]
		id := rs.Primary.Attributes["team_id"]

		// Try to find the team
		_, _, err := conn.Teams.Get(context.Background(), orgID, id)
		if err == nil {
			return fmt.Errorf("team (%s) still exists", id)
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasTeamStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		id := rs.Primary.Attributes["team_id"]

		return fmt.Sprintf("%s-%s", orgID, id), nil
	}
}

func testAccMongoDBAtlasTeamConfig(orgID, name string, usernames []string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_teams" "test" {
			org_id     = "%s"
			name       = "%s"
			usernames  = %s
		}`, orgID, name,
		strings.ReplaceAll(fmt.Sprintf("%+q", usernames), " ", ","),
	)
}
