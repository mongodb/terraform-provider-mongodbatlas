package team_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
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

func projectsPage(totalCount int, ids ...string) *admin20241113.PaginatedAtlasGroup {
	results := make([]admin20241113.Group, len(ids))
	for i := range ids {
		results[i] = admin20241113.Group{Id: &ids[i]}
	}
	return &admin20241113.PaginatedAtlasGroup{Results: &results, TotalCount: &totalCount}
}

func teamsPage(teamIDs ...string) *admin20241113.PaginatedTeamRole {
	results := make([]admin20241113.TeamRole, len(teamIDs))
	for i := range teamIDs {
		results[i] = admin20241113.TeamRole{TeamId: &teamIDs[i]}
	}
	return &admin20241113.PaginatedTeamRole{Results: &results}
}

func expectProjectTeams(t *testing.T, teamsAPI *mockadmin20241113.TeamsApi, projectID string, teamIDs ...string) {
	t.Helper()
	projectTeams := mockadmin20241113.NewTeamsApi(t)
	teamsAPI.EXPECT().ListProjectTeams(mock.Anything, projectID).Return(admin20241113.ListProjectTeamsApiRequest{ApiService: projectTeams}).Once()
	projectTeams.EXPECT().ListProjectTeamsExecute(mock.Anything).Return(teamsPage(teamIDs...), nil, nil).Once()
}

func TestGetProjectIDByTeamID(t *testing.T) {
	testCases := map[string]struct {
		expectations func(*testing.T, *mockadmin20241113.OrganizationsApi, *mockadmin20241113.TeamsApi)
		expectedID   string
		expectedErr  string
	}{
		"returns an error when no project has the team assigned": {
			expectations: func(t *testing.T, orgsAPI *mockadmin20241113.OrganizationsApi, teamsAPI *mockadmin20241113.TeamsApi) {
				t.Helper()
				orgsAPI.EXPECT().ListOrganizationProjects(mock.Anything, lookupOrgID).Return(admin20241113.ListOrganizationProjectsApiRequest{ApiService: orgsAPI}).Once()
				orgsAPI.EXPECT().ListOrganizationProjectsExecute(mock.Anything).Return(projectsPage(1, "project-1"), nil, nil).Once()
				expectProjectTeams(t, teamsAPI, "project-1", "other-team")
			},
			expectedErr: "no project in organization (org-id) has it assigned",
		},
		"returns an error when listing projects fails": {
			expectations: func(t *testing.T, orgsAPI *mockadmin20241113.OrganizationsApi, teamsAPI *mockadmin20241113.TeamsApi) {
				t.Helper()
				orgsAPI.EXPECT().ListOrganizationProjects(mock.Anything, lookupOrgID).Return(admin20241113.ListOrganizationProjectsApiRequest{ApiService: orgsAPI}).Once()
				orgsAPI.EXPECT().ListOrganizationProjectsExecute(mock.Anything).Return(nil, nil, errors.New("list failed")).Once()
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

// TestGetProjectIDByTeamIDReadsAllProjectPages uses a fake Atlas server that answers by the pageNum query
// parameter, so it fails if the lookup stops requesting later pages.
func TestGetProjectIDByTeamIDReadsAllProjectPages(t *testing.T) {
	var (
		mu             sync.Mutex
		requestedPages []string
	)
	projectPages := map[string]string{"1": "project-1", "2": "project-2"}
	projectTeams := map[string]string{"project-1": "other-team", "project-2": lookupTeamID}

	writeJSON := func(w http.ResponseWriter, body any) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(body))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/atlas/v2/orgs/"+lookupOrgID+"/groups", func(w http.ResponseWriter, r *http.Request) {
		pageNum := r.URL.Query().Get("pageNum")
		mu.Lock()
		requestedPages = append(requestedPages, pageNum)
		mu.Unlock()
		projectID, ok := projectPages[pageNum]
		if !ok {
			projectID = projectPages["1"] // Atlas returns the first page when pageNum is missing.
		}
		writeJSON(w, map[string]any{"results": []map[string]string{{"id": projectID}}, "totalCount": len(projectPages)})
	})
	mux.HandleFunc("GET /api/atlas/v2/groups/{projectID}/teams", func(w http.ResponseWriter, r *http.Request) {
		teamID := projectTeams[r.PathValue("projectID")]
		writeJSON(w, map[string]any{"results": []map[string]string{{"teamId": teamID}}, "totalCount": 1})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	client, err := admin20241113.NewClient(admin20241113.UseBaseURL(srv.URL))
	require.NoError(t, err)

	projectID, err := team.GetProjectIDByTeamID(t.Context(), client.OrganizationsApi, client.TeamsApi, lookupOrgID, lookupTeamID)
	require.NoError(t, err)
	assert.Equal(t, "project-2", projectID)
	assert.Equal(t, []string{"1", "2"}, requestedPages)
}
