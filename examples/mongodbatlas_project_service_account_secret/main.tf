resource "mongodbatlas_project_service_account" "this" {
  project_id                 = var.project_id
  name                       = "example-project-service-account"
  description                = "Example Project Service Account"
  roles                      = ["GROUP_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_project_service_account_secret" "this" {
  project_id                 = var.project_id
  client_id                  = mongodbatlas_project_service_account.this.client_id
  secret_expires_after_hours = 2160 # 90 days
}

data "mongodbatlas_project_service_account_secret" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
  secret_id  = mongodbatlas_project_service_account_secret.this.secret_id
}

output "secret_id" {
  value = mongodbatlas_project_service_account_secret.this.secret_id
}

output "secret" {
  sensitive = true
  value     = mongodbatlas_project_service_account_secret.this.secret
}

output "secret_expires_at" {
  value = data.mongodbatlas_project_service_account_secret.this.expires_at
}
