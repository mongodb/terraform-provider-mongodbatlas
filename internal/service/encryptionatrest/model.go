package encryptionatrest

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFEncryptionAtRestRSModel(ctx context.Context, projectID string, encryptionResp *admin.EncryptionAtRest) *TfEncryptionAtRestRSModel {
	enabledForSearchNodes := false
	if encryptionResp.EnabledForSearchNodes != nil {
		enabledForSearchNodes = encryptionResp.GetEnabledForSearchNodes()
	}
	return &TfEncryptionAtRestRSModel{
		ID:                    types.StringValue(projectID),
		ProjectID:             types.StringValue(projectID),
		AwsKmsConfig:          NewTFAwsKmsConfig(ctx, encryptionResp.AwsKms),
		AzureKeyVaultConfig:   NewTFAzureKeyVaultConfig(ctx, encryptionResp.AzureKeyVault),
		GoogleCloudKmsConfig:  NewTFGcpKmsConfig(ctx, encryptionResp.GoogleCloudKms),
		EnabledForSearchNodes: types.BoolValue(enabledForSearchNodes),
	}
}

func NewTFAwsKmsConfig(ctx context.Context, awsKms *admin.AWSKMSConfiguration) []TFAwsKmsConfigModel {
	if awsKms == nil {
		return []TFAwsKmsConfigModel{}
	}

	return []TFAwsKmsConfigModel{
		*NewTFAwsKmsConfigItem(awsKms),
	}
}

func NewTFAzureKeyVaultConfig(ctx context.Context, az *admin.AzureKeyVault) []TFAzureKeyVaultConfigModel {
	if az == nil {
		return []TFAzureKeyVaultConfigModel{}
	}

	return []TFAzureKeyVaultConfigModel{
		*NewTFAzureKeyVaultConfigItem(az),
	}
}

func NewTFGcpKmsConfig(ctx context.Context, gcpKms *admin.GoogleCloudKMS) []TFGcpKmsConfigModel {
	if gcpKms == nil {
		return []TFGcpKmsConfigModel{}
	}

	return []TFGcpKmsConfigModel{
		*NewTFGcpKmsConfigItem(gcpKms),
	}
}

func NewTFAwsKmsConfigItem(awsKms *admin.AWSKMSConfiguration) *TFAwsKmsConfigModel {
	if awsKms == nil {
		return nil
	}

	return &TFAwsKmsConfigModel{
		Enabled:                  types.BoolPointerValue(awsKms.Enabled),
		CustomerMasterKeyID:      types.StringValue(awsKms.GetCustomerMasterKeyID()),
		Region:                   types.StringValue(awsKms.GetRegion()),
		AccessKeyID:              conversion.StringNullIfEmpty(awsKms.GetAccessKeyID()),
		SecretAccessKey:          conversion.StringNullIfEmpty(awsKms.GetSecretAccessKey()),
		RoleID:                   conversion.StringNullIfEmpty(awsKms.GetRoleId()),
		Valid:                    types.BoolPointerValue(awsKms.Valid),
		RequirePrivateNetworking: types.BoolValue(awsKms.GetRequirePrivateNetworking()),
	}
}

func NewTFAzureKeyVaultConfigItem(az *admin.AzureKeyVault) *TFAzureKeyVaultConfigModel {
	if az == nil {
		return nil
	}

	return &TFAzureKeyVaultConfigModel{
		Enabled:                  types.BoolPointerValue(az.Enabled),
		ClientID:                 types.StringValue(az.GetClientID()),
		AzureEnvironment:         types.StringValue(az.GetAzureEnvironment()),
		SubscriptionID:           types.StringValue(az.GetSubscriptionID()),
		ResourceGroupName:        types.StringValue(az.GetResourceGroupName()),
		KeyVaultName:             types.StringValue(az.GetKeyVaultName()),
		KeyIdentifier:            types.StringValue(az.GetKeyIdentifier()),
		TenantID:                 types.StringValue(az.GetTenantID()),
		Secret:                   conversion.StringNullIfEmpty(az.GetSecret()),
		RequirePrivateNetworking: types.BoolValue(az.GetRequirePrivateNetworking()),
		Valid:                    types.BoolPointerValue(az.Valid),
	}
}

func NewTFGcpKmsConfigItem(gcpKms *admin.GoogleCloudKMS) *TFGcpKmsConfigModel {
	if gcpKms == nil {
		return nil
	}

	return &TFGcpKmsConfigModel{
		Enabled:              types.BoolPointerValue(gcpKms.Enabled),
		KeyVersionResourceID: types.StringValue(gcpKms.GetKeyVersionResourceID()),
		ServiceAccountKey:    conversion.StringNullIfEmpty(gcpKms.GetServiceAccountKey()),
		Valid:                types.BoolPointerValue(gcpKms.Valid),
		RoleID:               types.StringValue(gcpKms.GetRoleId()),
	}
}

