package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
)

var attributePlanModifiers = map[string]customplanmodifier.UnknownReplacementCall[PlanModifyResourceInfo]{
	"mongo_db_version":       mongoDBVersionReplaceUnknown,
	"auto_scaling":           autoScalingReplaceUnknown,
	"analytics_auto_scaling": autoScalingReplaceUnknown,
	// TODO: Add the other computed attributes
}

func autoScalingReplaceUnknown(ctx context.Context, state customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) attr.Value {
	// don't use auto_scaling or analytics_auto_scaling from state if it's not enabled as it doesn't need to be present in Update request payload
	if req.Info.AutoScalingComputedUsed || req.Info.AutoScalingDiskUsed {
		return state.AsObject()
	}
	return req.Unknown
}

func mongoDBVersionReplaceUnknown(ctx context.Context, state customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) attr.Value {
	if req.Changes.AttributeChanged("mongo_db_major_version") {
		return req.Unknown
	}
	return state.Value
}

type PlanModifyResourceInfo struct {
	AutoScalingComputedUsed bool
	AutoScalingDiskUsed     bool
	isShardingConfigUpgrade bool
}

func replicationSpecsKeepUnknownWhenChanged(ctx context.Context, state customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) []string {
	if !conversion.HasPathParent(req.Path, "replication_specs") {
		return nil
	}
	if req.Changes.AttributeChanged("replication_specs") {
		return []string{req.AttributeName}
	}
	return nil
}

func unknownReplacements(ctx context.Context, tfsdkState *tfsdk.State, tfsdkPlan *tfsdk.Plan, diags *diag.Diagnostics) {
	var plan, state TFModel
	diags.Append(tfsdkState.Get(ctx, &state)...)
	diags.Append(tfsdkPlan.Get(ctx, &plan)...)
	if diags.HasError() {
		return
	}
	diff := findClusterDiff(ctx, &state, &plan, diags)
	if diags.HasError() || diff.isAnyUpgrade() { // Don't do anything in upgrades
		return
	}
	computedUsed, diskUsed := autoScalingUsed(ctx, diags, &state, &plan)
	shardingConfigUpgrade := isShardingConfigUpgrade(ctx, &state, &plan, diags)
	info := PlanModifyResourceInfo{
		AutoScalingComputedUsed: computedUsed,
		AutoScalingDiskUsed:     diskUsed,
		isShardingConfigUpgrade: shardingConfigUpgrade,
	}
	unknownReplacements := customplanmodifier.NewUnknownReplacements(ctx, tfsdkState, tfsdkPlan, diags, ResourceSchema(ctx), info)
	for attrName, replacer := range attributePlanModifiers {
		unknownReplacements.AddReplacement(attrName, replacer)
	}
	unknownReplacements.AddKeepUnknownAlways("connection_strings", "state_name", "mongo_db_version") // Volatile attributes, should not be copied from state)
	unknownReplacements.AddKeepUnknownOnChanges(attributeRootChangeMapping)
	if computedUsed {
		unknownReplacements.AddKeepUnknownAlways("instance_size")
	}
	if diskUsed {
		unknownReplacements.AddKeepUnknownAlways("disk_size_gb")
	}
	unknownReplacements.AddKeepUnknownsExtraCall(replicationSpecsKeepUnknownWhenChanged)
	unknownReplacements.ApplyReplacements(ctx, diags)
}
