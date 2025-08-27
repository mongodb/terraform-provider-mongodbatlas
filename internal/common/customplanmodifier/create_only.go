package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// CreateOnlyStringPlanModifier creates a plan modifier that prevents updates to string attributes.
func CreateOnlyStringPlanModifier() planmodifier.String {
	return &createOnlyAttributePlanModifier{}
}

// CreateOnlyBoolPlanModifier creates a plan modifier that prevents updates to boolean attributes.
func CreateOnlyBoolPlanModifier() planmodifier.Bool {
	return &createOnlyAttributePlanModifier{}
}

// Plan modifier that implements create-only behavior for multiple attribute types
type createOnlyAttributePlanModifier struct{}

func (d *createOnlyAttributePlanModifier) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *createOnlyAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures that update operations fail when attempting to modify a create-only attribute."
}

func (d *createOnlyAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *createOnlyAttributePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

// validateCreateOnly checks if an attribute value has changed and adds an error if it has
func validateCreateOnly(planValue, stateValue attr.Value, attrPath path.Path, diagnostics *diag.Diagnostics,
) {
	if !stateValue.IsNull() && !stateValue.Equal(planValue) {
		diagnostics.AddError(
			fmt.Sprintf("%s cannot be updated", attrPath),
			fmt.Sprintf("%s cannot be updated", attrPath),
		)
	}
}
