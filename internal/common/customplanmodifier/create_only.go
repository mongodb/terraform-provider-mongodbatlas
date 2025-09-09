package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type CreateOnlyModifier interface {
	planmodifier.String
	planmodifier.Bool
}

// CreateOnlyAttributePlanModifier returns a plan modifier that ensures that update operations fails when the attribute is changed.
// This is useful for attributes only supported in create and not in update.
// It shows a helpful error message helping the user to update their config to match the state.
// Never use a schema.Default for create only attributes, instead use WithXXXDefault, the default will lead to plan changes that are not expected after import.
// Implement CopyFromPlan if the attribute is not in the API Response.
func CreateOnlyAttributePlanModifier() CreateOnlyModifier {
	return &createOnlyAttributePlanModifier{}
}

// CreateOnlyAttributePlanModifierWithBoolDefault sets a default value on create operation that will show in the plan.
// This avoids any custom logic in the resource "Create" handler.
// On update the default has no impact and the UseStateForUnknown behavior is observed instead.
// Always use Optional+Computed when using a default value.
func CreateOnlyAttributePlanModifierWithBoolDefault(b bool) CreateOnlyModifier {
	return &createOnlyAttributePlanModifierWithBoolDefault{defaultBool: &b}
}

type createOnlyAttributePlanModifierWithBoolDefault struct {
	defaultBool *bool
}

func (d *createOnlyAttributePlanModifierWithBoolDefault) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *createOnlyAttributePlanModifierWithBoolDefault) MarkdownDescription(ctx context.Context) string {
	return "Ensures the update operation fails when updating an attribute. If the read after import don't equal the configuration value it will also raise an error."
}

func isCreate(t *tfsdk.State) bool {
	return t.Raw.IsNull()
}

func (d *createOnlyAttributePlanModifierWithBoolDefault) UseDefault() bool {
	return d.defaultBool != nil
}

func (d *createOnlyAttributePlanModifierWithBoolDefault) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if isCreate(&req.State) {
		if !IsKnown(req.PlanValue) && d.UseDefault() {
			resp.PlanValue = types.BoolPointerValue(d.defaultBool)
		}
		return
	}
	if isUpdated(req.StateValue, req.PlanValue) {
		d.addDiags(&resp.Diagnostics, req.Path, req.StateValue)
	}
	if !IsKnown(req.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (d *createOnlyAttributePlanModifierWithBoolDefault) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
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

func isUpdated(state, plan attr.Value) bool {
	if !IsKnown(plan) {
		return false
	}
	return !state.Equal(plan)
}

func (d *createOnlyAttributePlanModifierWithBoolDefault) addDiags(diags *diag.Diagnostics, attrPath path.Path, stateValue attr.Value) {
	message := fmt.Sprintf("%s cannot be updated or set after import, remove it from the configuration or use the state value (see below).", attrPath)
	detail := fmt.Sprintf("The current state value is %s", stateValue)
	diags.AddError(message, detail)
}
