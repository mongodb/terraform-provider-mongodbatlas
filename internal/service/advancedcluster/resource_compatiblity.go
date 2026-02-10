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

func overrideAttributesWithPrevStateValue(ctx context.Context, modelIn, modelOut *TFModel) {
	if modelIn == nil || modelOut == nil {
		return
	}
	overrideMapStringWithPrevStateValue(&modelIn.Labels, &modelOut.Labels)
	overrideMapStringWithPrevStateValue(&modelIn.Tags, &modelOut.Tags)

	// Copy Terraform-only attributes which are not returned by Atlas API.
	// These fields are Optional-only or Optional+Computed (because of a default value),
	// so no need for more complex logic as they can't be Unknown in the plan.
	modelOut.Timeouts = modelIn.Timeouts
	modelOut.DeleteOnCreateTimeout = modelIn.DeleteOnCreateTimeout
	modelOut.RetainBackupsEnabled = modelIn.RetainBackupsEnabled

	// Preserve null values for Optional-only attribute bi_connector_config.
	// In v3.0.0, this attribute is Optional-only (not Computed), so if the user
	// didn't configure it, it should remain null in state even if the API returns values.
	if modelIn.BiConnectorConfig.IsNull() {
		modelOut.BiConnectorConfig = types.ObjectNull(biConnectorConfigObjType.AttrTypes)
	} else {
		// Preserve null values for partially configured bi_connector_config.
		modelOut.BiConnectorConfig = overrideBiConnectorConfigWithPrevStateValue(ctx, modelIn.BiConnectorConfig, modelOut.BiConnectorConfig)
	}
	// Note: advanced_configuration null preservation is handled by updateModelAdvancedConfig
	// because newTFModel always sets AdvancedConfiguration to null (it's populated from a separate API endpoint).

	// Preserve null values for other Optional-only attributes.
	// These were changed from Optional+Computed to Optional-only in v3.0.0.
	if modelIn.BackupEnabled.IsNull() {
		modelOut.BackupEnabled = types.BoolNull()
	}
	if modelIn.EncryptionAtRestProvider.IsNull() {
		modelOut.EncryptionAtRestProvider = types.StringNull()
	}
	if modelIn.GlobalClusterSelfManagedSharding.IsNull() {
		modelOut.GlobalClusterSelfManagedSharding = types.BoolNull()
	}
	if modelIn.MongoDBMajorVersion.IsNull() {
		modelOut.MongoDBMajorVersion = types.StringNull()
	}
	if modelIn.Paused.IsNull() {
		modelOut.Paused = types.BoolNull()
	}
	if modelIn.PitEnabled.IsNull() {
		modelOut.PitEnabled = types.BoolNull()
	}
	if modelIn.RedactClientLogData.IsNull() {
		modelOut.RedactClientLogData = types.BoolNull()
	}
	if modelIn.ReplicaSetScalingStrategy.IsNull() {
		modelOut.ReplicaSetScalingStrategy = types.StringNull()
	}
	if modelIn.TerminationProtectionEnabled.IsNull() {
		modelOut.TerminationProtectionEnabled = types.BoolNull()
	}
	if modelIn.VersionReleaseSystem.IsNull() {
		modelOut.VersionReleaseSystem = types.StringNull()
	}
	if modelIn.ConfigServerManagementMode.IsNull() {
		modelOut.ConfigServerManagementMode = types.StringNull()
	}
	if modelIn.PinnedFCV.IsNull() {
		modelOut.PinnedFCV = types.ObjectNull(pinnedFCVObjType.AttrTypes)
	}

	modelOut.ReplicationSpecs = overrideReplicationSpecsWithPrevStateValue(ctx, modelIn.ReplicationSpecs, modelOut.ReplicationSpecs)
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

// overrideBiConnectorConfigWithPrevStateValue preserves null values for Optional-only attributes
// within bi_connector_config when partially configured.
func overrideBiConnectorConfigWithPrevStateValue(ctx context.Context, biIn, biOut types.Object) types.Object {
	if biIn.IsNull() || biIn.IsUnknown() || biOut.IsNull() || biOut.IsUnknown() {
		return biOut
	}

	var configIn, configOut TFBiConnectorModel
	if diags := tfsdk.ValueAs(ctx, biIn, &configIn); diags.HasError() {
		return biOut
	}
	if diags := tfsdk.ValueAs(ctx, biOut, &configOut); diags.HasError() {
		return biOut
	}

	// Preserve null values for Optional-only attributes
	if configIn.Enabled.IsNull() {
		configOut.Enabled = types.BoolNull()
	}
	if configIn.ReadPreference.IsNull() {
		configOut.ReadPreference = types.StringNull()
	}

	newObj, diags := types.ObjectValueFrom(ctx, biConnectorConfigObjType.AttrTypes, configOut)
	if diags.HasError() {
		return biOut
	}
	return newObj
}

// overrideAdvancedConfigurationWithPrevStateValue preserves null values for Optional-only attributes
// within advanced_configuration when partially configured.
func overrideAdvancedConfigurationWithPrevStateValue(ctx context.Context, acIn, acOut types.Object) types.Object {
	if acIn.IsNull() || acIn.IsUnknown() || acOut.IsNull() || acOut.IsUnknown() {
		return acOut
	}

	var configIn, configOut TFAdvancedConfigurationModel
	if diags := tfsdk.ValueAs(ctx, acIn, &configIn); diags.HasError() {
		return acOut
	}
	if diags := tfsdk.ValueAs(ctx, acOut, &configOut); diags.HasError() {
		return acOut
	}

	// Preserve null values for Optional-only attributes
	if configIn.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds.IsNull() {
		configOut.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds = types.Int64Null()
	}
	if configIn.DefaultWriteConcern.IsNull() {
		configOut.DefaultWriteConcern = types.StringNull()
	}
	if configIn.JavascriptEnabled.IsNull() {
		configOut.JavascriptEnabled = types.BoolNull()
	}
	if configIn.MinimumEnabledTlsProtocol.IsNull() {
		configOut.MinimumEnabledTlsProtocol = types.StringNull()
	}
	if configIn.NoTableScan.IsNull() {
		configOut.NoTableScan = types.BoolNull()
	}
	if configIn.OplogMinRetentionHours.IsNull() {
		configOut.OplogMinRetentionHours = types.Float64Null()
	}
	if configIn.OplogSizeMb.IsNull() {
		configOut.OplogSizeMb = types.Int64Null()
	}
	if configIn.SampleRefreshIntervalBiconnector.IsNull() {
		configOut.SampleRefreshIntervalBiconnector = types.Int64Null()
	}
	if configIn.SampleSizeBiconnector.IsNull() {
		configOut.SampleSizeBiconnector = types.Int64Null()
	}
	if configIn.TransactionLifetimeLimitSeconds.IsNull() {
		configOut.TransactionLifetimeLimitSeconds = types.Int64Null()
	}
	if configIn.DefaultMaxTimeMS.IsNull() {
		configOut.DefaultMaxTimeMS = types.Int64Null()
	}
	if configIn.CustomOpensslCipherConfigTls12.IsNull() {
		configOut.CustomOpensslCipherConfigTls12 = types.SetNull(types.StringType)
	}
	if configIn.CustomOpensslCipherConfigTls13.IsNull() {
		configOut.CustomOpensslCipherConfigTls13 = types.SetNull(types.StringType)
	}
	if configIn.TlsCipherConfigMode.IsNull() {
		configOut.TlsCipherConfigMode = types.StringNull()
	}

	newObj, diags := types.ObjectValueFrom(ctx, advancedConfigurationObjType.AttrTypes, configOut)
	if diags.HasError() {
		return acOut
	}
	return newObj
}

// overrideReplicationSpecsWithPrevStateValue preserves null values for Optional-only attributes in replication_specs
// when the user didn't configure them. Without this, the API response values would replace null, causing an
// "inconsistent result after apply" error.
//
// Background: the provider (TF v3) always sends the Use-Effective-Fields-Replication-Specs header. When Atlas
// receives a Create or Update with this header, it changes behavior: replicationSpecs in the response starts
// echoing back only what Terraform sent in the request, while effectiveReplicationSpecs contains the full
// computed values (with disk_iops, ebs_volume_type, auto_scaling defaults, etc.). However, for existing clusters
// that have never been created or updated with this header, the Atlas API still returns replicationSpecs identical
// to effectiveReplicationSpecs (all fields populated). This means on the first Update in TF v3 — whether from a
// cluster-to-advanced_cluster migration, a provider v2-to-v3 upgrade, or any cluster not yet updated with the
// header — the API response includes fields the user didn't configure, causing plan/state inconsistencies.
func overrideReplicationSpecsWithPrevStateValue(ctx context.Context, specsIn, specsOut types.List) types.List {
	if specsOut.IsNull() || specsOut.IsUnknown() {
		return specsOut
	}
	if specsIn.IsNull() || specsIn.IsUnknown() {
		return nullifyOptionalOnlyInReplicationSpecs(ctx, specsOut)
	}

	elemsIn := specsIn.Elements()
	elemsOut := specsOut.Elements()
	if len(elemsIn) != len(elemsOut) {
		return specsOut
	}

	newElems := make([]TFReplicationSpecsModel, len(elemsOut))
	for i := range elemsOut {
		var specIn, specOut TFReplicationSpecsModel
		if diags := tfsdk.ValueAs(ctx, elemsIn[i], &specIn); diags.HasError() {
			return specsOut
		}
		if diags := tfsdk.ValueAs(ctx, elemsOut[i], &specOut); diags.HasError() {
			return specsOut
		}

		if specIn.ZoneName.IsNull() {
			specOut.ZoneName = types.StringNull()
		}

		specOut.RegionConfigs = overrideRegionConfigsWithPrevStateValue(ctx, specIn.RegionConfigs, specOut.RegionConfigs)

		newElems[i] = specOut
	}

	newList, diags := types.ListValueFrom(ctx, replicationSpecsObjType, newElems)
	if diags.HasError() {
		return specsOut
	}
	return newList
}

// nullifyOptionalOnlyInReplicationSpecs nullifies Optional-only attributes in replication_specs when the previous
// state has no replication_specs (e.g., after a state upgrade or move from cluster to advanced_cluster).
// This includes zone_name at the spec level and various sub-attributes within region_configs.
func nullifyOptionalOnlyInReplicationSpecs(ctx context.Context, specsOut types.List) types.List {
	elemsOut := specsOut.Elements()
	newElems := make([]TFReplicationSpecsModel, len(elemsOut))
	for i := range elemsOut {
		var specOut TFReplicationSpecsModel
		if diags := tfsdk.ValueAs(ctx, elemsOut[i], &specOut); diags.HasError() {
			return specsOut
		}
		specOut.ZoneName = types.StringNull()
		specOut.RegionConfigs = nullifyOptionalOnlyInRegionConfigs(ctx, specOut.RegionConfigs)
		newElems[i] = specOut
	}
	newList, diags := types.ListValueFrom(ctx, replicationSpecsObjType, newElems)
	if diags.HasError() {
		return specsOut
	}
	return newList
}

// nullifyOptionalOnlyInRegionConfigs nullifies Optional-only attributes within each region_config
// when there is no previous state to compare against (e.g., after a move from cluster to advanced_cluster).
func nullifyOptionalOnlyInRegionConfigs(ctx context.Context, regionsOut types.List) types.List {
	if regionsOut.IsNull() || regionsOut.IsUnknown() {
		return regionsOut
	}
	elemsOut := regionsOut.Elements()
	newElems := make([]TFRegionConfigsModel, len(elemsOut))
	for i := range elemsOut {
		var rcOut TFRegionConfigsModel
		if diags := tfsdk.ValueAs(ctx, elemsOut[i], &rcOut); diags.HasError() {
			return regionsOut
		}
		rcOut.AutoScaling = types.ObjectNull(autoScalingObjType.AttrTypes)
		rcOut.AnalyticsAutoScaling = types.ObjectNull(autoScalingObjType.AttrTypes)
		rcOut.BackingProviderName = types.StringNull()
		rcOut.AnalyticsSpecs = nullifySpecsIfNodeCountZero(ctx, rcOut.AnalyticsSpecs)
		rcOut.ReadOnlySpecs = nullifySpecsIfNodeCountZero(ctx, rcOut.ReadOnlySpecs)
		rcOut.ElectableSpecs = nullifyOptionalOnlyInSpecs(ctx, rcOut.ElectableSpecs)
		newElems[i] = rcOut
	}
	newList, diags := types.ListValueFrom(ctx, regionConfigsObjType, newElems)
	if diags.HasError() {
		return regionsOut
	}
	return newList
}

// nullifySpecsIfNodeCountZero returns null if node_count is 0 (not user-configured),
// otherwise delegates to nullifyOptionalOnlyInSpecs to nullify sub-fields.
func nullifySpecsIfNodeCountZero(ctx context.Context, specs types.Object) types.Object {
	if specs.IsNull() || specs.IsUnknown() {
		return specs
	}
	var model TFSpecsModel
	if diags := tfsdk.ValueAs(ctx, specs, &model); diags.HasError() {
		return specs
	}
	if model.NodeCount.ValueInt64() == 0 {
		return types.ObjectNull(specsObjType.AttrTypes)
	}
	return nullifyOptionalOnlyInSpecs(ctx, specs)
}

// nullifyOptionalOnlyInSpecs nullifies disk_iops and ebs_volume_type within a specs object,
// keeping node_count, instance_size, and disk_size_gb which are user-visible attributes.
func nullifyOptionalOnlyInSpecs(ctx context.Context, specs types.Object) types.Object {
	if specs.IsNull() || specs.IsUnknown() {
		return specs
	}
	var model TFSpecsModel
	if diags := tfsdk.ValueAs(ctx, specs, &model); diags.HasError() {
		return specs
	}
	model.DiskIops = types.Int64Null()
	model.EbsVolumeType = types.StringNull()
	newObj, diags := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, model)
	if diags.HasError() {
		return specs
	}
	return newObj
}

