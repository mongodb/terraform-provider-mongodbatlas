############################################################
# v2: Module using new mongodbatlas_cloud_user_org_assignment
############################################################

# NEW: Uses cloud_user_org_assignment with username lookup
data "mongodbatlas_cloud_user_org_assignment" "this" {
  org_id   = var.org_id
  username = var.username
}

# Role access is now simplified - no complex filtering needed
locals {
  # Project roles require finding the matching project in project_role_assignments
  matching_project_roles = var.project_id != null ? [
    for pra in data.mongodbatlas_cloud_user_org_assignment.this.roles.project_role_assignments :
    pra.project_roles if pra.project_id == var.project_id
  ] : []

  user_project_roles = length(local.matching_project_roles) > 0 ? local.matching_project_roles[0] : []
}
