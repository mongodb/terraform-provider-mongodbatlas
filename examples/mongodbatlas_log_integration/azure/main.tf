# Set up cloud provider access in Atlas for Azure
resource "mongodbatlas_cloud_provider_access_setup" "setup" {
  project_id    = mongodbatlas_project.project.id
  provider_name = "AZURE"

  azure_config {
    atlas_azure_app_id   = var.atlas_azure_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth" {
  project_id = mongodbatlas_project.project.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup.role_id

  azure {
    atlas_azure_app_id   = var.atlas_azure_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}

resource "mongodbatlas_log_integration" "example" {
  project_id             = mongodbatlas_project.project.id
  type                   = "AZURE_LOG_EXPORT"
  log_types              = ["MONGOD"]
  role_id                = mongodbatlas_cloud_provider_access_authorization.auth.role_id
  storage_account_name   = azurerm_storage_account.log_storage.name
  storage_container_name = azurerm_storage_container.log_container.name
  prefix_path            = "atlas-logs"
}
