package privatelinkendpointservicedatafederationonlinearchiveapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const endpointType = "DATA_LAKE"

var (
	_ autogen.PreCreateAPICallHook = (*rs)(nil)
	_ autogen.PreUpdateAPICallHook = (*rs)(nil)
	_ autogen.PreImportHook        = (*rs)(nil)
	_ autogen.ResourceSchemaHook   = (*rs)(nil)
)

type TFExpandedModel struct {
	ID types.String `tfsdk:"id" apiname:"id" autogen:"omitjson"`
}

func (r *rs) PreCreateAPICall(callParams config.APICallParams, bodyReq []byte) (modifiedParams config.APICallParams, modifiedBody []byte) {
	modifiedBody, ok := prepareBody(bodyReq)
	if !ok {
		return callParams, bodyReq
	}
	return callParams, modifiedBody
}

func (r *rs) PreUpdateAPICall(callParams config.APICallParams, bodyReq []byte) (modifiedParams config.APICallParams, modifiedBody []byte) {
	modifiedBody, ok := prepareBody(bodyReq)
	if !ok {
		return callParams, bodyReq
	}
	return callParams, modifiedBody
}

func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	requiresReplace := []string{
		"project_id",
		"endpoint_id",
		"provider_name",
		"region",
		"customer_endpoint_dns_name",
	}
	for _, name := range requiresReplace {
		attr, ok := baseSchema.Attributes[name].(schema.StringAttribute)
		if !ok {
			continue
		}
		// Override generated modifiers (CreateOnly) to mirror manual ForceNew behavior.
		attr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		}
		// Preserve stable planning for Optional+Computed replacement fields.
		if name == "region" || name == "customer_endpoint_dns_name" {
			attr.PlanModifiers = append(attr.PlanModifiers, stringplanmodifier.UseStateForUnknown())
		}
		baseSchema.Attributes[name] = attr
	}

	if regionAttr, ok := baseSchema.Attributes["region"].(schema.StringAttribute); ok {
		regionAttr.Validators = append(regionAttr.Validators,
			validate.ValidUppercaseString(),
		)
		baseSchema.Attributes["region"] = regionAttr
	}

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

	normalizeOptionalStringFields(obj)
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

func (r *rs) PreImport(id string) (string, error) {
	if strings.Contains(id, "/") {
		return id, nil
	}

	parts := strings.Split(id, "--")
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return fmt.Sprintf("%s/%s", parts[0], parts[1]), nil
	}

	return "", fmt.Errorf("use one of the formats: {project_id}/{endpoint_id} or {project_id}--{endpoint_id}")
}

func prepareBody(bodyReq []byte) ([]byte, bool) {
	var body map[string]any
	if err := json.Unmarshal(bodyReq, &body); err != nil {
		return bodyReq, false
	}

	body["type"] = endpointType
	if providerRaw, ok := body["provider"].(string); ok && providerRaw != "" {
		providerUpper := strings.ToUpper(providerRaw)
		body["provider"] = providerUpper
	}

	modifiedBody, err := json.Marshal(body)
	if err != nil {
		return bodyReq, false
	}
	return modifiedBody, true
}

func normalizeOptionalStringFields(obj map[string]any) {
	setEmptyStringIfMissing(obj, "comment")
	setEmptyStringIfMissing(obj, "region")
	setEmptyStringIfMissing(obj, "customerEndpointDNSName")
}

func setEmptyStringIfMissing(obj map[string]any, responseKey string) {
	if val, exists := obj[responseKey]; !exists || val == nil {
		obj[responseKey] = ""
	}
}
