package streamalias

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReversionPlanModifier prevents resources with canonical workspace_name-only
// state from reverting to deprecated instance_name.
type ReversionPlanModifier struct{}

func (ReversionPlanModifier) Description(_ context.Context) string {
	return "Prevents reverting from workspace_name to the deprecated instance_name attribute."
}

func (m ReversionPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (ReversionPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var planWorkspaceName, planInstanceName, stateWorkspaceName, stateInstanceName types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("workspace_name"), &planWorkspaceName)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("instance_name"), &planInstanceName)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("workspace_name"), &stateWorkspaceName)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("instance_name"), &stateInstanceName)...)
	if resp.Diagnostics.HasError() || !IsWorkspaceNameAliasReversion(stateWorkspaceName, stateInstanceName, planWorkspaceName, planInstanceName) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		path.Root("instance_name"),
		"Cannot revert stream workspace alias",
		"This resource already uses workspace_name in state. Use workspace_name instead of the deprecated instance_name attribute.",
	)
}
