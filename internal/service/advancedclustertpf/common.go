package advancedclustertpf

import (
	"context"
	"fmt"
	"strings"
	"time"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

const (
	LegacyIgnoredLabelKey = "Infrastructure Tool"
)

var (
	ErrLegacyIgnoreLabel = fmt.Errorf("label `%s` is not supported as it is reserved for internal purposes", LegacyIgnoredLabelKey)
)

type ProcessArgs struct {
	ArgsLegacy            *admin20240530.ClusterDescriptionProcessArgs
	ArgsDefault           *admin.ClusterDescriptionProcessArgs20240805
	ClusterAdvancedConfig *admin.ApiAtlasClusterAdvancedConfiguration
}

type OldShardConfigMeta struct {
	ID       string
	NumShard int
}

func FormatMongoDBMajorVersion(version string) string {
	if strings.Contains(version, ".") {
		return version
	}
	return fmt.Sprintf("%.1f", cast.ToFloat32(version))
}

func AddIDsToReplicationSpecs(replicationSpecs []admin.ReplicationSpec20240805, zoneToReplicationSpecsIDs map[string][]string) []admin.ReplicationSpec20240805 {
	for zoneName, availableIDs := range zoneToReplicationSpecsIDs {
		indexOfIDToUse := 0
		for i := range replicationSpecs {
			if indexOfIDToUse >= len(availableIDs) {
				break // all available ids for this zone have been used
			}
			if replicationSpecs[i].GetZoneName() == zoneName {
				newID := availableIDs[indexOfIDToUse]
				indexOfIDToUse++
				replicationSpecs[i].Id = &newID
			}
		}
	}
	return replicationSpecs
}

func GetAdvancedClusterContainerID(containers []admin.CloudProviderContainer, cluster *admin.CloudRegionConfig20240805) string {
	for i, container := range containers {
		gpc := cluster.GetProviderName() == constant.GCP
		azure := container.GetProviderName() == cluster.GetProviderName() && container.GetRegion() == cluster.GetRegionName()
		aws := container.GetRegionName() == cluster.GetRegionName()
		if gpc || azure || aws {
			return containers[i].GetId()
		}
	}
	return ""
}

func GenerateFCVPinningWarningForRead(fcvPresentInState bool, apiRespFCVExpirationDate *time.Time) []diag.Diagnostic {
	pinIsActive := apiRespFCVExpirationDate != nil
	if fcvPresentInState && !pinIsActive { // pin is not active but present in state (and potentially in config file)
		warning := diag.NewWarningDiagnostic(
			"FCV pin is no longer active",
			"Please remove `pinned_fcv` from the configuration and apply changes to avoid re-pinning the FCV. Warning can be ignored if `pinned_fcv` block has been removed from the configuration.")
		return []diag.Diagnostic{warning}
	}
	if fcvPresentInState && pinIsActive {
		if time.Now().After(*apiRespFCVExpirationDate) { // pin is active, present in state, but its expiration date has passed
			warning := diag.NewWarningDiagnostic(
				"FCV pin expiration date has expired",
				"During the next maintenance window FCV will be unpinned. FCV expiration date can be extended, or `pinned_fcv` block can be removed to trigger the unpin immediately.")
			return []diag.Diagnostic{warning}
		}
	}
	return nil
}

func IsFlex(replicationSpecs *[]admin.ReplicationSpec20240805) bool {
	return getProviderName(replicationSpecs) == flexcluster.FlexClusterType
}

func getProviderName(replicationSpecs *[]admin.ReplicationSpec20240805) string {
	regionConfig := getRegionConfig(replicationSpecs)
	if regionConfig == nil {
		return ""
	}
	return regionConfig.GetProviderName()
}

func getRegionConfig(replicationSpecs *[]admin.ReplicationSpec20240805) *admin.CloudRegionConfig20240805 {
	if replicationSpecs == nil || len(*replicationSpecs) == 0 {
		return nil
	}
	replicationSpec := (*replicationSpecs)[0]
	if replicationSpec.RegionConfigs == nil || len(replicationSpec.GetRegionConfigs()) == 0 {
		return nil
	}
	return &replicationSpec.GetRegionConfigs()[0]
}

func GetPriorityOfFlexReplicationSpecs(replicationSpecs *[]admin.ReplicationSpec20240805) *int {
	regionConfig := getRegionConfig(replicationSpecs)
	if regionConfig == nil {
		return nil
	}
	return regionConfig.Priority
}

// GetReplicationSpecAttributesFromOldAPI returns the id and num shard values of replication specs coming from old API. This is used to populate replication_specs.*.id and replication_specs.*.num_shard attributes for old sharding confirgurations.
// In the old API (2023-02-01), each replications spec has a 1:1 relation with each zone, so ids and num shards are stored in a struct OldShardConfigMeta and are returned in a map from zoneName to OldShardConfigMeta.
func GetReplicationSpecAttributesFromOldAPI(ctx context.Context, projectID, clusterName string, client20240530 admin20240530.ClustersApi) (map[string]OldShardConfigMeta, error) {
	clusterOldAPI, _, err := client20240530.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, err
	}
	specs := clusterOldAPI.GetReplicationSpecs()
	result := make(map[string]OldShardConfigMeta, len(specs))
	for _, spec := range specs {
		result[spec.GetZoneName()] = OldShardConfigMeta{spec.GetId(), spec.GetNumShards()}
	}
	return result, nil
}
