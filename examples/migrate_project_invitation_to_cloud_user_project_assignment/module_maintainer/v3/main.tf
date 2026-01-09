resource "mongodbatlas_cloud_user_project_assignment" "this" {
  project_id = var.project_id
  username   = var.username
  roles      = var.roles
}