func NewAtlasAwsKms(tfAwsKmsConfigSlice []TFAwsKmsConfigModel) *admin.AWSKMSConfiguration {
	if len(tfAwsKmsConfigSlice) == 0 {
		return &admin.AWSKMSConfiguration{}
	}
	v := tfAwsKmsConfigSlice[0]

	awsRegion, _ := conversion.ValRegion(v.Region.ValueString())

	return &admin.AWSKMSConfiguration{
		Enabled:                  v.Enabled.ValueBoolPointer(),
		AccessKeyID:              v.AccessKeyID.ValueStringPointer(),
		SecretAccessKey:          v.SecretAccessKey.ValueStringPointer(),
		CustomerMasterKeyID:      v.CustomerMasterKeyID.ValueStringPointer(),
		Region:                   conversion.StringPtr(awsRegion),
		RoleId:                   v.RoleID.ValueStringPointer(),
		RequirePrivateNetworking: v.RequirePrivateNetworking.ValueBoolPointer(),
	}
}

func NewAtlasGcpKms(tfGcpKmsConfigSlice []TFGcpKmsConfigModel) *admin.GoogleCloudKMS {
	if len(tfGcpKmsConfigSlice) == 0 {
		return &admin.GoogleCloudKMS{}
	}
	v := tfGcpKmsConfigSlice[0]

	return &admin.GoogleCloudKMS{
		Enabled:              v.Enabled.ValueBoolPointer(),
		ServiceAccountKey:    v.ServiceAccountKey.ValueStringPointer(),
		KeyVersionResourceID: v.KeyVersionResourceID.ValueStringPointer(),
		RoleId:               v.RoleID.ValueStringPointer(),
	}
}

func NewAtlasAzureKeyVault(tfAzKeyVaultConfigSlice []TFAzureKeyVaultConfigModel) *admin.AzureKeyVault {
	if len(tfAzKeyVaultConfigSlice) == 0 {
		return &admin.AzureKeyVault{}
	}
	v := tfAzKeyVaultConfigSlice[0]

	return &admin.AzureKeyVault{
		Enabled:                  v.Enabled.ValueBoolPointer(),
		ClientID:                 v.ClientID.ValueStringPointer(),
		AzureEnvironment:         v.AzureEnvironment.ValueStringPointer(),
		SubscriptionID:           v.SubscriptionID.ValueStringPointer(),
		ResourceGroupName:        v.ResourceGroupName.ValueStringPointer(),
		KeyVaultName:             v.KeyVaultName.ValueStringPointer(),
		KeyIdentifier:            v.KeyIdentifier.ValueStringPointer(),
		Secret:                   v.Secret.ValueStringPointer(),
		TenantID:                 v.TenantID.ValueStringPointer(),
		RequirePrivateNetworking: v.RequirePrivateNetworking.ValueBoolPointer(),
	}
}

func NewAtlasEncryptionAtRest(encryptionAtRestPlan, encryptionAtRestState *TfEncryptionAtRestRSModel, atlasEncryptionAtRest *admin.EncryptionAtRest) *admin.EncryptionAtRest {
	if hasAwsKmsConfigChanged(encryptionAtRestPlan.AwsKmsConfig, encryptionAtRestState.AwsKmsConfig) {
		atlasEncryptionAtRest.AwsKms = NewAtlasAwsKms(encryptionAtRestPlan.AwsKmsConfig)
	}
	if hasAzureKeyVaultConfigChanged(encryptionAtRestPlan.AzureKeyVaultConfig, encryptionAtRestState.AzureKeyVaultConfig) {
		atlasEncryptionAtRest.AzureKeyVault = NewAtlasAzureKeyVault(encryptionAtRestPlan.AzureKeyVaultConfig)
	}
	if hasGcpKmsConfigChanged(encryptionAtRestPlan.GoogleCloudKmsConfig, encryptionAtRestState.GoogleCloudKmsConfig) {
		atlasEncryptionAtRest.GoogleCloudKms = NewAtlasGcpKms(encryptionAtRestPlan.GoogleCloudKmsConfig)
	}
	if encryptionAtRestPlan.EnabledForSearchNodes != encryptionAtRestState.EnabledForSearchNodes {
		atlasEncryptionAtRest.EnabledForSearchNodes = encryptionAtRestPlan.EnabledForSearchNodes.ValueBoolPointer()
	}
	return atlasEncryptionAtRest
}
