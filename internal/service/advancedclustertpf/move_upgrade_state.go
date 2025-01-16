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

func (r *rs) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: {StateUpgrader: stateUpgraderFromV1},
	}
}

func stateMover(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if req.SourceTypeName != "mongodbatlas_cluster" || !strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas") {
		return
	}
	setStateResponse(ctx, &resp.Diagnostics, req.SourceRawState, &resp.TargetState)
}

func stateUpgraderFromV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	setStateResponse(ctx, &resp.Diagnostics, req.RawState, &resp.State)
}

func setStateResponse(ctx context.Context, diags *diag.Diagnostics, stateIn *tfprotov6.RawState, stateOut *tfsdk.State) {
	rawStateValue, err := stateIn.UnmarshalWithOpts(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"project_id":             tftypes.String,
			"name":                   tftypes.String,
			"retain_backups_enabled": tftypes.Bool,
			"mongo_db_major_version": tftypes.String,
			"timeouts": tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"create": tftypes.String,
					"update": tftypes.String,
					"delete": tftypes.String,
				},
			},
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

	projectID := getAttrFromRawState[string](diags, rawState, "project_id")
	name := getAttrFromRawState[string](diags, rawState, "name")
	if diags.HasError() {
		return
	}
	if !conversion.IsStringPresent(projectID) || !conversion.IsStringPresent(name) {
		diags.AddError("Unable to read project_id or name from state", fmt.Sprintf("project_id: %s, name: %s",
			conversion.SafeString(projectID), conversion.SafeString(name)))
		return
	}

	model := NewTFModel(ctx, &admin.ClusterDescription20240805{
		GroupId: projectID,
		Name:    name,
	}, getAttrTimeout(diags, rawState), diags, ExtraAPIInfo{})
	if diags.HasError() {
		return
	}

	if retainBackupsEnabled := getAttrFromRawState[bool](diags, rawState, "retain_backups_enabled"); retainBackupsEnabled != nil {
		model.RetainBackupsEnabled = types.BoolPointerValue(retainBackupsEnabled)
	}
	if mongoDBMajorVersion := getAttrFromRawState[string](diags, rawState, "mongo_db_major_version"); mongoDBMajorVersion != nil {
		model.MongoDBMajorVersion = types.StringPointerValue(mongoDBMajorVersion)
	}
	if diags.HasError() {
		return
	}

	AddAdvancedConfig(ctx, model, nil, nil, diags)
	if diags.HasError() {
		return
	}
	diags.Append(stateOut.Set(ctx, model)...)
}

func getAttrFromRawState[T any](diags *diag.Diagnostics, rawState map[string]tftypes.Value, attrName string) *T {
	var ret *T
	if err := rawState[attrName].As(&ret); err != nil {
		diags.AddAttributeError(path.Root(attrName), fmt.Sprintf("Unable to read cluster %s", attrName), err.Error())
		return nil
	}
	return ret
}

func getAttrTimeout(diags *diag.Diagnostics, rawState map[string]tftypes.Value) timeouts.Value {
	attrTypes := map[string]attr.Type{
		"create": types.StringType,
		"update": types.StringType,
		"delete": types.StringType,
	}
	nullObj := timeouts.Value{Object: types.ObjectNull(attrTypes)}
	timeoutState := getAttrFromRawState[map[string]tftypes.Value](diags, rawState, "timeouts")
	if diags.HasError() || timeoutState == nil {
		return nullObj
	}
	timeoutMap := make(map[string]attr.Value)
	for action := range attrTypes {
		actionTimeout := getAttrFromRawState[string](diags, *timeoutState, action)
		if actionTimeout == nil {
			timeoutMap[action] = types.StringNull()
		} else {
			timeoutMap[action] = types.StringPointerValue(actionTimeout)
		}
	}
	obj, d := types.ObjectValue(attrTypes, timeoutMap)
	diags.Append(d...)
	if diags.HasError() {
		return nullObj
	}
	return timeouts.Value{Object: obj}
}
