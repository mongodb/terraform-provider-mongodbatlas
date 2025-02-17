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

func useStateForUnknownsReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	// TF Models are used for CopyUnknows, Admin Models are used for PatchPayload (`json` annotations necessary)
	stateRepSpecs := newReplicationSpec20240805(ctx, state.ReplicationSpecs, diags)
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecs := newReplicationSpec20240805(ctx, plan.ReplicationSpecs, diags)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() || stateRepSpecs == nil || planRepSpecs == nil {
		return
	}
	planWithUnknowns := []TFReplicationSpecsModel{}
	useIss := clusterUseISS(planRepSpecs)
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
				if useIss {
					planRepSpecsTF[i].Id = types.StringValue("") // ISS receive ASYMMETRIC_SHARD_UNSUPPORTED error from older cluster API and therefore, the ID should be empty
				}
			} else {
				useStateForUnknownsRegionConfigs(ctx, diags, &stateSpec, &planSpec, &stateRepSpecsTF[i], &planRepSpecsTF[i])
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

func TFModelList[T any](ctx context.Context, diags *diag.Diagnostics, input types.List) []T {
	elements := make([]T, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return elements
}

func useStateForUnknownsRegionConfigs(ctx context.Context, diags *diag.Diagnostics, state, plan *admin.ReplicationSpec20240805, stateTF, planTF *TFReplicationSpecsModel) {
	stateRegionConfigs := state.GetRegionConfigs()
	stateRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, stateTF.RegionConfigs)
	planRegionConfigs := plan.GetRegionConfigs()
	planRegionConfigsTF := TFModelList[TFRegionConfigsModel](ctx, diags, planTF.RegionConfigs)
	if diags.HasError() || stateRegionConfigs == nil || planRegionConfigs == nil {
		return
	}
	for i := range planRegionConfigsTF {
		if i >= len(stateRegionConfigs) {
			continue
		}
		stateConfig := stateRegionConfigs[i]
		planConfig := planRegionConfigs[i]
		patchConfig, err := update.PatchPayload(&stateConfig, &planConfig)
		if err != nil {
			diags.AddError("error find diff useStateForUnknownsRegionConfigs", err.Error())
			return
		}
		if update.IsZeroValues(patchConfig) {
			schemafunc.CopyUnknowns(ctx, &stateRegionConfigsTF[i], &planRegionConfigsTF[i], nil)
		} else {
			useStateForUnknownsRegionConfig(ctx, diags, patchConfig, &stateRegionConfigsTF[i], &planRegionConfigsTF[i])
		}
	}
	listType, diagsLocal := types.ListValueFrom(ctx, RegionConfigsObjType, planRegionConfigsTF)
	diags.Append(diagsLocal...)
	if diags.HasError() {
		return
	}
	planTF.RegionConfigs = listType
}

func useStateForUnknownsRegionConfig(ctx context.Context, diags *diag.Diagnostics, patch *admin.CloudRegionConfig20240805, stateTF, planTF *TFRegionConfigsModel) {
	// Based on what is changed this impacts the other reginon configs too
	// In the same region config, the read_only_specs and electable_specs are the same, so changing instance_size in electable will also change it in read_only_specs
	// In the sibling region configs, the electable_specs are the same, so changing instance_size in electable will also change it in the sibling region configs
	planTF.AnalyticsSpecs = types.ObjectUnknown(SpecsObjType.AttrTypes)
}
