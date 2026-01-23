############################################################
# v3: Clean module using only new data source
############################################################

data "mongodbatlas_cloud_user_org_assignment" "this" {
  org_id   = var.org_id
  username = var.username
}

locals {
  # Project roles filtering (if project_id provided)
  matching_project_roles = var.project_id != null ? [
    for pra in data.mongodbatlas_cloud_user_org_assignment.this.roles.project_role_assignments :
    pra.project_roles if pra.project_id == var.project_id
  ] : []

  user_project_roles = length(local.matching_project_roles) > 0 ? local.matching_project_roles[0] : []
}
