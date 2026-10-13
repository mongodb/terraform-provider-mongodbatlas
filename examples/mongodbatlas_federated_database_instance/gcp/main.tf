# Example of using MongoDB Atlas Federated Database Instance.
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.project_id
  provider_name = "GCP"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
}

resource "google_storage_bucket" "example" {
  name          = var.gcs_bucket_name
  location      = var.gcs_bucket_location
  force_destroy = true
}

resource "google_storage_bucket_iam_member" "atlas_access" {
  bucket = google_storage_bucket.example.name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${mongodbatlas_cloud_provider_access_authorization.auth_role.gcp[0].service_account_for_atlas}"
}

resource "mongodbatlas_federated_database_instance" "gcp_example" {
  project_id = var.project_id
  name       = var.federated_instance_name

  cloud_provider_config {
    gcp {
      role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
    }
  }

  # Atlas Data Federation must run on the same cloud provider as the data store.
  data_process_region {
    cloud_provider = "GCP"
    region         = var.data_process_region
  }

  storage_databases {
    name = "VirtualDatabase0"
    collections {
      name = "VirtualCollection0"
      data_sources {
        store_name = google_storage_bucket.example.name
        path       = "/"
      }
    }
  }

  storage_stores {
    name     = google_storage_bucket.example.name
    provider = "gcs"
    bucket   = google_storage_bucket.example.name
    prefix   = "/"
    region   = var.gcs_bucket_location
  }
  depends_on = [google_storage_bucket_iam_member.atlas_access]
}

output "gcp_service_account" {
  description = "Email address of the Google Cloud Platform (GCP) service account created by Atlas, which should be authorized to allow Atlas to access Google Cloud Storage."
  value       = mongodbatlas_federated_database_instance.gcp_example.cloud_provider_config[0].gcp[0].gcp_service_account
}

output "gcs_bucket_name" {
  description = "Name of the Google Cloud Storage bucket the federated database instance queries."
  value       = google_storage_bucket.example.name
}
