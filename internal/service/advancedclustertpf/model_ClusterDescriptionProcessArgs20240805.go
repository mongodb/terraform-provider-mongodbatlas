package advancedclustertpf

import (
	"context"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	// "go.mongodb.org/atlas-sdk/v20241113003/admin"
	"github.com/mongodb/atlas-sdk-go/admin" // TODO: replace usage with latest once cipher config changes are in prod

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func AddAdvancedConfig(ctx context.Context, tfModel *TFModel, input *admin.ClusterDescriptionProcessArgs20240805, inputLegacy *admin20240530.ClusterDescriptionProcessArgs, diags *diag.Diagnostics) {
	var advancedConfig TFAdvancedConfigurationModel
	if input != nil && inputLegacy != nil {
		// Using the new API as source of Truth, only use `inputLegacy` for fields not in `input`
		changeStreamOptionsPreAndPostImagesExpireAfterSeconds := input.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds
		if changeStreamOptionsPreAndPostImagesExpireAfterSeconds == nil {
			// special behavior using -1 when it is unset by the user
			changeStreamOptionsPreAndPostImagesExpireAfterSeconds = conversion.Pointer(-1)
		}
		// When MongoDBMajorVersion is not 4.4 or lower, the API response for fail_index_key_too_long will always be null, to ensure no consistency issues, we need to match the config
		failIndexKeyTooLong := inputLegacy.GetFailIndexKeyTooLong()
		if tfModel != nil {
			stateConfig := tfModel.AdvancedConfiguration
			stateConfigSDK := NewAtlasReqAdvancedConfigurationLegacy(ctx, &stateConfig, diags)
			if diags.HasError() {
				return
			}
			if stateConfigSDK != nil && stateConfigSDK.GetFailIndexKeyTooLong() != failIndexKeyTooLong {
				failIndexKeyTooLong = stateConfigSDK.GetFailIndexKeyTooLong()
			}
		}
		advancedConfig = TFAdvancedConfigurationModel{
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: types.Int64PointerValue(conversion.IntPtrToInt64Ptr(changeStreamOptionsPreAndPostImagesExpireAfterSeconds)),
			DefaultWriteConcern:              types.StringValue(conversion.SafeValue(input.DefaultWriteConcern)),
			DefaultReadConcern:               types.StringValue(conversion.SafeValue(inputLegacy.DefaultReadConcern)),
			FailIndexKeyTooLong:              types.BoolValue(failIndexKeyTooLong),
			JavascriptEnabled:                types.BoolValue(conversion.SafeValue(input.JavascriptEnabled)),
			MinimumEnabledTlsProtocol:        types.StringValue(conversion.SafeValue(input.MinimumEnabledTlsProtocol)),
			NoTableScan:                      types.BoolValue(conversion.SafeValue(input.NoTableScan)),
			OplogMinRetentionHours:           types.Float64Value(conversion.SafeValue(input.OplogMinRetentionHours)),
			OplogSizeMb:                      types.Int64Value(conversion.SafeValue(conversion.IntPtrToInt64Ptr(input.OplogSizeMB))),
			SampleSizeBiconnector:            types.Int64Value(conversion.SafeValue(conversion.IntPtrToInt64Ptr(input.SampleSizeBIConnector))),
			SampleRefreshIntervalBiconnector: types.Int64Value(conversion.SafeValue(conversion.IntPtrToInt64Ptr(input.SampleRefreshIntervalBIConnector))),
			TransactionLifetimeLimitSeconds:  types.Int64Value(conversion.SafeValue(input.TransactionLifetimeLimitSeconds)),
		}
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, AdvancedConfigurationObjType.AttrTypes, advancedConfig)
	diags.Append(diagsLocal...)
	tfModel.AdvancedConfiguration = objType
}
