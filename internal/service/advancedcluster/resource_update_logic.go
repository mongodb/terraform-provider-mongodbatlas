package advancedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

func noIDsPopulatedInReplicationSpecs(replicationSpecs *[]admin.ReplicationSpec20240805) bool {
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

func populateIDValuesUsingNewAPI(ctx context.Context, projectID, clusterName string, connV2ClusterAPI admin.ClustersApi, replicationSpecs *[]admin.ReplicationSpec20240805) (*[]admin.ReplicationSpec20240805, diag.Diagnostics) {
	if replicationSpecs == nil || len(*replicationSpecs) == 0 {
		return replicationSpecs, nil
	}
	cluster, _, err := connV2ClusterAPI.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf(errorRead, clusterName, err))
	}

	zoneToReplicationSpecsIDs := groupIDsByZone(cluster.GetReplicationSpecs())
	result := advancedclustertpf.AddIDsToReplicationSpecs(*replicationSpecs, zoneToReplicationSpecsIDs)
	return &result, nil
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
func SyncAutoScalingConfigs(replicationSpecs *[]admin.ReplicationSpec20240805) {
	if replicationSpecs == nil || len(*replicationSpecs) == 0 {
		return
	}

	var defaultAnalyticsAutoScaling, defaultAutoScaling *admin.AdvancedAutoScalingSettings

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

func applyDefaultAutoScaling(replicationSpecs *[]admin.ReplicationSpec20240805, defaultAutoScaling, defaultAnalyticsAutoScaling *admin.AdvancedAutoScalingSettings) {
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
