package advancedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

// useStateForUnknowns should be called only in Update, because of findClusterDiff
func useStateForUnknowns(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	DebugPrintSpecs(ctx, "BEFORE state", state)
	DebugPrintSpecs(ctx, "BEFORE plan", plan)

	AdjustRegionConfigsChildren(ctx, diags, state, plan)

	DebugPrintSpecs(ctx, "MIDDLE state", state)
	DebugPrintSpecs(ctx, "MIDDLE plan", plan)

	diff := findClusterDiff(ctx, state, plan, diags)
	if diags.HasError() || diff.isAnyUpgrade() { // Don't do anything in upgrades
		return
	}
	attributeChanges := schemafunc.NewAttributeChanges(ctx, state, plan)
	keepUnknown := []string{"connection_strings", "state_name", "mongo_db_version"} // Volatile attributes, should not be copied from state
	// When a key attribute changes, dependent attributes may also change and must remain unknown. Attribute matching is by name, independent of nesting level.
	keepUnknown = append(keepUnknown, attributeChanges.KeepUnknown(map[string][]string{
		"replication_specs":      {},
		"tls_cipher_config_mode": {"custom_openssl_cipher_config_tls12", "custom_openssl_cipher_config_tls13"},
		// When switching between custom_openssl_cipher_config_tls12 and custom_openssl_cipher_config_tls13, both must remain unknown to avoid sending a PATCH with both values included.
		"custom_openssl_cipher_config_tls12": {"custom_openssl_cipher_config_tls13"},
		"custom_openssl_cipher_config_tls13": {"custom_openssl_cipher_config_tls12"},
		// Computed values of config server change when REPLICA_SET changes to SHARDED.
		"cluster_type": {"config_server_management_mode", "config_server_type"},
	})...)
	keepUnknown = append(keepUnknown, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
	keepUnknownFunc := determineKeepUnknownsUseEffectiveFields(state, plan)
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown, keepUnknownFunc)

	DebugPrintSpecs(ctx, "AFTER state", state)
	DebugPrintSpecs(ctx, "AFTER plan", plan)
}

// AdjustRegionConfigsChildren modifies the planned values of region configs based on the current state.
// This ensures proper handling of removing auto scaling and specs attributes by preserving state values.
func AdjustRegionConfigsChildren(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	// Not safe if use_effective_fields is changes
	if state.UseEffectiveFields.ValueBool() != plan.UseEffectiveFields.ValueBool() {
		return
	}
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
			stateElectableSpecs := TFModelObject[TFSpecsModel](ctx, stateRegionConfigsTF[j].ElectableSpecs)
			planElectableSpecs := TFModelObject[TFSpecsModel](ctx, planRegionConfigsTF[j].ElectableSpecs)
			if planElectableSpecs == nil && stateElectableSpecs != nil && stateElectableSpecs.NodeCount.ValueInt64() > 0 {
				planRegionConfigsTF[j].ElectableSpecs = stateRegionConfigsTF[j].ElectableSpecs
				planElectableSpecs = stateElectableSpecs
			}
			stateReadOnlySpecs := TFModelObject[TFSpecsModel](ctx, stateRegionConfigsTF[j].ReadOnlySpecs)
			planReadOnlySpecs := TFModelObject[TFSpecsModel](ctx, planRegionConfigsTF[j].ReadOnlySpecs)
			if stateReadOnlySpecs != nil { // read_only_specs is present in state
				// logic below ensures that if read only specs is present in state but not in the plan, plan will be populated so that read only spec configuration is not removed on update operations
				newPlanReadOnlySpecs := planReadOnlySpecs
				if newPlanReadOnlySpecs == nil {
					newPlanReadOnlySpecs = new(TFSpecsModel) // start with null attributes if not present plan
				}
				baseReadOnlySpecs := stateReadOnlySpecs        // using values directly from state if no electable specs are present in plan
				if planElectableSpecInReplicationSpec != nil { // ensures values are taken from a defined electable spec if not present in current region config
					baseReadOnlySpecs = planElectableSpecInReplicationSpec
				}
				if planElectableSpecs != nil {
					// we favor plan electable spec defined in same region config over one defined in replication spec
					// with current API this is redudant but is more future proof in case scaling between regions becomes independent in the future
					baseReadOnlySpecs = planElectableSpecs
				}
				copyAttrIfDestNotKnown(&baseReadOnlySpecs.DiskSizeGb, &newPlanReadOnlySpecs.DiskSizeGb)
				copyAttrIfDestNotKnown(&baseReadOnlySpecs.EbsVolumeType, &newPlanReadOnlySpecs.EbsVolumeType)
				copyAttrIfDestNotKnown(&baseReadOnlySpecs.InstanceSize, &newPlanReadOnlySpecs.InstanceSize)
				copyAttrIfDestNotKnown(&baseReadOnlySpecs.DiskIops, &newPlanReadOnlySpecs.DiskIops)
				// unknown node_count is always taken from state as it not dependent on electable_specs changes
				copyAttrIfDestNotKnown(&stateReadOnlySpecs.NodeCount, &newPlanReadOnlySpecs.NodeCount)
				objType, diagsLocal := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, newPlanReadOnlySpecs)
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
				objType, diagsLocal := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, newPlanAnalyticsSpecs)
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
		listRegionConfigs, diagsLocal := types.ListValueFrom(ctx, regionConfigsObjType, planRegionConfigsTF)
		diags.Append(diagsLocal...)
		if diags.HasError() {
			return
		}
		planRepSpecsTF[i].RegionConfigs = listRegionConfigs
	}
	listRepSpecs, diagsLocal := types.ListValueFrom(ctx, replicationSpecsObjType, planRepSpecsTF)
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

