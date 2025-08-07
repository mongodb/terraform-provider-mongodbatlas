package encryptionatrest_test

import (
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/encryptionatrest"
)

var (
	projectID                = "projectID"
	enabled                  = true
	requirePrivateNetworking = true
	customerMasterKeyID      = "CustomerMasterKeyID"
	region                   = "Region"
	accessKeyID              = "AccessKeyID"
	secretAccessKey          = "SecretAccessKey"
	roleID                   = "RoleID"
	clientID                 = "clientID"
	azureEnvironment         = "AzureEnvironment"
	subscriptionID           = "SubscriptionID"
	resourceGroupName        = "ResourceGroupName"
	keyVaultName             = "KeyVaultName"
	keyIdentifier            = "KeyIdentifier"
	tenantID                 = "TenantID"
	secret                   = "Secret"
	keyVersionResourceID     = "KeyVersionResourceID"
	serviceAccountKey        = "ServiceAccountKey"
	AWSKMSConfiguration      = &admin.AWSKMSConfiguration{
		Enabled:                  &enabled,
		CustomerMasterKeyID:      &customerMasterKeyID,
		Region:                   &region,
		AccessKeyID:              &accessKeyID,
		SecretAccessKey:          &secretAccessKey,
		RoleId:                   &roleID,
		RequirePrivateNetworking: &requirePrivateNetworking,
	}
	TfAwsKmsConfigModel = encryptionatrest.TFAwsKmsConfigModel{
		Enabled:                  types.BoolValue(enabled),
		CustomerMasterKeyID:      types.StringValue(customerMasterKeyID),
		Region:                   types.StringValue(region),
		AccessKeyID:              types.StringValue(accessKeyID),
		SecretAccessKey:          types.StringValue(secretAccessKey),
		RoleID:                   types.StringValue(roleID),
		RequirePrivateNetworking: types.BoolValue(requirePrivateNetworking),
	}
	AzureKeyVault = &admin.AzureKeyVault{
		Enabled:                  &enabled,
		ClientID:                 &clientID,
		AzureEnvironment:         &azureEnvironment,
		SubscriptionID:           &subscriptionID,
		ResourceGroupName:        &resourceGroupName,
		KeyVaultName:             &keyVaultName,
		KeyIdentifier:            &keyIdentifier,
		TenantID:                 &tenantID,
		Secret:                   &secret,
		RequirePrivateNetworking: &requirePrivateNetworking,
	}
	TfAzureKeyVaultConfigModel = encryptionatrest.TFAzureKeyVaultConfigModel{
		Enabled:                  types.BoolValue(enabled),
		ClientID:                 types.StringValue(clientID),
		AzureEnvironment:         types.StringValue(azureEnvironment),
		SubscriptionID:           types.StringValue(subscriptionID),
		ResourceGroupName:        types.StringValue(resourceGroupName),
		KeyVaultName:             types.StringValue(keyVaultName),
		KeyIdentifier:            types.StringValue(keyIdentifier),
		TenantID:                 types.StringValue(tenantID),
		Secret:                   types.StringValue(secret),
		RequirePrivateNetworking: types.BoolValue(requirePrivateNetworking),
	}
	GoogleCloudKMS = &admin.GoogleCloudKMS{
		Enabled:              &enabled,
		KeyVersionResourceID: &keyVersionResourceID,
		ServiceAccountKey:    &serviceAccountKey,
	}
	TfGcpKmsConfigModel = encryptionatrest.TFGcpKmsConfigModel{
		Enabled:              types.BoolValue(enabled),
		KeyVersionResourceID: types.StringValue(keyVersionResourceID),
		ServiceAccountKey:    types.StringValue(serviceAccountKey),
	}
	EncryptionAtRest = &admin.EncryptionAtRest{
		AwsKms:                AWSKMSConfiguration,
		AzureKeyVault:         AzureKeyVault,
		GoogleCloudKms:        GoogleCloudKMS,
		EnabledForSearchNodes: &enabled,
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
				ID:                    types.StringValue(projectID),
				ProjectID:             types.StringValue(projectID),
				AwsKmsConfig:          []encryptionatrest.TFAwsKmsConfigModel{TfAwsKmsConfigModel},
				AzureKeyVaultConfig:   []encryptionatrest.TFAzureKeyVaultConfigModel{TfAzureKeyVaultConfigModel},
				GoogleCloudKmsConfig:  []encryptionatrest.TFGcpKmsConfigModel{TfGcpKmsConfigModel},
				EnabledForSearchNodes: types.BoolValue(enabled),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewTFEncryptionAtRestRSModel(t.Context(), projectID, tc.sdkModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTFAwsKmsConfig(t *testing.T) {
	testCases := []struct {
		name           string
		sdkModel       *admin.AWSKMSConfiguration
		expectedResult []encryptionatrest.TFAwsKmsConfigModel
	}{
		{
			name:     "Success NewTFAwsKmsConfig",
			sdkModel: AWSKMSConfiguration,
			expectedResult: []encryptionatrest.TFAwsKmsConfigModel{
				TfAwsKmsConfigModel,
			},
		},
		{
			name:           "Empty sdkModel",
			sdkModel:       nil,
			expectedResult: []encryptionatrest.TFAwsKmsConfigModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewTFAwsKmsConfig(t.Context(), tc.sdkModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTFAzureKeyVaultConfig(t *testing.T) {
	testCases := []struct {
		name           string
		sdkModel       *admin.AzureKeyVault
		expectedResult []encryptionatrest.TFAzureKeyVaultConfigModel
	}{
		{
			name:     "Success NewTFAwsKmsConfig",
			sdkModel: AzureKeyVault,
			expectedResult: []encryptionatrest.TFAzureKeyVaultConfigModel{
				TfAzureKeyVaultConfigModel,
			},
		},
		{
			name:           "Empty sdkModel",
			sdkModel:       nil,
			expectedResult: []encryptionatrest.TFAzureKeyVaultConfigModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewTFAzureKeyVaultConfig(t.Context(), tc.sdkModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewTFGcpKmsConfig(t *testing.T) {
	testCases := []struct {
		name           string
		sdkModel       *admin.GoogleCloudKMS
		expectedResult []encryptionatrest.TFGcpKmsConfigModel
	}{
		{
			name:     "Success NewTFGcpKmsConfig",
			sdkModel: GoogleCloudKMS,
			expectedResult: []encryptionatrest.TFGcpKmsConfigModel{
				TfGcpKmsConfigModel,
			},
		},
		{
			name:           "Empty sdkModel",
			sdkModel:       nil,
			expectedResult: []encryptionatrest.TFGcpKmsConfigModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := encryptionatrest.NewTFGcpKmsConfig(t.Context(), tc.sdkModel)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestNewAtlasAwsKms(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult *admin.AWSKMSConfiguration
		tfModel        []encryptionatrest.TFAwsKmsConfigModel
	}{
		{
			name:           "Success NewAtlasAwsKms",
			tfModel:        []encryptionatrest.TFAwsKmsConfigModel{TfAwsKmsConfigModel},
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
		tfModel        []encryptionatrest.TFGcpKmsConfigModel
	}{
		{
			name:           "Success NewAtlasAwsKms",
			tfModel:        []encryptionatrest.TFGcpKmsConfigModel{TfGcpKmsConfigModel},
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
		tfModel        []encryptionatrest.TFAzureKeyVaultConfigModel
	}{
		{
			name:           "Success NewAtlasAwsKms",
			tfModel:        []encryptionatrest.TFAzureKeyVaultConfigModel{TfAzureKeyVaultConfigModel},
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
