package team_test

import (
	"errors"
	"testing"

	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
	mockadmin20241113 "go.mongodb.org/atlas-sdk/v20241113005/mockadmin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/team"
)

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
