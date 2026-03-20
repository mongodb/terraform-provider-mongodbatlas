package privatelinkendpointservicedatafederationonlinearchiveapi

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const endpointType = "DATA_LAKE"

var (
	_ autogen.PreCreateAPICallHook = (*rs)(nil)
	_ autogen.PreUpdateAPICallHook = (*rs)(nil)
	_ autogen.ResourceSchemaHook   = (*rs)(nil)
)

type TFExpandedModel struct {
	ID types.String `tfsdk:"id" apiname:"id" autogen:"omitjson"`
}

func (r *rs) PreCreateAPICall(callParams config.APICallParams, bodyReq []byte) (modifiedParams config.APICallParams, modifiedBody []byte) {
	modifiedBody, ok := injectEndpointType(bodyReq)
	if !ok {
		return callParams, bodyReq
	}
	return callParams, modifiedBody
}

func (r *rs) PreUpdateAPICall(callParams config.APICallParams, bodyReq []byte) (modifiedParams config.APICallParams, modifiedBody []byte) {
	modifiedBody, ok := injectEndpointType(bodyReq)
	if !ok {
		return callParams, bodyReq
	}
	return callParams, modifiedBody
}

func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	baseSchema.Attributes["id"] = schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	return baseSchema
}

func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}

	model, ok := req.State.(*TFModel)
	if !ok || model.ProjectId.IsNull() || model.EndpointId.IsNull() {
		return result
	}

	craftedID := conversion.EncodeStateID(map[string]string{
		"project_id":  model.ProjectId.ValueString(),
		"endpoint_id": model.EndpointId.ValueString(),
	})

	var obj map[string]any
	if err := json.Unmarshal(result.Body, &obj); err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}

	obj["id"] = craftedID

	body, err := json.Marshal(obj)
	if err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}

	return autogen.APICallResult{
		Body: body,
		Err:  nil,
		Resp: result.Resp,
	}
}

// Ensures POST request body includes type=DATA_LAKE.
func injectEndpointType(bodyReq []byte) ([]byte, bool) {
	var body map[string]any
	if err := json.Unmarshal(bodyReq, &body); err != nil {
		return bodyReq, false
	}
	body["type"] = endpointType
	modifiedBody, err := json.Marshal(body)
	if err != nil {
		return bodyReq, false
	}
	return modifiedBody, true
}
