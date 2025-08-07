package clouduserteamassignment_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

const (
	resourceTeamName           = "mongodbatlas_team.test"
	resourceTeamAssignmentName = "mongodbatlas_cloud_user_team_assignment.test"
)

func TestMigCloudUserTeamAssignmentRS_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0") // when resource 1st released
	mig.CreateAndRunTest(t, basicTestCase(t))
}

func TestMigCloudUserTeamAssignmentRS_migrationJourney(t *testing.T) {
	var (
		orgID    = os.Getenv("MONGODB_ATLAS_ORG_ID")
		teamName = fmt.Sprintf("team-test-%s", acc.RandomName())
		username = os.Getenv("MONGODB_ATLAS_USERNAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t); acc.PreCheckAtlasUsername(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				// NOTE: 'usernames' attribute (available v1.39.0, deprecated in v2.0.0) is used in this test in team resource,
				// which may be removed in future versions. This could cause the test to break - keep for version tracking.
				ExternalProviders: mig.ExternalProviders(),
				Config:            configTeamWithUsernamesFirst(orgID, teamName, username),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configWithTeamAssignmentsSecond(orgID, teamName, username), // expected to see 1 import in the plan
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceTeamName, "name", teamName),

					resource.TestCheckResourceAttrSet(resourceTeamAssignmentName, "user_id"),
					resource.TestCheckResourceAttr(resourceTeamAssignmentName, "username", username),
				),
			},
			mig.TestStepCheckEmptyPlan(configWithTeamAssignmentsSecond(orgID, teamName, username)),
		},
	})
}

// Step 1: Original configuration with usernames attribute
func configTeamWithUsernamesFirst(orgID, teamName, username string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_team" "test" {
		org_id    = %[2]q
		name      = %[3]q
		usernames = [%[1]q]
	}
	`, username, orgID, teamName)
}

// Step 2: Configuration with team assignments using import blocks

// NOTE: Using static resource assignment instead of for_each with multiple usernames
// due to a known limitation in Terraform's acceptance testing framework with indexed resources.
// The actual migration using for_each works correctly (verified locally).
func configWithTeamAssignmentsSecond(orgID, teamName, username string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_team" "test" {
		org_id    = %[1]q
		name      = %[2]q
	}

	data "mongodbatlas_team" "test" {
		org_id = %[1]q
		team_id = mongodbatlas_team.test.team_id
	}
 
	resource "mongodbatlas_cloud_user_team_assignment" "test" {
		org_id   = %[1]q
		team_id  = mongodbatlas_team.test.team_id
		user_id  = data.mongodbatlas_team.test.users[0].id
	}

	import {
		to = mongodbatlas_cloud_user_team_assignment.test
		id = "%[1]s/${mongodbatlas_team.test.team_id}/${data.mongodbatlas_team.test.users[0].id}"
	}

	`, orgID, teamName)
}
