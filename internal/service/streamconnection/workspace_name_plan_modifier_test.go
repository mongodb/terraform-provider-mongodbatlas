package streamconnection

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestIsWorkspaceNameAliasReversion(t *testing.T) {
	testCases := map[string]struct {
		stateWorkspaceName types.String
		planWorkspaceName  types.String
		planInstanceName   types.String
		reverted           bool
	}{
		"canonical_to_legacy_is_rejected": {
			stateWorkspaceName: types.StringValue("workspace"),
			planWorkspaceName:  types.StringNull(),
			planInstanceName:   types.StringValue("workspace"),
			reverted:           true,
		},
		"legacy_to_canonical_is_allowed": {
			stateWorkspaceName: types.StringNull(),
			planWorkspaceName:  types.StringValue("workspace"),
			planInstanceName:   types.StringNull(),
		},
		"unknown_instance_name_is_not_rejected": {
			stateWorkspaceName: types.StringValue("workspace"),
			planWorkspaceName:  types.StringNull(),
			planInstanceName:   types.StringUnknown(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.reverted, isWorkspaceNameAliasReversion(tc.stateWorkspaceName, tc.planWorkspaceName, tc.planInstanceName))
		})
	}
}

func TestRequiresWorkspaceNameReplacement(t *testing.T) {
	testCases := map[string]struct {
		stateWorkspaceName types.String
		stateInstanceName  types.String
		planWorkspaceName  types.String
		planInstanceName   types.String
		requiresReplace    bool
	}{
		"legacy_to_canonical_with_same_name": {
			stateWorkspaceName: types.StringNull(),
			stateInstanceName:  types.StringValue("workspace"),
			planWorkspaceName:  types.StringValue("workspace"),
			planInstanceName:   types.StringNull(),
		},
		"canonical_to_legacy_with_same_name": {
			stateWorkspaceName: types.StringValue("workspace"),
			stateInstanceName:  types.StringNull(),
			planWorkspaceName:  types.StringNull(),
			planInstanceName:   types.StringValue("workspace"),
		},
		"workspace_name_change_requires_replacement": {
			stateWorkspaceName: types.StringValue("old-workspace"),
			stateInstanceName:  types.StringNull(),
			planWorkspaceName:  types.StringValue("new-workspace"),
			planInstanceName:   types.StringNull(),
			requiresReplace:    true,
		},
		"initial_creation_does_not_require_replacement": {
			stateWorkspaceName: types.StringNull(),
			stateInstanceName:  types.StringNull(),
			planWorkspaceName:  types.StringValue("workspace"),
			planInstanceName:   types.StringNull(),
		},
		"unknown_planned_name_requires_replacement": {
			stateWorkspaceName: types.StringValue("workspace"),
			stateInstanceName:  types.StringNull(),
			planWorkspaceName:  types.StringUnknown(),
			planInstanceName:   types.StringNull(),
			requiresReplace:    true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.requiresReplace, requiresWorkspaceNameReplacement(tc.stateWorkspaceName, tc.stateInstanceName, tc.planWorkspaceName, tc.planInstanceName))
		})
	}
}
