package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20240805005/admin"
)

func AddAdvancedConfig(ctx context.Context, tfModel *TFModel, input *admin.ClusterDescriptionProcessArgs20240805, inputLegacy *admin20240530.ClusterDescriptionProcessArgs, diags *diag.Diagnostics) {
	var advancedConfig TFAdvancedConfigurationModel
	if input != nil && inputLegacy != nil {
		// Using the new API as source of Truth, only use `inputLegacy` for fields not in `input`
		advancedConfig = TFAdvancedConfigurationModel{
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds)),
			DefaultWriteConcern:              types.StringPointerValue(input.DefaultWriteConcern),
			DefaultReadConcern:               types.StringPointerValue(inputLegacy.DefaultReadConcern), // Legacy
			FailIndexKeyTooLong:              types.BoolPointerValue(inputLegacy.FailIndexKeyTooLong),  // TODO: Legacy and not set by the API if version higher than 4.4
			JavascriptEnabled:                types.BoolPointerValue(input.JavascriptEnabled),
			MinimumEnabledTlsProtocol:        types.StringPointerValue(input.MinimumEnabledTlsProtocol),
			NoTableScan:                      types.BoolPointerValue(input.NoTableScan),
			OplogMinRetentionHours:           types.Float64PointerValue(input.OplogMinRetentionHours),
			OplogSizeMb:                      types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.OplogSizeMB)),
			SampleSizeBiconnector:            types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.SampleSizeBIConnector)),
			SampleRefreshIntervalBiconnector: types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.SampleRefreshIntervalBIConnector)),
			TransactionLifetimeLimitSeconds:  types.Int64PointerValue(input.TransactionLifetimeLimitSeconds),
		}
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, AdvancedConfigurationObjType.AttrTypes, advancedConfig)
	diags.Append(diagsLocal...)
	tfModel.AdvancedConfiguration = objType
}
