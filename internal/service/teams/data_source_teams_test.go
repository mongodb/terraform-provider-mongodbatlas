package teams_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSTeam_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_teams.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name           = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		username       = os.Getenv("MONGODB_ATLAS_USERNAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasic(orgID, name, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "team_id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "usernames.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigDSTeamByName_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_teams.test2"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name           = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		username       = os.Getenv("MONGODB_ATLAS_USERNAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasicByName(orgID, name, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "team_id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "usernames.#", "1"),
				),
			},
		},
	})
}

func dataSourceConfigBasic(orgID, name, username string) string {
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

	`, orgID, name, username)
}

func dataSourceConfigBasicByName(orgID, name, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_teams" "test" {
			org_id     = "%s"
			name       = "%s"
			usernames  = ["%s"]
		}

		data "mongodbatlas_teams" "test2" {
			org_id     = mongodbatlas_teams.test.org_id
			name    = mongodbatlas_teams.test.name
		}
	`, orgID, name, username)
}
