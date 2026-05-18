package advancedcluster

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

type RegionAutoScaling struct {
	RCPrefix                string
	ComputeEnabled          bool
	DiskGBEnabled           bool
	AnalyticsComputeEnabled bool
}

var (
	// Spec fields that Atlas controls when auto-scaling is active.
	autoScalingManagedSpecFields = []string{"instance_size", "disk_size_gb", "disk_iops"}

	// Change mappings uses `attribute_name`, it doesn't care about the nested level.
	attributeRootChangeMapping = map[string][]string{
		"replication_specs":      {},
		"tls_cipher_config_mode": {"custom_openssl_cipher_config_tls12", "custom_openssl_cipher_config_tls13"},
		// When changing custom_openssl_cipher_config_tls12 -> custom_openssl_cipher_config_tls13 and vice versa, we CANNOT use the state value
		// it needs to be unknown to avoid sending a PATCH with both values included.
		"custom_openssl_cipher_config_tls12": {"custom_openssl_cipher_config_tls13"},
		"custom_openssl_cipher_config_tls13": {"custom_openssl_cipher_config_tls12"},
		"cluster_type":                       {"config_server_management_mode", "config_server_type"}, // computed values of config server change when REPLICA_SET changes to SHARDED
	}
)

// handleModifyPlan should be called only in Update, because of findClusterDiff
func handleModifyPlan(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	// Special logic for use_effective_fields changes, as normal optimization is not safe.
	if state.UseEffectiveFields.ValueBool() != plan.UseEffectiveFields.ValueBool() {
		if isReadOnlySpecsDeleted(ctx, diags, state, plan) {
			diags.AddError(
				"Cannot remove read_only_specs attributes while toggling use_effective_fields",
				"Your configuration previously had read_only_specs attributes that were removed. "+
					"To keep read-only nodes, add the attributes back. To delete them, add the attributes back with node_count = 0. "+
					"After adding the attributes back, apply without toggling use_effective_fields, then toggle the flag in a separate apply.",
			)
		}
		if isAnalyticsSpecsDeleted(ctx, diags, state, plan) {
			diags.AddError(
				"Cannot remove analytics_specs attributes while toggling use_effective_fields",
				"Your configuration previously had analytics_specs attributes that have been removed. "+
					"To keep analytics nodes, add the attributes back. To delete them, add the attributes back with node_count = 0. "+
					"After adding the attributes back, apply without toggling use_effective_fields, then toggle the flag in a separate apply.",
			)
		}
		return
	}

	adjustRegionConfigsChildren(ctx, diags, state, plan)

	diff := findClusterDiff(ctx, state, plan, diags)
	if diags.HasError() || diff.isAnyUpgrade() { // Don't do anything in upgrades
		return
	}
	attributeChanges := schemafunc.NewAttributeChanges(ctx, state, plan)
	keepUnknown := []string{"connection_strings", "state_name", "mongo_db_version", "config_server_type"} // Volatile attributes, should not be copied from state
	keepUnknown = append(keepUnknown, attributeChanges.KeepUnknown(attributeRootChangeMapping)...)
	keepUnknown = append(keepUnknown, determineKeepUnknownsAutoScaling(ctx, diags, state, plan)...)
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown, nil)
	configs := extractAutoScalingConfigs(ctx, diags, plan)
	if !diags.HasError() {
		WarnIgnoredSpecChange(diags, plan.UseEffectiveFields.ValueBool(), attributeChanges, configs)
	}
}

