package unit_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

func TestExtractVersion(t *testing.T) {
	version, err := unit.ExtractVersion("application/json; version=2022-06-01")
	require.NoError(t, err)
	require.Equal(t, "2022-06-01", version)
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

func request(method, path string) *http.Request {
	return &http.Request{
		Method: method,
		URL:    &url.URL{Path: path},
		Header: http.Header{
			"Accept": []string{"application/json; version=2024-08-05"},
		},
	}
}

const reqPoliciesCreateBody = `{
 "name": "test-policy",
 "policies": [
  {
   "body": "\t\t\t\n\tforbid (\n\tprincipal,\n\taction == cloud::Action::\"cluster.createEdit\",\n\tresource\n\t) when {\n\tcontext.cluster.cloudProviders.containsAny([cloud::cloudProvider::\"aws\"])\n\t};\n"
  }
 ]
}`
const reqPoliciesUpdateBody = `{
 "name": "updated-policy",
 "policies": [
  {
   "body": "\t\t\t\n\tforbid (\n\tprincipal,\n\taction == cloud::Action::\"cluster.createEdit\",\n\tresource\n\t) when {\n\tcontext.cluster.cloudProviders.containsAny([cloud::cloudProvider::\"aws\"])\n\t};\n"
  }
 ]
}`

const reqPoliciesManualValidateBody = `{
 "name": "overriding the filename by using _manual"
 }`

func TestMockRoundTripper(t *testing.T) {
	orgID := "123"
	resourcePolicyID := "456"
	vars := map[string]string{
		"orgId":            orgID,
		"resourcePolicyId": resourcePolicyID,
	}
	mockTransport, checkFunc := unit.MockRoundTripper(t, vars, &unit.MockHTTPDataConfig{AllowMissingRequests: true})
	client := &http.Client{
		Transport: mockTransport,
	}
	// Error check
	unknownRequest := request("GET", "/v1/cluster/123")
	resp, err := client.Do(unknownRequest)
	require.ErrorContains(t, err, "no matching request found")
	assert.Nil(t, resp)

	// Step 1
	createRequest := request("POST", fmt.Sprintf("/api/atlas/v2/orgs/%s/resourcePolicies", orgID))
	createRequest.Body = io.NopCloser(strings.NewReader(reqPoliciesCreateBody))
	resp, err = client.Do(createRequest)

	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)
	err = checkFunc(nil)
	require.NoError(t, err)
	// Step 2
	patchRequest := request("PATCH", fmt.Sprintf("/api/atlas/v2/orgs/%s/resourcePolicies/%s", orgID, resourcePolicyID))
	patchRequest.Body = io.NopCloser(strings.NewReader(reqPoliciesUpdateBody))
	resp, err = client.Do(patchRequest)
	require.NoError(t, err)
	err = checkFunc(nil)
	require.NoError(t, err)
	var policyResp admin.ApiAtlasResourcePolicy
	err = json.NewDecoder(resp.Body).Decode(&policyResp)
	require.NoError(t, err)
	assert.Equal(t, resourcePolicyID, policyResp.GetId())

	// Step 3
	validateRequest := request("POST", fmt.Sprintf("/api/atlas/v2/orgs/%s/resourcePolicies:validate", orgID))
	validateRequest.Body = io.NopCloser(strings.NewReader(reqPoliciesManualValidateBody))
	_, err = client.Do(validateRequest)
	require.NoError(t, err)
	err = checkFunc(nil)
	require.NoError(t, err)
}
