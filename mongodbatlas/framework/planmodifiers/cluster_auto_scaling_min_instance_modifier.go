package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type clusterAutoScalingMinInstanceModifier struct{}

func ClusterAutoScalingMinInstanceModifier() planmodifier.String {
	return clusterAutoScalingMinInstanceModifier{}
}

func (m clusterAutoScalingMinInstanceModifier) Description(_ context.Context) string {
	return "This planmodifier ensures that value of provider_auto_scaling_compute_min_instance_size is" +
		"only considered in the plan if both auto_scaling_compute_scale_down_enabled and auto_scaling_compute_enabled are set to true"
}

func (m clusterAutoScalingMinInstanceModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m clusterAutoScalingMinInstanceModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var canScaleDown types.Bool
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("auto_scaling_compute_scale_down_enabled"), &canScaleDown)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var canScale types.Bool
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("auto_scaling_compute_enabled"), &canScale)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is a known planned value AND computeEnabled is true
	if !req.PlanValue.IsUnknown() && canScale.ValueBool() && canScaleDown.ValueBool() {
		return
	}

	if !canScaleDown.ValueBool() || !canScale.ValueBool() {
		resp.PlanValue = req.StateValue
	}

	// if canScaleDown.ValueBool() && canScale.ValueBool() && req.PlanValue != resp.PlanValue {
	// 	resp.PlanValue = req.PlanValue
	// 	return // do nothing, let the change be detected, if any
	// }

	resp.PlanValue = req.StateValue // we want to ignore this value in the plan in this case
}
