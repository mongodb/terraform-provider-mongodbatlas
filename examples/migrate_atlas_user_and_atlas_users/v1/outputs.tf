# Single user outputs
output "user_email_by_id" {
  description = "User email retrieved by user ID"
  value       = local.user_email_by_id
}

output "user_email_by_username" {
  description = "User email retrieved by username"
  value       = local.user_email_by_username
}

output "user_org_roles" {
  description = "User's organization roles (filtered from consolidated roles)"
  value       = local.user_org_roles
}

output "user_project_roles" {
  description = "User's project roles (filtered from consolidated roles)"
  value       = local.user_project_roles
}

# User list outputs
output "org_user_emails" {
  description = "All organization user emails"
  value       = local.org_user_emails
}

output "project_user_emails" {
  description = "All project user emails"
  value       = local.project_user_emails
}

output "team_user_emails" {
  description = "All team user emails"
  value       = local.team_user_emails
}

# Count outputs
output "org_user_count" {
  description = "Number of organization users"
  value       = length(data.mongodbatlas_atlas_users.org_users.results)
}

output "project_user_count" {
  description = "Number of project users"
  value       = length(data.mongodbatlas_atlas_users.project_users.results)
}

output "team_user_count" {
  description = "Number of team users"
  value       = length(data.mongodbatlas_atlas_users.team_users.results)
}
