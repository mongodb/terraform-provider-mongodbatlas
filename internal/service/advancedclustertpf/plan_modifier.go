package advancedclustertpf

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

var (
	// Change mappings uses `attribute_name`, it doesn't care about the nested level.
	attributeRootChangeMapping = map[string][]string{
		"disk_size_gb":           {}, // disk_size_gb can be change at any level/spec
		"replication_specs":      {},
		"mongo_db_major_version": {"mongo_db_version"},
		"tls_cipher_config_mode": {"custom_openssl_cipher_config_tls12"},
		"cluster_type":           {"config_server_management_mode", "config_server_type"}, // computed values of config server change when REPLICA_SET changes to SHARDED
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
	keepUnknownsCalls = schemafunc.KeepUnknownFuncOr(keepUnkownFuncWithNodeCount, keepUnkownFuncWithNonEmptyAutoScaling)
)

func keepUnkownFuncWithNodeCount(name string, replacement attr.Value) bool {
	return name == "node_count" && !replacement.Equal(types.Int64Value(0))
}

func keepUnkownFuncWithNonEmptyAutoScaling(name string, replacement attr.Value) bool {
	autoScalingBoolValues := []string{"compute_enabled", "disk_gb_enabled", "compute_scale_down_enabled"}
	autoScalingStringValues := []string{"compute_min_instance_size", "compute_max_instance_size"}
	boolValues := slices.Contains(autoScalingBoolValues, name) && replacement.Equal(types.BoolValue(true))
	stringValues := slices.Contains(autoScalingStringValues, name) && replacement.(types.String).ValueString() != ""
	return boolValues || stringValues
}

// useStateForUnknowns should be called only in Update, because of findClusterDiff
func useStateForUnknowns(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	shardingConfigUpgrade := isShardingConfigUpgrade(ctx, state, plan, diags)
	if diags.HasError() {
		return
	}
	// Don't adjust region_configs upgrades if it's a sharding config upgrade because it will be done only in the first shard, because state only has the first shard with num_shards > 1.
	// This avoid errors like AUTO_SCALINGS_MUST_BE_IN_EVERY_REGION_CONFIG.
	if !shardingConfigUpgrade {
		AdjustRegionConfigsChildren(ctx, diags, state, plan)
	}
	diff := findClusterDiff(ctx, state, plan, diags)
	if diags.HasError() || diff.isAnyUpgrade() { // Don't do anything in upgrades
		return
	}
	attributeChanges := schemafunc.NewAttributeChanges(ctx, state, plan)
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	keepUnknown = append(keepUnknown, attributeChanges.KeepUnknown(attributeRootChangeMapping)...)
	keepUnknown = append(keepUnknown, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown, nil)
	/* pending revision if logic can be reincorporated safely:
	if slices.Contains(keepUnknown, "replication_specs") {
		useStateForUnknownsReplicationSpecs(ctx, diags, state, plan, &attributeChanges)
	}
	*/
}

func UseStateForUnknownsReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, attrChanges *schemafunc.AttributeChanges) {
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return
	}
	planWithUnknowns := []TFReplicationSpecsModel{}
	keepUnknownsUnchangedSpec := determineKeepUnknownsUnchangedReplicationSpecs(ctx, diags, state, plan, attrChanges)
	keepUnknownsUnchangedSpec = append(keepUnknownsUnchangedSpec, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
	if diags.HasError() {
		return
	}
	for i := range planRepSpecsTF {
		if i < len(stateRepSpecsTF) {
			keepUnknowns := keepUnknownsUnchangedSpec
			if attrChanges.ListIndexChanged("replication_specs", i) {
				keepUnknowns = determineKeepUnknownsChangedReplicationSpec(keepUnknownsUnchangedSpec, attrChanges, fmt.Sprintf("replication_specs[%d]", i))
			}
			schemafunc.CopyUnknowns(ctx, &stateRepSpecsTF[i], &planRepSpecsTF[i], keepUnknowns, keepUnknownsCalls)
		}
		planWithUnknowns = append(planWithUnknowns, planRepSpecsTF[i])
	}
	listType, diagsLocal := types.ListValueFrom(ctx, ReplicationSpecsObjType, planWithUnknowns)
	diags.Append(diagsLocal...)
	if diags.HasError() {
		return
	}
	plan.ReplicationSpecs = listType
}

