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

// CreateOnlyAttributePlanModifierWithBoolDefault sets a default value on create operation that will show in the plan.
// This avoids any custom logic in the resource "Create" handler.
// On update the default has no impact and the UseStateForUnknown behavior is observed instead.
// Always use Optional+Computed when using a default value.
func CreateOnlyAttributePlanModifierWithBoolDefault(b bool) planmodifier.Bool {
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
