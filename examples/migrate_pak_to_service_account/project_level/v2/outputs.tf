output "project_service_account_first_secret" {
  description = "The secret value of the first secret created with the project service account. Only available after initial creation."
  value       = try(mongodbatlas_project_service_account.this.secrets[0].secret, null)
  sensitive   = true
}