// WarnIgnoredSpecChange warns when use_effective_fields=true, auto-scaling is enabled, and the user
// changed instance_size, disk_size_gb, or disk_iops — fields Atlas silently ignores in that case.
func WarnIgnoredSpecChange(diags *diag.Diagnostics, useEffectiveFields bool, attributeChanges schemafunc.AttributeChanges, configs []RegionAutoScaling) {
	if !useEffectiveFields {
		return
	}
	ignoredFields := collectIgnoredSpecChanges(attributeChanges, configs)
	if len(ignoredFields) == 0 {
		return
	}
	diags.AddWarning(
		"Spec change ignored when use_effective_fields is true and auto-scaling is enabled",
		fmt.Sprintf("Your changes to %s will be stored in Terraform state but will not modify the actual cluster in Atlas. "+
			"When use_effective_fields is true and auto-scaling is enabled, Atlas controls instance_size, disk_size_gb, and disk_iops values. "+
			"To apply your changes, disable auto-scaling and apply, then re-enable auto-scaling in a separate apply. "+
			"See: https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#manually-updating-specs-with-use_effective_fields",
			strings.Join(ignoredFields, ", ")),
	)
}

func collectIgnoredSpecChanges(attributeChanges schemafunc.AttributeChanges, configs []RegionAutoScaling) []string {
	ignoredSet := map[string]struct{}{}
	for _, cfg := range configs {
		addIgnoredSpecChangesForRegionConfig(cfg, attributeChanges, ignoredSet)
	}
	if len(ignoredSet) == 0 {
		return nil
	}
	return slices.Sorted(maps.Keys(ignoredSet))
}

// addIgnoredSpecChangesForRegionConfig adds field names that will be silently ignored by Atlas.
// Per Atlas docs: when compute OR disk auto-scaling is enabled, Atlas ignores instanceSize, diskSizeGB, and diskIOPS
// for electable and read-only nodes. For analytics nodes, only instanceSize is ignored (and only when compute is enabled).
// https://www.mongodb.com/docs/atlas/cluster-autoscaling/#enhance-auto-scaling-with-effective-fields-in-terraform
func addIgnoredSpecChangesForRegionConfig(cfg RegionAutoScaling, attributeChanges schemafunc.AttributeChanges, ignored map[string]struct{}) {
	autoScalingActive := cfg.ComputeEnabled || cfg.DiskGBEnabled
	electablePrefix := cfg.RCPrefix + ".electable_specs."
	readOnlyPrefix := cfg.RCPrefix + ".read_only_specs."
	analyticsPrefix := cfg.RCPrefix + ".analytics_specs."
	for _, change := range attributeChanges {
		switch {
		case autoScalingActive && strings.HasPrefix(change, electablePrefix):
			if field := strings.TrimPrefix(change, electablePrefix); slices.Contains(autoScalingManagedSpecFields, field) {
				ignored[field] = struct{}{}
			}
		case autoScalingActive && strings.HasPrefix(change, readOnlyPrefix):
			if field := strings.TrimPrefix(change, readOnlyPrefix); slices.Contains(autoScalingManagedSpecFields, field) {
				ignored[field] = struct{}{}
			}
		case cfg.AnalyticsComputeEnabled && change == analyticsPrefix+"instance_size":
			ignored["instance_size"] = struct{}{}
		}
	}
}

// extractAutoScalingConfigs reads per-region-config auto-scaling flags from the plan model.
func extractAutoScalingConfigs(ctx context.Context, diags *diag.Diagnostics, plan *TFModel) []RegionAutoScaling {
	planRepSpecs := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return nil
	}
	var configs []RegionAutoScaling
	for i := range planRepSpecs {
		planRegionConfigs := TFModelList[TFRegionConfigsModel](ctx, diags, planRepSpecs[i].RegionConfigs)
		if diags.HasError() {
			return nil
		}
		for j := range planRegionConfigs {
			regionConfig := &planRegionConfigs[j]
			autoScaling := TFModelObject[TFAutoScalingModel](ctx, regionConfig.AutoScaling)
			analyticsAutoScaling := TFModelObject[TFAutoScalingModel](ctx, regionConfig.AnalyticsAutoScaling)
			cfg := RegionAutoScaling{
				RCPrefix: fmt.Sprintf("replication_specs[%d].region_configs[%d]", i, j),
			}
			if autoScaling != nil {
				cfg.ComputeEnabled = autoScaling.ComputeEnabled.ValueBool()
				cfg.DiskGBEnabled = autoScaling.DiskGBEnabled.ValueBool()
			}
			if analyticsAutoScaling != nil {
				cfg.AnalyticsComputeEnabled = analyticsAutoScaling.ComputeEnabled.ValueBool()
			}
			configs = append(configs, cfg)
		}
	}
	return configs
}

