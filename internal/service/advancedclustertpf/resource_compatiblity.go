package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

func normalizeReqModel(ctx context.Context, model *TFModel, diags *diag.Diagnostics) (legacyReq *admin20240805.ClusterDescription20240805, req *admin.ClusterDescription20240805) {
	// Ensure normal model is valid
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
		explodeNumShards(legacyModel, counts)
	}
	rootDiskSize := conversion.NilForUnknown(model.DiskSizeGB, model.DiskSizeGB.ValueFloat64Pointer())
	if rootDiskSize != nil {
		if usingLegacySchema {
			addRootDiskSizeLegacy(legacyModel, rootDiskSize)
		} else {
			addRootDiskSize(latestModel, rootDiskSize)
		}
	}
	if usingLegacySchema {
		return legacyModel, nil
	}
	return nil, latestModel
}

func explodeNumShards(req *admin20240805.ClusterDescription20240805, counts []int64) {
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

func numShardsGt1(counts []int64) bool {
	for _, count := range counts {
		if count > 1 {
			return true
		}
	}
	return false
}

// todo: Add validation for root disk size never set together with disk_size_gb
func addRootDiskSize(req *admin.ClusterDescription20240805, size *float64) {
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
