# GCS bucket for storing Atlas logs
resource "google_storage_bucket" "log_bucket" {
  name          = var.gcs_bucket_name
  location      = var.gcs_bucket_location
  force_destroy = true
}

# Grant the Atlas-managed service account object admin access to the bucket
resource "google_storage_bucket_iam_member" "atlas_access" {
  bucket = google_storage_bucket.log_bucket.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${mongodbatlas_cloud_provider_access_authorization.auth_role.gcp[0].service_account_for_atlas}"
}
