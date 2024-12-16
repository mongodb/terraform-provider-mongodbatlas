package advancedclustertpf

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/spf13/cast"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
)

type MajorVersionOperator int

const (
	EqualOrHigher MajorVersionOperator = iota
	Higher
	EqualOrLower
	Lower
)

func MajorVersionCompatible(input *string, version float64, operator MajorVersionOperator) *bool {
	if !conversion.IsStringPresent(input) {
		return nil
	}
	value, err := strconv.ParseFloat(*input, 64)
	if err != nil {
		return nil
	}
	var result bool
	switch operator {
	case EqualOrHigher:
		result = value >= version
	case Higher:
		result = value > version
	case EqualOrLower:
		result = value <= version
	case Lower:
		result = value < version
	default:
		return nil
	}
	return &result
}

func FormatMongoDBMajorVersion(version string) string {
	if strings.Contains(version, ".") {
		return version
	}
	return fmt.Sprintf("%.1f", cast.ToFloat32(version))
}

// based on flattenAdvancedReplicationSpecRegionConfigs in model_advanced_cluster.go
func resolveContainerIDs(ctx context.Context, projectID string, cluster *admin.ClusterDescription20240805, api admin.NetworkPeeringApi) (map[string]string, error) {
	containerIDs := map[string]string{}
	responseCache := map[string]*admin.PaginatedCloudProviderContainer{}
	for _, spec := range cluster.GetReplicationSpecs() {
		for _, regionConfig := range spec.GetRegionConfigs() {
			providerName := regionConfig.GetProviderName()
			if providerName == constant.TENANT {
				continue
			}
			params := &admin.ListPeeringContainerByCloudProviderApiParams{
				GroupId:      projectID,
				ProviderName: &providerName,
			}
			containerIDKey := fmt.Sprintf("%s:%s", providerName, regionConfig.GetRegionName())
			if _, ok := containerIDs[containerIDKey]; ok {
				continue
			}
			var containersResponse *admin.PaginatedCloudProviderContainer
			var err error
			if response, ok := responseCache[providerName]; ok {
				containersResponse = response
			} else {
				containersResponse, _, err = api.ListPeeringContainerByCloudProviderWithParams(ctx, params).Execute()
				if err != nil {
					return nil, err
				}
				responseCache[providerName] = containersResponse
			}
			if results := getAdvancedClusterContainerID(containersResponse.GetResults(), &regionConfig); results != "" {
				containerIDs[containerIDKey] = results
			} else {
				return nil, fmt.Errorf("container id not found for %s", containerIDKey)
			}
		}
	}
	return containerIDs, nil
}

// copied from model_advanced_cluster.go
func getAdvancedClusterContainerID(containers []admin.CloudProviderContainer, cluster *admin.CloudRegionConfig20240805) string {
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

func replicationSpecIDsFromOldAPI(clusterRespOld *admin20240530.AdvancedClusterDescription) map[string]string {
	specs := clusterRespOld.GetReplicationSpecs()
	zoneNameSpecIDs := make(map[string]string, len(specs))
	for _, spec := range specs {
		zoneNameSpecIDs[spec.GetZoneName()] = spec.GetId()
	}
	return zoneNameSpecIDs
}

func convertHardwareSpecToOldSDK(hwspec *admin.HardwareSpec20240805) *admin20240530.HardwareSpec {
	if hwspec == nil {
		return nil
	}
	return &admin20240530.HardwareSpec{
		DiskIOPS:      hwspec.DiskIOPS,
		EbsVolumeType: hwspec.EbsVolumeType,
		InstanceSize:  hwspec.InstanceSize,
		NodeCount:     hwspec.NodeCount,
	}
}

func convertAdvancedAutoScalingSettingsToOldSDK(settings *admin.AdvancedAutoScalingSettings) *admin20240530.AdvancedAutoScalingSettings {
	if settings == nil {
		return nil
	}
	return &admin20240530.AdvancedAutoScalingSettings{
		Compute: convertAdvancedComputeAutoScalingToOldSDK(settings.Compute),
		DiskGB:  convertDiskGBAutoScalingToOldSDK(settings.DiskGB),
	}
}

func convertAdvancedComputeAutoScalingToOldSDK(settings *admin.AdvancedComputeAutoScaling) *admin20240530.AdvancedComputeAutoScaling {
	if settings == nil {
		return nil
	}
	return &admin20240530.AdvancedComputeAutoScaling{
		Enabled:          settings.Enabled,
		MaxInstanceSize:  settings.MaxInstanceSize,
		MinInstanceSize:  settings.MinInstanceSize,
		ScaleDownEnabled: settings.ScaleDownEnabled,
	}
}

func convertDiskGBAutoScalingToOldSDK(settings *admin.DiskGBAutoScaling) *admin20240530.DiskGBAutoScaling {
	if settings == nil {
		return nil
	}
	return &admin20240530.DiskGBAutoScaling{
		Enabled: settings.Enabled,
	}
}

func convertDedicatedHardwareSpecToOldSDK(spec *admin.DedicatedHardwareSpec20240805) *admin20240530.DedicatedHardwareSpec {
	if spec == nil {
		return nil
	}
	return &admin20240530.DedicatedHardwareSpec{
		NodeCount:     spec.NodeCount,
		DiskIOPS:      spec.DiskIOPS,
		EbsVolumeType: spec.EbsVolumeType,
		InstanceSize:  spec.InstanceSize,
	}
}
