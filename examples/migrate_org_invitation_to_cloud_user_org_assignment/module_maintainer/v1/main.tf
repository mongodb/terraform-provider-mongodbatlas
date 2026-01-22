resource "mongodbatlas_org_invitation" "this" {
  org_id    = var.org_id
  username  = var.username
  roles     = var.roles
  teams_ids = var.team_ids
}

