output "user_id" {
  description = "User's Atlas ID"
  value       = data.mongodbatlas_atlas_user.this.user_id
}

output "username" {
  description = "User's username (email)"
  value       = data.mongodbatlas_atlas_user.this.username
}

output "email_address" {
  description = "User's email address"
  value       = data.mongodbatlas_atlas_user.this.email_address
}

output "first_name" {
  description = "User's first name"
  value       = data.mongodbatlas_atlas_user.this.first_name
}

output "last_name" {
  description = "User's last name"
  value       = data.mongodbatlas_atlas_user.this.last_name
}

output "org_roles" {
  description = "User's organization roles (filtered by org_id)"
  value       = local.user_org_roles
}

output "project_roles" {
  description = "User's project roles (filtered by project_id, if provided)"
  value       = local.user_project_roles
}
