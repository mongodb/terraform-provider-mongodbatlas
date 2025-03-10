package encryptionatrest_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	// TODO: update before merging to master: "go.mongodb.org/atlas-sdk/v20250219001/admin"
	"github.com/mongodb/atlas-sdk-go/admin"

	// TODO: update before merging to master: "go.mongodb.org/atlas-sdk/v20250219001/mockadmin"
	"github.com/mongodb/atlas-sdk-go/mockadmin"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/encryptionatrest"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_encryption_at_rest.test"
	datasourceName = "data.mongodbatlas_encryption_at_rest.test"
)

func TestAccEncryptionAtRest_basicAWS(t *testing.T) {
	var (
		projectID = os.Getenv("MONGODB_ATLAS_PROJECT_EAR_PE_AWS_ID") // to use RequirePrivateNetworking, Atlas Project is required to have FF enabled

		awsKms = admin.AWSKMSConfiguration{
			Enabled:                  conversion.Pointer(true),
			CustomerMasterKeyID:      conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:                   conversion.StringPtr(conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))),
			RoleId:                   conversion.StringPtr(os.Getenv("AWS_EAR_ROLE_ID")),
			RequirePrivateNetworking: conversion.Pointer(false),
		}
		awsKmsAttrMap = acc.ConvertToAwsKmsEARAttrMap(&awsKms)

		awsKmsUpdated = admin.AWSKMSConfiguration{
			Enabled:                  conversion.Pointer(true),
			CustomerMasterKeyID:      conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:                   conversion.StringPtr(conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))),
			RoleId:                   conversion.StringPtr(os.Getenv("AWS_EAR_ROLE_ID")),
			RequirePrivateNetworking: conversion.Pointer(true),
		}
		awsKmsUpdatedAttrMap = acc.ConvertToAwsKmsEARAttrMap(&awsKmsUpdated)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAwsEnv(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigAwsKms(projectID, &awsKms, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(resourceName, "aws_kms_config.0", awsKmsAttrMap),

					resource.TestCheckNoResourceAttr(resourceName, "azure_key_vault_config.#"),
					resource.TestCheckNoResourceAttr(resourceName, "google_cloud_kms_config.#"),

					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(datasourceName, "aws_kms_config.", awsKmsAttrMap),
				),
			},
			{
				Config: acc.ConfigAwsKms(projectID, &awsKmsUpdated, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(resourceName, "aws_kms_config.0", awsKmsUpdatedAttrMap),

					resource.TestCheckNoResourceAttr(resourceName, "azure_key_vault_config.#"),
					resource.TestCheckNoResourceAttr(resourceName, "google_cloud_kms_config.#"),

					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(datasourceName, "aws_kms_config", awsKmsUpdatedAttrMap),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: acc.EARImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEncryptionAtRest_basicAzure(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)

		azureKeyVault = admin.AzureKeyVault{
			Enabled:                  conversion.Pointer(true),
			ClientID:                 conversion.StringPtr(os.Getenv("AZURE_CLIENT_ID")),
			AzureEnvironment:         conversion.StringPtr("AZURE"),
			SubscriptionID:           conversion.StringPtr(os.Getenv("AZURE_SUBSCRIPTION_ID")),
			ResourceGroupName:        conversion.StringPtr(os.Getenv("AZURE_RESOURCE_GROUP_NAME")),
			KeyVaultName:             conversion.StringPtr(os.Getenv("AZURE_KEY_VAULT_NAME")),
			KeyIdentifier:            conversion.StringPtr(os.Getenv("AZURE_KEY_IDENTIFIER")),
			Secret:                   conversion.StringPtr(os.Getenv("AZURE_APP_SECRET")),
			TenantID:                 conversion.StringPtr(os.Getenv("AZURE_TENANT_ID")),
			RequirePrivateNetworking: conversion.Pointer(false),
		}

		azureKeyVaultAttrMap = acc.ConvertToAzureKeyVaultEARAttrMap(&azureKeyVault)

		azureKeyVaultUpdated = admin.AzureKeyVault{
			Enabled:                  conversion.Pointer(true),
			ClientID:                 conversion.StringPtr(os.Getenv("AZURE_CLIENT_ID")),
			AzureEnvironment:         conversion.StringPtr("AZURE"),
			SubscriptionID:           conversion.StringPtr(os.Getenv("AZURE_SUBSCRIPTION_ID")),
			ResourceGroupName:        conversion.StringPtr(os.Getenv("AZURE_RESOURCE_GROUP_NAME")),
			KeyVaultName:             conversion.StringPtr(os.Getenv("AZURE_KEY_VAULT_NAME_UPDATED")),
			KeyIdentifier:            conversion.StringPtr(os.Getenv("AZURE_KEY_IDENTIFIER_UPDATED")),
			Secret:                   conversion.StringPtr(os.Getenv("AZURE_APP_SECRET")),
			TenantID:                 conversion.StringPtr(os.Getenv("AZURE_TENANT_ID")),
			RequirePrivateNetworking: conversion.Pointer(false),
		}

		azureKeyVaultUpdatedAttrMap = acc.ConvertToAzureKeyVaultEARAttrMap(&azureKeyVaultUpdated)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckEncryptionAtRestEnvAzureWithUpdate(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigEARAzureKeyVault(projectID, &azureKeyVault, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(resourceName, "azure_key_vault_config.0", azureKeyVaultAttrMap),
					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(datasourceName, "azure_key_vault_config", azureKeyVaultAttrMap),
				),
			},
			{
				Config: acc.ConfigEARAzureKeyVault(projectID, &azureKeyVaultUpdated, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(resourceName, "azure_key_vault_config.0", azureKeyVaultUpdatedAttrMap),
					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(datasourceName, "azure_key_vault_config", azureKeyVaultUpdatedAttrMap),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: acc.EARImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// "azure_key_vault_config.0.secret" is a sensitive value not returned by the API
				ImportStateVerifyIgnore: []string{"azure_key_vault_config.0.secret"},
			},
		},
	})
}

