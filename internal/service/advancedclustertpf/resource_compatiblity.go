package advancedclustertpf

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

func overrideKnowTPFIssueFields(modelIn, modelOut *TFModel) {
	beforeVersion := conversion.NilForUnknown(modelIn.MongoDBMajorVersion, modelIn.MongoDBMajorVersion.ValueStringPointer())
	if beforeVersion != nil && !modelIn.MongoDBMajorVersion.Equal(modelOut.MongoDBMajorVersion) {
		modelOut.MongoDBMajorVersion = types.StringPointerValue(beforeVersion)
	}
}

func findNumShardsUpdates(ctx context.Context, state, plan *TFModel, diags *diag.Diagnostics) map[string]int64 {
	if !usingLegacySchema(ctx, plan.ReplicationSpecs, diags) {
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

func resolveLegacyInfo(ctx context.Context, plan *TFModel, diags *diag.Diagnostics, clusterLatest *admin.ClusterDescription20240805, api20240530 admin20240530.ClustersApi) *LegacySchemaInfo {
	if !usingLegacySchema(ctx, plan.ReplicationSpecs, diags) {
		return nil
	}
	rootDiskSize := conversion.NilForUnknown(plan.DiskSizeGB, plan.DiskSizeGB.ValueFloat64Pointer())
	zoneNameSpecIDs, err := getReplicationSpecIDsFromOldAPI(ctx, plan.ProjectID.ValueString(), plan.Name.ValueString(), api20240530)
	if err != nil {
		errMsg := err.Error()
		diags.AddError(errMsg, errMsg)
		return nil
	}
	if rootDiskSize == nil {
		rootDiskSize = findRegionRootDiskSize(clusterLatest.ReplicationSpecs)
	}
	return &LegacySchemaInfo{
		ZoneNameNumShards:          numShardsMap(ctx, plan.ReplicationSpecs, diags),
		RootDiskSize:               rootDiskSize,
		ZoneNameReplicationSpecIDs: zoneNameSpecIDs,
	}
}

// instead of using `num_shards` explode the replication specs, and set disk_size_gb
func normalizeFromTFModel(ctx context.Context, model *TFModel, diags *diag.Diagnostics, shoudlExplodeNumShards bool) *admin.ClusterDescription20240805 {
	latestModel := NewAtlasReq(ctx, model, diags)
	if diags.HasError() {
		return nil
	}
	counts := numShardsCounts(ctx, model.ReplicationSpecs, diags)
	if diags.HasError() {
		return nil
	}
	usingLegacySchema := numShardsGt1(counts)
	if usingLegacySchema && shoudlExplodeNumShards {
		explodeNumShards(latestModel, counts)
	}
	rootDiskSize := conversion.NilForUnknown(model.DiskSizeGB, model.DiskSizeGB.ValueFloat64Pointer())
	regionRootDiskSize := findRegionRootDiskSize(latestModel.ReplicationSpecs)
	if rootDiskSize != nil && regionRootDiskSize != nil && (*regionRootDiskSize-*rootDiskSize) > 0.01 {
		errMsg := "disk_size_gb @ root != disk_size_gb @ region (%.2f!=%.2f)"
		diags.AddError(errMsg, errMsg)
		return nil
	}
	if rootDiskSize != nil || regionRootDiskSize != nil {
		finalDiskSize := rootDiskSize
		if finalDiskSize == nil {
			finalDiskSize = regionRootDiskSize
		}
		setDiskSize(latestModel, finalDiskSize)
	}
	return latestModel
}

// Set "Computed" Specs to nil to avoid unnecessary diffs
func normalizePatchState(cluster *admin.ClusterDescription20240805) {
	for i, specCopy := range cluster.GetReplicationSpecs() {
		for j := range specCopy.GetRegionConfigs() {
			spec := cluster.GetReplicationSpecs()[i]
			regionConfigs := *spec.RegionConfigs
			actualConfig := &regionConfigs[j]
			analyticsSpecs := actualConfig.AnalyticsSpecs
			if analyticsSpecs != nil && analyticsSpecs.NodeCount != nil && *analyticsSpecs.NodeCount == 0 {
				actualConfig.AnalyticsSpecs = nil
			}
			readonly := actualConfig.ReadOnlySpecs
			if readonly != nil && readonly.NodeCount != nil && *readonly.NodeCount == 0 {
				actualConfig.ReadOnlySpecs = nil
			}
		}
	}
}

func explodeNumShards(req *admin.ClusterDescription20240805, counts []int64) {
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

func usingLegacySchema(ctx context.Context, input types.List, diags *diag.Diagnostics) bool {
	counts := numShardsCounts(ctx, input, diags)
	if diags.HasError() {
		return false
	}
	return numShardsGt1(counts)
}

func numShardsMap(ctx context.Context, input types.List, diags *diag.Diagnostics) map[string]int64 {
	elements := make([]TFReplicationSpecsModel, len(input.Elements()))
	if localDiags := input.ElementsAs(ctx, &elements, false); len(localDiags) > 0 {
		diags.Append(localDiags...)
		return nil
	}
	counts := map[string]int64{}
	for i := range elements {
		e := elements[i]
		counts[e.ZoneName.ValueString()] = e.NumShards.ValueInt64()
	}
	return counts
}

func numShardsGt1(counts []int64) bool {
	for _, count := range counts {
		if count > 1 {
			return true
		}
	}
	return false
}

// todo: Add validation for root disk size never set together with disk_size_gb?
func setDiskSize(req *admin.ClusterDescription20240805, size *float64) {
	for i, spec := range req.GetReplicationSpecs() {
		for j := range spec.GetRegionConfigs() {
			actualConfig := req.GetReplicationSpecs()[i].GetRegionConfigs()[j]
			analyticsSpecs := actualConfig.AnalyticsSpecs
			if analyticsSpecs != nil {
				analyticsSpecs.DiskSizeGB = size
			}
			electable := actualConfig.ElectableSpecs
			if electable != nil {
				electable.DiskSizeGB = size
			}
			readonly := actualConfig.ReadOnlySpecs
			if readonly != nil {
				readonly.DiskSizeGB = size
			}
		}
	}
}

func findRegionRootDiskSize(specs *[]admin.ReplicationSpec20240805) *float64 {
	if specs == nil {
		return nil
	}
	for _, spec := range *specs {
		for _, regionConfig := range spec.GetRegionConfigs() {
			analyticsSpecs := regionConfig.AnalyticsSpecs
			if analyticsSpecs != nil && analyticsSpecs.DiskSizeGB != nil {
				return analyticsSpecs.DiskSizeGB
			}
			electable := regionConfig.ElectableSpecs
			if electable != nil && electable.DiskSizeGB != nil {
				return electable.DiskSizeGB
			}
			readonly := regionConfig.ReadOnlySpecs
			if readonly != nil && readonly.DiskSizeGB != nil {
				return readonly.DiskSizeGB
			}
		}
	}
	return nil
}
