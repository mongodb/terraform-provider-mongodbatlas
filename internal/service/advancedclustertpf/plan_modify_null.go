package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseStateForNull() planmodifier.Object {
	return useStateForNull{}
}

// useStateForNull implements the plan modifier.
type useStateForNull struct{}

// Description returns a human-readable description of the plan modifier.
func (m useStateForNull) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m useStateForNull) MarkdownDescription(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// PlanModifyObject implements the plan modification logic.
func (m useStateForNull) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// Override PlanValue with Unknown if ConfigValue is null. (NOT ALLOWED)
	if req.PlanValue.IsNull() {
		model := TFAdvancedConfigurationModel{
			OplogMinRetentionHours: types.Float64Unknown(),
		}
		resp.PlanValue = NewObjectValueOfMust(ctx, &model).ObjectValue
		// resp.PlanValue = NewObjectValueOfUnknown[TFAdvancedConfigurationModel](ctx).ObjectValue ALSO failed
		return
	}
}
