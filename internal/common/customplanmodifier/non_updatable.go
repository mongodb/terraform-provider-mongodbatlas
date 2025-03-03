package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type Modifier interface {
	planmodifier.String
	planmodifier.Bool
}

func NonUpdatableAttributePlanModifier() Modifier {
	return &nonUpdatableAttributePlanModifier{}
}

type nonUpdatableAttributePlanModifier struct {
}

func (d *nonUpdatableAttributePlanModifier) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *nonUpdatableAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures that update operations fails when updating an attribute."
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if d.isUpdated(req.PlanValue, req.StateValue) {
		d.addDiags(&resp.Diagnostics, req.Path)
	}
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if d.isUpdated(req.PlanValue, req.StateValue) {
		d.addDiags(&resp.Diagnostics, req.Path)
	}
}

func (d *nonUpdatableAttributePlanModifier) isUpdated(planValue, stateValue attr.Value) bool {
	return !stateValue.IsNull() && !planValue.Equal(stateValue)
}

func (d *nonUpdatableAttributePlanModifier) addDiags(diags *diag.Diagnostics, attrPath path.Path) {
	message := fmt.Sprintf("%s cannot be updated", attrPath)
	diags.AddError(message, message)
}
