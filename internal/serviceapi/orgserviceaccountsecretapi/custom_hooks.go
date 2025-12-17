package orgserviceaccountsecretapi

import (
	"encoding/json"
	"errors"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// Resource hooks
var _ autogen.PostReadAPICallHook = (*rs)(nil)

type response struct {
	Secrets []map[string]any
}

func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	var responseJSON response
	if err := json.Unmarshal(result.Body, &responseJSON); err != nil {
		return autogen.APICallResult{Body: nil, Err: err}
	}

	id := req.State.(*TFModel).Id.ValueString()
	for _, secret := range responseJSON.Secrets {
		if secret["id"] == id {
			marshaledSecret, err := json.Marshal(secret)
			return autogen.APICallResult{Body: marshaledSecret, Err: err}
		}
	}

	return autogen.APICallResult{Body: nil, Err: errors.New("secret not found in service account response")}
}

// Data source hooks
var _ autogen.PreReadAPICallHook = (*ds)(nil)
var _ autogen.PostReadAPICallHook = (*ds)(nil)

func (d *ds) PreReadAPICall(callParams config.APICallParams) config.APICallParams {
	callParams.Method = "GET"
	callParams.RelativePath = "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}"
	return callParams
}

func (d *ds) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	var responseJSON response
	if err := json.Unmarshal(result.Body, &responseJSON); err != nil {
		return autogen.APICallResult{Body: nil, Err: err}
	}

	id := req.State.(*TFDSModel).Id.ValueString()
	for _, secret := range responseJSON.Secrets {
		if secret["id"] == id {
			marshaledSecret, err := json.Marshal(secret)
			return autogen.APICallResult{Body: marshaledSecret, Err: err}
		}
	}

	return autogen.APICallResult{Body: nil, Err: errors.New("secret not found in service account response")}
}
