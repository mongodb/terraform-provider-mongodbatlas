package projectserviceaccountsecret

import (
	"encoding/json"
	"errors"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var _ autogen.PostReadAPICallHook = (*rs)(nil)

type readAPIResponse struct {
	Secrets []map[string]any
}

const errImportFormat = "use one of the formats: {project_id}/{workspace_name}/{connection_name} or {workspace_name}-{project_id}-{connection_name}"

// PostReadAPICall Reads the secret from the API response Service Account secrets list and returns a new APICallResult with it.
func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	id := req.State.(*TFModel).SecretId.ValueString()
	return resourcePostReadAPICall(id, result)
}

func resourcePostReadAPICall(secretID string, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}
	var responseJSON readAPIResponse
	if err := json.Unmarshal(result.Body, &responseJSON); err != nil {
		return autogen.APICallResult{Body: nil, Err: err}
	}

	for _, secret := range responseJSON.Secrets {
		if secret["id"] == secretID {
			marshaledSecret, err := json.Marshal(secret)
			return autogen.APICallResult{Body: marshaledSecret, Err: err}
		}
	}

	return autogen.APICallResult{Body: nil, Err: errors.New("secret not found in service account response")}
}
