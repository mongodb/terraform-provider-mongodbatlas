output "atlas_role_id" {
  description = "The MongoDB Atlas cloud provider access role ID"
  value       = mongodbatlas_cloud_provider_access_authorization.this.role_id
}

output "gcp_service_account_email" {
  description = "The GCP service account email created by MongoDB Atlas for accessing KMS"
  value       = mongodbatlas_cloud_provider_access_authorization.this.gcp[0].service_account_for_atlas
}

output "kms_key_ring_id" {
  description = "The full ID of the created GCP KMS key ring"
  value       = google_kms_key_ring.key_ring.id
}

output "kms_crypto_key_id" {
  description = "The full ID of the created GCP KMS crypto key"
  value       = google_kms_crypto_key.crypto_key.id
}

output "kms_key_version_resource_id" {
  description = "The resource ID of the primary key version used for Atlas encryption"
  value       = google_kms_crypto_key.crypto_key.primary[0].name
}
