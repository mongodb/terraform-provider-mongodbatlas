package advancedcluster

import (
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/atlas-sdk/v20250312013/admin"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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

func containerIDKey(providerName, regionName string) string {
	return fmt.Sprintf("%s:%s", providerName, regionName)
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
			params := &admin.ListGroupContainersApiParams{
				GroupId:      projectID,
				ProviderName: &providerName,
			}
			key := containerIDKey(providerName, regionConfig.GetRegionName())
			if _, ok := containerIDs[key]; ok {
				continue
			}
			var containersResponse *admin.PaginatedCloudProviderContainer
			var err error
			if response, ok := responseCache[providerName]; ok {
				containersResponse = response
			} else {
				containersResponse, _, err = api.ListGroupContainersWithParams(ctx, params).Execute()
				if err != nil {
					return nil, err
				}
				responseCache[providerName] = containersResponse
			}
			if results := getAdvancedClusterContainerID(containersResponse.GetResults(), &regionConfig); results != "" {
				containerIDs[key] = results
			} else {
				return nil, fmt.Errorf("container id not found for %s", key)
			}
		}
	}
	return containerIDs, nil
}

func overrideAttributesWithPrevStateValue(modelIn, modelOut *TFModel) {
	if modelIn == nil || modelOut == nil {
		return
	}
	beforeVersion := conversion.NilForUnknown(modelIn.MongoDBMajorVersion, modelIn.MongoDBMajorVersion.ValueStringPointer())
	if beforeVersion != nil && !modelIn.MongoDBMajorVersion.Equal(modelOut.MongoDBMajorVersion) {
		modelOut.MongoDBMajorVersion = types.StringPointerValue(beforeVersion)
	}
	overrideMapStringWithPrevStateValue(&modelIn.Labels, &modelOut.Labels)
	overrideMapStringWithPrevStateValue(&modelIn.Tags, &modelOut.Tags)

	// Copy Terraform-only attributes which are not returned by Atlas API.
	// These fields are Optional-only or Optional+Computed (because of a default value),
	// so no need for more complex logic as they can't be Unknown in the plan.
	modelOut.Timeouts = modelIn.Timeouts
	modelOut.DeleteOnCreateTimeout = modelIn.DeleteOnCreateTimeout
	modelOut.RetainBackupsEnabled = modelIn.RetainBackupsEnabled

	// Preserve null values for Optional-only attributes in replication_specs.
	// In v3.0.0, these attributes are Optional-only (not Computed), so if the user
	// didn't configure them, they should remain null in state even if the API returns values.
	modelOut.ReplicationSpecs = overrideReplicationSpecsWithPrevStateValue(modelIn.ReplicationSpecs, modelOut.ReplicationSpecs)
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

// overrideReplicationSpecsWithPrevStateValue preserves null values for Optional-only attributes
// in replication_specs. This ensures that if the user didn't configure an attribute, it remains
// null in state even if the API returns a value for it.
func overrideReplicationSpecsWithPrevStateValue(specsIn, specsOut types.List) types.List {
	if specsIn.IsNull() || specsIn.IsUnknown() || specsOut.IsNull() || specsOut.IsUnknown() {
		return specsOut
	}

	elemsIn := specsIn.Elements()
	elemsOut := specsOut.Elements()
	if len(elemsIn) != len(elemsOut) {
		return specsOut
	}

	ctx := context.Background()
	newElems := make([]TFReplicationSpecsModel, len(elemsOut))

	for i := range elemsOut {
		var specIn, specOut TFReplicationSpecsModel
		if diags := tfsdk.ValueAs(ctx, elemsIn[i], &specIn); diags.HasError() {
			return specsOut
		}
		if diags := tfsdk.ValueAs(ctx, elemsOut[i], &specOut); diags.HasError() {
			return specsOut
		}

		// Preserve null zone_name if it was null in input
		if specIn.ZoneName.IsNull() {
			specOut.ZoneName = types.StringNull()
		}

		// Override region_configs with preserved null values
		specOut.RegionConfigs = overrideRegionConfigsWithPrevStateValue(specIn.RegionConfigs, specOut.RegionConfigs)

		newElems[i] = specOut
	}

	newList, diags := types.ListValueFrom(ctx, replicationSpecsObjType, newElems)
	if diags.HasError() {
		return specsOut
	}
	return newList
}

// overrideRegionConfigsWithPrevStateValue preserves null values for Optional-only attributes
// in region_configs (auto_scaling, analytics_auto_scaling, electable_specs, read_only_specs, analytics_specs).
func overrideRegionConfigsWithPrevStateValue(configsIn, configsOut types.List) types.List {
	if configsIn.IsNull() || configsIn.IsUnknown() || configsOut.IsNull() || configsOut.IsUnknown() {
		return configsOut
	}

	elemsIn := configsIn.Elements()
	elemsOut := configsOut.Elements()
	if len(elemsIn) != len(elemsOut) {
		return configsOut
	}

	ctx := context.Background()
	newElems := make([]TFRegionConfigsModel, len(elemsOut))

	for i := range elemsOut {
		var configIn, configOut TFRegionConfigsModel
		if diags := tfsdk.ValueAs(ctx, elemsIn[i], &configIn); diags.HasError() {
			return configsOut
		}
		if diags := tfsdk.ValueAs(ctx, elemsOut[i], &configOut); diags.HasError() {
			return configsOut
		}

		// Preserve null values for Optional-only attributes
		if configIn.AutoScaling.IsNull() {
			configOut.AutoScaling = types.ObjectNull(autoScalingObjType.AttrTypes)
		} else {
			// Preserve null values within auto_scaling for partially configured objects
			configOut.AutoScaling = overrideAutoScalingWithPrevStateValue(configIn.AutoScaling, configOut.AutoScaling)
		}
		if configIn.AnalyticsAutoScaling.IsNull() {
			configOut.AnalyticsAutoScaling = types.ObjectNull(autoScalingObjType.AttrTypes)
		} else {
			// Preserve null values within analytics_auto_scaling for partially configured objects
			configOut.AnalyticsAutoScaling = overrideAutoScalingWithPrevStateValue(configIn.AnalyticsAutoScaling, configOut.AnalyticsAutoScaling)
		}
		if configIn.ElectableSpecs.IsNull() {
			configOut.ElectableSpecs = types.ObjectNull(specsObjType.AttrTypes)
		} else {
			// Preserve null values within electable_specs for disk_size_gb, disk_iops, ebs_volume_type
			configOut.ElectableSpecs = overrideSpecsWithPrevStateValue(configIn.ElectableSpecs, configOut.ElectableSpecs)
		}
		if configIn.ReadOnlySpecs.IsNull() {
			configOut.ReadOnlySpecs = types.ObjectNull(specsObjType.AttrTypes)
		} else {
			configOut.ReadOnlySpecs = overrideSpecsWithPrevStateValue(configIn.ReadOnlySpecs, configOut.ReadOnlySpecs)
		}
		if configIn.AnalyticsSpecs.IsNull() {
			configOut.AnalyticsSpecs = types.ObjectNull(specsObjType.AttrTypes)
		} else {
			configOut.AnalyticsSpecs = overrideSpecsWithPrevStateValue(configIn.AnalyticsSpecs, configOut.AnalyticsSpecs)
		}

		newElems[i] = configOut
	}

	newList, diags := types.ListValueFrom(ctx, regionConfigsObjType, newElems)
	if diags.HasError() {
		return configsOut
	}
	return newList
}

// overrideAutoScalingWithPrevStateValue preserves null values for Optional-only attributes within auto_scaling.
func overrideAutoScalingWithPrevStateValue(asIn, asOut types.Object) types.Object {
	if asIn.IsNull() || asIn.IsUnknown() || asOut.IsNull() || asOut.IsUnknown() {
		return asOut
	}

	ctx := context.Background()
	var autoScaleIn, autoScaleOut TFAutoScalingModel
	if diags := tfsdk.ValueAs(ctx, asIn, &autoScaleIn); diags.HasError() {
		return asOut
	}
	if diags := tfsdk.ValueAs(ctx, asOut, &autoScaleOut); diags.HasError() {
		return asOut
	}

	// Preserve null values for Optional-only attributes
	if autoScaleIn.ComputeEnabled.IsNull() {
		autoScaleOut.ComputeEnabled = types.BoolNull()
	}
	if autoScaleIn.ComputeMaxInstanceSize.IsNull() {
		autoScaleOut.ComputeMaxInstanceSize = types.StringNull()
	}
	if autoScaleIn.ComputeMinInstanceSize.IsNull() {
		autoScaleOut.ComputeMinInstanceSize = types.StringNull()
	}
	if autoScaleIn.ComputeScaleDownEnabled.IsNull() {
		autoScaleOut.ComputeScaleDownEnabled = types.BoolNull()
	}
	if autoScaleIn.DiskGBEnabled.IsNull() {
		autoScaleOut.DiskGBEnabled = types.BoolNull()
	}

	newObj, diags := types.ObjectValueFrom(ctx, autoScalingObjType.AttrTypes, autoScaleOut)
	if diags.HasError() {
		return asOut
	}
	return newObj
}

// overrideSpecsWithPrevStateValue preserves null values for Optional-only attributes within specs
// (disk_size_gb, disk_iops, ebs_volume_type).
func overrideSpecsWithPrevStateValue(specsIn, specsOut types.Object) types.Object {
	if specsIn.IsNull() || specsIn.IsUnknown() || specsOut.IsNull() || specsOut.IsUnknown() {
		return specsOut
	}

	ctx := context.Background()
	var specIn, specOut TFSpecsModel
	if diags := tfsdk.ValueAs(ctx, specsIn, &specIn); diags.HasError() {
		return specsOut
	}
	if diags := tfsdk.ValueAs(ctx, specsOut, &specOut); diags.HasError() {
		return specsOut
	}

	// Preserve null values for Optional-only attributes
	if specIn.DiskSizeGb.IsNull() {
		specOut.DiskSizeGb = types.Float64Null()
	}
	if specIn.DiskIops.IsNull() {
		specOut.DiskIops = types.Int64Null()
	}
	if specIn.EbsVolumeType.IsNull() {
		specOut.EbsVolumeType = types.StringNull()
	}

	newObj, diags := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, specOut)
	if diags.HasError() {
		return specsOut
	}
	return newObj
}
