############################################################
# v1: Module using deprecated mongodbatlas_atlas_user
############################################################

# Using username (email) for lookup - more practical than user_id
data "mongodbatlas_atlas_user" "this" {
  username = var.username
}

# Complex role filtering required with deprecated data source
locals {
  # Filter roles to get only org-level roles for the specified org
  user_org_roles = [
    for r in data.mongodbatlas_atlas_user.this.roles : r.role_name
    if r.org_id == var.org_id
  ]

  # Filter roles to get project-level roles (requires project_id)
  user_project_roles = var.project_id != null ? [
    for r in data.mongodbatlas_atlas_user.this.roles : r.role_name
    if r.group_id == var.project_id
  ] : []
}
