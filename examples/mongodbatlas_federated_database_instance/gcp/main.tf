# Example of using MongoDB Atlas Federated Database Instance.
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.project_id
  provider_name = "GCP"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
}

resource "mongodbatlas_federated_database_instance" "gcp_example" {
  project_id = var.project_id
  name       = var.federated_instance_name

  cloud_provider_config {
    gcp {
      role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
    }
  }
}

output "gcp_service_account" {
  description = "Email address of the Google Cloud Platform (GCP) service account created by Atlas, which should be authorized to allow Atlas to access Google Cloud Storage."
  value       = mongodbatlas_federated_database_instance.gcp_example.cloud_provider_config[0].gcp[0].gcp_service_account
}
