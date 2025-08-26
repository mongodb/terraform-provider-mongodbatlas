package advancedclustertpf

import (
	"fmt"
	"strings"
	"time"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"

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
