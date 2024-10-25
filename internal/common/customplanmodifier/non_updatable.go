package customplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NonUpdatableAttributePlanModifier(attribute string) planmodifier.String {
	return &nonUpdatableAttributePlanModifier{
		Attribute: attribute,
	}
}

type nonUpdatableAttributePlanModifier struct {
	Attribute string
}

func (d *nonUpdatableAttributePlanModifier) Description(ctx context.Context) string {
	return "Ensures that update operations fails when updating an attribute."
}

func (d *nonUpdatableAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var planAttributeValue types.String
	diags := req.Plan.GetAttribute(ctx, path.Root(d.Attribute), &planAttributeValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateAttributeValue types.String
	req.State.GetAttribute(ctx, path.Root(d.Attribute), &stateAttributeValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !stateAttributeValue.IsNull() && stateAttributeValue.ValueString() != planAttributeValue.ValueString() {
		resp.Diagnostics.AddError(
			"attribute is not updatable",
			"attribute is not updatable",
		)
		return
	}
}
