# Set up cloud provider access in Atlas for GCP
resource "mongodbatlas_cloud_provider_access_setup" "setup" {
  project_id    = mongodbatlas_project.project.id
  provider_name = "GCP"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth" {
  project_id = mongodbatlas_project.project.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup.role_id
}

resource "mongodbatlas_log_integration" "example" {
  project_id  = mongodbatlas_project.project.id
  type        = "GCS_LOG_EXPORT"
  log_types   = ["MONGOD"]
  bucket_name = google_storage_bucket.log_bucket.name
  role_id     = mongodbatlas_cloud_provider_access_authorization.auth.role_id
  prefix_path = "atlas-logs"
}