func TestAccEncryptionAtRest_basicGCP(t *testing.T) {
	acc.SkipTestForCI(t) // needs GCP configuration

	var (
		projectID = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		googleCloudKms = admin.GoogleCloudKMS{
			Enabled:              conversion.Pointer(true),
			ServiceAccountKey:    conversion.StringPtr(os.Getenv("GCP_SERVICE_ACCOUNT_KEY")),
			KeyVersionResourceID: conversion.StringPtr(os.Getenv("GCP_KEY_VERSION_RESOURCE_ID")),
		}

		googleCloudKmsUpdated = admin.GoogleCloudKMS{
			Enabled:              conversion.Pointer(true),
			ServiceAccountKey:    conversion.StringPtr(os.Getenv("GCP_SERVICE_ACCOUNT_KEY_UPDATED")),
			KeyVersionResourceID: conversion.StringPtr(os.Getenv("GCP_KEY_VERSION_RESOURCE_ID_UPDATED")),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckGPCEnv(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				Config: configGoogleCloudKms(projectID, &googleCloudKms, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms_config.0.valid", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "google_cloud_kms_config.0.key_version_resource_id"),

					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(datasourceName, "google_cloud_kms_config.enabled", "true"),
					resource.TestCheckResourceAttr(datasourceName, "google_cloud_kms_config.valid", "true"),
					resource.TestCheckResourceAttrSet(datasourceName, "google_cloud_kms_config.key_version_resource_id"),
				),
			},
			{
				Config: configGoogleCloudKms(projectID, &googleCloudKmsUpdated, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms_config.0.valid", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "google_cloud_kms_config.0.key_version_resource_id"),

					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(datasourceName, "google_cloud_kms_config.enabled", "true"),
					resource.TestCheckResourceAttr(datasourceName, "google_cloud_kms_config.valid", "true"),
					resource.TestCheckResourceAttrSet(datasourceName, "google_cloud_kms_config.key_version_resource_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: acc.EARImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// "google_cloud_kms_config.0.service_account_key" is a sensitive value not returned by the API
				ImportStateVerifyIgnore: []string{"google_cloud_kms_config.0.service_account_key"},
			},
		},
	})
}

func TestAccEncryptionAtRestWithRole_basicAWS(t *testing.T) {
	acc.SkipTestForCI(t) // needs AWS configuration. This test case is similar to TestAccEncryptionAtRest_basicAWS except that it creates it's own AWS resources such as IAM roles, cloud provider access, etc so we don't need to run this in CI but may be used for local testing.

	resource.Test(t, *testCaseWithRoleBasicAWS(t))
}

func testCaseWithRoleBasicAWS(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID            = acc.ProjectIDExecution(t)
		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)
		awsKeyName           = acc.RandomName()
		awsKms               = admin.AWSKMSConfiguration{
			Enabled:             conversion.Pointer(true),
			Region:              conversion.StringPtr(conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))),
			CustomerMasterKeyID: conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
		}
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAwsEnv(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(projectID, awsIAMRoleName, awsIAMRolePolicyName, awsKeyName, &awsKms),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "aws_kms_config.0.role_id"),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.enabled", "true"),

					resource.TestCheckNoResourceAttr(resourceName, "azure_key_vault_config.#"),
					resource.TestCheckNoResourceAttr(resourceName, "google_cloud_kms_config.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: acc.EARImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

var (
	ServiceAccountKey                      = types.StringValue("service")
	googleCloudConfigWithServiceAccountKey = []encryptionatrest.TFGcpKmsConfigModel{
		{
			ServiceAccountKey: ServiceAccountKey,
		},
	}
	awsConfigWithRegion = []encryptionatrest.TFAwsKmsConfigModel{
		{
			Region: types.StringValue(region),
		},
	}
	awsConfigWithRegionAndSecretAccessKey = []encryptionatrest.TFAwsKmsConfigModel{
		{
			Region:          types.StringValue(region),
			SecretAccessKey: ServiceAccountKey,
		},
	}
	azureConfigWithSecret = []encryptionatrest.TFAzureKeyVaultConfigModel{
		{
			Secret: types.StringValue(secret),
		},
	}
)

type testHandleConfig struct {
	earRSCurrent      *encryptionatrest.TfEncryptionAtRestRSModel
	earRSNew          *encryptionatrest.TfEncryptionAtRestRSModel
	earRSConfig       *encryptionatrest.TfEncryptionAtRestRSModel
	expectedEarResult *encryptionatrest.TfEncryptionAtRestRSModel
	name              string
}

func TestHandleGcpKmsConfig(t *testing.T) {
	testCases := []testHandleConfig{
		{
			name: "Current GoogleCloudKmsConfig is nil",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: nil,
			},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: []encryptionatrest.TFGcpKmsConfigModel{},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: []encryptionatrest.TFGcpKmsConfigModel{},
			},
		},
		{
			name: "Current GoogleCloudKmsConfig not nil, GoogleCloudKmsConfig config is available",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: []encryptionatrest.TFGcpKmsConfigModel{},
			},
			earRSConfig: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: googleCloudConfigWithServiceAccountKey,
			},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: []encryptionatrest.TFGcpKmsConfigModel{{}},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: googleCloudConfigWithServiceAccountKey,
			},
		},
		{
			name: "Current GoogleCloudKmsConfig not nil, GoogleCloudKmsConfig config is not available",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: googleCloudConfigWithServiceAccountKey,
			},
			earRSConfig: &encryptionatrest.TfEncryptionAtRestRSModel{},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: []encryptionatrest.TFGcpKmsConfigModel{{}},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: googleCloudConfigWithServiceAccountKey,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encryptionatrest.HandleGcpKmsConfig(context.Background(), tc.earRSCurrent, tc.earRSNew, tc.earRSConfig)
			assert.Equal(t, tc.expectedEarResult, tc.earRSNew, "result did not match expected output")
		})
	}
}

