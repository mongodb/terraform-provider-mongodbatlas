# Output to capture the secret (new)
output "service_account_first_secret" {
  description = "The secret value of the first secret created with the service account. Only available after initial creation."
  value       = try(mongodbatlas_service_account.example.secrets[0].secret, null)
  sensitive   = true
}