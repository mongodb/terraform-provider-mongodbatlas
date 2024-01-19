package encryptionatrest_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/encryptionatrest"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mocksvc"
	"github.com/mwielbut/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

const (
	initialConfigEncryptionRestRoleAWS = `
provider "aws" {
	region     = lower(replace("%[1]s", "_", "-"))
	access_key = "%[2]s"
	secret_key = "%[3]s"
}

%[7]s

resource "mongodbatlas_cloud_provider_access" "test" {
	project_id = "%[4]s"
	provider_name = "AWS"
	%[8]s
		
}

resource "aws_iam_role_policy" "test_policy" {
  name = "%[5]s"
  role = aws_iam_role.test_role.id

  policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
		"Action": "*",
		"Resource": "*"
      }
    ]
  }
  EOF
}

resource "aws_iam_role" "test_role" {
 name = "%[6]s"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "${mongodbatlas_cloud_provider_access.test.atlas_aws_account_arn}"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${mongodbatlas_cloud_provider_access.test.atlas_assumed_role_external_id}"
        }
      }
    }
  ]
}
EOF

}

%[9]s

`
	configEncryptionRest = `
resource "mongodbatlas_encryption_at_rest" "test" {
	project_id = "%s"

	aws_kms_config {
		enabled                = %t
		customer_master_key_id = "%s"
		region                 = "%s"
		role_id = mongodbatlas_cloud_provider_access.test.role_id
	}
}`
	dataAWSARNConfig = `
data "aws_iam_role" "test" {
  name = "%s"
}

`
)

func TestAccAdvRSEncryptionAtRest_basicAWS(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		awsKms = admin.AWSKMSConfiguration{
			Enabled:             pointy.Bool(true),
			CustomerMasterKeyID: conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:              conversion.StringPtr(os.Getenv("AWS_REGION")),
			RoleId:              conversion.StringPtr(os.Getenv("AWS_ROLE_ID")),
		}

		awsKmsUpdated = admin.AWSKMSConfiguration{
			Enabled:             pointy.Bool(true),
			CustomerMasterKeyID: conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:              conversion.StringPtr(os.Getenv("AWS_REGION")),
			RoleId:              conversion.StringPtr(os.Getenv("AWS_ROLE_ID")),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckAwsEnv(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID, &awsKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.region", awsKms.GetRegion()),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.role_id", awsKms.GetRoleId()),

					resource.TestCheckNoResourceAttr(resourceName, "azure_key_vault_config.#"),
					resource.TestCheckNoResourceAttr(resourceName, "google_cloud_kms_config.#"),
				),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID, &awsKmsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.region", awsKmsUpdated.GetRegion()),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.role_id", awsKmsUpdated.GetRoleId()),

					resource.TestCheckNoResourceAttr(resourceName, "azure_key_vault_config.#"),
					resource.TestCheckNoResourceAttr(resourceName, "google_cloud_kms_config.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAdvRSEncryptionAtRest_basicAzure(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		azureKeyVault = admin.AzureKeyVault{
			Enabled:           pointy.Bool(true),
			ClientID:          conversion.StringPtr(os.Getenv("AZURE_CLIENT_ID")),
			AzureEnvironment:  conversion.StringPtr("AZURE"),
			SubscriptionID:    conversion.StringPtr(os.Getenv("AZURE_SUBSCRIPTION_ID")),
			ResourceGroupName: conversion.StringPtr(os.Getenv("AZURE_RESOURCE_GROUP_NAME")),
			KeyVaultName:      conversion.StringPtr(os.Getenv("AZURE_KEY_VAULT_NAME")),
			KeyIdentifier:     conversion.StringPtr(os.Getenv("AZURE_KEY_IDENTIFIER")),
			Secret:            conversion.StringPtr(os.Getenv("AZURE_SECRET")),
			TenantID:          conversion.StringPtr(os.Getenv("AZURE_TENANT_ID")),
		}

		azureKeyVaultUpdated = admin.AzureKeyVault{
			Enabled:           pointy.Bool(true),
			ClientID:          conversion.StringPtr(os.Getenv("AZURE_CLIENT_ID_UPDATED")),
			AzureEnvironment:  conversion.StringPtr("AZURE"),
			SubscriptionID:    conversion.StringPtr(os.Getenv("AZURE_SUBSCRIPTION_ID")),
			ResourceGroupName: conversion.StringPtr(os.Getenv("AZURE_RESOURCE_GROUP_NAME_UPDATED")),
			KeyVaultName:      conversion.StringPtr(os.Getenv("AZURE_KEY_VAULT_NAME_UPDATED")),
			KeyIdentifier:     conversion.StringPtr(os.Getenv("AZURE_KEY_IDENTIFIER_UPDATED")),
			Secret:            conversion.StringPtr(os.Getenv("AZURE_SECRET_UPDATED")),
			TenantID:          conversion.StringPtr(os.Getenv("AZURE_TENANT_ID")),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckEncryptionAtRestEnvAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(projectID, &azureKeyVault),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.azure_environment", azureKeyVault.GetAzureEnvironment()),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.resource_group_name", azureKeyVault.GetResourceGroupName()),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.key_vault_name", azureKeyVault.GetKeyVaultName()),
				),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(projectID, &azureKeyVaultUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.azure_environment", azureKeyVaultUpdated.GetAzureEnvironment()),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.resource_group_name", azureKeyVaultUpdated.GetResourceGroupName()),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.key_vault_name", azureKeyVaultUpdated.GetKeyVaultName()),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// "azure_key_vault_config.0.secret" is a sensitive value not returned by the API
				ImportStateVerifyIgnore: []string{"azure_key_vault_config.0.secret"},
			},
		},
	})
}

