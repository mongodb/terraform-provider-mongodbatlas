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
	attributeRootChangeMapping = map[string][]string{
		"disk_size_gb":           {}, // disk_size_gb can be change at any level/spec
		"replication_specs":      {},
		"mongo_db_major_version": {"mongo_db_version"},
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
	AdjustRegionConfigsChildren(ctx, diags, state, plan)
	diff := findClusterDiff(ctx, state, plan, diags)
	if diags.HasError() || diff.isAnyUpgrade() { // Don't do anything in upgrades
		return
	}
	attributeChanges := schemafunc.NewAttributeChanges(ctx, state, plan)
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	keepUnknown = append(keepUnknown, attributeChanges.KeepUnknown(attributeRootChangeMapping)...)
	// pending revision if logic can be reincorporated safely: keepUnknown = append(keepUnknown, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
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
	// TODO: pending if autoscaling logic is needed or provided that lifecycle ignored_changes is used, that is enough
	// keepUnknownsUnchangedSpec = append(keepUnknownsUnchangedSpec, DetermineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
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

func AdjustRegionConfigsChildren(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return
	}
	for i := range minLen(planRepSpecsTF, stateRepSpecsTF) {
		stateRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, stateRepSpecsTF[i].RegionConfigs)
		planRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, planRepSpecsTF[i].RegionConfigs)
		if diags.HasError() {
			return
		}
		for j := range minLen(planRegionConfigsTF, stateRegionConfigsTF) {
			stateReadOnlySpecs := TFModelObject[TFSpecsModel](ctx, stateRegionConfigsTF[j].ReadOnlySpecs)
			planReadOnlySpecs := TFModelObject[TFSpecsModel](ctx, planRegionConfigsTF[j].ReadOnlySpecs)
			planElectableSpecs := TFModelObject[TFSpecsModel](ctx, planRegionConfigsTF[j].ElectableSpecs)
			if stateReadOnlySpecs != nil && planElectableSpecs != nil { // read_only_specs is present in state and electable_specs in the plan
				newPlanReadOnlySpecs := planReadOnlySpecs
				if newPlanReadOnlySpecs == nil {
					newPlanReadOnlySpecs = new(TFSpecsModel) // start with null attributes if not present plan
				}
				// unknown node_count is got from state, all other unknowns are got from electable_specs plan
				copyAttrIfDestNotKnown(&planElectableSpecs.DiskSizeGb, &newPlanReadOnlySpecs.DiskSizeGb)
				copyAttrIfDestNotKnown(&planElectableSpecs.EbsVolumeType, &newPlanReadOnlySpecs.EbsVolumeType)
				copyAttrIfDestNotKnown(&planElectableSpecs.InstanceSize, &newPlanReadOnlySpecs.InstanceSize)
				copyAttrIfDestNotKnown(&planElectableSpecs.DiskIops, &newPlanReadOnlySpecs.DiskIops)
				copyAttrIfDestNotKnown(&stateReadOnlySpecs.NodeCount, &newPlanReadOnlySpecs.NodeCount)
				objType, diagsLocal := types.ObjectValueFrom(ctx, SpecsObjType.AttrTypes, newPlanReadOnlySpecs)
				diags.Append(diagsLocal...)
				if diags.HasError() {
					return
				}
				planRegionConfigsTF[j].ReadOnlySpecs = objType
			}
			if planRegionConfigsTF[j].AnalyticsSpecs.IsUnknown() && !stateRegionConfigsTF[j].AnalyticsSpecs.IsNull() {
				planRegionConfigsTF[j].AnalyticsSpecs = stateRegionConfigsTF[j].AnalyticsSpecs
			}
			if planRegionConfigsTF[j].AutoScaling.IsUnknown() && !stateRegionConfigsTF[j].AutoScaling.IsNull() {
				planRegionConfigsTF[j].AutoScaling = stateRegionConfigsTF[j].AutoScaling
			}
			if planRegionConfigsTF[j].AnalyticsAutoScaling.IsUnknown() && !stateRegionConfigsTF[j].AnalyticsAutoScaling.IsNull() {
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

func DetermineKeepUnknownsAutoScaling(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) []string {
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
					if autoScalingTF.IsNull() || autoScalingTF.IsUnknown() {
						continue
					}
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
