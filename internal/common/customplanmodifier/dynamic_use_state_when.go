package customplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DynamicEqualFunc compares two Dynamic values for semantic equality.
type DynamicEqualFunc func(a, b types.Dynamic) bool

// DynamicUseStateWhen suppresses no-op diffs on a Dynamic attribute. When the
// planned value is semantically equal to the prior state, the prior state is
// copied into the plan so Terraform does not flag formatting-only differences.
func DynamicUseStateWhen(eq DynamicEqualFunc) planmodifier.Dynamic {
	return &dynamicUseStateWhen{eq: eq}
}

type dynamicUseStateWhen struct {
	eq DynamicEqualFunc
}

func (m *dynamicUseStateWhen) Description(_ context.Context) string {
	return "Reuses prior state when the planned Dynamic value is semantically equal."
}

func (m *dynamicUseStateWhen) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m *dynamicUseStateWhen) PlanModifyDynamic(_ context.Context, req planmodifier.DynamicRequest, resp *planmodifier.DynamicResponse) {
	if req.StateValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if m.eq(req.PlanValue, req.StateValue) {
		resp.PlanValue = req.StateValue
	}
}
