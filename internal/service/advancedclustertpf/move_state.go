package advancedclustertpf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func (r *rs) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{{StateMover: stateMover}}
}

func stateMover(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	diags := &resp.Diagnostics
	if req.SourceTypeName != "mongodbatlas_cluster" || !strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas") {
		return
	}
	setStateResponse(ctx, diags, req.SourceRawState, &resp.TargetState)
}

// setStateResponse is used in Move State and Upgrade State
func setStateResponse(ctx context.Context, diags *diag.Diagnostics, stateIn *tfprotov6.RawState, stateOut *tfsdk.State) {
	rawStateValue, err := stateIn.UnmarshalWithOpts(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"project_id": tftypes.String,
			"name":       tftypes.String,
		},
	}, tfprotov6.UnmarshalOpts{ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}})
	if err != nil {
		diags.AddError("Unable to Unmarshal state", err.Error())
		return
	}
	var rawState map[string]tftypes.Value
	if err := rawStateValue.As(&rawState); err != nil {
		diags.AddError("Unable to Parse state", err.Error())
		return
	}
	var projectID *string
	if err := rawState["project_id"].As(&projectID); err != nil {
		diags.AddAttributeError(path.Root("project_id"), "Unable to read cluster project_id", err.Error())
		return
	}
	var name *string
	if err := rawState["name"].As(&name); err != nil {
		diags.AddAttributeError(path.Root("name"), "Unable to read cluster name", err.Error())
		return
	}
	if !conversion.IsStringPresent(projectID) || !conversion.IsStringPresent(name) {
		diags.AddError("Unable to read project_id or name", fmt.Sprintf("project_id: %s, name: %s",
			conversion.SafeString(projectID), conversion.SafeString(name)))
		return
	}
	validTimeout := timeouts.Value{
		Object: types.ObjectNull(
			map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"delete": types.StringType,
			}),
	}
	model := NewTFModel(ctx, &admin.ClusterDescription20240805{
		GroupId: projectID,
		Name:    name,
	}, validTimeout, diags, ExtraAPIInfo{})
	if diags.HasError() {
		return
	}
	AddAdvancedConfig(ctx, model, nil, nil, diags)
	if diags.HasError() {
		return
	}
	diags.Append(stateOut.Set(ctx, model)...)
}
