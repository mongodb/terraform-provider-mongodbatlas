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

var (
	attributePlanModifiers = map[string]customplanmodifier.UnknownReplacementCall[PlanModifyResourceInfo]{
		"read_only_specs":        readOnlyReplaceUnknown,
		"analytics_specs":        analyticsAndElectableSpecsReplaceUnknown,
		"electable_specs":        analyticsAndElectableSpecsReplaceUnknown,
		"auto_scaling":           autoScalingReplaceUnknown,
		"analytics_auto_scaling": autoScalingReplaceUnknown,
		// TODO: Add the other computed attributes
	}
	// Change mappings uses `attribute_name`, it doesn't care about the nested level.
	// However, it doesn't stop calling `replication_specs.**.attribute_name`.
	attributeRootChangeMapping = map[string][]string{
		"disk_size_gb":           {}, // disk_size_gb can be change at any level/spec
		"replication_specs":      {},
		"tls_cipher_config_mode": {"custom_openssl_cipher_config_tls12"},
		"cluster_type":           {"config_server_management_mode", "config_server_type"}, // computed values of config server change when REPLICA_SET changes to SHARDED
		"expiration_date":        {"version"},                                             // pinned_fcv
	}
	attributeReplicationSpecChangeMapping = map[string][]string{
		// All these fields can exist in specs that are computed, therefore, it is not safe to use them when they have changed.
		"disk_iops":       {},
		"ebs_volume_type": {},
		"disk_size_gb":    {},                  // disk_size_gb can be change at any level/spec
		"instance_size":   {"disk_iops"},       // disk_iops can change based on instance_size changes
		"provider_name":   {"ebs_volume_type"}, // AWS --> AZURE will change ebs_volume_type
		"region_name":     {"container_id"},    // container_id changes based on region_name changes
		"zone_name":       {"zone_id"},         // zone_id copy from state is not safe when
	}
)

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
	isUsingNewShardingConfig := usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags)
	if diags.HasError() {
		return
	}
	info := PlanModifyResourceInfo{
		AutoScalingComputedUsed: computedUsed,
		AutoScalingDiskUsed:     diskUsed,
		IsShardingConfigUpgrade: shardingConfigUpgrade,
		UsingNewShardingConfig:  isUsingNewShardingConfig,
	}
	unknownReplacements := customplanmodifier.NewUnknownReplacements(ctx, tfsdkState, tfsdkPlan, diags, ResourceSchema(ctx), info)
	for attrName, replacer := range attributePlanModifiers {
		unknownReplacements.AddReplacement(attrName, replacer)
	}
	unknownReplacements.AddKeepUnknownAlways("connection_strings", "state_name", "mongo_db_version") // Volatile attributes, should not be copied from state)
	unknownReplacements.AddKeepUnknownOnChanges(attributeRootChangeMapping)
	if computedUsed {
		unknownReplacements.AddKeepUnknownAlways("instance_size", "disk_iops")
	}
	if diskUsed {
		unknownReplacements.AddKeepUnknownAlways("disk_size_gb")
	}
	unknownReplacements.AddKeepUnknownsExtraCall(replicationSpecsKeepUnknownWhenChanged)
	unknownReplacements.ApplyReplacements(ctx, diags)
}

func autoScalingReplaceUnknown(ctx context.Context, state attr.Value, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) attr.Value {
	// don't use auto_scaling or analytics_auto_scaling from state if it's not enabled as it doesn't need to be present in Update request payload
	if req.Info.AutoScalingComputedUsed || req.Info.AutoScalingDiskUsed {
		return state.(types.Object)
	}
	return req.Unknown
}

type PlanModifyResourceInfo struct {
	AutoScalingComputedUsed bool
	AutoScalingDiskUsed     bool
	IsShardingConfigUpgrade bool
	UsingNewShardingConfig  bool
}

func parentRegionConfigs(ctx context.Context, path path.Path, differ *customplanmodifier.PlanModifyDiffer, diags *diag.Diagnostics) []TFRegionConfigsModel {
	regionConfigsPath := conversion.AncestorPathNoIndex(path, "region_configs", diags)
	if diags.HasError() {
		return nil
	}
	regionConfigs := customplanmodifier.ReadPlanStructValues[TFRegionConfigsModel](ctx, differ, regionConfigsPath, diags)
	if diags.HasError() {
		return nil
	}
	return regionConfigs
}

