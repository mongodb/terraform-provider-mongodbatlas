package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type clusterAutoScalingMaxInstanceModifier struct{}

func ClusterAutoScalingMaxInstanceModifier() planmodifier.String {
	return clusterAutoScalingMaxInstanceModifier{}
}

func (m clusterAutoScalingMaxInstanceModifier) Description(_ context.Context) string {
	return "This planmodifier ensures that value of provider_auto_scaling_compute_max_instance_size is only considered in the plan if auto_scaling_compute_enabled is set to true"
}

func (m clusterAutoScalingMaxInstanceModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m clusterAutoScalingMaxInstanceModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var canScalePlan types.Bool
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("auto_scaling_compute_enabled"), &canScalePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if canScalePlan.ValueBool() && req.PlanValue != resp.PlanValue {
		resp.PlanValue = req.PlanValue
		return // do nothing, let the change be detected, if any
	}

	resp.PlanValue = req.StateValue // we want to ignore this value in the plan in this case
}
