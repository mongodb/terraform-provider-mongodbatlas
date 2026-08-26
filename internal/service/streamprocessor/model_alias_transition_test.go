package streamprocessor

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.reverted, IsWorkspaceNameAliasReversion(tc.stateWorkspaceName, tc.planWorkspaceName, tc.planInstanceName))
		})
	}
}

func TestIsAliasOnlyTransition(t *testing.T) {
	legacyState := TFStreamProcessorRSModel{
		InstanceName:  types.StringValue("workspace"),
		WorkspaceName: types.StringNull(),
		ProjectID:     types.StringValue("project"),
		ProcessorName: types.StringValue("processor"),
		ProcessorID:   types.StringValue("processor-id"),
		Stats:         types.StringValue("{}"),
		Pipeline:      jsontypes.NewNormalizedValue(`[]`),
	}

	testCases := map[string]struct {
		plan                TFStreamProcessorRSModel
		state               TFStreamProcessorRSModel
		aliasOnlyTransition bool
	}{
		"legacy_to_canonical": {
			state: legacyState,
			plan: TFStreamProcessorRSModel{
				InstanceName:  types.StringNull(),
				WorkspaceName: types.StringValue("workspace"),
				ProjectID:     types.StringValue("project"),
				ProcessorName: types.StringValue("processor"),
				Pipeline:      jsontypes.NewNormalizedValue(`[]`),
			},
			aliasOnlyTransition: true,
		},
		"computed_values_unknown_in_plan": {
			state: legacyState,
			plan: TFStreamProcessorRSModel{
				InstanceName:  types.StringNull(),
				WorkspaceName: types.StringValue("workspace"),
				ProjectID:     types.StringValue("project"),
				ProcessorName: types.StringValue("processor"),
				ProcessorID:   types.StringUnknown(),
				Stats:         types.StringUnknown(),
				Pipeline:      jsontypes.NewNormalizedValue(`[]`),
			},
			aliasOnlyTransition: true,
		},
		"different_workspace": {
			state: legacyState,
			plan: TFStreamProcessorRSModel{
				InstanceName:  types.StringNull(),
				WorkspaceName: types.StringValue("different-workspace"),
				ProjectID:     types.StringValue("project"),
				ProcessorName: types.StringValue("processor"),
				Pipeline:      jsontypes.NewNormalizedValue(`[]`),
			},
		},
		"other_attribute_changed": {
			state: legacyState,
			plan: TFStreamProcessorRSModel{
				InstanceName:  types.StringNull(),
				WorkspaceName: types.StringValue("workspace"),
				ProjectID:     types.StringValue("project"),
				ProcessorName: types.StringValue("other-processor"),
				Pipeline:      jsontypes.NewNormalizedValue(`[]`),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.aliasOnlyTransition, IsAliasOnlyTransition(tc.plan, tc.state))
		})
	}
}
