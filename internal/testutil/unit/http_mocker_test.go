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
	require.Empty(t, version)
}

func asURL(t *testing.T, reqPath string) *url.URL {
	t.Helper()
	u, err := url.Parse("http://localhost" + reqPath)
	require.NoError(t, err)
	return u
}

func TestRequestInfo_Match(t *testing.T) {
	req := unit.RequestInfo{
		Version: "2022-06-01",
		Method:  "GET",
		Path:    "/v1/cluster/{cluster_id}",
	}
	mockData := unit.MockHTTPData{
		Variables: map[string]string{"cluster_id": "123"},
	}
	assert.True(t, req.Match(t, "GET", "2022-06-01", asURL(t, "/v1/cluster/123"), &mockData)) // Exact match
	mockData2 := unit.MockHTTPData{
		Variables: map[string]string{"cluster_id": "456"},
	}
	// Doesn't match the current request, but adds the new variable mapping cluster_id2=123
	assert.False(t, req.Match(t, "GET", "2022-06-01", asURL(t, "/v1/cluster/123"), &mockData2)) // API Spec match
	assert.Equal(t, map[string]string{"cluster_id": "456", "cluster_id2": "123"}, mockData2.Variables)
}
func TestRequestInfo_MatchQuery(t *testing.T) {
	reqAzure := unit.RequestInfo{
		Version: "2022-06-01",
		Method:  "GET",
		Path:    "/api/atlas/v2/groups/{groupId3}/containers?providerName=AZURE",
	}
	assert.Equal(t, []string{"providerName"}, reqAzure.QueryVars())
	expectedNormalized := "/api/atlas/v2/groups/6746cee66f62fc3c122a3b82/containers?providerName=AZURE"
	reqURLAzure := asURL(t, "/api/atlas/v2/groups/6746cee66f62fc3c122a3b82/containers?includeCount=true&itemsPerPage=100&pageNum=1&providerName=AZURE")
	assert.Equal(t, expectedNormalized, reqAzure.NormalizePath(reqURLAzure))
	mockData := unit.MockHTTPData{
		Variables: map[string]string{"groupId3": "6746cee66f62fc3c122a3b82"},
	}
	assert.True(t, reqAzure.Match(t, "GET", "2022-06-01", reqURLAzure, &mockData))

	assert.Equal(t, map[string]string{"groupId3": "6746cee66f62fc3c122a3b82"}, mockData.Variables)
	reqURLAws := asURL(t, "/api/atlas/v2/groups/6746cee66f62fc3c122a3b82/containers?includeCount=true&itemsPerPage=100&pageNum=1&providerName=AWS")
	assert.False(t, reqAzure.Match(t, "GET", "2022-06-01", reqURLAws, &mockData))
	assert.Equal(t, map[string]string{"groupId3": "6746cee66f62fc3c122a3b82"}, mockData.Variables)

	reqAws := unit.RequestInfo{
		Version: "2022-06-01",
		Method:  "GET",
		Path:    "/api/atlas/v2/groups/{groupId3}/containers?providerName=AWS",
	}
	assert.True(t, reqAws.Match(t, "GET", "2022-06-01", reqURLAws, &mockData))
}

func request(method, path, body string) *http.Request {
	reqURL, err := url.Parse("http://localhost" + path)
	if err != nil {
		panic(err)
	}
	req := http.Request{
		Method: method,
		URL:    reqURL,
		Header: http.Header{
			"Accept": []string{"application/json; version=2024-08-05"},
		},
	}
	if body != "" {
		req.Body = io.NopCloser(strings.NewReader(body))
	}
	return &req
}
