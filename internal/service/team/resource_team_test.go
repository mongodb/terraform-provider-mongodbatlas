package team_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
	mockadmin20241113 "go.mongodb.org/atlas-sdk/v20241113005/mockadmin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/team"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigRSTeam_basicNoUsernames(t *testing.T) {
	var (
		resourceName = "mongodbatlas_team.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name         = acc.RandomName()
		updatedName  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: configBasic(orgID, updatedName, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
			{
				Config: configBasic(orgID, updatedName, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
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
				Config: configBasic(orgID, name, &usernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
			{
				Config: configBasic(orgID, updatedName, &usernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
			{
				Config: configBasic(orgID, updatedName, &usernames),
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
				Config: configBasic(orgID, name, &usernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "usernames.*", usernames[0]),
				),
			},
			{
				Config: configBasic(orgID, name, &updatedSingleUsername),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "usernames.*", updatedSingleUsername[0]),
				),
			},
			{
				Config: configBasic(orgID, name, &updatedBothUsername),
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
		_, _, err := acc.ConnV2().TeamsAPI.GetOrgTeam(context.Background(), orgID, id).Execute()
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

func configBasic(orgID, name string, usernames *[]string) string {
	var usernamesAttr string
	if usernames != nil && len(*usernames) > 0 {
		usernamesStr := `"` + strings.Join(*usernames, `", "`) + `"`
		usernamesAttr = fmt.Sprintf(`
  usernames = [%s]`, usernamesStr)
	}

	return fmt.Sprintf(`
resource "mongodbatlas_team" "test" {
  org_id = "%s"
  name   = "%s"
  
  %s
}`, orgID, name, usernamesAttr)
}

const (
	lookupOrgID  = "org-id"
	lookupTeamID = "team-id"
)

// expectProjectsPage expects a request for exactly this page number, so the test fails if the lookup stops paginating.
func expectProjectsPage(t *testing.T, orgsAPI *mockadmin20241113.OrganizationsApi, pageNum, totalCount int, projectIDs []string, err error) {
	t.Helper()
	pageAPI := mockadmin20241113.NewOrganizationsApi(t)
	params := &admin20241113.ListOrganizationProjectsApiParams{OrgId: lookupOrgID, PageNum: &pageNum}
	orgsAPI.EXPECT().ListOrganizationProjectsWithParams(mock.Anything, params).Return(admin20241113.ListOrganizationProjectsApiRequest{ApiService: pageAPI}).Once()
	if err != nil {
		pageAPI.EXPECT().ListOrganizationProjectsExecute(mock.Anything).Return(nil, nil, err).Once()
		return
	}
	projects := make([]admin20241113.Group, len(projectIDs))
	for i := range projectIDs {
		projects[i] = admin20241113.Group{Id: &projectIDs[i]}
	}
	pageAPI.EXPECT().ListOrganizationProjectsExecute(mock.Anything).Return(&admin20241113.PaginatedAtlasGroup{Results: &projects, TotalCount: &totalCount}, nil, nil).Once()
}

func expectProjectTeams(t *testing.T, teamsAPI *mockadmin20241113.TeamsApi, projectID string, teamIDs ...string) {
	t.Helper()
	projectTeamsAPI := mockadmin20241113.NewTeamsApi(t)
	teams := make([]admin20241113.TeamRole, len(teamIDs))
	for i := range teamIDs {
		teams[i] = admin20241113.TeamRole{TeamId: &teamIDs[i]}
	}
	teamsAPI.EXPECT().ListProjectTeams(mock.Anything, projectID).Return(admin20241113.ListProjectTeamsApiRequest{ApiService: projectTeamsAPI}).Once()
	projectTeamsAPI.EXPECT().ListProjectTeamsExecute(mock.Anything).Return(&admin20241113.PaginatedTeamRole{Results: &teams}, nil, nil).Once()
}

func TestGetProjectIDByTeamID(t *testing.T) {
	testCases := map[string]struct {
		expectations func(*testing.T, *mockadmin20241113.OrganizationsApi, *mockadmin20241113.TeamsApi)
		expectedID   string
		expectedErr  string
	}{
		"finds the team in a project on a later page": {
			expectations: func(t *testing.T, orgsAPI *mockadmin20241113.OrganizationsApi, teamsAPI *mockadmin20241113.TeamsApi) {
				t.Helper()
				expectProjectsPage(t, orgsAPI, 1, 2, []string{"project-1"}, nil)
				expectProjectsPage(t, orgsAPI, 2, 2, []string{"project-2"}, nil)
				expectProjectTeams(t, teamsAPI, "project-1", "other-team")
				expectProjectTeams(t, teamsAPI, "project-2", "other-team", lookupTeamID)
			},
			expectedID: "project-2",
		},
		"returns an error when no project has the team assigned": {
			expectations: func(t *testing.T, orgsAPI *mockadmin20241113.OrganizationsApi, teamsAPI *mockadmin20241113.TeamsApi) {
				t.Helper()
				expectProjectsPage(t, orgsAPI, 1, 1, []string{"project-1"}, nil)
				expectProjectTeams(t, teamsAPI, "project-1", "other-team")
			},
			expectedErr: "no project in organization (org-id) has it assigned",
		},
		"returns an error when listing projects fails": {
			expectations: func(t *testing.T, orgsAPI *mockadmin20241113.OrganizationsApi, _ *mockadmin20241113.TeamsApi) {
				t.Helper()
				expectProjectsPage(t, orgsAPI, 1, 0, nil, errors.New("list failed"))
			},
			expectedErr: "error getting projects information: list failed",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			orgsAPI := mockadmin20241113.NewOrganizationsApi(t)
			teamsAPI := mockadmin20241113.NewTeamsApi(t)
			tc.expectations(t, orgsAPI, teamsAPI)

			projectID, err := team.GetProjectIDByTeamID(t.Context(), orgsAPI, teamsAPI, lookupOrgID, lookupTeamID)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedID, projectID)
		})
	}
}
