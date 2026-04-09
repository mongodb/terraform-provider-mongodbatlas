package acc

import "fmt"

func ConfigSetupAzure(projectID, atlasAzureAppID, servicePrincipalID, tenantID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_provider_access_setup" "test" {
		project_id = %[1]q
		provider_name = "AZURE"
		azure_config {
			atlas_azure_app_id = %[2]q
			service_principal_id = %[3]q
			tenant_id = %[4]q
		}
	 }
	`, projectID, atlasAzureAppID, servicePrincipalID, tenantID)
}

// ConfigAzureCloudProviderAccess returns Terraform config for mongodbatlas_cloud_provider_access_setup
// and mongodbatlas_cloud_provider_access_authorization resources with the "azure_setup" and "azure_auth" resource names.
func ConfigAzureCloudProviderAccess(projectID, atlasAzureAppID, servicePrincipalID, tenantID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cloud_provider_access_setup" "azure_setup" {
			project_id    = %[1]q
			provider_name = "AZURE"

			azure_config {
				atlas_azure_app_id   = %[2]q
				service_principal_id = %[3]q
				tenant_id            = %[4]q
			}
		}

		resource "mongodbatlas_cloud_provider_access_authorization" "azure_auth" {
			project_id = %[1]q
			role_id    = mongodbatlas_cloud_provider_access_setup.azure_setup.role_id

			azure {
				atlas_azure_app_id   = %[2]q
				service_principal_id = %[3]q
				tenant_id            = %[4]q
			}
		}
	`, projectID, atlasAzureAppID, servicePrincipalID, tenantID)
}

// ConfigAzureStorageResources returns Terraform config for Azure resource group, storage account,
// storage container, and role assignment resources. The prefix parameter is used as a prefix for the
// Terraform resource names (e.g., "blob" produces "azurerm_resource_group.blob_rg").
func ConfigAzureStorageResources(prefix, resourceGroupName, storageAccountName, storageContainerName, servicePrincipalID string) string {
	return fmt.Sprintf(`
		resource "azurerm_resource_group" "%[1]s_rg" {
			name     = %[2]q
			location = "East US"
		}

		resource "azurerm_storage_account" "%[1]s_storage" {
			name                     = %[3]q
			resource_group_name      = azurerm_resource_group.%[1]s_rg.name
			location                 = azurerm_resource_group.%[1]s_rg.location
			account_tier             = "Standard"
			account_replication_type = "LRS"
		}

		resource "azurerm_storage_container" "%[1]s_container" {
			name               = %[4]q
			storage_account_id = azurerm_storage_account.%[1]s_storage.id
		}

		resource "azurerm_role_assignment" "%[1]s_contributor" {
			scope                = azurerm_storage_account.%[1]s_storage.id
			role_definition_name = "Storage Blob Data Contributor"
			principal_id         = %[5]q
		}
	`, prefix, resourceGroupName, storageAccountName, storageContainerName, servicePrincipalID)
}
