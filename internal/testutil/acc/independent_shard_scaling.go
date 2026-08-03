package acc

import (
	"context"
	"fmt"
	"net/http"
)

// Deprecated AutoScalingConfiguration is the only Atlas API that returns autoScalingMode; used here
// (test-only) to verify independent shard scaling—no supported replacement exists.
func GetIndependentShardScalingMode(ctx context.Context, projectID, clusterName string) (*string, *http.Response, error) {
	cfg, resp, err := ConnV2().ClustersAPI.AutoScalingConfiguration(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, resp, err
	}
	if cfg == nil {
		return nil, resp, fmt.Errorf("auto scaling configuration is nil for cluster %s in project %s", clusterName, projectID)
	}
	mode := cfg.GetAutoScalingMode()
	return &mode, resp, nil
}