// overrideRegionConfigsWithPrevStateValue preserves null values for Optional-only attributes within each region_config.
func overrideRegionConfigsWithPrevStateValue(ctx context.Context, regionsIn, regionsOut types.List) types.List {
	if regionsIn.IsNull() || regionsIn.IsUnknown() || regionsOut.IsNull() || regionsOut.IsUnknown() {
		return regionsOut
	}

	elemsIn := regionsIn.Elements()
	elemsOut := regionsOut.Elements()
	if len(elemsIn) != len(elemsOut) {
		return regionsOut
	}

	newElems := make([]TFRegionConfigsModel, len(elemsOut))
	for i := range elemsOut {
		var rcIn, rcOut TFRegionConfigsModel
		if diags := tfsdk.ValueAs(ctx, elemsIn[i], &rcIn); diags.HasError() {
			return regionsOut
		}
		if diags := tfsdk.ValueAs(ctx, elemsOut[i], &rcOut); diags.HasError() {
			return regionsOut
		}

		if rcIn.AnalyticsSpecs.IsNull() {
			rcOut.AnalyticsSpecs = types.ObjectNull(specsObjType.AttrTypes)
		} else {
			rcOut.AnalyticsSpecs = overrideSpecsWithPrevStateValue(ctx, rcIn.AnalyticsSpecs, rcOut.AnalyticsSpecs)
		}
		if rcIn.ReadOnlySpecs.IsNull() {
			rcOut.ReadOnlySpecs = types.ObjectNull(specsObjType.AttrTypes)
		} else {
			rcOut.ReadOnlySpecs = overrideSpecsWithPrevStateValue(ctx, rcIn.ReadOnlySpecs, rcOut.ReadOnlySpecs)
		}
		if rcIn.ElectableSpecs.IsNull() {
			rcOut.ElectableSpecs = types.ObjectNull(specsObjType.AttrTypes)
		} else {
			rcOut.ElectableSpecs = overrideSpecsWithPrevStateValue(ctx, rcIn.ElectableSpecs, rcOut.ElectableSpecs)
		}
		if rcIn.AutoScaling.IsNull() {
			rcOut.AutoScaling = types.ObjectNull(autoScalingObjType.AttrTypes)
		} else {
			rcOut.AutoScaling = overrideAutoScalingWithPrevStateValue(ctx, rcIn.AutoScaling, rcOut.AutoScaling)
		}
		if rcIn.AnalyticsAutoScaling.IsNull() {
			rcOut.AnalyticsAutoScaling = types.ObjectNull(autoScalingObjType.AttrTypes)
		} else {
			rcOut.AnalyticsAutoScaling = overrideAutoScalingWithPrevStateValue(ctx, rcIn.AnalyticsAutoScaling, rcOut.AnalyticsAutoScaling)
		}
		if rcIn.BackingProviderName.IsNull() {
			rcOut.BackingProviderName = types.StringNull()
		}

		newElems[i] = rcOut
	}

	newList, diags := types.ListValueFrom(ctx, regionConfigsObjType, newElems)
	if diags.HasError() {
		return regionsOut
	}
	return newList
}

