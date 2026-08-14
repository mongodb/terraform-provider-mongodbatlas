package acc

import (
	"context"
	"fmt"
	"net/http"
)

// GetAIModelRateLimit fetches the current requests/tokens per minute limits for a model group.
func GetAIModelRateLimit(ctx context.Context, projectID, cloud, geography, modelGroupName string) (*http.Response, error) {
	baseURL, err := ConnV2().GetConfig().ServerURLWithContext(ctx, "")
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/api/atlas/v2/groups/%s/aiModelApiClouds/%s/geographies/%s/modelGroupNames/%s/rateLimits",
		baseURL, projectID, cloud, geography, modelGroupName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.atlas.2025-03-12+json")
	return ConnV2().GetConfig().HTTPClient.Do(req)
}
