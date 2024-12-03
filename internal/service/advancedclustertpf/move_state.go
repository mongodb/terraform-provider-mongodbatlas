package advancedclustertpf

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

func (r *rs) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{{StateMover: stateMover}}
}

func stateMover(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if req.SourceTypeName != "mongodbatlas_cluster" || !strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas") {
		return
	}
	rawStateValue, err := req.SourceRawState.UnmarshalWithOpts(tftypes.Object{
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
	setMoveStateResponse(ctx, *projectID, *name, resp)
}

func setMoveStateResponse(ctx context.Context, projectID, clusterName string, resp *resource.MoveStateResponse) {
	validTimeout := timeouts.Value{
		Object: types.ObjectValueMust(
			map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"delete": types.StringType,
			},
			map[string]attr.Value{
				"create": types.StringValue("30m"),
				"update": types.StringValue("30m"),
				"delete": types.StringValue("30m"),
			}),
	}
	model := NewTFModel(ctx, &admin.ClusterDescription20240805{
		GroupId: conversion.StringPtr(projectID),
		Name:    conversion.StringPtr(clusterName),
	}, validTimeout, &resp.Diagnostics, nil)
	if resp.Diagnostics.HasError() {
		return
	}
	AddAdvancedConfig(ctx, model, nil, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.TargetState.Set(ctx, model)...)
}
