############################################################
# v3: Final state - only new resource
############################################################

# Project user assignment using the new resource
resource "mongodbatlas_cloud_user_project_assignment" "user_assignment" {
  project_id = var.project_id
  username   = var.username
  roles      = var.roles
}

# Example of additional functionality available with new resource
data "mongodbatlas_cloud_user_project_assignment" "user_lookup" {
  project_id = var.project_id
  username   = mongodbatlas_cloud_user_project_assignment.user_assignment.username
}

# Clean, simplified local values
locals {
  # Basic assignment info
  assigned_user  = mongodbatlas_cloud_user_project_assignment.user_assignment.username
  assigned_roles = mongodbatlas_cloud_user_project_assignment.user_assignment.roles
  user_id        = mongodbatlas_cloud_user_project_assignment.user_assignment.user_id

  # Enhanced information from data source
  user_details = {
    username = data.mongodbatlas_cloud_user_project_assignment.user_lookup.username
    user_id  = data.mongodbatlas_cloud_user_project_assignment.user_lookup.user_id
    roles    = data.mongodbatlas_cloud_user_project_assignment.user_lookup.roles
  }

  # Assignment summary
  assignment_summary = {
    project_id = var.project_id
    user       = local.assigned_user
    roles      = local.assigned_roles
    user_id    = local.user_id
  }
}
