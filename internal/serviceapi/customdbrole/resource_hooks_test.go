package customdbrole_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/customdbrole"
)

const (
	projectID = "0123456789abcdef01234567"
	roleName  = "my-role-name"
	respBody  = `{
		"roleName": "` + roleName + `",
		"actions": [
			{"action": "INSERT", "resources": [{"cluster": false, "collection": "", "db": "db1"}]},
			{"action": "SERVER_STATUS", "resources": [{"cluster": true}]}
		],
		"inheritedRoles": []
	}`
)

func encodedID() string {
	return conversion.EncodeStateID(map[string]string{"project_id": projectID, "role_name": roleName})
}

func decodeBody(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var obj map[string]any
	require.NoError(t, json.Unmarshal(body, &obj))
	return obj
}

func respWithPath(p string) *http.Response {
	return &http.Response{Request: &http.Request{URL: &url.URL{Path: p}}}
}

func assertNormalizedActions(t *testing.T, obj map[string]any) {
	t.Helper()
	actions := obj["actions"].([]any)
	clusterRes := actions[1].(map[string]any)["resources"].([]any)[0]
	assert.Equal(t, map[string]any{"cluster": true, "collection": "", "db": ""}, clusterRes)
	dbRes := actions[0].(map[string]any)["resources"].([]any)[0]
	assert.Equal(t, map[string]any{"cluster": false, "collection": "", "db": "db1"}, dbRes)
}

func TestPostReadAPICall_normalizesAndCraftsID(t *testing.T) {
	hook, ok := customdbrole.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)
	state := &customdbrole.TFModel{ProjectId: types.StringValue(projectID)}
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Body: []byte(respBody)})
	require.NoError(t, result.Err)
	obj := decodeBody(t, result.Body)
	assert.Equal(t, encodedID(), obj["tfStateId"])
	assertNormalizedActions(t, obj)
}

func TestPostCreateAndUpdateAPICall_projectIDFromRequestPath(t *testing.T) {
	r := customdbrole.Resource()
	createHook, ok := r.(autogen.PostCreateAPICallHook)
	require.True(t, ok)
	updateHook, ok := r.(autogen.PostUpdateAPICallHook)
	require.True(t, ok)

	created := createHook.PostCreateAPICall(autogen.APICallResult{
		Body: []byte(respBody),
		Resp: respWithPath("/api/atlas/v2/groups/" + projectID + "/customDBRoles/roles"),
	})
	require.NoError(t, created.Err)
	assert.Equal(t, encodedID(), decodeBody(t, created.Body)["tfStateId"])
	assertNormalizedActions(t, decodeBody(t, created.Body))

	updated := updateHook.PostUpdateAPICall(autogen.APICallResult{
		Body: []byte(respBody),
		Resp: respWithPath("/api/atlas/v2/groups/" + projectID + "/customDBRoles/roles/" + roleName),
	})
	require.NoError(t, updated.Err)
	assert.Equal(t, encodedID(), decodeBody(t, updated.Body)["tfStateId"])
}

func TestPostCreateAPICall_errorAndEmptyBodyUnchanged(t *testing.T) {
	hook, ok := customdbrole.Resource().(autogen.PostCreateAPICallHook)
	require.True(t, ok)
	errResult := autogen.APICallResult{Err: errors.New("boom"), Body: []byte(respBody)}
	assert.Equal(t, errResult, hook.PostCreateAPICall(errResult))
	emptyResult := autogen.APICallResult{Body: []byte("")}
	assert.Equal(t, emptyResult, hook.PostCreateAPICall(emptyResult))
}

func TestPreImport(t *testing.T) {
	hook, ok := customdbrole.Resource().(autogen.PreImportHook)
	require.True(t, ok)
	testCases := map[string]struct {
		id       string
		expected string
	}{
		"legacy dash format":               {id: projectID + "-" + roleName, expected: projectID + "/" + roleName},
		"slash format unchanged":           {id: projectID + "/" + roleName, expected: projectID + "/" + roleName},
		"non-hex project id unchanged":     {id: "not-a-project-id", expected: "not-a-project-id"},
		"role name with slash in new form": {id: projectID + "/a-b", expected: projectID + "/a-b"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := hook.PreImport(tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestResourceSchema_roleNameValidators(t *testing.T) {
	resp := &resource.SchemaResponse{}
	customdbrole.Resource().Schema(t.Context(), resource.SchemaRequest{}, resp)
	require.False(t, resp.Diagnostics.HasError())
	attr := resp.Schema.Attributes["role_name"].(schema.StringAttribute)
	validate := func(value string) bool {
		for _, v := range attr.Validators {
			vResp := &validator.StringResponse{}
			v.ValidateString(t.Context(), validator.StringRequest{Path: path.Root("role_name"), ConfigValue: types.StringValue(value)}, vResp)
			if vResp.Diagnostics.HasError() {
				return false
			}
		}
		return true
	}
	assert.True(t, validate("my-role_1"))
	assert.False(t, validate("x-gen-role"), "prefix check ported as-is rejects x-gen")
	assert.True(t, validate("xgen-role"), "ported as-is: the message says xgen- but only x-gen is rejected")
	assert.True(t, validate("role.with.dots"), "ported as-is: the regex is unanchored")
	assert.False(t, validate("..."))
}
