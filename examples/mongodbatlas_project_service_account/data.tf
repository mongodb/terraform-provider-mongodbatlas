# Read the Project Service Account and the project's Service Accounts back from Atlas.

data "mongodbatlas_project_service_account" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
  depends_on = [mongodbatlas_project_service_account_secret.this]
}

data "mongodbatlas_project_service_accounts" "this" {
  project_id = var.project_id
  depends_on = [mongodbatlas_project_service_account_secret.this]
}

output "service_account_client_id" {
  description = "The Client ID of the Project Service Account. Use it with a secret to authenticate."
  value       = mongodbatlas_project_service_account.this.client_id
}

output "service_account_name" {
  description = "The name of the Project Service Account, read from the data source."
  value       = data.mongodbatlas_project_service_account.this.name
}

output "service_accounts_results" {
  description = "All Project Service Accounts in the project, read from the plural data source."
  value       = data.mongodbatlas_project_service_accounts.this.results
}
