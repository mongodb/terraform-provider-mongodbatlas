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

type Modifier interface {
	planmodifier.String
	planmodifier.Bool
}

// CreateOnlyAttributePlanModifier returns a plan modifier that ensures that update operations fails when the attribute is changed.
// This is useful for attributes only supported in create and not in update.
// It shows a helpful error message helping the user to update their config to match the state.
// Never use a schema.Default for create only attributes, instead use WithXXXDefault, the default will lead to plan changes that are not expected after import.
// Implement CopyFromPlan if the attribute is not in the API Response.
func CreateOnlyAttributePlanModifier() Modifier {
	return &createOnlyAttributePlanModifier{}
}

func CreateOnlyAttributePlanModifierWithBoolDefault(b bool) Modifier {
	return &createOnlyAttributePlanModifier{defaultBool: &b}
}

type createOnlyAttributePlanModifier struct {
	defaultBool *bool
}

func (d *createOnlyAttributePlanModifier) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *createOnlyAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures the update operation fails when updating an attribute. If the read after import don't equal the configuration value it will also raise an error."
}

func isCreate(t *tfsdk.State) bool {
	return t.Raw.IsNull()
}

func (d *createOnlyAttributePlanModifier) UseDefault() bool {
	return d.defaultBool != nil
}

func (d *createOnlyAttributePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
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

func (d *createOnlyAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
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

func (d *createOnlyAttributePlanModifier) addDiags(diags *diag.Diagnostics, attrPath path.Path, stateValue attr.Value) {
	message := fmt.Sprintf("%s cannot be updated or set after import, remove it from the configuration or use the state value (see below).", attrPath)
	detail := fmt.Sprintf("The current state value is %s", stateValue)
	diags.AddError(message, detail)
}
