resource "mongodbatlas_project_service_account" "this" {
  project_id                 = var.project_id
  name                       = "example-project-service-account"
  description                = "Example Project Service Account"
  roles                      = ["GROUP_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

data "mongodbatlas_project_service_account" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
}

data "mongodbatlas_project_service_accounts" "this" {
  project_id = var.project_id
}

output "service_account_client_id" {
  value = mongodbatlas_project_service_account.this.client_id
}

output "service_account_name" {
  value = data.mongodbatlas_project_service_account.this.name
}

output "service_account_first_secret" {
  description = "The secret value of the first secret created with the Project Service Account. Available only immediately after initial creation."
  value       = try(mongodbatlas_project_service_account.this.secrets[0].secret, null)
  sensitive   = true
}

output "service_accounts_results" {
  value = data.mongodbatlas_project_service_accounts.this.results
}
