############################################################
# v2: Migration phase - re-create with new resource and remove old
############################################################

# NEW: Project assignment using the new resource
# This re-creates the pending invitation using the new resource
resource "mongodbatlas_cloud_user_project_assignment" "user_assignment" {
  project_id = var.project_id
  username   = var.username
  roles      = var.roles
}

# REMOVE: Clean removal of deprecated resource from state
removed {
  from = mongodbatlas_project_invitation.pending_user

  lifecycle {
    destroy = false
  }
}

# Migration validation
locals {
  # Verify the new resource works correctly
  new_assignment_user  = mongodbatlas_cloud_user_project_assignment.user_assignment.username
  new_assignment_roles = mongodbatlas_cloud_user_project_assignment.user_assignment.roles

  # Basic validation
  username_matches = var.username == local.new_assignment_user
  roles_match      = toset(var.roles) == toset(local.new_assignment_roles)

  migration_successful = local.username_matches && local.roles_match
}
