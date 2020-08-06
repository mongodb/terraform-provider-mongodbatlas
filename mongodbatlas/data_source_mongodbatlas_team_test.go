package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMongoDBAtlasTeam_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_teams.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name           = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		username       = "mongodbatlas.testing@gmail.com"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasTeamConfig(orgID, name, username),
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

func TestAccDataSourceMongoDBAtlasTeamByName_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_teams.test2"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name           = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		username       = "mongodbatlas.testing@gmail.com"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasTeamConfigByName(orgID, name, username),
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

	`, orgID, name, username)
}

func testAccDataSourceMongoDBAtlasTeamConfigByName(orgID, name, username string) string {
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
