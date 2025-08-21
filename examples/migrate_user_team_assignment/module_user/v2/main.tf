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
  team_assigments = {
    for user in data.mongodbatlas_team.this.users :
    user.id => {
      org_id  = var.org_id
      team_id = data.mongodbatlas_team.this.team_id
      user_id = user.id
    }
  }
}

module "user_team_assignment" {
  source          = "../../module_maintainer/v2"
  org_id          = var.org_id
  team_name       = var.team_name
  team_assigments = local.team_assigments
}

import {
  for_each = local.team_assigments

  to = module.user_team_assignment.mongodbatlas_cloud_user_team_assignment.this[each.key]
  id = "${var.org_id}/${data.mongodbatlas_team.this.team_id}/${each.value.user_id}"
}
