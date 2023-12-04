package project

import (
	"context"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

type GroupProjectService interface {
	UpdateProject(ctx context.Context, groupID string, groupName *admin.GroupName) (*admin.Group, *http.Response, error)
}

type GroupProjectServiceFromClient struct {
	client *admin.APIClient
}

func (a *GroupProjectServiceFromClient) UpdateProject(ctx context.Context, groupID string, groupName *admin.GroupName) (*admin.Group, *http.Response, error) {
	return a.client.ProjectsApi.UpdateProject(ctx, groupID, groupName).Execute()
}

func ServiceFromClient(client *admin.APIClient) GroupProjectService {
	return &GroupProjectServiceFromClient{
		client: client,
	}
}
