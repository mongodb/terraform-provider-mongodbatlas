package teamprojectassignment_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var resourceName = "mongodbatlas_team_project_assignment.test"

func TestAccTeamProjectAssignmentRS_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	teamName := acc.RandomName()
	roles := []string{"GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"}
	updatedRoles := []string{"GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN", "GROUP_DATA_ACCESS_READ_ONLY"}

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, teamName, projectID, roles),
				Check:  checks(projectID, roles),
			},
			{
				Config: configBasic(orgID, teamName, projectID, updatedRoles),
				Check:  checks(projectID, updatedRoles),
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "team_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceName].Primary.Attributes
					teamID := attrs["team_id"]
					projectID := attrs["project_id"]
					return projectID + "/" + teamID, nil
				},
			},
		},
	}
}

func configBasic(orgID, teamName, projectID string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	return fmt.Sprintf(`
		resource "mongodbatlas_team" "test" {
			org_id     = %[1]q
			name       = %[2]q
		}

		resource "mongodbatlas_team_project_assignment" "test" {
			project_id = %[3]q
			team_id    = mongodbatlas_team.test.team_id
			role_names      = [%[4]s]
		}
	
	`, orgID, teamName, projectID, rolesStr)
}

func checks(projectID string, roles []string) resource.TestCheckFunc {
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
		resource.TestCheckResourceAttrSet(resourceName, "team_id"),
		resource.TestCheckResourceAttr(resourceName, "role_names.#", fmt.Sprint(len(roles))),
	}

	for _, role := range roles {
		checkFuncs = append(checkFuncs, resource.TestCheckTypeSetElemAttr(resourceName, "role_names.*", role))
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_team_project_assignment" {
			continue
		}
		teamID := rs.Primary.Attributes["team_id"]
		projectID := rs.Primary.Attributes["project_id"]
		conn := acc.ConnV2()
		apiListResp, _, err := conn.TeamsApi.ListProjectTeams(context.Background(), projectID).Execute()
		if err != nil {
			continue
		}

		if apiListResp != nil && apiListResp.Results != nil {
			results := *apiListResp.Results
			for i := range results {
				if results[i].GetTeamId() == teamID {
					return fmt.Errorf("team %s still exists", teamID)
				}
			}
		}
	}
	return nil
}
