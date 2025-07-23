resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.project_id
  provider_name = "AZURE"

  azure_config {
    atlas_azure_app_id   = var.azure_atlas_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  azure {
    atlas_azure_app_id   = var.azure_atlas_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}
resource "mongodbatlas_federated_database_instance" "azure_example" {
  project_id = var.project_id
  name       = var.federated_instance_name
  cloud_provider_config {
    azure {
      role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
    }
  }
}