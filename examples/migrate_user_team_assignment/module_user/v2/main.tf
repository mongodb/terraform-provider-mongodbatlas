provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

# New module usage
data "mongodbatlas_team" "this" {
  org_id = var.org_id
  name   = var.team_name
}

locals {
  user_ids = toset([
    for user in data.mongodbatlas_team.this.users : user.id
  ])
}

module "user_team_assignment" {
  source          = "../../module_maintainer/v2"
  org_id          = var.org_id
  team_name       = var.team_name
  user_ids = local.user_ids
}

import {
  for_each = local.user_ids

  to = module.user_team_assignment.mongodbatlas_cloud_user_team_assignment.this[each.key]
  id = "${var.org_id}/${data.mongodbatlas_team.this.team_id}/${each.value.user_id}"
}
