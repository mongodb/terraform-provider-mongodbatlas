resource "mongodbatlas_team_project_assignment" "this" {
  project_id = var.project_id
  team_id    = var.team_id
  role_names = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]
}

data "mongodbatlas_team_project_assignment" "this" {
  project_id = mongodbatlas_team_project_assignment.this.project_id
  team_id    = mongodbatlas_team_project_assignment.this.team_id
}
