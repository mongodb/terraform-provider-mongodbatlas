package advancedclustertpf

import (
	"fmt"
	"strings"

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
