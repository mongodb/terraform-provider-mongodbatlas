resource "mongodbatlas_cloud_provider_access_setup" "this" {
  project_id    = var.atlas_project_id
  provider_name = "GCP"
}

resource "mongodbatlas_cloud_provider_access_authorization" "this" {
  project_id = var.atlas_project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.this.role_id
}

