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
	_, upgradeRequest := findClusterDiff(ctx, state, plan, diags, &update.PatchOptions{})
	if diags.HasError() {
		return
	}
	attributeChanges := schemafunc.FindAttributeChanges(ctx, state, plan)
	keepUnknown := determineKeepUnknownsRoot(upgradeRequest, &attributeChanges)
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown)
	// `replication_specs` is handled by index to allow:
	// 1. Using full state for "unchanged" specs
	// 2. Using partial state for "changed" specs
	if slices.Contains(keepUnknown, "replication_specs") {
		useStateForUnknownsReplicationSpecs(ctx, diags, state, plan, &attributeChanges, upgradeRequest != nil)
	}
}

var attributeRootChangeMapping = map[string][]string{
	"disk_size_gb":      {},
	"replication_specs": {},
	"mongo_db_version":  {"mongo_db_major_version"},
}
var attributeReplicationSpecChangeMapping = map[string][]string{
	"disk_size_gb":  {},
	"provider_name": {"ebs_volume_type"},
	"instance_size": {"disk_iops"},
	"region_name":   {"container_id"},
	"zone_name":     {"zone_id"},
}

func determineKeepUnknownsRoot(upgradeRequest *admin.LegacyAtlasTenantClusterUpgradeRequest, attributeChanges *schemafunc.AttributeChanges) []string {
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	if upgradeRequest != nil {
		// TenantUpgrade changes a few root level fields that are normally ok to use state values for
		keepUnknown = append(keepUnknown, "disk_size_gb", "cluster_id", "replication_specs", "backup_enabled", "create_date")
	}
	if attributeChanges != nil {
		keepUnknown = append(keepUnknown, attributeChanges.KeepUnknown(attributeRootChangeMapping)...)
	}
	return keepUnknown
}

func useStateForUnknownsReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, attrChanges *schemafunc.AttributeChanges, isTenantUpgrade bool) {
	// TF Models are used for CopyUnknows, Admin Models are used for PatchPayload (`json` annotations necessary)
	stateRepSpecs := newReplicationSpec20240805(ctx, state.ReplicationSpecs, diags)
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecs := newReplicationSpec20240805(ctx, plan.ReplicationSpecs, diags)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() || stateRepSpecs == nil || planRepSpecs == nil {
		return
	}
	planWithUnknowns := []TFReplicationSpecsModel{}
	keepUnknownsUnchangedSpec := determineKeepUnknownsUnchangedReplicationSpecs(ctx, diags, state, plan)
	if diags.HasError() {
		return
	}
	for i := range planRepSpecsTF {
		if i < len(*stateRepSpecs) {
			stateSpec := (*stateRepSpecs)[i]
			planSpec := (*planRepSpecs)[i]
			patchSpec, err := update.PatchPayload(&stateSpec, &planSpec) // TODO: Replace with attrChanges.listChanges(name, index)
			if err != nil {
				diags.AddError("error find diff useStateForUnknownsReplicationSpecs", err.Error())
				return
			}
			if update.IsZeroValues(patchSpec) {
				schemafunc.CopyUnknowns(ctx, &stateRepSpecsTF[i], &planRepSpecsTF[i], keepUnknownsUnchangedSpec)
			} else {
				keepUnknownsSpec := determineKeepUnknownsChangedReplicationSpec(keepUnknownsUnchangedSpec, isTenantUpgrade, attrChanges)
				schemafunc.CopyUnknowns(ctx, &stateRepSpecsTF[i], &planRepSpecsTF[i], keepUnknownsSpec)
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

func determineKeepUnknownsChangedReplicationSpec(keepUnknownsAlways []string, isTenantUpgrade bool, attributeChanges *schemafunc.AttributeChanges) []string {
	// These fields must be kept unknown in the replication_specs[index_of_changes]
	// *_specs are kept unknown as not having them in the config means that changes in "sibling" region_configs can impact the "computed" spec
	// read_only_specs also reacts to changes in the electable_specs
	// disk_size_gb can be change at any level/spec
	// disk_iops can change based on instance_size changes
	// auto_scaling can not use state value when a new region_spec/replication_spec is added, the auto_scaling will be empty and we get the AUTO_SCALINGS_MUST_BE_IN_EVERY_REGION_CONFIG error
	// 	potentially could be included if we check that the region_spec count is the same
	var keepUnknowns = []string{}
	if isTenantUpgrade {
		// TenantUpgrade changes many extra fields that are normally ok to use state values for
		keepUnknowns = append(keepUnknowns, "zone_id", "id", "container_id", "external_id")
	}
	keepUnknowns = append(keepUnknowns, keepUnknownsAlways...)
	return append(keepUnknowns, attributeChanges.KeepUnknown(attributeReplicationSpecChangeMapping)...)
}

func determineKeepUnknownsUnchangedReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) []string {
	keepUnknowns := []string{}
	// Could be set to "" if we are using an ISS cluster
	if usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags) { // When using new sharding config, the legacy id must never be copied
		keepUnknowns = append(keepUnknowns, "id")
	}
	if isShardingConfigUpgrade(ctx, state, plan, diags) {
		keepUnknowns = append(keepUnknowns, "external_id") // Will be empty in the plan, so we need to keep it unknown
	}
	return keepUnknowns
}

func TFModelList[T any](ctx context.Context, diags *diag.Diagnostics, input types.List) []T {
	elements := make([]T, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return elements
}
