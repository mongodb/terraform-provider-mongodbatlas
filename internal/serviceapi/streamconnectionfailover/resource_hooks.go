package streamconnectionfailover

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// ResourceSchema marks `type` and `region` as requiring replacement, since both are immutable on the
// failover connection PATCH. The autogen schema has no create-only knob for body fields.
func (r *rs) ResourceSchema(ctx context.Context, s schema.Schema) schema.Schema {
	for _, name := range []string{"type", "region"} {
		attr, ok := s.Attributes[name].(schema.StringAttribute)
		if !ok {
			continue
		}
		attr.PlanModifiers = append(attr.PlanModifiers, stringplanmodifier.RequiresReplace())
		s.Attributes[name] = attr
	}
	return s
}

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
