package advancedcluster

import (
	"context"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func replicationSpecsIfUpdated(ctx context.Context, planVal, stateVal types.List) (bool, []*matlas.AdvancedReplicationSpec, diag.Diagnostics) {
	var d diag.Diagnostics

	if !planVal.IsUnknown() && !stateVal.IsUnknown() { // check if updated/added
		updated, d := hasReplicationSpecsUpdated(ctx, planVal, stateVal)

		if d.HasError() {
			return false, nil, d
		}

		if updated {
			updatedSpecs, d := getUpdatedReplicationSpecs(ctx, planVal, stateVal)
			return true, updatedSpecs, d
			// return true, getUpdatedReplicationSpecs(), d
		}

		return false, nil, d
	}

	return false, nil, d
}

// getUpdatedReplicationSpecs This method creates API request objects to update replication specs.
// The API request objects are cerated by iterating over state replication_specs and replaces attribute
// values that are known in the plan. This is because the state replication_specs can have other Computed values
// while the plan only has values configured by the user or values from any plan modifiers or defaults set.
func getUpdatedReplicationSpecs(ctx context.Context, planVal, stateVal types.List) ([]*matlas.AdvancedReplicationSpec, diag.Diagnostics) {
	var diags diag.Diagnostics
	var res []*matlas.AdvancedReplicationSpec

	var planRepSpecs, stateRepSpecs []tfReplicationSpecRSModel
	if diags = planVal.ElementsAs(ctx, planRepSpecs, false); diags.HasError() {
		return nil, diags
	}
	if diags = stateVal.ElementsAs(ctx, stateRepSpecs, false); diags.HasError() {
		return nil, diags
	}

	var i int
	for i = range planRepSpecs {
		if i < len(stateRepSpecs) { // if rep_specs removed from config, we don't take them in API object list
			ss := stateRepSpecs[i]
			ps := planRepSpecs[i]

			// updatedRepSpec, diags := getUpdatedReplicationSpec(ctx, &ps, &ss)
			// if diags.HasError() {
			// 	return nil, diags
			// }
			updatedRepSpec := newReplicationSpec(ctx, &ps)
			updatedRepSpec.ID = ss.ID.ValueString()

			res = append(res, updatedRepSpec)
		}
	}

	if i <= len(planRepSpecs) { // new replication_specs added
		for i < len(planRepSpecs) {
			tmp := newReplicationSpec(ctx, &planRepSpecs[i])
			res = append(res, tmp)
			i++
		}
	}

	return res, diags
}

func getUpdatedReplicationSpec(ctx context.Context, ps, ss *tfReplicationSpecRSModel) (*matlas.AdvancedReplicationSpec, diag.Diagnostics) {
	var diags diag.Diagnostics // TODO remove

	newSpec := *newReplicationSpec(ctx, ps)

	// if v := ps.NumShards; !v.IsUnknown() {
	// 	newSpec.NumShards = int(v.ValueInt64())
	// }
	// if v := ps.ZoneName; !v.IsUnknown() {
	// 	newSpec.ZoneName = v.ValueString()
	// }

	// newSpec.RegionConfigs = newRegionConfigs(ctx, ps.RegionsConfigs)
	return &newSpec, diags
}

// func getUpdatedRegionConfigs(ctx context.Context, planVal, stateVal types.List) ([]*matlas.AdvancedRegionConfig, diag.Diagnostics) {
// 	var diags diag.Diagnostics
// 	var regionConfigs []*matlas.AdvancedRegionConfig

// 	regionConfigs = newRegionConfigs(ctx, planVal)
// 	return

// }

// hasReplicationSpecsUpdated This method checks for any attribute if known in the replication_specs plan
// should be same as it's state value. This needs to be checked as plan attributes unless configured by the user
// or any plan modifiers or defaults will always be unknown (this happens because most values in these objects
// are Optional & Computed) and hence, cannot be compared with their corresponding state value.
func hasReplicationSpecsUpdated(ctx context.Context, planVal, stateVal types.List) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	var planRepSpecs, stateRepSpecs []tfReplicationSpecRSModel
	if diags = planVal.ElementsAs(ctx, planRepSpecs, false); diags.HasError() {
		return false, diags
	}
	if diags = stateVal.ElementsAs(ctx, stateRepSpecs, false); diags.HasError() {
		return false, diags
	}

	if len(planRepSpecs) != len(stateRepSpecs) {
		return false, diags
	}

	for i := range planRepSpecs {
		if updated, d := hasReplicationSpecUpdated(ctx, &planRepSpecs[i], &stateRepSpecs[i]); d.HasError() || updated {
			return updated, append(diags, d...)
		}
	}

	return false, diags
}

func hasReplicationSpecUpdated(ctx context.Context, planRepSpec, stateRepSpec *tfReplicationSpecRSModel) (bool, diag.Diagnostics) {
	var hasUpdated bool
	var diags diag.Diagnostics

	if !planRepSpec.NumShards.IsUnknown() {
		hasUpdated = planRepSpec.NumShards.Equal(stateRepSpec.NumShards)
	}

	if !planRepSpec.ZoneName.IsUnknown() { // if user has defined zone_name in config
		hasUpdated = hasUpdated || planRepSpec.ZoneName.Equal(stateRepSpec.ZoneName)

		// if user has NOT defined zone_name in config, we set it to defaultZoneName during create so we check against that:
	} else if planRepSpec.ZoneName.IsUnknown() && stateRepSpec.ZoneName.ValueString() != defaultZoneName {
		return true, diags
	}

	// TODO refactor
	if updated, d := hasRegionConfigsUpdated(ctx, planRepSpec.RegionsConfigs, stateRepSpec.RegionsConfigs); d.HasError() || updated {
		return updated, d
	}

	return hasUpdated, diags
}

