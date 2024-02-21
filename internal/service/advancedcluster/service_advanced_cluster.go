package advancedcluster

import (
	"context"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

type ClusterService interface {
	Get(ctx context.Context, groupID, clusterName string) (*admin.AdvancedClusterDescription, *http.Response, error)
	List(ctx context.Context, options *admin.ListClustersApiParams) (*admin.PaginatedAdvancedClusterDescription, *http.Response, error)
}

type ClusterServiceFromClient struct {
	client *admin.APIClient
}

func (a *ClusterServiceFromClient) Get(ctx context.Context, groupID, clusterName string) (*admin.AdvancedClusterDescription, *http.Response, error) {
	return a.client.ClustersApi.GetCluster(ctx, groupID, clusterName).Execute()
}

func (a *ClusterServiceFromClient) List(ctx context.Context, options *admin.ListClustersApiParams) (*admin.PaginatedAdvancedClusterDescription, *http.Response, error) {
	return a.client.ClustersApi.ListClustersWithParams(ctx, options).Execute()
}

func ServiceFromClient(client *admin.APIClient) ClusterService {
	return &ClusterServiceFromClient{
		client: client,
	}
}
