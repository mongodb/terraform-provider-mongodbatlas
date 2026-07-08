package streamconnectionfailover

import (
	"encoding/json"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// A failover connection's body `name` must equal its parent connection name (the `connectionName`
// path param). The generated schema omits the body `name` (it would collide with the
// `connection_name` path param), so inject it into the create/update request body. Without it, the
// server fails converting the request body and returns a 500 UNEXPECTED_ERROR.

func (r *rs) PreCreateAPICall(callParams config.APICallParams, bodyReq []byte) (params config.APICallParams, body []byte) {
	return callParams, injectConnectionName(callParams, bodyReq)
}

func (r *rs) PreUpdateAPICall(callParams config.APICallParams, bodyReq []byte) (params config.APICallParams, body []byte) {
	return callParams, injectConnectionName(callParams, bodyReq)
}

func injectConnectionName(callParams config.APICallParams, bodyReq []byte) []byte {
	name := callParams.PathParams["connectionName"]
	if name == "" {
		return bodyReq
	}
	var body map[string]any
	if err := json.Unmarshal(bodyReq, &body); err != nil {
		return bodyReq
	}
	body["name"] = name
	updated, err := json.Marshal(body)
	if err != nil {
		return bodyReq
	}
	return updated
}
