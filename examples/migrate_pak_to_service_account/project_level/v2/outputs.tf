output "project_service_account_secret" {
  description = "The secret value for the Project Service Account. Returned only when the secret is created."
  value       = mongodbatlas_project_service_account_secret.this.secret
  sensitive   = true
}
