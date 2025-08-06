package clouduserteamassignment_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigCloudUserTeamAssignmentRS_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0") // when resource 1st released
	mig.CreateAndRunTest(t, basicTestCase(t))
}

func TestMigCloudUserTeamAssignmentRS_migrationJourney(t *testing.T) {
	var (
		orgID     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		teamName  = fmt.Sprintf("team-test-%s", acc.RandomName())
		usernames = []string{os.Getenv("MONGODB_ATLAS_USERNAME")}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				// NOTE: 'usernames' attribute (available v1.39.0, deprecated in v2.0.0) is used in this test in team resource.
				// May be removed in future versions.
				ExternalProviders: mig.ExternalProviders(),
				Config:            configTeamWithUsernamesFirst(orgID, teamName, usernames),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configWithTeamAssignmentsSecond(orgID, teamName, usernames),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodbatlas_team.test", "name", teamName),

					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_user_team_assignment.test", "user_id"),
					resource.TestCheckResourceAttr("mongodbatlas_cloud_user_team_assignment.test", "username", usernames[0]),
				),
			},
			mig.TestStepCheckEmptyPlan(configWithTeamAssignmentsSecond(orgID, teamName, usernames)),
		},
	})
}

// Step 1: Original configuration with usernames attribute
func configTeamWithUsernamesFirst(orgID, teamName string, usernames []string) string {
	return fmt.Sprintf(`
	locals {
		usernames = [%[1]q]
	}

	resource "mongodbatlas_team" "test" {
	org_id    = %[2]q
	name      = %[3]q
	usernames = local.usernames
	}
	`, usernames[0], orgID, teamName)
}

// Step 2: Configuration with team assignments using import blocks
func configWithTeamAssignmentsSecond(orgID, teamName string, usernames []string) string {
	return fmt.Sprintf(`
	locals {
		usernames = [%[1]q]
	}

	resource "mongodbatlas_team" "test" {
		org_id    = %[3]q
		name      = %[4]q
	}

	data "mongodbatlas_team" "test" {
		org_id = %[5]q  
		team_id = mongodbatlas_team.test.team_id
	}
 
	resource "mongodbatlas_cloud_user_team_assignment" "test" {
		org_id   = %[3]q
		team_id  = mongodbatlas_team.test.team_id
		user_id  = data.mongodbatlas_team.test.users[0].id
	}

	import {
		to = mongodbatlas_cloud_user_team_assignment.test
		id = "%[3]s/${mongodbatlas_team.test.team_id}/${data.mongodbatlas_team.test.users[0].id}"
	}

	`, usernames, orgID, orgID, teamName, orgID)
}
