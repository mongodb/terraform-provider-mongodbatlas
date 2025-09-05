package flexcluster

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFModel(ctx context.Context, apiResp *admin.FlexClusterDescription20241113) (*TFModel, diag.Diagnostics) {
	connectionStrings, diags := ConvertConnectionStringsToTF(ctx, apiResp.ConnectionStrings)
	if diags.HasError() {
		return nil, diags
	}
	backupSettings, diags := ConvertBackupSettingsToTF(ctx, apiResp.BackupSettings)
	if diags.HasError() {
		return nil, diags
	}
	providerSettings, diags := ConvertProviderSettingsToTF(ctx, apiResp.ProviderSettings)
	if diags.HasError() {
		return nil, diags
	}
	return &TFModel{
		ProviderSettings:             *providerSettings,
		ConnectionStrings:            *connectionStrings,
		Tags:                         conversion.NewTFTags(apiResp.GetTags()),
		CreateDate:                   types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.CreateDate)),
		ProjectId:                    types.StringPointerValue(apiResp.GroupId),
		Id:                           types.StringPointerValue(apiResp.Id),
		MongoDbversion:               types.StringPointerValue(apiResp.MongoDBVersion),
		Name:                         types.StringPointerValue(apiResp.Name),
		ClusterType:                  types.StringPointerValue(apiResp.ClusterType),
		StateName:                    types.StringPointerValue(apiResp.StateName),
		VersionReleaseSystem:         types.StringPointerValue(apiResp.VersionReleaseSystem),
		BackupSettings:               *backupSettings,
		TerminationProtectionEnabled: types.BoolPointerValue(apiResp.TerminationProtectionEnabled),
	}, nil
}

func NewTFModelDSP(ctx context.Context, projectID string, input []admin.FlexClusterDescription20241113) (*TFModelDSP, diag.Diagnostics) {
	diags := &diag.Diagnostics{}
	tfModels := make([]TFModel, len(input))
	for i := range input {
		item := &input[i]
		tfModel, diagsLocal := NewTFModel(ctx, item)
		diags.Append(diagsLocal...)
		if tfModel != nil {
			tfModels[i] = *tfModel
		}
	}
	if diags.HasError() {
		return nil, *diags
	}
	return &TFModelDSP{
		ProjectId: types.StringValue(projectID),
		Results:   tfModels,
	}, *diags
}

