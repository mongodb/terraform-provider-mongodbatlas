resource "mongodbatlas_team_project_assignment" "example" {
  project_id = var.project_id
  team_id    = var.team_id
  role_names = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]
}

data "mongodbatlas_team_project_assignment" "example_username" {
  project_id = mongodbatlas_team_project_assignment.example.project_id
  team_id    = mongodbatlas_team_project_assignment.example.team_id
}
