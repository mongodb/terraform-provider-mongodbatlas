package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func (r *rs) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: {StateUpgrader: stateUpgraderFromV1},
	}
}

func stateUpgraderFromV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	rawStateValue, err := req.RawState.UnmarshalWithOpts(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"project_id": tftypes.String,
			"name":       tftypes.String,
		},
	}, tfprotov6.UnmarshalOpts{ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Unmarshal Source State", err.Error())
		return
	}
	var rawState map[string]tftypes.Value
	if err := rawStateValue.As(&rawState); err != nil {
		resp.Diagnostics.AddError("Unable to Convert Source State", err.Error())
		return
	}
	var projectID *string
	if err := rawState["project_id"].As(&projectID); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("project_id"), "Unable to read cluster project_id", err.Error())
		return
	}
	var name *string
	if err := rawState["name"].As(&name); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Unable to Convert read cluster name", err.Error())
		return
	}

	if !conversion.IsStringPresent(projectID) || !conversion.IsStringPresent(name) {
		resp.Diagnostics.AddError("Unable to read project_id or name", "")
		return
	}
	setUpgradeStateResponse(ctx, *projectID, *name, resp)
}

func setUpgradeStateResponse(ctx context.Context, projectID, clusterName string, resp *resource.UpgradeStateResponse) {
	validTimeout := timeouts.Value{
		Object: types.ObjectNull(
			map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"delete": types.StringType,
			}),
	}
	model := NewTFModel(ctx, &admin.ClusterDescription20240805{
		GroupId: conversion.StringPtr(projectID),
		Name:    conversion.StringPtr(clusterName),
	}, validTimeout, &resp.Diagnostics, ExtraAPIInfo{})
	if resp.Diagnostics.HasError() {
		return
	}
	AddAdvancedConfig(ctx, model, nil, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
