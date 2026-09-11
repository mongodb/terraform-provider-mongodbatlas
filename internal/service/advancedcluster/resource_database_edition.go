package advancedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
)

// resolveDatabaseEdition prefers effectiveDatabaseEdition: databaseEdition is only requested intent.
func resolveDatabaseEdition(databaseEdition, effectiveDatabaseEdition string) string {
	if effectiveDatabaseEdition != "" {
		return effectiveDatabaseEdition
	}
	return databaseEdition
}

// prepareUpdateDatabaseEdition omits the empty auto-scaling children that INFINITE rejects. It resolves the
// edition from Atlas when state or plan leaves it unset, e.g. after an import, because only Atlas knows it then.
func (r *rs) prepareUpdateDatabaseEdition(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, patch *admin.ClusterDescription20240805) {
	if !state.DatabaseEdition.IsNull() && !plan.DatabaseEdition.IsNull() {
		return
	}
	if !hasEmptyAutoScalingDiskGB(patch.GetReplicationSpecs()) {
		return
	}
	cluster, flexCluster := GetClusterDetails(ctx, diags, state.ProjectID.ValueString(), state.Name.ValueString(), r.Client, false, state.UseEffectiveFields.ValueBool())
	if diags.HasError() || flexCluster != nil {
		return
	}
	if cluster == nil {
		diags.AddError("Unable to verify cluster database edition", "The cluster no longer exists. Refresh the Terraform state before trying again.")
		return
	}
	if resolveDatabaseEdition(cluster.GetDatabaseEdition(), cluster.GetEffectiveDatabaseEdition()) == "INFINITE" {
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
