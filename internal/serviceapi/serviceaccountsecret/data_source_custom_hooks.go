package serviceaccountsecret

import (
	"encoding/json"
	"errors"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ autogen.PreReadAPICallHook = (*ds)(nil)
var _ autogen.PostReadAPICallHook = (*ds)(nil)

// config.yml configures POST operation for data source read operation to have an accurate response schema
// this hook adjusts the call params to use the correct GET operation
func (d *ds) PreReadAPICall(callParams config.APICallParams) config.APICallParams {
	callParams.Method = "GET"
	callParams.RelativePath = "/api/atlas/v2/orgs/{orgId}/serviceAccounts/{clientId}"
	return callParams
}

func (d *ds) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}
	var responseJSON response
	if err := json.Unmarshal(result.Body, &responseJSON); err != nil {
		return autogen.APICallResult{Body: nil, Err: err}
	}

	id := req.State.(*TFDSModel).SecretId.ValueString()
	for _, secret := range responseJSON.Secrets {
		if secret["id"] == id {
			marshaledSecret, err := json.Marshal(secret)
			return autogen.APICallResult{Body: marshaledSecret, Err: err}
		}
	}

	return autogen.APICallResult{Body: nil, Err: errors.New("secret not found in service account response")}
}
