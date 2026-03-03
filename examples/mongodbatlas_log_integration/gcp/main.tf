resource "mongodbatlas_project" "project" {
  name   = var.atlas_project_name
  org_id = var.atlas_org_id
}

# Set up cloud provider access in Atlas for GCP
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = mongodbatlas_project.project.id
  provider_name = "GCP"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = mongodbatlas_project.project.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
}

# Set up log integration to export logs to GCS
resource "mongodbatlas_log_integration" "example" {
  project_id  = mongodbatlas_project.project.id
  type        = "GCS_LOG_EXPORT"
  log_types   = ["MONGOD"]
  bucket_name = google_storage_bucket.log_bucket.name
  role_id     = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  prefix_path = "atlas-logs"
}

data "mongodbatlas_log_integration" "example" {
  project_id     = mongodbatlas_log_integration.example.project_id
  integration_id = mongodbatlas_log_integration.example.integration_id
}

data "mongodbatlas_log_integrations" "example" {
  project_id = mongodbatlas_log_integration.example.project_id
  depends_on = [mongodbatlas_log_integration.example]
}

output "log_integration_bucket_name" {
  value = data.mongodbatlas_log_integration.example.bucket_name
}

output "log_integration_ids" {
  value = [for r in data.mongodbatlas_log_integrations.example.results : r.integration_id]
}
