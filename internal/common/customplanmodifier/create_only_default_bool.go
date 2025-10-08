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

// CreateOnlyBool creates a plan modifier that prevents updates to boolean attributes.
// This is useful for attributes only supported in create and not in update.
// It shows a helpful error message helping the user to update their config to match the state.
// Never use a schema.Default for create only attributes, instead use `WithDefault`, the default will lead to plan changes that are not expected after import.
// If the attribute is not in the API Response implement CopyFromPlan behavior when converting API Model to TF Model.
func CreateOnlyBool() planmodifier.Bool {
	return &createOnlyBoolPlanModifier{}
}

// CreateOnlyBoolWithDefault sets a default value on create operation that will show in the plan.
// This avoids any custom logic in the resource "Create" handler.
// On update the default has no impact and the UseStateForUnknown behavior is observed instead.
// Always use Optional+Computed when using a default value.
// If the attribute is not in the API Response implement CopyFromPlan behavior when converting API Model to TF Model.
func CreateOnlyBoolWithDefault(b bool) planmodifier.Bool {
	return &createOnlyBoolPlanModifier{defaultBool: &b}
}

type createOnlyBoolPlanModifier struct {
	defaultBool *bool
}

func (d *createOnlyBoolPlanModifier) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *createOnlyBoolPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures the update operation fails when updating an attribute. If the read after import doesn't equal the configuration value it will also raise an error."
}

// isCreate uses the full state to check if this is a create operation
func isCreate(t *tfsdk.State) bool {
	return t.Raw.IsNull()
}

func (d *createOnlyBoolPlanModifier) UseDefault() bool {
	return d.defaultBool != nil
}

func (d *createOnlyBoolPlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
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

// isUpdated checks if the attribute was updated.
// Special case when the attribute is removed/set to null in the plan:
// Computed Attribute: returns false (unknown in the plan)
// Optional Attribute: returns true if the state has a value
func isUpdated(state, plan attr.Value) bool {
	if !IsKnown(plan) {
		return false
	}
	return !state.Equal(plan)
}

func (d *createOnlyBoolPlanModifier) addDiags(diags *diag.Diagnostics, attrPath path.Path, stateValue attr.Value) {
	message := fmt.Sprintf("%s cannot be updated or set after import, remove it from the configuration or use the state value (see below).", attrPath)
	detail := fmt.Sprintf("The current state value is %s", stateValue)
	diags.AddError(message, detail)
}
