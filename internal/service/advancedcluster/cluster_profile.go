package advancedcluster

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// =============================================================================
// cluster_profile PROTOTYPE  (experimentation / demo branch — not production)
//
// This file is the single, self-contained home for the `cluster_profile`
// prototype. It demonstrates how one attribute (`cluster_profile`) can drive
// *conditional* defaults for other attributes (here: compute auto-scaling)
// without hardcoding them as static schema defaults — because the default is
// conditional on the value of another field, it lives in the plan-modification
// path instead.
//
// HOW THE PIECES FIT TOGETHER
//   - schema.go    declares the `cluster_profile` string attribute + validator
//                  and excludes it from the data-source schema.
//   - resource.go  calls applyClusterProfileDefaults from ModifyPlan (the
//                  plan-modification path) on both create and update.
//   - this file    holds the constants, the ordered tier list, the tier-math
//                  helper, and the actual default-resolution logic.
//
// TO FORK / TWEAK  (start here)
//   - Add/rename a profile:         edit the ClusterProfile* constants below and
//                                   the OneOf validator in schema.go.
//   - Change which defaults apply:  edit ClusterProfileAutoScaling — it is a
//                                   pure function that, given a profile + the
//                                   configured instance size, returns the
//                                   auto-scaling values to inject.
//   - Change the tier progression:  edit instanceSizeTiers / the "+2" in
//                                   InstanceSizeTwoTiersUp.
//
// SCOPE: prototype only. Edge-cases (multi-shard nuances, drift / effective
// fields, migrations, reverse-compatibility) are intentionally NOT handled.
// =============================================================================

const (
	// ClusterProfileCore is the baseline: clusters behave exactly as they do today.
	ClusterProfileCore = "CORE"
	// ClusterProfileInfinite turns on the profile-driven auto-scaling defaults.
	ClusterProfileInfinite = "INFINITE"
)

// instanceSizeTiers is the ordered list of dedicated instance sizes, smallest to
// largest. "Two tiers up" means two positions up in THIS list (the progression
// is not uniform, e.g. M60 -> M80 -> M140), so we index by position rather than
// doing arithmetic on the number.
var instanceSizeTiers = []string{
	"M10", "M20", "M30", "M40", "M50", "M60",
	"M80", "M140", "M200", "M300", "M400", "M700",
}

// InstanceSizeTwoTiersUp returns the instance size two tiers above instanceSize
// using instanceSizeTiers. If instanceSize is at (or within one of) the top, it
// is capped at the largest available size. An instance size that isn't in the
// known list is returned unchanged (prototype: we don't guess at non-standard
// tiers such as the NVMe/low-CPU families).
func InstanceSizeTwoTiersUp(instanceSize string) string {
	idx := slices.Index(instanceSizeTiers, instanceSize)
	if idx < 0 {
		return instanceSize // not in our known list: leave it as-is
	}
	twoUp := min(idx+2, len(instanceSizeTiers)-1) // cap at the max available tier
	return instanceSizeTiers[twoUp]
}

// ClusterProfileAutoScaling is the ONE place that decides what compute
// auto-scaling defaults a profile injects. It is a pure function, so it is
// trivial to unit-test and to fork.
//
//   - INFINITE, and the user has NOT set auto_scaling explicitly: enable compute
//     auto-scaling (scale up + down), min = configured instance size, max = two
//     tiers up. Returns apply=true.
//   - any other profile (CORE / unset), OR the user already configured
//     auto_scaling: returns apply=false -> baseline, nothing is changed.
//
// disk_gb_enabled is explicitly false: this prototype only turns on *compute*
// auto-scaling, and a concrete value avoids a null-vs-false plan/state mismatch.
func ClusterProfileAutoScaling(profile, instanceSize string, userConfiguredAutoScaling bool) (model TFAutoScalingModel, apply bool) {
	if profile != ClusterProfileInfinite || userConfiguredAutoScaling || instanceSize == "" {
		return TFAutoScalingModel{}, false
	}
	return TFAutoScalingModel{
		ComputeEnabled:          types.BoolValue(true),
		ComputeScaleDownEnabled: types.BoolValue(true),
		ComputeMinInstanceSize:  types.StringValue(instanceSize),
		ComputeMaxInstanceSize:  types.StringValue(InstanceSizeTwoTiersUp(instanceSize)),
		DiskGBEnabled:           types.BoolValue(false),
	}, true
}

// applyClusterProfileDefaults resolves profile-driven defaults into the plan. It
// is called from ModifyPlan on both create and update, and returns true if it
// modified the plan (so the caller knows it must write the plan back).
//
// For INFINITE clusters, every region config whose auto_scaling block was NOT
// set by the user (it is null in the *config*) gets the defaults from
// ClusterProfileAutoScaling. Explicit user auto_scaling always wins. CORE (or an
// unset profile) is a no-op so today's behavior is unchanged.
func applyClusterProfileDefaults(ctx context.Context, diags *diag.Diagnostics, config, plan *TFModel) bool {
	profile := plan.ClusterProfile.ValueString()
	if profile != ClusterProfileInfinite {
		return false // CORE / unset: baseline, nothing to do.
	}

	planRepSpecs := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	configRepSpecs := TFModelList[TFReplicationSpecsModel](ctx, diags, config.ReplicationSpecs)
	if diags.HasError() {
		return false
	}

	changed := false
	for i := range minLen(planRepSpecs, configRepSpecs) {
		planRegionConfigs := TFModelList[TFRegionConfigsModel](ctx, diags, planRepSpecs[i].RegionConfigs)
		configRegionConfigs := TFModelList[TFRegionConfigsModel](ctx, diags, configRepSpecs[i].RegionConfigs)
		if diags.HasError() {
			return false
		}
		specChanged := false
		for j := range minLen(planRegionConfigs, configRegionConfigs) {
			// "User set auto_scaling explicitly" == the block is present (non-null) in the config.
			userConfigured := TFModelObject[TFAutoScalingModel](ctx, configRegionConfigs[j].AutoScaling) != nil
			// The configured instance size comes from electable_specs (the base nodes).
			instanceSize := electableInstanceSize(ctx, planRegionConfigs[j])

			asModel, apply := ClusterProfileAutoScaling(profile, instanceSize, userConfigured)
			if !apply {
				continue
			}
			asObj, diagsLocal := types.ObjectValueFrom(ctx, autoScalingObjType.AttrTypes, asModel)
			diags.Append(diagsLocal...)
			if diags.HasError() {
				return false
			}
			planRegionConfigs[j].AutoScaling = asObj
			specChanged = true
			changed = true
		}
		if specChanged {
			listRegionConfigs, diagsLocal := types.ListValueFrom(ctx, regionConfigsObjType, planRegionConfigs)
			diags.Append(diagsLocal...)
			if diags.HasError() {
				return false
			}
			planRepSpecs[i].RegionConfigs = listRegionConfigs
		}
	}
	if !changed {
		return false
	}
	listRepSpecs, diagsLocal := types.ListValueFrom(ctx, replicationSpecsObjType, planRepSpecs)
	diags.Append(diagsLocal...)
	if diags.HasError() {
		return false
	}
	plan.ReplicationSpecs = listRepSpecs
	return true
}

// electableInstanceSize returns a region config's configured
// electable_specs.instance_size, or "" if it is not set/known.
func electableInstanceSize(ctx context.Context, rc TFRegionConfigsModel) string {
	specs := TFModelObject[TFSpecsModel](ctx, rc.ElectableSpecs)
	if specs == nil || !isKnown(specs.InstanceSize) {
		return ""
	}
	return specs.InstanceSize.ValueString()
}
