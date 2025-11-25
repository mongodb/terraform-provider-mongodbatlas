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
  description = "User's organization roles (structured)"
  value       = local.user_org_roles
}

output "user_project_roles" {
  description = "User's roles for specific project"
  value       = local.user_project_roles
}

output "user_all_project_roles" {
  description = "User's roles across all projects"
  value       = local.user_all_project_roles
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
  value       = length(data.mongodbatlas_organization.org.users)
}

output "project_user_count" {
  description = "Number of project users"
  value       = length(data.mongodbatlas_project.proj.users)
}

output "team_user_count" {
  description = "Number of team users"
  value       = length(data.mongodbatlas_team.team.users)
}

# User details from different scopes
output "org_users_with_roles" {
  description = "Organization users with their roles"
  value = [
    for user in data.mongodbatlas_organization.org.users : {
      username = user.username
      user_id  = user.id
      # Although the API defines roles as an object, in the organization and team data sources roles are represented as a list. This is due to SDK v2 only supporting blocks for nested elements.
      org_roles           = user.roles[0].org_roles
      project_assignments = user.roles[0].project_role_assignments
    }
  ]
}

output "project_users_with_roles" {
  description = "Project users with their roles"
  value = [
    for user in data.mongodbatlas_project.proj.users : {
      username = user.username
      user_id  = user.id
      roles    = user.roles
    }
  ]
}