// AdjustRegionConfigsChildren modifies the planned values of region configs based on the current state.
// This ensures proper handling of removing auto scaling and specs attributes by preserving state values.
func AdjustRegionConfigsChildren(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return
	}
	for i := range minLen(planRepSpecsTF, stateRepSpecsTF) {
		stateRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, stateRepSpecsTF[i].RegionConfigs)
		planRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, planRepSpecsTF[i].RegionConfigs)
		planElectableSpecInReplicationSpec := findDefinedElectableSpecInReplicationSpec(ctx, planRegionConfigsTF)
		if diags.HasError() {
			return
		}
		for j := range minLen(planRegionConfigsTF, stateRegionConfigsTF) {
			stateReadOnlySpecs := TFModelObject[TFSpecsModel](ctx, stateRegionConfigsTF[j].ReadOnlySpecs)
			planReadOnlySpecs := TFModelObject[TFSpecsModel](ctx, planRegionConfigsTF[j].ReadOnlySpecs)
			planElectableSpecs := TFModelObject[TFSpecsModel](ctx, planRegionConfigsTF[j].ElectableSpecs)
			if stateReadOnlySpecs != nil { // read_only_specs is present in state
				newPlanReadOnlySpecs := planReadOnlySpecs
				if newPlanReadOnlySpecs == nil {
					newPlanReadOnlySpecs = new(TFSpecsModel) // start with null attributes if not present plan
				}
				baseReadOnlySpecs := stateReadOnlySpecs        // using values directly from state if no electable specs are present in plan
				if planElectableSpecInReplicationSpec != nil { // ensures values are taken from a defined electable spec if not present in current region config
					baseReadOnlySpecs = planElectableSpecInReplicationSpec
				}
				if planElectableSpecs != nil { // we favor plan electable spec defined in same region config over one defined in replication spec to be more future proof
					baseReadOnlySpecs = planElectableSpecs
				}
				copyAttrIfDestNotKnown(&baseReadOnlySpecs.DiskSizeGb, &newPlanReadOnlySpecs.DiskSizeGb)
				copyAttrIfDestNotKnown(&baseReadOnlySpecs.EbsVolumeType, &newPlanReadOnlySpecs.EbsVolumeType)
				copyAttrIfDestNotKnown(&baseReadOnlySpecs.InstanceSize, &newPlanReadOnlySpecs.InstanceSize)
				copyAttrIfDestNotKnown(&baseReadOnlySpecs.DiskIops, &newPlanReadOnlySpecs.DiskIops)
				// unknown node_count is always taken from state as it not dependent on electable_specs changes
				copyAttrIfDestNotKnown(&stateReadOnlySpecs.NodeCount, &newPlanReadOnlySpecs.NodeCount)
				objType, diagsLocal := types.ObjectValueFrom(ctx, SpecsObjType.AttrTypes, newPlanReadOnlySpecs)
				diags.Append(diagsLocal...)
				if diags.HasError() {
					return
				}
				planRegionConfigsTF[j].ReadOnlySpecs = objType
			}

			stateAnalyticsSpecs := TFModelObject[TFSpecsModel](ctx, stateRegionConfigsTF[j].AnalyticsSpecs)
			planAnalyticsSpecs := TFModelObject[TFSpecsModel](ctx, planRegionConfigsTF[j].AnalyticsSpecs)
			// don't get analytics_specs from state if node_count is 0 to avoid possible ANALYTICS_INSTANCE_SIZE_MUST_MATCH errors
			if planAnalyticsSpecs == nil && stateAnalyticsSpecs != nil && stateAnalyticsSpecs.NodeCount.ValueInt64() > 0 {
				newPlanAnalyticsSpecs := TFModelObject[TFSpecsModel](ctx, stateRegionConfigsTF[j].AnalyticsSpecs)
				// if disk_size_gb is defined at root level we cannot use analytics_specs.disk_size_gb from state as it can be outdated
				// read_only_specs implicitly covers this as it uses value from electable_specs which is unknown if not defined.
				if plan.DiskSizeGB.ValueFloat64() > 0 { // has known value in config
					newPlanAnalyticsSpecs.DiskSizeGb = types.Float64Unknown()
				}
				objType, diagsLocal := types.ObjectValueFrom(ctx, SpecsObjType.AttrTypes, newPlanAnalyticsSpecs)
				diags.Append(diagsLocal...)
				if diags.HasError() {
					return
				}
				planRegionConfigsTF[j].AnalyticsSpecs = objType
			}

			// don't use auto_scaling or analytics_auto_scaling from state if it's not enabled as it doesn't need to be present in Update request payload
			stateAutoScaling := TFModelObject[TFAutoScalingModel](ctx, stateRegionConfigsTF[j].AutoScaling)
			planAutoScaling := TFModelObject[TFAutoScalingModel](ctx, planRegionConfigsTF[j].AutoScaling)
			if planAutoScaling == nil && stateAutoScaling != nil && (stateAutoScaling.ComputeEnabled.ValueBool() || stateAutoScaling.DiskGBEnabled.ValueBool()) {
				planRegionConfigsTF[j].AutoScaling = stateRegionConfigsTF[j].AutoScaling
			}
			stateAnalyticsAutoScaling := TFModelObject[TFAutoScalingModel](ctx, stateRegionConfigsTF[j].AnalyticsAutoScaling)
			planAnalyticsAutoScaling := TFModelObject[TFAutoScalingModel](ctx, planRegionConfigsTF[j].AnalyticsAutoScaling)
			if planAnalyticsAutoScaling == nil && stateAnalyticsAutoScaling != nil && (stateAnalyticsAutoScaling.ComputeEnabled.ValueBool() || stateAnalyticsAutoScaling.DiskGBEnabled.ValueBool()) {
				planRegionConfigsTF[j].AnalyticsAutoScaling = stateRegionConfigsTF[j].AnalyticsAutoScaling
			}
		}
		listRegionConfigs, diagsLocal := types.ListValueFrom(ctx, RegionConfigsObjType, planRegionConfigsTF)
		diags.Append(diagsLocal...)
		if diags.HasError() {
			return
		}
		planRepSpecsTF[i].RegionConfigs = listRegionConfigs
	}
	listRepSpecs, diagsLocal := types.ListValueFrom(ctx, ReplicationSpecsObjType, planRepSpecsTF)
	diags.Append(diagsLocal...)
	if diags.HasError() {
		return
	}
	plan.ReplicationSpecs = listRepSpecs
}

