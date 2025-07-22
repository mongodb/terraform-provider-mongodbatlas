package clouduserteamassignment_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var dataSourceName = "data.mongodbatlas_cloud_user_team_assignment.test"

func TestAccCloudUserTeamAssignmentDS_withUsername(t *testing.T) {
	var (
		orgID    = os.Getenv("MONGODB_ATLAS_ORG_ID")
		teamID   = os.Getenv("MONGODB_ATLAS_TEAM_ID")
		username = os.Getenv("MONGODB_ATLAS_USERNAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configUsername(orgID, teamID, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(dataSourceName, "user_id", acc.IsUserID()),
					resource.TestCheckResourceAttr(dataSourceName, "username", username),
					resource.TestCheckResourceAttr(dataSourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(dataSourceName, "org_id", orgID),
				),
			},
		},
	})
}

func TestAccCloudUserTeamAssignmentDS_withUserID(t *testing.T) {
	var (
		orgID  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		teamID = os.Getenv("MONGODB_ATLAS_TEAM_ID")
		userID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configUserID(orgID, teamID, userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "user_id", userID),
					resource.TestCheckResourceAttrWith(dataSourceName, "username", acc.IsUsername()),
					resource.TestCheckResourceAttr(dataSourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(dataSourceName, "org_id", orgID),
				),
			},
		},
	})
}

func TestAccCloudUserTeamAssignmentDS_multiple(t *testing.T) {
	var (
		orgID  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		teamID = os.Getenv("MONGODB_ATLAS_TEAM_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configError(orgID, teamID),
				ExpectError: regexp.MustCompile("either username or user_id must be provided"),
			},
		},
	})
}

func configUsername(orgID, teamID, username string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cloud_user_team_assignment" "test" {
		org_id   = %q
		team_id  = %q
		username = %q
		}
		`, orgID, teamID, username)
}

func configUserID(orgID, teamID, userID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cloud_user_team_assignment" "test" {
		org_id  = %q
		team_id = %q
		user_id = %q
		}
		`, orgID, teamID, userID)
}

func configError(orgID, teamID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cloud_user_team_assignment" "test" {
		org_id  = %q
		team_id = %q
		}
		`, orgID, teamID)
}
