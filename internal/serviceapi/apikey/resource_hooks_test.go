package apikey_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/apikey"
)

const (
	orgID      = "aaaaaaaaaaaaaaaaaaaaaaaa"
	otherOrgID = "bbbbbbbbbbbbbbbbbbbbbbbb"
	apiKeyID   = "cccccccccccccccccccccccc"
	respBody   = `{
		"id": "` + apiKeyID + `",
		"desc": "d",
		"publicKey": "abcdefgh",
		"privateKey": "********-****-****-1234",
		"roles": [
			{"orgId": "` + orgID + `", "roleName": "ORG_MEMBER"},
			{"orgId": "` + orgID + `", "roleName": "ORG_BILLING_ADMIN"},
			{"orgId": "` + otherOrgID + `", "roleName": "ORG_OWNER"},
			{"groupId": "dddddddddddddddddddddddd", "roleName": "GROUP_READ_ONLY"}
		]
	}`
)

func encodedID() string {
	return conversion.EncodeStateID(map[string]string{"org_id": orgID, "api_key_id": apiKeyID})
}

func decodeBody(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var obj map[string]any
	require.NoError(t, json.Unmarshal(body, &obj))
	return obj
}

func TestPostReadAPICall_flattensOrgRolesAndCraftsID(t *testing.T) {
	hook, ok := apikey.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)
	state := &apikey.TFModel{OrgId: types.StringValue(orgID)}
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Body: []byte(respBody)})
	require.NoError(t, result.Err)
	obj := decodeBody(t, result.Body)
	assert.Equal(t, []any{"ORG_MEMBER", "ORG_BILLING_ADMIN"}, obj["roles"])
	assert.Equal(t, encodedID(), obj["tfStateId"])
}

func TestPostReadAPICall_badRequestIsNotFound(t *testing.T) {
	hook, ok := apikey.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)
	result := hook.PostReadAPICall(autogen.HandleReadReq{}, autogen.APICallResult{
		Resp: &http.Response{StatusCode: http.StatusBadRequest},
		Err:  errors.New("API_KEY_NOT_FOUND"),
	})
	assert.ErrorIs(t, result.Err, autogen.ErrNotFound)
}

func TestPostCreateAPICall_takesOrgFromRoles(t *testing.T) {
	hook, ok := apikey.Resource().(autogen.PostCreateAPICallHook)
	require.True(t, ok)
	result := hook.PostCreateAPICall(autogen.APICallResult{Body: []byte(respBody)})
	require.NoError(t, result.Err)
	obj := decodeBody(t, result.Body)
	assert.Equal(t, []any{"ORG_MEMBER", "ORG_BILLING_ADMIN"}, obj["roles"])
	assert.Equal(t, encodedID(), obj["tfStateId"])
}

func TestPreImport(t *testing.T) {
	hook, ok := apikey.Resource().(autogen.PreImportHook)
	require.True(t, ok)
	testCases := map[string]string{
		orgID + "-" + apiKeyID: orgID + "/" + apiKeyID, // legacy manual format
		orgID + "/" + apiKeyID: orgID + "/" + apiKeyID,
		"invalid":              "invalid",
	}
	for input, expected := range testCases {
		got, err := hook.PreImport(input)
		require.NoError(t, err)
		assert.Equal(t, expected, got, input)
	}
}

func TestPostReadAPICall_unmarshalKeepsPrivateKey(t *testing.T) {
	hook, ok := apikey.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)
	state := &apikey.TFModel{
		OrgId:      types.StringValue(orgID),
		ApiKeyId:   types.StringValue(apiKeyID),
		PrivateKey: types.StringValue("unredacted-from-create"),
	}
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Body: []byte(respBody)})
	require.NoError(t, result.Err)
	require.NoError(t, autogen.Unmarshal(result.Body, state))
	assert.Equal(t, "unredacted-from-create", state.PrivateKey.ValueString())
	assert.Equal(t, encodedID(), state.ID.ValueString())
	assert.Equal(t, "d", state.Description.ValueString())
	assert.Equal(t, "abcdefgh", state.PublicKey.ValueString())
	var roles []string
	require.False(t, state.RoleNames.ElementsAs(t.Context(), &roles, false).HasError())
	assert.ElementsMatch(t, []string{"ORG_MEMBER", "ORG_BILLING_ADMIN"}, roles)
}
