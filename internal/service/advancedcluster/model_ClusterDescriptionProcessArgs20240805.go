package advancedcluster

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func buildAdvancedConfigObjType(ctx context.Context, input *ProcessArgs, diags *diag.Diagnostics) types.Object {
	var advancedConfig TFAdvancedConfigurationModel
	var customCipherConfigTLS12 *[]string
	var customCipherConfigTLS13 *[]string

	if input.ArgsDefault != nil {
		// Using the new API as source of Truth, only use `inputLegacy` for fields not in `input`
		changeStreamOptionsPreAndPostImagesExpireAfterSeconds := input.ArgsDefault.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds
		if changeStreamOptionsPreAndPostImagesExpireAfterSeconds == nil {
			// special behavior using -1 when it is unset by the user
			changeStreamOptionsPreAndPostImagesExpireAfterSeconds = new(-1)
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
		customCipherConfigTLS12 = input.ArgsDefault.CustomOpensslCipherConfigTls12
		customCipherConfigTLS13 = input.ArgsDefault.CustomOpensslCipherConfigTls13
	}
	advancedConfig.CustomOpensslCipherConfigTls12 = customOpensslCipherConfig(ctx, diags, customCipherConfigTLS12)
	advancedConfig.CustomOpensslCipherConfigTls13 = customOpensslCipherConfig(ctx, diags, customCipherConfigTLS13)

	overrideTLSIfClusterAdvancedConfigPresent(ctx, diags, &advancedConfig, input.ClusterAdvancedConfig)

	objType, diagsLocal := types.ObjectValueFrom(ctx, advancedConfigurationObjType.AttrTypes, advancedConfig)
	diags.Append(diagsLocal...)
	return objType
}

func overrideTLSIfClusterAdvancedConfigPresent(ctx context.Context, diags *diag.Diagnostics, tfAdvConfig *TFAdvancedConfigurationModel, conf *admin.ApiAtlasClusterAdvancedConfiguration) {
	if conf == nil {
		return
	}
	tfAdvConfig.MinimumEnabledTlsProtocol = types.StringValue(conversion.SafeValue(conf.MinimumEnabledTlsProtocol))
	tfAdvConfig.TlsCipherConfigMode = types.StringValue(conversion.SafeValue(conf.TlsCipherConfigMode))
	tfAdvConfig.CustomOpensslCipherConfigTls12 = customOpensslCipherConfig(ctx, diags, conf.CustomOpensslCipherConfigTls12)
	tfAdvConfig.CustomOpensslCipherConfigTls13 = customOpensslCipherConfig(ctx, diags, conf.CustomOpensslCipherConfigTls13)
}

func customOpensslCipherConfig(ctx context.Context, diags *diag.Diagnostics, customOpensslCipherConfig *[]string) types.Set {
	if customOpensslCipherConfig == nil {
		return types.SetNull(types.StringType)
	}
	val, d := types.SetValueFrom(ctx, types.StringType, customOpensslCipherConfig)
	diags.Append(d...)
	return val
}
