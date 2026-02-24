provider "azurerm" {
  features {}

  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}

# Set up log integration with authorized IAM role
resource "mongodbatlas_log_integration" "azure" {
  project_id  = var.project_id
  type        = "AZURE_LOG_EXPORT"
  log_types   = ["MONGOS_AUDIT"]
  prefix_path            = "logs/mongodb/"
  service_principal_id   = mongodbatlas_cloud_provider_access_authorization.azure_auth.role_id
  storage_account_name   = azurerm_storage_account.log_storage.name
  storage_container_name = azurerm_storage_container.log_container.name
}

// Set up cloud provider access in Atlas for Azure
resource "mongodbatlas_cloud_provider_access_setup" "azure_setup" {
  project_id    = var.project_id
  provider_name = "AZURE"

  azure_config {
    atlas_azure_app_id   = var.atlas_azure_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}

resource "mongodbatlas_cloud_provider_access_authorization" "azure_auth" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.azure_setup.role_id

  azure {
    atlas_azure_app_id   = var.atlas_azure_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}

// Create Azure resource group (Needs Contributor role to create a resource group)
resource "azurerm_resource_group" "log_rg" {
  name     = var.azure_resource_group_name
  location = "East US"
}

// Create Azure storage account
resource "azurerm_storage_account" "log_storage" {
  name                     = var.azure_storage_account_name
  resource_group_name      = azurerm_resource_group.log_rg.name
  location                 = azurerm_resource_group.log_rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

// Create Azure storage container
resource "azurerm_storage_container" "log_container" {
  name                  = var.azure_storage_container_name
  storage_account_id    = azurerm_storage_account.log_storage.id
}
