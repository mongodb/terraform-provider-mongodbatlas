package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseNullForUnknown* returns a plan modifier that applies if the attribute has null value in the state and configuration
// and the plan value is unknown. This scenario typically occurs for Optional & Computed attributes that are null in state and configuration
// and the API doesn't return a value for these as well. The Framework will mark values for such attributes to an unknown value "(known after apply)"
// during planning, this plan modifier changes that behavior to not detect any unexpected drift allowing the value to stay null.
func UseNullForUnknownString() planmodifier.String {
	return useNullForUnknownStringModifier{}
}

// useStateForUnknownModifier implements the plan modifier.
type useNullForUnknownStringModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m useNullForUnknownStringModifier) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m useNullForUnknownStringModifier) MarkdownDescription(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// PlanModifyInt64 implements the plan modification logic.
func (m useNullForUnknownStringModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() && req.ConfigValue.IsNull() && req.PlanValue.IsUnknown() {
		resp.PlanValue = types.StringNull()
		return
	}

	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue = req.StateValue
}
