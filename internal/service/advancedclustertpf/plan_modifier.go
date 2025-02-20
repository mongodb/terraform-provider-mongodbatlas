package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"go.mongodb.org/atlas-sdk/v20241113005/admin"
)

var keepUnknownTenantUpgrade = []string{"disk_size_gb", "cluster_id", "replication_specs", "backup_enabled", "create_date"}
var keepUnknownFlexUpgrade = []string{"disk_size_gb", "encryption_at_rest_provider", "replication_specs", "backup_enabled", "cluster_id", "create_date", "root_cert_type", "bi_connector_config"}

func useStateForUnknowns(ctx context.Context, diags *diag.Diagnostics, plan, state *TFModel) {
	if !schemafunc.HasUnknowns(plan) {
		return
	}
	stateReq := normalizeFromTFModel(ctx, state, diags, false)
	planReq := normalizeFromTFModel(ctx, plan, diags, false)
	if diags.HasError() {
		return
	}
	flexUpgrade, _ := flexChanges(planReq, stateReq, diags)
	if diags.HasError() {
		return
	}
	if flexUpgrade {
		keepUnknown := []string{"connection_strings", "state_name", "advanced_configuration", "encryption_at_rest_provider", "root_cert_type", "bi_connector_config"}
		keepUnknown = append(keepUnknown, keepUnknownTenantUpgrade...)
		schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown)
		return
	}

	patchReq, upgradeRequest, upgradeFlexRequest := findClusterDiff(ctx, state, plan, diags, &update.PatchOptions{})
	if diags.HasError() {
		return
	}
	keepUnknown := determineKeepUnknowns(upgradeRequest, upgradeFlexRequest, patchReq)
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown)
}

func determineKeepUnknowns(upgradeRequest *admin.LegacyAtlasTenantClusterUpgradeRequest, upgradeFlexRequest *admin.AtlasTenantClusterUpgradeRequest20240805, patchReq *admin.ClusterDescription20240805) []string {
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	if upgradeRequest != nil {
		// TenantUpgrade changes a few root level fields that are normally ok to use state values for
		keepUnknown = append(keepUnknown, keepUnknownTenantUpgrade...)
	}
	if upgradeFlexRequest != nil {
		// FlexToDedicatedUpgrade changes a few root level fields that are normally ok to use state values for
		keepUnknown = append(keepUnknown, keepUnknownFlexUpgrade...)
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
