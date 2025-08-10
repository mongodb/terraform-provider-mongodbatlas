resource "mongodbatlas_cloud_user_project_assignment" "example" {
  project_id = var.project_id
  username   = var.user_email
  roles      = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]
}

data "mongodbatlas_cloud_user_project_assignment" "example_username" {
  project_id = var.project_id
  username   = var.user_email
}

data "mongodbatlas_cloud_user_project_assignment" "example_user_id" {
  project_id = var.project_id
  user_id    = var.user_id
}
