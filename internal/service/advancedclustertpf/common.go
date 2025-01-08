package advancedclustertpf

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func FormatMongoDBMajorVersion(version string) string {
	if strings.Contains(version, ".") {
		return version
	}
	return fmt.Sprintf("%.1f", cast.ToFloat32(version))
}

func AddIDsToReplicationSpecs(replicationSpecs []admin.ReplicationSpec20240805, zoneToReplicationSpecsIDs map[string][]string) []admin.ReplicationSpec20240805 {
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
