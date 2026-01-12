############################################################
# v1: Original configuration using deprecated resource
############################################################

# Pending project invitation using deprecated resource
resource "mongodbatlas_project_invitation" "pending_user" {
  project_id = var.project_id
  username   = var.username
  roles      = var.roles
}

# Example usage of the invitation
locals {
  invitation_id  = mongodbatlas_project_invitation.pending_user.invitation_id
  invited_user   = mongodbatlas_project_invitation.pending_user.username
  assigned_roles = mongodbatlas_project_invitation.pending_user.roles

  # This shows how the deprecated resource was typically used
  invitation_details = {
    id       = local.invitation_id
    username = local.invited_user
    roles    = local.assigned_roles
    project  = var.project_id
  }
}
