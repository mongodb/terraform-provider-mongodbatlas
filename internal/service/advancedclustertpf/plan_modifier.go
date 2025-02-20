package advancedclustertpf

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
)

func useStateForUnknowns(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel) {
	if !schemafunc.HasUnknowns(plan) {
		return
	}
	_, upgradeRequest := findClusterDiff(ctx, state, plan, diags, &update.PatchOptions{})
	if diags.HasError() {
		return
	}
	isTenantUpgrade := upgradeRequest != nil
	attributeChanges := schemafunc.FindAttributeChanges(ctx, state, plan)
	keepUnknown := determineKeepUnknownsRoot(attributeChanges, isTenantUpgrade)
	schemafunc.CopyUnknowns(ctx, state, plan, keepUnknown)
	if slices.Contains(keepUnknown, "replication_specs") {
		useStateForUnknownsReplicationSpecs(ctx, diags, state, plan, &attributeChanges, isTenantUpgrade)
	}
}

var (
	// TenantUpgrade changes many extra fields that are normally ok to use state values for
	tenantUpgradeRootKeepUnknown            = []string{"disk_size_gb", "cluster_id", "replication_specs", "backup_enabled", "create_date"}
	tenantUpgradeReplicationSpecKeepUnknown = []string{"disk_size_gb", "zone_id", "id", "container_id", "external_id", "auto_scaling", "analytics_specs", "read_only_specs"}
	attributeRootChangeMapping              = map[string][]string{
		"disk_size_gb":           {}, // disk_size_gb can be change at any level/spec
		"replication_specs":      {},
		"mongo_db_major_version": {"mongo_db_version"},
	}
	attributeReplicationSpecChangeMapping = map[string][]string{
		// All these fields can exist in specs that are computed, therefore, it is not safe to use them when they have changed.
		"disk_iops":       {},
		"ebs_volume_type": {},
		"instance_size":   {"disk_iops"},       // disk_iops can change based on instance_size changes
		"disk_size_gb":    {},                  // disk_size_gb can be change at any level/spec
		"provider_name":   {"ebs_volume_type"}, // AWS --> AZURE will change ebs_volume_type
		"region_name":     {"container_id"},    // container_id changes based on region_name changes
		"zone_name":       {"zone_id"},         // zone_id copy from state is not safe when
	}
)

func determineKeepUnknownsRoot(attributeChanges schemafunc.AttributeChanges, isTenantUpgrade bool) []string {
	keepUnknown := []string{"connection_strings", "state_name"} // Volatile attributes, should not be copied from state
	if isTenantUpgrade {
		// TenantUpgrade changes a few root level fields that are normally ok to use state values for
		keepUnknown = append(keepUnknown, tenantUpgradeRootKeepUnknown...)
	}
	return append(keepUnknown, attributeChanges.KeepUnknown(attributeRootChangeMapping)...)
}

func useStateForUnknownsReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, attrChanges *schemafunc.AttributeChanges, isTenantUpgrade bool) {
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() {
		return
	}
	planWithUnknowns := []TFReplicationSpecsModel{}
	keepUnknownsUnchangedSpec := determineKeepUnknownsUnchangedReplicationSpecs(ctx, diags, state, plan, attrChanges)
	if diags.HasError() {
		return
	}
	for i := range planRepSpecsTF {
		if i < len(stateRepSpecsTF) {
			if attrChanges.ListIndexChanged("replication_specs", i) {
				keepUnknownsSpec := determineKeepUnknownsChangedReplicationSpec(keepUnknownsUnchangedSpec, isTenantUpgrade, attrChanges, fmt.Sprintf("replication_specs[%d]", i))
				schemafunc.CopyUnknowns(ctx, &stateRepSpecsTF[i], &planRepSpecsTF[i], keepUnknownsSpec)
			} else {
				schemafunc.CopyUnknowns(ctx, &stateRepSpecsTF[i], &planRepSpecsTF[i], keepUnknownsUnchangedSpec)
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
func determineKeepUnknownsChangedReplicationSpec(keepUnknownsAlways []string, isTenantUpgrade bool, attributeChanges *schemafunc.AttributeChanges, parentPath string) []string {
	var keepUnknowns = slices.Clone(keepUnknownsAlways)
	if isTenantUpgrade {
		keepUnknowns = append(keepUnknowns, tenantUpgradeReplicationSpecKeepUnknown...)
	}
	if attributeChanges.NestedListLenChanges(parentPath + ".region_configs") {
		keepUnknowns = append(keepUnknowns, "container_id")
	}
	return append(keepUnknowns, attributeChanges.KeepUnknown(attributeReplicationSpecChangeMapping)...)
}

func determineKeepUnknownsUnchangedReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, attributeChanges *schemafunc.AttributeChanges) []string {
	keepUnknowns := []string{}
	// Could be set to "" if we are using an ISS cluster
	if usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags) { // When using new sharding config, the legacy id must never be copied
		keepUnknowns = append(keepUnknowns, "id")
	}
	// for isShardingConfigUpgrade, it will be empty in the plan, so we need to keep it unknown
	// for listLenChanges, it might be an insertion in the middle of replication spec leading to wrong value from state copied
	if isShardingConfigUpgrade(ctx, state, plan, diags) || attributeChanges.ListLenChanges("replication_specs") {
		keepUnknowns = append(keepUnknowns, "external_id")
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
