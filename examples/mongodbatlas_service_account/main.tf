resource "mongodbatlas_service_account" "this" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account"
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

data "mongodbatlas_service_account" "this" {
  org_id    = var.org_id
  client_id = mongodbatlas_service_account.this.client_id
}

data "mongodbatlas_service_accounts" "this" {
  org_id = var.org_id
}

output "service_account_client_id" {
  value = mongodbatlas_service_account.this.client_id
}

output "service_account_name" {
  value = data.mongodbatlas_service_account.this.name
}

output "service_account_first_secret" {
  description = "The secret value of the first secret created with the service account. Only available after initial creation."
  value       = try(mongodbatlas_service_account.this.secrets[0].secret, null)
  sensitive   = true
}

output "service_accounts_results" {
  value = data.mongodbatlas_service_accounts.this.results
}