func TestHandleAwsKmsConfigDefaults(t *testing.T) {
	testCases := []testHandleConfig{
		{
			name: "Current AwsKmsConfig is nil",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: nil,
			},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: []encryptionatrest.TFAwsKmsConfigModel{},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: []encryptionatrest.TFAwsKmsConfigModel{},
			},
		},
		{
			name: "Current AwsKmsConfig not nil, AwsKmsConfig config is available",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: []encryptionatrest.TFAwsKmsConfigModel{},
			},
			earRSConfig: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: awsConfigWithRegion,
			},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: []encryptionatrest.TFAwsKmsConfigModel{{}},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: awsConfigWithRegion,
			},
		},
		{
			name: "Current AwsKmsConfig not nil, AwsKmsConfig config is not available",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: awsConfigWithRegionAndSecretAccessKey,
			},
			earRSConfig: &encryptionatrest.TfEncryptionAtRestRSModel{},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: []encryptionatrest.TFAwsKmsConfigModel{{}},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: awsConfigWithRegionAndSecretAccessKey,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encryptionatrest.HandleAwsKmsConfigDefaults(context.Background(), tc.earRSCurrent, tc.earRSNew, tc.earRSConfig)
			assert.Equal(t, tc.expectedEarResult, tc.earRSNew, "result did not match expected output")
		})
	}
}

