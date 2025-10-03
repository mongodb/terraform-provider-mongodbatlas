package advancedcluster

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312008/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func AddAdvancedConfig(ctx context.Context, tfModel *TFModel, input *ProcessArgs, diags *diag.Diagnostics) {
	var advancedConfig TFAdvancedConfigurationModel
	var customCipherConfig *[]string

	if input.ArgsDefault != nil {
		// Using the new API as source of Truth, only use `inputLegacy` for fields not in `input`
		changeStreamOptionsPreAndPostImagesExpireAfterSeconds := input.ArgsDefault.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds
		if changeStreamOptionsPreAndPostImagesExpireAfterSeconds == nil {
			// special behavior using -1 when it is unset by the user
			changeStreamOptionsPreAndPostImagesExpireAfterSeconds = conversion.Pointer(-1)
		}

		advancedConfig = TFAdvancedConfigurationModel{
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: types.Int64PointerValue(conversion.IntPtrToInt64Ptr(changeStreamOptionsPreAndPostImagesExpireAfterSeconds)),
			DefaultWriteConcern:              types.StringValue(conversion.SafeValue(input.ArgsDefault.DefaultWriteConcern)),
			JavascriptEnabled:                types.BoolValue(conversion.SafeValue(input.ArgsDefault.JavascriptEnabled)),
			NoTableScan:                      types.BoolValue(conversion.SafeValue(input.ArgsDefault.NoTableScan)),
			OplogMinRetentionHours:           types.Float64Value(conversion.SafeValue(input.ArgsDefault.OplogMinRetentionHours)),
			OplogSizeMb:                      types.Int64Value(conversion.SafeValue(conversion.IntPtrToInt64Ptr(input.ArgsDefault.OplogSizeMB))),
			SampleSizeBiconnector:            types.Int64Value(conversion.SafeValue(conversion.IntPtrToInt64Ptr(input.ArgsDefault.SampleSizeBIConnector))),
			SampleRefreshIntervalBiconnector: types.Int64Value(conversion.SafeValue(conversion.IntPtrToInt64Ptr(input.ArgsDefault.SampleRefreshIntervalBIConnector))),
			TransactionLifetimeLimitSeconds:  types.Int64Value(conversion.SafeValue(input.ArgsDefault.TransactionLifetimeLimitSeconds)),
			DefaultMaxTimeMS:                 types.Int64PointerValue(conversion.IntPtrToInt64Ptr(input.ArgsDefault.DefaultMaxTimeMS)),
			MinimumEnabledTlsProtocol:        types.StringValue(conversion.SafeValue(input.ArgsDefault.MinimumEnabledTlsProtocol)),
			TlsCipherConfigMode:              types.StringValue(conversion.SafeValue(input.ArgsDefault.TlsCipherConfigMode)),
		}
		customCipherConfig = input.ArgsDefault.CustomOpensslCipherConfigTls12
	}
	advancedConfig.CustomOpensslCipherConfigTls12 = customOpensslCipherConfigTLS12(ctx, diags, customCipherConfig)

	overrideTLSIfClusterAdvancedConfigPresent(ctx, diags, &advancedConfig, input.ClusterAdvancedConfig)

	objType, diagsLocal := types.ObjectValueFrom(ctx, AdvancedConfigurationObjType.AttrTypes, advancedConfig)
	diags.Append(diagsLocal...)
	tfModel.AdvancedConfiguration = objType
}

func overrideTLSIfClusterAdvancedConfigPresent(ctx context.Context, diags *diag.Diagnostics, tfAdvConfig *TFAdvancedConfigurationModel, conf *admin.ApiAtlasClusterAdvancedConfiguration) {
	if conf == nil {
		return
	}
	tfAdvConfig.MinimumEnabledTlsProtocol = types.StringValue(conversion.SafeValue(conf.MinimumEnabledTlsProtocol))
	tfAdvConfig.TlsCipherConfigMode = types.StringValue(conversion.SafeValue(conf.TlsCipherConfigMode))
	tfAdvConfig.CustomOpensslCipherConfigTls12 = customOpensslCipherConfigTLS12(ctx, diags, conf.CustomOpensslCipherConfigTls12)
}

func customOpensslCipherConfigTLS12(ctx context.Context, diags *diag.Diagnostics, customOpensslCipherConfigTLS12 *[]string) types.Set {
	if customOpensslCipherConfigTLS12 == nil {
		return types.SetNull(types.StringType)
	}
	val, d := types.SetValueFrom(ctx, types.StringType, customOpensslCipherConfigTLS12)
	diags.Append(d...)
	return val
}
