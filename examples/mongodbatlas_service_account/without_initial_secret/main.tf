# Creates a Service Account without an initial secret by omitting `secret_expires_after_hours`.
# This is the recommended flow: the bootstrap secret Atlas would otherwise generate is never created.
# Add a secret later with the `mongodbatlas_service_account_secret` resource or the Atlas API.
resource "mongodbatlas_service_account" "this" {
  org_id      = var.org_id
  name        = "example-service-account"
  description = "Example Service Account without an initial secret"
  roles       = ["ORG_READ_ONLY"]
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

output "service_account_secrets" {
  description = "Secrets on the Service Account. Empty after create, since no initial secret is generated."
  value       = mongodbatlas_service_account.this.secrets
}

output "service_accounts_results" {
  value = data.mongodbatlas_service_accounts.this.results
}
