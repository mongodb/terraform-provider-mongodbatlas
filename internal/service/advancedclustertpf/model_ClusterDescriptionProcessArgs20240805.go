package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20241023001/admin"
)

func AddAdvancedConfig(ctx context.Context, tfModel *TFModel, input *admin.ClusterDescriptionProcessArgs20240805, diags *diag.Diagnostics) {
	var advancedConfig TFAdvancedConfigurationModel
	if input != nil {
		advancedConfig = TFAdvancedConfigurationModel{
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds)),
			ChunkMigrationConcurrency:                             types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.ChunkMigrationConcurrency)),
			DefaultMaxTimeMs:                                      types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.DefaultMaxTimeMS)),
			DefaultWriteConcern:                                   types.StringPointerValue(input.DefaultWriteConcern),
			DefaultReadConcern:                                    types.StringNull(), // TODO: static
			FailIndexKeyTooLong:                                   types.BoolNull(),   // TODO: static,
			JavascriptEnabled:                                     types.BoolPointerValue(input.JavascriptEnabled),
			MinimumEnabledTlsProtocol:                             types.StringPointerValue(input.MinimumEnabledTlsProtocol),
			NoTableScan:                                           types.BoolPointerValue(input.NoTableScan),
			OplogMinRetentionHours:                                types.Float64PointerValue(input.OplogMinRetentionHours),
			OplogSizeMb:                                           types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.OplogSizeMB)),
			QueryStatsLogVerbosity:                                types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.QueryStatsLogVerbosity)),
			SampleSizeBiconnector:                                 types.Int64Null(), // TODO: static
			SampleRefreshIntervalBiconnector:                      types.Int64Null(), // TODO: static
			TransactionLifetimeLimitSeconds:                       types.Int64PointerValue(input.TransactionLifetimeLimitSeconds),
		}
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, AdvancedConfigurationObjType.AttrTypes, advancedConfig)
	diags.Append(diagsLocal...)
	tfModel.AdvancedConfiguration = objType
}
