resource "mongodbatlas_team" "this" {
  org_id = var.org_id
  name   = var.team_name
}

resource "mongodbatlas_cloud_user_team_assignment" "this" {
  for_each = var.team_assigments

  org_id  = each.value.org_id
  team_id = each.value.team_id
  user_id = each.value.user_id
}
