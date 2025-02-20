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

var (
	// TenantUpgrade changes many extra fields that are normally ok to use state values for
	tenantUpgradeRootKeepUnknown            = []string{"disk_size_gb", "cluster_id", "replication_specs", "backup_enabled", "create_date"}
	tenantUpgradeReplicationSpecKeepUnknown = []string{"disk_size_gb", "zone_id", "id", "container_id", "external_id", "auto_scaling", "analytics_specs", "read_only_specs"}
	attributeRootChangeMapping              = map[string][]string{
		// disk_size_gb can be change at any level/spec
		"disk_size_gb":      {},
		"replication_specs": {},
		"mongo_db_version":  {"mongo_db_major_version"},
	}
	attributeReplicationSpecChangeMapping = map[string][]string{
		"disk_size_gb":  {},
		"provider_name": {"ebs_volume_type"},
		"instance_size": {"disk_iops"}, // disk_iops can change based on instance_size changes
		"region_name":   {"container_id"},
		"zone_name":     {"zone_id"},
	}
)

func determineKeepUnknownsRoot(upgradeRequest *admin.LegacyAtlasTenantClusterUpgradeRequest, attributeChanges *schemafunc.AttributeChanges) []string {
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	if upgradeRequest != nil {
		// TenantUpgrade changes a few root level fields that are normally ok to use state values for
		keepUnknown = append(keepUnknown, tenantUpgradeRootKeepUnknown...)
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
	keepUnknownsUnchangedSpec := determineKeepUnknownsUnchangedReplicationSpecs(ctx, diags, state, plan, attrChanges)
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

// determineKeepUnknownsChangedReplicationSpec: These fields must be kept unknown in the replication_specs[index_of_changes]
func determineKeepUnknownsChangedReplicationSpec(keepUnknownsAlways []string, isTenantUpgrade bool, attributeChanges *schemafunc.AttributeChanges) []string {
	var keepUnknowns = slices.Clone(keepUnknownsAlways)
	if isTenantUpgrade {
		keepUnknowns = append(keepUnknowns, tenantUpgradeReplicationSpecKeepUnknown...)
	}
	return append(keepUnknowns, attributeChanges.KeepUnknown(attributeReplicationSpecChangeMapping)...)
}

func determineKeepUnknownsUnchangedReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, attributeChanges *schemafunc.AttributeChanges) []string {
	keepUnknowns := []string{}
	// Could be set to "" if we are using an ISS cluster
	if usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags) { // When using new sharding config, the legacy id must never be copied
		keepUnknowns = append(keepUnknowns, "id")
	}
	if isShardingConfigUpgrade(ctx, state, plan, diags) || attributeChanges.ListLenChanges("replication_specs") {
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
