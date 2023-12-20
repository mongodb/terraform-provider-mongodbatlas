package advancedcluster

import (
	"context"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

type ClusterService interface {
	Get(ctx context.Context, groupID, clusterName string) (*matlas.Cluster, *matlas.Response, error)
	List(ctx context.Context, groupID string, options *matlas.ListOptions) (*matlas.AdvancedClustersResponse, *matlas.Response, error)
	GetAdvancedCluster(ctx context.Context, groupID, clusterName string) (*matlas.AdvancedCluster, *matlas.Response, error)
}

type ClusterServiceFromClient struct {
	client *matlas.Client
}

func (a *ClusterServiceFromClient) Get(ctx context.Context, groupID, clusterName string) (*matlas.Cluster, *matlas.Response, error) {
	return a.client.Clusters.Get(ctx, groupID, clusterName)
}

func (a *ClusterServiceFromClient) GetAdvancedCluster(ctx context.Context, groupID, clusterName string) (*matlas.AdvancedCluster, *matlas.Response, error) {
	return a.client.AdvancedClusters.Get(ctx, groupID, clusterName)
}

func (a *ClusterServiceFromClient) List(ctx context.Context, groupID string, options *matlas.ListOptions) (*matlas.AdvancedClustersResponse, *matlas.Response, error) {
	return a.client.AdvancedClusters.List(ctx, groupID, options)
}

func ServiceFromClient(client *matlas.Client) ClusterService {
	return &ClusterServiceFromClient{
		client: client,
	}
}
