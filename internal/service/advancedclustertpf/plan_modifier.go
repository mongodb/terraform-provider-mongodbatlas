package advancedclustertpf

import (
	"context"
	"reflect"
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
	// `replication_specs` is handled by index to allow:
	// 1. Using full state for "unchanged" specs
	// 2. Using partial state for "changed" specs
	if slices.Contains(keepUnknown, "replication_specs") {
		// These fields must be kept unknown in the replication_specs[index_of_changes]
		// *_specs are kept unknown as not having them in the config means that changes in "sibling" region_configs can impact the "computed" spec
		// read_only_specs also reacts to changes in the electable_specs
		// disk_size_gb can be change at any level/spec
		// disk_iops can change based on instance_size changes
		// auto_scaling can not use state value when a new region_spec/replication_spec is added, the auto_scaling will be empty and we get the AUTO_SCALINGS_MUST_BE_IN_EVERY_REGION_CONFIG error
		// 	potentially could be included if we check that the region_spec count is the same
		var keepUnknownReplicationSpecs = []string{"disk_size_gb", "disk_iops", "read_only_specs", "analytics_specs", "electable_specs", "auto_scaling"}
		if upgradeRequest != nil {
			// TenantUpgrade changes many extra fields that are normally ok to use state values for
			keepUnknownReplicationSpecs = append(keepUnknownReplicationSpecs, "zone_id", "id", "container_id", "external_id")
		}
		useStateForUnknownsReplicationSpecs(ctx, diags, state, plan, keepUnknownReplicationSpecs)
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

func useStateForUnknownsReplicationSpecs(ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, keepUnknowns []string) {
	// TF Models are used for CopyUnknows, Admin Models are used for PatchPayload (`json` annotations necessary)
	stateRepSpecs := newReplicationSpec20240805(ctx, state.ReplicationSpecs, diags)
	stateRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, state.ReplicationSpecs)
	planRepSpecs := newReplicationSpec20240805(ctx, plan.ReplicationSpecs, diags)
	planRepSpecsTF := TFModelList[TFReplicationSpecsModel](ctx, diags, plan.ReplicationSpecs)
	if diags.HasError() || stateRepSpecs == nil || planRepSpecs == nil {
		return
	}
	planWithUnknowns := []TFReplicationSpecsModel{}
	keepUnknownsAlways := []string{}
	if usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags) { // When using new sharding config, the legacy id must never be copied
		keepUnknownsAlways = append(keepUnknownsAlways, "id")
	}
	if isShardingConfigUpgrade(ctx, state, plan, diags) {
		keepUnknownsAlways = append(keepUnknownsAlways, "external_id") // Will be empty in the plan, so we need to keep it unknown
	}
	if diags.HasError() {
		return
	}
	if !zoneNamesMatch(stateRepSpecs, planRepSpecs) { // Using zone_id that doesn't match the previous zone_name can lead to errors
		keepUnknownsAlways = append(keepUnknownsAlways, "zone_id")
	}
	if !clusterUseAutoScaling(planRepSpecs) {
		keepUnknownsAlways = append(keepUnknownsAlways, "auto_scaling")
	}
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
				schemafunc.CopyUnknowns(ctx, &stateRepSpecsTF[i], &planRepSpecsTF[i], keepUnknownsAlways)
			} else {
				keepUnknownsSpec := slices.Clone(keepUnknowns)
				keepUnknownsSpec = append(keepUnknownsSpec, keepUnknownsAlways...)
				if !regionsMatch(&stateSpec, &planSpec) { // If regions are different, we need to keep the container_id unknown
					keepUnknownsSpec = append(keepUnknownsSpec, "container_id")
				}
				if !providersMatch(&stateSpec, &planSpec) { // If providers are different, we need to keep the ebs_volume_type unknown
					keepUnknownsSpec = append(keepUnknownsSpec, "ebs_volume_type")
				}
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

func TFModelList[T any](ctx context.Context, diags *diag.Diagnostics, input types.List) []T {
	elements := make([]T, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	return elements
}

func regionsMatch(state, plan *admin.ReplicationSpec20240805) bool {
	regionsState := getRegions(state)
	regionsPlan := getRegions(plan)
	return reflect.DeepEqual(regionsState, regionsPlan)
}

func getRegions(spec *admin.ReplicationSpec20240805) []string {
	regions := []string{}
	for _, region := range spec.GetRegionConfigs() {
		regions = append(regions, region.GetRegionName())
	}
	return regions
}

func providersMatch(state, plan *admin.ReplicationSpec20240805) bool {
	providersState := getProviders(state)
	providersPlan := getProviders(plan)
	return reflect.DeepEqual(providersState, providersPlan)
}

func getProviders(spec *admin.ReplicationSpec20240805) []string {
	providers := []string{}
	for _, region := range spec.GetRegionConfigs() {
		providers = append(providers, region.GetProviderName())
	}
	return providers
}

func clusterUseAutoScaling(specs *[]admin.ReplicationSpec20240805) bool {
	if specs == nil {
		return false
	}
	for _, spec := range *specs {
		for _, regionConfig := range spec.GetRegionConfigs() {
			if regionConfig.AutoScaling != nil {
				return true
			}
		}
	}
	return false
}

func zoneNamesMatch(state, plan *[]admin.ReplicationSpec20240805) bool {
	zonesState := getZoneNames(state)
	zonesPlan := getZoneNames(plan)
	return reflect.DeepEqual(zonesState, zonesPlan)
}

func getZoneNames(specs *[]admin.ReplicationSpec20240805) []string {
	zones := []string{}
	if specs == nil {
		return zones
	}
	for _, spec := range *specs {
		zones = append(zones, spec.GetZoneName())
	}
	return zones
}
