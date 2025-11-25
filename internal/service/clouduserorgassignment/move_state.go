package clouduserorgassignment

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

// MoveState is used with moved block to migrate from mongodbatlas_org_invitation to mongodbatlas_cloud_user_org_assignment
func (r *rs) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{{StateMover: stateMover}}
}

func stateMover(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if req.SourceTypeName != "mongodbatlas_org_invitation" || !strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas") {
		return
	}

	setStateResponse(ctx, &resp.Diagnostics, req.SourceRawState, &resp.TargetState)
}

var stateAttrs = map[string]tftypes.Type{
	"org_id":   tftypes.String,
	"username": tftypes.String,
	"roles":    tftypes.List{ElementType: tftypes.String},
}

func setStateResponse(ctx context.Context, diags *diag.Diagnostics, stateIn *tfprotov6.RawState, stateOut *tfsdk.State) {
	rawStateValue, err := stateIn.UnmarshalWithOpts(tftypes.Object{
		AttributeTypes: stateAttrs,
	}, tfprotov6.UnmarshalOpts{ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}})
	if err != nil {
		diags.AddError("Unable to Unmarshal state", err.Error())
		return
	}
	var stateObj map[string]tftypes.Value
	if err := rawStateValue.As(&stateObj); err != nil {
		diags.AddError("Unable to Parse state", err.Error())
		return
	}
	orgID, username := getOrgIDUsernameRolesFromStateObj(diags, stateObj)
	if diags.HasError() {
		return
	}

	model := TFModel{
		OrgId:    types.StringPointerValue(orgID),
		Username: types.StringPointerValue(username),
		Roles:    types.ObjectNull(RolesObjectAttrTypes),               // Let roles be populated during Read
		TeamIds:  types.SetValueMust(types.StringType, []attr.Value{}), // Empty set for team IDs, will be populated during Read
	}

	diags.Append(stateOut.Set(ctx, model)...)
}

func getOrgIDUsernameRolesFromStateObj(diags *diag.Diagnostics, stateObj map[string]tftypes.Value) (orgID, username *string) {
	orgID = schemafunc.GetAttrFromStateObj[string](stateObj, "org_id")
	username = schemafunc.GetAttrFromStateObj[string](stateObj, "username")
	if !conversion.IsStringPresent(orgID) || !conversion.IsStringPresent(username) {
		diags.AddError("Unable to read org_id or username from state", fmt.Sprintf("org_id: %s, username: %s",
			conversion.SafeString(orgID), conversion.SafeString(username)))
		return
	}

	return orgID, username
}
