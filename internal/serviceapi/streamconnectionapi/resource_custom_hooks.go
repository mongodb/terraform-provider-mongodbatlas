package streamconnectionapi

import (
	"encoding/json"
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var _ autogen.PostReadAPICallHook = (*rs)(nil)
var _ autogen.PostCreateAPICallHook = (*rs)(nil)
var _ autogen.PostUpdateAPICallHook = (*rs)(nil)

func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}

	m := req.State.(*TFModel)
	return injectCraftedID(result, m.WorkspaceName.ValueString(), m.ProjectId.ValueString(), m.ConnectionName.ValueString())
}

func (r *rs) PostCreateAPICall(result autogen.APICallResult) autogen.APICallResult {
	return injectCraftedIDFromBody(result)
}

func (r *rs) PostUpdateAPICall(result autogen.APICallResult) autogen.APICallResult {
	return injectCraftedIDFromBody(result)
}

func injectCraftedIDFromBody(result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}

	var obj map[string]any
	if err := json.Unmarshal(result.Body, &obj); err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}

	workspace, okWorkspace := getStringFromKeys(obj, "workspaceName", "tenantName")
	projectID, okProjectID := getStringFromKeys(obj, "projectId", "groupId")
	connection, okConnection := getStringFromKeys(obj, "connectionName", "name")
	if !okWorkspace || !okProjectID || !okConnection {
		return result
	}

	return injectCraftedID(result, workspace, projectID, connection)
}

func injectCraftedID(result autogen.APICallResult, workspace, projectID, connection string) autogen.APICallResult {
	if result.Err != nil {
		return result
	}

	var obj map[string]any
	if err := json.Unmarshal(result.Body, &obj); err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}

	obj["id"] = fmt.Sprintf("%s-%s-%s", workspace, projectID, connection)

	b, err := json.Marshal(obj)
	if err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}
	return autogen.APICallResult{Body: b, Err: nil, Resp: result.Resp}
}

func getStringFromKeys(body map[string]any, keys ...string) (string, bool) {
	for _, key := range keys {
		if value, ok := body[key]; ok {
			if strVal, ok := value.(string); ok && strVal != "" {
				return strVal, true
			}
		}
	}
	return "", false
}
