provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# Final module usage (no import blocks needed after migration)
module "project" {
  source       = "../../module_maintainer/v3"
  org_id       = var.org_id
  project_name = var.project_name
  team_map     = var.team_map
}
