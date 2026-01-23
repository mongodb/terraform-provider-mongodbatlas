resource "mongodbatlas_project_invitation" "this" {
  project_id = var.project_id
  username   = var.username
  roles      = var.roles
}
