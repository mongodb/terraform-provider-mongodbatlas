package advancedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
)

var _ resource.ResourceWithValidateConfig = &rs{}

func (r *rs) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var clusterType, databaseEdition types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("cluster_type"), &clusterType)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("database_edition"), &databaseEdition)...)
	if !resp.Diagnostics.HasError() {
		validateInfiniteClusterType(&resp.Diagnostics, clusterType.ValueString(), databaseEdition.ValueString(), "")
	}
}

func isShardedClusterType(clusterType string) bool {
	return clusterType == "SHARDED" || clusterType == "GEOSHARDED"
}

// resolveDatabaseEdition prefers effectiveDatabaseEdition: databaseEdition is only requested intent.
func resolveDatabaseEdition(databaseEdition, effectiveDatabaseEdition string) string {
	if effectiveDatabaseEdition != "" {
		return effectiveDatabaseEdition
	}
	return databaseEdition
}

// unsupportedInfiniteTopology reports whether the INFINITE edition runs on a topology this provider version rejects.
func unsupportedInfiniteTopology(clusterType, databaseEdition, effectiveDatabaseEdition string) bool {
	return isShardedClusterType(clusterType) && resolveDatabaseEdition(databaseEdition, effectiveDatabaseEdition) == "INFINITE"
}

func addUnsupportedInfiniteTopologyError(diags *diag.Diagnostics, clusterType string) {
	diags.AddAttributeError(path.Root("cluster_type"), "Unsupported INFINITE cluster type",
		fmt.Sprintf("INFINITE clusters with cluster_type = %q are not supported by this version of the MongoDB Atlas Terraform provider. Support for SHARDED and GEOSHARDED requires a newer provider version that explicitly enables these topologies once available.", clusterType))
}

func validateInfiniteClusterType(diags *diag.Diagnostics, clusterType, databaseEdition, effectiveDatabaseEdition string) {
	if unsupportedInfiniteTopology(clusterType, databaseEdition, effectiveDatabaseEdition) {
		addUnsupportedInfiniteTopologyError(diags, clusterType)
	}
}

func (r *rs) prepareUpdateDatabaseEdition(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, patch *admin.ClusterDescription20240805) {
	clusterType := plan.ClusterType.ValueString()
	if !state.DatabaseEdition.IsNull() && !plan.DatabaseEdition.IsNull() {
		return
	}
	if !isShardedClusterType(clusterType) && !hasEmptyAutoScalingDiskGB(patch.GetReplicationSpecs()) {
		return
	}
	// State or plan edition may be unset (e.g. just imported); the effective edition is authoritative either way.
	cluster, flexCluster := GetClusterDetails(ctx, diags, state.ProjectID.ValueString(), state.Name.ValueString(), r.Client, false, state.UseEffectiveFields.ValueBool())
	if diags.HasError() || flexCluster != nil {
		return
	}
	if cluster == nil {
		diags.AddError("Unable to verify cluster database edition", "The cluster no longer exists. Refresh the Terraform state before trying again.")
		return
	}
	edition, effectiveEdition := cluster.GetDatabaseEdition(), cluster.GetEffectiveDatabaseEdition()
	for _, current := range []string{clusterType, cluster.GetClusterType()} {
		if unsupportedInfiniteTopology(current, edition, effectiveEdition) {
			addUnsupportedInfiniteTopologyError(diags, current)
			return
		}
	}
	if resolveDatabaseEdition(edition, effectiveEdition) == "INFINITE" {
		omitEmptyAutoScalingChildren(patch.GetReplicationSpecs())
	}
}

func hasEmptyAutoScalingDiskGB(replicationSpecs []admin.ReplicationSpec20240805) bool {
	for _, spec := range replicationSpecs {
		for _, region := range spec.GetRegionConfigs() {
			for _, autoScaling := range []*admin.AdvancedAutoScalingSettings{region.AutoScaling, region.AnalyticsAutoScaling} {
				if autoScaling != nil && autoScaling.DiskGB != nil && !autoScaling.DiskGB.HasEnabled() && len(autoScaling.DiskGB.NullFields) == 0 {
					return true
				}
			}
		}
	}
	return false
}
