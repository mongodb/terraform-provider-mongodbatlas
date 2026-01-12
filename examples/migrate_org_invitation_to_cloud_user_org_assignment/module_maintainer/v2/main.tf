resource "mongodbatlas_cloud_user_org_assignment" "this" {
  org_id   = var.org_id
  username = var.username
  roles    = { org_roles = var.roles }
}

resource "mongodbatlas_cloud_user_team_assignment" "team" {
  for_each = var.team_ids

  org_id  = var.org_id
  team_id = each.key
  user_id = mongodbatlas_cloud_user_org_assignment.this.user_id
}

moved {
  from = mongodbatlas_org_invitation.this
  to   = mongodbatlas_cloud_user_org_assignment.this
}

