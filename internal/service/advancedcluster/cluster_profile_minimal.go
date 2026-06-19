package advancedcluster

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// =============================================================================
// cluster_profile MINIMAL-CONFIG PROTOTYPE  (experimentation / demo — not production)
//
// Builds on the cluster_profile prototype (see cluster_profile.go). Goal: a user
// can deploy a working cluster with only a project_id, a name, a provider region,
// and a cluster_profile, e.g.:
//
//     resource "mongodbatlas_advanced_cluster" "example" {
//       project_id      = "<your-project-id>"
//       name            = "my-cluster"
//       provider_region = "AWS:US_EAST_1"
//       cluster_profile = "INFINITE"
//     }
//
// To make that work we made the previously-required inputs cluster_type and
// replication_specs Optional+Computed in the schema, and fill them here in the
// plan-modification path (the Framework requires Computed to let the provider set
// a value the user omitted; conditional defaults can't be schema Defaults). This
// is purely INPUT-side defaulting — no effective-fields / no surfacing of
// server-computed values.
//
// (project_id stays Required: a cluster needs a real project and there is no
// sensible default, so the user always supplies it.)
//
// PRECEDENCE: explicit user input always wins. Each default is filled ONLY when
// the field is null in the *config* (i.e. the user omitted it).
//
// TO FORK / TWEAK  (start here)
//   - Change which fields are defaulted:  edit applyMinimalConfigDefaults.
//   - Change the default VALUES:           edit the constants below and
//                                          DefaultInstanceSizeForProfile.
//   - Change the "PROVIDER:REGION" syntax: edit ParseProviderRegion.
//
// SCOPE: prototype only. Edge-cases, multi-shard, drift, migrations are NOT handled.
// =============================================================================

const (
	// defaultClusterType: per the task, cluster_type defaults to a replica set.
	defaultClusterType = "REPLICASET"
	// defaultProviderName / defaultRegionName: used when provider_region is omitted entirely.
	defaultProviderName = "AWS"
	defaultRegionName   = "US_EAST_1"
	// defaultNodeCount: a standard 3-node replica set.
	defaultNodeCount = 3
	// defaultRegionPriority: primary-election priority for a single-region cluster (max is 7).
	defaultRegionPriority = 7
)

// DefaultInstanceSizeForProfile returns the default electable instance size filled in when the
// user omits replication_specs. Keyed on cluster_profile so each profile gets a sensible base:
//   - INFINITE: M30, so the branch-1 "+2" auto-scaling default resolves to min=M30 / max=M50.
//   - CORE / unset: M10, the entry-level dedicated tier.
func DefaultInstanceSizeForProfile(profile string) string {
	if profile == ClusterProfileInfinite {
		return "M30"
	}
	return "M10"
}

// ParseProviderRegion splits the provider_region convenience input "PROVIDER:REGION"
// (e.g. "AWS:US_EAST_1") into its parts. It falls back to defaultProviderName /
// defaultRegionName when the input is empty or malformed (prototype: no strict validation).
func ParseProviderRegion(providerRegion string) (providerName, regionName string) {
	providerName, regionName = defaultProviderName, defaultRegionName
	if parts := strings.SplitN(providerRegion, ":", 2); len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		providerName, regionName = parts[0], parts[1]
	}
	return providerName, regionName
}

// applyMinimalConfigDefaults fills profile-driven defaults for the now-optional required inputs
// when the user omits them, so a minimal config resolves to a full cluster spec. It is called
// from ModifyPlan and returns true if it changed the plan. Explicit user input is never
// overwritten (each default is gated on the field being null in the config).
func applyMinimalConfigDefaults(ctx context.Context, diags *diag.Diagnostics, config, plan *TFModel) bool {
	profile := plan.ClusterProfile.ValueString() // "" / CORE / INFINITE
	changed := false

	// cluster_type: universal static default (not profile-dependent).
	if config.ClusterType.IsNull() {
		plan.ClusterType = types.StringValue(defaultClusterType)
		changed = true
	}

	// project_id is intentionally NOT defaulted — it stays Required (the user always supplies it).

	// replication_specs: synthesize a single-region spec from provider_region + profile defaults.
	if config.ReplicationSpecs.IsNull() {
		plan.ReplicationSpecs = synthesizeReplicationSpecs(ctx, diags, profile, plan.ProviderRegion.ValueString())
		changed = true
	}

	return changed
}

// synthesizeReplicationSpecs builds a one-shard, one-region replication_specs list using the
// provider_region input and profile-driven defaults. Fields the user would normally leave for
// Atlas to compute (disk size/iops, ids, zone) are set Unknown -> "known after apply".
func synthesizeReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, profile, providerRegion string) types.List {
	providerName, regionName := ParseProviderRegion(providerRegion)
	instanceSize := DefaultInstanceSizeForProfile(profile)

	electable := TFSpecsModel{
		InstanceSize:  types.StringValue(instanceSize),
		NodeCount:     types.Int64Value(defaultNodeCount),
		DiskSizeGb:    types.Float64Unknown(), // let Atlas pick the tier default
		DiskIops:      types.Int64Unknown(),
		EbsVolumeType: types.StringUnknown(),
	}
	electableObj, d := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, electable)
	diags.Append(d...)

	// auto_scaling: INFINITE injects the branch-1 compute defaults; CORE leaves it
	// "known after apply" (unset), matching a normal CORE cluster with no auto_scaling block.
	autoScaling := types.ObjectUnknown(autoScalingObjType.AttrTypes)
	if asModel, apply := ClusterProfileAutoScaling(profile, instanceSize, false); apply {
		asObj, dAS := types.ObjectValueFrom(ctx, autoScalingObjType.AttrTypes, asModel)
		diags.Append(dAS...)
		autoScaling = asObj
	}

	regionConfig := TFRegionConfigsModel{
		AnalyticsAutoScaling: types.ObjectNull(autoScalingObjType.AttrTypes),
		AnalyticsSpecs:       types.ObjectNull(specsObjType.AttrTypes),
		AutoScaling:          autoScaling,
		BackingProviderName:  types.StringNull(),
		ElectableSpecs:       electableObj,
		Priority:             types.Int64Value(defaultRegionPriority),
		ProviderName:         types.StringValue(providerName),
		ReadOnlySpecs:        types.ObjectNull(specsObjType.AttrTypes),
		RegionName:           types.StringValue(regionName),
	}
	regionConfigsList, dRC := types.ListValueFrom(ctx, regionConfigsObjType, []TFRegionConfigsModel{regionConfig})
	diags.Append(dRC...)

	repSpec := TFReplicationSpecsModel{
		RegionConfigs: regionConfigsList,
		ContainerId:   types.MapUnknown(types.StringType), // known after apply
		ExternalId:    types.StringUnknown(),
		ZoneId:        types.StringUnknown(),
		ZoneName:      types.StringUnknown(),
	}
	repSpecsList, dRS := types.ListValueFrom(ctx, replicationSpecsObjType, []TFReplicationSpecsModel{repSpec})
	diags.Append(dRS...)
	return repSpecsList
}
