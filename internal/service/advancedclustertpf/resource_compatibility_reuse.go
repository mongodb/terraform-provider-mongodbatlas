package advancedclustertpf

import (
	"context"
	"fmt"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
)

func getReplicationSpecIDsFromOldAPI(ctx context.Context, projectID, clusterName string, api admin20240530.ClustersApi) (map[string]string, error) {
	clusterOldAPI, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
	if apiError, ok := admin20240530.AsError(err); ok {
		if apiError.GetErrorCode() == "ASYMMETRIC_SHARD_UNSUPPORTED" {
			return nil, nil // if its the case of an asymmetric shard an error is expected in old API, replication_specs.*.id attribute will not be populated
		}
		readErrorMsg := "error reading  advanced cluster with 2023-02-01 API (%s): %s"
		return nil, fmt.Errorf(readErrorMsg, clusterName, err)
	}
	specs := clusterOldAPI.GetReplicationSpecs()
	result := make(map[string]string, len(specs))
	for _, spec := range specs {
		result[spec.GetZoneName()] = spec.GetId()
	}
	return result, nil
}