// adjustRegionConfigsChildren modifies the planned values of region configs based on the current state.
// This ensures proper handling of removing auto scaling and specs attributes by preserving state values.
func adjustRegionConfigsChildren(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
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
	if !autoScalingUsed(ctx, diags, state) && !autoScalingUsed(ctx, diags, plan) {
		return nil
	}
	// When either compute or disk auto-scaling is enabled, all three fields may be adjusted by Atlas
	return slices.Clone(autoScalingManagedSpecFields)
}

// autoScalingUsed checks if auto-scaling is enabled in the given cluster model.
func autoScalingUsed(ctx context.Context, diags *diag.Diagnostics, model *TFModel) bool {
	repSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, model.ReplicationSpecs)
	for i := range repSpecsTF {
		regiongConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, repSpecsTF[i].RegionConfigs)
		for j := range regiongConfigsTF {
			for _, autoScalingTF := range []types.Object{regiongConfigsTF[j].AutoScaling, regiongConfigsTF[j].AnalyticsAutoScaling} {
				autoscaling := TFModelObject[TFAutoScalingModel](ctx, autoScalingTF)
				if autoscaling == nil {
					continue
				}
				if autoscaling.ComputeEnabled.ValueBool() || autoscaling.DiskGBEnabled.ValueBool() {
					return true
				}
			}
		}
	}
	return false
}

// isReadOnlySpecsDeleted detects if any read_only_specs block with node_count > 0 was deleted from the plan.
func isReadOnlySpecsDeleted(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) bool {
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return false
	}
	for i := range minLen(planRepSpecsTF, stateRepSpecsTF) {
		stateRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, stateRepSpecsTF[i].RegionConfigs)
		planRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, planRepSpecsTF[i].RegionConfigs)
		if diags.HasError() {
			return false
		}
		for j := range minLen(planRegionConfigsTF, stateRegionConfigsTF) {
			stateReadOnlySpecs := TFModelObject[TFSpecsModel](ctx, stateRegionConfigsTF[j].ReadOnlySpecs)
			planReadOnlySpecs := TFModelObject[TFSpecsModel](ctx, planRegionConfigsTF[j].ReadOnlySpecs)
			if stateReadOnlySpecs != nil && stateReadOnlySpecs.NodeCount.ValueInt64() > 0 && planReadOnlySpecs == nil {
				return true
			}
		}
	}
	return false
}

// isAnalyticsSpecsDeleted detects if any analytics_specs block with node_count > 0 was deleted from the plan.
func isAnalyticsSpecsDeleted(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) bool {
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return false
	}
	for i := range minLen(planRepSpecsTF, stateRepSpecsTF) {
		stateRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, stateRepSpecsTF[i].RegionConfigs)
		planRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, planRepSpecsTF[i].RegionConfigs)
		if diags.HasError() {
			return false
		}
		for j := range minLen(planRegionConfigsTF, stateRegionConfigsTF) {
			stateAnalyticsSpecs := TFModelObject[TFSpecsModel](ctx, stateRegionConfigsTF[j].AnalyticsSpecs)
			planAnalyticsSpecs := TFModelObject[TFSpecsModel](ctx, planRegionConfigsTF[j].AnalyticsSpecs)
			if stateAnalyticsSpecs != nil && stateAnalyticsSpecs.NodeCount.ValueInt64() > 0 && planAnalyticsSpecs == nil {
				return true
			}
		}
	}
	return false
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
