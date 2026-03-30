package privatelinkendpointservicedatafederationonlinearchiveapi

import (
	"context"
	"encoding/json"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

var _ autogen.PostReadAPICallHook = (*ds)(nil)
var _ autogen.DataSourceSchemaHook = (*ds)(nil)
var _ autogen.PostReadAggregatedListAPICallHook = (*pluralDS)(nil)
var _ autogen.DataSourceSchemaHook = (*pluralDS)(nil)

type TFDSExpandedModel struct {
	ID types.String `tfsdk:"id" apiname:"id" autogen:"omitjson"`
}

type TFPluralDSExpandedModel struct {
	ID types.String `tfsdk:"id" apiname:"id" autogen:"omitjson"`
}

func (d *ds) DataSourceSchema(_ context.Context, baseSchema datasourceschema.Schema) datasourceschema.Schema {
	baseSchema.Attributes["id"] = datasourceschema.StringAttribute{
		Computed: true,
	}
	return baseSchema
}

func (d *pluralDS) DataSourceSchema(_ context.Context, baseSchema datasourceschema.Schema) datasourceschema.Schema {
	baseSchema.Attributes["id"] = datasourceschema.StringAttribute{
		Computed: true,
	}
	return baseSchema
}

// PostReadAPICall injects a crafted ID into the singular data source response body.
func (d *ds) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}

	model, ok := req.State.(*TFDSModel)
	if !ok || model.ProjectId.IsNull() || model.ProjectId.IsUnknown() || model.EndpointId.IsNull() || model.EndpointId.IsUnknown() {
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
	// Mirror SDKv2 behavior for omitted optional strings in state.
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

// PostReadAggregatedListAPICall injects a generated ID after list pagination has completed.
func (d *pluralDS) PostReadAggregatedListAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}
	model, ok := req.State.(*TFPluralDSModel)
	if !ok || model.ProjectId.IsNull() || model.ProjectId.IsUnknown() {
		return result
	}

	var obj map[string]any
	if err := json.Unmarshal(result.Body, &obj); err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}

	// Mirror SDKv2 behavior for omitted optional strings
	if results, ok := obj["results"].([]any); ok {
		for i := range results {
			if entry, ok := results[i].(map[string]any); ok {
				normalizeOptionalStringFields(entry)
				results[i] = entry
			}
		}
		obj["results"] = results
	}
	// Injects a generated ID for the plural data source, keeps same behavior as SDKv2 manual data source
	obj["id"] = id.UniqueId()

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
