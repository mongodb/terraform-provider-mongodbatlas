package teams_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115005/admin"
)

func TestAccConfigRSTeam_basic(t *testing.T) {
	var (
		team         admin.Team
		resourceName = "mongodbatlas_teams.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name         = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		updatedName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		username     = os.Getenv("MONGODB_ATLAS_USERNAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name,
					[]string{
						username,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &team),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
			{
				Config: configBasic(orgID, updatedName,
					[]string{
						username,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &team),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
			{
				Config: configBasic(orgID, updatedName,
					[]string{
						username,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &team),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSTeam_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_teams.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username     = os.Getenv("MONGODB_ATLAS_USERNAME")
		name         = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name, []string{username}),
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
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkExists(resourceName string, team *admin.Team) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		orgID := ids["org_id"]
		id := ids["id"]
		if orgID == "" && id == "" {
			return fmt.Errorf("no ID is set")
		}
		log.Printf("[DEBUG] orgID: %s", orgID)
		log.Printf("[DEBUG] teamID: %s", id)
		teamResp, _, err := acc.ConnV2().TeamsApi.GetTeamById(context.Background(), orgID, id).Execute()
		if err == nil {
			team.Id = teamResp.Id
			team.Name = teamResp.GetName()
			return nil
		}
		return fmt.Errorf("team(%s) does not exist", id)
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["org_id"], rs.Primary.Attributes["team_id"]), nil
	}
}

func configBasic(orgID, name string, usernames []string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_teams" "test" {
			org_id     = "%s"
			name       = "%s"
			usernames  = %s
		}`, orgID, name,
		strings.ReplaceAll(fmt.Sprintf("%+q", usernames), " ", ","),
	)
}
