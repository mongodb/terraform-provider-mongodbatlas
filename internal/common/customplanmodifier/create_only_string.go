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
// This is useful for attributes only supported in create and not in update.
// It shows a helpful error message helping the user to update their config to match the state.
// Never use a schema.Default for create only attributes, instead use `WithDefault`, the default will lead to plan changes that are not expected after import.
// No default value implemented for string until we have a use case.
// Implement CopyFromPlan if the attribute is not in the API Response.
func CreateOnlyStringPlanModifier() planmodifier.String {
	return &createOnlyStringPlanModifier{}
}

type createOnlyStringPlanModifier struct{}

func (d *createOnlyStringPlanModifier) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *createOnlyStringPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures the update operation fails when updating an attribute. If the read after import don't equal the configuration value it will also raise an error."
}

func (d *createOnlyStringPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if isCreate(&req.State) {
		return
	}
	if isUpdated(req.StateValue, req.PlanValue) {
		d.addDiags(&resp.Diagnostics, req.Path, req.StateValue)
	}
	if !IsKnown(req.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (d *createOnlyStringPlanModifier) addDiags(diags *diag.Diagnostics, attrPath path.Path, stateValue attr.Value) {
	message := fmt.Sprintf("%s cannot be updated or set after import, remove it from the configuration or use the state value (see below).", attrPath)
	detail := fmt.Sprintf("The current state value is %s", stateValue)
	diags.AddError(message, detail)
}
