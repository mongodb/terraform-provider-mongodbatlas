output "atlas_role_id" {
  description = "The MongoDB Atlas cloud provider access role ID"
  value       = mongodbatlas_cloud_provider_access_authorization.this.role_id
}

output "gcp_service_account_email" {
  description = "The GCP service account email created by MongoDB Atlas"
  value       = mongodbatlas_cloud_provider_access_authorization.this.gcp[0].service_account_for_atlas
}
