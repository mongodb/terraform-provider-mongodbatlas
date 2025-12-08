package orgserviceaccountsecretapi

import (
	"encoding/json"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

type response struct {
	Secrets []map[string]any
}

func (r *rs) PostReadAPICall(result autogen.APICallResult) autogen.APICallResult {
	var responseJSON response
	if err := json.Unmarshal(result.Body, &responseJSON); err != nil {
		return autogen.APICallResult{Body: nil, Err: err}
	}

	// id := req.State.(*TFModel).Id.ValueString()
	for _, secret := range responseJSON.Secrets {
		if secret["id"] == "" {
			marshaledSecret, err := json.Marshal(secret)
			return autogen.APICallResult{Body: marshaledSecret, Err: err}
		}
	}

	return autogen.APICallResult{Body: nil, Err: nil}
}
