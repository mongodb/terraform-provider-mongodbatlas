package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSTeam_basic(t *testing.T) {
	var (
		team         matlas.Team
		resourceName = "mongodbatlas_teams.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name         = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		updatedName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		username     = os.Getenv("MONGODB_ATLAS_USERNAME_CLOUD_DEV")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, name,
					[]string{
						username,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasTeamExists(resourceName, &team),
					testAccCheckMongoDBAtlasTeamAttributes(&team, name),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, updatedName,
					[]string{
						"edgar.lopez@digitalonus.com",
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
						"edgar.lopez@digitalonus.com",
						"mongodbatlas.testing@gmail.com",
						"antonio.cabrera@digitalonus.com",
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

func TestAccConfigRSTeam_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_teams.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username     = os.Getenv("MONGODB_ATLAS_USERNAME_CLOUD_DEV")
		name         = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasTeamConfig(orgID, name, []string{username}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "usernames.#"),
					resource.TestCheckResourceAttrSet(resourceName, "team_id"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
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
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		orgID := ids["org_id"]
		id := ids["id"]

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
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_teams" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		orgID := ids["org_id"]
		id := ids["id"]

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
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["org_id"], rs.Primary.Attributes["team_id"]), nil
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