func TestHandleAzureKeyVaultConfigDefaults(t *testing.T) {
	testCases := []testHandleConfig{
		{
			name: "Current AzureKeyVaultConfig is nil",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: nil,
			},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: []encryptionatrest.TFAzureKeyVaultConfigModel{},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: []encryptionatrest.TFAzureKeyVaultConfigModel{},
			},
		},
		{
			name: "Current AzureKeyVaultConfig not nil, AzureKeyVaultConfig config is available",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: []encryptionatrest.TFAzureKeyVaultConfigModel{},
			},
			earRSConfig: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: azureConfigWithSecret,
			},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: []encryptionatrest.TFAzureKeyVaultConfigModel{{}},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: azureConfigWithSecret,
			},
		},
		{
			name: "Current AzureKeyVaultConfig not nil, AzureKeyVaultConfig config is not available",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: azureConfigWithSecret,
			},
			earRSConfig: &encryptionatrest.TfEncryptionAtRestRSModel{},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: []encryptionatrest.TFAzureKeyVaultConfigModel{{}},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: azureConfigWithSecret,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encryptionatrest.HandleAzureKeyVaultConfigDefaults(context.Background(), tc.earRSCurrent, tc.earRSNew, tc.earRSConfig)
			assert.Equal(t, tc.expectedEarResult, tc.earRSNew, "result did not match expected output")
		})
	}
}

func TestResourceMongoDBAtlasEncryptionAtRestCreateRefreshFunc(t *testing.T) {
	var projectID = "projectID"
	testCases := []struct {
		name                  string
		mockResponse          *admin.EncryptionAtRest
		mockError             error
		expectedResponse      *admin.EncryptionAtRest
		expectedRetrystrategy string
		expectedError         bool
	}{
		{
			name:                  "Successful API call",
			mockResponse:          &admin.EncryptionAtRest{},
			mockError:             nil,
			expectedResponse:      &admin.EncryptionAtRest{},
			expectedRetrystrategy: retrystrategy.RetryStrategyCompletedState,
			expectedError:         false,
		},
		{
			name:                  "Failed API call: Error not one of CANNOT_ASSUME_ROLE, INVALID_AWS_CREDENTIALS, CLOUD_PROVIDER_ACCESS_ROLE_NOT_AUTHORIZED",
			mockResponse:          nil,
			mockError:             errors.New("random error"),
			expectedResponse:      nil,
			expectedRetrystrategy: retrystrategy.RetryStrategyErrorState,
			expectedError:         true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := mockadmin.NewEncryptionAtRestUsingCustomerKeyManagementApi(t)

			m.EXPECT().UpdateEncryptionAtRest(mock.Anything, mock.Anything, mock.Anything).Return(admin.UpdateEncryptionAtRestApiRequest{ApiService: m})
			m.EXPECT().UpdateEncryptionAtRestExecute(mock.Anything).Return(tc.mockResponse, nil, tc.mockError).Once()

			response, strategy, err := encryptionatrest.ResourceMongoDBAtlasEncryptionAtRestCreateRefreshFunc(context.Background(), projectID, m, &admin.EncryptionAtRest{})()

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}

			assert.Equal(t, tc.expectedResponse, response)
			assert.Equal(t, tc.expectedRetrystrategy, strategy)
		})
	}
}

