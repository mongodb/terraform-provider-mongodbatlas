provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# New module usage
module "org_membership" {
  source   = "../../module_maintainer/v2"
  org_id   = var.org_id
  username = var.username
  roles    = var.roles
  team_ids = var.team_ids
}

# Team assignments (must be imported)
import {
  for_each = toset(var.team_ids)

  to = module.org_membership.mongodbatlas_cloud_user_team_assignment.team[each.key]
  id = "${var.org_id}/${each.key}/${var.username}" # or user_id
}

