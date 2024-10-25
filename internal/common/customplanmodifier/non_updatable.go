package customplanmodifier

import (
	"context"
	"fmt"
	"strings"

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
	planFullPath, diag := req.Plan.PathMatches(ctx, GetFullPathExpression(ctx, d.Attribute))
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}
	diags := req.Plan.GetAttribute(ctx, planFullPath[0], &planAttributeValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateAttributeValue types.String
	stateFullPath, diag := req.Plan.PathMatches(ctx, GetFullPathExpression(ctx, d.Attribute))
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}
	req.State.GetAttribute(ctx, stateFullPath[0], &stateAttributeValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !stateAttributeValue.IsNull() && stateAttributeValue.ValueString() != planAttributeValue.ValueString() {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s cannot be updated", d.Attribute),
			fmt.Sprintf("%s cannot be updated", d.Attribute),
		)
		return
	}
}

func GetFullPathExpression(ctx context.Context, attribute string) path.Expression {
	parts := strings.Split(attribute, ".")
	pathExpression := path.MatchRelative()
	for _, part := range parts {
		pathExpression = pathExpression.AtName(part)
	}
	return pathExpression
}
