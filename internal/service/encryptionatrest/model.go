package encryptionatrest

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20240805001/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFEncryptionAtRestRSModel(ctx context.Context, projectID string, encryptionResp *admin.EncryptionAtRest) *TfEncryptionAtRestRSModel {
	return &TfEncryptionAtRestRSModel{
		ID:                   types.StringValue(projectID),
		ProjectID:            types.StringValue(projectID),
		AwsKmsConfig:         NewTFAwsKmsConfig(ctx, encryptionResp.AwsKms),
		AzureKeyVaultConfig:  NewTFAzureKeyVaultConfig(ctx, encryptionResp.AzureKeyVault),
		GoogleCloudKmsConfig: NewTFGcpKmsConfig(ctx, encryptionResp.GoogleCloudKms),
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

func NewAtlasAwsKms(tfAwsKmsConfigSlice []TFAwsKmsConfigModel) *admin.AWSKMSConfiguration {
	if len(tfAwsKmsConfigSlice) == 0 {
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

func NewAtlasGcpKms(tfGcpKmsConfigSlice []TFGcpKmsConfigModel) *admin.GoogleCloudKMS {
	if len(tfGcpKmsConfigSlice) == 0 {
		return &admin.GoogleCloudKMS{}
	}
	v := tfGcpKmsConfigSlice[0]

	return &admin.GoogleCloudKMS{
		Enabled:              v.Enabled.ValueBoolPointer(),
		ServiceAccountKey:    v.ServiceAccountKey.ValueStringPointer(),
		KeyVersionResourceID: v.KeyVersionResourceID.ValueStringPointer(),
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
