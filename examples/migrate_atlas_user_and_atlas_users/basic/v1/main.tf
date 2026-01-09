############################################################
# v1: Original configuration using deprecated data sources
############################################################

# Single user read using deprecated data source
data "mongodbatlas_atlas_user" "single_user_by_id" {
  user_id = var.user_id
}

data "mongodbatlas_atlas_user" "single_user_by_username" {
  username = var.username
}

# User lists using deprecated data source
data "mongodbatlas_atlas_users" "org_users" {
  org_id = var.org_id
}

data "mongodbatlas_atlas_users" "project_users" {
  project_id = var.project_id
}

data "mongodbatlas_atlas_users" "team_users" {
  team_id = var.team_id
  org_id  = var.org_id
}

# Example usage of deprecated data sources
locals {
  # Single user examples
  user_email_by_id       = data.mongodbatlas_atlas_user.single_user_by_id.email_address
  user_email_by_username = data.mongodbatlas_atlas_user.single_user_by_username.email_address

  # User list examples
  org_user_emails     = data.mongodbatlas_atlas_users.org_users.results[*].email_address
  project_user_emails = data.mongodbatlas_atlas_users.project_users.results[*].email_address
  team_user_emails    = data.mongodbatlas_atlas_users.team_users.results[*].email_address

  # Role filtering examples (complex expressions)
  user_org_roles = [
    for r in data.mongodbatlas_atlas_user.single_user_by_id.roles : r.role_name
    if r.org_id == var.org_id
  ]

  user_project_roles = [
    for r in data.mongodbatlas_atlas_user.single_user_by_id.roles : r.role_name
    if r.group_id == var.project_id
  ]
}
