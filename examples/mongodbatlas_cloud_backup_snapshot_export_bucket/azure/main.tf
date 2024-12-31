resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.project_id
  provider_name = "AZURE"
  azure_config {
    atlas_azure_app_id   = var.azure_atlas_app_id
    service_principal_id = azuread_service_principal.mongo.id
    tenant_id            = var.tenant_id
  }
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  azure {
    atlas_azure_app_id   = var.azure_atlas_app_id
    service_principal_id = azuread_service_principal.mongo.id
    tenant_id            = var.tenant_id
  }
}


resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id     = var.project_id
  bucket_name    = azurerm_storage_container.test_storage_container.name
  cloud_provider = "AZURE"
  service_url    = azurerm_storage_account.test_storage_account.primary_blob_endpoint
  role_id        = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
}