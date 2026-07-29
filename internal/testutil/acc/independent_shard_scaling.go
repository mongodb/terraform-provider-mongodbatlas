package acc

import (
	"context"
	"net/http"
)

func GetIndependentShardScalingMode(ctx context.Context, projectID, clusterName string) (*string, *http.Response, error) {
	cfg, resp, err := ConnV2().ClustersAPI.AutoScalingConfiguration(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, resp, err
	}
	mode := cfg.GetAutoScalingMode()
	return &mode, resp, nil
}
