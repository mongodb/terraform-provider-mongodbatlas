# Outputs maintain backward compatibility with v1
# Internal implementation changed but output names stay the same

output "user_id" {
  description = "User's Atlas ID"
  value       = data.mongodbatlas_cloud_user_org_assignment.this.user_id
}

output "username" {
  description = "User's email address"
  value       = data.mongodbatlas_cloud_user_org_assignment.this.username
}

output "email_address" {
  description = "User's email address"
  value       = data.mongodbatlas_cloud_user_org_assignment.this.username
}

output "first_name" {
  description = "User's first name"
  value       = data.mongodbatlas_cloud_user_org_assignment.this.first_name
}

output "last_name" {
  description = "User's last name"
  value       = data.mongodbatlas_cloud_user_org_assignment.this.last_name
}

output "org_roles" {
  description = "User's organization roles"
  value       = data.mongodbatlas_cloud_user_org_assignment.this.roles.org_roles
}

output "project_roles" {
  description = "User's project roles (filtered by project_id, if provided)"
  value       = local.user_project_roles
}