func NewAtlasCreateReq(ctx context.Context, plan *TFModel) (*admin.FlexClusterDescriptionCreate20241113, diag.Diagnostics) {
	providerSettings := &TFProviderSettings{}
	if diags := plan.ProviderSettings.As(ctx, providerSettings, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	return &admin.FlexClusterDescriptionCreate20241113{
		Name: plan.Name.ValueString(),
		ProviderSettings: admin.FlexProviderSettingsCreate20241113{
			BackingProviderName: providerSettings.BackingProviderName.ValueString(),
			RegionName:          providerSettings.RegionName.ValueString(),
		},
		TerminationProtectionEnabled: plan.TerminationProtectionEnabled.ValueBoolPointer(),
		Tags:                         conversion.NewResourceTags(ctx, plan.Tags),
	}, nil
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel) (*admin.FlexClusterDescriptionUpdate20241113, diag.Diagnostics) {
	updateRequest := &admin.FlexClusterDescriptionUpdate20241113{
		TerminationProtectionEnabled: plan.TerminationProtectionEnabled.ValueBoolPointer(),
		Tags:                         conversion.NewResourceTags(ctx, plan.Tags),
	}

	return updateRequest, nil
}

func ConvertBackupSettingsToTF(ctx context.Context, backupSettings *admin.FlexBackupSettings20241113) (*types.Object, diag.Diagnostics) {
	if backupSettings == nil {
		backupSettingsTF := types.ObjectNull(BackupSettingsType.AttributeTypes())
		return &backupSettingsTF, nil
	}

	backupSettingsTF := &TFBackupSettings{
		Enabled: types.BoolPointerValue(backupSettings.Enabled),
	}
	backupSettingsObject, diags := types.ObjectValueFrom(ctx, BackupSettingsType.AttributeTypes(), backupSettingsTF)
	if diags.HasError() {
		return nil, diags
	}
	return &backupSettingsObject, nil
}

func ConvertConnectionStringsToTF(ctx context.Context, connectionStrings *admin.FlexConnectionStrings20241113) (*types.Object, diag.Diagnostics) {
	if connectionStrings == nil {
		connectionStringsTF := types.ObjectNull(ConnectionStringsType.AttributeTypes())
		return &connectionStringsTF, nil
	}

	connectionStringsTF := &TFConnectionStrings{
		Standard:    types.StringPointerValue(connectionStrings.Standard),
		StandardSrv: types.StringPointerValue(connectionStrings.StandardSrv),
	}
	connectionStringsObject, diags := types.ObjectValueFrom(ctx, ConnectionStringsType.AttributeTypes(), connectionStringsTF)
	if diags.HasError() {
		return nil, diags
	}
	return &connectionStringsObject, nil
}

func ConvertProviderSettingsToTF(ctx context.Context, providerSettings admin.FlexProviderSettings20241113) (*types.Object, diag.Diagnostics) {
	providerSettingsTF := &TFProviderSettings{
		ProviderName:        types.StringPointerValue(providerSettings.ProviderName),
		RegionName:          types.StringPointerValue(providerSettings.RegionName),
		BackingProviderName: types.StringPointerValue(providerSettings.BackingProviderName),
		DiskSizeGb:          types.Float64PointerValue(providerSettings.DiskSizeGB),
	}
	providerSettingsObject, diags := types.ObjectValueFrom(ctx, ProviderSettingsType.AttributeTypes(), providerSettingsTF)
	if diags.HasError() {
		return nil, diags
	}
	return &providerSettingsObject, nil
}

func FlattenFlexConnectionStrings(str *admin.FlexConnectionStrings20241113) []map[string]any {
	return []map[string]any{
		{
			"standard":     str.GetStandard(),
			"standard_srv": str.GetStandardSrv(),
		},
	}
}

func FlattenFlexProviderSettingsIntoReplicationSpecs(providerSettings admin.FlexProviderSettings20241113, priority *int, zoneName *string) []map[string]any {
	tfMaps := []map[string]any{{}}
	tfMaps[0]["num_shards"] = 1 // default value
	tfMaps[0]["zone_name"] = zoneName
	tfMaps[0]["region_configs"] = []map[string]any{
		{
			"provider_name":         providerSettings.GetProviderName(),
			"backing_provider_name": providerSettings.GetBackingProviderName(),
			"region_name":           providerSettings.GetRegionName(),
			"priority":              priority, // no-op for flex clusters, value from config is set in the state to avoid plan changes
		},
	}
	return tfMaps
}

func FlattenFlexClustersToAdvancedClusters(flexClusters *[]admin.FlexClusterDescription20241113) []map[string]any {
	if flexClusters == nil {
		return nil
	}
	results := make([]map[string]any, len(*flexClusters))
	for i := range *flexClusters {
		flexCluster := &(*flexClusters)[i]
		results[i] = map[string]any{
			"cluster_type":                   flexCluster.GetClusterType(),
			"backup_enabled":                 flexCluster.BackupSettings.GetEnabled(),
			"connection_strings":             FlattenFlexConnectionStrings(flexCluster.ConnectionStrings),
			"create_date":                    conversion.TimePtrToStringPtr(flexCluster.CreateDate),
			"mongo_db_version":               flexCluster.GetMongoDBVersion(),
			"replication_specs":              FlattenFlexProviderSettingsIntoReplicationSpecs(flexCluster.ProviderSettings, nil, nil),
			"name":                           flexCluster.GetName(),
			"state_name":                     flexCluster.GetStateName(),
			"tags":                           conversion.FlattenTags(flexCluster.GetTags()),
			"termination_protection_enabled": flexCluster.GetTerminationProtectionEnabled(),
			"version_release_system":         flexCluster.GetVersionReleaseSystem(),
		}
	}
	return results
}
