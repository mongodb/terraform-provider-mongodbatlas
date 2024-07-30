package advancedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

func noIDsPopulatedInReplicationSpecs(replicationSpecs *[]admin.ReplicationSpec20250101) bool {
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

func populateIDValuesUsingNewAPI(ctx context.Context, projectID, clusterName string, connV2CLusterAPI admin.ClustersApi, replicationSpecs *[]admin.ReplicationSpec20250101) (*[]admin.ReplicationSpec20250101, diag.Diagnostics) {
	if replicationSpecs == nil || len(*replicationSpecs) == 0 {
		return replicationSpecs, nil
	}
	cluster, _, err := connV2CLusterAPI.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf(errorRead, clusterName, err))
	}

	zoneToReplicationSpecsIDs := groupIDsByZone(cluster.GetReplicationSpecs())
	result := AddIDsToReplicationSpecs(*replicationSpecs, zoneToReplicationSpecsIDs)
	return &result, nil
}

func AddIDsToReplicationSpecs(replicationSpecs []admin.ReplicationSpec20250101, zoneToReplicationSpecsIDs map[string][]string) []admin.ReplicationSpec20250101 {
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

func groupIDsByZone(specs []admin.ReplicationSpec20250101) map[string][]string {
	result := make(map[string][]string)
	for _, spec := range specs {
		result[spec.GetZoneName()] = append(result[spec.GetZoneName()], spec.GetId())
	}
	return result
}
