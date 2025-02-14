package advancedclustertpf

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"go.mongodb.org/atlas-sdk/v20241113005/admin"
)

func useStateForUnknowns(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	if !schemafunc.HasUnknowns(plan) {
		return
	}
	patchReq, upgradeRequest := findClusterDiff(ctx, state, plan, diags, &update.PatchOptions{})
	if diags.HasError() {
		return
	}
	keepUnknown := determineKeepUnknowns(upgradeRequest, patchReq)
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown)
	if slices.Contains(keepUnknown, "replication_specs") {
		useStateForUnknownsReplicationSpecs(ctx, diags, state, plan)
	}
}

func useStateForUnknownsReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	// TF Models are used for CopyUnknows, Admin Models are used for PatchPayload (`json` annotations necessary)
	stateRepSpecs := newReplicationSpec20240805(ctx, state.ReplicationSpecs, diags)
	stateRepSpecsTF := replicationSpecTFModel(ctx, diags, state.ReplicationSpecs)
	planRepSpecs := newReplicationSpec20240805(ctx, plan.ReplicationSpecs, diags)
	planRepSpecsTF := replicationSpecTFModel(ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() || stateRepSpecs == nil || planRepSpecs == nil {
		return
	}
	planWithUnknowns := []TFReplicationSpecsModel{}
	for i := range planRepSpecsTF {
		if i < len(*stateRepSpecs) {
			stateSpec := (*stateRepSpecs)[i]
			planSpec := (*planRepSpecs)[i]
			patchSpec, err := update.PatchPayload(&stateSpec, &planSpec)
			if err != nil {
				diags.AddError("error find diff useStateForUnknownsReplicationSpecs", err.Error())
				return
			}
			if update.IsZeroValues(patchSpec) {
				schemafunc.CopyUnknowns(ctx, &stateRepSpecsTF[i], &planRepSpecsTF[i], nil)
			}
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

func replicationSpecTFModel(ctx context.Context, diags *diag.Diagnostics, input types.List) []TFReplicationSpecsModel {
	elements := make([]TFReplicationSpecsModel, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return elements
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