func findDefinedElectableSpecInReplicationSpec(ctx context.Context, regionConfigs []TFRegionConfigsModel) *TFSpecsModel {
	for i := range regionConfigs {
		electableSpecs := TFModelObject[TFSpecsModel](ctx, regionConfigs[i].ElectableSpecs)
		if electableSpecs != nil {
			return electableSpecs
		}
	}
	return nil
}

// determineKeepUnknownsChangedReplicationSpec: These fields must be kept unknown in the replication_specs[index_of_changes]
func determineKeepUnknownsChangedReplicationSpec(keepUnknownsAlways []string, attributeChanges *schemafunc.AttributeChanges, parentPath string) []string {
	var keepUnknowns = slices.Clone(keepUnknownsAlways)
	if attributeChanges.NestedListLenChanges(parentPath + ".region_configs") {
		keepUnknowns = append(keepUnknowns, "container_id")
	}
	return append(keepUnknowns, attributeChanges.KeepUnknown(attributeReplicationSpecChangeMapping)...)
}

func determineKeepUnknownsUnchangedReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, attributeChanges *schemafunc.AttributeChanges) []string {
	keepUnknowns := []string{}
	// Could be set to "" if we are using an ISS cluster
	if usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags) { // When using new sharding config, the legacy id must never be copied
		keepUnknowns = append(keepUnknowns, "id")
	}
	// for isShardingConfigUpgrade, it will be empty in the plan, so we need to keep it unknown
	// for listLenChanges, it might be an insertion in the middle of replication spec leading to wrong value from state copied
	if isShardingConfigUpgrade(ctx, state, plan, diags) || attributeChanges.ListLenChanges("replication_specs") {
		keepUnknowns = append(keepUnknowns, "external_id")
	}
	return keepUnknowns
}

func determineKeepUnknownsAutoScaling(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) []string {
	var keepUnknown []string
	computedUsed, diskUsed := autoScalingUsed(ctx, diags, state, plan)
	if computedUsed {
		keepUnknown = append(keepUnknown, "instance_size")
		keepUnknown = append(keepUnknown, attributeReplicationSpecChangeMapping["instance_size"]...)
	}
	if diskUsed {
		keepUnknown = append(keepUnknown, "disk_size_gb")
		keepUnknown = append(keepUnknown, attributeReplicationSpecChangeMapping["disk_size_gb"]...)
	}
	return keepUnknown
}

// autoScalingUsed checks is auto-scaling was enabled (state) or will be enabled (plan).
func autoScalingUsed(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) (computedUsed, diskUsed bool) {
	for _, model := range []*TFModel{state, plan} {
		repSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, model.ReplicationSpecs)
		for i := range repSpecsTF {
			regiongConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, repSpecsTF[i].RegionConfigs)
			for j := range regiongConfigsTF {
				for _, autoScalingTF := range []types.Object{regiongConfigsTF[j].AutoScaling, regiongConfigsTF[j].AnalyticsAutoScaling} {
					autoscaling := TFModelObject[TFAutoScalingModel](ctx, autoScalingTF)
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

func TFModelList[T any](ctx context.Context, diags *diag.Diagnostics, input types.List) []T {
	elements := make([]T, len(input.Elements()))
	diags.Append(input.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return nil
	}
	return elements
}

// TFModelObject returns nil if the Terraform object is null or unknown, or casting to T is not valid. However object attributes can be null or unknown.
func TFModelObject[T any](ctx context.Context, input types.Object) *T {
	item := new(T)
	if diags := input.As(ctx, item, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil
	}
	return item
}

func copyAttrIfDestNotKnown[T attr.Value](src, dest *T) {
	if !isKnown(*dest) {
		*dest = *src
	}
}

// isKnown returns true if the attribute is known (not null or unknown). Note that !isKnown is not the same as IsUnknown because null is !isKnown but not IsUnknown.
func isKnown(attribute attr.Value) bool {
	return !attribute.IsNull() && !attribute.IsUnknown()
}

func minLen[T any](a, b []T) int {
	la, lb := len(a), len(b)
	if la < lb {
		return la
	}
	return lb
}
