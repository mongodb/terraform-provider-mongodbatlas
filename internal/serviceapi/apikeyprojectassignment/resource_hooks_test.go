package apikeyprojectassignment_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/apikeyprojectassignment"
)

const (
	projectID      = "aaaaaaaaaaaaaaaaaaaaaaaa"
	otherProjectID = "bbbbbbbbbbbbbbbbbbbbbbbb"
	apiKeyID       = "cccccccccccccccccccccccc"
	otherAPIKeyID  = "dddddddddddddddddddddddd"
	listResponse   = `{
		"results": [
			{
				"id": "` + otherAPIKeyID + `",
				"roles": [{"groupId": "` + projectID + `", "roleName": "GROUP_OWNER"}]
			},
			{
				"id": "` + apiKeyID + `",
				"roles": [
					{"orgId": "eeeeeeeeeeeeeeeeeeeeeeee", "roleName": "ORG_MEMBER"},
					{"groupId": "` + projectID + `", "roleName": "GROUP_OWNER"},
					{"groupId": "` + projectID + `", "roleName": "GROUP_READ_ONLY"},
					{"groupId": "` + otherProjectID + `", "roleName": "GROUP_OWNER"}
				]
			}
		],
		"totalCount": 2
	}`
)

func newState(t *testing.T) *apikeyprojectassignment.TFModel {
	t.Helper()
	return &apikeyprojectassignment.TFModel{
		ProjectId: types.StringValue(projectID),
		ApiKeyId:  types.StringValue(apiKeyID),
		Roles: customtypes.NewSetValue[types.String](t.Context(), []attr.Value{
			types.StringValue("GROUP_OWNER"),
		}),
	}
}

func TestMarshalCreateBody_onlyRoles(t *testing.T) {
	body, err := autogen.Marshal(newState(t), false)
	require.NoError(t, err)
	assert.JSONEq(t, `{"roles":["GROUP_OWNER"]}`, string(body))
}

func TestPreCreateAPICall_wrapsRolesWithUserID(t *testing.T) {
	hook, ok := apikeyprojectassignment.Resource().(autogen.PreCreateAPICallHook)
	require.True(t, ok)
	callParams := config.APICallParams{PathParams: map[string]string{"apiKeyId": apiKeyID}}
	gotParams, gotBody := hook.PreCreateAPICall(callParams, []byte(`{"roles":["GROUP_OWNER","GROUP_READ_ONLY"]}`))
	assert.Equal(t, callParams, gotParams)
	assert.JSONEq(t, `[{"roles":["GROUP_OWNER","GROUP_READ_ONLY"],"userId":"`+apiKeyID+`"}]`, string(gotBody))
}

func TestPreReadAPICall_setsMaximumPageSize(t *testing.T) {
	hook, ok := apikeyprojectassignment.Resource().(autogen.PreReadAPICallHook)
	require.True(t, ok)
	got := hook.PreReadAPICall(config.APICallParams{})
	assert.Equal(t, map[string]string{"itemsPerPage": "500"}, got.QueryParams)
}

func TestPostReadAPICall_selectsAssignmentAndKeepsProjectRoles(t *testing.T) {
	hook, ok := apikeyprojectassignment.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)
	state := newState(t)
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Body: []byte(listResponse)})
	require.NoError(t, result.Err)
	require.NoError(t, autogen.Unmarshal(result.Body, state))
	var roles []string
	require.False(t, state.Roles.ElementsAs(t.Context(), &roles, false).HasError())
	assert.ElementsMatch(t, []string{"GROUP_OWNER", "GROUP_READ_ONLY"}, roles)
	assert.Equal(t, projectID, state.ProjectId.ValueString())
	assert.Equal(t, apiKeyID, state.ApiKeyId.ValueString())
}

func TestPostReadAPICall_notAssignedIsNotFound(t *testing.T) {
	hook, ok := apikeyprojectassignment.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)
	state := newState(t)
	state.ApiKeyId = types.StringValue("ffffffffffffffffffffffff")
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Body: []byte(listResponse)})
	assert.ErrorIs(t, result.Err, autogen.ErrNotFound)
}

func TestPostUpdateAPICall_discardsRoleObjects(t *testing.T) {
	hook, ok := apikeyprojectassignment.Resource().(autogen.PostUpdateAPICallHook)
	require.True(t, ok)
	result := hook.PostUpdateAPICall(autogen.APICallResult{Body: []byte(`{"id":"` + apiKeyID + `","roles":[{"groupId":"` + projectID + `","roleName":"GROUP_OWNER"}]}`)})
	require.NoError(t, result.Err)
	assert.True(t, autogen.IsEmptyJSON(result.Body))
	// The update response must not break unmarshaling into the model; the plan is kept.
	state := newState(t)
	require.NoError(t, autogen.Unmarshal(result.Body, state))
	assert.Equal(t, []attr.Value{types.StringValue("GROUP_OWNER")}, state.Roles.Elements())
}

func TestPostReadAPICall_apiErrorIsReturned(t *testing.T) {
	hook, ok := apikeyprojectassignment.Resource().(autogen.PostReadAPICallHook)
	require.True(t, ok)
	errAPI := assert.AnError
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: newState(t)}, autogen.APICallResult{Err: errAPI})
	assert.ErrorIs(t, result.Err, errAPI)
}
