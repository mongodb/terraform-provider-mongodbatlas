package team_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigRSTeam_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_team.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		usernames    = []string{os.Getenv("MONGODB_ATLAS_USERNAME")}
		name         = acc.RandomName()
		updatedName  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name, usernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
			{
				Config: configBasic(orgID, updatedName, usernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
			{
				Config: configBasic(orgID, updatedName, usernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
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

func TestAccConfigRSTeam_updatingUsernames(t *testing.T) {
	var (
		resourceName          = "mongodbatlas_team.test"
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		firstUser             = os.Getenv("MONGODB_ATLAS_USERNAME")
		secondUser            = os.Getenv("MONGODB_ATLAS_USERNAME_2")
		usernames             = []string{firstUser}
		updatedSingleUsername = []string{secondUser}
		updatedBothUsername   = []string{firstUser, secondUser}
		name                  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAtlasUsernames(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name, usernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "usernames.*", usernames[0]),
				),
			},
			{
				Config: configBasic(orgID, name, updatedSingleUsername),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "usernames.*", updatedSingleUsername[0]),
				),
			},
			{
				Config: configBasic(orgID, name, updatedBothUsername),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "usernames.*", updatedBothUsername[0]),
					resource.TestCheckTypeSetElemAttr(resourceName, "usernames.*", updatedBothUsername[1]),
				),
			},
		},
	})
}

func TestAccConfigRSTeam_legacyName(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_teams.test"
		dataSourceName = "data.mongodbatlas_teams.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		usernames      = []string{os.Getenv("MONGODB_ATLAS_USERNAME")}
		name           = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: configBasicLegacyNames(orgID, name, usernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "usernames.#", "1"),
				),
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
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
		_, _, err := acc.ConnV2().TeamsApi.GetGroupTeam(context.Background(), orgID, id).Execute()
		if err == nil {
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
		resource "mongodbatlas_team" "test" {
			org_id     = "%s"
			name       = "%s"
			usernames  = %s
		}`, orgID, name,
		strings.ReplaceAll(fmt.Sprintf("%+q", usernames), " ", ","),
	)
}

func configBasicLegacyNames(orgID, name string, usernames []string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_teams" "test" {
			org_id     = %[1]q
			name       = %[2]q
			usernames  = %[3]s
		}
		
		data "mongodbatlas_teams" "test" {
			org_id = %[1]q
			name   = mongodbatlas_teams.test.name
		}
		`, orgID, name,
		strings.ReplaceAll(fmt.Sprintf("%+q", usernames), " ", ","),
	)
}
