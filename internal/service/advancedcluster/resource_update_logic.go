package advancedcluster

import (
	"context"
	"fmt"

	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"

	// "go.mongodb.org/atlas-sdk/v20241113003/admin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mongodb/atlas-sdk-go/admin" // TODO: replace usage with latest once cipher config changes are in prod

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func noIDsPopulatedInReplicationSpecs(replicationSpecs *[]admin20240805.ReplicationSpec20240805) bool {
	if replicationSpecs == nil || len(*replicationSpecs) == 0 {
		return false
	}
	for _, spec := range *replicationSpecs {
		if conversion.IsStringPresent(spec.Id) {
			return false
		}
	}
	return true
}

func populateIDValuesUsingNewAPI(ctx context.Context, projectID, clusterName string, connV2ClusterAPI admin.ClustersApi, replicationSpecs *[]admin20240805.ReplicationSpec20240805) (*[]admin20240805.ReplicationSpec20240805, diag.Diagnostics) {
	if replicationSpecs == nil || len(*replicationSpecs) == 0 {
		return replicationSpecs, nil
	}
	cluster, _, err := connV2ClusterAPI.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf(errorRead, clusterName, err))
	}

	zoneToReplicationSpecsIDs := groupIDsByZone(cluster.GetReplicationSpecs())
	result := AddIDsToReplicationSpecs(*replicationSpecs, zoneToReplicationSpecsIDs)
	return &result, nil
}

func AddIDsToReplicationSpecs(replicationSpecs []admin20240805.ReplicationSpec20240805, zoneToReplicationSpecsIDs map[string][]string) []admin20240805.ReplicationSpec20240805 {
	for zoneName, availableIDs := range zoneToReplicationSpecsIDs {
		var indexOfIDToUse = 0
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

func groupIDsByZone(specs []admin.ReplicationSpec20240805) map[string][]string {
	result := make(map[string][]string)
	for _, spec := range specs {
		result[spec.GetZoneName()] = append(result[spec.GetZoneName()], spec.GetId())
	}
	return result
}

// Having the following considerations:
// - Existing replication specs can have the autoscaling values present in the state with default values even if not defined in the config (case when cluster is imported)
// - API expects autoScaling and analyticsAutoScaling aligned cross all region configs in the PATCH request
// This function is needed to avoid errors if a new replication spec is added, ensuring the PATCH request will have the auto scaling aligned with other replication specs when not present in config.
func SyncAutoScalingConfigs(replicationSpecs *[]admin20240805.ReplicationSpec20240805) {
	if replicationSpecs == nil || len(*replicationSpecs) == 0 {
		return
	}

	var defaultAnalyticsAutoScaling, defaultAutoScaling *admin20240805.AdvancedAutoScalingSettings

	for _, spec := range *replicationSpecs {
		for i := range *spec.RegionConfigs {
			regionConfig := &(*spec.RegionConfigs)[i]
			if regionConfig.AutoScaling != nil && defaultAutoScaling == nil {
				defaultAutoScaling = regionConfig.AutoScaling
			}
			if regionConfig.AnalyticsAutoScaling != nil && defaultAnalyticsAutoScaling == nil {
				defaultAnalyticsAutoScaling = regionConfig.AnalyticsAutoScaling
			}
		}
	}
	applyDefaultAutoScaling(replicationSpecs, defaultAutoScaling, defaultAnalyticsAutoScaling)
}

func applyDefaultAutoScaling(replicationSpecs *[]admin20240805.ReplicationSpec20240805, defaultAutoScaling, defaultAnalyticsAutoScaling *admin20240805.AdvancedAutoScalingSettings) {
	for _, spec := range *replicationSpecs {
		for i := range *spec.RegionConfigs {
			regionConfig := &(*spec.RegionConfigs)[i]
			if regionConfig.AutoScaling == nil && defaultAutoScaling != nil {
				regionConfig.AutoScaling = defaultAutoScaling
			}
			if regionConfig.AnalyticsAutoScaling == nil && defaultAnalyticsAutoScaling != nil {
				regionConfig.AnalyticsAutoScaling = defaultAnalyticsAutoScaling
			}
		}
	}
}
