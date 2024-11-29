package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

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
		rootDiskSize = findRootDiskSize(clusterLatest)
	}
	return &LegacySchemaInfo{
		ZoneNameNumShards:          numShardsMap(ctx, plan.ReplicationSpecs, diags),
		RootDiskSize:               rootDiskSize,
		ZoneNameReplicationSpecIDs: zoneNameSpecIDs,
	}
}

// instead of using `num_shards` explode the replication specs
func normalizeFromTFModel(ctx context.Context, model *TFModel, diags *diag.Diagnostics) (legacyReq *admin20240805.ClusterDescription20240805, req *admin.ClusterDescription20240805) {
	latestModel := NewAtlasReq(ctx, model, diags)
	var legacyModel *admin20240805.ClusterDescription20240805
	if diags.HasError() {
		return nil, nil
	}
	counts := numShardsCounts(ctx, model.ReplicationSpecs, diags)
	if diags.HasError() {
		return nil, nil
	}
	usingLegacySchema := numShardsGt1(counts)
	if usingLegacySchema {
		legacyModel = newLegacyModel(latestModel)
		explodeNumShardsLegacy(legacyModel, counts)
		explodeNumShards(latestModel, counts)
	}
	rootDiskSize := conversion.NilForUnknown(model.DiskSizeGB, model.DiskSizeGB.ValueFloat64Pointer())
	if rootDiskSize != nil {
		if usingLegacySchema {
			addRootDiskSizeLegacy(legacyModel, rootDiskSize)
		} else {
			addRootDiskSize(latestModel, rootDiskSize)
		}
	}
	return legacyModel, latestModel
}

func explodeNumShards(req *admin.ClusterDescription20240805, counts []int64) {
	specs := req.GetReplicationSpecs()
	newSpecs := []admin.ReplicationSpec20240805{}
	for i, spec := range specs {
		newSpecs = append(newSpecs, spec)
		for range counts[i] - 1 {
			newSpecs = append(newSpecs, spec)
		}
	}
	req.ReplicationSpecs = &newSpecs
}

func explodeNumShardsLegacy(req *admin20240805.ClusterDescription20240805, counts []int64) {
	specs := req.GetReplicationSpecs()
	newSpecs := []admin20240805.ReplicationSpec20240805{}
	for i, spec := range specs {
		newSpecs = append(newSpecs, spec)
		for range counts[i] - 1 {
			newSpecs = append(newSpecs, spec)
		}
	}
	req.ReplicationSpecs = &newSpecs
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
func addRootDiskSize(req *admin.ClusterDescription20240805, size *float64) {
	for i, spec := range req.GetReplicationSpecs() {
		for j := range spec.GetRegionConfigs() {
			actualConfig := req.GetReplicationSpecs()[i].GetRegionConfigs()[j]
			analyticsSpecs := actualConfig.AnalyticsSpecs
			if analyticsSpecs != nil && analyticsSpecs.DiskSizeGB == nil {
				analyticsSpecs.DiskSizeGB = size
			}
			electable := actualConfig.ElectableSpecs
			if electable != nil && electable.DiskSizeGB == nil {
				electable.DiskSizeGB = size
			}
			readonly := actualConfig.ReadOnlySpecs
			if readonly != nil && readonly.DiskSizeGB == nil {
				readonly.DiskSizeGB = size
			}
		}
	}
}

func findRootDiskSize(req *admin.ClusterDescription20240805) *float64 {
	for _, spec := range req.GetReplicationSpecs() {
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

func addRootDiskSizeLegacy(req *admin20240805.ClusterDescription20240805, size *float64) {
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
