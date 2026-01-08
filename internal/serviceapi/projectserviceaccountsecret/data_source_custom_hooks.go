package projectserviceaccountsecret

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ autogen.PreReadAPICallHook = (*ds)(nil)
var _ autogen.PostReadAPICallHook = (*ds)(nil)

// config.yml configures POST operation for data source read operation to have an accurate response schema
// this hook adjusts the call params to use the correct GET operation
func (d *ds) PreReadAPICall(callParams config.APICallParams) config.APICallParams {
	callParams.Method = "GET"
	callParams.RelativePath = "/api/atlas/v2/groups/{projectId}/serviceAccounts/{clientId}"
	return callParams
}

// PostReadAPICall Reads the secret from the API response Service Account secrets list and returns a new APICallResult with it.
func (d *ds) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	id := req.State.(*TFDSModel).SecretId.ValueString()
	return resourcePostReadAPICall(id, result)
}
