package customplanmodifier

import (
	"context"
	"fmt"

	planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func NonUpdatableStringAttributePlanModifier(attribute string) planmodifier.String {
	return &nonUpdatableStringAttributePlanModifier{
		Attribute: attribute,
	}
}

type nonUpdatableStringAttributePlanModifier struct {
	Attribute string
}

func (d *nonUpdatableStringAttributePlanModifier) Description(ctx context.Context) string {
	return "Ensures that update operations fails when updating an attribute."
}

func (d *nonUpdatableStringAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *nonUpdatableStringAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	planAttributeValue := req.PlanValue
	stateAttributeValue := req.StateValue

	if !stateAttributeValue.IsNull() && stateAttributeValue.ValueString() != planAttributeValue.ValueString() {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s cannot be updated", d.Attribute),
			fmt.Sprintf("%s cannot be updated", d.Attribute),
		)
		return
	}
}
