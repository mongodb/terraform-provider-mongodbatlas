package advancedclustertpf

import (
	"context"
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/spf13/cast"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

func FormatMongoDBMajorVersion(version string) string {
	if strings.Contains(version, ".") {
		return version
	}
	return fmt.Sprintf("%.1f", cast.ToFloat32(version))
}

func resolveExtraAPIInfo(ctx context.Context, projectID string, cluster *admin.ClusterDescription20240805, api admin.NetworkPeeringApi) (*ExtraAPIInfo, error) {
	containerIDs := map[string]string{}
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
			containers, _, err := api.ListPeeringContainerByCloudProviderWithParams(ctx, params).Execute()
			if err != nil {
				return nil, err
			}
			if results := getAdvancedClusterContainerID(containers.GetResults(), &regionConfig); results != "" {
				containerIDs[containerIDKey] = results
			} else {
				return nil, fmt.Errorf("container id not found for %s", containerIDKey)
			}
		}
	}
	return &ExtraAPIInfo{
		ContainerIDs: containerIDs,
	}, nil
}

// copied from model_advanced_cluster.go
func getAdvancedClusterContainerID(containers []admin.CloudProviderContainer, cluster *admin.CloudRegionConfig20240805) string {
	for i := range containers {
		gpc := cluster.GetProviderName() == constant.GCP
		azure := containers[i].GetProviderName() == cluster.GetProviderName() && containers[i].GetRegion() == cluster.GetRegionName()
		aws := containers[i].GetRegionName() == cluster.GetRegionName()
		if gpc || azure || aws {
			return containers[i].GetId()
		}
	}
	return ""
}

func getReplicationSpecIDsFromOldAPI(ctx context.Context, projectID, clusterName string, api admin20240530.ClustersApi) (map[string]string, error) {
	clusterOldAPI, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if apiError, ok := admin20240530.AsError(err); ok {
			if apiError.GetErrorCode() == "ASYMMETRIC_SHARD_UNSUPPORTED" {
				return nil, nil // if its the case of an asymmetric shard an error is expected in old API, replication_specs.*.id attribute will not be populated
			}
		}
		return nil, fmt.Errorf("error reading  advanced cluster with 2023-02-01 API (%s): %s", clusterName, err)
	}
	specs := clusterOldAPI.GetReplicationSpecs()
	result := make(map[string]string, len(specs))
	for _, spec := range specs {
		result[spec.GetZoneName()] = spec.GetId()
	}
	return result, nil
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
