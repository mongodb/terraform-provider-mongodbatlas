resource "mongodbatlas_service_account" "this" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account"
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_service_account_secret" "this" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.this.client_id
  secret_expires_after_hours = 2160 # 90 days
}

data "mongodbatlas_service_account_secret" "this" {
  org_id    = var.org_id
  client_id = mongodbatlas_service_account.this.client_id
  secret_id = mongodbatlas_service_account_secret.this.secret_id
}

output "secret_id" {
  value = mongodbatlas_service_account_secret.this.secret_id
}

output "secret" {
  sensitive = true
  value     = mongodbatlas_service_account_secret.this.secret
}

output "secret_expires_at" {
  value = data.mongodbatlas_service_account_secret.this.expires_at
}
