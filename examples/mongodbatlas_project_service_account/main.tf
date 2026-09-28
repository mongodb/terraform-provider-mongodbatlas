# Create a Project Service Account without an Atlas-generated secret, then create the first secret
# explicitly with mongodbatlas_project_service_account_secret so this configuration owns it.

resource "mongodbatlas_project_service_account" "this" {
  project_id             = var.project_id
  name                   = "example-project-service-account"
  description            = "Example Project Service Account"
  roles                  = ["GROUP_READ_ONLY"]
  without_initial_secret = true
}

resource "mongodbatlas_project_service_account_secret" "this" {
  project_id                 = var.project_id
  client_id                  = mongodbatlas_project_service_account.this.client_id
  secret_expires_after_hours = 2160 # 90 days
}

data "mongodbatlas_project_service_account" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
}

data "mongodbatlas_project_service_accounts" "this" {
  project_id = var.project_id
  depends_on = [mongodbatlas_project_service_account.this]
}

output "service_account_client_id" {
  description = "The Client ID of the Project Service Account. Use it with a secret to authenticate."
  value       = mongodbatlas_project_service_account.this.client_id
}

output "service_account_name" {
  description = "The name of the Project Service Account, read from the data source."
  value       = data.mongodbatlas_project_service_account.this.name
}

output "secret_id" {
  description = "The ID of the Project Service Account secret."
  value       = mongodbatlas_project_service_account_secret.this.secret_id
}

output "secret" {
  description = "The secret value for the Project Service Account. Returned only when the secret is created."
  sensitive   = true
  value       = mongodbatlas_project_service_account_secret.this.secret
}

output "service_accounts_results" {
  description = "All Project Service Accounts in the project, read from the plural data source."
  value       = data.mongodbatlas_project_service_accounts.this.results
}
