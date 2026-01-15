# Module: creates a project with team assignments using new resource
# Note: keep ignore_changes until the provider removes the teams attribute

resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.org_id
  lifecycle {
    ignore_changes = [teams]
  }
}

resource "mongodbatlas_team_project_assignment" "this" {
  for_each = var.team_map

  project_id = mongodbatlas_project.this.id
  team_id    = each.key
  role_names = each.value
}
