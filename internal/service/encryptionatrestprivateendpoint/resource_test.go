package encryptionatrestprivateendpoint_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/encryptionatrestprivateendpoint"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName         = "mongodbatlas_encryption_at_rest_private_endpoint.test"
	dataSourceName       = "data.mongodbatlas_encryption_at_rest_private_endpoint.test"
	pluralDataSourceName = "data.mongodbatlas_encryption_at_rest_private_endpoints.test"
	earResourceName      = "mongodbatlas_encryption_at_rest.test"
	earDatasourceName    = "data.mongodbatlas_encryption_at_rest.test"
)

func TestAccEncryptionAtRestPrivateEndpoint_Azure_basic(t *testing.T) {
	resource.Test(t, *basicTestCaseAzure(t))
}

func TestAccEncryptionAtRestPrivateEndpoint_createTimeoutWithDeleteOnCreate(t *testing.T) {
	// This test is skipped because it creates a race condition with other tests:
	// 1. This test creates an encryption at rest private endpoint with a 1s timeout, causing it to fail and trigger cleanup
	// 2. The private endpoint deletion doesn't complete immediately
	// 3. Other tests share the same project and attempt to disable encryption at rest during cleanup
	// 4. MongoDB Atlas returns "CANNOT_DISABLE_ENCRYPTION_AT_REST_REQUIRE_PRIVATE_NETWORKING_WHILE_PRIVATE_ENDPOINTS_EXIST"
	//    because the private endpoint from this test is still being deleted
	// This race condition occurs even when tests don't run in parallel due to the async nature of private endpoint deletion.
	acc.SkipTestForCI(t)
	var (
		createTimeout         = "1s"
		deleteOnCreateTimeout = true
		region                = conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))
		// Create encryption at rest configuration outside of test configuration to avoid cleanup issues
		projectID = acc.EncryptionAtRestExecution(t)
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckEncryptionAtRestEnvAWS(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configEARPrivateEndpointWithTimeout(projectID, region, acc.TimeoutConfig(&createTimeout, nil, nil), &deleteOnCreateTimeout),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
			},
		},
	})
}

func basicTestCaseAzure(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID     = os.Getenv("MONGODB_ATLAS_PROJECT_EAR_PE_ID")
		azureKeyVault = &admin.AzureKeyVault{
			Enabled:                  conversion.Pointer(true),
			RequirePrivateNetworking: conversion.Pointer(true),
			AzureEnvironment:         conversion.StringPtr("AZURE"),
			ClientID:                 conversion.StringPtr(os.Getenv("AZURE_CLIENT_ID")),
			SubscriptionID:           conversion.StringPtr(os.Getenv("AZURE_SUBSCRIPTION_ID")),
			ResourceGroupName:        conversion.StringPtr(os.Getenv("AZURE_RESOURCE_GROUP_NAME")),
			KeyVaultName:             conversion.StringPtr(os.Getenv("AZURE_KEY_VAULT_NAME")),
			KeyIdentifier:            conversion.StringPtr(os.Getenv("AZURE_KEY_IDENTIFIER")),
			Secret:                   conversion.StringPtr(os.Getenv("AZURE_APP_SECRET")),
			TenantID:                 conversion.StringPtr(os.Getenv("AZURE_TENANT_ID")),
		}
		region = os.Getenv("AZURE_PRIVATE_ENDPOINT_REGION")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb); acc.PreCheckEncryptionAtRestEnvAzure(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigEARAzureKeyVault(projectID, azureKeyVault, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(earResourceName, "azure_key_vault_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(earResourceName, "azure_key_vault_config.0.require_private_networking", "false"),
				),
			},
			{
				Config: configAzureBasic(projectID, azureKeyVault, region, false),
				Check:  checkBasic(projectID, *azureKeyVault.AzureEnvironment, region, false),
			},
			{
				Config:            configAzureBasic(projectID, azureKeyVault, region, false),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func TestAccEncryptionAtRestPrivateEndpoint_approveEndpointWithAzureProvider(t *testing.T) {
	acc.SkipTestForCI(t) // uses azure/azapi Terraform provider which can log sensitive information in CI like Azure subscriptionID used in parent_id of the resource

	var (
		subscriptionID    = os.Getenv("AZURE_SUBSCRIPTION_ID")
		resourceGroupName = os.Getenv("AZURE_RESOURCE_GROUP_NAME")
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_EAR_PE_ID")
		keyVaultName      = os.Getenv("AZURE_KEY_VAULT_NAME")
		azureKeyVault     = &admin.AzureKeyVault{
			Enabled:                  conversion.Pointer(true),
			RequirePrivateNetworking: conversion.Pointer(true),
			AzureEnvironment:         conversion.StringPtr("AZURE"),
			ClientID:                 conversion.StringPtr(os.Getenv("AZURE_CLIENT_ID")),
			SubscriptionID:           conversion.StringPtr(subscriptionID),
			ResourceGroupName:        conversion.StringPtr(resourceGroupName),
			KeyVaultName:             conversion.StringPtr(keyVaultName),
			KeyIdentifier:            conversion.StringPtr(os.Getenv("AZURE_KEY_IDENTIFIER")),
			Secret:                   conversion.StringPtr(os.Getenv("AZURE_APP_SECRET")),
			TenantID:                 conversion.StringPtr(os.Getenv("AZURE_TENANT_ID")),
		}
		region = os.Getenv("AZURE_PRIVATE_ENDPOINT_REGION")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckEncryptionAtRestEnvAzure(t) },
		ExternalProviders:        acc.ExternalProvidersOnlyAzapi(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configAzureBasic(projectID, azureKeyVault, region, false),
				Check:  checkBasic(projectID, *azureKeyVault.AzureEnvironment, region, false),
			},
			{
				Config: configAzureBasic(projectID, azureKeyVault, region, true),
			},
			{
				PreConfig:    waitForStatusUpdate,
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkBasic(projectID, *azureKeyVault.AzureEnvironment, region, true),
				),
			},
		},
	})
}

