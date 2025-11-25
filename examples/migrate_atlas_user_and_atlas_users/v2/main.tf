############################################################
# v2: Migration phase - both old and new data sources
############################################################

# OLD: Single user reads (keep temporarily for comparison)
data "mongodbatlas_atlas_user" "single_user_by_id" {
  user_id = var.user_id
}

data "mongodbatlas_atlas_user" "single_user_by_username" {
  username = var.username
}

# NEW: Single user reads using cloud_user_org_assignment
data "mongodbatlas_cloud_user_org_assignment" "user_by_id" {
  org_id  = var.org_id
  user_id = var.user_id
}

data "mongodbatlas_cloud_user_org_assignment" "user_by_username" {
  org_id   = var.org_id
  username = var.username
}

# OLD: User lists (keep temporarily for comparison)
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

# NEW: User lists using organization/project/team data sources
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

# Migration examples showing attribute mapping
locals {
  # Single user attribute mapping examples

  # Email address mapping
  old_user_email_by_id = data.mongodbatlas_atlas_user.single_user_by_id.email_address
  new_user_email_by_id = data.mongodbatlas_cloud_user_org_assignment.user_by_id.username

  # Organization roles mapping
  old_user_org_roles = [
    for r in data.mongodbatlas_atlas_user.single_user_by_id.roles : r.role_name
    if r.org_id == var.org_id
  ]
  new_user_org_roles = data.mongodbatlas_cloud_user_org_assignment.user_by_id.roles.org_roles

  # Project roles mapping (more complex for old, simpler for new)
  old_user_project_roles = [
    for r in data.mongodbatlas_atlas_user.single_user_by_id.roles : r.role_name
    if r.group_id == var.project_id
  ]
  # Find project role assignments that match the project_id
  matching_project_roles = [
    for pra in data.mongodbatlas_cloud_user_org_assignment.user_by_id.roles.project_role_assignments :
    pra.project_roles if pra.project_id == var.project_id
  ]

  # Use the first match if available, otherwise empty list
  new_user_project_roles = length(local.matching_project_roles) > 0 ? local.matching_project_roles[0] : []

  # User list attribute mapping examples

  # Organization users
  old_org_user_emails = data.mongodbatlas_atlas_users.org_users.results[*].email_address
  new_org_user_emails = data.mongodbatlas_organization.org.users[*].username

  # Project users  
  old_project_user_emails = data.mongodbatlas_atlas_users.project_users.results[*].email_address
  new_project_user_emails = data.mongodbatlas_project.proj.users[*].username

  # Team users
  old_team_user_emails = data.mongodbatlas_atlas_users.team_users.results[*].email_address
  new_team_user_emails = data.mongodbatlas_team.team.users[*].username

  # Validation: Compare old vs new results
  email_mapping_matches   = local.old_user_email_by_id == local.new_user_email_by_id
  org_users_count_matches = length(local.old_org_user_emails) == length(local.new_org_user_emails)
}
