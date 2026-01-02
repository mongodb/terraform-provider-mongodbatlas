resource "mongodbatlas_service_account" "example" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account"
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

data "mongodbatlas_service_account" "example" {
  org_id    = var.org_id
  client_id = mongodbatlas_service_account.example.client_id
}

data "mongodbatlas_service_accounts" "example" {
  org_id = var.org_id
}

output "service_account_client_id" {
  value = mongodbatlas_service_account.example.client_id
}

output "service_account_name" {
  value = data.mongodbatlas_service_account.example.name
}

output "service_accounts_results" {
  value = data.mongodbatlas_service_accounts.example.results
}
