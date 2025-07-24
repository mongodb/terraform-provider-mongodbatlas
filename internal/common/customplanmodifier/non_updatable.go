package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// NonUpdatableStringAttributePlanModifier creates a plan modifier that prevents updates to string attributes.
func NonUpdatableStringAttributePlanModifier() planmodifier.String {
	return &nonUpdatableAttributePlanModifier{}
}

// NonUpdatableBoolAttributePlanModifier creates a plan modifier that prevents updates to boolean attributes.
func NonUpdatableBoolAttributePlanModifier() planmodifier.Bool {
	return &nonUpdatableAttributePlanModifier{}
}

// Plan modifier that implements non-updatable behavior for multiple attribute types
type nonUpdatableAttributePlanModifier struct{}

func (d *nonUpdatableAttributePlanModifier) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *nonUpdatableAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures that update operations fails when updating an attribute."
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func validateNonUpdatable(planValue, stateValue attr.Value, attrPath path.Path, diagnostics *diag.Diagnostics,
) {
	if !stateValue.IsNull() && !stateValue.Equal(planValue) {
		diagnostics.AddError(
			fmt.Sprintf("%s cannot be updated", attrPath),
			fmt.Sprintf("%s cannot be updated", attrPath),
		)
	}
}
