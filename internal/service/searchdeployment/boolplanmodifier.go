package boolplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefaultBoolValue is a plan modifier that sets a default boolean value if none is provided.
type DefaultBoolValue struct {
	Value bool
}

var _ planmodifier.Bool = DefaultBoolValue{}

// Description returns a plain text description of the plan modifier.
func (m DefaultBoolValue) Description(ctx context.Context) string {
	return "If the attribute is not specified, default to the given boolean value."
}

// MarkdownDescription returns a markdown formatted description of the plan modifier.
func (m DefaultBoolValue) MarkdownDescription(ctx context.Context) string {
	return "If the attribute is not specified, default to the given boolean value."
}

// Modify sets the plan value to the default if the config value is null or unknown.
func (m DefaultBoolValue) Modify(ctx context.Context, req planmodifier.ModifyAttributeRequest, resp *planmodifier.ModifyAttributeResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = types.BoolValue(m.Value)
	}
}
