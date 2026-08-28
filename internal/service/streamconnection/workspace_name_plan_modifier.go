package streamconnection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamalias"
)

// WorkspaceNameRequiresReplace requires replacement only when the effective workspace
// identifier changes. instance_name and workspace_name are mutually exclusive aliases,
// so changing between aliases with the same value must not replace the connection.
type WorkspaceNameRequiresReplace struct{}

func (m WorkspaceNameRequiresReplace) Description(_ context.Context) string {
	return "Requires replacement when the effective stream workspace name changes."
}

func (m WorkspaceNameRequiresReplace) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m WorkspaceNameRequiresReplace) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var planWorkspaceName, planInstanceName, stateWorkspaceName, stateInstanceName types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("workspace_name"), &planWorkspaceName)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("instance_name"), &planInstanceName)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("workspace_name"), &stateWorkspaceName)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("instance_name"), &stateInstanceName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if req.Path.Equal(path.Root("instance_name")) && streamalias.IsWorkspaceNameAliasReversion(stateWorkspaceName, stateInstanceName, planWorkspaceName, planInstanceName) {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_name"),
			"Cannot revert stream workspace alias",
			"This resource already uses workspace_name in state. Use workspace_name instead of the deprecated instance_name attribute.",
		)
		return
	}

	if RequiresWorkspaceNameReplacement(stateWorkspaceName, stateInstanceName, planWorkspaceName, planInstanceName) {
		resp.RequiresReplace = true
	}
}

func RequiresWorkspaceNameReplacement(stateWorkspaceName, stateInstanceName, planWorkspaceName, planInstanceName types.String) bool {
	stateName, stateKnown := effectiveWorkspaceName(stateWorkspaceName, stateInstanceName)
	if !stateKnown || stateName == "" {
		return false
	}

	planName, planKnown := effectiveWorkspaceName(planWorkspaceName, planInstanceName)
	if !planKnown {
		return true
	}

	return planName != stateName
}

func effectiveWorkspaceName(workspaceName, instanceName types.String) (string, bool) {
	if workspaceName.IsUnknown() {
		return "", false
	}
	if !workspaceName.IsNull() {
		return workspaceName.ValueString(), true
	}
	if instanceName.IsUnknown() {
		return "", false
	}
	if !instanceName.IsNull() {
		return instanceName.ValueString(), true
	}
	return "", true
}
