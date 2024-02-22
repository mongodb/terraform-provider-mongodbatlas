package searchdeployment

import (
	"context"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

type DeploymentService interface {
	GetAtlasSearchDeployment(ctx context.Context, groupID, clusterName string) (*admin.ApiSearchDeploymentResponse, *http.Response, error)
}

type DeploymentServiceFromClient struct {
	client *admin.APIClient
}

func (a *DeploymentServiceFromClient) GetAtlasSearchDeployment(ctx context.Context, groupID, clusterName string) (*admin.ApiSearchDeploymentResponse, *http.Response, error) {
	return a.client.AtlasSearchApi.GetAtlasSearchDeployment(ctx, groupID, clusterName).Execute()
}

func ServiceFromClient(client *admin.APIClient) DeploymentService {
	return &DeploymentServiceFromClient{
		client: client,
	}
}
