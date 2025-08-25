package advancedclustertpf

import (
	"context"
	"fmt"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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

func resolveAPIInfo(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, clusterLatest *admin.ClusterDescription20240805, useReplicationSpecPerShard bool) *ExtraAPIInfo {
	var (
		// api20240530                = client.AtlasV220240530.ClustersApi
		projectID                  = clusterLatest.GetGroupId()
		clusterName                = clusterLatest.GetName()
		useOldShardingConfigFailed = false
	)
	// clusterRespOld, _, err := api20240530.GetCluster(ctx, projectID, clusterName).Execute()
	// if err != nil {
	// 	if validate.ErrorClusterIsAsymmetrics(err) {
	// 		useOldShardingConfigFailed = !useReplicationSpecPerShard
	// 	} else {
	// 		diags.AddError(errorReadLegacy20240530, defaultAPIErrorDetails(clusterName, err))
	// 		return nil
	// 	}
	// }
	containerIDs, err := resolveContainerIDs(ctx, projectID, clusterLatest, client.AtlasV2.NetworkPeeringApi)
	if err != nil {
		diags.AddError(errorResolveContainerIDs, fmt.Sprintf("cluster name = %s, error details: %s", clusterName, err.Error()))
		return nil
	}
	return &ExtraAPIInfo{
		ContainerIDs:               containerIDs,
		// ZoneNameReplicationSpecIDs: replicationSpecIDsFromOldAPI(clusterRespOld),
		UseOldShardingConfigFailed: useOldShardingConfigFailed,
		// ZoneNameNumShards:          numShardsMapFromOldAPI(clusterRespOld),
		UseNewShardingConfig:       useReplicationSpecPerShard,
	}
}

func normalizeFromTFModel(ctx context.Context, model *TFModel, diags *diag.Diagnostics, shouldExpandNumShards bool) *admin.ClusterDescription20240805 {
	latestModel := NewAtlasReq(ctx, model, diags)
	if diags.HasError() {
		return nil
	}
	return latestModel
}

func numShardsMapFromOldAPI(clusterRespOld *admin20240530.AdvancedClusterDescription) map[string]int64 {
	ret := make(map[string]int64)
	for i := range clusterRespOld.GetReplicationSpecs() {
		spec := &clusterRespOld.GetReplicationSpecs()[i]
		ret[spec.GetZoneName()] = int64(spec.GetNumShards())
	}
	return ret
}
