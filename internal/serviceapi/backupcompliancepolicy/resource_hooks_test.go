package backupcompliancepolicy_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/backupcompliancepolicy"
)

const (
	hookProjectID = "0123456789abcdef01234567"
	hookRespBody  = `{
		"authorizedEmail": "test@example.com",
		"copyProtectionEnabled": false,
		"onDemandPolicyItem": {"frequencyType": "ondemand", "frequencyInterval": 0, "retentionUnit": "days", "retentionValue": 3},
		"scheduledPolicyItems": [
			{"frequencyType": "hourly", "frequencyInterval": 6, "retentionUnit": "days", "retentionValue": 7}
		],
		"state": "ACTIVE"
	}`
)

func respWithPath(p string) *http.Response {
	return &http.Response{Request: &http.Request{URL: &url.URL{Path: p}}}
}

func decodeBody(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var obj map[string]any
	require.NoError(t, json.Unmarshal(body, &obj))
	return obj
}

func TestPreCreateAPICall_preparesPutBody(t *testing.T) {
	hook, ok := backupcompliancepolicy.Resource().(autogen.PreCreateAPICallHook)
	require.True(t, ok)
	params, body := hook.PreCreateAPICall(config.APICallParams{
		PathParams: map[string]string{"projectId": hookProjectID},
	}, []byte(`{"authorizedEmail":"test@example.com"}`))
	assert.Equal(t, "false", params.QueryParams["overwriteBackupPolicies"])
	obj := decodeBody(t, body)
	assert.Equal(t, hookProjectID, obj["projectId"])
	assert.Equal(t, "test@example.com", obj["authorizedEmail"])
}

func TestPreUpdateAPICall_preparesPutBody(t *testing.T) {
	hook, ok := backupcompliancepolicy.Resource().(autogen.PreUpdateAPICallHook)
	require.True(t, ok)
	params, body := hook.PreUpdateAPICall(config.APICallParams{
		PathParams: map[string]string{"projectId": hookProjectID},
	}, []byte(`{"authorizedEmail":"test@example.com"}`))
	assert.Equal(t, "false", params.QueryParams["overwriteBackupPolicies"])
	obj := decodeBody(t, body)
	assert.Equal(t, hookProjectID, obj["projectId"])
}

func TestPreDeleteAPICall_pinsDeleteVersion(t *testing.T) {
	hook, ok := backupcompliancepolicy.Resource().(autogen.PreDeleteAPICallHook)
	require.True(t, ok)
	params := hook.PreDeleteAPICall(config.APICallParams{VersionHeader: "application/vnd.atlas.2023-10-01+json"})
	assert.Equal(t, "application/vnd.atlas.2024-11-13+json", params.VersionHeader)
}

func TestPostCreateAPICall_craftsIDFromRequestPath(t *testing.T) {
	hook, ok := backupcompliancepolicy.Resource().(autogen.PostCreateAPICallHook)
	require.True(t, ok)
	result := hook.PostCreateAPICall(autogen.APICallResult{
		Body: []byte(hookRespBody),
		Resp: respWithPath("/api/atlas/v2/groups/" + hookProjectID + "/backupCompliancePolicy"),
	})
	require.NoError(t, result.Err)
	obj := decodeBody(t, result.Body)
	assert.Equal(t, conversion.EncodeStateID(map[string]string{"project_id": hookProjectID}), obj["tfStateId"])
}

func TestPostReadAPICall_craftsIDFromState(t *testing.T) {
	hook, ok := backupcompliancepolicy.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)
	state := &backupcompliancepolicy.TFModel{ProjectId: types.StringValue(hookProjectID)}
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Body: []byte(hookRespBody)})
	require.NoError(t, result.Err)
	obj := decodeBody(t, result.Body)
	assert.Equal(t, conversion.EncodeStateID(map[string]string{"project_id": hookProjectID}), obj["tfStateId"])
}

func TestPostReadAPICall_emptyBodyKeptAsIs(t *testing.T) {
	hook, ok := backupcompliancepolicy.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: &backupcompliancepolicy.TFModel{}}, autogen.APICallResult{Body: []byte("")})
	require.NoError(t, result.Err)
	assert.Empty(t, result.Body)
}
