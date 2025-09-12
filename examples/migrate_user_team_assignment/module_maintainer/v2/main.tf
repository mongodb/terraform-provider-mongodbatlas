resource "mongodbatlas_team" "this" {
  org_id = var.org_id
  name   = var.team_name
}

resource "mongodbatlas_cloud_user_team_assignment" "this" {
  for_each = var.user_ids

  org_id  = var.org_id
  team_id = mongodbatlas_team.this.team_id
  user_id = each.value
}
