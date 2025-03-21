package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
)

var attributePlanModifiers = map[string]customplanmodifier.UnknownReplacementCall[PlanModifyResourceInfo]{
	"mongo_db_version": mongoDBVersionReplaceUnknown,
}

func mongoDBVersionReplaceUnknown(ctx context.Context, state customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) attr.Value {
	if req.Changes.AttributeChanged("mongo_db_version") {
		return req.Unknown
	}
	return state.Value
}

type PlanModifyResourceInfo struct {
	AutoScalingComputedUsed bool
	AutoScalingDiskUsed     bool
	isShardingConfigUpgrade bool
}

func unknownReplacements(ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state TFModel
	diags := &resp.Diagnostics
	diags.Append(req.Plan.Get(ctx, &plan)...)
	diags.Append(req.State.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	computedUsed, diskUsed := autoScalingUsed(ctx, diags, &state, &plan)
	shardingConfigUpgrade := isShardingConfigUpgrade(ctx, &state, &plan, diags)
	info := PlanModifyResourceInfo{
		AutoScalingComputedUsed: computedUsed,
		AutoScalingDiskUsed:     diskUsed,
		isShardingConfigUpgrade: shardingConfigUpgrade,
	}
	unknownReplacements := customplanmodifier.NewUnknownReplacements(ctx, req, resp, resourceSchema(ctx), info)
	for attrName, replacer := range attributePlanModifiers {
		unknownReplacements.AddReplacement(attrName, replacer)
	}
	unknownReplacements.ApplyReplacments(ctx, diags)
}
