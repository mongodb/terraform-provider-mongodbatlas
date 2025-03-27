package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
)

var attributePlanModifiers = map[string]customplanmodifier.UnknownReplacementCall[PlanModifyResourceInfo]{
	"mongo_db_version":       mongoDBVersionReplaceUnknown,
	"read_only_specs":        readOnlyReplaceUnknown,
	"analytics_specs":        analyticsAndElectableSpecsReplaceUnknown,
	"electable_specs":        analyticsAndElectableSpecsReplaceUnknown,
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

func parentRegionConfigs(ctx context.Context, path path.Path, differ *customplanmodifier.PlanModifyDiffer, diags *diag.Diagnostics) []TFRegionConfigsModel {
	regionConfigsPath := conversion.ParentPathNoIndex(path, "region_configs", diags)
	if diags.HasError() {
		return nil
	}
	regionConfigs := customplanmodifier.ReadPlanStructValues[TFRegionConfigsModel](ctx, differ, regionConfigsPath, diags)
	if conversion.DiagsNonEmpty(diags) {
		return nil
	}
	return regionConfigs
}

func readOnlyReplaceUnknown(ctx context.Context, state customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) attr.Value {
	if req.Info.isShardingConfigUpgrade {
		return req.Unknown
	}
	stateParsed := conversion.TFModelObject[TFSpecsModel](ctx, state.AsObject())
	if stateParsed == nil {
		return req.Unknown
	}
	electablePath := req.Path.ParentPath().AtName("electable_specs")
	electable := customplanmodifier.ReadPlanStructValue[TFSpecsModel](ctx, req.Differ, electablePath)
	if electable == nil {
		electable = customplanmodifier.ReadStateStructValue[TFSpecsModel](ctx, req.Differ, electablePath)
	}
	if electable == nil {
		regionConfigs := parentRegionConfigs(ctx, req.Path, req.Differ, req.Diags)
		if conversion.DiagsNonEmpty(req.Diags) {
			return req.Unknown
		}
		// ensures values are taken from a defined electable spec if not present in current region config
		electable = findDefinedElectableSpecInReplicationSpec(ctx, regionConfigs)
	}
	var newReadOnly *TFSpecsModel
	if electable == nil {
		// using values directly from state if no electable specs are present in plan
		newReadOnly = stateParsed
	} else {
		// node_count is from state, all others are from electable_specs plan
		newReadOnly = &TFSpecsModel{
			NodeCount:     stateParsed.NodeCount,
			InstanceSize:  electable.InstanceSize,
			DiskSizeGb:    electable.DiskSizeGb,
			EbsVolumeType: electable.EbsVolumeType,
			DiskIops:      electable.DiskIops,
		}
	}
	return conversion.AsObjectValue(ctx, newReadOnly, SpecsObjType.AttrTypes)
}

func analyticsAndElectableSpecsReplaceUnknown(ctx context.Context, state customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) attr.Value {
	if req.Info.isShardingConfigUpgrade {
		return req.Unknown
	}
	stateParsed := conversion.TFModelObject[TFSpecsModel](ctx, state.AsObject())
	// don't get analytics_specs from state if node_count is 0 to avoid possible ANALYTICS_INSTANCE_SIZE_MUST_MATCH and INSTANCE_SIZE_MUST_MATCH errors
	if stateParsed == nil || stateParsed.NodeCount.ValueInt64() == 0 {
		return req.Unknown
	}
	// if disk_size_gb is defined at root level we cannot use analytics_specs.disk_size_gb from state as it can be outdated
	if req.Changes.AttributeChanged("disk_size_gb") {
		stateParsed.DiskSizeGb = types.Float64Unknown()
	}
	return conversion.AsObjectValue(ctx, stateParsed, SpecsObjType.AttrTypes)
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