func TestAccAdvRSEncryptionAtRest_basicGCP(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		googleCloudKms = admin.GoogleCloudKMS{
			Enabled:              pointy.Bool(true),
			ServiceAccountKey:    conversion.StringPtr(os.Getenv("GCP_SERVICE_ACCOUNT_KEY")),
			KeyVersionResourceID: conversion.StringPtr(os.Getenv("GCP_KEY_VERSION_RESOURCE_ID")),
		}

		googleCloudKmsUpdated = admin.GoogleCloudKMS{
			Enabled:              pointy.Bool(true),
			ServiceAccountKey:    conversion.StringPtr(os.Getenv("GCP_SERVICE_ACCOUNT_KEY_UPDATED")),
			KeyVersionResourceID: conversion.StringPtr(os.Getenv("GCP_KEY_VERSION_RESOURCE_ID_UPDATED")),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckGPCEnv(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(projectID, &googleCloudKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms_config.0.enabled", "true"),
				),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(projectID, &googleCloudKmsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms_config.0.enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// "google_cloud_kms_config.0.service_account_key" is a sensitive value not returned by the API
				ImportStateVerifyIgnore: []string{"google_cloud_kms_config.0.service_account_key"},
			},
		},
	})
}

func TestAccAdvRSEncryptionAtRestWithRole_basicAWS(t *testing.T) {
	acc.SkipTest(t) // For now it will skipped because of aws errors reasons, already made another test using terratest.
	acc.SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		accessKeyID  = os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey    = os.Getenv("AWS_SECRET_ACCESS_KEY")
		policyName   = acctest.RandomWithPrefix("test-aws-policy")
		roleName     = acctest.RandomWithPrefix("test-aws-role")

		awsKms = admin.AWSKMSConfiguration{
			Enabled:             pointy.Bool(true),
			CustomerMasterKeyID: conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:              conversion.StringPtr(os.Getenv("AWS_REGION")),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckAwsEnv(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(awsKms.GetRegion(), accessKeyID, secretKey, projectID, policyName, roleName, false, &awsKms),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(awsKms.GetRegion(), accessKeyID, secretKey, projectID, policyName, roleName, true, &awsKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var (
	ServiceAccountKey                      = types.StringValue("service")
	googleCloudConfigWithServiceAccountKey = []encryptionatrest.TfGcpKmsConfigModel{
		{
			ServiceAccountKey: ServiceAccountKey,
		},
	}
	awsConfigWithRegion = []encryptionatrest.TfAwsKmsConfigModel{
		{
			Region: types.StringValue(region),
		},
	}
	awsConfigWithRegionAndSecretAccessKey = []encryptionatrest.TfAwsKmsConfigModel{
		{
			Region:          types.StringValue(region),
			SecretAccessKey: ServiceAccountKey,
		},
	}
	azureConfigWithSecret = []encryptionatrest.TfAzureKeyVaultConfigModel{
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
				GoogleCloudKmsConfig: []encryptionatrest.TfGcpKmsConfigModel{},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: []encryptionatrest.TfGcpKmsConfigModel{},
			},
		},
		{
			name: "Current GoogleCloudKmsConfig not nil, GoogleCloudKmsConfig config is available",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: []encryptionatrest.TfGcpKmsConfigModel{},
			},
			earRSConfig: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: googleCloudConfigWithServiceAccountKey,
			},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				GoogleCloudKmsConfig: []encryptionatrest.TfGcpKmsConfigModel{{}},
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
				GoogleCloudKmsConfig: []encryptionatrest.TfGcpKmsConfigModel{{}},
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
				AwsKmsConfig: []encryptionatrest.TfAwsKmsConfigModel{},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: []encryptionatrest.TfAwsKmsConfigModel{},
			},
		},
		{
			name: "Current AwsKmsConfig not nil, AwsKmsConfig config is available",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: []encryptionatrest.TfAwsKmsConfigModel{},
			},
			earRSConfig: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: awsConfigWithRegion,
			},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				AwsKmsConfig: []encryptionatrest.TfAwsKmsConfigModel{{}},
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
				AwsKmsConfig: []encryptionatrest.TfAwsKmsConfigModel{{}},
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
				AzureKeyVaultConfig: []encryptionatrest.TfAzureKeyVaultConfigModel{},
			},
			expectedEarResult: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: []encryptionatrest.TfAzureKeyVaultConfigModel{},
			},
		},
		{
			name: "Current AzureKeyVaultConfig not nil, AzureKeyVaultConfig config is available",
			earRSCurrent: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: []encryptionatrest.TfAzureKeyVaultConfigModel{},
			},
			earRSConfig: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: azureConfigWithSecret,
			},
			earRSNew: &encryptionatrest.TfEncryptionAtRestRSModel{
				AzureKeyVaultConfig: []encryptionatrest.TfAzureKeyVaultConfigModel{{}},
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
				AzureKeyVaultConfig: []encryptionatrest.TfAzureKeyVaultConfigModel{{}},
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
			testObject := mocksvc.NewEarService(t)

			testObject.On("UpdateEncryptionAtRest", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockResponse, nil, tc.mockError)

			response, strategy, err := encryptionatrest.ResourceMongoDBAtlasEncryptionAtRestCreateRefreshFunc(context.Background(), projectID, testObject, &admin.EncryptionAtRest{})()

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}

			assert.Equal(t, tc.expectedResponse, response)
			assert.Equal(t, tc.expectedRetrystrategy, strategy)
		})
	}
}

func testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		connV2 := acc.TestMongoDBClient.(*config.MongoDBClient).AtlasV2

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if _, _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRest(context.Background(), rs.Primary.ID).Execute(); err == nil {
			return nil
		}

		return fmt.Errorf("encryptionAtRest (%s) does not exist", rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasEncryptionAtRestDestroy(s *terraform.State) error {
	connV2 := acc.TestMongoDBClient.(*config.MongoDBClient).AtlasV2

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_encryption_at_rest" {
			continue
		}

		res, _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRest(context.Background(), rs.Primary.ID).Execute()

		if err != nil ||
			(*res.AwsKms.Enabled != false &&
				*res.AzureKeyVault.Enabled != false &&
				*res.GoogleCloudKms.Enabled != false) {
			return fmt.Errorf("encryptionAtRest (%s) still exists: err: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID string, aws *admin.AWSKMSConfiguration) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "%s"

		  aws_kms_config {
				enabled                = %t
				customer_master_key_id = "%s"
				region                 = "%s"
				role_id              = "%s"
			}
		}
	`, projectID, aws.GetEnabled(), aws.GetCustomerMasterKeyID(), aws.GetRegion(), aws.GetRoleId())
}

func testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(projectID string, azure *admin.AzureKeyVault) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "%s"

		  azure_key_vault_config {
				enabled             = %t
				client_id           = "%s"
				azure_environment   = "%s"
				subscription_id     = "%s"
				resource_group_name = "%s"
				key_vault_name  	  = "%s"
				key_identifier  	  = "%s"
				secret  						= "%s"
				tenant_id  					= "%s"
			}
		}
	`, projectID, *azure.Enabled, azure.GetClientID(), azure.GetAzureEnvironment(), azure.GetSubscriptionID(), azure.GetResourceGroupName(),
		azure.GetKeyVaultName(), azure.GetKeyIdentifier(), azure.GetSecret(), azure.GetTenantID())
}

func testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(projectID string, google *admin.GoogleCloudKMS) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "%s"

		  google_cloud_kms_config {
				enabled                 = %t
				service_account_key     = "%s"
				key_version_resource_id = "%s"
			}
		}
	`, projectID, *google.Enabled, google.GetServiceAccountKey(), google.GetKeyVersionResourceID())
}

func testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(region, awsAccesKey, awsSecretKey, projectID, policyName, awsRoleName string, isUpdate bool, aws *admin.AWSKMSConfiguration) string {
	cfg := fmt.Sprintf(initialConfigEncryptionRestRoleAWS, region, awsAccesKey, awsSecretKey, projectID, policyName, awsRoleName, "", "", "")
	if isUpdate {
		configEncrypt := fmt.Sprintf(configEncryptionRest, projectID, *aws.Enabled, aws.GetCustomerMasterKeyID(), aws.GetRegion())
		dataAWSARN := fmt.Sprintf(dataAWSARNConfig, awsRoleName)
		dataARN := `iam_assumed_role_arn = data.aws_iam_role.test.arn`
		cfg = fmt.Sprintf(initialConfigEncryptionRestRoleAWS, region, awsAccesKey, awsSecretKey, projectID, policyName, awsRoleName, dataAWSARN, dataARN, configEncrypt)
	}
	return cfg
}

func testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