func TestAccEncryptionAtRestPrivateEndpoint_AWS_basic(t *testing.T) {
	resource.Test(t, *basicTestCaseAWS(t))
}

func basicTestCaseAWS(tb testing.TB) *resource.TestCase {
	tb.Helper()
	var (
		projectID = acc.ProjectIDExecution(tb)

		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)

		awsKms = admin.AWSKMSConfiguration{
			Enabled:                  conversion.Pointer(true),
			CustomerMasterKeyID:      conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:                   conversion.StringPtr(conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))),
			RequirePrivateNetworking: conversion.Pointer(false),
		}
		awsKmsPrivateNetworking = admin.AWSKMSConfiguration{
			Enabled:                  conversion.Pointer(true),
			CustomerMasterKeyID:      conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:                   conversion.StringPtr(conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))),
			RequirePrivateNetworking: conversion.Pointer(true),
		}
		region = conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckEncryptionAtRestEnvAWS(tb) },
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigAwsKmsWithRole(projectID, awsIAMRoleName, awsIAMRolePolicyName, &awsKms, false, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(earResourceName, "aws_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(earResourceName, "aws_kms_config.0.require_private_networking", "false"),
				),
			},
			{
				Config: configAWSBasic(projectID, awsIAMRoleName, awsIAMRolePolicyName, &awsKmsPrivateNetworking),
				Check:  checkBasic(projectID, "AWS", region, true),
			},
			{
				Config:            configAWSBasic(projectID, awsIAMRoleName, awsIAMRolePolicyName, &awsKms),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

type errMsgTestCase struct {
	SDKResp *admin.EARPrivateEndpoint
	diags   diag.Diagnostics
}

func TestCheckErrorMessageAndStatus(t *testing.T) {
	var defaultDiags diag.Diagnostics

	testCases := map[string]errMsgTestCase{
		"FAILED status with no error_message": {
			SDKResp: &admin.EARPrivateEndpoint{
				ErrorMessage: nil,
				Status:       admin.PtrString(retrystrategy.RetryStrategyFailedState),
			},
			diags: append(defaultDiags, diag.NewErrorDiagnostic(encryptionatrestprivateendpoint.FailedStatusErrorMessageSummary, "")),
		},
		"FAILED status with error_message": {
			SDKResp: &admin.EARPrivateEndpoint{
				ErrorMessage: admin.PtrString("test err message"),
				Status:       admin.PtrString(retrystrategy.RetryStrategyFailedState),
			},
			diags: append(defaultDiags, diag.NewErrorDiagnostic(encryptionatrestprivateendpoint.FailedStatusErrorMessageSummary, "test err message")),
		},
		"non-empty error_message": {
			SDKResp: &admin.EARPrivateEndpoint{
				ErrorMessage: admin.PtrString("private endpoint was rejected"),
				Status:       admin.PtrString(retrystrategy.RetryStrategyPendingRecreationState),
			},
			diags: append(defaultDiags, diag.NewWarningDiagnostic(encryptionatrestprivateendpoint.NonEmptyErrorMessageFieldSummary, "private endpoint was rejected")),
		},
		"nil error_message": {
			SDKResp: &admin.EARPrivateEndpoint{
				ErrorMessage: nil,
				Status:       admin.PtrString(retrystrategy.RetryStrategyActiveState),
			},
			diags: defaultDiags,
		},
		"empty error_message": {
			SDKResp: &admin.EARPrivateEndpoint{
				ErrorMessage: admin.PtrString(""),
				Status:       admin.PtrString(retrystrategy.RetryStrategyActiveState),
			},
			diags: defaultDiags,
		},
		"pending acceptance status": {
			SDKResp: &admin.EARPrivateEndpoint{
				ErrorMessage: admin.PtrString(""),
				Status:       admin.PtrString(retrystrategy.RetryStrategyPendingAcceptanceState),
			},
			diags: append(defaultDiags, diag.NewWarningDiagnostic(encryptionatrestprivateendpoint.PendingAcceptanceWarnMsgSummary, encryptionatrestprivateendpoint.PendingAcceptanceWarnMsg)),
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			diags := encryptionatrestprivateendpoint.CheckErrorMessageAndStatus(tc.SDKResp)
			assert.Equal(t, tc.diags, diags, "diagnostics did not match expected output")
		})
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cloud_provider"], rs.Primary.Attributes["id"]), nil
	}
}

