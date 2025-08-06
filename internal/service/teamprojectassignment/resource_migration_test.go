package teamprojectassignment_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

const (
	resourceProjectName     = "mongodbatlas_project.migration_path_project1"
	resourceAssignmentName1 = "mongodbatlas_team_project_assignment.team1"
	resourceAssignmentName2 = "mongodbatlas_team_project_assignment.team2"
)

func TestMigCloudUserTeamAssignmentRS_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0") // when resource 1st released
	mig.CreateAndRunTest(t, basicTestCase(t))
}

func TestMigTeamProjectAssignment_migrationJourney(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		teamName1   = acc.RandomName()
		teamName2   = acc.RandomName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1: Create project with `teams`
				ExternalProviders: mig.ExternalProviders(),
				Config:            originalConfigFirst(projectName, orgID, teamName1, teamName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceProjectName, "name", projectName),
					resource.TestCheckResourceAttr(resourceProjectName, "teams.#", "2"),
				),
			},
			{
				// Step 2: Ignore `teams` attribute & import new resource
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   ignoreTeamsImportConfigSecond(projectName, orgID, teamName1, teamName2), // expected to see 2 import in the plan
				Check:                    secondChecks(),
			},
			mig.TestStepCheckEmptyPlan(ignoreTeamsImportConfigSecond(projectName, orgID, teamName1, teamName2)),
		},
	})
}

func originalConfigFirst(projectName, orgID, teamName1, teamName2 string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_team" "team1" {
			name   = %[1]q
			org_id = %[3]q
			usernames = []
		}
				
		resource "mongodbatlas_team" "team2" {
			name   = %[2]q
			org_id = %[3]q
			usernames = []
		}

		locals {
			team_map = {
				(mongodbatlas_team.team1.team_id) = ["GROUP_OWNER"]
				(mongodbatlas_team.team2.team_id) = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_WRITE"]
			}
		}

		resource "mongodbatlas_project" "migration_path_project1" {
			name   = %[4]q
			org_id = %[3]q

			dynamic "teams" {
				for_each = local.team_map
				content {
				team_id = teams.key
				role_names = teams.value
				}
				
			}
		}`, teamName1, teamName2, orgID, projectName)
}

func ignoreTeamsImportConfigSecond(projectName, orgID, teamName1, teamName2 string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_team" "team1" {
			name   = %[1]q
			org_id = %[3]q
		}

		resource "mongodbatlas_team" "team2" {
			name   = %[2]q
			org_id = %[3]q
		}

		locals {
			team_map = {
				(mongodbatlas_team.team1.team_id) = ["GROUP_OWNER"]
				(mongodbatlas_team.team2.team_id) = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_WRITE"]
			}
		}

		resource "mongodbatlas_project" "migration_path_project1" {
			name   = %[4]q
			org_id = %[3]q

			lifecycle {
				ignore_changes = [teams]
			}
		}

		resource "mongodbatlas_team_project_assignment" "team1" {  
			project_id = mongodbatlas_project.migration_path_project1.id  
			team_id    = mongodbatlas_team.team1.team_id  
			role_names = local.team_map[mongodbatlas_team.team1.team_id]  
		}

		import {  
			to = mongodbatlas_team_project_assignment.team1  
			id = "${mongodbatlas_project.migration_path_project1.id}/${mongodbatlas_team.team1.team_id}"  
		}

		resource "mongodbatlas_team_project_assignment" "team2" {  
			project_id = mongodbatlas_project.migration_path_project1.id  
			team_id    = mongodbatlas_team.team2.team_id  
			role_names = local.team_map[mongodbatlas_team.team2.team_id]  
		}
		
		import {  
			to = mongodbatlas_team_project_assignment.team2  
			id = "${mongodbatlas_project.migration_path_project1.id}/${mongodbatlas_team.team2.team_id}"  
		}
		`, teamName1, teamName2, orgID, projectName)
}

func secondChecks() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceAssignmentName1, "project_id"),
		resource.TestCheckResourceAttrSet(resourceAssignmentName2, "project_id"),
		resource.TestCheckResourceAttrSet(resourceAssignmentName1, "team_id"),
		resource.TestCheckResourceAttrSet(resourceAssignmentName2, "team_id"),
		resource.TestCheckResourceAttr(resourceAssignmentName1, "role_names.#", "1"),
		resource.TestCheckResourceAttr(resourceAssignmentName2, "role_names.#", "2"),
	)
}
