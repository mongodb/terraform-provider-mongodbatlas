package encryptionatrestprivateendpoint_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

const (
	resourceName    = "mongodbatlas_encryption_at_rest_private_endpoint.test"
	dataSourceName  = "data.mongodbatlas_encryption_at_rest_private_endpoint.test"
	earResourceName = "mongodbatlas_encryption_at_rest.test"
)

func TestAccEncryptionAtRestPrivateEndpoint_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	acc.SkipTestForCI(tb) // needs Azure configuration
	var (
		projectID     = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		azureKeyVault = &admin.AzureKeyVault{
			Enabled:                  conversion.Pointer(true),
			RequirePrivateNetworking: conversion.Pointer(true),
			AzureEnvironment:         conversion.StringPtr("AZURE"),
			ClientID:                 conversion.StringPtr(os.Getenv("AZURE_CLIENT_ID")),
			SubscriptionID:           conversion.StringPtr(os.Getenv("AZURE_SUBSCRIPTION_ID")),
			ResourceGroupName:        conversion.StringPtr(os.Getenv("AZURE_RESOURCE_GROUP_NAME")),
			KeyVaultName:             conversion.StringPtr(os.Getenv("AZURE_KEY_VAULT_NAME")),
			KeyIdentifier:            conversion.StringPtr(os.Getenv("AZURE_KEY_IDENTIFIER")),
			Secret:                   conversion.StringPtr(os.Getenv("AZURE_SECRET")),
			TenantID:                 conversion.StringPtr(os.Getenv("AZURE_TENANT_ID")),
		}
		region = os.Getenv("AZURE_PRIVATE_ENDPOINT_REGION")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(tb); acc.PreCheckEncryptionAtRestEnvAzure(tb); acc.PreCheckPreviewFlag(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configPrivateEndpointAzureBasic(projectID, azureKeyVault, region),
				Check:  checkPrivateEndpointAzureBasic(projectID, azureKeyVault, region),
			},
			{
				Config:            configPrivateEndpointAzureBasic(projectID, azureKeyVault, region),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func TestAccEncryptionAtRestPrivateEndpoint_transitionPublicToPrivateNetwork(t *testing.T) {
	acc.SkipTestForCI(t) // needs Azure configuration
	var (
		projectID     = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		azureKeyVault = &admin.AzureKeyVault{
			Enabled:                  conversion.Pointer(true),
			RequirePrivateNetworking: conversion.Pointer(true),
			AzureEnvironment:         conversion.StringPtr("AZURE"),
			ClientID:                 conversion.StringPtr(os.Getenv("AZURE_CLIENT_ID")),
			SubscriptionID:           conversion.StringPtr(os.Getenv("AZURE_SUBSCRIPTION_ID")),
			ResourceGroupName:        conversion.StringPtr(os.Getenv("AZURE_RESOURCE_GROUP_NAME")),
			KeyVaultName:             conversion.StringPtr(os.Getenv("AZURE_KEY_VAULT_NAME")),
			KeyIdentifier:            conversion.StringPtr(os.Getenv("AZURE_KEY_IDENTIFIER")),
			Secret:                   conversion.StringPtr(os.Getenv("AZURE_SECRET")),
			TenantID:                 conversion.StringPtr(os.Getenv("AZURE_TENANT_ID")),
		}
		region = os.Getenv("AZURE_PRIVATE_ENDPOINT_REGION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckEncryptionAtRestEnvAzure(t); acc.PreCheckPreviewFlag(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigEARAzureKeyVault(projectID, azureKeyVault, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(earResourceName, "azure_key_vault_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(earResourceName, "azure_key_vault_config.0.require_private_networking", "false"),
				),
			},
			{
				Config: configPrivateEndpointAzureBasic(projectID, azureKeyVault, region),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(earResourceName, "azure_key_vault_config.0.require_private_networking", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status", "PENDING_ACCEPTANCE"),
				),
			},
		},
	})
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

func configPrivateEndpointAzureBasic(projectID string, azure *admin.AzureKeyVault, region string) string {
	encryptionAtRestConfig := acc.ConfigEARAzureKeyVault(projectID, azure, true)
	return fmt.Sprintf(`
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
	`, encryptionAtRestConfig, region)
}

func checkPrivateEndpointAzureBasic(projectID string, azure *admin.AzureKeyVault, region string) resource.TestCheckFunc {
	return acc.CheckRSAndDS(
		resourceName,
		admin.PtrString(dataSourceName),
		nil,
		[]string{"id", "private_endpoint_connection_name"},
		map[string]string{
			"project_id":     projectID,
			"status":         retrystrategy.RetryStrategyPendingAcceptanceState,
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
