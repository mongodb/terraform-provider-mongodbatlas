//nolint:gocritic
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
	diags := &diag.Diagnostics{}
	providerSettings := newProviderSettings(apiResp.ProviderSettings)
	connectionStrings := newConnectionStrings(apiResp.ConnectionStrings)
	tags := newTags(apiResp.Tags)
	backupSettings := newBackupSettings(apiResp.BackupSettings)
	if diags.HasError() {
		return nil, *diags
	}
	return &TFModel{
		ProviderSettings:             providerSettings,
		ConnectionStrings:            connectionStrings,
		Tags:                         tags,
		CreateDate:                   types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.CreateDate)),
		ProjectId:                    types.StringPointerValue(apiResp.GroupId),
		Id:                           types.StringPointerValue(apiResp.Id),
		MongoDbversion:               types.StringPointerValue(apiResp.MongoDBVersion),
		Name:                         types.StringPointerValue(apiResp.Name),
		ClusterType:                  types.StringPointerValue(apiResp.ClusterType),
		StateName:                    types.StringPointerValue(apiResp.StateName),
		VersionReleaseSystem:         types.StringPointerValue(apiResp.VersionReleaseSystem),
		BackupSettings:               backupSettings,
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
	return &admin.FlexClusterDescription20250101{
		BackupSettings: &admin.FlexBackupSettings20250101{
			Enabled: plan.BackupSettings.Enabled.ValueBoolPointer(),
		},
		ClusterType: plan.ClusterType.ValueStringPointer(),
		ConnectionStrings: &admin.FlexConnectionStrings20250101{
			Standard:    plan.ConnectionStrings.Standard.ValueStringPointer(),
			StandardSrv: plan.ConnectionStrings.StandardSrv.ValueStringPointer(),
		},
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
	}, nil
}

func newProviderSettings(providerSettings admin.FlexProviderSettings20250101) TFProviderSettings {
	return TFProviderSettings{
		ProviderName:        types.StringPointerValue(providerSettings.ProviderName),
		RegionName:          types.StringPointerValue(providerSettings.RegionName),
		BackingProviderName: types.StringPointerValue(providerSettings.BackingProviderName),
		DiskSizeGb:          types.Float64PointerValue(providerSettings.DiskSizeGB),
	}
}

func newConnectionStrings(connectionStrings *admin.FlexConnectionStrings20250101) TFConnectionStrings {
	if connectionStrings == nil {
		return TFConnectionStrings{}
	}
	return TFConnectionStrings{
		Standard:    types.StringPointerValue(connectionStrings.Standard),
		StandardSrv: types.StringPointerValue(connectionStrings.StandardSrv),
	}
}

func newTags(tags *[]admin.ResourceTag) basetypes.MapValue {
	if tags == nil || len(*tags) == 0 {
		return basetypes.MapValue{}
	}
	typesTags := make(map[string]attr.Value, len(*tags))
	for _, tag := range *tags {
		typesTags[tag.Key] = types.StringValue(tag.Value)
	}
	return types.MapValueMust(types.StringType, typesTags)
}

func newBackupSettings(backupSettings *admin.FlexBackupSettings20250101) TFBackupSettings {
	if backupSettings == nil {
		return TFBackupSettings{}
	}
	return TFBackupSettings{
		Enabled: types.BoolPointerValue(backupSettings.Enabled),
	}
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
