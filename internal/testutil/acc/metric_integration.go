package acc

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// GetMetricIntegration performs an authenticated GET against the metric integration preview endpoint.
func GetMetricIntegration(ctx context.Context, projectID, integrationID string) (*http.Response, error) {
	baseURL := config.NormalizeBaseURL(os.Getenv("MONGODB_ATLAS_BASE_URL"))
	url := fmt.Sprintf("%s/api/atlas/v2/groups/%s/metricIntegrations/%s", baseURL, projectID, integrationID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.atlas.preview+json")
	return ConnV2().GetConfig().HTTPClient.Do(req)
}