func hasRegionConfigsUpdated(ctx context.Context, planVal, stateVal types.List) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	var planRegionConfigs, stateRegionConfigs []tfRegionsConfigModel
	planVal.ElementsAs(ctx, planRegionConfigs, false)
	stateVal.ElementsAs(ctx, stateRegionConfigs, false)

	if len(planRegionConfigs) != len(stateRegionConfigs) {
		return true, diags
	}

	for i := range planRegionConfigs {
		hasUpdated, d := hasRegionConfigUpdated(ctx, &planRegionConfigs[i], &stateRegionConfigs[i])

		if hasUpdated || d.HasError() {
			return hasUpdated, d
		}
	}

	return false, diags
}

func hasRegionConfigUpdated(ctx context.Context, planRegionConfig, stateRegionConfig *tfRegionsConfigModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// TODO refactor
	if !planRegionConfig.BackingProviderName.IsUnknown() && !planRegionConfig.BackingProviderName.Equal(stateRegionConfig.BackingProviderName) {
		return true, diags
	}
	if !planRegionConfig.Priority.IsUnknown() && !planRegionConfig.Priority.Equal(stateRegionConfig.Priority) {
		return true, diags
	}
	if !planRegionConfig.RegionName.IsUnknown() && !planRegionConfig.RegionName.Equal(stateRegionConfig.RegionName) {
		return true, diags
	}
	if !planRegionConfig.ProviderName.IsUnknown() && !planRegionConfig.ProviderName.Equal(stateRegionConfig.ProviderName) {
		return true, diags
	}

	if updated, d := hasRegionConfigSpecUpdated(ctx, planRegionConfig.AnalyticsSpecs, stateRegionConfig.AnalyticsSpecs); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigSpecUpdated(ctx, planRegionConfig.ElectableSpecs, stateRegionConfig.ElectableSpecs); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigSpecUpdated(ctx, planRegionConfig.ReadOnlySpecs, stateRegionConfig.ReadOnlySpecs); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigAutoScalingSpecUpdated(ctx, planRegionConfig.AutoScaling, stateRegionConfig.AutoScaling); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigAutoScalingSpecUpdated(ctx, planRegionConfig.AnalyticsAutoScaling, stateRegionConfig.AnalyticsAutoScaling); d.HasError() || updated {
		return updated, append(diags, d...)
	}

	return false, diags
}

func hasRegionConfigAutoScalingSpecUpdated(ctx context.Context, planVal, stateVal types.List) (bool, diag.Diagnostics) {
	var planSpecs, stateSpecs []tfRegionsConfigAutoScalingSpecsModel
	var diags diag.Diagnostics

	if d := planVal.ElementsAs(ctx, &planSpecs, false); diags.HasError() {
		return true, append(diags, d...)
	}
	if d := stateVal.ElementsAs(ctx, &stateSpecs, false); diags.HasError() {
		return true, append(diags, d...)
	}

	if len(planSpecs) != len(stateSpecs) {
		return true, diags
	}

	ps := planSpecs[0]
	ss := stateSpecs[0]

	if !ps.ComputeEnabled.IsUnknown() && !ps.ComputeEnabled.Equal(ss.ComputeEnabled) ||
		!ps.ComputeMaxInstanceSize.IsUnknown() && !ps.ComputeMaxInstanceSize.Equal(ss.ComputeMaxInstanceSize) ||
		!ps.ComputeMinInstanceSize.IsUnknown() && !ps.ComputeMinInstanceSize.Equal(ss.ComputeMinInstanceSize) ||
		!ps.ComputeScaleDownEnabled.IsUnknown() && !ps.ComputeScaleDownEnabled.Equal(ss.ComputeScaleDownEnabled) ||
		!ps.DiskGBEnabled.IsUnknown() && !ps.DiskGBEnabled.Equal(ss.DiskGBEnabled) {
		return true, diags
	}

	return false, diags
}

func hasRegionConfigSpecUpdated(ctx context.Context, planVal, stateVal types.List) (bool, diag.Diagnostics) {
	var planSpecs, stateSpecs []tfRegionsConfigSpecsModel
	var diags diag.Diagnostics

	if d := planVal.ElementsAs(ctx, &planSpecs, false); diags.HasError() {
		return true, append(diags, d...)
	}
	if d := stateVal.ElementsAs(ctx, &stateSpecs, false); diags.HasError() {
		return true, append(diags, d...)
	}

	if len(planSpecs) != len(stateSpecs) {
		return true, diags
	}

	ps := planSpecs[0]
	ss := stateSpecs[0]

	// TODO refactor
	var hasUpdated bool
	if v := ps.DiskIOPS; !v.IsUnknown() {
		hasUpdated = v.Equal(ss.DiskIOPS)
	}
	if v := ps.EBSVolumeType; !v.IsUnknown() {
		hasUpdated = hasUpdated || v.Equal(ss.EBSVolumeType)
	}
	if v := ps.InstanceSize; !v.IsUnknown() {
		hasUpdated = hasUpdated || v.Equal(ss.InstanceSize)
	}
	if v := ps.NodeCount; !v.IsUnknown() {
		hasUpdated = hasUpdated || v.Equal(ss.NodeCount)
	}

	return hasUpdated, diags
}
