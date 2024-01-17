resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.project_id
  provider_name = var.cloud_provider_access_name
  azure_config {
    atlas_azure_app_id   = var.atlas_azure_app_id
    service_principal_id = azuread_service_principal.example.object_id
    tenant_id            = var.azure_tenant_id
  }

}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  azure {
    atlas_azure_app_id   = var.atlas_azure_app_id
    service_principal_id = azuread_service_principal.example.object_id
    tenant_id            = var.azure_tenant_id
  }
}
