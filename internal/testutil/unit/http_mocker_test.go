package unit_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractVersion(t *testing.T) {
	version, err := unit.ExtractVersion("application/json; version=2022-06-01")
	require.NoError(t, err)
	require.Equal(t, "2022-06-01", version)
}

func TestExtractVersionRequestResponse(t *testing.T) {
	version := unit.ExtractVersionRequestResponse("application/json;", "application/vnd.atlas.2023-01-01+json;charset=utf-8")
	require.Equal(t, "2023-01-01", version)
}

func TestExtractVersionRequestResponseNotFound(t *testing.T) {
	version := unit.ExtractVersionRequestResponse("application/json;", "application/vnd.atlas.2023-01+json;charset=utf-8")
	require.Equal(t, "", version)
}

func TestRequestInfo_Match(t *testing.T) {
	req := unit.RequestInfo{
		Version: "2022-06-01",
		Method:  "GET",
		Path:    "/v1/cluster/{cluster_id}",
	}
	assert.True(t, req.Match("GET", "/v1/cluster/123", "2022-06-01", map[string]string{"cluster_id": "123"}))
	assert.False(t, req.Match("GET", "/v1/cluster/123", "2022-06-01", map[string]string{"cluster_id": "456"}))
}

func request(method, path, body string) *http.Request {
	req := http.Request{
		Method: method,
		URL:    &url.URL{Path: path},
		Header: http.Header{
			"Accept": []string{"application/json; version=2024-08-05"},
		},
	}
	if body != "" {
		req.Body = io.NopCloser(strings.NewReader(body))
	}
	return &req
}
