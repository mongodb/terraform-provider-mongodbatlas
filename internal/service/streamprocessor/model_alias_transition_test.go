package streamprocessor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
)

func TestIsAliasOnlyTransition(t *testing.T) {
	legacyState := streamprocessor.TFStreamProcessorRSModel{
		InstanceName:  types.StringValue("workspace"),
		WorkspaceName: types.StringNull(),
		ProjectID:     types.StringValue("project"),
		ProcessorName: types.StringValue("processor"),
		ProcessorID:   types.StringValue("processor-id"),
		Stats:         types.StringValue("{}"),
		Pipeline:      jsontypes.NewNormalizedValue(`[]`),
	}

	testCases := map[string]struct {
		plan                streamprocessor.TFStreamProcessorRSModel
		state               streamprocessor.TFStreamProcessorRSModel
		aliasOnlyTransition bool
	}{
		"legacy_to_canonical": {
			state: legacyState,
			plan: streamprocessor.TFStreamProcessorRSModel{
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
			plan: streamprocessor.TFStreamProcessorRSModel{
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
		"optional_computed_values_unknown_in_plan": {
			state: legacyState,
			plan: streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringNull(),
				WorkspaceName: types.StringValue("workspace"),
				ProjectID:     types.StringValue("project"),
				ProcessorName: types.StringValue("processor"),
				Pipeline:      jsontypes.NewNormalizedValue(`[]`),
				State:         types.StringUnknown(),
				Tier:          types.StringUnknown(),
			},
			aliasOnlyTransition: true,
		},
		"legacy_dual_alias_to_canonical": {
			state: streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringValue("workspace"),
				WorkspaceName: types.StringValue("workspace"),
				ProjectID:     types.StringValue("project"),
				ProcessorName: types.StringValue("processor"),
				Pipeline:      jsontypes.NewNormalizedValue(`[]`),
			},
			plan: streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringNull(),
				WorkspaceName: types.StringValue("workspace"),
				ProjectID:     types.StringValue("project"),
				ProcessorName: types.StringValue("processor"),
				Pipeline:      jsontypes.NewNormalizedValue(`[]`),
			},
			aliasOnlyTransition: true,
		},
		"different_workspace": {
			state: legacyState,
			plan: streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringNull(),
				WorkspaceName: types.StringValue("different-workspace"),
				ProjectID:     types.StringValue("project"),
				ProcessorName: types.StringValue("processor"),
				Pipeline:      jsontypes.NewNormalizedValue(`[]`),
			},
		},
		"other_attribute_changed": {
			state: legacyState,
			plan: streamprocessor.TFStreamProcessorRSModel{
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
			assert.Equal(t, tc.aliasOnlyTransition, streamprocessor.IsAliasOnlyTransition(&tc.plan, &tc.state))
		})
	}
}
