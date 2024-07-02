package encryptionatrest_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/encryptionatrest"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

var (
	projectID            = "projectID"
	enabled              = true
	customerMasterKeyID  = "CustomerMasterKeyID"
	region               = "Region"
	accessKeyID          = "AccessKeyID"
	secretAccessKey      = "SecretAccessKey"
	roleID               = "RoleID"
	clientID             = "clientID"
	azureEnvironment     = "AzureEnvironment"
	subscriptionID       = "SubscriptionID"
	resourceGroupName    = "ResourceGroupName"
	keyVaultName         = "KeyVaultName"
	keyIdentifier        = "KeyIdentifier"
	tenantID             = "TenantID"
	secret               = "Secret"
	keyVersionResourceID = "KeyVersionResourceID"
	serviceAccountKey    = "ServiceAccountKey"
	AWSKMSConfiguration  = &admin.AWSKMSConfiguration{
		Enabled:             &enabled,
		CustomerMasterKeyID: &customerMasterKeyID,
		Region:              &region,
		AccessKeyID:         &accessKeyID,
		SecretAccessKey:     &secretAccessKey,
		RoleId:              &roleID,
	}
	TfAwsKmsConfigModel = encryptionatrest.TfAwsKmsConfigModel{
		Enabled:             types.BoolValue(enabled),
		CustomerMasterKeyID: types.StringValue(customerMasterKeyID),
		Region:              types.StringValue(region),
		AccessKeyID:         types.StringValue(accessKeyID),
		SecretAccessKey:     types.StringValue(secretAccessKey),
		RoleID:              types.StringValue(roleID),
	}
	AzureKeyVault = &admin.AzureKeyVault{
		Enabled:           &enabled,
		ClientID:          &clientID,
		AzureEnvironment:  &azureEnvironment,
		SubscriptionID:    &subscriptionID,
		ResourceGroupName: &resourceGroupName,
		KeyVaultName:      &keyVaultName,
		KeyIdentifier:     &keyIdentifier,
		TenantID:          &tenantID,
		Secret:            &secret,
	}
	TfAzureKeyVaultConfigModel = encryptionatrest.TfAzureKeyVaultConfigModel{
		Enabled:           types.BoolValue(enabled),
		ClientID:          types.StringValue(clientID),
		AzureEnvironment:  types.StringValue(azureEnvironment),
		SubscriptionID:    types.StringValue(subscriptionID),
		ResourceGroupName: types.StringValue(resourceGroupName),
		KeyVaultName:      types.StringValue(keyVaultName),
		KeyIdentifier:     types.StringValue(keyIdentifier),
		TenantID:          types.StringValue(tenantID),
		Secret:            types.StringValue(secret),
	}
	GoogleCloudKMS = &admin.GoogleCloudKMS{
		Enabled:              &enabled,
		KeyVersionResourceID: &keyVersionResourceID,
		ServiceAccountKey:    &serviceAccountKey,
	}
	TfGcpKmsConfigModel = encryptionatrest.TfGcpKmsConfigModel{
		Enabled:              types.BoolValue(enabled),
		KeyVersionResourceID: types.StringValue(keyVersionResourceID),
		ServiceAccountKey:    types.StringValue(serviceAccountKey),
	}
	EncryptionAtRest = &admin.EncryptionAtRest{
		AwsKms:         AWSKMSConfiguration,
		AzureKeyVault:  AzureKeyVault,
		GoogleCloudKms: GoogleCloudKMS,
	}
)