func readOnlyReplaceUnknown(ctx context.Context, state attr.Value, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) attr.Value {
	if req.Info.IsShardingConfigUpgrade {
		return req.Unknown
	}
	stateParsed := conversion.TFModelObject[TFSpecsModel](ctx, state.(types.Object))
	if stateParsed == nil {
		return req.Unknown
	}
	electablePath := req.Path.ParentPath().AtName("electable_specs")
	electable := customplanmodifier.ReadPlanStructValue[TFSpecsModel](ctx, req.Differ, electablePath)
	if electable == nil {
		electableState := customplanmodifier.ReadStateStructValue[TFSpecsModel](ctx, req.Differ, electablePath)
		if electableState.NodeCount.ValueInt64() > 0 {
			electable = electableState
		}
	}
	if electable == nil {
		regionConfigs := parentRegionConfigs(ctx, req.Path, req.Differ, req.Diags)
		if req.Diags.HasError() {
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
	return ensureSpecRespectChanges(ctx, newReadOnly, req)
}

func analyticsAndElectableSpecsReplaceUnknown(ctx context.Context, state attr.Value, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) attr.Value {
	if req.Info.IsShardingConfigUpgrade {
		return req.Unknown
	}
	stateParsed := conversion.TFModelObject[TFSpecsModel](ctx, state.(types.Object))
	// don't get analytics_specs from state if node_count is 0 to avoid possible ANALYTICS_INSTANCE_SIZE_MUST_MATCH and INSTANCE_SIZE_MUST_MATCH errors
	if stateParsed == nil || stateParsed.NodeCount.ValueInt64() == 0 {
		return req.Unknown
	}
	return ensureSpecRespectChanges(ctx, stateParsed, req)
}

func ensureSpecRespectChanges(ctx context.Context, spec *TFSpecsModel, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) types.Object {
	// if disk_size_gb is defined at root level we cannot use (analytics|electable|read_only)_specs.disk_size_gb from state as it can be outdated
	if req.Changes.AttributeChanged("disk_size_gb") || req.Info.AutoScalingDiskUsed {
		spec.DiskSizeGb = types.Float64Unknown()
	}
	if req.Info.AutoScalingComputedUsed {
		spec.DiskIops = types.Int64Unknown()
	}
	return conversion.AsObjectValue(ctx, spec, SpecsObjType.AttrTypes)
}

func replicationSpecsKeepUnknownWhenChanged(ctx context.Context, state attr.Value, req *customplanmodifier.UnknownReplacementRequest[PlanModifyResourceInfo]) []string {
	rootPath := path.Root("replication_specs")
	if !conversion.HasAncestor(req.Path, rootPath) {
		return nil
	}
	if !req.Changes.PathChanged(rootPath) {
		return nil
	}
	keepUnknowns := []string{}
	if req.Info.UsingNewShardingConfig {
		keepUnknowns = append(keepUnknowns, "id") // When using new sharding config, the legacy id must never be copied
	}
	// for isShardingConfigUpgrade, it will be empty in the plan, so we need to keep it unknown
	// for listLenChanges, it might be an insertion in the middle of replication spec leading to wrong value from state copied
	if req.Info.IsShardingConfigUpgrade || req.Changes.ListLenChanged(rootPath) {
		keepUnknowns = append(keepUnknowns, "external_id")
	}
	replicationSpecAncestor := conversion.AncestorPathWithIndex(req.Path, "replication_specs", req.Diags)
	if req.Diags.HasError() {
		return keepUnknowns
	}
	if !req.Changes.PathChanged(replicationSpecAncestor) {
		return keepUnknowns
	}
	if req.Changes.ListLenChanged(replicationSpecAncestor.AtName("region_configs")) {
		keepUnknowns = append(keepUnknowns, "container_id")
	}
	keepUnknowns = append(keepUnknowns, req.Changes.KeepUnknown(attributeReplicationSpecChangeMapping)...)
	return keepUnknowns
}

func findDefinedElectableSpecInReplicationSpec(ctx context.Context, regionConfigs []TFRegionConfigsModel) *TFSpecsModel {
	for i := range regionConfigs {
		electableSpecs := conversion.TFModelObject[TFSpecsModel](ctx, regionConfigs[i].ElectableSpecs)
		if electableSpecs != nil {
			return electableSpecs
		}
	}
	return nil
}

func autoScalingUsed(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) (computedUsed, diskUsed bool) {
	for _, model := range []*TFModel{state, plan} {
		repSpecsTF := conversion.TFModelList[TFReplicationSpecsModel](ctx, diags, model.ReplicationSpecs)
		for i := range repSpecsTF {
			regiongConfigsTF := conversion.TFModelList[TFRegionConfigsModel](ctx, diags, repSpecsTF[i].RegionConfigs)
			for j := range regiongConfigsTF {
				for _, autoScalingTF := range []types.Object{regiongConfigsTF[j].AutoScaling, regiongConfigsTF[j].AnalyticsAutoScaling} {
					autoscaling := conversion.TFModelObject[TFAutoScalingModel](ctx, autoScalingTF)
					if autoscaling == nil {
						continue
					}
					if autoscaling.ComputeEnabled.ValueBool() {
						computedUsed = true
					}
					if autoscaling.DiskGBEnabled.ValueBool() {
						diskUsed = true
					}
				}
			}
		}
	}
	return
}
