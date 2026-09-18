# Create a Service Account without an Atlas-generated secret, then create the first secret
# explicitly with mongodbatlas_service_account_secret so this configuration owns it.

resource "mongodbatlas_service_account" "this" {
  org_id                 = var.org_id
  name                   = "example-service-account"
  description            = "Example Service Account"
  roles                  = ["ORG_READ_ONLY"]
  without_initial_secret = true
}

resource "mongodbatlas_service_account_secret" "this" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.this.client_id
  secret_expires_after_hours = 2160 # 90 days
}

output "service_account_client_id" {
  description = "The Client ID of the Service Account. Use it with a secret to authenticate."
  value       = mongodbatlas_service_account.this.client_id
}

output "secret_id" {
  description = "The ID of the Service Account secret."
  value       = mongodbatlas_service_account_secret.this.secret_id
}

output "secret" {
  description = "The secret value for the Service Account. Returned only when the secret is created."
  sensitive   = true
  value       = mongodbatlas_service_account_secret.this.secret
}
