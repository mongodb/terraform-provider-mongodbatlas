# Module: creates a project with team assignments using deprecated teams block
provider "mongodbatlas" {
  base_url = "https://cloud-dev.mongodb.com/"
}
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.org_id

  dynamic "teams" {
    for_each = var.team_map
    content {
      team_id    = teams.key
      role_names = teams.value
    }
  }
}