func TestNewTfEncryptionAtRestRSModel(t *testing.T) {
	testCases := []struct {
		expectedResult *encryptionatrest.TfEncryptionAtRestRSModel
		sdkModel       *admin.EncryptionAtRest
		name           string
	}{
		{
			name:     "Success NewTFAwsKmsConfig",
			sdkModel: EncryptionAtRest,
			expectedResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				ID:                   types.StringValue(projectID),
				ProjectID:            types.StringValue(projectID),
				AwsKmsConfig:         []encryptionatrest.TfAwsKmsConfigModel{TfAwsKmsConfigModel},
				AzureKeyVaultConfig:  []encryptionatrest.TfAzureKeyVaultConfigModel{TfAzureKeyVaultConfigModel},
				GoogleCloudKmsConfig: []encryptionatrest.TfGcpKmsConfigModel{TfGcpKmsConfigModel},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewTfEncryptionAtRestRSModel(context.Background(), projectID, tc.sdkModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTFAwsKmsConfig(t *testing.T) {
	testCases := []struct {
		name           string
		sdkModel       *admin.AWSKMSConfiguration
		expectedResult []encryptionatrest.TfAwsKmsConfigModel
	}{
		{
			name:     "Success NewTFAwsKmsConfig",
			sdkModel: AWSKMSConfiguration,
			expectedResult: []encryptionatrest.TfAwsKmsConfigModel{
				TfAwsKmsConfigModel,
			},
		},
		{
			name:           "Empty sdkModel",
			sdkModel:       nil,
			expectedResult: []encryptionatrest.TfAwsKmsConfigModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewTFAwsKmsConfig(context.Background(), tc.sdkModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTFAzureKeyVaultConfig(t *testing.T) {
	testCases := []struct {
		name           string
		sdkModel       *admin.AzureKeyVault
		expectedResult []encryptionatrest.TfAzureKeyVaultConfigModel
	}{
		{
			name:     "Success NewTFAwsKmsConfig",
			sdkModel: AzureKeyVault,
			expectedResult: []encryptionatrest.TfAzureKeyVaultConfigModel{
				TfAzureKeyVaultConfigModel,
			},
		},
		{
			name:           "Empty sdkModel",
			sdkModel:       nil,
			expectedResult: []encryptionatrest.TfAzureKeyVaultConfigModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewTFAzureKeyVaultConfig(context.Background(), tc.sdkModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTFGcpKmsConfig(t *testing.T) {
	testCases := []struct {
		name           string
		sdkModel       *admin.GoogleCloudKMS
		expectedResult []encryptionatrest.TfGcpKmsConfigModel
	}{
		{
			name:     "Success NewTFGcpKmsConfig",
			sdkModel: GoogleCloudKMS,
			expectedResult: []encryptionatrest.TfGcpKmsConfigModel{
				TfGcpKmsConfigModel,
			},
		},
		{
			name:           "Empty sdkModel",
			sdkModel:       nil,
			expectedResult: []encryptionatrest.TfGcpKmsConfigModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewTFGcpKmsConfig(context.Background(), tc.sdkModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewAtlasAwsKms(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult *admin.AWSKMSConfiguration
		tfModel        []encryptionatrest.TfAwsKmsConfigModel
	}{
		{
			name:           "Success NewAtlasAwsKms",
			tfModel:        []encryptionatrest.TfAwsKmsConfigModel{TfAwsKmsConfigModel},
			expectedResult: AWSKMSConfiguration,
		},
		{
			name:           "Empty tfModel",
			tfModel:        nil,
			expectedResult: &admin.AWSKMSConfiguration{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewAtlasAwsKms(tc.tfModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewAtlasGcpKms(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult *admin.GoogleCloudKMS
		tfModel        []encryptionatrest.TfGcpKmsConfigModel
	}{
		{
			name:           "Success NewAtlasAwsKms",
			tfModel:        []encryptionatrest.TfGcpKmsConfigModel{TfGcpKmsConfigModel},
			expectedResult: GoogleCloudKMS,
		},
		{
			name:           "Empty tfModel",
			tfModel:        nil,
			expectedResult: &admin.GoogleCloudKMS{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewAtlasGcpKms(tc.tfModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewAtlasAzureKeyVault(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult *admin.AzureKeyVault
		tfModel        []encryptionatrest.TfAzureKeyVaultConfigModel
	}{
		{
			name:           "Success NewAtlasAwsKms",
			tfModel:        []encryptionatrest.TfAzureKeyVaultConfigModel{TfAzureKeyVaultConfigModel},
			expectedResult: AzureKeyVault,
		},
		{
			name:           "Empty tfModel",
			tfModel:        nil,
			expectedResult: &admin.AzureKeyVault{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewAtlasAzureKeyVault(tc.tfModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}
