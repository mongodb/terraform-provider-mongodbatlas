############################################################
# v3: Final state - only new data sources
############################################################

# Single user reads using cloud_user_org_assignment
data "mongodbatlas_cloud_user_org_assignment" "user_by_id" {
  org_id  = var.org_id
  user_id = var.user_id
}

data "mongodbatlas_cloud_user_org_assignment" "user_by_username" {
  org_id   = var.org_id
  username = var.username
}

# User lists using organization/project/team data sources
data "mongodbatlas_organization" "org" {
  org_id = var.org_id
}

data "mongodbatlas_project" "proj" {
  project_id = var.project_id
}

data "mongodbatlas_team" "team" {
  team_id = var.team_id
  org_id  = var.org_id
}

# Clean, simplified local values using new data sources
locals {
  # Single user examples (simplified)
  user_email_by_id       = data.mongodbatlas_cloud_user_org_assignment.user_by_id.username
  user_email_by_username = data.mongodbatlas_cloud_user_org_assignment.user_by_username.username

  # User list examples (simplified)
  org_user_emails     = data.mongodbatlas_organization.org.users[*].username
  project_user_emails = data.mongodbatlas_project.proj.users[*].username
  team_user_emails    = data.mongodbatlas_team.team.users[*].username

  # Role examples (much cleaner than v1)
  user_org_roles = data.mongodbatlas_cloud_user_org_assignment.user_by_id.roles.org_roles

  # Find project role assignments that match the project_id
  matching_project_roles = [
    for pra in data.mongodbatlas_cloud_user_org_assignment.user_by_id.roles.project_role_assignments :
    pra.project_roles if pra.project_id == var.project_id
  ]
  # Use the first match if available, otherwise empty list
  user_project_roles = length(local.matching_project_roles) > 0 ? local.matching_project_roles[0] : []

  # All project role assignments
  user_all_project_roles = {
    for pra in data.mongodbatlas_cloud_user_org_assignment.user_by_id.roles.project_role_assignments :
    pra.project_id => pra.project_roles
  }
}
