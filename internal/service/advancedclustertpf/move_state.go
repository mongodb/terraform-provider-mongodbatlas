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
	projectID, name := getProjectIDClusterNameFromRawState(diags, req.SourceRawState)
	if diags.HasError() {
		return
	}
	setStateResponse(ctx, diags, &resp.TargetState, projectID, name)
}

// getProjectIDClusterNameFromRawState is used in Move State and Upgrade State
func getProjectIDClusterNameFromRawState(diags *diag.Diagnostics, state *tfprotov6.RawState) (projectID, name string) {
	rawStateValue, err := state.UnmarshalWithOpts(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"project_id": tftypes.String,
			"name":       tftypes.String,
		},
	}, tfprotov6.UnmarshalOpts{ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}})
	if err != nil {
		diags.AddError("Unable to Unmarshal state", err.Error())
		return "", ""
	}
	var rawState map[string]tftypes.Value
	if err := rawStateValue.As(&rawState); err != nil {
		diags.AddError("Unable to Parse state", err.Error())
		return "", ""
	}
	var projectIDPtr *string
	if err := rawState["project_id"].As(&projectIDPtr); err != nil {
		diags.AddAttributeError(path.Root("project_id"), "Unable to read cluster project_id", err.Error())
		return "", ""
	}
	var namePtr *string
	if err := rawState["name"].As(&namePtr); err != nil {
		diags.AddAttributeError(path.Root("name"), "Unable to read cluster name", err.Error())
		return "", ""
	}
	projectID, name = conversion.SafeString(projectIDPtr), conversion.SafeString(namePtr)
	if projectID == "" || name == "" {
		diags.AddError("Unable to read project_id or name", fmt.Sprintf("project_id: %s, name: %s", projectID, name))
		return "", ""
	}
	return projectID, name
}

// setStateResponse is used in Move State and Upgrade State
func setStateResponse(ctx context.Context, diags *diag.Diagnostics, state *tfsdk.State, projectID, clusterName string) {
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
	}, validTimeout, diags, ExtraAPIInfo{})
	if diags.HasError() {
		return
	}
	AddAdvancedConfig(ctx, model, nil, nil, diags)
	if diags.HasError() {
		return
	}
	diags.Append(state.Set(ctx, model)...)
}
