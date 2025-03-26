package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
)

var attributePlanModifiers = map[string]customplanmodifier.UnknownReplacementCall[PlanModifyResourceInfo]{
	"mongo_db_version": mongoDBVersionReplaceUnknown,
	"read_only_specs":  readOnlyReplaceUnknown,
	"analytics_specs":  analyticsSpecsReplaceUnknown,
	// TODO: Add the other computed attributes
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
	if stateParsed == nil || stateParsed.NodeCount.ValueInt64() == 0 {
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
	if req.Changes.AttributeChanged("disk_size_gb") {
		newReadOnly.DiskSizeGb = types.Float64Unknown()
	}
	return conversion.AsObjectValue(ctx, newReadOnly, SpecsObjType.AttrTypes)
}

func analyticsSpecsReplaceUnknown(ctx context.Context, state customplanmodifier.ParsedAttrValue, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) attr.Value {
	if req.Info.isShardingConfigUpgrade {
		return req.Unknown
	}
	stateParsed := conversion.TFModelObject[TFSpecsModel](ctx, state.AsObject())
	// don't get analytics_specs from state if node_count is 0 to avoid possible ANALYTICS_INSTANCE_SIZE_MUST_MATCH errors
	if stateParsed == nil || stateParsed.NodeCount.ValueInt64() == 0 {
		return req.Unknown
	}
	// if disk_size_gb is defined at root level we cannot use analytics_specs.disk_size_gb from state as it can be outdated
	if req.Changes.AttributeChanged("disk_size_gb") {
		stateParsed.DiskSizeGb = types.Float64Unknown()
	}
	return conversion.AsObjectValue(ctx, stateParsed, SpecsObjType.AttrTypes)
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
	unknownReplacements := customplanmodifier.NewUnknownReplacements(ctx, req, resp, ResourceSchema(ctx), info)
	for attrName, replacer := range attributePlanModifiers {
		unknownReplacements.AddReplacement(attrName, replacer)
	}
	unknownReplacements.ApplyReplacements(ctx, diags)
}
