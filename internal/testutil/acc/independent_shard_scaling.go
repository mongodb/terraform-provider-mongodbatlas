package acc

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func GetIndependentShardScalingMode(ctx context.Context, projectID, clusterName string) (*string, *http.Response, error) {
	baseURL := config.NormalizeBaseURL(os.Getenv("MONGODB_ATLAS_BASE_URL"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/test/utils/auth/groups/"+projectID+"/clusters/"+clusterName+"/independentShardScalingMode", http.NoBody)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "*/*")
	resp, err := ConnV2().GetConfig().HTTPClient.Do(req)
	if err != nil || resp == nil {
		return nil, resp, err
	}

	var result *string
	result, err = decode(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

func decode(body io.ReadCloser) (*string, error) {
	buf, err := io.ReadAll(body)
	_ = body.Close()
	if err != nil {
		return nil, err
	}
	result := string(buf)
	return &result, nil
}
