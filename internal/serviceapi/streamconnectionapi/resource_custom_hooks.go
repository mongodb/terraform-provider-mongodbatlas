package streamconnectionapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

// TFExpandedModel holds Terraform-only attributes (e.g. crafted id) not returned by the API.
type TFExpandedModel struct {
	Id types.String `tfsdk:"id" autogen:"omitjson"`
}

var _ autogen.PostReadAPICallHook = (*rs)(nil)
var _ autogen.SchemaExtensionHook = (*rs)(nil)

// PostReadAPICall injects the crafted id into the API response. It runs on normal Read and during
// Create/Update wait (refresh). The model always has workspace/project/connection from plan or state.
func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err != nil {
		return result
	}
	m, ok := req.State.(*TFModel)
	if !ok || m.WorkspaceName.IsNull() || m.ProjectId.IsNull() || m.ConnectionName.IsNull() {
		return result
	}
	craftedID := fmt.Sprintf("%s-%s-%s", m.WorkspaceName.ValueString(), m.ProjectId.ValueString(), m.ConnectionName.ValueString())
	var obj map[string]any
	if err := json.Unmarshal(result.Body, &obj); err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}
	obj["id"] = craftedID
	b, err := json.Marshal(obj)
	if err != nil {
		return autogen.APICallResult{Body: nil, Err: err, Resp: result.Resp}
	}
	return autogen.APICallResult{Body: b, Err: nil, Resp: result.Resp}
}

func (r *rs) ExtendSchema(ctx context.Context, base schema.Schema) schema.Schema {
	base.Attributes["id"] = schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "ID composed of workspace_name, project_id, and connection_name.",
		PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	return base
}
