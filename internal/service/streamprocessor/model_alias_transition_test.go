package streamprocessor

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestIsAliasOnlyTransition(t *testing.T) {
	legacyState := TFStreamProcessorRSModel{
		InstanceName:  types.StringValue("workspace"),
		WorkspaceName: types.StringNull(),
		ProjectID:     types.StringValue("project"),
		ProcessorName: types.StringValue("processor"),
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
