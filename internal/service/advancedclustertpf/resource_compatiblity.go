package advancedclustertpf

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func overrideAttributesWithPrevStateValue(modelIn, modelOut *TFModel) {
	beforeVersion := conversion.NilForUnknown(modelIn.MongoDBMajorVersion, modelIn.MongoDBMajorVersion.ValueStringPointer())
	if beforeVersion != nil && !modelIn.MongoDBMajorVersion.Equal(modelOut.MongoDBMajorVersion) {
		modelOut.MongoDBMajorVersion = types.StringPointerValue(beforeVersion)
	}
	retainBackups := conversion.NilForUnknown(modelIn.RetainBackupsEnabled, modelIn.RetainBackupsEnabled.ValueBoolPointer())
	if retainBackups != nil && !modelIn.RetainBackupsEnabled.Equal(modelOut.RetainBackupsEnabled) {
		modelOut.RetainBackupsEnabled = types.BoolPointerValue(retainBackups)
	}
	if modelIn.DeleteOnCreateTimeout.ValueBoolPointer() != nil {
		modelOut.DeleteOnCreateTimeout = modelIn.DeleteOnCreateTimeout
	}
	overrideMapStringWithPrevStateValue(&modelIn.Labels, &modelOut.Labels)
	overrideMapStringWithPrevStateValue(&modelIn.Tags, &modelOut.Tags)
}
func overrideMapStringWithPrevStateValue(mapIn, mapOut *types.Map) {
	if mapIn == nil || mapOut == nil || len(mapOut.Elements()) > 0 {
		return
	}
	if mapIn.IsNull() {
		*mapOut = types.MapNull(types.StringType)
	} else {
		*mapOut = types.MapValueMust(types.StringType, nil)
	}
}

func findNumShardsUpdates(ctx context.Context, state, plan *TFModel, diags *diag.Diagnostics) map[string]int64 {
	if usingNewShardingConfig(ctx, plan.ReplicationSpecs, diags) {
		return nil
	}
	stateCounts := numShardsMap(ctx, state.ReplicationSpecs, diags)
	planCounts := numShardsMap(ctx, plan.ReplicationSpecs, diags)
	if diags.HasError() {
		return nil
	}
	if reflect.DeepEqual(stateCounts, planCounts) {
		return nil
	}
	return planCounts
}

func resolveAPIInfo(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, clusterLatest *admin.ClusterDescription20240805, useReplicationSpecPerShard bool) *ExtraAPIInfo {
	var (
		api20240530                = client.AtlasV220240530.ClustersApi
		projectID                  = clusterLatest.GetGroupId()
		clusterName                = clusterLatest.GetName()
		useOldShardingConfigFailed = false
	)
	clusterRespOld, _, err := api20240530.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if validate.ErrorClusterIsAsymmetrics(err) {
			useOldShardingConfigFailed = !useReplicationSpecPerShard
		} else {
			diags.AddError(errorReadLegacy20240530, defaultAPIErrorDetails(clusterName, err))
			return nil
		}
	}
	containerIDs, err := resolveContainerIDs(ctx, projectID, clusterLatest, client.AtlasV2.NetworkPeeringApi)
	if err != nil {
		diags.AddError(errorResolveContainerIDs, fmt.Sprintf("cluster name = %s, error details: %s", clusterName, err.Error()))
		return nil
	}
	return &ExtraAPIInfo{
		ContainerIDs:               containerIDs,
		ZoneNameReplicationSpecIDs: replicationSpecIDsFromOldAPI(clusterRespOld),
		UseOldShardingConfigFailed: useOldShardingConfigFailed,
		ZoneNameNumShards:          numShardsMapFromOldAPI(clusterRespOld),
		UseNewShardingConfig:       useReplicationSpecPerShard,
	}
}

// instead of using `num_shards` expand the replication specs, and set disk_size_gb
func normalizeFromTFModel(ctx context.Context, model *TFModel, diags *diag.Diagnostics, shouldExpandNumShards bool) *admin.ClusterDescription20240805 {
	latestModel := NewAtlasReq(ctx, model, diags)
	if diags.HasError() {
		return nil
	}
	counts := numShardsCounts(ctx, model.ReplicationSpecs, diags)
	if diags.HasError() {
		return nil
	}
	usingLegacySchema := isNumShardsGreaterThanOne(counts)
	if usingLegacySchema && shouldExpandNumShards {
		expandNumShards(latestModel, counts)
	}
	normalizeDiskSize(model, latestModel, diags)
	if diags.HasError() {
		return nil
	}
	return latestModel
}

func normalizeDiskSize(model *TFModel, latestModel *admin.ClusterDescription20240805, diags *diag.Diagnostics) {
	rootDiskSize := conversion.NilForUnknown(model.DiskSizeGB, model.DiskSizeGB.ValueFloat64Pointer())
	regionRootDiskSize := findFirstRegionDiskSizeGB(latestModel.ReplicationSpecs)
	if rootDiskSize != nil && regionRootDiskSize != nil && (*regionRootDiskSize-*rootDiskSize) > 0.01 {
		errMsg := fmt.Sprintf("disk_size_gb @ root != disk_size_gb @ region (%.2f!=%.2f)", *rootDiskSize, *regionRootDiskSize)
		diags.AddError(errMsg, errMsg)
		return
	}
	diskSize := rootDiskSize
	// Prefer regionRootDiskSize over rootDiskSize
	if regionRootDiskSize != nil {
		diskSize = regionRootDiskSize
	}
	if diskSize != nil {
		setDiskSize(latestModel, diskSize)
	}
}

