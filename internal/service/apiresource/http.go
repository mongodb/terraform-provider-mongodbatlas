package apiresource

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
		// Atlas Admin API endpoints return JSON objects at the top level. If the
		// body is non-empty but does not unmarshal into a map, surface a clear
		// error rather than silently producing a nil Parsed (which would break
		// id_attribute derivation and look like empty output).
		if err := json.Unmarshal(result.Raw, &result.Parsed); err != nil || result.Parsed == nil {
			var probe any
			if jsonErr := json.Unmarshal(result.Raw, &probe); jsonErr == nil {
				result.Err = fmt.Errorf("expected JSON object at top level, got %T (raw: %s)", probe, truncateForError(result.Raw))
			} else if err != nil {
				result.Err = fmt.Errorf("response is not valid JSON: %w (raw: %s)", err, truncateForError(result.Raw))
			}
		}
	}
	return result
}

func truncateForError(b []byte) string {
	const maxLen = 200
	if len(b) <= maxLen {
		return string(b)
	}
	return string(b[:maxLen]) + "..."
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
