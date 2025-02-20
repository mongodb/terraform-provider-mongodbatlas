package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"go.mongodb.org/atlas-sdk/v20241113005/admin"
)

func useStateForUnknowns(ctx context.Context, diags *diag.Diagnostics, plan, state *TFModel) {
	if !schemafunc.HasUnknowns(plan) {
		return
	}
	patchReq, upgradeRequest, _ := findClusterDiff(ctx, state, plan, diags, &update.PatchOptions{})
	if diags.HasError() {
		return
	}
	keepUnknown := determineKeepUnknowns(upgradeRequest, patchReq)
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown)
}

func determineKeepUnknowns(upgradeRequest *admin.LegacyAtlasTenantClusterUpgradeRequest, patchReq *admin.ClusterDescription20240805) []string {
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	if upgradeRequest != nil {
		// TenantUpgrade changes a few root level fields that are normally ok to use state values for
		keepUnknown = append(keepUnknown, "disk_size_gb", "cluster_id", "replication_specs", "backup_enabled", "create_date")
	}
	if !update.IsZeroValues(patchReq) {
		if patchReq.MongoDBMajorVersion != nil {
			keepUnknown = append(keepUnknown, "mongo_db_version") // Not safe to set MongoDBVersion when updating MongoDBMajorVersion
		}
		if patchReq.ReplicationSpecs != nil {
			keepUnknown = append(keepUnknown, "replication_specs", "disk_size_gb") // Not safe to use root value of DiskSizeGB when updating replication specs
		}
	}
	return keepUnknown
}
