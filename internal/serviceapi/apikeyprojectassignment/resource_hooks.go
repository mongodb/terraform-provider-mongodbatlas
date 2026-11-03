package apikeyprojectassignment

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// Manual-parity custom code for mongodbatlas_api_key_project_assignment (lift-and-shift from the
// handwritten Terraform Plugin Framework resource).

var (
	_ autogen.PreCreateAPICallHook  = (*rs)(nil)
	_ autogen.PreReadAPICallHook    = (*rs)(nil)
	_ autogen.PostReadAPICallHook   = (*rs)(nil)
	_ autogen.PostUpdateAPICallHook = (*rs)(nil)
)

// PreCreateAPICall wraps the generated {"roles": [...]} body into the API's
// [{roles, userId}] array. The generated body cannot contain userId because it is a read-only
// property of the array item, so it is taken from the apiKeyId path parameter.
func (r *rs) PreCreateAPICall(callParams config.APICallParams, bodyReq []byte) (params config.APICallParams, body []byte) {
	var bodyMap map[string]any
	if err := autogen.Decode(bodyReq, &bodyMap); err != nil {
		return callParams, bodyReq
	}
	wrapped, err := json.Marshal([]map[string]any{
		{
			"roles":  bodyMap["roles"],
			"userId": callParams.PathParams["apiKeyId"],
		},
	})
	if err != nil {
		return callParams, bodyReq
	}
	return callParams, wrapped
}

// PreReadAPICall requests the maximum page size. The read operation lists the project API keys
// (there is no GET for a single assignment) and the manual resource paginated through all of them;
// 500 is the API maximum and the maximum number of organization API keys, so one page covers them.
func (r *rs) PreReadAPICall(callParams config.APICallParams) config.APICallParams {
	if callParams.QueryParams == nil {
		callParams.QueryParams = map[string]string{}
	}
	callParams.QueryParams["itemsPerPage"] = "500"
	return callParams
}

// PostReadAPICall selects the assignment from the project API key list and reduces its role
// objects to the project-level role names, as the manual resource did.
func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}
	model, ok := req.State.(*TFModel)
	if !ok {
		return result
	}
	projectID := model.ProjectId.ValueString()
	apiKeyID := model.ApiKeyId.ValueString()
	body, err := adaptResponse(result.Body, projectID, apiKeyID)
	if err != nil {
		// The manual read returned an error diagnostic when the key was not assigned to the project;
		// autogen removes the resource from state when the hook reports ErrNotFound.
		return autogen.APICallResult{Resp: result.Resp, Err: err}
	}
	return autogen.APICallResult{Body: body, Resp: result.Resp}
}

// PostUpdateAPICall discards the update response. It returns the API key details with role objects
// for every project the key is assigned to, which cannot be unmarshaled into the roles set; the plan
// already holds the roles that were accepted, and the next read filters them per project.
func (r *rs) PostUpdateAPICall(result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}
	return autogen.APICallResult{Resp: result.Resp}
}

type readAPIResponse struct {
	Results []map[string]any
}

func adaptResponse(body []byte, projectID, apiKeyID string) ([]byte, error) {
	var response readAPIResponse
	if err := autogen.Decode(body, &response); err != nil {
		// Leave the body untouched so the generated unmarshal reports the API response error.
		return body, nil
	}
	for _, apiKey := range response.Results {
		if id, _ := apiKey["id"].(string); id != apiKeyID {
			continue
		}
		return json.Marshal(map[string]any{
			"roles": projectRoleNames(apiKey["roles"], projectID),
		})
	}
	return nil, fmt.Errorf("API key %s is not assigned to project %s: %w", apiKeyID, projectID, autogen.ErrNotFound)
}

func projectRoleNames(roles any, projectID string) []string {
	roleList, _ := roles.([]any)
	roleNames := make([]string, 0, len(roleList))
	for _, role := range roleList {
		roleObj, _ := role.(map[string]any)
		roleName, _ := roleObj["roleName"].(string)
		groupID, _ := roleObj["groupId"].(string)
		if strings.HasPrefix(roleName, "ORG_") || groupID != projectID {
			continue
		}
		roleNames = append(roleNames, roleName)
	}
	return roleNames
}
