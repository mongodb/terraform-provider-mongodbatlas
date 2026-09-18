# Create a Service Account without an Atlas-generated secret.
# Create secrets separately with mongodbatlas_service_account_secret, so this configuration owns them.

resource "mongodbatlas_service_account" "this" {
  org_id                 = var.org_id
  name                   = "example-service-account"
  description            = "Example Service Account"
  roles                  = ["ORG_READ_ONLY"]
  without_initial_secret = true
}

output "service_account_client_id" {
  description = "The Client ID of the Service Account. Use it with a secret to authenticate."
  value       = mongodbatlas_service_account.this.client_id
}
