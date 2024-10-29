package flexcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241023001/admin"
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
		Tags:                         newTFTags(apiResp.Tags),
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

func NewTFModelDSP(ctx context.Context, projectID string, input *admin.PaginatedFlexClusters20250101) (*TFModelDSP, diag.Diagnostics) {
	diags := &diag.Diagnostics{}
	tfModels := make([]TFModel, len(input.GetResults()))
	for i, item := range input.GetResults() {
		tfModel, diagsLocal := NewTFModel(ctx, &item) // currently complains that item is of type *any
		diags.Append(diagsLocal...)
		tfModels[i] = *tfModel
	}
	if diags.HasError() {
		return nil, *diags
	}
	return &TFModelDSP{
		ProjectId: types.StringValue(projectID),
		Results:   tfModels,
	}, *diags
}

func NewAtlasCreateReq(ctx context.Context, plan *TFModel) (*admin.FlexClusterDescriptionCreate20250101, diag.Diagnostics) {
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
		Tags:                         newResourceTags(ctx, plan.Tags),
	}, nil
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel) (*admin.FlexClusterDescriptionUpdate20241113, diag.Diagnostics) {
	updateRequest := &admin.FlexClusterDescriptionUpdate20241113{
		TerminationProtectionEnabled: plan.TerminationProtectionEnabled.ValueBoolPointer(),
		Tags:                         newResourceTags(ctx, plan.Tags),
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

func newTFTags(tags *[]admin.ResourceTag) basetypes.MapValue {
	if len(*tags) == 0 {
		return types.MapNull(types.StringType)
	}
	typesTags := make(map[string]attr.Value, len(*tags))
	for _, tag := range *tags {
		typesTags[tag.Key] = types.StringValue(tag.Value)
	}
	return types.MapValueMust(types.StringType, typesTags)
}

func newResourceTags(ctx context.Context, tags types.Map) *[]admin.ResourceTag {
	if tags.IsNull() || len(tags.Elements()) == 0 {
		return &[]admin.ResourceTag{}
	}
	elements := make(map[string]types.String, len(tags.Elements()))
	_ = tags.ElementsAs(ctx, &elements, false)
	var tagsAdmin []admin.ResourceTag
	for key, tagValue := range elements {
		tagsAdmin = append(tagsAdmin, admin.ResourceTag{
			Key:   key,
			Value: tagValue.ValueString(),
		})
	}
	return &tagsAdmin
}
