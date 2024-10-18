package flexcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
)

func NewTFModel(ctx context.Context, apiResp *admin.FlexClusterDescription20250101) (*TFModel, diag.Diagnostics) {
	connectionStrings, diags := ConvertConnectionStringsToTF(ctx, apiResp.ConnectionStrings)
	if diags.HasError() {
		return nil, diags
	}
	backupSettings, diags := ConvertBackupSettingsToTF(ctx, apiResp.BackupSettings)
	if diags.HasError() {
		return nil, diags
	}
	return &TFModel{
		ProviderSettings:             newProviderSettings(apiResp.ProviderSettings),
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

func NewAtlasCreateReq(ctx context.Context, plan *TFModel) (*admin.FlexClusterDescriptionCreate20250101, diag.Diagnostics) {
	return &admin.FlexClusterDescriptionCreate20250101{
		Name: plan.Name.ValueString(),
		ProviderSettings: admin.FlexProviderSettingsCreate20250101{
			BackingProviderName: plan.ProviderSettings.BackingProviderName.ValueString(),
			RegionName:          plan.ProviderSettings.RegionName.ValueString(),
		},
		TerminationProtectionEnabled: plan.TerminationProtectionEnabled.ValueBoolPointer(),
		Tags:                         newResourceTags(ctx, plan.Tags),
	}, nil
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel) (*admin.FlexClusterDescription20250101, diag.Diagnostics) {
	createDateAsTime, _ := conversion.StringToTime(plan.CreateDate.ValueString())

	updateRequest := &admin.FlexClusterDescription20250101{
		ClusterType:    plan.ClusterType.ValueStringPointer(),
		CreateDate:     &createDateAsTime,
		GroupId:        plan.ProjectId.ValueStringPointer(),
		Id:             plan.Id.ValueStringPointer(),
		MongoDBVersion: plan.MongoDbversion.ValueStringPointer(),
		Name:           plan.Name.ValueStringPointer(),
		ProviderSettings: admin.FlexProviderSettings20250101{
			BackingProviderName: plan.ProviderSettings.BackingProviderName.ValueStringPointer(),
			DiskSizeGB:          plan.ProviderSettings.DiskSizeGb.ValueFloat64Pointer(),
			ProviderName:        plan.ProviderSettings.ProviderName.ValueStringPointer(),
			RegionName:          plan.ProviderSettings.RegionName.ValueStringPointer(),
		},
		StateName:                    plan.StateName.ValueStringPointer(),
		TerminationProtectionEnabled: plan.TerminationProtectionEnabled.ValueBoolPointer(),
		Tags:                         newResourceTags(ctx, plan.Tags),
		VersionReleaseSystem:         plan.VersionReleaseSystem.ValueStringPointer(),
	}

	if !plan.BackupSettings.IsNull() && !plan.BackupSettings.IsUnknown() {
		backupSettings := &TFBackupSettings{}
		if diags := plan.BackupSettings.As(ctx, backupSettings, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		updateRequest.BackupSettings = &admin.FlexBackupSettings20250101{
			Enabled: backupSettings.Enabled.ValueBoolPointer(),
		}
	}

	if !plan.ConnectionStrings.IsNull() && !plan.ConnectionStrings.IsUnknown() {
		connectionStrings := &TFConnectionStrings{}
		if diags := plan.ConnectionStrings.As(ctx, connectionStrings, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		updateRequest.ConnectionStrings = &admin.FlexConnectionStrings20250101{
			Standard:    connectionStrings.Standard.ValueStringPointer(),
			StandardSrv: connectionStrings.StandardSrv.ValueStringPointer(),
		}
	}

	return updateRequest, nil
}

func newProviderSettings(providerSettings admin.FlexProviderSettings20250101) TFProviderSettings {
	return TFProviderSettings{
		ProviderName:        types.StringPointerValue(providerSettings.ProviderName),
		RegionName:          types.StringPointerValue(providerSettings.RegionName),
		BackingProviderName: types.StringPointerValue(providerSettings.BackingProviderName),
		DiskSizeGb:          types.Float64PointerValue(providerSettings.DiskSizeGB),
	}
}

func ConvertBackupSettingsToTF(ctx context.Context, backupSettings *admin.FlexBackupSettings20250101) (*types.Object, diag.Diagnostics) {
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

func ConvertConnectionStringsToTF(ctx context.Context, connectionStrings *admin.FlexConnectionStrings20250101) (*types.Object, diag.Diagnostics) {
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