func determineKeepUnknownsAutoScaling(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) []string {
	var keepUnknown []string
	computedUsed, diskUsed := autoScalingUsed(ctx, diags, state, plan)
	if computedUsed {
		keepUnknown = append(keepUnknown, "instance_size", "disk_iops") // disk_iops can change based on instance_size changes
	}
	if diskUsed {
		keepUnknown = append(keepUnknown, "disk_size_gb")
	}
	return keepUnknown
}

// determineKeepUnknownsUseEffectiveFields returns a function that keeps spec fields unknown when use_effective_fields changes.
func determineKeepUnknownsUseEffectiveFields(state, plan *TFModel) func(string, attr.Value) bool {
	// If use_effective_fields is changing, we need to keep spec fields unknown
	if state.UseEffectiveFields.ValueBool() == plan.UseEffectiveFields.ValueBool() {
		return nil // No change, don't filter anything
	}

	// List of spec object names and their field names
	specFields := []string{
		// Spec object names
		"electable_specs", "read_only_specs", "analytics_specs",
		// Field names within specs
		"node_count", "instance_size", "disk_size_gb", "disk_iops", "ebs_volume_type",
	}

	return func(name string, value attr.Value) bool {
		// Keep unknown if the field is spec-related
		for _, field := range specFields {
			if name == field {
				return true
			}
		}
		return false
	}
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

// DebugPrintSpecs prints spec information for the first replication spec's first region config for debugging.
func DebugPrintSpecs(ctx context.Context, label string, model *TFModel) {
	if model == nil {
		fmt.Printf("[DEBUG %s] model is nil\n", label)
		return
	}

	repSpecs := TFModelList[TFReplicationSpecsModel](ctx, &diag.Diagnostics{}, model.ReplicationSpecs)
	if len(repSpecs) == 0 {
		fmt.Printf("[DEBUG %s] no replication specs\n", label)
		return
	}

	regionConfigs := TFModelList[TFRegionConfigsModel](ctx, &diag.Diagnostics{}, repSpecs[0].RegionConfigs)
	if len(regionConfigs) == 0 {
		fmt.Printf("[DEBUG %s] no region configs\n", label)
		return
	}

	fmt.Printf("[DEBUG 12345 %s] Specs %s\n", label, attrValueString(model.UseEffectiveFields))

	printSpecDetails := func(specName string, specObj types.Object) {
		if specObj.IsNull() {
			fmt.Printf("  %s: null\n", specName)
			return
		}
		if specObj.IsUnknown() {
			fmt.Printf("  %s: unknown\n", specName)
			return
		}

		spec := TFModelObject[TFSpecsModel](ctx, specObj)
		if spec == nil {
			fmt.Printf("  %s: <failed to convert>\n", specName)
			return
		}

		fmt.Printf("  %s: node_count=%s, instance_size=%s, disk_size_gb=%s, disk_iops=%s, ebs_volume_type=%s\n",
			specName,
			attrValueString(spec.NodeCount),
			attrValueString(spec.InstanceSize),
			attrValueString(spec.DiskSizeGb),
			attrValueString(spec.DiskIops),
			attrValueString(spec.EbsVolumeType),
		)
	}

	printSpecDetails("electable_specs", regionConfigs[0].ElectableSpecs)
	printSpecDetails("read_only_specs", regionConfigs[0].ReadOnlySpecs)
	printSpecDetails("analytics_specs", regionConfigs[0].AnalyticsSpecs)
}

// attrValueString returns a string representation of an attribute value showing if it's known, null, or unknown.
func attrValueString(value attr.Value) string {
	if value.IsNull() {
		return "null"
	}
	if value.IsUnknown() {
		return "unknown"
	}
	return value.String()
}
