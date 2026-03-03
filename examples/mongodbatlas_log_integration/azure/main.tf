resource "mongodbatlas_project" "project" {
  name   = var.atlas_project_name
  org_id = var.atlas_org_id
}

# Set up cloud provider access in Atlas for Azure
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = mongodbatlas_project.project.id
  provider_name = "AZURE"

  azure_config {
    atlas_azure_app_id   = var.atlas_azure_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = mongodbatlas_project.project.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  azure {
    atlas_azure_app_id   = var.atlas_azure_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}

# Set up log integration to export logs to Azure Blob Storage
resource "mongodbatlas_log_integration" "example" {
  project_id             = mongodbatlas_project.project.id
  type                   = "AZURE_LOG_EXPORT"
  log_types              = ["MONGOD"]
  role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  storage_account_name   = azurerm_storage_account.log_storage.name
  storage_container_name = azurerm_storage_container.log_container.name
  prefix_path            = "atlas-logs"
}

data "mongodbatlas_log_integration" "example" {
  project_id     = mongodbatlas_log_integration.example.project_id
  integration_id = mongodbatlas_log_integration.example.integration_id
}

data "mongodbatlas_log_integrations" "example" {
  project_id = mongodbatlas_log_integration.example.project_id
  depends_on = [mongodbatlas_log_integration.example]
}

output "log_integration_storage_container_name" {
  value = data.mongodbatlas_log_integration.example.storage_container_name
}

output "log_integration_ids" {
  value = [for r in data.mongodbatlas_log_integrations.example.results : r.integration_id]
}
