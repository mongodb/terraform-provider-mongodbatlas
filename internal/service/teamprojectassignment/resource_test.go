package teamprojectassignment_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var resourceName = "mongodbatlas_team_project_assignment.test"
var dataSourceName = "data.mongodbatlas_team_project_assignment.test"

func TestAccTeamProjectAssignment_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func TestAccTeamProjectAssignmentDS_error(t *testing.T) {
	resource.ParallelTest(t, *errorTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectID := acc.ProjectIDExecution(t)
	teamName := acc.RandomName()
	roles := []string{"GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"}
	updatedRoles := []string{"GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN", "GROUP_DATA_ACCESS_READ_ONLY"}

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
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
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
			},
		},
	}
}

func errorTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	projectID := acc.ProjectIDExecution(t)

	return &resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configError1(projectID),
				ExpectError: regexp.MustCompile("The argument \"team_id\" is required"),
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

		data "mongodbatlas_team_project_assignment" "test" {
			project_id = %[3]q
			team_id    = mongodbatlas_team_project_assignment.test.team_id
		}
	
	`, orgID, teamName, projectID, rolesStr)
}

func configError1(projectID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_team_project_assignment" "test" {
			project_id = %[1]q
		}
	`, projectID)
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

	dataCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
		resource.TestCheckResourceAttrSet(dataSourceName, "team_id"),
		resource.TestCheckResourceAttr(dataSourceName, "role_names.#", fmt.Sprint(len(roles))),
		resource.TestCheckResourceAttrPair(dataSourceName, "team_id", resourceName, "team_id"),
	}
	checkFuncs = append(checkFuncs, dataCheckFuncs...)
	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func importStateIDFunc(resourceName string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		attrs := s.RootModule().Resources[resourceName].Primary.Attributes
		teamID := attrs["team_id"]
		projectID := attrs["project_id"]
		return projectID + "/" + teamID, nil
	}
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