func expandNumShards(req *admin.ClusterDescription20240805, counts []int64) {
	specs := req.GetReplicationSpecs()
	newSpecs := []admin.ReplicationSpec20240805{}
	for i, spec := range specs {
		newSpecs = append(newSpecs, spec)
		for range counts[i] - 1 {
			newSpecs = append(newSpecs, *repSpecNoIDs(spec))
		}
	}
	req.ReplicationSpecs = &newSpecs
}

func repSpecNoIDs(repspec admin.ReplicationSpec20240805) *admin.ReplicationSpec20240805 {
	repspec.Id = nil
	repspec.ZoneId = nil
	return &repspec
}

func numShardsCounts(ctx context.Context, input types.List, diags *diag.Diagnostics) []int64 {
	elements := make([]TFReplicationSpecsModel, len(input.Elements()))
	if len(elements) == 0 {
		return nil
	}
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	counts := make([]int64, len(elements))
	for i := range elements {
		item := &elements[i]
		counts[i] = item.NumShards.ValueInt64()
	}
	return counts
}

func usingNewShardingConfig(ctx context.Context, input types.List, diags *diag.Diagnostics) bool {
	counts := numShardsCounts(ctx, input, diags)
	if diags.HasError() {
		return true
	}
	return !isNumShardsGreaterThanOne(counts)
}

func numShardsMap(ctx context.Context, input types.List, diags *diag.Diagnostics) map[string]int64 {
	elements := make([]TFReplicationSpecsModel, len(input.Elements()))
	if len(elements) == 0 {
		return nil
	}
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	counts := map[string]int64{}
	for i := range elements {
		e := elements[i]
		zoneName := resolveZoneNameOrUseDefault(&e)
		counts[zoneName] = e.NumShards.ValueInt64()
	}
	return counts
}

func numShardsMapFromOldAPI(clusterRespOld *admin20240530.AdvancedClusterDescription) map[string]int64 {
	ret := make(map[string]int64)
	for i := range clusterRespOld.GetReplicationSpecs() {
		spec := &clusterRespOld.GetReplicationSpecs()[i]
		ret[spec.GetZoneName()] = int64(spec.GetNumShards())
	}
	return ret
}

func isNumShardsGreaterThanOne(counts []int64) bool {
	for _, count := range counts {
		if count > 1 {
			return true
		}
	}
	return false
}

// setDiskSize use most specific disk size, prefer region > spec > root disk size
func setDiskSize(req *admin.ClusterDescription20240805, defaultSize *float64) {
	for i, spec := range req.GetReplicationSpecs() {
		specSizeDefault := findFirstRegionDiskSizeGB(&[]admin.ReplicationSpec20240805{spec})
		if specSizeDefault == nil {
			specSizeDefault = defaultSize
		}
		for j := range spec.GetRegionConfigs() {
			actualConfig := req.GetReplicationSpecs()[i].GetRegionConfigs()[j]
			regionSize := findRegionDiskSizeGB(&actualConfig)
			if regionSize == nil {
				regionSize = specSizeDefault
			}
			analyticsSpecs := actualConfig.AnalyticsSpecs
			if analyticsSpecs != nil {
				analyticsSpecs.DiskSizeGB = regionSize
			}
			electable := actualConfig.ElectableSpecs
			if electable != nil {
				electable.DiskSizeGB = regionSize
			}
			readonly := actualConfig.ReadOnlySpecs
			if readonly != nil {
				readonly.DiskSizeGB = regionSize
			}
		}
	}
}

func findFirstRegionDiskSizeGB(specs *[]admin.ReplicationSpec20240805) *float64 {
	if specs == nil {
		return nil
	}
	for _, spec := range *specs {
		for _, regionConfig := range spec.GetRegionConfigs() {
			diskSizeGB := findRegionDiskSizeGB(&regionConfig)
			if diskSizeGB != nil {
				return diskSizeGB
			}
		}
	}
	return nil
}

func findRegionDiskSizeGB(regionConfig *admin.CloudRegionConfig20240805) *float64 {
	electable := regionConfig.ElectableSpecs
	if electable != nil && electable.DiskSizeGB != nil {
		return electable.DiskSizeGB
	}
	analyticsSpecs := regionConfig.AnalyticsSpecs
	if analyticsSpecs != nil && analyticsSpecs.DiskSizeGB != nil {
		return analyticsSpecs.DiskSizeGB
	}
	readonly := regionConfig.ReadOnlySpecs
	if readonly != nil && readonly.DiskSizeGB != nil {
		return readonly.DiskSizeGB
	}
	return nil
}

func externalIDToLegacyID(ctx context.Context, input types.List, diags *diag.Diagnostics) map[string]string {
	elements := make([]TFReplicationSpecsModel, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	idsMapped := map[string]string{}
	for i := range elements {
		e := elements[i]
		externalID := e.ExternalId.ValueString()
		legacyID := e.Id.ValueString()
		if externalID != "" && legacyID != "" {
			idsMapped[externalID] = legacyID
		}
	}
	return idsMapped
}
