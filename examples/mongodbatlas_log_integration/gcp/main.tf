provider "google" {
  project = var.gcp_project_id
  # region  = "us-central1"
}

# Set up log integration with authorized IAM role
resource "mongodbatlas_log_integration" "gcs" {
  project_id  = var.project_id

  type        = "GCS_LOG_EXPORT"
  log_types   = ["MONGOS_AUDIT"]
  prefix_path = "logs/mongodb/"
  role_id      = mongodbatlas_cloud_provider_access_authorization.gcp_auth.role_id
  bucket_name  = google_storage_bucket.log_bucket.name
}

# Set up cloud provider access in Atlas for GCP
resource "mongodbatlas_cloud_provider_access_setup" "gcp_setup" {
  project_id    = var.project_id
  provider_name = "GCP"
}

resource "mongodbatlas_cloud_provider_access_authorization" "gcp_auth" {
  project_id =  mongodbatlas_cloud_provider_access_setup.gcp_setup.project_id
  role_id    =  mongodbatlas_cloud_provider_access_setup.gcp_setup.role_id
}

# Create GCS bucket
resource "google_storage_bucket" "log_bucket" {
  name          = var.gcs_bucket_name
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_iam_member" "atlas_access" {
  bucket = google_storage_bucket.log_bucket.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${mongodbatlas_cloud_provider_access_authorization.gcp_auth.gcp[0].service_account_for_atlas}"
}
