package streamalias

import "github.com/hashicorp/terraform-plugin-framework/types"

// IsWorkspaceNameAliasReversion reports whether a resource with canonical
// workspace_name-only state is being changed back to deprecated instance_name.
func IsWorkspaceNameAliasReversion(stateWorkspaceName, stateInstanceName, planWorkspaceName, planInstanceName types.String) bool {
	// Legacy imports populated both aliases. Only workspace_name-only state is
	// canonical and therefore subject to the one-way migration restriction.
	return !stateWorkspaceName.IsNull() && !stateWorkspaceName.IsUnknown() &&
		stateInstanceName.IsNull() &&
		planWorkspaceName.IsNull() &&
		!planInstanceName.IsNull() && !planInstanceName.IsUnknown()
}
