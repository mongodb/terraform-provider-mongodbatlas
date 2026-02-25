provider "azurerm" {
  features {}

  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}

# Set up cloud provider access in Atlas for Azure
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

# Set up log integration with Azure
resource "mongodbatlas_log_integration" "example" {
  project_id  = var.project_id
  type        = "AZURE_LOG_EXPORT"
  log_types   = [""MONGOD", "MONGOS", "MONGOD_AUDIT", "MONGOS_AUDIT""]
  prefix_path            = "logs/mongodb/"
  service_principal_id   = mongodbatlas_cloud_provider_access_authorization.azure_auth.role_id
  storage_account_name   = azurerm_storage_account.log_storage.name
  storage_container_name = azurerm_storage_container.log_container.name
}
