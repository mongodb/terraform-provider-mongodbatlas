resource "mongodbatlas_cloud_user_project_assignment" "this" {
  project_id = var.project_id
  username   = var.username
  roles      = var.roles
}

# Removal of deprecated resource from state (keeps invitation in Atlas)
removed {
  from = mongodbatlas_project_invitation.this
  lifecycle {
    destroy = false
  }
}
