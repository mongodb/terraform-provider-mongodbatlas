package encryptionatrest

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

func NewTfEncryptionAtRestRSModel(ctx context.Context, projectID string, encryptionResp *admin.EncryptionAtRest, plan *TfEncryptionAtRestRSModel) *TfEncryptionAtRestRSModel {
	return &TfEncryptionAtRestRSModel{
		ID:                   types.StringValue(projectID),
		ProjectID:            types.StringValue(projectID),
		AwsKmsConfig:         NewTFAwsKmsConfig(ctx, encryptionResp.AwsKms, plan.AwsKmsConfig),
		AzureKeyVaultConfig:  NewTFAzureKeyVaultConfig(ctx, encryptionResp.AzureKeyVault, plan.AzureKeyVaultConfig),
		GoogleCloudKmsConfig: NewTFGcpKmsConfig(ctx, encryptionResp.GoogleCloudKms, plan.GoogleCloudKmsConfig),
	}
}

func NewTFAwsKmsConfig(ctx context.Context, awsKms *admin.AWSKMSConfiguration, currStateSlice []TfAwsKmsConfigModel) []TfAwsKmsConfigModel {
	if awsKms == nil {
		return []TfAwsKmsConfigModel{}
	}
	newState := TfAwsKmsConfigModel{}

	newState.Enabled = types.BoolPointerValue(awsKms.Enabled)
	newState.CustomerMasterKeyID = types.StringValue(awsKms.GetCustomerMasterKeyID())
	newState.Region = types.StringValue(awsKms.GetRegion())
	newState.AccessKeyID = conversion.StringNullIfEmpty(awsKms.GetAccessKeyID())
	newState.SecretAccessKey = conversion.StringNullIfEmpty(awsKms.GetSecretAccessKey())
	newState.RoleID = conversion.StringNullIfEmpty(awsKms.GetRoleId())

	return []TfAwsKmsConfigModel{newState}
}

func NewTFAzureKeyVaultConfig(ctx context.Context, az *admin.AzureKeyVault, currStateSlice []TfAzureKeyVaultConfigModel) []TfAzureKeyVaultConfigModel {
	if az == nil {
		return []TfAzureKeyVaultConfigModel{}
	}
	newState := TfAzureKeyVaultConfigModel{}

	newState.Enabled = types.BoolPointerValue(az.Enabled)
	newState.ClientID = types.StringValue(az.GetClientID())
	newState.AzureEnvironment = types.StringValue(az.GetAzureEnvironment())
	newState.SubscriptionID = types.StringValue(az.GetSubscriptionID())
	newState.ResourceGroupName = types.StringValue(az.GetResourceGroupName())
	newState.KeyVaultName = types.StringValue(az.GetKeyVaultName())
	newState.KeyIdentifier = types.StringValue(az.GetKeyIdentifier())
	newState.TenantID = types.StringValue(az.GetTenantID())
	newState.Secret = conversion.StringNullIfEmpty(az.GetSecret())

	return []TfAzureKeyVaultConfigModel{newState}
}

func NewTFGcpKmsConfig(ctx context.Context, gcpKms *admin.GoogleCloudKMS, currStateSlice []TfGcpKmsConfigModel) []TfGcpKmsConfigModel {
	if gcpKms == nil {
		return []TfGcpKmsConfigModel{}
	}
	newState := TfGcpKmsConfigModel{}

	newState.Enabled = types.BoolPointerValue(gcpKms.Enabled)
	newState.KeyVersionResourceID = types.StringValue(gcpKms.GetKeyVersionResourceID())
	newState.ServiceAccountKey = conversion.StringNullIfEmpty(gcpKms.GetServiceAccountKey())

	return []TfGcpKmsConfigModel{newState}
}

func NewAtlasAwsKms(tfAwsKmsConfigSlice []TfAwsKmsConfigModel) *admin.AWSKMSConfiguration {
	if tfAwsKmsConfigSlice == nil || len(tfAwsKmsConfigSlice) < 1 {
		return &admin.AWSKMSConfiguration{}
	}
	v := tfAwsKmsConfigSlice[0]

	awsRegion, _ := conversion.ValRegion(v.Region.ValueString())

	return &admin.AWSKMSConfiguration{
		Enabled:             v.Enabled.ValueBoolPointer(),
		AccessKeyID:         v.AccessKeyID.ValueStringPointer(),
		SecretAccessKey:     v.SecretAccessKey.ValueStringPointer(),
		CustomerMasterKeyID: v.CustomerMasterKeyID.ValueStringPointer(),
		Region:              conversion.StringPtr(awsRegion),
		RoleId:              v.RoleID.ValueStringPointer(),
	}
}

func NewAtlasGcpKms(tfGcpKmsConfigSlice []TfGcpKmsConfigModel) *admin.GoogleCloudKMS {
	if tfGcpKmsConfigSlice == nil || len(tfGcpKmsConfigSlice) < 1 {
		return &admin.GoogleCloudKMS{}
	}
	v := tfGcpKmsConfigSlice[0]

	return &admin.GoogleCloudKMS{
		Enabled:              v.Enabled.ValueBoolPointer(),
		ServiceAccountKey:    v.ServiceAccountKey.ValueStringPointer(),
		KeyVersionResourceID: v.KeyVersionResourceID.ValueStringPointer(),
	}
}

func NewAtlasAzureKeyVault(tfAzKeyVaultConfigSlice []TfAzureKeyVaultConfigModel) *admin.AzureKeyVault {
	if tfAzKeyVaultConfigSlice == nil || len(tfAzKeyVaultConfigSlice) < 1 {
		return &admin.AzureKeyVault{}
	}
	v := tfAzKeyVaultConfigSlice[0]

	return &admin.AzureKeyVault{
		Enabled:           v.Enabled.ValueBoolPointer(),
		ClientID:          v.ClientID.ValueStringPointer(),
		AzureEnvironment:  v.AzureEnvironment.ValueStringPointer(),
		SubscriptionID:    v.SubscriptionID.ValueStringPointer(),
		ResourceGroupName: v.ResourceGroupName.ValueStringPointer(),
		KeyVaultName:      v.KeyVaultName.ValueStringPointer(),
		KeyIdentifier:     v.KeyIdentifier.ValueStringPointer(),
		Secret:            v.Secret.ValueStringPointer(),
		TenantID:          v.TenantID.ValueStringPointer(),
	}
}
