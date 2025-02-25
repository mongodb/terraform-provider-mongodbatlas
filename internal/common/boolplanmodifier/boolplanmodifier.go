package boolplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefaultBoolValue is a plan modifier that sets a default boolean value
// if the configuration is null or unknown.
type DefaultBoolValue struct {
	Value bool
}

var _ planmodifier.Bool = DefaultBoolValue{}

// Description returns a plain text description of the plan modifier.
func (m DefaultBoolValue) Description(ctx context.Context) string {
	return "If the attribute is not specified, default to the given boolean value."
}

// MarkdownDescription returns a markdown formatted description.
func (m DefaultBoolValue) MarkdownDescription(ctx context.Context) string {
	return "If the attribute is not specified, default to the given boolean value."
}

// ModifyBool sets the planned value to the default if the configuration is null or unknown.
func (m DefaultBoolValue) ModifyBool(ctx context.Context, req planmodifier.BoolModifierRequest, resp *planmodifier.BoolModifierResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = types.BoolValue(m.Value)
	}
}
