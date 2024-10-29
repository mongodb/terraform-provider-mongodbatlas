package encryptionatrestprivateendpoint_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20241023001/admin"

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

func TestAccEncryptionAtRestPrivateEndpoint_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
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
				Config: configAzureBasic(projectID, azureKeyVault, region, false),
				Check:  checkAzureBasic(projectID, azureKeyVault, region, false),
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
				Check:  checkAzureBasic(projectID, azureKeyVault, region, false),
			},
			{
				Config: configAzureBasic(projectID, azureKeyVault, region, true),
			},
			{
				PreConfig:    waitForStatusUpdate,
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAzureBasic(projectID, azureKeyVault, region, true),
				),
			},
		},
	})
}

func TestAccEncryptionAtRestPrivateEndpoint_transitionPublicToPrivateNetwork(t *testing.T) {
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckEncryptionAtRestEnvAzure(t) },
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(earResourceName, "azure_key_vault_config.0.require_private_networking", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status", "PENDING_ACCEPTANCE"),
				),
			},
		},
	})
}

func TestAccEncryptionAtRest_azure_requirePrivateNetworking(t *testing.T) {
	var (
		projectID = os.Getenv("MONGODB_ATLAS_PROJECT_EAR_PE_ID")

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
			RequirePrivateNetworking: conversion.Pointer(true),
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
		PreCheck: func() {
			acc.PreCheckEncryptionAtRestPrivateEndpoint(t)
			acc.PreCheckEncryptionAtRestEnvAzureWithUpdate(t)
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigEARAzureKeyVault(projectID, &azureKeyVault, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(earResourceName),
					resource.TestCheckResourceAttr(earResourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(earResourceName, "azure_key_vault_config.0", azureKeyVaultAttrMap),

					resource.TestCheckResourceAttr(earDatasourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(earDatasourceName, "azure_key_vault_config", azureKeyVaultAttrMap),
				),
			},
			{
				Config: acc.ConfigEARAzureKeyVault(projectID, &azureKeyVaultUpdated, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(earResourceName),
					resource.TestCheckResourceAttr(earResourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(earResourceName, "azure_key_vault_config.0", azureKeyVaultUpdatedAttrMap),

					resource.TestCheckResourceAttr(earDatasourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(earDatasourceName, "azure_key_vault_config.", azureKeyVaultUpdatedAttrMap),
				),
			},
			{
				ResourceName:      earResourceName,
				ImportStateIdFunc: acc.EARImportStateIDFunc(earResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// "azure_key_vault_config.0.secret" is a sensitive value not returned by the API
				ImportStateVerifyIgnore: []string{"azure_key_vault_config.0.secret"},
			},
		},
	})
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

		data "mongodbatlas_encryption_at_rest_private_endpoint" "test" {
		    project_id = mongodbatlas_encryption_at_rest_private_endpoint.test.project_id
			cloud_provider = mongodbatlas_encryption_at_rest_private_endpoint.test.cloud_provider
		    id = mongodbatlas_encryption_at_rest_private_endpoint.test.id
		}

		data "mongodbatlas_encryption_at_rest_private_endpoints" "test" {
		    project_id = mongodbatlas_encryption_at_rest_private_endpoint.test.project_id
			cloud_provider = mongodbatlas_encryption_at_rest_private_endpoint.test.cloud_provider
		}

	`, encryptionAtRestConfig, region)

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

func checkAzureBasic(projectID string, azure *admin.AzureKeyVault, region string, expectApproved bool) resource.TestCheckFunc {
	expectedStatus := retrystrategy.RetryStrategyPendingAcceptanceState
	if expectApproved {
		expectedStatus = retrystrategy.RetryStrategyActiveState
	}

	return acc.CheckRSAndDS(
		resourceName,
		admin.PtrString(dataSourceName),
		admin.PtrString(pluralDataSourceName),
		[]string{"id", "private_endpoint_connection_name"},
		map[string]string{
			"project_id":     projectID,
			"status":         expectedStatus,
			"region_name":    region,
			"cloud_provider": *azure.AzureEnvironment,
		})
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_encryption_at_rest_private_endpoint" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		cloudProvider := rs.Primary.Attributes["cloud_provider"]
		endpointID := rs.Primary.Attributes["id"]
		_, _, err := acc.ConnV2().EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRestPrivateEndpoint(context.Background(), projectID, cloudProvider, endpointID).Execute()
		if err == nil {
			return fmt.Errorf("EAR private endpoint (%s:%s:%s) still exists", projectID, cloudProvider, endpointID)
		}
	}
	return nil
}

func waitForStatusUpdate() {
	time.Sleep(4 * time.Minute)
}
