package customplanmodifier_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
	"github.com/stretchr/testify/assert"
)

func TestIsAliasReversion(t *testing.T) {
	testCases := map[string]struct {
		stateCanonical   types.String
		stateDeprecated  types.String
		planCanonical    types.String
		planDeprecated   types.String
		isAliasReversion bool
	}{
		"canonical_to_deprecated_is_rejected": {
			stateCanonical:   types.StringValue("value"),
			stateDeprecated:  types.StringNull(),
			planCanonical:    types.StringNull(),
			planDeprecated:   types.StringValue("value"),
			isAliasReversion: true,
		},
		"deprecated_to_canonical_is_allowed": {
			stateCanonical:  types.StringNull(),
			stateDeprecated: types.StringValue("value"),
			planCanonical:   types.StringValue("value"),
			planDeprecated:  types.StringNull(),
		},
		"legacy_state_with_both_aliases_is_allowed": {
			stateCanonical:  types.StringValue("value"),
			stateDeprecated: types.StringValue("value"),
			planCanonical:   types.StringNull(),
			planDeprecated:  types.StringValue("value"),
		},
		"unknown_deprecated_plan_value_is_allowed": {
			stateCanonical:  types.StringValue("value"),
			stateDeprecated: types.StringNull(),
			planCanonical:   types.StringNull(),
			planDeprecated:  types.StringUnknown(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.isAliasReversion, customplanmodifier.IsAliasReversion(tc.stateCanonical, tc.stateDeprecated, tc.planCanonical, tc.planDeprecated))
		})
	}
}

func TestAliasReversionDescription(t *testing.T) {
	modifier := customplanmodifier.AliasReversion("canonical_name", "deprecated_name")

	assert.Equal(t, "Prevents reverting from canonical_name to deprecated_name.", modifier.Description(context.Background()))
}
