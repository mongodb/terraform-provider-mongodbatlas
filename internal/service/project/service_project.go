package project

import (
	"context"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

type GroupProjectService interface {
	UpdateProject(ctx context.Context, groupID string, groupName *admin.GroupName) (*admin.Group, *http.Response, error)
	ListProjectTeams(ctx context.Context, groupID string) (*admin.PaginatedTeamRole, *http.Response, error)
	GetProjectSettings(ctx context.Context, groupID string) (*admin.GroupSettings, *http.Response, error)
	DeleteProjectLimit(ctx context.Context, limitName, projectID string) (map[string]interface{}, *http.Response, error)
	SetProjectLimit(ctx context.Context, limitName, groupID string, dataFederationLimit *admin.DataFederationLimit) (*admin.DataFederationLimit, *http.Response, error)
	ListProjectLimits(ctx context.Context, groupID string) ([]admin.DataFederationLimit, *http.Response, error)
	RemoveProjectTeam(ctx context.Context, groupID, teamID string) (*http.Response, error)
	UpdateTeamRoles(ctx context.Context, groupID, teamID string, teamRole *admin.TeamRole) (*admin.PaginatedTeamRole, *http.Response, error)
	AddAllTeamsToProject(ctx context.Context, groupID string, teamRole *[]admin.TeamRole) (*admin.PaginatedTeamRole, *http.Response, error)
	ListClusters(ctx context.Context, groupID string) (*admin.PaginatedAdvancedClusterDescription, *http.Response, error)
}

type GroupProjectServiceFromClient struct {
	client *admin.APIClient
}

func (a *GroupProjectServiceFromClient) UpdateProject(ctx context.Context, groupID string, groupName *admin.GroupName) (*admin.Group, *http.Response, error) {
	return a.client.ProjectsApi.UpdateProject(ctx, groupID, groupName).Execute()
}

func (a *GroupProjectServiceFromClient) ListProjectLimits(ctx context.Context, groupID string) ([]admin.DataFederationLimit, *http.Response, error) {
	return a.client.ProjectsApi.ListProjectLimits(ctx, groupID).Execute()
}

func (a *GroupProjectServiceFromClient) GetProjectSettings(ctx context.Context, groupID string) (*admin.GroupSettings, *http.Response, error) {
	return a.client.ProjectsApi.GetProjectSettings(ctx, groupID).Execute()
}

func (a *GroupProjectServiceFromClient) DeleteProjectLimit(ctx context.Context, limitName, projectID string) (map[string]interface{}, *http.Response, error) {
	return a.client.ProjectsApi.DeleteProjectLimit(ctx, limitName, projectID).Execute()
}

func (a *GroupProjectServiceFromClient) SetProjectLimit(ctx context.Context, limitName, groupID string,
	dataFederationLimit *admin.DataFederationLimit) (*admin.DataFederationLimit, *http.Response, error) {
	return a.client.ProjectsApi.SetProjectLimit(ctx, limitName, groupID, dataFederationLimit).Execute()
}

func (a *GroupProjectServiceFromClient) ListProjectTeams(ctx context.Context, groupID string) (*admin.PaginatedTeamRole, *http.Response, error) {
	return a.client.TeamsApi.ListProjectTeams(ctx, groupID).Execute()
}

func (a *GroupProjectServiceFromClient) RemoveProjectTeam(ctx context.Context, groupID, teamID string) (*http.Response, error) {
	return a.client.TeamsApi.RemoveProjectTeam(ctx, groupID, teamID).Execute()
}

func (a *GroupProjectServiceFromClient) UpdateTeamRoles(ctx context.Context, groupID, teamID string, teamRole *admin.TeamRole) (*admin.PaginatedTeamRole, *http.Response, error) {
	return a.client.TeamsApi.UpdateTeamRoles(ctx, groupID, teamID, teamRole).Execute()
}

func (a *GroupProjectServiceFromClient) AddAllTeamsToProject(ctx context.Context, groupID string, teamRole *[]admin.TeamRole) (*admin.PaginatedTeamRole, *http.Response, error) {
	return a.client.TeamsApi.AddAllTeamsToProject(ctx, groupID, teamRole).Execute()
}

func (a *GroupProjectServiceFromClient) ListClusters(ctx context.Context, groupID string) (*admin.PaginatedAdvancedClusterDescription, *http.Response, error) {
	return a.client.ClustersApi.ListClusters(ctx, groupID).Execute()
}

// users
func (a *GroupProjectServiceFromClient) AddUserToProject(ctx context.Context, groupID string, groupInvitationRequest *admin.GroupInvitationRequest) (*admin.OrganizationInvitation, *http.Response, error) {
	return a.client.ProjectsApi.AddUserToProject(ctx, groupID, groupInvitationRequest).Execute()
}
func (a *GroupProjectServiceFromClient) UpdateProjectUserRoles(ctx context.Context, groupID, userID string, updateGroupRolesForUser *admin.UpdateGroupRolesForUser) (*admin.UpdateGroupRolesForUser, *http.Response, error) {
	return a.client.ProjectsApi.UpdateProjectRoles(ctx, groupID, userID, updateGroupRolesForUser).Execute()
}
func (a *GroupProjectServiceFromClient) ListProjectUsers(ctx context.Context, groupID string) (*admin.PaginatedAppUser, *http.Response, error) {
	return a.client.ProjectsApi.ListProjectUsers(ctx, groupID).Execute()
}
func (a *GroupProjectServiceFromClient) RemoveProjectUser(ctx context.Context, groupID, userID string) (*http.Response, error) {
	return a.client.ProjectsApi.RemoveProjectUser(ctx, groupID, userID).Execute()
}

// org invitations
func (a *GroupProjectServiceFromClient) GetOrganizationInvitation(ctx context.Context, orgID string, invitationID string) (*admin.OrganizationInvitation, *http.Response, error) {
	return a.client.OrganizationsApi.GetOrganizationInvitation(ctx, orgID, invitationID).Execute()
}

// TODO confirm if this is required
func (a *GroupProjectServiceFromClient) UpdateOrganizationInvitation(ctx context.Context, orgId string, organizationInvitationRequest *admin.OrganizationInvitationRequest) (*admin.OrganizationInvitation, *http.Response, error) {
	return a.client.OrganizationsApi.UpdateOrganizationInvitation(ctx, orgId, organizationInvitationRequest).Execute()
}

func ServiceFromClient(client *admin.APIClient) GroupProjectService {
	return &GroupProjectServiceFromClient{
		client: client,
	}
}
