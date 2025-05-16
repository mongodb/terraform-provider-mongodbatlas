package unit_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.mongodb.org/atlas-sdk/v20250312003/admin"
)

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

const reqPoliciesManualValidateDelete = `{}`

func TestMockRoundTripper(t *testing.T) {
	orgID := "123"
	resourcePolicyID := "456"
	data := unit.ReadMockData(t, []string{"", "", ""})
	data.Variables = map[string]string{}
	mockTransport, tracker := unit.NewMockRoundTripper(t, &unit.MockHTTPDataConfig{AllowMissingRequests: true}, data)
	client := &http.Client{
		Transport: mockTransport,
	}
	// Error check
	tracker.IncreaseStepNumberAndInit()
	unknownRequest := request("GET", "/v1/cluster/123", "")
	resp, err := client.Do(unknownRequest)
	require.ErrorContains(t, err, "no matching request found")
	assert.Nil(t, resp)

	// Step 1
	createRequest := request("POST", fmt.Sprintf("/api/atlas/v2/orgs/%s/resourcePolicies", orgID), reqPoliciesCreateBody)
	resp, err = client.Do(createRequest)

	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)
	err = tracker.CheckStepRequests(nil)
	require.NoError(t, err)
	// Step 2
	tracker.IncreaseStepNumberAndInit()
	patchRequest := request("PATCH", fmt.Sprintf("/api/atlas/v2/orgs/%s/resourcePolicies/%s", orgID, resourcePolicyID), reqPoliciesUpdateBody)
	resp, err = client.Do(patchRequest)
	require.NoError(t, err)
	err = tracker.CheckStepRequests(nil)
	require.NoError(t, err)
	var policyResp admin.ApiAtlasResourcePolicy
	err = json.NewDecoder(resp.Body).Decode(&policyResp)
	require.NoError(t, err)
	assert.Equal(t, resourcePolicyID, policyResp.GetId())

	// Step 3
	tracker.IncreaseStepNumberAndInit()
	// First GET request OK
	// Second GET request OK
	getRequest := request("GET", fmt.Sprintf("/api/atlas/v2/orgs/%s/resourcePolicies/%s", orgID, resourcePolicyID), "")
	_, err = client.Do(getRequest)
	require.NoError(t, err)
	_, err = client.Do(getRequest)
	require.NoError(t, err)
	// Third GET request is re-read, since we have not gotten the diff
	require.NoError(t, err)
	okResp, err := client.Do(getRequest)
	require.NoError(t, err)
	require.Equal(t, 200, okResp.StatusCode)

	// Test _manual diff file (set to {} instead of '')
	validateRequest := request("DELETE", fmt.Sprintf("/api/atlas/v2/orgs/%s/resourcePolicies/%s", orgID, resourcePolicyID), reqPoliciesManualValidateDelete)
	_, err = client.Do(validateRequest)
	require.NoError(t, err)
	// Fourth GET request OK, since we have gotten the diff
	notFoundResp, err := client.Do(getRequest)
	require.NoError(t, err)
	notFoundMap := parseMapStringAny(t, notFoundResp)
	assert.Equal(t, "RESOURCE_POLICY_NOT_FOUND", notFoundMap["errorCode"])

	err = tracker.CheckStepRequests(nil)
	require.NoError(t, err)
}

func parseMapStringAny(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	stringMap := map[string]any{}
	err := json.NewDecoder(resp.Body).Decode(&stringMap)
	require.NoError(t, err)
	return stringMap
}

func TestMockRoundTripperAllowReRead(t *testing.T) {
	orgID := "123"
	data := unit.ReadMockData(t, []string{""})
	data.Variables = map[string]string{}
	mockTransport, tracker := unit.NewMockRoundTripper(t, &unit.MockHTTPDataConfig{AllowMissingRequests: true}, data)
	client := &http.Client{
		Transport: mockTransport,
	}
	tracker.IncreaseStepNumberAndInit()
	for range []int{0, 1, 2} {
		getRequest := request("GET", fmt.Sprintf("/api/atlas/v2/orgs/%s/resourcePolicies", orgID), "")
		resp, err := client.Do(getRequest)
		require.NoError(t, err)
		assert.Equal(t, "returned again", parseMapStringAny(t, resp)["expect"])
	}
	createRequest := request("POST", fmt.Sprintf("/api/atlas/v2/orgs/%s/resourcePolicies", orgID), reqPoliciesCreateBody)
	resp, err := client.Do(createRequest)

	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)
	err = tracker.CheckStepRequests(nil)
	require.NoError(t, err)
}
