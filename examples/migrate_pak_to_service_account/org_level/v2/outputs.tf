output "service_account_secret" {
  description = "The secret value for the Service Account. Returned only when the secret is created."
  value       = mongodbatlas_service_account_secret.this.secret
  sensitive   = true
}
