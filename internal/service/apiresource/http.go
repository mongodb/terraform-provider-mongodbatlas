package apiresource

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// callResult bundles everything the resource needs to inspect after a call.
type callResult struct {
	Parsed   map[string]any
	Err      error
	Raw      []byte
	Status   int
	NotFound bool
}

// callAPI performs an HTTP request via the Atlas client.
// If parsed is non-nil, the response body was a JSON object.
func callAPI(ctx context.Context, client *config.MongoDBClient, method, path, versionHeader string, body []byte) callResult {
	params := config.APICallParams{
		Method:        method,
		RelativePath:  path,
		VersionHeader: versionHeader,
	}
	resp, err := client.UntypedAPICall(ctx, params, body)
	result := callResult{}
	if resp != nil {
		result.Status = resp.StatusCode
	}
	if resp != nil && resp.Body != nil {
		result.Raw, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
	}
	result.NotFound = isNotFound(resp, err, result.Raw)
	if err != nil {
		result.Err = err
		return result
	}
	if len(bytes.TrimSpace(result.Raw)) > 0 {
		_ = json.Unmarshal(result.Raw, &result.Parsed)
	}
	return result
}

func isNotFound(resp *http.Response, _ error, raw []byte) bool {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	// Some Atlas endpoints answer 2xx with empty body when the entity is gone.
	if resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		trimmed := bytes.TrimSpace(raw)
		if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("{}")) {
			return true
		}
	}
	return false
}
