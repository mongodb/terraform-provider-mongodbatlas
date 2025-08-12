resource "mongodbatlas_cloud_user_team_assignment" "example" {
  org_id  = var.org_id
  team_id = var.team_id
  user_id = var.user_id
}

data "mongodbatlas_cloud_user_team_assignment" "example_user_id" {
  org_id  = var.org_id
  team_id = var.team_id
  user_id = var.user_id
  depends_on = [mongodbatlas_cloud_user_team_assignment.example]
}

data "mongodbatlas_cloud_user_team_assignment" "example_username" {
  org_id   = var.org_id
  team_id  = var.team_id
  username = var.user_email
  depends_on = [mongodbatlas_cloud_user_team_assignment.example]
}
