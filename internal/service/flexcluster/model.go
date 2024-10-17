//nolint:gocritic
package flexcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
)

// TODO: `ctx` parameter and `diags` return value can be removed if tf schema has no complex data types (e.g., schema.ListAttribute, schema.SetAttribute)
func NewTFModel(ctx context.Context, apiResp *admin.FlexClusterDescription20250101) (*TFModel, diag.Diagnostics) {
	diags := &diag.Diagnostics{}
	providerSettings := newProviderSettings()
	connectionStrings := newConnectionStrings()
	tags := newTags()
	backupSettings := newBackupSettings()
	if diags.HasError() {
		return nil, nil //TODO
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

// TODO: If SDK defined different models for create and update separate functions will need to be defined.
// TODO: `ctx` parameter and `diags` in return value can be removed if tf schema has no complex data types (e.g., schema.ListAttribute, schema.SetAttribute)
func NewAtlasCreateReq(ctx context.Context, plan *TFModel) (*admin.FlexClusterDescriptionCreate20250101, diag.Diagnostics) {
	// var tfList []complexArgumentData
	// resp.Diagnostics.Append(plan.ComplexArgument.ElementsAs(ctx, &tfList, false)...)
	// if resp.Diagnostics.HasError() {
	// 	return nil, diagnostics
	// }
	return &admin.FlexClusterDescriptionCreate20250101{}, nil
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel) (*admin.FlexClusterDescription20250101, diag.Diagnostics) {
	return &admin.FlexClusterDescription20250101{}, nil
}

func newProviderSettings() TFProviderSettings {
	return TFProviderSettings{}
}

func newConnectionStrings() TFConnectionStrings {
	return TFConnectionStrings{}
}

func newTags() basetypes.MapValue {
	return basetypes.MapValue{}
}

func newBackupSettings() TFBackupSettings {
	return TFBackupSettings{}
}
