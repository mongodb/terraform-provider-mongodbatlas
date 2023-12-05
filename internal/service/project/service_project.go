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

func ServiceFromClient(client *admin.APIClient) GroupProjectService {
	return &GroupProjectServiceFromClient{
		client: client,
	}
}