// overrideSpecsWithPrevStateValue preserves null values for Optional-only attributes within a specs block
// (electable_specs, read_only_specs, analytics_specs) when partially configured.
func overrideSpecsWithPrevStateValue(ctx context.Context, specsIn, specsOut types.Object) types.Object {
	if specsIn.IsNull() || specsIn.IsUnknown() || specsOut.IsNull() || specsOut.IsUnknown() {
		return specsOut
	}

	var in, out TFSpecsModel
	if diags := tfsdk.ValueAs(ctx, specsIn, &in); diags.HasError() {
		return specsOut
	}
	if diags := tfsdk.ValueAs(ctx, specsOut, &out); diags.HasError() {
		return specsOut
	}

	if in.DiskIops.IsNull() {
		out.DiskIops = types.Int64Null()
	}
	if in.DiskSizeGb.IsNull() {
		out.DiskSizeGb = types.Float64Null()
	}
	if in.EbsVolumeType.IsNull() {
		out.EbsVolumeType = types.StringNull()
	}
	if in.InstanceSize.IsNull() {
		out.InstanceSize = types.StringNull()
	}
	if in.NodeCount.IsNull() {
		out.NodeCount = types.Int64Null()
	}

	newObj, diags := types.ObjectValueFrom(ctx, specsObjType.AttrTypes, out)
	if diags.HasError() {
		return specsOut
	}
	return newObj
}

