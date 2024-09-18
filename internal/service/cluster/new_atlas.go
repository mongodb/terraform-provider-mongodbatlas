package cluster

import (
	"context"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
)

func newAtlasUpdate(ctx context.Context, timeout time.Duration, connV2 *admin.APIClient, projectID, clusterName string, redactClientLogData bool) error {
	current, err := newAtlasGet(ctx, connV2, projectID, clusterName)
	if err != nil {
		return err
	}
	if current == redactClientLogData {
		return nil
	}
	req := &admin.ClusterDescription20240805{
		RedactClientLogData: &redactClientLogData,
	}
	if _, _, err = connV2.ClustersApi.UpdateCluster(ctx, projectID, clusterName, req).Execute(); err != nil {
		return err
	}
	stateConf := advancedcluster.CreateStateChangeConfig(ctx, connV2, projectID, clusterName, timeout)
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func newAtlasGet(ctx context.Context, connV2 *admin.APIClient, projectID, clusterName string) (redactClientLogData bool, err error) {
	cluster, _, err := connV2.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	return cluster.GetRedactClientLogData(), err
}

func newAtlasList(ctx context.Context, connV2 *admin.APIClient, projectID string) (map[string]bool, error) {
	clusters, _, err := connV2.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		return nil, err
	}
	results := clusters.GetResults()
	list := make(map[string]bool)
	for i := range results {
		list[results[i].GetName()] = results[i].GetRedactClientLogData()
	}
	return list, nil
}
