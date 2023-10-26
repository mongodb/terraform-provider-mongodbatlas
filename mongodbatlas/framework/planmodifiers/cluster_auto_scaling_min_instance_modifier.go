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
	// Do nothing if there is no state value.
	// if req.StateValue.IsNull() {
	// 	return
	// }

	// Do nothing if there is a known planned value.
	// if !req.PlanValue.IsUnknown() {
	// 	return
	// }

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	// if req.ConfigValue.IsUnknown() {
	// 	return
	// }
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

	if canScaleDown.ValueBool() && canScale.ValueBool() && req.PlanValue != resp.PlanValue {
		resp.PlanValue = req.PlanValue
		return // do nothing, let the change be detected, if any
	}

	resp.PlanValue = req.StateValue // we want to ignore this value in the plan in this case
}
