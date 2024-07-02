package encryptionatrest

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

func NewTfEncryptionAtRestRSModel(ctx context.Context, projectID string, encryptionResp *admin.EncryptionAtRest) *TfEncryptionAtRestRSModel {
	return &TfEncryptionAtRestRSModel{
		ID:                   types.StringValue(projectID),
		ProjectID:            types.StringValue(projectID),
		AwsKmsConfig:         NewTFAwsKmsConfig(ctx, encryptionResp.AwsKms),
		AzureKeyVaultConfig:  NewTFAzureKeyVaultConfig(ctx, encryptionResp.AzureKeyVault),
		GoogleCloudKmsConfig: NewTFGcpKmsConfig(ctx, encryptionResp.GoogleCloudKms),
	}
}

func NewTFAwsKmsConfig(ctx context.Context, awsKms *admin.AWSKMSConfiguration) []TfAwsKmsConfigModel {
	if awsKms == nil {
		return []TfAwsKmsConfigModel{}
	}

	return []TfAwsKmsConfigModel{
		{
			Enabled:             types.BoolPointerValue(awsKms.Enabled),
			CustomerMasterKeyID: types.StringValue(awsKms.GetCustomerMasterKeyID()),
			Region:              types.StringValue(awsKms.GetRegion()),
			AccessKeyID:         conversion.StringNullIfEmpty(awsKms.GetAccessKeyID()),
			SecretAccessKey:     conversion.StringNullIfEmpty(awsKms.GetSecretAccessKey()),
			RoleID:              conversion.StringNullIfEmpty(awsKms.GetRoleId()),
		},
	}
}

func NewTFAzureKeyVaultConfig(ctx context.Context, az *admin.AzureKeyVault) []TfAzureKeyVaultConfigModel {
	if az == nil {
		return []TfAzureKeyVaultConfigModel{}
	}

	return []TfAzureKeyVaultConfigModel{
		{
			Enabled:           types.BoolPointerValue(az.Enabled),
			ClientID:          types.StringValue(az.GetClientID()),
			AzureEnvironment:  types.StringValue(az.GetAzureEnvironment()),
			SubscriptionID:    types.StringValue(az.GetSubscriptionID()),
			ResourceGroupName: types.StringValue(az.GetResourceGroupName()),
			KeyVaultName:      types.StringValue(az.GetKeyVaultName()),
			KeyIdentifier:     types.StringValue(az.GetKeyIdentifier()),
			TenantID:          types.StringValue(az.GetTenantID()),
			Secret:            conversion.StringNullIfEmpty(az.GetSecret()),
		},
	}
}

func NewTFGcpKmsConfig(ctx context.Context, gcpKms *admin.GoogleCloudKMS) []TfGcpKmsConfigModel {
	if gcpKms == nil {
		return []TfGcpKmsConfigModel{}
	}

	return []TfGcpKmsConfigModel{
		{
			Enabled:              types.BoolPointerValue(gcpKms.Enabled),
			KeyVersionResourceID: types.StringValue(gcpKms.GetKeyVersionResourceID()),
			ServiceAccountKey:    conversion.StringNullIfEmpty(gcpKms.GetServiceAccountKey()),
		},
	}
}

func NewAtlasAwsKms(tfAwsKmsConfigSlice []TfAwsKmsConfigModel) *admin.AWSKMSConfiguration {
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

func NewAtlasGcpKms(tfGcpKmsConfigSlice []TfGcpKmsConfigModel) *admin.GoogleCloudKMS {
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

func NewAtlasAzureKeyVault(tfAzKeyVaultConfigSlice []TfAzureKeyVaultConfigModel) *admin.AzureKeyVault {
	if len(tfAzKeyVaultConfigSlice) == 0 {
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