// overrideAutoScalingWithPrevStateValue preserves null values for Optional-only attributes within an auto_scaling block
// when partially configured.
func overrideAutoScalingWithPrevStateValue(ctx context.Context, asIn, asOut types.Object) types.Object {
	if asIn.IsNull() || asIn.IsUnknown() || asOut.IsNull() || asOut.IsUnknown() {
		return asOut
	}

	var in, out TFAutoScalingModel
	if diags := tfsdk.ValueAs(ctx, asIn, &in); diags.HasError() {
		return asOut
	}
	if diags := tfsdk.ValueAs(ctx, asOut, &out); diags.HasError() {
		return asOut
	}

	if in.ComputeEnabled.IsNull() {
		out.ComputeEnabled = types.BoolNull()
	}
	if in.ComputeScaleDownEnabled.IsNull() {
		out.ComputeScaleDownEnabled = types.BoolNull()
	}
	if in.ComputeMaxInstanceSize.IsNull() {
		out.ComputeMaxInstanceSize = types.StringNull()
	}
	if in.ComputeMinInstanceSize.IsNull() {
		out.ComputeMinInstanceSize = types.StringNull()
	}
	if in.DiskGBEnabled.IsNull() {
		out.DiskGBEnabled = types.BoolNull()
	}

	newObj, diags := types.ObjectValueFrom(ctx, autoScalingObjType.AttrTypes, out)
	if diags.HasError() {
		return asOut
	}
	return newObj
}
