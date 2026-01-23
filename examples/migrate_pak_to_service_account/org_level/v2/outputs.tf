output "service_account_first_secret" {
  description = "The secret value of the first secret created with the Service Account. Only available after initial creation."
  value       = try(mongodbatlas_service_account.this.secrets[0].secret, null)
  sensitive   = true
}
