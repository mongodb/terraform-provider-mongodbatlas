provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# New module usage (team_project_assignment resource)
module "project" {
  source       = "../../module_maintainer/v2"
  org_id       = var.org_id
  project_name = var.project_name
  team_map     = var.team_map
}

# Import existing team-project assignments (must be at root level)
import {
  for_each = var.team_map

  to = module.project.mongodbatlas_team_project_assignment.this[each.key]
  id = "${module.project.project_id}/${each.key}"
}