func configGoogleCloudKms(projectID string, google *admin.GoogleCloudKMS, useDatasource bool) string {
	config := fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "%s"

		  google_cloud_kms_config {
				enabled                 = %t
				service_account_key     = "%s"
				key_version_resource_id = "%s"
			}
		}
	`, projectID, *google.Enabled, google.GetServiceAccountKey(), google.GetKeyVersionResourceID())

	if useDatasource {
		return fmt.Sprintf(`%s %s`, config, acc.EARDatasourceConfig())
	}
	return config
}

func testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(projectID, awsIAMRoleName, awsIAMRolePolicyName, awsKeyName string, awsEar *admin.AWSKMSConfiguration) string {
	test := fmt.Sprintf(`
	locals {
		project_id = %[1]q
		aws_iam_role_policy_name = %[2]q
		aws_iam_role_name        = %[3]q
		aws_kms_key_name         = %[4]q
	  }

		  %[5]s	
`, projectID, awsIAMRolePolicyName, awsIAMRoleName, awsKeyName, awsIAMroleAuthAndEarConfigUsingLocals(awsEar))
	return test
}

func awsIAMroleAuthAndEarConfigUsingLocals(awsEar *admin.AWSKMSConfiguration) string {
	return fmt.Sprintf(`  
	resource "aws_iam_role_policy" "test_policy" {
		name = local.aws_iam_role_policy_name
		role = aws_iam_role.test_role.id
	  
		policy = jsonencode({
		  "Version" : "2012-10-17",
		  "Statement" : [
			{
			  "Effect" : "Allow",
			  "Action" : [
				"kms:Decrypt",
				"kms:Encrypt",
				"kms:DescribeKey"
			  ],
			  "Resource" : [
				%[3]q
			  ]
			}
		  ]
		})
	  }
	  
	resource "aws_iam_role" "test_role" {
		name = local.aws_iam_role_name
	  
		assume_role_policy = jsonencode({
		  "Version" : "2012-10-17",
		  "Statement" : [
			{
			  "Effect" : "Allow",
			  "Principal" : {
				"AWS" : "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_aws_account_arn}"
			  },
			  "Action" : "sts:AssumeRole",
			  "Condition" : {
				"StringEquals" : {
				  "sts:ExternalId" : "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_assumed_role_external_id}"
				}
			  }
			}
		  ]
		})
	  }

	resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
		project_id    = local.project_id
		provider_name = "AWS"
	  }
	  
	  resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
		project_id = local.project_id
		role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
	  
		aws {
		  iam_assumed_role_arn = aws_iam_role.test_role.arn
		}
	  }

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = local.project_id

  aws_kms_config {
    enabled                = %[1]t
    customer_master_key_id = %[3]q
	region                 = %[2]q
    role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
	require_private_networking = %[4]t
  }
}
	`, awsEar.GetEnabled(), awsEar.GetRegion(), awsEar.GetCustomerMasterKeyID(), awsEar.GetRequirePrivateNetworking())
}
