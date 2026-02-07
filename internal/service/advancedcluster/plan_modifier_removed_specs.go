package advancedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ModifyPlan validates that read_only_specs and analytics_specs are not silently removed during migration from schema v2 to v3.
func (r *rs) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var state, plan TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	validateSpecsNotRemoved(ctx, &resp.Diagnostics, &state, &plan)
}

func validateSpecsNotRemoved(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	stateRepSpecs := tfModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecs := tfModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return
	}
	for i := range min(len(stateRepSpecs), len(planRepSpecs)) {
		stateRegionConfigs := tfModelList[TFRegionConfigsModel](ctx, diags, stateRepSpecs[i].RegionConfigs)
		planRegionConfigs := tfModelList[TFRegionConfigsModel](ctx, diags, planRepSpecs[i].RegionConfigs)
		if diags.HasError() {
			return
		}
		for j := range min(len(stateRegionConfigs), len(planRegionConfigs)) {
			if specsRemovedWithActiveNodes(ctx, stateRegionConfigs[j].ReadOnlySpecs, planRegionConfigs[j].ReadOnlySpecs) {
				diags.AddError(
					"read_only_specs cannot be removed when read-only nodes are active",
					"Your cluster has active read-only nodes but the read_only_specs block was removed from the configuration. "+
						"To keep read-only nodes, add the read_only_specs block back. "+
						"To delete them, set node_count = 0 instead of removing the block.",
				)
				return
			}
			if specsRemovedWithActiveNodes(ctx, stateRegionConfigs[j].AnalyticsSpecs, planRegionConfigs[j].AnalyticsSpecs) {
				diags.AddError(
					"analytics_specs cannot be removed when analytics nodes are active",
					"Your cluster has active analytics nodes but the analytics_specs block was removed from the configuration. "+
						"To keep analytics nodes, add the analytics_specs block back. "+
						"To delete them, set node_count = 0 instead of removing the block.",
				)
				return
			}
		}
	}
}

// specsRemovedWithActiveNodes returns true if state has specs with node_count > 0 but plan has null specs.
func specsRemovedWithActiveNodes(ctx context.Context, stateSpecs, planSpecs types.Object) bool {
	stateModel := tfModelObject[TFSpecsModel](ctx, stateSpecs)
	if stateModel == nil || stateModel.NodeCount.ValueInt64() == 0 {
		return false
	}
	return tfModelObject[TFSpecsModel](ctx, planSpecs) == nil
}

func tfModelList[T any](ctx context.Context, diags *diag.Diagnostics, input types.List) []T {
	elements := make([]T, len(input.Elements()))
	diags.Append(input.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return nil
	}
	return elements
}

func tfModelObject[T any](ctx context.Context, input types.Object) *T {
	item := new(T)
	if diags := input.As(ctx, item, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil
	}
	return item
}