func configAzureBasic(projectID string, azure *admin.AzureKeyVault, region string, approveWithAzapi bool) string {
	encryptionAtRestConfig := acc.ConfigEARAzureKeyVault(projectID, azure, true, true)
	config := fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_encryption_at_rest_private_endpoint" "test" {
		    project_id = mongodbatlas_encryption_at_rest.test.project_id
		    cloud_provider = "AZURE"
		    region_name = %[2]q
		}

		%[3]s

	`, encryptionAtRestConfig, region, configDS())

	if approveWithAzapi {
		azKeyVaultResourceID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s", azure.GetSubscriptionID(), azure.GetResourceGroupName(), azure.GetKeyVaultName())

		return fmt.Sprintf(`
	%[1]s 
	
	%[2]s 

	resource "azapi_update_resource" "approval" {
		type      = "Microsoft.KeyVault/Vaults/PrivateEndpointConnections@2023-07-01"
		name      = mongodbatlas_encryption_at_rest_private_endpoint.test.private_endpoint_connection_name
		parent_id = %[3]q
	  
		body = jsonencode({
		  properties = {
			privateLinkServiceConnectionState = {
			  description = "Approved via Terraform"
			  status      = "Approved"
			}
		  }
		})
	  }
	`, config, acc.ConfigAzapiProvider(azure.GetSubscriptionID(), azure.GetClientID(), azure.GetSecret(), azure.GetTenantID()), azKeyVaultResourceID)
	}

	return config
}

func checkBasic(projectID, cloudProvider, region string, expectApproved bool) resource.TestCheckFunc {
	expectedStatus := retrystrategy.RetryStrategyPendingAcceptanceState
	if expectApproved {
		expectedStatus = retrystrategy.RetryStrategyActiveState
	}
	attrsSet := []string{"id"}
	if cloudProvider == "AZURE" {
		attrsSet = append(attrsSet, "private_endpoint_connection_name")
	}

	return acc.CheckRSAndDS(
		resourceName,
		admin.PtrString(dataSourceName),
		admin.PtrString(pluralDataSourceName),
		attrsSet,
		map[string]string{
			"project_id":     projectID,
			"status":         expectedStatus,
			"region_name":    region,
			"cloud_provider": cloudProvider,
		})
}

func configAWSBasic(projectID, awsIAMRoleName, awsIAMRolePolicyName string, awsKms *admin.AWSKMSConfiguration) string {
	encryptionAtRestConfig := acc.ConfigAwsKmsWithRole(projectID, awsIAMRoleName, awsIAMRolePolicyName, awsKms, false, true, false)

	config := fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_encryption_at_rest_private_endpoint" "test" {
		    project_id = mongodbatlas_encryption_at_rest.test.project_id
		    cloud_provider = "AWS"
		    region_name = %[2]q
		}

		%[3]s
	`, encryptionAtRestConfig, awsKms.GetRegion(), configDS())

	return config
}

func configEARPrivateEndpointWithTimeout(projectID, region, timeoutConfig string, deleteOnCreateTimeout *bool) string {
	deleteOnCreateTimeoutConfig := ""
	if deleteOnCreateTimeout != nil {
		deleteOnCreateTimeoutConfig = fmt.Sprintf(`
			delete_on_create_timeout = %[1]t
		`, *deleteOnCreateTimeout)
	}

	config := fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest_private_endpoint" "test" {
		    project_id = %[1]q
		    cloud_provider = "AWS"
		    region_name = %[2]q
		    %[3]s
		    %[4]s
		}

		%[5]s

	`, projectID, region, deleteOnCreateTimeoutConfig, timeoutConfig, configDS())

	return config
}

func configDS() string {
	return `
	data "mongodbatlas_encryption_at_rest_private_endpoint" "test" {
		project_id = mongodbatlas_encryption_at_rest_private_endpoint.test.project_id
		cloud_provider = mongodbatlas_encryption_at_rest_private_endpoint.test.cloud_provider
		id = mongodbatlas_encryption_at_rest_private_endpoint.test.id
	}

	data "mongodbatlas_encryption_at_rest_private_endpoints" "test" {
		project_id = mongodbatlas_encryption_at_rest_private_endpoint.test.project_id
		cloud_provider = mongodbatlas_encryption_at_rest_private_endpoint.test.cloud_provider
	}`
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_encryption_at_rest_private_endpoint" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		cloudProvider := rs.Primary.Attributes["cloud_provider"]
		endpointID := rs.Primary.Attributes["id"]
		_, _, err := acc.ConnV2().EncryptionAtRestUsingCustomerKeyManagementApi.GetRestPrivateEndpoint(context.Background(), projectID, cloudProvider, endpointID).Execute()
		if err == nil {
			return fmt.Errorf("EAR private endpoint (%s:%s:%s) still exists", projectID, cloudProvider, endpointID)
		}
	}
	return nil
}

func waitForStatusUpdate() {
	time.Sleep(4 * time.Minute)
}
