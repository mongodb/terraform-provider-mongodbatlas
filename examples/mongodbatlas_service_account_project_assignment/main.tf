resource "mongodbatlas_service_account" "this" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account"
  roles                      = ["ORG_MEMBER"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_service_account_project_assignment" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_service_account.this.client_id
  roles      = ["GROUP_READ_ONLY"]
}

data "mongodbatlas_service_account_project_assignment" "this" {
  project_id = mongodbatlas_service_account_project_assignment.this.project_id
  client_id  = mongodbatlas_service_account_project_assignment.this.client_id
}

data "mongodbatlas_service_account_project_assignments" "this" {
  org_id    = var.org_id
  client_id = mongodbatlas_service_account.this.client_id
}

output "service_account_project_roles" {
  value = data.mongodbatlas_service_account_project_assignment.this.roles
}

output "service_account_assigned_projects" {
  value = data.mongodbatlas_service_account_project_assignments.this.results
}
