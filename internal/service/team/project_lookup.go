package team

import (
	"context"
	"fmt"
	"net/http"

	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
)

// GetProjectIDByTeamID returns the ID of a project in the organization that the team is assigned to.
// It reads every page of the organization's projects, so teams assigned to projects beyond the first page are found.
func GetProjectIDByTeamID(ctx context.Context, orgsAPI admin20241113.OrganizationsApi, teamsAPI admin20241113.TeamsApi, orgID, teamID string) (string, error) {
	projects, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin20241113.Group], *http.Response, error) {
		return orgsAPI.ListOrganizationProjectsWithParams(ctx, &admin20241113.ListOrganizationProjectsApiParams{OrgId: orgID, PageNum: &pageNum}).Execute()
	})
	if err != nil {
		return "", fmt.Errorf("error getting projects information: %s", err)
	}

	for _, project := range projects {
		// A project has at most 100 teams, which fits in the default page size.
		teams, _, err := teamsAPI.ListProjectTeams(ctx, project.GetId()).Execute()
		if err != nil {
			return "", fmt.Errorf("error getting teams from project information: %s", err)
		}

		for _, team := range teams.GetResults() {
			if team.GetTeamId() == teamID {
				return project.GetId(), nil
			}
		}
	}

	return "", fmt.Errorf("team (%s) is assigned to a project, but no project in organization (%s) has it assigned", teamID, orgID)
}
