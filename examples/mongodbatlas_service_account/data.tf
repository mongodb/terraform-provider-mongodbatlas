# Read the Service Account and the organization's Service Accounts back from Atlas.

data "mongodbatlas_service_account" "this" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.this.client_id
  depends_on = [mongodbatlas_service_account_secret.this]
}

data "mongodbatlas_service_accounts" "this" {
  org_id     = var.org_id
  depends_on = [mongodbatlas_service_account_secret.this]
}

output "service_account_client_id" {
  description = "The Client ID of the Service Account. Use it with a secret to authenticate."
  value       = mongodbatlas_service_account.this.client_id
}

output "service_account_name" {
  description = "The name of the Service Account, read from the data source."
  value       = data.mongodbatlas_service_account.this.name
}

output "service_accounts_results" {
  description = "All Service Accounts in the organization, read from the plural data source."
  value       = data.mongodbatlas_service_accounts.this.results
}
